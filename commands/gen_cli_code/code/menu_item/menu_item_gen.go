package menu_item

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

//go:embed menu_item.template
var menuItemTemplate string

var menuItemTmpl *template.Template

var _IN_DEBUG bool

const (
	genDir             = "base-cli-gen/cmd/gen"
	DEBUG_MENU_NAME    = "lbaas"
	DEBUG_SUBMENU_NAME = "NetworkACLs"
	DEBUG_METHOD_NAME  = "Replace"
)

func init() {
	_IN_DEBUG = os.Getenv("IN_DEBUG") != ""

	if _IN_DEBUG {
		fmt.Printf("\t\t\tWARNING!!! DEBUG MODE ENABLED\n\t\t\t%s\n\t\t\t%s\n\t\t\t%s\n", DEBUG_MENU_NAME, DEBUG_SUBMENU_NAME, DEBUG_METHOD_NAME)
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

		for _, param := range method.Parameters {
			if param.Name == "ctx" {
				menuItem.AddParam("ctx")
				continue
			}
			menuItem = ProcessCobraFlagsDefinition(menuItem, param)
			menuItem = ProcessServiceSDKParamCreate(menuItem, param, sdkName)
			menuItem = ProcessCobraStructInitialize(menuItem, param, sdkName)
		}
		menuItem = ProcessCobraFlagsAssign(menuItem, sdkName)
		menuItem = ProcessCobraFlagsCreation(menuItem, sdkName)
		assignResult, menuItem := ProcessAssignResult(menuItem, method, sdkName)
		menuItem = ProcessPrintResult(menuItem, assignResult)
		menuItem.SetServiceCall(fmt.Sprintf("%s.%s", moduleServiceName, method.Name))
		menuItem.SetServiceSDKParam(strings.Join(menuItem.GetParams(), ","))

		menuItem.Save()

	}
	for _, ssmenu := range menu.Menus {
		GenMenuItem(cfg, ssmenu)
	}
	return nil
}

func prepareName(name string) string {
	name = strings.ReplaceAll(name, "*", "")
	name = strings.ReplaceAll(name, "[]", "")
	return name
}

