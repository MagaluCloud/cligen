package menu_item

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

//go:embed menu_item.template
var menuItemTemplate string

var menuItemTmpl *template.Template

const (
	genDir             = "base-cli-gen/cmd/gen"
	_IN_DEBUG          = false
	DEBUG_MENU_NAME    = "lbaas"
	DEBUG_SUBMENU_NAME = "NetworkBackends"
	DEBUG_METHOD_NAME  = "Get"
)

func init() {
	if _IN_DEBUG {
		fmt.Printf("\t\t\tWARNING!!! DEBUG MODE ENABLED\n%s\n%s\n%s\n", DEBUG_MENU_NAME, DEBUG_SUBMENU_NAME, DEBUG_METHOD_NAME)
	}
	var err error
	menuItemTmpl, err = template.New("menu_item").Parse(menuItemTemplate)
	if err != nil {
		panic(err)
	}
}

func GenerateMenuItem(cfg *config.Config) error {
	for _, menu := range cfg.Menus {
		if _IN_DEBUG && menu.Name != DEBUG_MENU_NAME {
			continue
		}
		for _, submenu := range menu.Menus {
			if _IN_DEBUG && submenu.Name != DEBUG_SUBMENU_NAME {
				continue
			}
			GenMenuItem(cfg, submenu)
		}
	}

	return nil
}

func GenMenuItem(cfg *config.Config, menu *config.Menu) error {
	for _, method := range menu.Methods {
		if _IN_DEBUG && method.Name != DEBUG_METHOD_NAME {
			continue
		}
		parents := FindParents(cfg, menu.ParentMenuID)

		parents = append(parents, strings.ToLower(menu.Name))
		parentsPath := strings.Join(parents, "/")
		nameFromParent, sdkPackageFromParent := FindSDKPackageFromParents(cfg, menu.ParentMenuID)

		menuItem := NewMenuItem()
		sdkName := fmt.Sprintf("%sSdk", strings.ToLower(nameFromParent))

		moduleServiceName := fmt.Sprintf("%sService", strings.ToLower(menu.Name))
		menuItem.SetServiceParam(fmt.Sprintf("%s %s.%s", moduleServiceName, sdkName, menu.ServiceInterface))
		menuItem.AddImport(fmt.Sprintf("%s \"%s\"", sdkName, sdkPackageFromParent))
		menuItem.SetPathSaveToFile(filepath.Join(genDir, parentsPath, fmt.Sprintf("%s.go", strings.ToLower(method.Name))))
		menuItem.SetPackageName(strings.ToLower(menu.Name))
		menuItem.SetFunctionName(strutils.FirstUpper(method.Name))
		menuItem.SetUseName(strutils.ToSnakeCasePreserveID(method.Name, "-"))
		menuItem.SetShortDescription(method.Description)
		menuItem.SetLongDescription(method.LongDescription)
		menuItem.AddImport(fmt.Sprintf("%s \"%s\"", "flags", "github.com/magaluCloud/mgccli/cobra_utils/flags"))
		for _, param := range method.Parameters {
			menuItem = ProcessParams(menu, method, menuItem, sdkName, param, []config.Parameter{}...)
		}
		menuItem = printResult(menuItem, method)
		menuItem.SetServiceCall(fmt.Sprintf("%s.%s", moduleServiceName, method.Name))
		menuItem.SetServiceSDKParam(strings.Join(menuItem.GetParams(), ","))
		// menuItem.AddCobraFlagsAssign()
		// menuItem.AddCobraStructInitialize()
		// menuItem.AddCobraArrayParse()
		// menuItem.SetConfirmation()
		// menuItem.SetAssignResult()
		// menuItem.SetPrintResult()
		// menuItem.SetServiceCall()
		// menuItem.AddPositionalArgs()
		// menuItem.SetServiceSDKParam()

		menuItem.Save()

	}
	for _, ssmenu := range menu.Menus {
		GenMenuItem(cfg, ssmenu)
	}
	return nil
}
func addPrintError() string {
	return "\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\t"
}
func genResponseNameFromType(typ string, i int) string {
	if typ == "error" {
		return "err"
	}
	return fmt.Sprintf("result%d", i)
}
func printResult(menuItem MenuItem, method *config.Method) MenuItem {
	assignResult := make([]string, 0, len(method.Returns))
	printRst := make([]string, 0, len(method.Returns))
	hasNonError := false

	for i, response := range method.Returns {
		if response.Name == "" {
			response.Name = genResponseNameFromType(response.Type, i)
		}
		assignResult = append(assignResult, response.Name)

		if response.Type == "error" {
			printRst = append(printRst, addPrintError())
			continue
		}

		hasNonError = true
	}

	if hasNonError {
		printRst = append(printRst, "\t\t\traw, _ := cmd.Root().PersistentFlags().GetBool(\"raw\")")
		menuItem.AddImport("\"github.com/magaluCloud/mgccli/beautiful\"")

		for i, response := range method.Returns {
			if response.Type != "error" {
				if response.Name == "" {
					response.Name = genResponseNameFromType(response.Type, i)
				}
				printRst = append(printRst, fmt.Sprintf("\t\t\tbeautiful.NewOutput(raw).PrintData(%s)", response.Name))
			}
		}
	}

	menuItem.SetAssignResult(strings.Join(assignResult, ", "))
	menuItem.SetPrintResult(strings.Join(printRst, "\n"))

	return menuItem
}

