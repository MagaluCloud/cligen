package generate

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/magaluCloud/cligen/config"
	"golang.org/x/tools/go/packages"
)

const (
	SDK_DIR = "./tmp-sdk"
)

var (
	DIRS_TO_SKIP = []string{"internal", "client", "cmd", "helpers", "docs"}
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar configuração: %w", err))
	}
	cfg.Menus = nil

	dirs, err := ListDir(SDK_DIR)
	if err != nil {
		panic(fmt.Errorf("erro ao listar diretório %s: %w", SDK_DIR, err))
	}

	for _, dir := range dirs {
		pkgs, fset, err := GetPkgWithFset(filepath.Join(SDK_DIR, dir))
		if err != nil {
			log.Println(fmt.Errorf("1.erro ao analisar diretório %s: %w", filepath.Join(SDK_DIR, dir), err))
			continue
		}
		if pkgs[0].Errors != nil {
			log.Println(fmt.Errorf("2.erro ao analisar diretório %s: %s", filepath.Join(SDK_DIR, dir), pkgs[0].Errors[0].Msg))
			continue
		}
		if len(pkgs) > 1 {
			log.Println(fmt.Errorf("3.erro ao analisar diretório %s: mais de um pacote encontrado", filepath.Join(SDK_DIR, dir)))
			continue
		}
		if len(pkgs) == 0 {
			log.Println(fmt.Errorf("4.erro ao analisar diretório %s: nenhum pacote encontrado", filepath.Join(SDK_DIR, dir)))
			continue
		}
		if !IsValidModule(pkgs[0].GoFiles) {
			log.Println(fmt.Errorf("5.erro ao analisar diretório %s: modulo invalido", filepath.Join(SDK_DIR, dir)))
			continue
		}

		menu := &config.Menu{
			ID:               GenerateRandomID(),
			Name:             pkgs[0].Name,
			Enabled:          true,
			Description:      "",
			LongDescription:  "",
			SDKPackage:       pkgs[0].ID,
			CliGroup:         "",
			Alias:            nil,
			Menus:            nil,
			ServiceInterface: "",
			Methods:          nil,
			SDKFile:          "",
			CustomFile:       "",
			IsGroup:          false,
			ParentMenu:       nil,
			Pkgs:             pkgs[0],
			Fset:             fset,
			MapFile:          make(map[string]*ast.File),
		}
		cfg.Menus = append(cfg.Menus, menu)
	}

	for _, menu := range cfg.Menus {
		for _, astFile := range menu.Pkgs.Syntax {
			fileName := menu.Fset.File(astFile.Pos()).Name()
			menu.MapFile[fileName] = astFile
		}
	}
	for _, menu := range cfg.Menus {
		// Primeiro nivel sempre virá de client.go
		for _, filePath := range menu.Pkgs.GoFiles {
			astFile, exists := menu.MapFile[filePath]
			if !exists {
				log.Printf("Arquivo AST não encontrado para %s", filePath)
				continue
			}
			ast.Inspect(astFile, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					if funcDecl.Name != nil && len(funcDecl.Name.Name) > 0 {
						firstChar := rune(funcDecl.Name.Name[0])
						if unicode.IsUpper(firstChar) {
							returnType := funcDecl.Type.Results.List[0].Type
							if IsServiceFunction(returnType) {
								subMenu := &config.Menu{
									ID:               GenerateRandomID(),
									Name:             funcDecl.Name.Name,
									Enabled:          true,
									Description:      "",
									LongDescription:  "",
									SDKPackage:       "",
									CliGroup:         "",
									Alias:            nil,
									Menus:            nil,
									ServiceInterface: returnType.(*ast.Ident).Name,
									Methods:          nil,
									SDKFile:          filePath,
									CustomFile:       "",
									IsGroup:          false,
									ParentMenu:       menu,
									Pkgs:             menu.Pkgs,
									Fset:             menu.Fset,
									MapFile:          menu.MapFile,
								}
								menu.Menus = append(menu.Menus, subMenu)

							}
						}
					}
				}
				return true
			})
		}
	}

	// Aqui virá de interfaces
	for _, menu := range cfg.Menus {
		fmt.Printf("%sMenu: %s\n", Ident(ParentMenuCount(menu, nil)), menu.Name)
		for _, subMenu := range menu.Menus {
			fmt.Printf("%sSubMenu: %s\n", Ident(ParentMenuCount(subMenu, nil)), subMenu.Name)
			ProcessMenu(subMenu)
		}

	}

	err = cfg.SaveConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao salvar configuração: %w", err))
	}

	fmt.Println("fim")
}

