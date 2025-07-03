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

func genProductParameters(productData *PackageGroupData, params []sdk_structure.Parameter) []string {
	var serviceCallParams []string
	if len(params) > 0 {
		productData.AddImport("flags \"mgccli/cobra_utils/flags\"")
	}

	var initialChars = make([]string, 0)

	for i, param := range params {
		if i != param.Position {
			fmt.Printf("   ❌ Parâmetro %s não está na posição %d\n", param.Name, param.Position)
		}
		if param.Type == "context.Context" {
			serviceCallParams = append(serviceCallParams, param.Name)
			continue
		}

		if param.IsPrimitive {
			productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, param.Type))
			productData.AddCobraFlagsDefinition(fmt.Sprintf("var %sFlag *flags.%s", param.Name, translateTypeToCobraFlag(param.Type)))
			initialChar := strutils.FirstUnusedChar(param.Name, &initialChars)
			productData.AddCobraFlagsCreation(
				fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\", %s, \"%s\")",
					param.Name,
					translateTypeToCobraFlagCreate(param.Type, true),
					strutils.ToSnakeCase(param.Name, "-"),
					initialChar,
					defaultByType(param.Type),
					strutils.RemoveNewLine(strutils.EscapeQuotes(param.Description)),
				),
			)
			productData.AddCobraFlagsAssign(createPrimitiveFlagToAssign(param.Name, param.IsPointer))
			if !param.IsPointer {
				productData.AddCobraFlagsRequired(fmt.Sprintf("cmd.MarkFlagRequired(\"%s\")", param.Name))
			}
		}

		if !param.IsPrimitive {
			productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, param.Type))
			for _, field := range param.Struct {
				if field.Struct == nil {
					productData.AddCobraFlagsDefinition(fmt.Sprintf("var %s_%sFlag *flags.%s", param.Name, field.Name, translateTypeToCobraFlag(field.Type)))
					initialChar := strutils.FirstUnusedChar(field.Name, &initialChars)
					productData.AddCobraFlagsCreation(
						fmt.Sprintf("%s_%sFlag = flags.New%s(cmd, \"%s\", \"%s\", %s, \"%s\")",
							param.Name,
							field.Name,
							translateTypeToCobraFlagCreate(field.Type, true),
							strutils.ToSnakeCase(field.Name, "-"),
							initialChar,
							defaultByType(field.Type),
							strutils.RemoveNewLine(strutils.EscapeQuotes(field.Description)),
						),
					)
					productData.AddCobraFlagsAssign(createStructFlagToAssign(param.Name, field.Name, field.IsPointer))
					if !field.IsPointer {
						productData.AddCobraFlagsRequired(fmt.Sprintf("cmd.MarkFlagRequired(\"%s\")", strutils.ToSnakeCase(field.Name, "-")))
					}
				}
			}
		}

		serviceCallParams = append(serviceCallParams, param.Name)
	}
	return serviceCallParams
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
	default:
		return "\"\""
	}
}

func translateTypeToCobraFlag(paramType string) string {
	// StrFlag
	// Int64Flag
	// BoolFlag
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
	default:
		return "StrFlag"
	}
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
	default:
		if withChar {
			return "StrP"
		}
		return "Str"
	}
}
