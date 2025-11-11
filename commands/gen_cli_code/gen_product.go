package gen_cli_code

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/commands/sdk_structure"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

func genProductCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genProductCodeRecursive(&pkg, nil)
	}
	return nil
}

func genProductCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	for _, service := range pkg.Services {
		for _, method := range service.Methods {
			filePath := buildProductFilePath(pkg, parentPkg, service.Name, method.Name)

			productData := NewPackageGroupData()
			productData.SetFileID(filePath)

			if productData.HasCustomFile {
				if err := productData.WriteProductCustomToFile(filePath); err != nil {
					return fmt.Errorf("erro ao escrever arquivo customizado para método %s do serviço %s do pacote %s: %w", method.Name, service.Name, pkg.Name, err)
				}
				continue
			}

			if err := setupProductData(productData, pkg, service, method); err != nil {
				return fmt.Errorf("erro ao configurar dados do produto para método %s: %w", method.Name, err)
			}

			serviceCallParams := genProductParameters(productData, method.Parameters)
			productData.SetServiceSDKParam(strings.Join(serviceCallParams, ", "))
			printResult(productData, method)

			if err := productData.WriteProductToFile(filePath); err != nil {
				return fmt.Errorf("erro ao escrever arquivo para método %s do serviço %s do pacote %s: %w", method.Name, service.Name, pkg.Name, err)
			}
		}
	}

	for _, subPkg := range pkg.SubPkgs {
		if err := genProductCodeRecursive(&subPkg, pkg); err != nil {
			return err
		}
	}
	return nil
}

func buildProductFilePath(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package, serviceName, methodName string) string {
	parts := []string{genDir}
	if parentPkg != nil {
		parts = append(parts, strings.ToLower(parentPkg.Name))
	}
	parts = append(parts, strings.ToLower(pkg.Name), strings.ToLower(serviceName))
	return filepath.Join(parts...) + fmt.Sprintf("/%s.go", strings.ToLower(methodName))
}

func setupProductData(productData *PackageGroupData, pkg *sdk_structure.Package, service sdk_structure.Service, method sdk_structure.Method) error {
	productData.SetPackageName(service.Name)
	productData.AddImport(importCobra)
	productData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
	productData.SetFunctionName(method.Name)
	productData.SetUseName(strutils.FirstLower(method.Name))
	productData.SetDescriptions(pkg.Description, method.Description)
	productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
	productData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
	productData.SetServiceCall(fmt.Sprintf("%s.%s", strutils.FirstLower(service.Interface), method.Name))
	return nil
}

func addPrintError() string {
	return "\n\t\t\tif err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\t"
}

func printResult(productData *PackageGroupData, method sdk_structure.Method) {
	assignResult := make([]string, 0, len(method.Returns))
	printRst := make([]string, 0, len(method.Returns))
	hasNonError := false

	for _, response := range method.Returns {
		assignResult = append(assignResult, response.Name)

		if response.Type == "error" {
			printRst = append(printRst, addPrintError())
			continue
		}

		hasNonError = true
	}

	if hasNonError {
		printRst = append(printRst, "\t\t\traw, _ := cmd.Root().PersistentFlags().GetBool(\"raw\")")
		productData.AddImport("\"github.com/magaluCloud/mgccli/beautiful\"")

		for _, response := range method.Returns {
			if response.Type != "error" {
				printRst = append(printRst, fmt.Sprintf("\t\t\tbeautiful.NewOutput(raw).PrintData(%s)", response.Name))
			}
		}
	}

	productData.AddAssignResult(strings.Join(assignResult, ", "))
	productData.AddPrintResult(strings.Join(printRst, "\n"))
}

func genProductParameters(productData *PackageGroupData, params []sdk_structure.Parameter) []string {
	serviceCallParams := make([]string, 0, len(params))
	flagsImportAdded := false

	for _, param := range params {
		if param.Type == "context.Context" {
			serviceCallParams = append(serviceCallParams, param.Name)
			continue
		}

		if !flagsImportAdded && len(params) > 0 {
			productData.AddImport("flags \"github.com/magaluCloud/mgccli/cobra_utils/flags\"")
			flagsImportAdded = true
		}

		if param.IsPrimitive {
			typeName := getParamTypeName(param)
			productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", param.Name, typeName))
			processFieldsRecursive(productData, []sdk_structure.Parameter{param}, "", nil)
		} else {
			typeName := strings.Replace(param.Type, "*", "", 1)
			productData.AddServiceSDKParamCreate(fmt.Sprintf("%s := %s{}", param.Name, typeName))
			processFieldsRecursive(productData, mapToSlice(param.Struct), param.Name, &param)
		}

		callName := getParamCallName(param)
		serviceCallParams = append(serviceCallParams, callName)
	}

	return serviceCallParams
}

