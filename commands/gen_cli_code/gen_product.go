package gen_cli_code

import (
	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
	"fmt"
	"path/filepath"
	"strings"
)

func genProductCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genProductCodeRecursive(&pkg, nil)
	}
	return nil
}

func genProductCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	if len(pkg.Services) > 0 {
		for _, service := range pkg.Services {
			for _, method := range service.Methods {
				productData := NewPackageGroupData()
				productData.SetPackageName(service.Name)
				productData.AddImport(importCobra)
				productData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
				productData.SetFunctionName(method.Name)
				productData.SetUseName(strutils.FirstLower(method.Name))
				productData.SetDescriptions(defaultShortDesc, defaultLongDesc)
				productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
				productData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
				productData.SetServiceCall(fmt.Sprintf("%s.%s", strutils.FirstLower(service.Interface), method.Name))

				serviceCallParams := genProductParameters(productData, method.Parameters)
				productData.SetServiceSDKParam(strings.Join(serviceCallParams, ", "))
				printResult(productData, method)
				dir := genDir
				if parentPkg != nil {
					dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(method.Name)))
				} else {
					dir = filepath.Join(dir, strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(method.Name)))
				}
				err := productData.WriteProductToFile(dir)
				if err != nil {
					return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
				}
			}
		}
	}
	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			genProductCodeRecursive(&subPkg, pkg)
		}
	}
	return nil
}

func printResult(productData *PackageGroupData, method sdk_structure.Method) {
	// first check responses and create a AssignResult
	assignResult := []string{}
	printResult := []string{}
	for _, response := range method.Returns {
		assignResult = append(assignResult, response.Name)
		if response.Type == "error" {
			printResult = append(printResult, "\t\t\tif err != nil {\n\t\t\t\tfmt.Println(err.Error())\n\t\t\t}")
			continue
		}
		printResult = append(printResult, fmt.Sprintf("sdkResult, err := json.MarshalIndent(%s, \"\", \"  \")\n\t\t\tif err != nil {\n\t\t\t\tfmt.Println(err.Error())\n\t\t\t}\n\t\t\tfmt.Println(string(sdkResult))", response.Name))
		productData.AddImport("\"encoding/json\"")
	}
	productData.AddAssignResult(strings.Join(assignResult, ", "))
	productData.AddPrintResult(strings.Join(printResult, "\n"))
}

func genProductParameters(productData *PackageGroupData, params []sdk_structure.Parameter) []string {
	var serviceCallParams []string

	for i, param := range params {
		if i != param.Position {
			fmt.Printf("   ❌ Parâmetro %s não está na posição %d\n", param.Name, param.Position)
		}
		if param.Type == "context.Context" {
			serviceCallParams = append(serviceCallParams, param.Name)
			continue
		}
		if len(params) > 0 {
			productData.AddImport("flags \"mgccli/cobra_utils/flags\"")
		}

		if param.IsPrimitive {
			productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, param.Type))
			productData.AddCobraFlagsDefinition(fmt.Sprintf("var %sFlag *flags.%s", param.Name, translateTypeToCobraFlag(param.Type)))
			initialChar := strutils.FirstUnusedChar(param.Name, &productData.UsedChars)
			cobraFlagName := strutils.ToSnakeCase(param.Name, "-")
			productData.AddCobraFlagsCreation(
				fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", %s, \"%s\")",
					param.Name,
					translateTypeToCobraFlagCreate(param.Type, true),
					cobraFlagName,
					initialChar,
					defaultByType(param.Type),
					strutils.RemoveNewLine(strutils.EscapeQuotes(param.Description)),
				),
			)
			productData.AddCobraFlagsAssign(createPrimitiveFlagToAssign(param.Name, param.IsPointer))
			if !param.IsPointer {
				addRequiredFlag(productData, param, param.Name)
			}
		}

		if !param.IsPrimitive {
			productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, param.Type))
			for _, field := range param.Struct {
				if field.IsPrimitive {
					productData.AddCobraFlagsDefinition(fmt.Sprintf("var %s_%sFlag *flags.%s", param.Name, field.Name, translateTypeToCobraFlag(field.Type)))
					initialChar := strutils.FirstUnusedChar(field.Name, &productData.UsedChars)
					cobraFlagName := strutils.ToSnakeCase(field.Name, "-")
					productData.AddCobraFlagsCreation(
						fmt.Sprintf("%s_%sFlag = flags.New%s(cmd, \"%s\", \"%s\", %s, \"%s\")",
							param.Name,
							field.Name,
							translateTypeToCobraFlagCreate(field.Type, true),
							cobraFlagName,
							initialChar,
							defaultByType(field.Type),
							strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
						),
					)
					productData.AddCobraFlagsAssign(createStructFlagToAssign(param.Name, field.Name, field.IsPointer))
					if !field.IsPointer {
						addRequiredFlag(productData, field, cobraFlagName)
					}
				}
				// Here is a struct, we need some recursive call to generate the code for the struct
				if !field.IsPrimitive && !field.IsArray {
					genProductParametersRecursive(productData, field, param.Name)
				}
				if !field.IsPrimitive && field.IsArray {
					// if canUseStrAsJson(field, param) {
					currentPath := param.Name + "." + field.Name
					varFlagName := strings.ReplaceAll(currentPath, ".", "_")
					varCommandName := prepareCommandFlag(varFlagName)

					productData.AddCobraFlagsDefinition(fmt.Sprintf("var %sFlag *flags.%s", varFlagName, translateTypeToCobraFlagComplex(field))) // translateTypeToCobraFlag(field.Type)))
					initialChar := strutils.FirstUnusedChar(field.Name, &productData.UsedChars)

					productData.AddCobraFlagsCreation(
						fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", \"%s\",)",
							varFlagName,
							translateTypeToCobraFlagCreateComplex(field),
							varCommandName,
							initialChar,
							strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
						),
					)

					productData.AddCobraFlagsAssign(createPrimitiveFlagToAssignStruct(varFlagName, currentPath, field, param))

					// }

				}
			}
		}

		serviceCallParams = append(serviceCallParams, param.Name)
	}
	return serviceCallParams
}

