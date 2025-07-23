package sdk_structure

// Service representa um serviço individual com seus métodos
type Service struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Interface   string             `json:"interface"`
	Methods     []Method           `json:"methods"`
	SubServices map[string]Service `json:"sub_services,omitempty"` // Para subserviços aninhados
}

// Parameter representa um parâmetro de método
type Parameter struct {
	Position    int                  `json:"position"`
	Name        string               `json:"name"`
	Type        string               `json:"type"`
	Description string               `json:"description"`
	IsPrimitive bool                 `json:"is_primitive"`
	IsPointer   bool                 `json:"is_pointer"`
	IsArray     bool                 `json:"is_array"`
	Struct      map[string]Parameter `json:"struct,omitempty"`
}

// Method representa um método de um serviço
type Method struct {
	Description string      `json:"description"`
	Name        string      `json:"name"`
	Parameters  []Parameter `json:"parameters"` // nome -> tipo
	Returns     []Parameter `json:"returns"`    // nome -> tipo
	Comments    string      `json:"comments"`
}

// Package representa um pacote do SDK com seus serviços
type Package struct {
	MenuName        string             `json:"menu_name"`
	Description     string             `json:"description"`
	LongDescription string             `json:"long_description"`
	Aliases         []string           `json:"aliases"`
	Name            string             `json:"name"`
	Services        []Service          `json:"services"`
	SubPkgs         map[string]Package `json:"sub_packages,omitempty"` // Para suporte recursivo
}

// SDKStructure representa a estrutura completa do SDK
type SDKStructure struct {
	Packages map[string]Package `json:"packages"`
}

// ClientMethod representa um método do cliente que retorna um serviço
type ClientMethod struct {
	Name            string
	ReturnType      string
	ServiceName     string
	LongDescription string
}
