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
	DEBUG_MENU_NAME    = "compute"
	DEBUG_SUBMENU_NAME = "Instances"
	DEBUG_METHOD_NAME  = "Create"
)

func init() {
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
		menuItem.AddImport("flags \"github.com/magaluCloud/mgccli/cobra_utils/flags\"")
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
			// menuItem = ProcessCobraFlagsCreation(menuItem, param, sdkName)
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

		if !flag.param.IsArray && !flag.parentIsArray {
			pointer := ""
			if !flag.param.IsPointer {
				pointer = "*"
			}
			cfa = fmt.Sprintf("%s				%s = %s%sFlag.Value\n\n", cfa, flag.cobraAssign, pointer, flag.Name)
		}

		if flag.param.IsArray && !flag.parentIsArray {

			if flag.param.IsPrimitive && flag.param.AliasType == "" {
				cfa = fmt.Sprintf("%s				for _, v := range *%sFlag.Value {\n", cfa, flag.Name)
				cfa = fmt.Sprintf("%s					*%s = append(*%s, v)//dd\n", cfa, flag.cobraAssign, flag.cobraAssign)
				cfa = fmt.Sprintf("%s				}\n", cfa)
			} else { //aqui ok!
				cfa = fmt.Sprintf("%sfor _, v := range *%sFlag.Value {\n", cfa, flag.Name)
				if !hasSDKPackage(flag.param.AliasType, sdkName) {
					cfa = fmt.Sprintf("%s%s = append(%s, %s.%s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, sdkName, flag.param.AliasType)
				} else {
					cfa = fmt.Sprintf("%s%s = append(%s, %s(v))\n", cfa, flag.cobraAssign, flag.cobraAssign, flag.param.AliasType)
				}
				cfa = fmt.Sprintf("%s}\n", cfa)
			}
		}

		if !flag.param.IsArray && flag.parentIsArray {
			if flag.param.IsPrimitive && flag.param.AliasType == "" {
				if flag.parentIsStruct && flag.arrayMake {
					csi := menuItem.GetCobraStructInitialize()
					for _, csi := range csi {
						if csi.Name == removeSuffix(flag.cobraAssign, flag.param.Name) {
							cfa = fmt.Sprintf("%s				for _, v := range *%sFlag.Value {\n", cfa, flag.Name)
							cfa = fmt.Sprintf("%s					*%s = append(*%s,%s{%s: v})//kk\n", cfa, removeSuffix(flag.cobraAssign, flag.param.Name), removeSuffix(flag.cobraAssign, flag.param.Name), csi.ParamType, flag.ParamName)
							cfa = fmt.Sprintf("%s				}\n", cfa)
							break
						}
					}
				}
			}
		}

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
		}
		menuItem.AddCobraFlagsDefinition(cfd)
	} else {
		for _, sparam := range param.Struct {
			menuItem = ProcessCobraFlagsDefinition(menuItem, sparam, append(parents, param)...)
		}
	}

	return menuItem
}

func hasSDKPackage(paramType string, sdkName string) bool {
	return strings.Contains(paramType, sdkName)
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
	for _, flag := range menuItem.GetCobraFlagsDefinition() {
		menuItem.AddCobraFlagsCreation(
			fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", %s,\"%s\")", flag.Name, flag.cobraType, flag.cobraVar, flag.cobraDefaultValue, removeNewLine(flag.param.Description)),
		)
	}
	return menuItem
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//

// anterior

// func addPrintError() string {
// 	return "\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\t"
// }
// func genResponseNameFromType(typ string, i int) string {
// 	if typ == "error" {
// 		return "err"
// 	}
// 	return fmt.Sprintf("result%d", i)
// }

// func printResult(menuItem MenuItem, method *config.Method) MenuItem {
// 	assignResult := make([]string, 0, len(method.Returns))
// 	printRst := make([]string, 0, len(method.Returns))
// 	hasNonError := false

// 	for i, response := range method.Returns {
// 		if response.Name == "" || response.Name == "error" {
// 			response.Name = genResponseNameFromType(response.Type, i)
// 		}
// 		response.Name = strings.ReplaceAll(response.Name, "*", "")
// 		response.Name = strings.ReplaceAll(response.Name, ".", "")
// 		response.Name = strings.ReplaceAll(response.Name, "[]", "")
// 		response.Name = strings.ReplaceAll(response.Name, "_", "")
// 		assignResult = append(assignResult, response.Name)

