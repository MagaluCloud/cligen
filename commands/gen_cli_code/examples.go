package gen_cli_code

import (
	"fmt"
	"strings"
)

// Exemplo de como criar dados para o template do grupo "compute"
func ExampleComputePackageGroup() *PackageGroupData {
	data := NewPackageGroupData()

	// Configurar informações básicas
	data.SetPackageName("compute")
	data.SetFunctionName("Compute")
	data.SetUseName("compute")
	data.SetDescriptions("Comandos relacionados a computação", "Comandos para gerenciar recursos de computação")
	data.SetGroupID("products")
	data.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	// Adicionar imports necessários
	data.AddImport("mgccli/cmd/gen/compute/instances")
	data.AddImport("mgccli/cmd/gen/compute/images")
	data.AddImport("mgccli/cmd/gen/compute/instancetypes")
	data.AddImport("mgccli/cmd/gen/compute/snapshots")
	data.AddImport("sdk \"github.com/MagaluCloud/mgc-sdk-go/client\"")
	data.AddImport("\"github.com/MagaluCloud/mgc-sdk-go/compute\"")
	data.AddImport("\"github.com/spf13/cobra\"")

	// Adicionar subcomandos
	data.AddSubCommand("instances", "Instances", "computeService.Instances()")
	data.AddSubCommand("images", "Images", "computeService.Images()")
	data.AddSubCommand("instancetypes", "InstanceTypes", "computeService.InstanceTypes()")
	data.AddSubCommand("snapshots", "Snapshots", "computeService.Snapshots()")

	return data
}

// Exemplo de como criar dados para o template do subgrupo "instances"
func ExampleInstancesPackageGroup() *PackageGroupData {
	data := NewPackageGroupData()

	// Configurar informações básicas
	data.SetPackageName("instances")
	data.SetFunctionName("Instances")
	data.SetUseName("instances")
	data.SetDescriptions("Comandos relacionados a instâncias", "Comandos para gerenciar instâncias de computação")
	data.SetServiceParam("instanceService compute.InstanceService")

	// Adicionar imports necessários
	data.AddImport("\"github.com/MagaluCloud/mgc-sdk-go/compute\"")
	data.AddImport("\"github.com/spf13/cobra\"")

	// Adicionar subcomandos (métodos do serviço)
	data.AddSubCommand("", "List", "instanceService")
	data.AddSubCommand("", "Get", "instanceService")
	data.AddSubCommand("", "Create", "instanceService")
	data.AddSubCommand("", "Delete", "instanceService")

	return data
}

// Exemplo de como criar dados para o template do grupo "kubernetes"
func ExampleKubernetesPackageGroup() *PackageGroupData {
	data := NewPackageGroupData()

	// Configurar informações básicas
	data.SetPackageName("kubernetes")
	data.SetFunctionName("Kubernetes")
	data.SetUseName("kubernetes")
	data.SetDescriptions("Comandos relacionados a Kubernetes", "Comandos para gerenciar clusters Kubernetes")
	data.SetGroupID("products")
	data.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	// Adicionar imports necessários
	data.AddImport("mgccli/cmd/gen/kubernetes/clusters")
	data.AddImport("mgccli/cmd/gen/kubernetes/flavors")
	data.AddImport("mgccli/cmd/gen/kubernetes/nodepools")
	data.AddImport("mgccli/cmd/gen/kubernetes/versions")
	data.AddImport("sdk \"github.com/MagaluCloud/mgc-sdk-go/client\"")
	data.AddImport("\"github.com/MagaluCloud/mgc-sdk-go/kubernetes\"")
	data.AddImport("\"github.com/spf13/cobra\"")

	// Adicionar subcomandos
	data.AddSubCommand("clusters", "Clusters", "kubernetesService.Clusters()")
	data.AddSubCommand("flavors", "Flavors", "kubernetesService.Flavors()")
	data.AddSubCommand("nodepools", "Nodepools", "kubernetesService.Nodepools()")
	data.AddSubCommand("versions", "Versions", "kubernetesService.Versions()")

	return data
}