func ProcessParams(menu *config.Menu, method *config.Method, menuItem MenuItem, sdkName string, param config.Parameter, parents ...config.Parameter) MenuItem {
	if param.Name == "ctx" {
		return menuItem
	}

	if len(param.Struct) == 0 && param.IsPrimitive {
		menuItem = ProcessPrimitiveParam(menu, method, menuItem, sdkName, param, parents...)
		if len(parents) == 0 {
			ptr := ""
			if param.IsPointer {
				ptr = "&"
			}
			menuItem.AddParam(ptr + param.Name)
		}
	}
	if len(param.Struct) > 0 {
		for _, sparam := range param.Struct {
			menuItem = ProcessParams(menu, method, menuItem, sdkName, sparam, append(parents, param)...)
		}
		ptr := ""
		if param.IsPointer {
			ptr = "&"
		}

		if len(parents) == 0 {
			menuItem.AddParam(ptr + param.Name)
			menuItem.AddServiceSDKParamCreate(ServiceSDKParamCreate(param, sdkName))
		}
	}

	return menuItem
}

func parentsNamesDefinition(parents []config.Parameter) []string {
	names := make([]string, len(parents))
	for i, parent := range parents {
		names[i] = parent.Name
	}
	return names
}

func ProcessPrimitiveParam(menu *config.Menu, method *config.Method, menuItem MenuItem, sdkName string, param config.Parameter, parents ...config.Parameter) MenuItem {
	// Definition
	varName := fmt.Sprintf("%sFlag", param.Name)
	if len(parents) > 0 {
		varName = fmt.Sprintf("%s_%s", strings.Join(parentsNamesDefinition(parents), "_"), varName)
	}
	menuItem.AddCobraFlagsDefinition(fmt.Sprintf("var %s *flags.%s", varName, translateTypeToCobraFlag(param.Type)))

	if len(parents) == 0 && !param.IsOptional && !param.IsArray {
		param.Description = fmt.Sprintf("%s (required)", param.Description)
	}

	// FlagCreation
	flagName := strutils.ToSnakeCasePreserveID(param.Name, "-")
	menuItem.AddCobraFlagsCreation(fmt.Sprintf("%s= flags.New%s(cmd, \"%s\", %s,\"%s\")", varName, translateTypeToCobraFlagCreate(param.Type), flagName, defaultByType(param.Type), param.Description))

	// ServiceSDKParamCreate
	if len(parents) == 0 {
		typeName := getParamTypeName(param, sdkName)
		menuItem.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, typeName))
	}

	// CobraFlagsAssign
	if len(parents) > 0 {
		cfa := ""
		if param.IsPrimitive && !param.IsArray {
			varNameCfa := fmt.Sprintf("%s.%s", strings.Join(parentsNamesDefinition(parents), "."), param.Name)
			pointer := ""
			if !param.IsPointer {
				pointer = "*"
			}
			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeA
			%s = %s%s.Value
			}`, varName, varNameCfa, pointer, varName)

		}
		if param.IsPrimitive && param.IsArray {
			varNameCfa := fmt.Sprintf("%s.%s", strings.Join(parentsNamesDefinition(parents), "."), param.Name)

			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeb
			  %s = make([]%s, len(*%s.Value))
			  for i, v := range *%s.Value {
			    %s[i] = %s(v)
			  }
			}`, varName, varNameCfa, param.AliasType, varName, varName, varNameCfa, param.AliasType)
		}
		if !param.IsPrimitive && param.IsArray {
			fmt.Println("asasd")
		}
		menuItem.AddCobraFlagsAssign(cfa)
	}

	if len(parents) == 0 {
		cfa := ""
		if param.IsPrimitive && !param.IsArray && param.IsOptional {
			pointer := ""
			if !param.IsPointer {
				pointer = "*"
			}
			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeC
			%s = %s%s.Value
			}`, varName, flagName, pointer, varName)
		}

		if param.IsPrimitive && !param.IsArray && !param.IsOptional {
			cfa = fmt.Sprintf(`if len(args) > 0{
			cmd.Flags().Set("%s", args[0])
			}`, flagName)

			pointer := ""
			if !param.IsPointer {
				pointer = "*"
			}

			cfa = fmt.Sprintf(`%s
			if %s.IsChanged(){//typeD
			%s = %s%s.Value
			}else{
			return fmt.Errorf("é necessário fornecer o %s como argumento ou usar a flag --%s")
			}`, cfa, varName, param.Name, pointer, varName, param.Name, flagName)

			menuItem.AddImport("\"fmt\"")
		}

		if param.IsPrimitive && param.IsArray {
			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeE
			  %s = make([]%s, len(*%s.Value))
			  for i, v := range *%s.Value {
			    %s[i] = %s(v)
			  }
			}`, varName, flagName, checkHasSdkName(param.AliasType, sdkName), varName, varName, flagName, checkHasSdkName(param.AliasType, sdkName))
		}
		if !param.IsPrimitive && param.IsArray {
			fmt.Println("asasd")
		}
		menuItem.AddCobraFlagsAssign(cfa)
	}

	return menuItem
}