func ProcessCobraStructInitialize(menuItem MenuItem, param config.Parameter, sdkName string, parents ...config.Parameter) MenuItem {
	parentControl := false
	for _, sparam := range param.Struct {
		if sparam.IsPointer && !sparam.IsPrimitive {
			if !parentControl {
				parentControl = true
				parents = append(parents, param)
			}

			typePrefix := "&"
			if sparam.IsArray {
				typePrefix = "&[]"
			}
			sci := ServiceCobraStructInitialize{
				ParamName:  sparam.Name,
				Name:       genStructName(sparam, parents...),
				TypePrefix: typePrefix,
				TypeSuffix: "{}",
				ParamType:  prepareName(sparam.Type),
			}

			menuItem.AddCobraStructInitialize(sci)
			menuItem = ProcessCobraStructInitialize(menuItem, sparam, sdkName, parents...)
		}
	}
	return menuItem
}
func ProcessCobraFlagsAssign(menuItem MenuItem, sdkName string) MenuItem {
	for _, flag := range menuItem.GetCobraFlagsDefinition() {
		cfa := ""
		if !flag.param.IsOptional {
			cfa = fmt.Sprintf(`if len(args) > 0{
				cmd.Flags().Set("%s", args[0])
			}
			`, flag.cobraVar)

		}
		cfa = fmt.Sprintf("%sif %sFlag.IsChanged(){\n", cfa, flag.Name)

		pointer := ""
		if !flag.param.IsPointer && !flag.param.IsArray {
			pointer = "*"
		}

		if flag.isComplex {
			switch {
			case strings.HasPrefix(flag.FlagType, "JSONArrayValue"):
				if !flag.param.IsPointer {
					pointer = "*"
				}
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case strings.HasPrefix(flag.FlagType, "JSONValue"):
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)

			}
			// cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
		}

		if !flag.isComplex {
			switch flag.FlagType {
			case "StrSliceFlag":
				if flag.param.AliasType != "" {
					// expand = flags.StrSliceFlagToSlice[blockstorageSdk.ExpandSchedulers](expandFlag)
					if hasSDKPackage(flag.param.AliasType, sdkName) {
						cfa = fmt.Sprintf("%s				%s = flags.StrSliceFlagToSlice[%s](%sFlag)\n\n", cfa, flag.cobraAssign, flag.param.AliasType, flag.Name)
					} else {
						cfa = fmt.Sprintf("%s				%s = flags.StrSliceFlagToSlice[%s.%s](%sFlag)\n\n", cfa, flag.cobraAssign, sdkName, flag.param.AliasType, flag.Name)
					}
				} else if hasSDKPackage(flag.param.Type, sdkName) {
					cfa = fmt.Sprintf("%s				%s = flags.StrSliceFlagToSlice[%s](%sFlag)\n\n", cfa, flag.cobraAssign, removeArrayFromString(flag.param.Type), flag.Name)
				} else {
					if !flag.param.IsPointer {
						pointer = "*"
					}
					cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
				}
			case "StrFlag":
				ptrValue := ""
				if flag.param.AliasType != "" {
					if flag.param.IsPointer {
						pointer = "*"
						ptrValue = "&"
					}
					if hasSDKPackage(flag.param.AliasType, sdkName) {
						cfa = fmt.Sprintf("%s				localVar := %s(%s%sFlag.Value)\n\n", cfa, flag.param.AliasType, pointer, flag.Name)
						cfa = fmt.Sprintf("%s				%s = %slocalVar\n\n", cfa, flag.cobraAssign, ptrValue)
					} else {
						cfa = fmt.Sprintf("%s				localVar := %s.%s(%s%sFlag.Value)\n\n", cfa, sdkName, flag.param.AliasType, pointer, flag.Name)
						cfa = fmt.Sprintf("%s				%s = %slocalVar\n\n", cfa, flag.cobraAssign, ptrValue)
					}
				} else if hasSDKPackage(flag.param.Type, sdkName) {
					if flag.param.IsPointer {
						pointer = "*"
					}
					cfa = fmt.Sprintf("%s				localVar := %s(%s%sFlag.Value)\n\n", cfa, flag.param.Type, pointer, flag.Name)
					cfa = fmt.Sprintf("%s				%s = %slocalVar\n\n", cfa, flag.cobraAssign, ptrValue)
				} else {
					cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
				}
			case "Int64Flag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case "BoolFlag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case "IntFlag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case "Float64Flag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case "StrMapFlag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			case "TimeFlag":
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			default:
				cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
			}
		}
		// aqui ruim
		// if flag.isComplex {
		// 	cfa = fmt.Sprintf("%s				%s = %sFlag.Value\n\n", cfa, flag.cobraAssign, flag.Name)
		// } else {

		// 	if !flag.param.IsArray && !flag.parentIsArray && flag.param.IsPrimitive && flag.param.AliasType == "" {
		// 		pointer := ""
		// 		if !flag.param.IsPointer {
		// 			pointer = "*"
		// 		}
		// 		cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
		// 	}
		// 	if !flag.param.IsArray && !flag.parentIsArray && flag.param.IsPrimitive && flag.param.AliasType != "" {
		// 		pointer := ""
		// 		if flag.param.IsPointer {
		// 			pointer = "*"
		// 		}
		// 		cfa = fmt.Sprintf("%s				%s%s = %s(*%sFlag.Value)\n\n", cfa, pointer, flag.cobraAssign, flag.param.AliasType, flag.Name)
		// 	}
		// 	if flag.param.IsArray && !flag.parentIsArray {

		// 		pointer := ""
		// 		if flag.param.IsPointer {
		// 			pointer = "*"
		// 		}

		// 		if flag.param.IsPrimitive && flag.param.AliasType == "" {
		// 			cfa = fmt.Sprintf("%s				for _, v := range *%sFlag.Value {\n", cfa, flag.Name)
		// 			cfa = fmt.Sprintf("%s					%s%s = append(%s%s, v)//dd\n", cfa, pointer, flag.cobraAssign, pointer, flag.cobraAssign)
		// 			cfa = fmt.Sprintf("%s				}\n", cfa)
		// 		}
		// 		if !flag.param.IsPrimitive && flag.param.AliasType != "" {
		// 			cfa = fmt.Sprintf("%sfor _, v := range *%sFlag.Value {\n", cfa, flag.Name)
		// 			if !hasSDKPackage(flag.param.AliasType, sdkName) {
		// 				cfa = fmt.Sprintf("%s%s = append(%s, %s.%s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, sdkName, flag.param.AliasType)
		// 			} else {
		// 				cfa = fmt.Sprintf("%s%s = append(%s, %s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, flag.param.AliasType)
		// 			}
		// 			cfa = fmt.Sprintf("%s}\n", cfa)
		// 		}
		// 		if !flag.param.IsPrimitive && flag.param.AliasType == "" {
		// 			cfa = fmt.Sprintf("%sfor _, v := range *%sFlag.Value {\n", cfa, flag.Name)
		// 			if !hasSDKPackage(flag.param.Type, sdkName) {
		// 				cfa = fmt.Sprintf("%s%s = append(%s, %s.%s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, sdkName, flag.param.Type)
		// 			} else {
		// 				cfa = fmt.Sprintf("%s%s = append(%s, %s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, removeArrayFromString(flag.param.Type))
		// 			}
		// 			cfa = fmt.Sprintf("%s}\n", cfa)
		// 		}
		// 	}

		// 	if !flag.param.IsArray && flag.parentIsArray {
		// 		if flag.param.IsPrimitive && flag.param.AliasType == "" {
		// 			if flag.parentIsStruct && flag.arrayMake {
		// 				csi := menuItem.GetCobraStructInitialize()
		// 				for _, csi := range csi {
		// 					if csi.Name == removeSuffix(flag.cobraAssign, flag.param.Name) {
		// 						cfa = fmt.Sprintf("%s				for _, v := range *%sFlag.Value {\n", cfa, flag.Name)
		// 						cfa = fmt.Sprintf("%s					*%s = append(*%s,%s{%s: v})//kk\n", cfa, removeSuffix(flag.cobraAssign, flag.param.Name), removeSuffix(flag.cobraAssign, flag.param.Name), csi.ParamType, flag.ParamName)
		// 						cfa = fmt.Sprintf("%s				}\n", cfa)
		// 						break
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// }

		// aqui bom:
		if !flag.param.IsOptional {
			cfa = fmt.Sprintf(`%s			} else {
				return fmt.Errorf("é necessário fornecer o %s como argumento ou usar a flag --%s")
			`, cfa, flag.Name, flag.param.Name)
			menuItem.AddImport("\"fmt\"")

		}

		cfa = fmt.Sprintf("%s			}\n", cfa)
		menuItem.AddCobraFlagsAssign(cfa)
	}
	return menuItem
}

