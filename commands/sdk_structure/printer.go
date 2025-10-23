package sdk_structure

import (
	"fmt"
)

// printSDKStructure exibe a estrutura do SDK encontrada
func PrintSDKStructure(sdk *SDKStructure) {
	fmt.Println("=== Estrutura do SDK Encontrada ===")
	for pkgName, pkg := range sdk.Packages {
		fmt.Printf("\n📦 Pacote: %s\n", pkgName)
		fmt.Printf("   Menu Name: %s\n", pkg.MenuName)
		fmt.Printf("   Serviços encontrados: %d\n", len(pkg.Services))
		fmt.Printf("   Subpacotes encontrados: %d\n", len(pkg.SubPkgs))

		// Exibir serviços
		for _, service := range pkg.Services {
			printService(service, "   ")
		}

		// Exibir subpacotes
		if len(pkg.SubPkgs) > 0 {
			fmt.Printf("   📁 Subpacotes:\n")
			for subPkgName, subPkg := range pkg.SubPkgs {
				printPackage(subPkg, "      ", subPkgName)
			}
		}
	}
}

// printPackage exibe um pacote e seus subpacotes de forma recursiva
func printPackage(pkg Package, indent string, pkgName string) {
	fmt.Printf("%s📦 Subpacote: %s\n", indent, pkgName)
	fmt.Printf("%s   Menu Name: %s\n", indent, pkg.MenuName)
	fmt.Printf("%s   Serviços encontrados: %d\n", indent, len(pkg.Services))
	fmt.Printf("%s   Subpacotes encontrados: %d\n", indent, len(pkg.SubPkgs))

	// Exibir serviços
	for _, service := range pkg.Services {
		printService(service, indent+"   ")
	}

	// Exibir subpacotes recursivamente
	if len(pkg.SubPkgs) > 0 {
		fmt.Printf("%s   📁 Subpacotes:\n", indent)
		for subPkgName, subPkg := range pkg.SubPkgs {
			printPackage(subPkg, indent+"      ", subPkgName)
		}
	}
}

// printService exibe um serviço e seus subserviços de forma recursiva
func printService(service Service, indent string) {
	fmt.Printf("%s🔧 Serviço: %s\n", indent, service.Name)
	fmt.Printf("%s   Métodos: %d\n", indent, len(service.Methods))

	for _, method := range service.Methods {
		fmt.Printf("%s   - %s(", indent, method.Name)

		// Exibir parâmetros
		paramCount := 0
		for _, param := range method.Parameters {
			if paramCount > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s %s", param.Name, param.Type)
			paramCount++
		}
		fmt.Print(")")

		// Exibir retornos
		if len(method.Returns) > 0 {
			fmt.Print(" -> ")
			returnCount := 0
			for _, ret := range method.Returns {
				if returnCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s %s", ret.Name, ret.Type)
				returnCount++
			}
		}

		// Exibir detalhes dos parâmetros com structs
		if len(method.Parameters) > 0 {
			fmt.Printf("%s     📋 Parâmetros detalhados:\n", indent)
			for _, param := range method.Parameters {
				printParameterDetails(param, indent+"       ")
			}
		}

		// Exibir detalhes dos retornos com structs
		if len(method.Returns) > 0 {
			fmt.Printf("%s     📤 Retornos detalhados:\n", indent)
			for _, ret := range method.Returns {
				printParameterDetails(ret, indent+"       ")
			}
		}
	}

	// Exibir subserviços
	if len(service.SubServices) > 0 {
		fmt.Printf("%s   Subserviços: %d\n", indent, len(service.SubServices))
		for subServiceName, subService := range service.SubServices {
			fmt.Printf("%s   📋 Subserviço: %s\n", indent, subServiceName)
			printService(subService, indent+"      ")
		}
	}
}

// printParameterDetails exibe detalhes de um parâmetro, incluindo campos de struct
func printParameterDetails(param Parameter, indent string) {
	fmt.Printf("%s- %s (%s)", indent, param.Name, param.Type)
	if param.Description != "" {
		fmt.Printf(" - %s", param.Description)
	}

	// Se tem campos de struct, exibir recursivamente
	if param.Struct != nil {
		fmt.Printf("%s  📋 Campos da struct:\n", indent)
		for fieldName, field := range param.Struct {
			fmt.Printf("%s    - %s (%s)", indent, fieldName, field.Type)
			if field.Description != "" {
				fmt.Printf(" - %s", field.Description)
			}

			// Recursão para campos aninhados
			if field.Struct != nil {
				printStructFields(field.Struct, indent+"      ")
			}
		}
	}
}

// printStructFields exibe campos de uma struct de forma recursiva
func printStructFields(fields map[string]Parameter, indent string) {
	for fieldName, field := range fields {
		fmt.Printf("%s- %s (%s)", indent, fieldName, field.Type)
		if field.Description != "" {
			fmt.Printf(" - %s", field.Description)
		}

		// Recursão para campos aninhados
		if field.Struct != nil && len(field.Struct) > 0 {
			printStructFields(field.Struct, indent+"  ")
		}
	}
}