func addRequiredFlag(productData *PackageGroupData, param sdk_structure.Parameter, flagName string) {
	if !param.IsPointer && !param.IsArray {
		productData.AddCobraFlagsRequired(fmt.Sprintf("cmd.MarkFlagRequired(\"%s\")", flagName))
	}
}

func prepareCommandFlag(str string) string {
	strSplit := strings.Split(str, "_")[1:]
	for i, s := range strSplit {
		if len(s) > 2 {
			strSplit[i] = strutils.ToSnakeCase(s, "-")
		}
	}
	result := strings.Join(strSplit, ".")
	result = strings.ToLower(result)
	return result
}

func genProductParametersRecursive(productData *PackageGroupData, parentField sdk_structure.Parameter, parentStructName string) {
	for _, field := range parentField.Struct {
		if field.IsPrimitive {
			currentPath := ""
			currentPath = parentStructName + "." + parentField.Name + "." + field.Name
			varFlagName := strings.ReplaceAll(currentPath, ".", "_")
			varCommandName := prepareCommandFlag(varFlagName)

			productData.AddCobraFlagsDefinition(fmt.Sprintf("var %sFlag *flags.%s", varFlagName, translateTypeToCobraFlagStruct(field, parentField))) // translateTypeToCobraFlag(field.Type)))
			initialChar := strutils.FirstUnusedChar(field.Name, &productData.UsedChars)

			if canUseSliceFlag(field, parentField) {
				productData.AddCobraFlagsCreation(
					fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", \"%s\",)",
						varFlagName,
						translateTypeToCobraFlagCreateStruct(field, parentField),
						varCommandName,
						initialChar,
						strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
					),
				)
			}

			if canUseStrAsJson(field, parentField) {
				productData.AddCobraFlagsCreation(
					fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", \"%s\",)",
						varFlagName,
						translateTypeToCobraFlagCreateStruct(field, parentField),
						varCommandName,
						initialChar,
						strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
					),
				)
			}
			if !canUseSliceFlag(field, parentField) && !canUseStrAsJson(field, parentField) {
				productData.AddCobraFlagsCreation(
					fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", %s, \"%s\")",
						varFlagName,
						translateTypeToCobraFlagCreate(field.Type, true),
						varCommandName,
						initialChar,
						defaultByType(field.Type),
						strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
					),
				)
			}
			productData.AddCobraFlagsAssign(createPrimitiveFlagToAssignStruct(varFlagName, currentPath, field, parentField))

		}

		if !field.IsPrimitive {
			localVar := parentStructName + "." + parentField.Name
			if parentField.IsPointer {
				productData.AddCobraStructInitialize(fmt.Sprintf("%s = &%s{}", localVar, strings.Replace(parentField.Type, "*", "", 1)))
			}

			localVar = localVar + "." + field.Name
			if field.IsPointer {
				productData.AddCobraStructInitialize(fmt.Sprintf("%s = &%s{}", localVar, strings.Replace(field.Type, "*", "", 1)))
			}
			genProductParametersRecursive(productData, field, parentStructName+"."+parentField.Name)

		}
	}
}

