package sdk_structure

import (
	"cligen/config"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Service representa um serviço individual com seus métodos
type Service struct {
	Name        string             `json:"name"`
	Interface   string             `json:"interface"`
	Methods     []Method           `json:"methods"`
	SubServices map[string]Service `json:"sub_services,omitempty"` // Para subserviços aninhados
}

// Method representa um método de um serviço
type Method struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"` // nome -> tipo
	Returns    map[string]string `json:"returns"`    // nome -> tipo
	Comments   string            `json:"comments"`
}

// Package representa um pacote do SDK com seus serviços
type Package struct {
	Name     string             `json:"name"`
	Services []Service          `json:"services"`
	SubPkgs  map[string]Package `json:"sub_packages,omitempty"` // Para suporte recursivo
}

// SDKStructure representa a estrutura completa do SDK
type SDKStructure struct {
	Packages map[string]Package `json:"packages"`
}

// ClientMethod representa um método do cliente que retorna um serviço
type ClientMethod struct {
	Name        string
	ReturnType  string
	ServiceName string
}

func GenCliSDKStructure() (SDKStructure, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	sdkStructure := &SDKStructure{
		Packages: make(map[string]Package),
	}

	for _, menu := range config.Menus {
		pkg := genCliCodeFromSDK(menu.SDKPackage)
		sdkStructure.Packages[menu.SDKPackage] = pkg
	}

	return *sdkStructure, nil
}

// Agora iremos utilizar go/ast e go/parser para analisar o código fonte do SDK e gerar o código da CLI
// O SDK foi antereiormente clonado no diretório tmp-sdk/
func genCliCodeFromSDK(packageName string) Package {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter diretório atual: %v", err)
	}

	sdkDir := filepath.Join(dir, "tmp-sdk", packageName)

	files, err := os.ReadDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório do SDK: %v", err)
	}

	pkg := Package{
		Name:     packageName,
		Services: []Service{},
		SubPkgs:  make(map[string]Package),
	}

	for _, file := range files {
		if file.Name() == "client.go" {
			services := genCliCodeFromClient(sdkDir, filepath.Join(sdkDir, file.Name()))
			pkg.Services = services
		}
	}

	return pkg
}

var ignoredFunctions = []string{"newRequest", "newResponse"}

// Por padrão do SDK, cada pacote possui um arquivo client.go que contém a estrutura do cliente e os serviços disponíveis
// Vamos utilizar o go/ast e go/parser para analisar o código fonte do arquivo client.go e gerar o código da CLI
func genCliCodeFromClient(sdkDir, filePath string) []Service {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Erro ao analisar o arquivo %s: %v", filePath, err)
	}

	var services []Service
	var clientMethods []ClientMethod

	// Primeiro, vamos encontrar os métodos do cliente que retornam serviços
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				// É um método do cliente
				if !slices.Contains(ignoredFunctions, funcDecl.Name.Name) {
					// Verificar se retorna um tipo de serviço
					if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
						resultType := funcDecl.Type.Results.List[0].Type
						if typeName, ok := resultType.(*ast.Ident); ok {
							clientMethod := ClientMethod{
								Name:        funcDecl.Name.Name,
								ReturnType:  typeName.Name,
								ServiceName: funcDecl.Name.Name,
							}
							clientMethods = append(clientMethods, clientMethod)
						}
					}
				}
			}
		}
		return true
	})

	// Agora vamos analisar cada serviço encontrado
	for _, clientMethod := range clientMethods {
		service := analyzeService(sdkDir, filePath, clientMethod.ServiceName)
		services = append(services, service)

	}

	return services
}