// 		if response.Type == "error" {
// 			printRst = append(printRst, addPrintError())
// 			continue
// 		}

// 		hasNonError = true
// 	}

// 	if hasNonError {
// 		printRst = append(printRst, "\t\t\traw, _ := cmd.Root().PersistentFlags().GetBool(\"raw\")")
// 		menuItem.AddImport("\"github.com/magaluCloud/mgccli/beautiful\"")

// 		for i, response := range method.Returns {
// 			if response.Type != "error" {
// 				if response.Name == "" || response.Name == "error" {
// 					response.Name = genResponseNameFromType(response.Type, i)
// 				}
// 				response.Name = strings.ReplaceAll(response.Name, "*", "")
// 				response.Name = strings.ReplaceAll(response.Name, ".", "")
// 				response.Name = strings.ReplaceAll(response.Name, "[]", "")
// 				response.Name = strings.ReplaceAll(response.Name, "_", "")
// 				printRst = append(printRst, fmt.Sprintf("\t\t\tbeautiful.NewOutput(raw).PrintData(%s)", response.Name))
// 			}
// 		}
// 	}

// 	menuItem.SetAssignResult(strings.Join(assignResult, ", "))
// 	menuItem.SetPrintResult(strings.Join(printRst, "\n"))

// 	return menuItem
// }

// func ProcessParams(menu *config.Menu, method *config.Method, menuItem MenuItem, sdkName string, param config.Parameter, parents ...config.Parameter) MenuItem {
// 	if param.Name == "ctx" {
// 		return menuItem
// 	}

// 	if len(param.Struct) == 0 && param.IsPrimitive {
// 		menuItem = ProcessPrimitiveParam(menu, method, menuItem, sdkName, param, parents...)
// 		if len(parents) == 0 {
// 			ptr := ""
// 			if param.IsPointer {
// 				ptr = "&"
// 			}
// 			menuItem.AddParam(ptr + param.Name)
// 		}
// 		return menuItem
// 	}
// 	if len(param.Struct) == 0 && param.IsArray {
// 		menuItem = ProcessPrimitiveParam(menu, method, menuItem, sdkName, param, parents...)
// 		menuItem.AddParam(param.Name)
// 		return menuItem
// 	}
// 	if len(param.Struct) > 0 {
// 		for _, sparam := range param.Struct {
// 			menuItem = ProcessParams(menu, method, menuItem, sdkName, sparam, append(parents, param)...)
// 		}
// 		ptr := ""
// 		if param.IsPointer {
// 			ptr = "&"
// 		}

// 		if len(parents) == 0 {
// 			menuItem.AddParam(ptr + param.Name)
// 			menuItem.AddServiceSDKParamCreate(ServiceSDKParamCreate(param, sdkName))
// 		}
// 	}

// 	return menuItem
// }

// func parentsNamesDefinition(parents []config.Parameter) []string {
// 	names := make([]string, len(parents))
// 	for i, parent := range parents {
// 		names[i] = parent.Name
// 	}
// 	return names
// }

// func processParamDescription(param config.Parameter) string {
// 	if param.Description == "" {
// 		return ""
// 	}
// 	result := param.Description
// 	result = strings.ReplaceAll(result, "\n", "")
// 	result = strings.ReplaceAll(result, "\t", "")
// 	return result
// }

// func ProcessPrimitiveParam(menu *config.Menu, method *config.Method, menuItem MenuItem, sdkName string, param config.Parameter, parents ...config.Parameter) MenuItem {
// 	// Definition
// 	varName := fmt.Sprintf("%sFlag", param.Name)
// 	if len(parents) > 0 {
// 		varName = fmt.Sprintf("%s_%s", strings.Join(parentsNamesDefinition(parents), "_"), varName)
// 	}
// 	menuItem.AddCobraFlagsDefinition(fmt.Sprintf("var %s *flags.%s", varName, translateTypeToCobraFlag(param.Type)))

// 	if len(parents) == 0 && !param.IsOptional && !param.IsArray {
// 		param.Description = fmt.Sprintf("%s (required)", param.Description)
// 	}

