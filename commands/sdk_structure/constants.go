package sdk_structure

import (
	"fmt"
	"strings"

	strutils "cligen/str_utils"
)

// Sufixos comuns para serviços
var serviceSuffixes = []string{"Service", "API", "Client"}

// Palavras-chave de serviço
var serviceKeywords = []string{"service", "api", "client"}

// Tipos primitivos do Go
var primitiveTypes = []string{
	"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
	"int", "int8", "int16", "int32", "int64", "rune", "string", "uint",
	"uint8", "uint16", "uint32", "uint64", "uintptr", "string",
}

// Nomes de arquivos esperados para serviços
func getExpectedServiceFileNames(serviceName string) []string {
	return []string{
		fmt.Sprintf("%s.go", strings.ToLower(serviceName)),
		fmt.Sprintf("%s.go", strutils.ToSnakeCase(serviceName, "")),
	}
}

// Nomes possíveis de interface para um serviço
func getPossibleInterfaceNames(serviceName string) []string {
	names := []string{
		serviceName + "Service", // Ex: InstancesService
		serviceName + "API",     // Ex: InstancesAPI
		serviceName + "Client",  // Ex: InstancesClient
		"Service",
		serviceName, // Ex: Instances (sem sufixo)
	}

	// Adicionar variações para serviços que podem usar singular
	if strings.HasSuffix(serviceName, "s") {
		singularName := strings.TrimSuffix(serviceName, "s")
		names = append(names,
			singularName+"Service", // Ex: InstanceService
			singularName+"API",     // Ex: InstanceAPI
			singularName+"Client",  // Ex: InstanceClient
			singularName,           // Ex: Instance (sem sufixo)
		)
	}

	// Adicionar variações para serviços que podem usar singular
	singularNamePascal := strutils.RemovePlural(serviceName)
	names = append(names,
		singularNamePascal+"Service", // Ex: InstanceService
		singularNamePascal+"API",     // Ex: InstanceAPI
		singularNamePascal+"Client",  // Ex: InstanceClient
		singularNamePascal,           // Ex: Instance (sem sufixo)
	)

	return names
}