// analyzeService analisa um serviço específico para extrair seus métodos
func analyzeService(sdkDir, clientFilePath, serviceName string) Service {
	service := Service{
		Name:        serviceName,
		Methods:     []Method{},
		SubServices: make(map[string]Service),
		Interface:   serviceName,
	}

	// Aqui vamos mapear os arquivos que existem no pacote, e entao buscaremos pela interface do serviço
	files, err := os.ReadDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório do SDK: %v", err)
	}

	fmt.Printf("🔍 Procurando serviço: %s\n", serviceName)
	fmt.Printf("📁 Diretório: %s\n", sdkDir)
	fmt.Printf("📄 Total de arquivos no pacote: %d\n", len(files))

	// Possíveis nomes de interface para o serviço
	possibleInterfaceNames := []string{
		serviceName + "Service", // Ex: InstancesService
		serviceName + "API",     // Ex: InstancesAPI
		serviceName + "Client",  // Ex: InstancesClient
		serviceName,             // Ex: Instances (sem sufixo)
	}

	// Adicionar variações para serviços que podem usar singular
	if strings.HasSuffix(serviceName, "s") {
		singularName := strings.TrimSuffix(serviceName, "s")
		possibleInterfaceNames = append(possibleInterfaceNames,
			singularName+"Service", // Ex: InstanceService
			singularName+"API",     // Ex: InstanceAPI
			singularName+"Client",  // Ex: InstanceClient
			singularName,           // Ex: Instance (sem sufixo)
		)
	}

	// Primeiro, tentar encontrar o arquivo específico do serviço
	fileName := fmt.Sprintf("%s.go", strings.ToLower(serviceName))
	serviceFilePath := filepath.Join(sdkDir, fileName)

	fmt.Printf("📄 Arquivo esperado: %s\n", fileName)

	// Verificar se o arquivo específico existe
	if _, err := os.Stat(serviceFilePath); err == nil {
		fmt.Printf("✅ Arquivo encontrado: %s\n", serviceFilePath)
		if found := analyzeFileForService(serviceFilePath, possibleInterfaceNames, &service); found {
			return service
		}
	} else {
		fmt.Printf("❌ Arquivo não encontrado: %s\n", serviceFilePath)
	}

	// Se não encontrou no arquivo específico, procurar em todos os arquivos do pacote
	fmt.Printf("🔍 Procurando interface em outros arquivos do pacote...\n")

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		filePath := filepath.Join(sdkDir, file.Name())
		fmt.Printf("🔍 Verificando arquivo: %s\n", file.Name())

		if found := analyzeFileForService(filePath, possibleInterfaceNames, &service); found {
			fmt.Printf("✅ Interface encontrada no arquivo: %s\n", file.Name())
			break
		}
	}

	if len(service.Methods) == 0 {
		fmt.Printf("⚠️  Nenhum método encontrado para o serviço: %s\n", serviceName)
	} else {
		fmt.Printf("✅ Total de métodos encontrados: %d\n", len(service.Methods))
	}

	return service
}