func removeSuffix(str string, suffix string) string {
	return strings.TrimSuffix(str, "."+suffix)
}

func ProcessPrintResult(menuItem MenuItem, assignResult string) MenuItem {
	var printRst []string
	printRst = append(printRst, "raw, _ := cmd.Root().PersistentFlags().GetBool(\"raw\")")
	printRst = append(printRst, fmt.Sprintf("\t\t\tbeautiful.NewOutput(raw).PrintData(%s)", assignResult))
	menuItem.SetPrintResult(strings.Join(printRst, "\n"))
	menuItem.AddImport("\"github.com/magaluCloud/mgccli/beautiful\"")
	return menuItem
}

func ProcessAssignResult(menuItem MenuItem, method *config.Method, sdkName string) (string, MenuItem) {
	assignResult := make([]string, 0, len(method.Returns))
	for i, response := range method.Returns {
		if response.Name == "" || response.Name == "error" || response.Type == "error" {
			assignResult = append(assignResult, "err")
			if response.Type == "error" {
				menuItem.SetErrorResult("\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\t")
			}
			continue
		}
		assignResult = append(assignResult, fmt.Sprintf("result%d", i))
	}
	menuItem.SetAssignResult(strings.Join(assignResult, ", "))
	return assignResult[0], menuItem
}

