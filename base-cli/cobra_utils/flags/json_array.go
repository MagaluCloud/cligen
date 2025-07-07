package cobrautils

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

// JSONValue é um tipo genérico para parsear flags JSON usando pflag/cobra
// Exemplo de uso:
//
//	var myFlag cobrautils.JSONValue[MyType]
//	flags.Var(&myFlag, "myflag", "Some JSON input")
type JSONArrayValue[T any] struct {
	baseFlag
	Value []T
}

// Set faz o parse do JSON recebido na flag para o tipo T
func (j *JSONArrayValue[T]) Set(val string) error {
	var v []T
	if err := json.Unmarshal([]byte(val), &v); err != nil {
		return fmt.Errorf("invalid JSON for flag: %w", err)
	}
	j.Value = append(j.Value, v...)
	return nil
}

// String serializa o valor atual para JSON
func (j *JSONArrayValue[T]) String() string {
	b, err := json.Marshal(j.Value)
	if err != nil || len(b) == 0 {
		return "[]"
	}
	return string(b)
}

// Type retorna o nome do tipo para pflag
func (j *JSONArrayValue[T]) Type() string {
	return "json-array"
}

func NewJSONArrayValue[T any](cmd *cobra.Command, name string, usage string) *JSONArrayValue[T] {
	var value JSONArrayValue[T]
	cmd.Flags().Var(&value, name, usage)
	return &JSONArrayValue[T]{baseFlag: baseFlag{cmd, name}, Value: value.Value}
}

func NewJSONArrayValueP[T any](cmd *cobra.Command, name string, shorthand string, usage string) *JSONArrayValue[T] {
	var value JSONArrayValue[T]
	cmd.Flags().VarP(&value, name, shorthand, usage)
	return &JSONArrayValue[T]{baseFlag: baseFlag{cmd, name}, Value: value.Value}
}
