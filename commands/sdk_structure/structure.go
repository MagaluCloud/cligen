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

	strutils "cligen/str_utils"
)

// Service representa um serviço individual com seus métodos
type Service struct {
	Name        string             `json:"name"`
	Interface   string             `json:"interface"`
	Methods     []Method           `json:"methods"`
	SubServices map[string]Service `json:"sub_services,omitempty"` // Para subserviços aninhados
}

type Parameter struct {
	Position    int                  `json:"position"`
	Name        string               `json:"name"`
	Type        string               `json:"type"`
	Description string               `json:"description"`
	IsPrimitive bool                 `json:"is_primitive"`
	Struct      map[string]Parameter `json:"struct,omitempty"`
}

// Method representa um método de um serviço
type Method struct {
	Name       string      `json:"name"`
	Parameters []Parameter `json:"parameters"` // nome -> tipo
	Returns    []Parameter `json:"returns"`    // nome -> tipo
	Comments   string      `json:"comments"`
}

// Package representa um pacote do SDK com seus serviços
type Package struct {
	MenuName string             `json:"menu_name"`
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

	// Processar menus principais e seus submenus
	for _, menu := range config.Menus {
		processMenu(menu, sdkStructure)
	}

	return *sdkStructure, nil
}

// processMenu processa um menu e seus submenus recursivamente
func processMenu(menu config.Menu, sdkStructure *SDKStructure) {
	processMenuRecursive(menu, "", sdkStructure)
}

// processMenuRecursive processa um menu e seus submenus recursivamente com suporte a hierarquia
func processMenuRecursive(menu config.Menu, parentPath string, sdkStructure *SDKStructure) {
	fmt.Printf("🔄 Processando menu: %s (caminho pai: %s)\n", menu.Name, parentPath)

	// Se o menu tem submenus, criar um pacote de agrupamento
	if len(menu.Menus) > 0 {
		fmt.Printf("📁 Menu '%s' é um agrupador com %d submenus\n", menu.Name, len(menu.Menus))

		// Criar um pacote vazio para o menu de agrupamento
		groupPkg := Package{
			MenuName: menu.Name,
			Name:     menu.Name,
			Services: []Service{},
			SubPkgs:  make(map[string]Package),
		}

		// Construir o caminho atual para este menu
		currentPath := menu.Name
		if parentPath != "" {
			currentPath = filepath.Join(parentPath, menu.Name)
		}

		fmt.Printf("📍 Caminho atual para menu '%s': %s\n", menu.Name, currentPath)

		// Adicionar subpacotes para cada submenu
		for _, submenu := range menu.Menus {
			fmt.Printf("  🔍 Processando submenu: %s\n", submenu.Name)

			if submenu.SDKPackage != "" {
				fmt.Printf("  📦 Submenu '%s' tem SDK Package: %s\n", submenu.Name, submenu.SDKPackage)
				// Para menus filhos, o diretório será dentro do diretório pai
				subPkg := genCliCodeFromSDK(currentPath, submenu.SDKPackage)
				subPkg.MenuName = submenu.Name
				groupPkg.SubPkgs[submenu.SDKPackage] = subPkg
			} else if len(submenu.Menus) > 0 {
				fmt.Printf("  📁 Submenu '%s' é um agrupador com %d sub-submenus\n", submenu.Name, len(submenu.Menus))
				// Se o submenu também tem submenus, processar recursivamente
				// Criar um subpacote de agrupamento
				subGroupPkg := Package{
					MenuName: submenu.Name,
					Name:     submenu.Name,
					Services: []Service{},
					SubPkgs:  make(map[string]Package),
				}

				// Processar submenus do submenu
				for _, subSubmenu := range submenu.Menus {
					fmt.Printf("    🔍 Processando sub-submenu: %s\n", subSubmenu.Name)

					if subSubmenu.SDKPackage != "" {
						fmt.Printf("    📦 Sub-submenu '%s' tem SDK Package: %s\n", subSubmenu.Name, subSubmenu.SDKPackage)
						// Para sub-submenus, o diretório será dentro do diretório do submenu pai
						subSubPkg := genCliCodeFromSDK(filepath.Join(currentPath, submenu.Name), subSubmenu.SDKPackage)
						subSubPkg.MenuName = subSubmenu.Name
						subGroupPkg.SubPkgs[subSubmenu.SDKPackage] = subSubPkg
					} else if len(subSubmenu.Menus) > 0 {
						fmt.Printf("    📁 Sub-submenu '%s' é um agrupador com %d sub-sub-submenus\n", subSubmenu.Name, len(subSubmenu.Menus))
						// Recursão para níveis mais profundos
						processMenuRecursive(subSubmenu, filepath.Join(currentPath, submenu.Name), sdkStructure)
					}
				}

				groupPkg.SubPkgs[submenu.Name] = subGroupPkg
			}
		}

		// Adicionar o pacote ao nível apropriado
		if parentPath == "" {
			// Menu principal - adicionar diretamente ao SDKStructure
			fmt.Printf("✅ Adicionando menu principal '%s' ao SDKStructure\n", menu.Name)
			sdkStructure.Packages[menu.Name] = groupPkg
		} else {
			// Submenu - adicionar ao pacote pai
			// Nota: Aqui precisamos adicionar ao pacote pai correto
			// Por enquanto, vamos adicionar diretamente ao SDKStructure com um nome único
			packageKey := filepath.Join(parentPath, menu.Name)
			fmt.Printf("✅ Adicionando submenu '%s' ao SDKStructure com chave: %s\n", menu.Name, packageKey)
			sdkStructure.Packages[packageKey] = groupPkg
		}
	} else if menu.SDKPackage != "" {
		fmt.Printf("📦 Menu '%s' tem SDK Package: %s\n", menu.Name, menu.SDKPackage)
		// Se o menu não tem submenus mas tem SDKPackage, processá-lo como um pacote normal
		pkg := genCliCodeFromSDK(parentPath, menu.SDKPackage)
		pkg.MenuName = menu.Name

		// Adicionar ao nível apropriado
		if parentPath == "" {
			// Menu principal
			fmt.Printf("✅ Adicionando menu principal com SDK '%s' ao SDKStructure\n", menu.SDKPackage)
			sdkStructure.Packages[menu.SDKPackage] = pkg
		} else {
			// Submenu - adicionar com nome único
			packageKey := filepath.Join(parentPath, menu.SDKPackage)
			fmt.Printf("✅ Adicionando submenu com SDK '%s' ao SDKStructure com chave: %s\n", menu.SDKPackage, packageKey)
			sdkStructure.Packages[packageKey] = pkg
		}
	} else {
		fmt.Printf("⚠️  Menu '%s' não tem submenus nem SDK Package (menu vazio)\n", menu.Name)
	}
}