func canUseSliceFlag(field, parentField sdk_structure.Parameter) bool {
	if !parentField.IsArray {
		return false
	}
	if len(parentField.Struct) == 1 {
		return true
	}
	return false
}

func canUseStrAsJson(field, parentField sdk_structure.Parameter) bool {
	if !parentField.IsArray {
		return false
	}
	if len(parentField.Struct) > 1 {
		return true
	}
	return false
}

func createPrimitiveFlagToAssignStruct(flagName string, parentStructName string, field, parentField sdk_structure.Parameter) string {
	if canUseSliceFlag(field, parentField) {
		return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = %sFlag.Value\n\t\t\t}", flagName, strings.TrimSuffix(parentStructName, "."+field.Name), flagName)
	}
	if canUseStrAsJson(field, parentField) {
		return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = %sFlag.Value\n\t\t\t}", flagName, strings.TrimSuffix(parentStructName, "."+field.Name), flagName)
	}

	if field.IsPointer {
		return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = %sFlag.Value\n\t\t\t}", flagName, parentStructName, flagName)
	}
	return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = *%sFlag.Value\n\t\t\t}", flagName, parentStructName, flagName)
}

func createPrimitiveFlagToAssign(flagName string, isPointer bool) string {
	if isPointer {
		return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = %sFlag.Value\n\t\t\t}", flagName, flagName, flagName)
	}
	return fmt.Sprintf("if %sFlag.IsChanged() {\n\t\t\t\t%s = *%sFlag.Value\n\t\t\t}", flagName, flagName, flagName)
}

func createStructFlagToAssign(paramName, fieldName string, isPointer bool) string {
	if isPointer {
		return fmt.Sprintf("if %s_%sFlag.IsChanged() {\n\t\t\t\t%s.%s = %s_%sFlag.Value\n\t\t\t}", paramName, fieldName, paramName, fieldName, paramName, fieldName)
	}
	return fmt.Sprintf("if %s_%sFlag.IsChanged() {\n\t\t\t\t%s.%s = *%s_%sFlag.Value\n\t\t\t}", paramName, fieldName, paramName, fieldName, paramName, fieldName)
}

func defaultByType(paramType string) string {
	paramType = strings.TrimPrefix(paramType, "*")
	switch paramType {
	case "string":
		return "\"\""
	case "int64", "int32", "int16", "int8", "int":
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

func translateTypeToCobraFlagComplex(field sdk_structure.Parameter) string {
	typeName := strings.TrimPrefix(field.Type, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	if field.IsArray {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}

	return fmt.Sprintf("JSONValue[%s]", typeName)

}

func translateTypeToCobraFlagStruct(field, parentField sdk_structure.Parameter) string {
	typeName := strings.TrimPrefix(parentField.Type, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	if canUseSliceFlag(field, parentField) {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	if canUseStrAsJson(field, parentField) {
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
	case "[]string":
		return "StrSliceFlag"
	case "map[string]string":
		return "StrMapFlag"
	default:
		return "StrFlag"
	}
}

func translateTypeToCobraFlagCreateComplex(field sdk_structure.Parameter) string {
	typeName := strings.TrimPrefix(field.Type, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	if field.IsArray {
		return fmt.Sprintf("JSONArrayValueP[%s]", typeName)
	}

	return fmt.Sprintf("JSONValueP[%s]", typeName)
}

func translateTypeToCobraFlagCreateStruct(field, parentField sdk_structure.Parameter) string {
	typeName := strings.TrimPrefix(parentField.Type, "*")
	typeName = strings.TrimPrefix(typeName, "[]")
	if canUseSliceFlag(field, parentField) {
		return fmt.Sprintf("JSONArrayValueP[%s]", typeName)
	}
	if canUseStrAsJson(field, parentField) {
		return fmt.Sprintf("JSONValueP[%s]", typeName)
	}
	return translateTypeToCobraFlag(field.Type)
}

func translateTypeToCobraFlagCreate(paramType string, withChar bool) string {
	// StrFlag
	// Int64Flag
	// BoolFlag
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
	default:
		if withChar {
			return "StrP"
		}
		return "Str"
	}
}

// PrintResult