func GenerateRandomID() string {
	return uuid.New().String()
}

func ParentMenuCount(menu *config.Menu, count *int) int {
	if count == nil {
		count = new(int)
	}
	if menu.ParentMenu != nil {
		*count++
		ParentMenuCount(menu.ParentMenu, count)
	}
	return *count
}

func Ident(count int) string {
	return strings.Repeat("\t", count)
}

func ProcessMenu(menu *config.Menu) {
	doBreak := false
	for _, filePath := range menu.Pkgs.GoFiles {
		astFile, exists := menu.MapFile[filePath]
		if !exists {
			log.Printf("Arquivo AST não encontrado para %s", filePath)
			continue
		}
		ast.Inspect(astFile, func(n ast.Node) bool {
			if typeDecl, ok := n.(*ast.TypeSpec); ok {
				if interfaceType, ok := typeDecl.Type.(*ast.InterfaceType); ok {
					if typeDecl.Name.Name == menu.ServiceInterface {
						for _, method := range interfaceType.Methods.List {
							// verificar se o método tem assinatura com retorno do tipo *Service
							methodIsService := false
							if funcType, ok := method.Type.(*ast.FuncType); ok && funcType.Results != nil {
								for _, result := range funcType.Results.List {
									if IsServiceFunction(result.Type) {
										fmt.Printf("%sSubMenu: %s\n", Ident(ParentMenuCount(menu, nil)+1), result.Type.(*ast.Ident).Name)
										subSubMenu := &config.Menu{
											ID:               GenerateRandomID(),
											Name:             result.Type.(*ast.Ident).Name,
											Enabled:          true,
											Description:      "",
											LongDescription:  "",
											SDKPackage:       "",
											CliGroup:         "",
											Alias:            nil,
											Menus:            nil,
											ServiceInterface: result.Type.(*ast.Ident).Name,
											Methods:          nil,
											SDKFile:          filePath,
											CustomFile:       "",
											IsGroup:          false,
											ParentMenu:       menu,
											Pkgs:             menu.Pkgs,
											Fset:             menu.Fset,
											MapFile:          menu.MapFile,
										}
										ProcessMenu(subSubMenu)
										menu.Menus = append(menu.Menus, subSubMenu)
										methodIsService = true
									}
								}
							}
							if !methodIsService {
								for _, name := range method.Names {
									fmt.Printf("%sMethod:%s\n", Ident(ParentMenuCount(menu, nil)+1), name)
									methodItem := &config.Method{
										Description:     "",
										LongDescription: "",
										Name:            name.Name,
										Parameters:      []config.Parameter{},
										Returns:         []config.Parameter{},
										Comments:        "",
										Confirmation:    nil,
										IsService:       false,
										ServiceImport:   "",
										SDKFile:         filePath,
										CustomFile:      "",
									}

									params := RetrieveParameters(method.Type.(*ast.FuncType), menu.Pkgs)
									if params != nil {
										methodItem.Parameters = *params
									}

									returns := RetrieveReturns(method.Type.(*ast.FuncType), menu.Pkgs)
									if returns != nil {
										methodItem.Returns = *returns
									}

									menu.Methods = append(menu.Methods, methodItem)
								}
							}
							doBreak = true
						}
						return !doBreak
					}
				}
			}
			return !doBreak
		})
		if doBreak {
			break
		}
	}
}