// 	// FlagCreation
// 	flagName := strutils.ToSnakeCasePreserveID(param.Name, "-")
// 	menuItem.AddCobraFlagsCreation(fmt.Sprintf("%s= flags.New%s(cmd, \"%s\", %s,\"%s\")", varName, translateTypeToCobraFlagCreate(param.Type), flagName, defaultByType(param.Type), processParamDescription(param)))
// 	menuItem.AddImport(fmt.Sprintf("%s \"%s\"", "flags", "github.com/magaluCloud/mgccli/cobra_utils/flags"))

// 	// ServiceSDKParamCreate
// 	if len(parents) == 0 {
// 		typeName := getParamTypeName(param, sdkName)
// 		menuItem.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, typeName))
// 	}

// 	// CobraFlagsAssign
// 	if len(parents) > 0 {
// 		cfa := ""
// 		if param.IsPrimitive && !param.IsArray {
// 			varNameCfa := fmt.Sprintf("%s.%s", strings.Join(parentsNamesDefinition(parents), "."), param.Name)
// 			pointer := ""
// 			if !param.IsPointer {
// 				pointer = "*"
// 			}
// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeA
// 			%s = %s%s.Value
// 			}`, varName, varNameCfa, pointer, varName)

// 		}
// 		if param.IsPrimitive && param.IsArray {
// 			varNameCfa := fmt.Sprintf("%s.%s", strings.Join(parentsNamesDefinition(parents), "."), param.Name)

// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeb
// 			  %s = make([]%s, len(*%s.Value))
// 			  for i, v := range *%s.Value {
// 			    %s[i] = %s(v)
// 			  }
// 			}`, varName, varNameCfa, processAliasType(param, sdkName), varName, varName, varNameCfa, param.AliasType)
// 		}
// 		if !param.IsPrimitive && param.IsArray {
// 			fmt.Println("asasd")
// 		}
// 		menuItem.AddCobraFlagsAssign(cfa)
// 	}

// 	if len(parents) == 0 {
// 		cfa := ""
// 		if param.IsPrimitive && !param.IsArray && param.IsOptional {
// 			pointer := ""
// 			if !param.IsPointer {
// 				pointer = "*"
// 			}
// 			if strings.Contains(flagName, "-") {
// 				fmt.Println("asdkfasdf")
// 			}
// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeC
// 			%s = %s%s.Value
// 			}`, varName, param.Name, pointer, varName)
// 		}

// 		if param.IsPrimitive && !param.IsArray && !param.IsOptional {
// 			cfa = fmt.Sprintf(`if len(args) > 0{
// 			cmd.Flags().Set("%s", args[0])
// 			}`, flagName)

// 			pointer := ""
// 			if !param.IsPointer {
// 				pointer = "*"
// 			}

// 			cfa = fmt.Sprintf(`%s
// 			if %s.IsChanged(){//typeD
// 			%s = %s%s.Value
// 			}else{
// 			return fmt.Errorf("é necessário fornecer o %s como argumento ou usar a flag --%s")
// 			}`, cfa, varName, param.Name, pointer, varName, param.Name, flagName)

// 			menuItem.AddImport("\"fmt\"")
// 		}

// 		if param.IsPrimitive && param.IsArray && param.AliasType == "" {
// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeE
// 			  %s = make([]%s, len(*%s.Value))
// 			  for i, v := range *%s.Value {
// 			    %s[i] = v
// 			  }
// 			}`, varName, flagName, processStringType(param.Type), varName, varName, flagName)
// 		}
// 		if param.IsPrimitive && param.IsArray && param.AliasType != "" {
// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeF
// 			  %s = make([]%s, len(*%s.Value))
// 			  for i, v := range *%s.Value {
// 			    %s[i] = %s(v)
// 			  }
// 			}`, varName, flagName, processAliasType(param, sdkName), varName, varName, flagName, processAliasType(param, sdkName))
// 		}
// 		if !param.IsPrimitive && param.IsArray {
// 			cfa = fmt.Sprintf(`if %s.IsChanged(){//typeG
// 			%s = make([]%s, len(*%s.Value))
// 			for i, v := range *%s.Value {
// 			  %s[i] = v
// 			}
// 		  }`, varName, flagName, checkHasSdkName(param, sdkName), varName, varName, flagName)
// 		}
// 		menuItem.AddCobraFlagsAssign(cfa)
// 	}

// 	return menuItem
// }

