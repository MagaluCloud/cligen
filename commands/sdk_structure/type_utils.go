package sdk_structure

import (
	"fmt"
	"go/ast"
	"strings"
)

// getTypeString converte um ast.Expr para string representando o tipo
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

// getTypeStringWithPackage converte um ast.Expr para string representando o tipo, incluindo o pacote quando necessário
func getTypeStringWithPackage(expr ast.Expr, packageName string) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		// Verificar se é um tipo primitivo
		if isPrimitiveType(t.Name) {
			return t.Name, true
		}
		// Se não for primitivo, adicionar o pacote
		return packageName + "Sdk." + t.Name, false
	case *ast.StarExpr:
		subType, isPrimitive := getTypeStringWithPackage(t.X, packageName)
		return "*" + subType, isPrimitive
	case *ast.ArrayType:
		// Para arrays, verificar se o tipo do elemento é primitivo
		elementType, isPrimitive := getTypeStringWithPackage(t.Elt, packageName)
		// Se o elemento é um tipo primitivo, não adicionar o packageName
		if isPrimitive {
			return "[]" + elementType, isPrimitive
		}
		// Se o elemento já tem o packageName, usar como está
		if strings.Contains(elementType, ".") {
			return "[]" + elementType, isPrimitive
		}
		// Caso contrário, adicionar o packageName
		return "[]" + packageName + "." + elementType, isPrimitive
	case *ast.MapType:
		// Para maps, analisar chave e valor
		keyType, keyPrimitive := getTypeStringWithPackage(t.Key, packageName)
		valueType, valuePrimitive := getTypeStringWithPackage(t.Value, packageName)
		// Map é considerado primitivo se ambos chave e valor são primitivos
		isPrimitive := keyPrimitive && valuePrimitive
		return fmt.Sprintf("map[%s]%s", keyType, valueType), isPrimitive
	case *ast.ChanType:
		// Para channels, analisar o tipo do elemento
		elementType, elementPrimitive := getTypeStringWithPackage(t.Value, packageName)
		// Channel é considerado primitivo se o elemento é primitivo
		var chanType string
		switch t.Dir {
		case ast.SEND:
			chanType = "chan<-"
		case ast.RECV:
			chanType = "<-chan"
		default:
			chanType = "chan"
		}
		return chanType + " " + elementType, elementPrimitive
	case *ast.FuncType:
		// Para function types, gerar uma representação simplificada
		// Function types são considerados não primitivos
		return "func()", false
	case *ast.SelectorExpr:
		// SelectorExpr já tem o pacote qualificado (ex: context.Context)
		elementType, isPrimitive := getTypeString(t.X)
		return elementType + "." + t.Sel.Name, isPrimitive
	case *ast.InterfaceType:
		return "interface{}", true
	default:
		return fmt.Sprintf("%T", expr), false
	}
}

// isPrimitiveType verifica se um tipo é primitivo do Go
func isPrimitiveType(typeName string) bool {
	fmt.Println(typeName)
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