// analyzeFileForService analisa um arquivo específico procurando por interfaces de serviço
func analyzeFileForService(filePath string, possibleInterfaceNames []string, service *Service) bool {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Erro ao analisar o arquivo %s: %v", filePath, err)
		return false
	}

	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeDecl.Type.(*ast.InterfaceType); ok {
				// Verificar se é uma das interfaces que estamos procurando
				for _, interfaceName := range possibleInterfaceNames {
					if typeDecl.Name.Name == interfaceName {
						fmt.Printf("✅ Interface encontrada: %s\n", interfaceName)
						found = true

						// Extrair métodos da interface
						if interfaceType.Methods != nil {
							fmt.Printf("📋 Métodos encontrados: %d\n", len(interfaceType.Methods.List))
							for _, method := range interfaceType.Methods.List {
								if funcType, ok := method.Type.(*ast.FuncType); ok {
									methodName := method.Names[0].Name

									// Extrair comentários
									var comments string
									if method.Doc != nil {
										comments = method.Doc.Text()
									}

									// Extrair parâmetros
									params := make(map[string]string)
									if funcType.Params != nil {
										for i, param := range funcType.Params.List {
											paramType := getTypeString(param.Type)
											// Se o parâmetro tem nome, usar o nome, senão gerar um nome baseado no tipo
											if len(param.Names) > 0 {
												for _, name := range param.Names {
													params[name.Name] = paramType
												}
											} else {
												// Parâmetro sem nome - gerar nome baseado no tipo
												paramName := generateParamName(paramType, i)
												params[paramName] = paramType
											}
										}
									}

									// Extrair retornos
									returns := make(map[string]string)
									if funcType.Results != nil {
										for i, result := range funcType.Results.List {
											returnType := getTypeString(result.Type)
											// Se o retorno tem nome, usar o nome, senão gerar um nome baseado no tipo
											if len(result.Names) > 0 {
												for _, name := range result.Names {
													returns[name.Name] = returnType
												}
											} else {
												// Retorno sem nome - gerar nome baseado no tipo
												returnName := generateReturnName(returnType, i)
												returns[returnName] = returnType
											}
										}
									}

									method := Method{
										Name:       methodName,
										Parameters: params,
										Returns:    returns,
										Comments:   comments,
									}
									service.Interface = interfaceName
									service.Methods = append(service.Methods, method)
									fmt.Printf("   ✅ Método adicionado: %s\n", methodName)

									// Verificar se este método retorna um subserviço
									if len(returns) == 1 {
										for _, returnType := range returns {
											if isSubServiceType(returnType) {
												fmt.Printf("   🔍 Detectado possível subserviço: %s -> %s\n", methodName, returnType)
												subServiceName := extractSubServiceName(returnType, methodName)
												if subServiceName != "" {
													// Analisar o subserviço recursivamente
													subService := analyzeService(filepath.Dir(filePath), filePath, subServiceName)
													if len(subService.Methods) > 0 {
														service.SubServices[subServiceName] = subService
														fmt.Printf("   ✅ Subserviço adicionado: %s (%d métodos)\n", subServiceName, len(subService.Methods))
													}
												}
											}
										}
									}
								}
							}
						}
						break
					}
				}
			}
		}
		return true
	})

	return found
}

// getTypeString converte um ast.Expr para string representando o tipo
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// isSubServiceType verifica se um tipo de retorno representa um subserviço
func isSubServiceType(returnType string) bool {
	// Remover ponteiros e arrays para análise
	baseType := strings.TrimPrefix(returnType, "*")
	baseType = strings.TrimPrefix(baseType, "[]")

	// Verificar se o tipo termina com sufixos comuns de serviço
	serviceSuffixes := []string{"Service", "API", "Client"}
	for _, suffix := range serviceSuffixes {
		if strings.HasSuffix(baseType, suffix) {
			return true
		}
	}

	// Verificar se contém palavras-chave de serviço
	serviceKeywords := []string{"service", "api", "client"}
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
	suffixes := []string{"Service", "API", "Client"}
	for _, suffix := range suffixes {
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
func generateParamName(paramType string, index int) string {
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
func generateReturnName(returnType string, index int) string {
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

// printSDKStructure exibe a estrutura do SDK encontrada
func PrintSDKStructure(sdk *SDKStructure) {
	fmt.Println("=== Estrutura do SDK Encontrada ===")
	for pkgName, pkg := range sdk.Packages {
		fmt.Printf("\n📦 Pacote: %s\n", pkgName)
		fmt.Printf("   Serviços encontrados: %d\n", len(pkg.Services))

		for _, service := range pkg.Services {
			printService(service, "   ")
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
		for paramName, paramType := range method.Parameters {
			if paramCount > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%s %s", paramName, paramType)
			paramCount++
		}
		fmt.Print(")")

		// Exibir retornos
		if len(method.Returns) > 0 {
			fmt.Print(" -> ")
			returnCount := 0
			for retName, retType := range method.Returns {
				if returnCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s %s", retName, retType)
				returnCount++
			}
		}
		fmt.Println()
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
