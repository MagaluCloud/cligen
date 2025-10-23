package sdk_structure

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// analyzeParameterType analisa um tipo de parâmetro e retorna informações detalhadas incluindo campos de struct
func analyzeParameterType(expr ast.Expr, packageName string, sdkDir string) (string, bool, map[string]Parameter) {
	// Verificar se é uma struct inline (anônima)
	if structType, ok := expr.(*ast.StructType); ok {
		// fmt.Printf("   🔍 Struct inline detectada\n")
		structFields := extractStructFields(structType, packageName, sdkDir)
		return "struct{}", false, structFields
	}

	// Verificar se é uma interface inline (anônima)
	if _, ok := expr.(*ast.InterfaceType); ok {
		// fmt.Printf("   🔍 Interface inline detectada\n")
		// Para interfaces, podemos extrair métodos se necessário
		// Por enquanto, retornamos como interface{}
		return "interface{}", true, nil
	}

	paramType, isPrimitive := getTypeStringWithPackage(expr, packageName)

	// Se é primitivo, não precisa analisar struct
	if isPrimitive {
		return paramType, isPrimitive, nil
	}

	// Verificar se é um tipo próprio que pode ser uma struct
	structFields := analyzeStructType(expr, packageName, sdkDir)

	return paramType, isPrimitive, structFields
}

// analyzeStructType analisa um tipo para verificar se é uma struct e extrai seus campos
func analyzeStructType(expr ast.Expr, packageName string, sdkDir string) map[string]Parameter {
	// Extrair o nome do tipo base
	typeName := extractTypeName(expr, packageName)
	if typeName == "" {
		return nil
	}

	structFields := extractTypeFieldsFromIdent(expr, packageName)
	if structFields != nil {
		// fmt.Printf("   🔍 Struct encontrada: %s com %d campos\n", typeName, len(structFields))
		return structFields
	}

	// Procurar pela definição da struct no diretório do SDK
	structFields = findStructDefinition(typeName, sdkDir, packageName)
	if structFields != nil {
		// fmt.Printf("   🔍 Struct encontrada: %s com %d campos\n", typeName, len(structFields))
	}

	return structFields
}

func extractTypeFieldsFromIdent(expr ast.Expr, packageName string) map[string]Parameter {
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

// findStructDefinition procura pela definição de uma struct no diretório do SDK
func findStructDefinition(typeName string, sdkDir string, packageName string) map[string]Parameter {
	// Remover o prefixo do pacote se presente
	cleanTypeName := typeName
	if strings.Contains(typeName, ".") {
		parts := strings.Split(typeName, ".")
		cleanTypeName = parts[len(parts)-1]
	}

	// Procurar em todos os arquivos .go do diretório
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
			return structFields
		}
	}

	return nil
}

// analyzeFileForStruct analisa um arquivo procurando por uma definição de struct específica
func analyzeFileForStruct(filePath string, typeName string, packageName string) map[string]Parameter {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	var structFields map[string]Parameter

	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeDecl.Type.(*ast.StructType); ok {
				// Verificar se é a struct que estamos procurando
				if typeDecl.Name.Name == typeName {
					// fmt.Printf("   ✅ Struct encontrada: %s\n", typeName)
					structFields = extractStructFields(structType, packageName, filepath.Dir(filePath))
					return false // Parar a busca
				}
			}
		}
		return true
	})

	return structFields
}

// extractStructFields extrai os campos de uma struct
func extractStructFields(structType *ast.StructType, packageName string, sdkDir string) map[string]Parameter {
	fields := make(map[string]Parameter)

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

		isOptional := false
		if field.Tag != nil {
			tagValue := extractJSONTag(field.Tag.Value)
			if slices.Contains(tagValue, "omitempty") {
				isOptional = true
			}
		}

		// Extrair tipo do campo
		fieldType, isPrimitive, structFields := analyzeParameterType(field.Type, packageName, sdkDir)

		isPointer := false
		if strings.HasPrefix(fieldType, "*") {
			isPointer = true
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

				fields[jsonName] = Parameter{
					Name:        name.Name,
					Type:        fieldType,
					Description: description,
					IsPrimitive: isPrimitive,
					Struct:      structFields,
					IsPointer:   isPointer,
					IsArray:     isArray,
					IsOptional:  isOptional,
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

// extractJSONTag extrai o nome do campo da tag JSON
func extractJSONTag(tagValue string) []string {
	// Remover aspas
	tagValue = strings.Trim(tagValue, "`\"")

	// Procurar pela tag json
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
