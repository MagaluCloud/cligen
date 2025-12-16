package generate_config

import (
	"fmt"
	"go/ast"
	"strings"
)

var primitiveTypes = []string{
	"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64", "rune", "string", "uint",
	"uint8", "uint16", "uint32", "uint64", "uintptr", "string", "time", "context",
}
var serviceSuffixes = []string{"Service", "API", "Client"}
var serviceKeywords = []string{"service", "api", "client"}

func getTypeString(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, isPrimitiveType(t.Name)
	case *ast.StarExpr:
		subType, isPrimitive := getTypeString(t.X)
		return "*" + subType, isPrimitive
	case *ast.ArrayType:
		subType, isPrimitive := getTypeString(t.Elt)
		return "[]" + subType, isPrimitive
	case *ast.SelectorExpr:
		subType, isPrimitive := getTypeString(t.X)
		return subType + "." + t.Sel.Name, isPrimitive
	case *ast.InterfaceType:
		return "interface{}", true
	default:
		return fmt.Sprintf("%T", expr), false
	}
}

func resolveUnderlyingType(ident *ast.Ident) string {
	if ident.Obj == nil {
		return ""
	}

	typeSpec, ok := ident.Obj.Decl.(*ast.TypeSpec)
	if !ok {
		return ""
	}
	underlyingType, _ := getTypeString(typeSpec.Type)
	return underlyingType
}
func getTypeStringWithPackage(expr ast.Expr, packageName string) (string, string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		if isPrimitiveType(t.Name) {
			return t.Name, "", true
		}

		underlyingType := resolveUnderlyingType(t)
		if underlyingType != "" && isPrimitiveType(underlyingType) {
			return underlyingType, t.Name, true
		}

		return packageName + "Sdk." + t.Name, "", false
	case *ast.StarExpr:
		subType, aliasType, isPrimitive := getTypeStringWithPackage(t.X, packageName)
		return "*" + subType, aliasType, isPrimitive
	case *ast.ArrayType:
		elementType, aliasType, isPrimitive := getTypeStringWithPackage(t.Elt, packageName)
		if isPrimitive {
			return "[]" + elementType, aliasType, isPrimitive
		}
		if strings.Contains(elementType, ".") {
			return "[]" + elementType, aliasType, isPrimitive
		}
		return "[]" + packageName + "." + elementType, aliasType, isPrimitive
	case *ast.MapType:
		keyType, keyAliasType, keyPrimitive := getTypeStringWithPackage(t.Key, packageName)
		valueType, valueAliasType, valuePrimitive := getTypeStringWithPackage(t.Value, packageName)
		isPrimitive := keyPrimitive && valuePrimitive
		return fmt.Sprintf("map[%s]%s", keyType, valueType), keyAliasType + valueAliasType, isPrimitive
	case *ast.ChanType:
		elementType, elementAliasType, elementPrimitive := getTypeStringWithPackage(t.Value, packageName)
		var chanType string
		switch t.Dir {
		case ast.SEND:
			chanType = "chan<-"
		case ast.RECV:
			chanType = "<-chan"
		default:
			chanType = "chan"
		}
		return chanType + " " + elementType, elementAliasType, elementPrimitive
	case *ast.FuncType:
		return "func()", "", false
	case *ast.SelectorExpr:
		elementType, elementAliasType, isPrimitive := getTypeStringWithPackage(t.X, packageName)
		return elementType + "." + t.Sel.Name, elementAliasType, isPrimitive
	case *ast.InterfaceType:
		return "interface{}", "", true
	default:
		return fmt.Sprintf("%T", expr), "", false
	}
}

func isPrimitiveType(typeName string) bool {
	if strings.Contains(typeName, ".") {
		typeName = strings.Split(typeName, ".")[len(strings.Split(typeName, "."))-1]
	}
	return contains(primitiveTypes, typeName)
}

