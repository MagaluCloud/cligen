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
	genDir              = "base-cli-gen/cmd/gen"
	importCobra         = "\"github.com/spf13/cobra\""
	importSDK           = "sdk \"github.com/MagaluCloud/mgc-sdk-go/client\""
	serviceParamPattern = "sdkCoreConfig sdk.CoreClient"
)

func init() {
	var err error
	menuItemTmpl, err = template.New("menu_item").Parse(menuItemTemplate)
	if err != nil {
		panic(err)
	}
}

func GenerateMenuItem(cfg *config.Config) error {
	for _, menu := range cfg.Menus {
		for _, submenu := range menu.Menus {
			GenMenuItem(cfg, submenu)
		}
	}

	return nil
}

func GenMenuItem(cfg *config.Config, menu *config.Menu) error {
	for _, method := range menu.Methods {
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
		menuItem.SetUseName(method.Name)
		menuItem.SetShortDescription(method.Description)
		menuItem.SetLongDescription(method.LongDescription)
		menuItem.AddImport(fmt.Sprintf("%s \"%s\"", "flags", "github.com/magaluCloud/mgccli/cobra_utils/flags"))
		for _, param := range method.Parameters {
			fd := FlagsDefinition(param)
			for _, item := range fd {
				menuItem.AddCobraFlagsDefinition(item)
			}
			//
			spc := ServiceSDKParamCreate(param, sdkName)
			menuItem.AddServiceSDKParamCreate(spc)

			cfc := CobraFlagsCreation(param)
			menuItem.AddCobraFlagsCreation(cfc)

		}

		// menuItem.AddCobraFlagsAssign()
		// menuItem.AddCobraFlagsRequired()
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

func CobraFlagsCreation(param config.Parameter) string {
	if param.Name == "ctx" {
		return ""
	}

	return ""
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

func FlagsDefinition(param config.Parameter) []string {
	if param.Name == "ctx" {
		return nil
	}
	if param.IsPrimitive {
		return []string{fmt.Sprintf("var %sFlag *flags.%s", param.Name, translateTypeToCobraFlag(param.Type))}
	}
	if !param.IsPrimitive {
		if !param.IsArray {
			if len(param.Struct) > 0 {
				var results []string
				for _, sparam := range param.Struct {
					sparam.Name = fmt.Sprintf("%s_%s", param.Name, sparam.Name)
					results = append(results, FlagsDefinition(sparam)...)
				}
				return results
			}
		}

	}

	return nil
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

func translateTypeToCobraFlagCreate(paramType string, withChar bool) string {
	paramType = strings.TrimPrefix(paramType, "*")

	switch paramType {
	case "string":
		if withChar {
			return "StrP"
		}
		return "Str"
	case "int64", "int32", "int16", "int8":
		if withChar {
			return "Int64P"
		}
		return "Int64"
	case "bool":
		if withChar {
			return "BoolP"
		}
		return "Bool"
	case "int":
		if withChar {
			return "IntP"
		}
		return "Int"
	case "float64":
		if withChar {
			return "Float64P"
		}
		return "Float64"
	case "[]string":
		if withChar {
			return "StrSliceP"
		}
		return "StrSlice"
	case "map[string]string":
		if withChar {
			return "StrMapP"
		}
		return "StrMap"
	case "time.Time":
		if withChar {
			return "TimeP"
		}
		return "Time"
	default:
		if withChar {
			return "StrP"
		}
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