func getParamTypeName(param sdk_structure.Parameter) string {
	if param.AliasType == "" {
		return param.Type
	}
	if param.IsArray {
		return "[]" + param.AliasType
	}
	return param.AliasType
}

func getParamCallName(param sdk_structure.Parameter) string {
	if param.IsPointer && !param.IsPrimitive {
		return "&" + param.Name
	}
	return param.Name
}

func prepareCommandFlag(str string) string {
	strSplit := strings.Split(str, "_")[1:]
	for i, s := range strSplit {
		if len(s) > 2 {
			strSplit[i] = strutils.ToSnakeCasePreserveID(s, "-")
		}
	}
	result := strings.Join(strSplit, ".")
	result = strings.ToLower(result)
	return result
}

// mapToSlice converte um map de parâmetros em slice
func mapToSlice(paramMap map[string]sdk_structure.Parameter) []sdk_structure.Parameter {
	result := make([]sdk_structure.Parameter, 0, len(paramMap))
	for _, param := range paramMap {
		result = append(result, param)
	}
	return result
}

// processFieldsRecursive processa campos de forma recursiva (unifica a lógica de genProductParameters e genProductParametersRecursive)
func processFieldsRecursive(productData *PackageGroupData, fields []sdk_structure.Parameter, pathPrefix string, parentField *sdk_structure.Parameter) {
	for _, field := range fields {
		productData.SetAllowedPositionalArgs()
		// Constrói o caminho atual
		var currentPath string
		if pathPrefix == "" {
			currentPath = field.Name
		} else {
			currentPath = pathPrefix + "." + field.Name
		}
		varFlagName := strings.ReplaceAll(currentPath, ".", "_")
		// Verifica se está em nível profundo (aninhado) baseado no número de pontos no path
		isDeepNested := strings.Count(pathPrefix, ".") > 0

		if field.IsPrimitive {
			field.IsPositional = false

			var cobraFlagName string
			var flagTypeGetter, flagCreationGetter func() string
			var defaultGetter func() string

			// Determina getters baseado no contexto
			if parentField == nil || !isDeepNested {
				// Parâmetros primitivos diretos ou campos de primeiro nível de struct
				cobraFlagName = strutils.ToSnakeCasePreserveID(field.Name, "-")
				flagTypeGetter = func() string { return translateTypeToCobraFlag(field.Type) }
				flagCreationGetter = func() string { return translateTypeToCobraFlagCreate(field.Type, false) }
				defaultGetter = func() string { return defaultByType(field.Type) }

				if !field.IsOptional && !field.IsArray && parentField == nil {
					if parentField != nil && !parentField.IsPrimitive {
						productData.SetNotAllowedPositionalArgs()
					}
					canPositionalArgs := productData.CanAddPositionalArgs(productData.UseName)
					if canPositionalArgs {
						field.PositionalIndex = productData.AddPositionalArgs(cobraFlagName)
						field.IsPositional = true
					}
					if !canPositionalArgs {
						productData.SetNotAllowedPositionalArgs()
					}

					if field.IsPositional {
						field.Description = fmt.Sprintf("%s (required)", field.Description)
					}
				}

			} else {
				// Campos aninhados profundos - usa lógica de struct
				cobraFlagName = prepareCommandFlag(varFlagName)
				flagTypeGetter = func() string { return translateTypeToCobraFlagStruct(field, *parentField) }
				flagCreationGetter = func() string {
					if canUseSliceFlag(*parentField) || canUseStrAsJson(*parentField) {
						return translateTypeToCobraFlagCreateStruct(field, *parentField)
					}
					return translateTypeToCobraFlagCreate(field.Type, false)
				}
				// Só usa default se não for slice ou JSON
				if !canUseSliceFlag(*parentField) && !canUseStrAsJson(*parentField) {
					defaultGetter = func() string { return defaultByType(field.Type) }
				}
			}

			addPrimitiveFlag(productData, FlagCreationConfig{
				FlagName:           varFlagName,
				TargetVar:          currentPath,
				CobraFlagName:      cobraFlagName,
				Field:              field,
				ParentField:        parentField,
				FlagTypeGetter:     flagTypeGetter,
				FlagCreationGetter: flagCreationGetter,
				DefaultValueGetter: defaultGetter,
				AliasType:          field.AliasType,
				IsPositional:       field.IsPositional,
				PositionalIndex:    field.PositionalIndex,
			})
		}

		if !field.IsPrimitive && !field.IsArray {
			// Inicializa structs aninhados se necessário (apenas em níveis profundos)
			if isDeepNested && parentField != nil {
				// Inicializa o parent se for pointer
				if parentField.IsPointer {
					productData.AddCobraStructInitialize(fmt.Sprintf("%s = &%s{}", pathPrefix, strings.Replace(parentField.Type, "*", "", 1)))
				}
			}
			// Inicializa o field atual se for pointer
			if field.IsPointer {
				productData.AddCobraStructInitialize(fmt.Sprintf("%s = &%s{}", currentPath, strings.Replace(field.Type, "*", "", 1)))
			}
			// Recursão para processar os campos do struct aninhado
			processFieldsRecursive(productData, mapToSlice(field.Struct), currentPath, &field)
		}

		if !field.IsPrimitive && field.IsArray {
			varCommandName := prepareCommandFlag(varFlagName)

			addPrimitiveFlag(productData, FlagCreationConfig{
				FlagName:           varFlagName,
				TargetVar:          currentPath,
				CobraFlagName:      varCommandName,
				Field:              field,
				ParentField:        parentField,
				FlagTypeGetter:     func() string { return translateTypeToCobraFlagComplex(field) },
				FlagCreationGetter: func() string { return translateTypeToCobraFlagCreateComplex(field) },
				DefaultValueGetter: nil, // Arrays complexos não têm valor default
				AliasType:          field.AliasType,
			})
		}
	}
}