// contains verifica se um slice contém um elemento
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isSubServiceType verifica se um tipo de retorno representa um subserviço
func isSubServiceType(returnType string) bool {
	// Remover ponteiros e arrays para análise
	baseType := strings.TrimPrefix(returnType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	// Verificar se o tipo termina com sufixos comuns de serviço
	for _, suffix := range serviceSuffixes {
		if strings.HasSuffix(baseType, suffix) {
			return true
		}
	}

	// Verificar se contém palavras-chave de serviço
	lowerType := strings.ToLower(baseType)
	for _, keyword := range serviceKeywords {
		if strings.Contains(lowerType, keyword) {
			return true
		}
	}

	// Verificar padrões específicos como "networkBackendTargetService"
	if strings.Contains(lowerType, "service") && len(baseType) > 10 {
		return true
	}

	return false
}

// extractSubServiceName extrai o nome do subserviço a partir do tipo de retorno
func extractSubServiceName(returnType string, methodName string) string {
	// Remover ponteiros e arrays
	baseType := strings.TrimPrefix(returnType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	// Se o tipo contém um ponto (pacote.tipo), extrair apenas o nome do tipo
	if strings.Contains(baseType, ".") {
		parts := strings.Split(baseType, ".")
		baseType = parts[len(parts)-1]
	}

	// Remover sufixos comuns de serviço para obter o nome base
	for _, suffix := range serviceSuffixes {
		if strings.HasSuffix(baseType, suffix) {
			baseType = strings.TrimSuffix(baseType, suffix)
			break
		}
	}

	// Se o nome base estiver vazio, usar o nome do método
	if baseType == "" {
		baseType = methodName
	}

	// Converter para PascalCase se necessário
	if len(baseType) > 0 {
		// Se já está em PascalCase, manter como está
		if baseType[0] >= 'A' && baseType[0] <= 'Z' {
			return baseType
		}
		// Converter para PascalCase
		baseType = strings.ToUpper(baseType[:1]) + baseType[1:]
	}

	return baseType
}

// generateParamName gera um nome para um parâmetro baseado no tipo
func generateParamName(paramType string) string {
	// Converter o tipo para um nome de variável em camelCase
	switch paramType {
	case "context.Context":
		return "ctx"
	case "string":
		return "str"
	case "int":
		return "num"
	case "bool":
		return "flag"
	case "error":
		return "err"
	case "interface{}":
		return "data"
	default:
		// Para tipos complexos, tentar extrair o nome base
		if strings.Contains(paramType, ".") {
			parts := strings.Split(paramType, ".")
			return strings.ToLower(parts[len(parts)-1])
		}
		// Remover ponteiros e arrays para gerar nome base
		baseType := strings.TrimPrefix(paramType, "*")
		baseType = strings.TrimPrefix(baseType, "[]")
		return strings.ToLower(baseType)
	}
}

// generateReturnName gera um nome para um retorno baseado no tipo
func generateReturnName(returnType string) string {
	// Para retornos, usar nomes mais descritivos
	switch returnType {
	case "error":
		return "err"
	case "bool":
		return "success"
	case "string":
		return "result"
	case "int":
		return "count"
	default:
		// Para tipos complexos, tentar extrair o nome base
		if strings.Contains(returnType, ".") {
			parts := strings.Split(returnType, ".")
			return strings.ToLower(parts[len(parts)-1])
		}
		// Remover ponteiros e arrays para gerar nome base
		baseType := strings.TrimPrefix(returnType, "*")
		baseType = strings.TrimPrefix(baseType, "[]")
		return strings.ToLower(baseType)
	}
}

// generateFieldName gera um nome para um campo anônimo baseado no tipo
func generateFieldName(fieldType string) string {
	// Remover ponteiros e arrays para gerar nome base
	baseType := strings.TrimPrefix(fieldType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	// Se contém ponto (pacote.tipo), extrair apenas o nome do tipo
	if strings.Contains(baseType, ".") {
		parts := strings.Split(baseType, ".")
		baseType = parts[len(parts)-1]
	}

	return strings.ToLower(baseType)
}

// extractTypeName extrai o nome do tipo de um ast.Expr
func extractTypeName(expr ast.Expr, packageName string) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return extractTypeName(t.X, packageName)
	case *ast.ArrayType:
		return extractTypeName(t.Elt, packageName)
	case *ast.SelectorExpr:
		// Para tipos de outros pacotes, retornar apenas o nome do tipo
		return t.Sel.Name
	default:
		return ""
	}
}