func genStructName(param config.Parameter, parents ...config.Parameter) string {
	if len(parents) == 0 {
		return param.Name
	}
	parts := make([]string, 0, len(parents)+1)
	for _, parent := range parents {
		parts = append(parts, parent.Name)
	}
	parts = append(parts, param.Name)
	return strings.Join(parts, ".")
}

func genFlagName(param config.Parameter, parents []config.Parameter) string {
	if len(parents) == 0 {
		return param.Name
	}
	parts := make([]string, 0, len(parents)+1)
	for _, parent := range parents {
		parts = append(parts, parent.Name)
	}
	parts = append(parts, param.Name)
	return strings.Join(parts, "_")
}

func getParentsNames(parents []config.Parameter) []string {
	parts := make([]string, 0, len(parents))
	for _, parent := range parents {
		parts = append(parts, parent.Name)
	}
	return parts
}

func genCobraVarName(param config.Parameter, parents []config.Parameter) string {
	if len(parents) == 0 {
		return strutils.ToSnakeCasePreserveID(param.Name, "-")
	}
	par := parents[1:]
	if len(par) == 0 {
		return strutils.ToSnakeCasePreserveID(param.Name, "-")
	}
	parts := make([]string, 0, len(par))
	for _, parent := range par {
		parts = append(parts, strutils.ToSnakeCasePreserveID(parent.Name, "-"))
	}
	parts = append(parts, strutils.ToSnakeCasePreserveID(param.Name, "-"))

	return strings.Join(parts, ".")
}

func genCobraAssignName(param config.Parameter, parents []config.Parameter) string {
	if len(parents) == 0 {
		return param.Name
	}
	parts := make([]string, 0, len(parents))
	for _, parent := range parents {
		parts = append(parts, parent.Name)
	}
	parts = append(parts, param.Name)

	return strings.Join(parts, ".")
}

func ProcessCobraFlagsDefinition(menuItem MenuItem, param config.Parameter, parents ...config.Parameter) MenuItem {
	if param.Struct == nil {
		cfd := CobraFlagsDefinition{
			ParamName:         param.Name,
			Name:              genFlagName(param, parents),
			FlagType:          translateTypeToCobraFlag(param, parents),
			param:             param,
			parents:           getParentsNames(parents),
			cobraVar:          genCobraVarName(param, parents),
			cobraAssign:       genCobraAssignName(param, parents),
			cobraType:         translateTypeToCobraFlagCreate(param, parents),
			cobraDefaultValue: defaultByType(param, parents),
			parentIsArray:     len(parents) > 0 && parents[len(parents)-1].IsArray,
			parentIsStruct:    len(parents) > 0 && len(parents[len(parents)-1].Struct) > 0,
			arrayMake:         len(parents) > 0 && parents[len(parents)-1].IsArray && strings.HasPrefix(parents[len(parents)-1].Type, "*[]"),
			isComplex:         false,
		}
		menuItem.AddCobraFlagsDefinition(cfd)
		return menuItem
	}

	parents = append(parents, param)
	for _, sparam := range param.Struct {
		if sparam.Struct != nil {
			cfd := CobraFlagsDefinition{
				ParamName:         sparam.Name,
				Name:              genFlagName(sparam, parents),
				FlagType:          complexType(sparam),
				param:             sparam,
				parents:           getParentsNames(parents),
				cobraVar:          genCobraVarName(sparam, parents),
				cobraAssign:       genCobraAssignName(sparam, parents),
				cobraType:         complexType(sparam),
				cobraDefaultValue: "",
				parentIsArray:     len(parents) > 0 && parents[len(parents)-1].IsArray,
				parentIsStruct:    len(parents) > 0 && len(parents[len(parents)-1].Struct) > 0,
				arrayMake:         len(parents) > 0 && parents[len(parents)-1].IsArray && strings.HasPrefix(parents[len(parents)-1].Type, "*[]"),
				isComplex:         true,
			}
			menuItem.AddCobraFlagsDefinition(cfd)
			continue
		}
		menuItem = ProcessCobraFlagsDefinition(menuItem, sparam, parents...)
	}

	return menuItem
}

func removeArrayFromString(str string) string {
	cprefix, _ := strings.CutPrefix(str, "[]")
	return cprefix
}

