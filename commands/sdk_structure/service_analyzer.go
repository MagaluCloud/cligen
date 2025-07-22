package sdk_structure

import (
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

// analyzeService analisa um serviço específico para extrair seus métodos
func analyzeService(sdkDir, serviceName string) Service {
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

	// Obter nomes possíveis de interface
	possibleInterfaceNames := getPossibleInterfaceNames(serviceName)

	// Obter nomes de arquivos esperados
	fileNamesToTry := getExpectedServiceFileNames(serviceName)

	// Primeiro, tentar encontrar o arquivo específico do serviço
	for _, fileName := range fileNamesToTry {
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

	// Extrair o nome do pacote do arquivo
	packageName := file.Name.Name + "Sdk"

	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeDecl.Type.(*ast.InterfaceType); ok {
				// Verificar se é uma das interfaces que estamos procurando
				for _, interfaceName := range possibleInterfaceNames {
					if typeDecl.Name.Name == interfaceName || strings.EqualFold(typeDecl.Name.Name, interfaceName) {
						fmt.Printf("✅ Interface encontrada: %s\n", interfaceName)
						found = true
						service.Interface = typeDecl.Name.Name

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
									params := make([]Parameter, 0)
									if funcType.Params != nil {
										for i, param := range funcType.Params.List {
											paramType, isPrimitive, structFields := analyzeParameterType(param.Type, packageName, filepath.Dir(filePath))
											// Se o parâmetro tem nome, usar o nome, senão gerar um nome baseado no tipo
											if len(param.Names) > 0 {
												for _, name := range param.Names {
													isPointer := false
													if strings.HasPrefix(paramType, "*") {
														isPointer = true
													}
													isArray := false
													if strings.HasPrefix(paramType, "[]") || strings.HasPrefix(paramType, "*[]") {
														isArray = true
													}
													params = append(params, Parameter{
														Position:    i,
														Name:        name.Name,
														Type:        paramType,
														IsPrimitive: isPrimitive,
														Description: param.Comment.Text(),
														Struct:      structFields,
														IsPointer:   isPointer,
														IsArray:     isArray,
													})
												}
											} else {
												// Parâmetro sem nome - gerar nome baseado no tipo
												paramName := generateParamName(paramType)
												isPointer := false
												if strings.HasPrefix(paramType, "*") {
													isPointer = true
												}
												isArray := false
												if strings.HasPrefix(paramType, "[]") || strings.HasPrefix(paramType, "*[]") {
													isArray = true
												}
												params = append(params, Parameter{
													Position:    i,
													Name:        paramName,
													Type:        paramType,
													IsPrimitive: isPrimitive,
													Description: param.Comment.Text(),
													Struct:      structFields,
													IsPointer:   isPointer,
													IsArray:     isArray,
												})
											}
										}
									}

									// Extrair retornos
									returns := make([]Parameter, 0)
									if funcType.Results != nil {
										for i, result := range funcType.Results.List {
											returnType, isPrimitive, structFields := analyzeParameterType(result.Type, packageName, filepath.Dir(filePath))
											// Se o retorno tem nome, usar o nome, senão gerar um nome baseado no tipo
											if len(result.Names) > 0 {
												for _, name := range result.Names {
													isPointer := false
													if strings.HasPrefix(returnType, "*") {
														isPointer = true
													}
													isArray := false
													if strings.HasPrefix(returnType, "[]") || strings.HasPrefix(returnType, "*[]") {
														isArray = true
													}
													returns = append(returns, Parameter{
														Position:    i,
														Name:        name.Name,
														Type:        returnType,
														IsPrimitive: isPrimitive,
														Description: result.Comment.Text(),
														Struct:      structFields,
														IsPointer:   isPointer,
														IsArray:     isArray,
													})
												}
											} else {
												// Retorno sem nome - gerar nome baseado no tipo
												returnName := generateReturnName(returnType)
												isPointer := false
												if strings.HasPrefix(returnType, "*") {
													isPointer = true
												}
												isArray := false
												if strings.HasPrefix(returnType, "[]") || strings.HasPrefix(returnType, "*[]") {
													isArray = true
												}
												returns = append(returns, Parameter{
													Position:    i,
													Name:        returnName,
													Type:        returnType,
													IsPrimitive: isPrimitive,
													Description: result.Comment.Text(),
													Struct:      structFields,
													IsPointer:   isPointer,
													IsArray:     isArray,
												})
											}
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

									// Verificar se este método retorna um subserviço
									if len(returns) == 1 {
										for _, returnType := range returns {
											if isSubServiceType(returnType.Type) {
												fmt.Printf("   🔍 Detectado possível subserviço: %s -> %s\n", methodName, returnType.Type)
												subServiceName := extractSubServiceName(returnType.Type, methodName)
												if subServiceName != "" {
													// Analisar o subserviço recursivamente
													subService := analyzeService(filepath.Dir(filePath), subServiceName)
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

var ignoredFunctions = []string{"newRequest", "newResponse"}

// genCliCodeFromClient analisa o arquivo client.go para extrair serviços
func genCliCodeFromClient(pkg *Package, sdkDir, filePath string) []Service {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Erro ao analisar o arquivo %s: %v", filePath, err)
	}

	var services []Service
	var clientMethods []ClientMethod

	// Primeiro, vamos encontrar os métodos do cliente que retornam serviços
	ast.Inspect(file, func(n ast.Node) bool {
		// Collect the comment of header of the file
		if file.Doc != nil {
			pkg.LongDescription = file.Doc.Text()
		}

		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				// É um método do cliente
				if !slices.Contains(ignoredFunctions, funcDecl.Name.Name) {
					// Verificar se retorna um tipo de serviço
					if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
						resultType := funcDecl.Type.Results.List[0].Type
						if typeName, ok := resultType.(*ast.Ident); ok {
							clientMethod := ClientMethod{
								Name:            funcDecl.Name.Name,
								ReturnType:      typeName.Name,
								ServiceName:     funcDecl.Name.Name,
								LongDescription: funcDecl.Doc.Text(),
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
		service := analyzeService(sdkDir, clientMethod.ServiceName)
		services = append(services, service)
	}

	return services
}
