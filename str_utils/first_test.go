package strutils

import "testing"

func TestRemovePlural(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ParametersGroupService", "ParameterGroupService"},
		{"VirtualMachines", "VirtualMachine"},
		{"NetworkInterfaces", "NetworkInterface"},
		{"LoadBalancers", "LoadBalancer"},
		{"SecurityGroups", "SecurityGroup"},
		{"Status", "Status"}, // Não deve alterar palavras que já são singulares
		{"Process", "Process"},
		{"Access", "Access"},
		{"", ""}, // String vazia
	}

	for _, test := range tests {
		result := RemovePlural(test.input)
		if result != test.expected {
			t.Errorf("RemovePlural(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ParametersGroupService", "parameters_group_service"},
		{"VirtualMachine", "virtual_machine"},
		{"API", "api"},
		{"", ""},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("ToSnakeCase(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