// func processAliasType(param config.Parameter, sdkName string) string {
// 	if param.AliasType != "" {
// 		return checkHasSdkName(param, sdkName)
// 	}

// 	return processStringType(param.Type)
// }

// func processStringType(t string) string {
// 	result := strings.ReplaceAll(t, "*", "")
// 	result = strings.ReplaceAll(result, "[]", "")
// 	result = strings.ReplaceAll(result, "_", "")
// 	return result
// }

// // fix in config gen
// func checkHasSdkName(param config.Parameter, sdkName string) string {
// 	if param.AliasType != "" {
// 		if strings.HasPrefix(param.AliasType, sdkName) {
// 			return param.AliasType
// 		}
// 		return fmt.Sprintf("%s.%s", sdkName, processStringType(param.AliasType))
// 	}
// 	return fmt.Sprintf("%s.%s", sdkName, processStringType(param.Type))
// }

// func ServiceSDKParamCreate(param config.Parameter, sdkName string) string {
// 	if param.Name == "ctx" {
// 		return ""
// 	}

// 	if param.IsPrimitive && !param.IsArray {
// 		typeName := getParamTypeName(param, sdkName)
// 		return fmt.Sprintf("var %s %s", param.Name, typeName)
// 	}
// 	if param.IsPrimitive && param.IsArray {
// 		typeName := getParamTypeName(param, sdkName)
// 		return fmt.Sprintf("var %s %s", param.Name, typeName)
// 	}

// 	if !param.IsPrimitive && !param.IsArray {
// 		typeName := strings.Replace(param.Type, "*", "", 1)
// 		return fmt.Sprintf("%s := %s{}", param.Name, typeName)
// 	}

// 	return ""
// }

// func getParamTypeName(param config.Parameter, sdkName string) string {
// 	if param.AliasType == "" {
// 		return param.Type
// 	}
// 	if param.IsArray {
// 		if param.AliasType != "" && sdkName != "" {
// 			return "[]" + sdkName + "." + param.AliasType
// 		}
// 		return "[]" + param.AliasType
// 	}
// 	return param.AliasType
// }

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
	case "map[string]string":
		return "map[string]string{}"
	default:
		return "\"\""
	}
}

// func translateTypeToCobraFlagComplex(field config.Parameter) string {
// 	typeName := cleanTypeName(field.Type)
// 	if field.IsArray {
// 		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
// 	}
// 	return fmt.Sprintf("JSONValue[%s]", typeName)
// }

// func cleanTypeName(typeName string) string {
// 	typeName = strings.TrimPrefix(typeName, "*")
// 	typeName = strings.TrimPrefix(typeName, "[]")
// 	return typeName
// }

// func translateTypeToCobraFlagStruct(field, parentField config.Parameter) string {
// 	typeName := cleanTypeName(parentField.Type)
// 	if canUseSliceFlag(parentField) {
// 		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
// 	}
// 	if canUseStrAsJson(parentField) {
// 		return fmt.Sprintf("JSONValue[%s]", typeName)
// 	}
// 	return translateTypeToCobraFlag(field.Type)
// }

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

// func translateTypeToCobraFlagCreateComplex(field config.Parameter) string {
// 	typeName := cleanTypeName(field.Type)
// 	if field.IsArray {
// 		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
// 	}
// 	return fmt.Sprintf("JSONValue[%s]", typeName)
// }

// func translateTypeToCobraFlagCreateStruct(field, parentField config.Parameter) string {
// 	typeName := cleanTypeName(parentField.Type)
// 	if canUseSliceFlag(parentField) {
// 		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
// 	}
// 	if canUseStrAsJson(parentField) {
// 		return fmt.Sprintf("JSONValue[%s]", typeName)
// 	}
// 	return translateTypeToCobraFlag(field.Type)
// }

func translateTypeToCobraFlagCreate(param config.Parameter, parents []config.Parameter) string {
	if param.IsArray {
		return "StrSlice"
	}

	if len(parents) > 0 {
		if parents[len(parents)-1].IsArray {
			return "StrSlice"
		}
	}

	paramType := strings.TrimPrefix(param.Type, "*")

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

// func canUseSliceFlag(parentField config.Parameter) bool {
// 	return parentField.IsArray && len(parentField.Struct) == 1
// }

// func canUseStrAsJson(parentField config.Parameter) bool {
// 	return parentField.IsArray && len(parentField.Struct) > 1
// }

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
