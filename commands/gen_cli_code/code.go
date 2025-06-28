package genclicode

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
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
}

// Method representa um método de um serviço
type Method struct {
	Name       string   `json:"name"`
	Parameters []string `json:"parameters"`
	Returns    []string `json:"returns"`
	Comments   string   `json:"comments"`
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

// ServiceInterface representa uma interface de serviço encontrada no código
type ServiceInterface struct {
	Name    string
	Methods []Method
}

// ClientMethod representa um método do cliente que retorna um serviço
type ClientMethod struct {
	Name        string
	ReturnType  string
	ServiceName string
}

func GenCliCode() {
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

	// Exibir a estrutura encontrada
	printSDKStructure(sdkStructure)
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
		Name:    serviceName,
		Methods: []Method{},
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
									var params []string
									if funcType.Params != nil {
										for _, param := range funcType.Params.List {
											paramType := getTypeString(param.Type)
											params = append(params, paramType)
										}
									}

									// Extrair retornos
									var returns []string
									if funcType.Results != nil {
										for _, result := range funcType.Results.List {
											returnType := getTypeString(result.Type)
											returns = append(returns, returnType)
										}
									}

									method := Method{
										Name:       methodName,
										Parameters: params,
										Returns:    returns,
										Comments:   comments,
									}
									service.Methods = append(service.Methods, method)
									fmt.Printf("   ✅ Método adicionado: %s\n", methodName)
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

// printSDKStructure exibe a estrutura do SDK encontrada
func printSDKStructure(sdk *SDKStructure) {
	fmt.Println("=== Estrutura do SDK Encontrada ===")
	for pkgName, pkg := range sdk.Packages {
		fmt.Printf("\n📦 Pacote: %s\n", pkgName)
		fmt.Printf("   Serviços encontrados: %d\n", len(pkg.Services))

		for _, service := range pkg.Services {
			fmt.Printf("   🔧 Serviço: %s\n", service.Name)
			fmt.Printf("      Métodos: %d\n", len(service.Methods))

			for _, method := range service.Methods {
				fmt.Printf("      - %s(", method.Name)
				for i, param := range method.Parameters {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(param)
				}
				fmt.Print(")")

				if len(method.Returns) > 0 {
					fmt.Print(" -> ")
					for i, ret := range method.Returns {
						if i > 0 {
							fmt.Print(", ")
						}
						fmt.Print(ret)
					}
				}
				fmt.Println()
			}
		}
	}
}