// Agora iremos utilizar go/ast e go/parser para analisar o código fonte do SDK e gerar o código da CLI
// O SDK foi antereiormente clonado no diretório tmp-sdk/
// parentDir é o diretório pai (pode ser vazio para menus principais)
// packageName é o nome do pacote SDK
func genCliCodeFromSDK(parentDir, packageName string) Package {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter diretório atual: %v", err)
	}

	// Construir o caminho do SDK baseado na hierarquia
	sdkDir := filepath.Join(dir, "tmp-sdk", packageName)
	fmt.Printf("🔍 Procurando SDK em diretório principal: %s\n", sdkDir)

	pkg := Package{
		MenuName: packageName,
		Name:     packageName,
		Services: []Service{},
		SubPkgs:  make(map[string]Package),
	}

	// Verificar se o diretório do SDK existe
	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {
		// Se o diretório não existe, retornar um pacote vazio (para menus de agrupamento)
		fmt.Printf("⚠️  Diretório do SDK não encontrado: %s (menu de agrupamento)\n", sdkDir)
		return pkg
	}

	fmt.Printf("✅ Diretório do SDK encontrado: %s\n", sdkDir)

	files, err := os.ReadDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório do SDK: %v", err)
	}

	fmt.Printf("📄 Total de arquivos no diretório: %d\n", len(files))

	for _, file := range files {
		if file.Name() == "client.go" {
			fmt.Printf("🔧 Processando arquivo client.go em: %s\n", sdkDir)
			services := genCliCodeFromClient(sdkDir, filepath.Join(sdkDir, file.Name()))
			pkg.Services = services
			fmt.Printf("✅ Processados %d serviços do pacote %s\n", len(services), packageName)
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
		"Service",
		serviceName, // Ex: Instances (sem sufixo)
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

	// Adicionar variações para serviços que podem usar singular
	singularNamePascal := strutils.RemovePlural(serviceName)
	possibleInterfaceNames = append(possibleInterfaceNames,
		singularNamePascal+"Service", // Ex: InstanceService
		singularNamePascal+"API",     // Ex: InstanceAPI
		singularNamePascal+"Client",  // Ex: InstanceClient
		singularNamePascal,           // Ex: Instance (sem sufixo)
	)

	fileNamesToTry := []string{fmt.Sprintf("%s.go", strings.ToLower(serviceName))}

	fileName := fmt.Sprintf("%s.go", strutils.ToSnakeCase(serviceName, ""))
	if !slices.Contains(fileNamesToTry, fileName) {
		fileNamesToTry = append(fileNamesToTry, fileName)
	}

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
													params = append(params, Parameter{
														Position:    i,
														Name:        name.Name,
														Type:        paramType,
														IsPrimitive: isPrimitive,
														Description: param.Comment.Text(),
														Struct:      structFields,
													})
												}
											} else {
												// Parâmetro sem nome - gerar nome baseado no tipo
												paramName := generateParamName(paramType, i)
												params = append(params, Parameter{
													Position:    i,
													Name:        paramName,
													Type:        paramType,
													IsPrimitive: isPrimitive,
													Description: param.Comment.Text(),
													Struct:      structFields,
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
													returns = append(returns, Parameter{
														Position:    i,
														Name:        name.Name,
														Type:        returnType,
														IsPrimitive: isPrimitive,
														Description: result.Comment.Text(),
														Struct:      structFields,
													})
												}
											} else {
												// Retorno sem nome - gerar nome baseado no tipo
												returnName := generateReturnName(returnType, i)
												returns = append(returns, Parameter{
													Position:    i,
													Name:        returnName,
													Type:        returnType,
													IsPrimitive: isPrimitive,
													Description: result.Comment.Text(),
													Struct:      structFields,
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
func getTypeString(expr ast.Expr) (string, bool) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, false
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
		return packageName + "." + t.Name, false
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
	if strings.Contains(typeName, ".") {
		typeName = strings.Split(typeName, ".")[len(strings.Split(typeName, "."))-1]
	}
	primitiveTypes := []string{
		"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr", "string",
	}
	return slices.Contains(primitiveTypes, typeName)
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
		fmt.Println()

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
	fmt.Println()

	// Se tem campos de struct, exibir recursivamente
	if param.Struct != nil && len(param.Struct) > 0 {
		fmt.Printf("%s  📋 Campos da struct:\n", indent)
		for fieldName, field := range param.Struct {
			fmt.Printf("%s    - %s (%s)", indent, fieldName, field.Type)
			if field.Description != "" {
				fmt.Printf(" - %s", field.Description)
			}
			fmt.Println()

			// Recursão para campos aninhados
			if field.Struct != nil && len(field.Struct) > 0 {
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
		fmt.Println()

		// Recursão para campos aninhados
		if field.Struct != nil && len(field.Struct) > 0 {
			printStructFields(field.Struct, indent+"  ")
		}
	}
}

// analyzeParameterType analisa um tipo de parâmetro e retorna informações detalhadas incluindo campos de struct
func analyzeParameterType(expr ast.Expr, packageName string, sdkDir string) (string, bool, map[string]Parameter) {
	// Verificar se é uma struct inline (anônima)
	if structType, ok := expr.(*ast.StructType); ok {
		fmt.Printf("   🔍 Struct inline detectada\n")
		structFields := extractStructFields(structType, packageName, sdkDir)
		return "struct{}", false, structFields
	}

	// Verificar se é uma interface inline (anônima)
	if _, ok := expr.(*ast.InterfaceType); ok {
		fmt.Printf("   🔍 Interface inline detectada\n")
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

	// Procurar pela definição da struct no diretório do SDK
	structFields := findStructDefinition(typeName, sdkDir, packageName)
	if structFields != nil {
		fmt.Printf("   🔍 Struct encontrada: %s com %d campos\n", typeName, len(structFields))
	}

	return structFields
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
					fmt.Printf("   ✅ Struct encontrada: %s\n", typeName)
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

	for i, field := range structType.Fields.List {
		// Extrair comentários do campo
		var description string
		if field.Doc != nil {
			description = field.Doc.Text()
		} else if field.Comment != nil {
			description = field.Comment.Text()
		}

		// Extrair tipo do campo
		fieldType, isPrimitive, structFields := analyzeParameterType(field.Type, packageName, sdkDir)

		// Se o campo tem nome, usar o nome, senão gerar um nome baseado no tipo
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				// Verificar se há tags JSON
				var jsonName string
				if field.Tag != nil {
					jsonName = extractJSONTag(field.Tag.Value)
				}
				if jsonName == "" {
					jsonName = name.Name
				}

				fields[jsonName] = Parameter{
					Position:    i,
					Name:        name.Name,
					Type:        fieldType,
					Description: description,
					IsPrimitive: isPrimitive,
					Struct:      structFields,
				}
			}
		} else {
			// Campo anônimo (embedded struct)
			fieldName := generateFieldName(fieldType, i)
			fields[fieldName] = Parameter{
				Position:    i,
				Name:        fieldName,
				Type:        fieldType,
				Description: description,
				IsPrimitive: isPrimitive,
				Struct:      structFields,
			}
		}
	}

	return fields
}

// extractJSONTag extrai o nome do campo da tag JSON
func extractJSONTag(tagValue string) string {
	// Remover aspas
	tagValue = strings.Trim(tagValue, "`\"")

	// Procurar pela tag json
	if strings.Contains(tagValue, "json:") {
		parts := strings.Split(tagValue, " ")
		for _, part := range parts {
			if strings.HasPrefix(part, "json:") {
				jsonValue := strings.TrimPrefix(part, "json:")
				jsonValue = strings.Trim(jsonValue, "\"")

				// Se há vírgula, pegar apenas a primeira parte (nome do campo)
				if strings.Contains(jsonValue, ",") {
					jsonValue = strings.Split(jsonValue, ",")[0]
				}

				// Se o valor é "-", ignorar o campo
				if jsonValue == "-" {
					return ""
				}

				return jsonValue
			}
		}
	}

	return ""
}

// generateFieldName gera um nome para um campo anônimo baseado no tipo
func generateFieldName(fieldType string, index int) string {
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