func canUseSliceFlag(parentField sdk_structure.Parameter) bool {
	return parentField.IsArray && len(parentField.Struct) == 1
}

func canUseStrAsJson(parentField sdk_structure.Parameter) bool {
	return parentField.IsArray && len(parentField.Struct) > 1
}

// FlagAssignmentConfig contém as configurações para gerar código de atribuição de flags
type FlagAssignmentConfig struct {
	FlagName          string
	TargetVar         string
	CobraFlagName     string
	Field             sdk_structure.Parameter
	ParentField       *sdk_structure.Parameter
	AliasType         string
	IsOptional        bool
	RequirePositional bool
	PositionalIndex   int
}

// FlagCreationConfig contém todas as configurações necessárias para criar uma flag completa
type FlagCreationConfig struct {
	FlagName           string
	TargetVar          string
	CobraFlagName      string
	Field              sdk_structure.Parameter
	ParentField        *sdk_structure.Parameter
	FlagTypeGetter     func() string
	FlagCreationGetter func() string
	DefaultValueGetter func() string
	AliasType          string
	IsPositional       bool
	PositionalIndex    int
}

// addPrimitiveFlag cria definição, criação e atribuição de uma flag primitiva
func addPrimitiveFlag(productData *PackageGroupData, config FlagCreationConfig) {
	// Adiciona definição da flag
	productData.AddCobraFlagsDefinition(fmt.Sprintf("var %sFlag *flags.%s", config.FlagName, config.FlagTypeGetter()))

	// Adiciona criação da flag
	defaultValue := ""
	if config.DefaultValueGetter != nil {
		defaultValue = config.DefaultValueGetter()
	}

	if defaultValue != "" {
		productData.AddCobraFlagsCreation(
			fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", %s, \"%s\")",
				config.FlagName,
				config.FlagCreationGetter(),
				config.CobraFlagName,
				defaultValue,
				strutils.RemoveNewLine(strutils.EscapeQuotes(config.Field.Description)),
			),
		)
	} else {
		productData.AddCobraFlagsCreation(
			fmt.Sprintf("%sFlag = flags.New%s(cmd, \"%s\", \"%s\",)",
				config.FlagName,
				config.FlagCreationGetter(),
				config.CobraFlagName,
				strutils.RemoveNewLine(strutils.EscapeQuotes(config.Field.Description)),
			),
		)
	}

	// Calcula targetVar se houver ParentField com casos especiais
	targetVar := config.TargetVar
	if config.ParentField != nil {
		if canUseSliceFlag(*config.ParentField) || canUseStrAsJson(*config.ParentField) {
			targetVar = strings.TrimSuffix(config.TargetVar, "."+config.Field.Name)
		}
	}

	// Cria e adiciona o código de atribuição
	command := createFlagAssignment(FlagAssignmentConfig{
		FlagName:          config.FlagName,
		TargetVar:         targetVar,
		CobraFlagName:     config.CobraFlagName,
		Field:             config.Field,
		ParentField:       config.ParentField,
		IsOptional:        config.Field.IsOptional,
		RequirePositional: config.IsPositional,
		AliasType:         config.AliasType,
		PositionalIndex:   config.PositionalIndex,
	})

	productData.AddCobraFlagsAssign(command)

	// Adiciona import do fmt se necessário
	if strings.Contains(command, "fmt") {
		productData.AddImport("\"fmt\"")
	}

}