// fix in config gen
func checkHasSdkName(name, sdkName string) string {
	if strings.HasPrefix(name, sdkName) {
		return name
	}
	return fmt.Sprintf("%s.%s", sdkName, name)
}

func ServiceSDKParamCreate(param config.Parameter, sdkName string) string {
	if param.Name == "ctx" {
		return ""
	}

	if param.IsPrimitive && !param.IsArray {
		typeName := getParamTypeName(param, sdkName)
		return fmt.Sprintf("var %s %s", param.Name, typeName)
	}
	if param.IsPrimitive && param.IsArray {
		typeName := getParamTypeName(param, sdkName)
		return fmt.Sprintf("var %s %s", param.Name, typeName)
	}

	if !param.IsPrimitive && !param.IsArray {
		typeName := strings.Replace(param.Type, "*", "", 1)
		return fmt.Sprintf("%s := %s{}", param.Name, typeName)
	}

	return ""
}

func getParamTypeName(param config.Parameter, sdkName string) string {
	if param.AliasType == "" {
		return param.Type
	}
	if param.IsArray {
		if param.AliasType != "" && sdkName != "" {
			return "[]" + sdkName + "." + param.AliasType
		}
		return "[]" + param.AliasType
	}
	return param.AliasType
}

func defaultByType(paramType string) string {
	paramType = strings.TrimPrefix(paramType, "*")
	switch paramType {
	case "string":
		return "\"\""
	case "int64", "int32", "int16", "int8", "int", "float64":
		return "0"
	case "bool":
		return "false"
	case "[]string":
		return "[]string{}"
	case "map[string]string":
		return "map[string]string{}"
	default:
		return "\"\""
	}
}

