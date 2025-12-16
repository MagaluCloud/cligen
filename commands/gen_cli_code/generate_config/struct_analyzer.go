package generate_config

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/magaluCloud/cligen/config"
)

type fileCache struct {
	mu    sync.RWMutex
	cache map[string]*ast.File
}

type structCache struct {
	mu    sync.RWMutex
	cache map[string]map[string]config.Parameter
}

var (
	globalFileCache   = &fileCache{cache: make(map[string]*ast.File)}
	globalStructCache = &structCache{cache: make(map[string]map[string]config.Parameter)}
)

func extractJSONTag(tagValue string) []string {
	tagValue = strings.Trim(tagValue, "`\"")

	if strings.Contains(tagValue, "json:") {
		parts := strings.Split(tagValue, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "json:") {
				jsonValue := strings.TrimPrefix(part, "json:")
				jsonValue = strings.Trim(jsonValue, "\"")

				if jsonValue == "-" {
					return []string{}
				}

				jsonValues := strings.Split(jsonValue, ",")
				return jsonValues

			}
		}
	}

	return []string{}
}

// analyzeStructType analisa um tipo para verificar se é uma struct e extrai seus campos
func analyzeStructType(expr ast.Expr, packageName string, sdkDir string) map[string]config.Parameter {
	// Extrair o nome do tipo base
	typeName := extractTypeName(expr, packageName)
	if typeName == "" {
		return nil
	}

	structFields := extractTypeFieldsFromIdent(expr, packageName)
	if structFields != nil {
		return structFields
	}

	structFields = findStructDefinition(typeName, sdkDir, packageName)
	if structFields != nil {
	}

	return structFields
}

func extractTypeFieldsFromIdent(expr ast.Expr, packageName string) map[string]config.Parameter {
	// Tratar ponteiros desembrulhando recursivamente
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		return extractTypeFieldsFromIdent(starExpr.X, packageName)
	}

	if ident, ok := expr.(*ast.Ident); ok {
		if ident.Obj != nil {
			if typeDecl, ok := ident.Obj.Decl.(*ast.TypeSpec); ok {
				if structType, ok := typeDecl.Type.(*ast.StructType); ok {
					return extractStructFields(structType, packageName, "")
				}
			}
		}
	}

	return nil
}

// findStructDefinition procura pela definição de uma struct no diretório do SDK com cache
func findStructDefinition(typeName string, sdkDir string, packageName string) map[string]config.Parameter {
	// Remover o prefixo do pacote se presente
	cleanTypeName := typeName
	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		cleanTypeName = parts[len(parts)-1]
	}

	// Verificar cache de structs
	cacheKey := sdkDir + ":" + cleanTypeName
	globalStructCache.mu.RLock()
	if cached, exists := globalStructCache.cache[cacheKey]; exists {
		globalStructCache.mu.RUnlock()
		return cached
	}
	globalStructCache.mu.RUnlock()

	// Cache miss - procurar struct
	files, err := os.ReadDir(sdkDir)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(sdkDir, file.Name())
		structFields := analyzeFileForStruct(filePath, cleanTypeName, packageName)
		if structFields != nil {
			// Armazenar no cache
			globalStructCache.mu.Lock()
			globalStructCache.cache[cacheKey] = structFields
			globalStructCache.mu.Unlock()
			return structFields
		}
	}

	// Armazenar nil no cache para evitar buscas futuras
	globalStructCache.mu.Lock()
	globalStructCache.cache[cacheKey] = nil
	globalStructCache.mu.Unlock()

	return nil
}

func extractStructFields(structType *ast.StructType, packageName string, sdkDir string) map[string]config.Parameter {
	fields := make(map[string]config.Parameter)

	if structType.Fields == nil {
		return fields
	}

	for _, field := range structType.Fields.List {
		// Extrair comentários do campo
		var description string
		if field.Doc != nil {
			description = field.Doc.Text()
		} else if field.Comment != nil {
			description = field.Comment.Text()
		}

		// Extrair tipo do campo
		fieldType, isPrimitive, aliasType, structFields := analyzeParameterType(field.Type, packageName, sdkDir)

		isOptional := false
		if field.Tag != nil {
			tagValue := extractJSONTag(field.Tag.Value)
			if slices.Contains(tagValue, "omitempty") {
				isOptional = true
			}
		}
		if !isPrimitive {
			isOptional = true
		}

		isPointer := false
		if strings.HasPrefix(fieldType, "*") {
			isPointer = true
			isOptional = true
		}

		isArray := false
		if strings.HasPrefix(fieldType, "[]") || strings.HasPrefix(fieldType, "*[]") {
			isArray = true
		}

		// Se o campo tem nome, usar o nome, senão gerar um nome baseado no tipo
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				// Verificar se há tags JSON
				var jsonName string
				if field.Tag != nil {
					jsonNames := extractJSONTag(field.Tag.Value)
					if len(jsonNames) > 0 {
						jsonName = jsonNames[0]
					}
				}
				if jsonName == "" {
					jsonName = name.Name
				}

				fields[jsonName] = config.Parameter{
					Name:        name.Name,
					Type:        fieldType,
					Description: description,
					IsPrimitive: isPrimitive,
					Struct:      structFields,
					IsPointer:   isPointer,
					IsArray:     isArray,
					IsOptional:  isOptional,
					AliasType:   aliasType,
				}
			}
		} else {

			for _, field := range structFields {
				fields[field.Name] = field
			}

		}
	}

	return fields
}
func analyzeParameterType(expr ast.Expr, packageName string, sdkDir string) (string, bool, string, map[string]config.Parameter) {
	if structType, ok := expr.(*ast.StructType); ok {
		structFields := extractStructFields(structType, packageName, sdkDir)
		return "struct{}", false, "", structFields
	}

	if _, ok := expr.(*ast.InterfaceType); ok {
		return "interface{}", true, "", nil
	}

	paramType, aliasType, isPrimitive := getTypeStringWithPackage(expr, packageName)

	if aliasType != "" {
		aliasType = packageName + "Sdk." + aliasType
	}
	if isPrimitive {
		return paramType, isPrimitive, aliasType, nil
	}

	structFields := analyzeStructType(expr, packageName, sdkDir)
	return paramType, isPrimitive, aliasType, structFields
}
func analyzeFileForStruct(filePath string, typeName string, packageName string) map[string]config.Parameter {
	// Verificar cache de arquivo
	globalFileCache.mu.RLock()
	file, exists := globalFileCache.cache[filePath]
	globalFileCache.mu.RUnlock()

	if !exists {
		// Cache miss - fazer parsing
		globalFileCache.mu.Lock()
		// Double-check após adquirir lock exclusivo
		if cached, exists := globalFileCache.cache[filePath]; exists {
			globalFileCache.mu.Unlock()
			file = cached
		} else {
			fset := token.NewFileSet()
			parsedFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
			if err != nil {
				globalFileCache.mu.Unlock()
				return nil
			}
			globalFileCache.cache[filePath] = parsedFile
			file = parsedFile
			globalFileCache.mu.Unlock()
		}
	}

	var structFields map[string]config.Parameter

	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeDecl.Type.(*ast.StructType); ok {
				if typeDecl.Name.Name == typeName {
					structFields = extractStructFields(structType, packageName, filepath.Dir(filePath))
					return false
				}
			}
		}
		return true
	})

	return structFields
}