// createFlagAssignment gera código para atribuir o valor de uma flag a uma variável
func createFlagAssignment(config FlagAssignmentConfig) string {
	flagVar := config.FlagName + "Flag"

	// Casos especiais para structs aninhados com slice ou JSON
	if config.ParentField != nil {
		if canUseSliceFlag(*config.ParentField) || canUseStrAsJson(*config.ParentField) {
			return fmt.Sprintf("if %s.IsChanged() {\n\t\t\t\t%s = %s.Value\n\t\t\t}", flagVar, config.TargetVar, flagVar)
		}
	}

	// Detecta slice de tipo não-primitivo (primitivo com alias): []string -> []ImageExpand
	if config.AliasType != "" && config.Field.IsArray && config.Field.IsPrimitive {
		return fmt.Sprintf(`if %s.IsChanged() {
				%s = make([]%s, len(*%s.Value))
				for i, v := range *%s.Value {
					%s[i] = %s(v)
				}
			}`, flagVar, config.TargetVar, config.AliasType, flagVar, flagVar, config.TargetVar, config.AliasType)
	}

	// Se tem AliasType (mas não é slice), usa casting
	if config.AliasType != "" {
		// Se é pointer, faz cast com pointer
		if config.Field.IsPointer {
			return fmt.Sprintf("if %s.IsChanged() {\n\t\t\t\t%s = (*%s)(%s.Value)\n\t\t\t}", flagVar, config.TargetVar, config.AliasType, flagVar)
		}
		// Se é opcional ou já existe flag posicional, faz cast desreferenciando
		if config.IsOptional || !config.RequirePositional {
			return fmt.Sprintf("if %s.IsChanged() {\n\t\t\t\t%s = (%s)(*%s.Value)\n\t\t\t}", flagVar, config.TargetVar, config.AliasType, flagVar)
		}
		// Caso requer argumento posicional: verifica args ou flag com cast
		return fmt.Sprintf(`if len(args) > 0{
				cmd.Flags().Set("%s", args[%d])
			}
			if %s.IsChanged() {
				%s = (%s)(*%s.Value)
			} else {
				return fmt.Errorf("é necessário fornecer o %s como argumento ou usar a flag --%s")
			}`, config.CobraFlagName, config.PositionalIndex, flagVar, config.TargetVar, config.AliasType, flagVar, config.CobraFlagName, config.CobraFlagName)
	}

	// Se é pointer, atribui diretamente o valor
	if config.Field.IsPointer {
		return fmt.Sprintf("if %s.IsChanged() {\n\t\t\t\t%s = %s.Value\n\t\t\t}", flagVar, config.TargetVar, flagVar)
	}

	// Se é opcional ou já existe flag posicional, atribui com desreferência
	if config.IsOptional || !config.RequirePositional {
		return fmt.Sprintf("if %s.IsChanged() {\n\t\t\t\t%s = *%s.Value\n\t\t\t}", flagVar, config.TargetVar, flagVar)
	}

	// Caso requer argumento posicional: verifica args ou flag
	return fmt.Sprintf(`if len(args) > 0{
				cmd.Flags().Set("%s", args[%d])
			}
			if %s.IsChanged() {
				%s = *%s.Value
			} else {
				return fmt.Errorf("é necessário fornecer o %s como argumento ou usar a flag --%s")
			}`, config.CobraFlagName, config.PositionalIndex, flagVar, config.TargetVar, flagVar, config.CobraFlagName, config.CobraFlagName)
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

func translateTypeToCobraFlagComplex(field sdk_structure.Parameter) string {
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

func translateTypeToCobraFlagStruct(field, parentField sdk_structure.Parameter) string {
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

func translateTypeToCobraFlagCreateComplex(field sdk_structure.Parameter) string {
	typeName := cleanTypeName(field.Type)
	if field.IsArray {
		return fmt.Sprintf("JSONArrayValue[%s]", typeName)
	}
	return fmt.Sprintf("JSONValue[%s]", typeName)
}

func translateTypeToCobraFlagCreateStruct(field, parentField sdk_structure.Parameter) string {
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
