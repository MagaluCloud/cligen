package sdk_structure

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/config"
	"golang.org/x/tools/go/packages"
)

// analyzePackageWithParseDir analisa todo o diretório do package usando parser.ParseDir com cache
func analyzePackageWithParseDir(sdkDir string) ([]*packages.Package, *token.FileSet, error) {

	fset := token.NewFileSet()

	cfg := &packages.Config{Fset: fset, Dir: sdkDir, Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedName | packages.NeedFiles | packages.NeedImports}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao analisar diretório %s: %w", sdkDir, err)
	}

	for _, pkg := range pkgs {
		fmt.Println(pkg)
	}

	return pkgs, fset, nil
}

func isCompatible(a, b string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z]+$`)
	return regex.MatchString(a) && regex.MatchString(b) && strings.EqualFold(a, b)
}

// analyzeServiceWithPackage analisa um serviço usando o package já parseado
func analyzeServiceWithPackage(menu *config.Menu, pkgs []*packages.Package, fset *token.FileSet, serviceName string, sdkDir string) Service {
	description := menu.Description
	longDescription := menu.LongDescription
	for _, menuItem := range menu.Menus {
		if isCompatible(menuItem.Name, serviceName) {
			description = menuItem.Description
			longDescription = menuItem.LongDescription
			break
		}
	}

	service := Service{
		Name:            serviceName,
		Description:     description,
		LongDescription: longDescription,
		Methods:         []Method{},
		SubServices:     make(map[string]Service),
		Interface:       serviceName,
	}

	possibleInterfaceNames := getPossibleInterfaceNames(serviceName)

	for _, pkg := range pkgs {
		doBreak := false
		for _, astFile := range pkg.Syntax {
			fileName := fset.File(astFile.Pos()).Name()
			if strings.HasSuffix(fileName, "_test.go") {
				continue
			}

			if found := analyzeFileForServiceWithAST(menu, astFile, possibleInterfaceNames, &service, pkg.Name, sdkDir); found {
				doBreak = true
				break
			}
		}
		if doBreak {
			break
		}
	}

	return service
}

// analyzeFileForServiceWithAST analisa um arquivo AST procurando por interfaces de serviço
func analyzeFileForServiceWithAST(menu *config.Menu, file *ast.File, possibleInterfaceNames []string, service *Service, packageName string, sdkDir string) bool {
	found := false

	ast.Inspect(file, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			if interfaceType, ok := typeDecl.Type.(*ast.InterfaceType); ok {
				for _, interfaceName := range possibleInterfaceNames {
					if typeDecl.Name.Name == interfaceName || strings.EqualFold(typeDecl.Name.Name, interfaceName) {
						service.Interface = typeDecl.Name.Name
						found = true
						for _, method := range interfaceType.Methods.List {
							if funcType, ok := method.Type.(*ast.FuncType); ok {
								methodName := method.Names[0].Name

								if strings.ToLower(methodName) == "listall" {
									continue
								}

								methodDescription := ""
								methodLongDescription := ""
								var confirmation *config.Confirmation

								for _, menu := range menu.Menus {
									if menu.Name == methodName {
										methodDescription = menu.Description
										methodLongDescription = menu.LongDescription
										if menu.Confirmation != nil {
											confirmation = menu.Confirmation
										}
										break
									}
								}

								if methodDescription == "" {
									doBreak := false
									for _, menuLvl2 := range menu.Menus {
										for _, menu := range menuLvl2.Menus {
											if menu.Name == methodName {
												methodDescription = menu.Description
												methodLongDescription = menu.LongDescription
												if menu.Confirmation != nil {
													confirmation = menu.Confirmation
												}
												doBreak = true
												break
											}
										}
										if doBreak {
											break
										}
									}
								}

								params := []Parameter{}
								returns := []Parameter{}

								if funcType.Params != nil {
									for _, param := range funcType.Params.List {
										paramType, aliasType, isPrimitive := getTypeStringWithPackage(param.Type, packageName)
										if aliasType != "" {
											aliasType = packageName + "Sdk." + aliasType
										}
										structFields := analyzeStructType(param.Type, packageName, sdkDir)
										isPointer := strings.HasPrefix(paramType, "*")
										isArray := strings.HasPrefix(paramType, "[]")
										isOptional := false
										if param.Tag != nil {
											tagValue := extractJSONTag(param.Tag.Value)
											if slices.Contains(tagValue, "omitempty") {
												isOptional = true
											}
										}
										for _, paramName := range param.Names {
											params = append(params, Parameter{
												Name:        paramName.Name,
												Type:        paramType,
												IsPrimitive: isPrimitive,
												IsPointer:   isPointer,
												IsArray:     isArray,
												IsOptional:  isOptional,
												Struct:      structFields,
												AliasType:   aliasType,
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
												IsOptional:  isOptional,
												Struct:      structFields,
												AliasType:   aliasType,
											})
										}
									}
								}

								if funcType.Results != nil {
									for _, result := range funcType.Results.List {
										returnType, aliasType, _ := getTypeStringWithPackage(result.Type, packageName)
										if aliasType != "" {
											aliasType = packageName + "Sdk." + aliasType
										}
										structFields := analyzeStructType(result.Type, packageName, sdkDir)
										isPointer := strings.HasPrefix(returnType, "*")
										isArray := strings.HasPrefix(returnType, "[]")
										isOptional := false
										if result.Tag != nil {
											tagValue := extractJSONTag(result.Tag.Value)
											if slices.Contains(tagValue, "omitempty") {
												isOptional = true
											}
										}
										for _, resultName := range result.Names {
											returns = append(returns, Parameter{
												Name:       resultName.Name,
												Type:       returnType,
												IsPointer:  isPointer,
												IsArray:    isArray,
												Struct:     structFields,
												IsOptional: isOptional,
												AliasType:  aliasType,
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
												AliasType: aliasType,
											})
										}
									}
								}

								method := Method{
									Name:            methodName,
									Parameters:      params,
									Returns:         returns,
									Comments:        methodDescription,
									Description:     methodDescription,
									LongDescription: methodLongDescription,
									Confirmation:    confirmation,
								}
								service.Methods = append(service.Methods, method)
								if len(returns) == 1 {
									for _, returnType := range returns {
										if isSubServiceType(returnType.Type) {
											_ = extractSubServiceName(returnType.Type, methodName)
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

func genCliCodeFromClient(menu *config.Menu, pkg *Package, sdkDir, filePath string) []Service {
	astPkg, fset, err := analyzePackageWithParseDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao analisar package: %v", err)
	}

	var services []Service
	var clientMethods []ClientMethod

	var clientFile *ast.File
	for _, pkg := range astPkg {
		for _, file := range pkg.Syntax {
			fileName := fset.File(file.Pos()).Name()
			if filepath.Base(fileName) == "client.go" {
				clientFile = file
				break
			}
		}
		if clientFile != nil {
			break
		}
	}

	if clientFile == nil {
		log.Fatalf("Arquivo client.go não encontrado no package")
	}

	ast.Inspect(clientFile, func(n ast.Node) bool {
		if clientFile.Doc != nil {
			pkg.LongDescription = clientFile.Doc.Text()
		}

		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if !slices.Contains(ignoredFunctions, funcDecl.Name.Name) {
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

	for _, clientMethod := range clientMethods {
		service := analyzeServiceWithPackage(menu, astPkg, fset, clientMethod.ServiceName, sdkDir)
		services = append(services, service)
	}

	return services
}
