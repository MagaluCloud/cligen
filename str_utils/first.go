package strutils

import "strings"

func FirstLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func FirstUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// input: "VirtualMachine"
// output: "virtual_machine"
func ToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	result.WriteByte(s[0])

	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteByte(s[i])
	}

	return strings.ToLower(result.String())
}

// input: "ParametersGroupService"
// output: "ParameterGroupService"
func RemovePlural(s string) string {
	if len(s) == 0 {
		return s
	}

	// Converter para snake_case para separar as palavras
	snake := ToSnakeCase(s)
	words := strings.Split(snake, "_")

	// Processar cada palavra para remover plurais
	for i, word := range words {
		words[i] = singularize(word)
	}

	// Reconstruir em PascalCase
	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(FirstUpper(word))
		}
	}

	return result.String()
}

// singularize converte uma palavra do plural para singular
func singularize(word string) string {
	// Regras comuns para remoção de plurais em inglês
	if strings.HasSuffix(word, "ies") && len(word) > 3 {
		// Ex: "parameters" -> "parameter"
		return word[:len(word)-3] + "y"
	}
	if strings.HasSuffix(word, "s") && len(word) > 1 {
		// Ex: "groups" -> "group", "services" -> "service"
		// Mas não remover 's' de palavras que terminam com 's' no singular
		// como "status", "process", etc.
		singularExceptions := map[string]bool{
			"status": true, "process": true, "access": true, "address": true,
			"class": true, "glass": true, "grass": true, "mass": true,
			"pass": true, "press": true, "stress": true, "success": true,
		}

		singular := word[:len(word)-1]
		if !singularExceptions[singular] {
			return singular
		}
	}

	return word
}