// Exemplo de como criar dados para o template do grupo "lbaas"
func ExampleLbaasPackageGroup() *PackageGroupData {
	data := NewPackageGroupData()

	// Configurar informações básicas
	data.SetPackageName("lbaas")
	data.SetFunctionName("Lbaas")
	data.SetUseName("lbaas")
	data.SetDescriptions("Comandos relacionados a Load Balancer", "Comandos para gerenciar load balancers")
	data.SetGroupID("products")
	data.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	// Adicionar imports necessários
	data.AddImport("mgccli/cmd/gen/lbaas/networkacls")
	data.AddImport("mgccli/cmd/gen/lbaas/networkbackends")
	data.AddImport("mgccli/cmd/gen/lbaas/networkbackendtarget")
	data.AddImport("mgccli/cmd/gen/lbaas/networkcertificates")
	data.AddImport("mgccli/cmd/gen/lbaas/networkhealthchecks")
	data.AddImport("mgccli/cmd/gen/lbaas/networklisteners")
	data.AddImport("mgccli/cmd/gen/lbaas/networkloadbalancers")
	data.AddImport("sdk \"github.com/MagaluCloud/mgc-sdk-go/client\"")
	data.AddImport("\"github.com/MagaluCloud/mgc-sdk-go/lbaas\"")
	data.AddImport("\"github.com/spf13/cobra\"")

	// Adicionar subcomandos
	data.AddSubCommand("networkacls", "NetworkACLs", "lbaasService.NetworkACLs()")
	data.AddSubCommand("networkbackends", "NetworkBackends", "lbaasService.NetworkBackends()")
	data.AddSubCommand("networkbackendtarget", "NetworkBackendTarget", "lbaasService.NetworkBackendTarget()")
	data.AddSubCommand("networkcertificates", "NetworkCertificates", "lbaasService.NetworkCertificates()")
	data.AddSubCommand("networkhealthchecks", "NetworkHealthChecks", "lbaasService.NetworkHealthChecks()")
	data.AddSubCommand("networklisteners", "NetworkListeners", "lbaasService.NetworkListeners()")
	data.AddSubCommand("networkloadbalancers", "NetworkLoadBalancers", "lbaasService.NetworkLoadBalancers()")

	return data
}

// PrintPackageGroupData exibe os dados do grupo de pacotes de forma legível
func PrintPackageGroupData(data *PackageGroupData) {
	fmt.Printf("=== Dados do Grupo de Pacotes: %s ===\n", data.PackageName)
	fmt.Printf("Pacote: %s\n", data.PackageName)
	fmt.Printf("Função: %s\n", data.FunctionName)
	fmt.Printf("Uso: %s\n", data.UseName)
	fmt.Printf("Descrição Curta: %s\n", data.ShortDescription)
	fmt.Printf("Descrição Longa: %s\n", data.LongDescription)
	if data.GroupID != "" {
		fmt.Printf("ID do Grupo: %s\n", data.GroupID)
	}
	fmt.Printf("Parâmetro do Serviço: %s\n", data.ServiceParam)

	fmt.Printf("\nImports (%d):\n", len(data.Imports))
	for _, imp := range data.Imports {
		fmt.Printf("  %s\n", imp)
	}

	fmt.Printf("\nSubcomandos (%d):\n", len(data.SubCommands))
	for i, subCmd := range data.SubCommands {
		fmt.Printf("  %d. %s.%sCmd(cmd, %s)\n", i+1, subCmd.PackageName, subCmd.FunctionName, subCmd.ServiceCall)
	}
	fmt.Println()
}

// GenerateExampleData gera dados de exemplo para todos os grupos principais
func GenerateExampleData() map[string]*PackageGroupData {
	examples := make(map[string]*PackageGroupData)

	examples["compute"] = ExampleComputePackageGroup()
	examples["kubernetes"] = ExampleKubernetesPackageGroup()
	examples["lbaas"] = ExampleLbaasPackageGroup()

	return examples
}

// PrintAllExamples exibe todos os exemplos de dados
func PrintAllExamples() {
	fmt.Println("=== EXEMPLOS DE DADOS PARA TEMPLATES ===")

	examples := GenerateExampleData()
	for name, data := range examples {
		fmt.Printf("\n--- %s ---\n", strings.ToUpper(name))
		PrintPackageGroupData(data)
	}
}