func IsServiceFunction(returnType ast.Expr) bool {
	if returnType == nil {
		return false
	}
	if _, ok := returnType.(*ast.Ident); ok {
		if strings.HasSuffix(returnType.(*ast.Ident).Name, "Service") {
			return true
		}
	}
	return false
}

func ListDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar diretório %s: %w", dir, err)
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			dirName := file.Name()
			if slices.Contains(DIRS_TO_SKIP, dirName) {
				continue
			}
			names = append(names, file.Name())
		}
	}
	return names, nil
}

func GetPkgWithFset(dir string) ([]*packages.Package, *token.FileSet, error) {
	fset := token.NewFileSet()
	pkgConfig := &packages.Config{Fset: fset, Dir: dir, Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedDeps | packages.NeedName | packages.NeedFiles | packages.NeedImports}
	pkgs, err := packages.Load(pkgConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("erro ao analisar diretório %s: %w", dir, err)
	}
	return pkgs, fset, nil
}

func IsValidModule(goFiles []string) bool {
	for _, file := range goFiles {
		if strings.HasSuffix(file, "client.go") {
			return true
		}
	}
	return false
}

func RetrieveParameters(funcType *ast.FuncType, pkgs *packages.Package) *[]config.Parameter {
	if funcType == nil {
		return nil
	}

	if funcType.Params == nil {
		return nil
	}
	params := make([]config.Parameter, 0)
	for _, param := range funcType.Params.List {
		returnType, aliasType, isPrimitive, isPointer, isArray, isOptional := RetrieveType(param.Type, param, pkgs)

		params = append(params, config.Parameter{
			Name:            RetrieveName(param.Names),
			Type:            returnType,
			Description:     "",
			IsPrimitive:     isPrimitive,
			IsPointer:       isPointer,
			IsOptional:      isOptional,
			IsArray:         isArray,
			IsPositional:    false,
			PositionalIndex: 0,
			Struct:          analyzeStructType(param.Type, pkgs.Name, pkgs.Dir),
			AliasType:       aliasType,
		})
	}

	return &params
}

func RetrieveReturns(funcType *ast.FuncType, pkgs *packages.Package) *[]config.Parameter {
	if funcType == nil {
		return nil
	}

	if funcType.Results == nil {
		return nil
	}
	results := make([]config.Parameter, 0)
	for _, resultItem := range funcType.Results.List {
		returnType, aliasType, isPrimitive, isPointer, isArray, isOptional := RetrieveType(resultItem.Type, resultItem, pkgs)
		results = append(results, config.Parameter{
			Name:            RetrieveName(resultItem.Names),
			Type:            returnType,
			Description:     "",
			IsPrimitive:     isPrimitive,
			IsPointer:       isPointer,
			IsOptional:      isOptional,
			IsArray:         isArray,
			IsPositional:    false,
			PositionalIndex: 0,
			Struct:          analyzeStructType(resultItem.Type, pkgs.Name, pkgs.Dir),
			AliasType:       aliasType,
		})
	}

	return &results
}

func RetrieveType(expr ast.Expr, param *ast.Field, pkgs *packages.Package) (typeString string, aliasType string, isPrimitive bool, isPointer bool, isArray bool, isOptional bool) {
	typeString, aliasType, isPrimitive = getTypeStringWithPackage(expr, pkgs.Name)
	isPointer = strings.HasPrefix(typeString, "*")
	isArray = strings.HasPrefix(typeString, "[]")
	isOptional = false
	if param.Tag != nil {
		tagValue := extractJSONTag(param.Tag.Value)
		if slices.Contains(tagValue, "omitempty") {
			isOptional = true
		}
	}
	return
}

func RetrieveName(names []*ast.Ident) string {
	if len(names) == 0 {
		return ""
	}
	return names[0].Name
}