func translateTypeToCobraFlagComplex(field config.Parameter) string {
	typeName := cleanTypeName(field.Type)
	if field.IsArray {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	return fmt.Sprintf("JSONValue[%s]", typeName)
}

func cleanTypeName(typeName string) string {
	typeName = strings.TrimPrefix(typeName, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	return typeName
}

func translateTypeToCobraFlagStruct(field, parentField config.Parameter) string {
	typeName := cleanTypeName(parentField.Type)
	if canUseSliceFlag(parentField) {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	if canUseStrAsJson(parentField) {
		return fmt.Sprintf("JSONValue[%s]", typeName)
	}
	return translateTypeToCobraFlag(field.Type)
}

func translateTypeToCobraFlag(paramType string) string {

	paramType = strings.TrimPrefix(paramType, "*")
	switch paramType {
	case "string":
		return "StrFlag"
	case "int64", "int32", "int16", "int8":
		return "Int64Flag"
	case "bool":
		return "BoolFlag"
	case "int":
		return "IntFlag"
	case "float64":
		return "Float64Flag"
	case "[]string":
		return "StrSliceFlag"
	case "map[string]string":
		return "StrMapFlag"
	case "time.Time":
		return "TimeFlag"
	default:
		return "StrFlag"
	}
}

func translateTypeToCobraFlagCreateComplex(field config.Parameter) string {
	typeName := cleanTypeName(field.Type)
	if field.IsArray {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	return fmt.Sprintf("JSONValue[%s]", typeName)
}

func translateTypeToCobraFlagCreateStruct(field, parentField config.Parameter) string {
	typeName := cleanTypeName(parentField.Type)
	if canUseSliceFlag(parentField) {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	if canUseStrAsJson(parentField) {
		return fmt.Sprintf("JSONValue[%s]", typeName)
	}
	return translateTypeToCobraFlag(field.Type)
}

func translateTypeToCobraFlagCreate(paramType string) string {
	paramType = strings.TrimPrefix(paramType, "*")

	switch paramType {
	case "string":
		return "Str"
	case "int64", "int32", "int16", "int8":
		return "Int64"
	case "bool":
		return "Bool"
	case "int":
		return "Int"
	case "float64":
		return "Float64"
	case "[]string":
		return "StrSlice"
	case "map[string]string":
		return "StrMap"
	case "time.Time":
		return "Time"
	default:

		return "Str"
	}
}

func canUseSliceFlag(parentField config.Parameter) bool {
	return parentField.IsArray && len(parentField.Struct) == 1
}

func canUseStrAsJson(parentField config.Parameter) bool {
	return parentField.IsArray && len(parentField.Struct) > 1
}

// mover isso pra uma common
func FindParents(cfg *config.Config, menuID string) []string {
	parents := []string{}
	menu := FindMenuByID(cfg.Menus, menuID)
	if menu != nil {
		parents = append(parents, strings.ToLower(menu.Name))
		if menu.ParentMenuID != "" {
			parents = append(FindParents(cfg, menu.ParentMenuID), parents...)
		}
	}
	return parents
}

func FindMenuByID(menus []*config.Menu, id string) *config.Menu {
	for _, menu := range menus {
		if menu.ID == id {
			return menu
		}
		if len(menu.Menus) > 0 {
			menu := FindMenuByID(menu.Menus, id)
			if menu != nil {
				return menu
			}
		}
	}
	return nil
}

func FindSDKPackageFromParents(cfg *config.Config, menuID string) (name string, sdkPackage string) {
	menu := FindMenuByID(cfg.Menus, menuID)
	if menu == nil {
		return "", ""
	}
	if menu.SDKPackage != "" {
		return menu.Name, menu.SDKPackage
	}
	if menu.ParentMenuID != "" {
		return FindSDKPackageFromParents(cfg, menu.ParentMenuID)
	}
	return "", ""
}
