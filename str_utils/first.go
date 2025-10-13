package strutils

import (
	"slices"
	"strings"
)

func FirstLower(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}

func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// input: "VirtualMachine"
// output: "virtual_machine"
func ToSnakeCase(s string, char string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder
	result.WriteByte(s[0])

	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			if char != "" {
				result.WriteByte(char[0])
			} else {
				result.WriteByte('_')
			}
		}
		result.WriteByte(s[i])
	}

	return strings.ToLower(result.String())
}

func ToSnakeCasePreserveID(s string, char string) string {
	if len(s) == 0 {
		return s
	}

	var result strings.Builder

	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			if len(result.String()) == 0 {
				result.WriteByte(s[i])
				continue
			}
			lastResultChar := result.String()[len(result.String())-1]
			if lastResultChar >= 'a' && lastResultChar <= 'z' {
				result.WriteByte('-')
			}
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
	snake := ToSnakeCase(s, "")
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

// ToPascalCase converte uma string com hífens para PascalCase
// Ex: "availability-zones" -> "AvailabilityZones"
func ToPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}

	// Dividir por hífens
	parts := strings.Split(s, "-")
	var result strings.Builder

	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(FirstUpper(part))
		}
	}

	return result.String()
}

func FirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1])
}

func ForbiddenChars(str string) bool {
	// Flags globais que não podem ser usadas como nomes de flags
	return slices.Contains([]string{"k", "d", "o", "n", "r", "h"}, str)
}

func FirstUnusedChar(s string, usedChars *[]string) string {
	// Se a string estiver vazia, buscar no alfabeto
	if len(s) == 0 {
		return findNextAvailableAlphabetChar(usedChars)
	}

	firstChar := FirstChar(s)
	for _, usedChar := range *usedChars {
		if firstChar == usedChar || ForbiddenChars(firstChar) {
			return FirstUnusedChar(s[1:], usedChars)
		}
	}
	*usedChars = append(*usedChars, firstChar)
	return firstChar
}

// findNextAvailableAlphabetChar busca a próxima letra disponível do alfabeto
func findNextAvailableAlphabetChar(usedChars *[]string) string {
	// Criar um mapa para verificação rápida de caracteres usados
	usedMap := make(map[string]bool)
	for _, char := range *usedChars {
		usedMap[char] = true
	}

	// Buscar a próxima letra disponível do alfabeto (a-z)
	for char := 'a'; char <= 'z'; char++ {
		charStr := string(char)
		if !usedMap[charStr] && !ForbiddenChars(charStr) {
			*usedChars = append(*usedChars, charStr)
			return charStr
		}
	}

	// Se todas as letras do alfabeto estiverem usadas, retornar uma string vazia
	// ou lançar um erro, dependendo do comportamento desejado
	return ""
}

func RemoveNewLine(s string) string {
	return strings.ReplaceAll(s, "\n", "")
}

// Input: "xpto "some" xpto"
// Output: "xpto \"some\" xpto"
func EscapeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}
