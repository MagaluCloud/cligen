package strutils

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		char     string
		expected string
	}{

		{
			name:     "BlockStorage com separador padrão",
			input:    "BlockStorage",
			char:     "",
			expected: "block_storage",
		},
		{
			name:     "ContainerRegistry com separador padrão",
			input:    "ContainerRegistry",
			char:     "",
			expected: "container_registry",
		},
		{
			name:     "VirtualMachine com separador padrão",
			input:    "VirtualMachine",
			char:     "",
			expected: "virtual_machine",
		},
		{
			name:     "attachToPort	 com separador padrão",
			input:    "attachToPort",
			char:     "",
			expected: "attach_to_port",
		},
		{
			name:     "VPCs com separador padrão",
			input:    "VPCs",
			char:     "",
			expected: "vpcs",
		},
		{
			name:     "networkACLs com hífen",
			input:    "networkACLs",
			char:     "-",
			expected: "network-acls",
		},
		{
			name:     "abobinhaSdkQuery com hífen",
			input:    "abobinhaSdkQuery",
			char:     "-",
			expected: "abobinha-sdk-query",
		},

		// {
		// 	name:     "Sequência de maiúsculas no início",
		// 	input:    "XMLParser",
		// 	char:     "-",
		// 	expected: "xml-parser",
		// },
		// {
		// 	name:     "Sequência de maiúsculas no meio",
		// 	input:    "parseXMLFile",
		// 	char:     "-",
		// 	expected: "parse-xml-file",
		// },
		{
			name:     "Uma letra maiúscula isolada",
			input:    "testeZ",
			char:     "-",
			expected: "teste-z",
		},
		{
			name:     "CamelCase simples",
			input:    "camelCase",
			char:     "_",
			expected: "camel_case",
		},
		{
			name:     "PascalCase com underscore",
			input:    "PascalCase",
			char:     "_",
			expected: "pascal_case",
		},
		{
			name:     "ID no final",
			input:    "userID",
			char:     "-",
			expected: "user-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSnakeCase(tt.input, tt.char)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q, %q) = %q, esperado %q", tt.input, tt.char, result, tt.expected)
			}
		})
	}
}
