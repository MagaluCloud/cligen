package sdk_structure

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"slices"
	"strings"
)

// analyzePackageWithParseDir analisa todo o diretório do package usando parser.ParseDir
func analyzePackageWithParseDir(sdkDir string) (*ast.Package, *token.FileSet, error) {
	fset := token.NewFileSet()

	// ParseDir retorna um map[string]*ast.Package, onde a chave é o nome do package
	pkgs, err := parser.ParseDir(fset, sdkDir, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao analisar diretório %s: %v", sdkDir, err)
	}

	// Como estamos analisando um único package, pegamos o primeiro (e único) package
	var pkg *ast.Package
	for _, p := range pkgs {
		pkg = p
		break
	}

	if pkg == nil {
		return nil, nil, fmt.Errorf("nenhum package encontrado em %s", sdkDir)
	}

	fmt.Printf("✅ Package analisado: %s (%d arquivos)\n", pkg.Name, len(pkg.Files))
	return pkg, fset, nil
}

// analyzeServiceWithPackage analisa um serviço usando o package já parseado
func analyzeServiceWithPackage(pkg *ast.Package, fset *token.FileSet, serviceName string, sdkDir string) Service {
	service := Service{
		Name:        serviceName,
		Description: "Dqui1",
		Methods:     []Method{},
		SubServices: make(map[string]Service),
		Interface:   serviceName,
	}

	// fmt.Printf("🔍 Procurando serviço: %s\n", serviceName)
	// fmt.Printf("📄 Total de arquivos no package: %d\n", len(pkg.Files))

	// Obter nomes possíveis de interface
	possibleInterfaceNames := getPossibleInterfaceNames(serviceName)

	// Analisar todos os arquivos do package procurando pela interface
	// found := false
	for fileName, file := range pkg.Files {
		// fmt.Printf("🔍 Verificando arquivo: %s\n", filepath.Base(fileName))

		if strings.HasSuffix(fileName, "test.go") {
			continue
		}

		if lfound := analyzeFileForServiceWithAST(file, possibleInterfaceNames, &service, pkg.Name, sdkDir); lfound {
			// fmt.Printf("✅ Interface encontrada no arquivo: %s\n", filepath.Base(fileName))
			// found = true
			break
		}
	}

	// if !found {
	// 	fmt.Printf("⚠️  Interface não encontrada para o serviço: %s\n", serviceName)
	// } else {
	// 	fmt.Printf("✅ Total de métodos encontrados: %d\n", len(service.Methods))
	// }

	return service
}

// analyzeFileForServiceWithAST analisa um arquivo AST procurando por interfaces de serviço
func analyzeFileForServiceWithAST(file *ast.File, possibleInterfaceNames []string, service *Service, packageName string, sdkDir string) bool {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeDecl.Type.(*ast.InterfaceType); ok {
				// Verificar se é uma das interfaces que estamos procurando
				for _, interfaceName := range possibleInterfaceNames {
					if typeDecl.Name.Name == interfaceName || strings.EqualFold(typeDecl.Name.Name, interfaceName) {
						// fmt.Printf("   ✅ Interface encontrada: %s\n", typeDecl.Name.Name)
						service.Interface = typeDecl.Name.Name
						found = true

						// Analisar os métodos da interface
						for _, method := range interfaceType.Methods.List {
							if funcType, ok := method.Type.(*ast.FuncType); ok {
								// É um método direto da interface
								methodName := method.Names[0].Name
								methodDescription := "doto3"
								if method.Doc != nil {
									methodDescription = method.Doc.Text()
								}

								params := []Parameter{}
								returns := []Parameter{}

								// Analisar parâmetros
								if funcType.Params != nil {
									for k, param := range funcType.Params.List {
										paramType, isPrimitive := getTypeStringWithPackage(param.Type, packageName)
										structFields := analyzeStructType(param.Type, packageName, sdkDir)
										isPointer := strings.HasPrefix(paramType, "*")
										isArray := strings.HasPrefix(paramType, "[]")
										for _, paramName := range param.Names {
											params = append(params, Parameter{
												Name:        paramName.Name,
												Type:        paramType,
												IsPrimitive: isPrimitive,
												IsPointer:   isPointer,
												IsArray:     isArray,
												Position:    k,
												Struct:      structFields,
											})
										}
										if len(param.Names) == 0 {
											paramName := generateParamName(paramType)
											params = append(params, Parameter{
												Name:        paramName,
												Type:        paramType,
												IsPrimitive: isPrimitive,
												IsPointer:   isPointer,
												IsArray:     isArray,
												Position:    k,
												Struct:      structFields,
											})
										}
									}
								}

								// Analisar retornos
								if funcType.Results != nil {
									for _, result := range funcType.Results.List {
										returnType, _ := getTypeStringWithPackage(result.Type, packageName)
										structFields := analyzeStructType(result.Type, packageName, sdkDir)
										isPointer := strings.HasPrefix(returnType, "*")
										isArray := strings.HasPrefix(returnType, "[]")
										for _, resultName := range result.Names {
											returns = append(returns, Parameter{
												Name:      resultName.Name,
												Type:      returnType,
												IsPointer: isPointer,
												IsArray:   isArray,
												Struct:    structFields,
											})
										}
										if len(result.Names) == 0 {
											returnName := generateReturnName(returnType)
											returns = append(returns, Parameter{
												Name:      returnName,
												Type:      returnType,
												IsPointer: isPointer,
												IsArray:   isArray,
												Struct:    structFields,
											})
										}
									}
								}

								method := Method{
									Name:        methodName,
									Parameters:  params,
									Returns:     returns,
									Comments:    methodDescription,
									Description: methodDescription,
								}
								service.Methods = append(service.Methods, method)
								// fmt.Printf("   ✅ Método adicionado: %s\n", methodName)

								// Verificar se este método retorna um subserviço
								if len(returns) == 1 {
									for _, returnType := range returns {
										if isSubServiceType(returnType.Type) {
											fmt.Printf("   🔍 Detectado possível subserviço: %s -> %s\n", methodName, returnType.Type)
											subServiceName := extractSubServiceName(returnType.Type, methodName)
											if subServiceName != "" {
												// Analisar o subserviço recursivamente usando o mesmo package
												// Nota: aqui precisaríamos passar o ast.Package, mas como estamos dentro de analyzeFileForServiceWithAST,
												// vamos usar uma abordagem diferente - analisar o subserviço depois
												fmt.Printf("   🔍 Subserviço detectado: %s (será analisado posteriormente)\n", subServiceName)
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
	// Usar a nova abordagem com ParseDir
	astPkg, fset, err := analyzePackageWithParseDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao analisar package: %v", err)
	}

	var services []Service
	var clientMethods []ClientMethod

	// Procurar pelo arquivo client.go no package
	var clientFile *ast.File
	for fileName, file := range astPkg.Files {
		if filepath.Base(fileName) == "client.go" {
			clientFile = file
			break
		}
	}

	if clientFile == nil {
		log.Fatalf("Arquivo client.go não encontrado no package")
	}

	// Primeiro, vamos encontrar os métodos do cliente que retornam serviços
	ast.Inspect(clientFile, func(n ast.Node) bool {
		// Collect the comment of header of the file
		if clientFile.Doc != nil {
			pkg.LongDescription = clientFile.Doc.Text()
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
		service := analyzeServiceWithPackage(astPkg, fset, clientMethod.ServiceName, sdkDir)
		services = append(services, service)
	}

	return services
}