func hasSDKPackage(paramType string, sdkName string) bool {
	return strings.Contains(removeArrayFromString(paramType), sdkName)
}

func ProcessServiceSDKParamCreate(menuItem MenuItem, param config.Parameter, sdkName string, parents ...config.Parameter) MenuItem {
	var sspc ServiceSDKParamCreate
	sspc.ParamName = param.Name
	sspc.Name = param.Name

	if param.IsPrimitive && param.AliasType == "" {
		sspc.ParamType = param.Type
	}
	if param.IsPrimitive && param.AliasType != "" {
		if !hasSDKPackage(param.AliasType, sdkName) {
			sspc.ParamType = fmt.Sprintf("%s.%s", sdkName, param.AliasType)
		} else {
			sspc.ParamType = param.AliasType
		}
	}
	if !param.IsPrimitive && param.Struct != nil {
		sspc.ParamType = param.Type
	}
	if !param.IsPrimitive && param.Struct == nil {
		sspc.ParamType = param.Type
	}

	if param.IsArray && !strings.HasPrefix(sspc.ParamType, "[]") {
		sspc.ParamType = fmt.Sprintf("[]%s", sspc.ParamType)
	}
	menuItem.AddServiceSDKParamCreate(sspc)
	menuItem.AddParam(sspc.Name)
	return menuItem
}

func removeNewLine(str string) string {
	return strings.ReplaceAll(str, "\n", "")
}

func ProcessCobraFlagsCreation(menuItem MenuItem, sdkName string) MenuItem {
	doImport := false
	for _, flag := range menuItem.GetCobraFlagsDefinition() {
		if flag.cobraDefaultValue != "" {
			menuItem.AddCobraFlagsCreation(
				fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", %s,\"%s\")", flag.Name, flag.cobraType, flag.cobraVar, flag.cobraDefaultValue, removeNewLine(flag.param.Description)),
			)
		} else {
			menuItem.AddCobraFlagsCreation(
				fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\")", flag.Name, flag.cobraType, flag.cobraVar, removeNewLine(flag.param.Description)),
			)
		}
		doImport = true
	}
	if doImport {
		menuItem.AddImport("flags \"github.com/magaluCloud/mgccli/cobra_utils/flags\"")
	}
	return menuItem
}

func defaultByType(param config.Parameter, parents []config.Parameter) string {
	paramType := strings.TrimPrefix(param.Type, "*")
	if param.IsArray {
		return "[]string{}"
	}
	if len(parents) > 0 {
		if parents[len(parents)-1].IsArray {
			return "[]string{}"
		}
	}
	switch paramType {
	case "string":
		return "\"\""
	case "int64", "int32", "int16", "int8", "int", "float64":
		return "0"
	case "bool":
		return "false"
	case "[]string":
		return "[]string{}"
	case "[]int64":
		return "[]int64{}"
	case "map[string]string":
		return "map[string]string{}"
	default:
		return "\"\""
	}
}

func complexType(param config.Parameter) string {
	if param.IsArray {
		return fmt.Sprintf("JSONArrayValue[%s]", cleanTypeName(param.Type))
	}
	if !param.IsArray {
		return fmt.Sprintf("JSONValue[%s]", cleanTypeName(param.Type))
	}
	return "StrFlag"
}

func translateTypeToCobraFlag(param config.Parameter, parents []config.Parameter) string {
	if param.IsArray {
		return "StrSliceFlag"
	}

	if len(parents) > 0 {
		if parents[len(parents)-1].IsArray {
			return "StrSliceFlag"
		}
	}

	paramType := strings.ReplaceAll(param.Type, "*", "")
	switch paramType {
	case "string":
		if param.IsArray {
			return "StrSliceFlag"
		}
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

func translateTypeToCobraFlagCreate(param config.Parameter, parents []config.Parameter) string {
	paramType := strings.TrimPrefix(param.Type, "*")

	if !param.IsPrimitive {
		if param.IsArray {
			return "StrSlice"
		}
		return "Str"
	}

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

func cleanTypeName(typeName string) string {
	typeName = strings.TrimPrefix(typeName, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	return typeName
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
