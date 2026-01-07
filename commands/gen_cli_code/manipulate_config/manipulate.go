package manipulate_config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magaluCloud/cligen/config"
)

// StartServer inicia o servidor web para manipulação do config.json
func StartServer(port string) error {
	if port == "" {
		port = "9080"
	}

	// Configurar Gin em modo release para produção
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Obter o diretório base do projeto
	baseDir, err := getBaseDir()
	if err != nil {
		return fmt.Errorf("erro ao obter diretório base: %w", err)
	}

	staticDir := filepath.Join(baseDir, "commands", "gen_cli_code", "manipulate_config", "static")
	templatesDir := filepath.Join(baseDir, "commands", "gen_cli_code", "manipulate_config", "templates")

	// Servir arquivos estáticos (HTML, CSS, JS)
	r.Static("/static", staticDir)
	r.LoadHTMLGlob(filepath.Join(templatesDir, "*"))

	// Rota principal
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Config Manipulator",
		})
	})

	// API: Carregar config
	r.GET("/api/config", loadConfig)

	// API: Salvar config manipulado
	r.POST("/api/config/save", saveConfig)

	// API: Reordenar menus
	r.POST("/api/config/reorder", reorderConfig)

	// API: Criar novo menu
	r.POST("/api/menu/create", createMenu)

	// API: Remover menu/submenu
	r.DELETE("/api/menu/:id", deleteMenu)

	// API: Mover elemento
	r.POST("/api/menu/move", moveElement)

	// API: Atualizar menu/submenu/method
	r.PUT("/api/menu/:id", updateMenu)
	r.PUT("/api/method/:menuId/:submenuId/:methodIndex", updateMethod)

	// API: Recriar config.json
	r.POST("/api/config/regenerate", regenerateConfig)

	fmt.Printf("Servidor iniciado em http://localhost:%s\n", port)
	fmt.Println("Pressione Ctrl+C para parar o servidor")

	return r.Run(":" + port)
}

// loadConfig carrega o config.json ou config.json
func loadConfig(c *gin.Context) {

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, cfg)
}

// saveConfig salva o config manipulado como config.json
func saveConfig(c *gin.Context) {
	var cfg config.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar JSON: %v", err),
		})
		return
	}

	// Salvar como config.json
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao serializar config: %v", err),
		})
		return
	}

	configPath := filepath.Join("config", "config.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar arquivo: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config salvo com sucesso",
		"path":    configPath,
	})
}

// reorderConfig recebe a nova ordem e salva
func reorderConfig(c *gin.Context) {
	var cfg config.Config
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar JSON: %v", err),
		})
		return
	}

	// Salvar como config.json
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao serializar config: %v", err),
		})
		return
	}

	configPath := filepath.Join("config", "config.json")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar arquivo: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config reordenado e salvo com sucesso",
		"path":    configPath,
	})
}

// getBaseDir retorna o diretório raiz do projeto
func getBaseDir() (string, error) {
	// Tentar encontrar o diretório raiz procurando por go.mod
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Chegamos à raiz do sistema de arquivos
			break
		}
		dir = parent
	}

	// Se não encontrou, retorna o diretório atual
	return wd, nil
}

// createMenuRequest representa a requisição para criar um novo menu
type createMenuRequest struct {
	Name        string `json:"name" binding:"required"`
	SDKPackage  string `json:"sdk_package"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// createMenu cria um novo menu e adiciona ao config
func createMenu(c *gin.Context) {
	var req createMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	// Carregar config atual
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	// Criar novo menu
	newMenu := &config.Menu{
		ID:          uuid.New().String(),
		Name:        req.Name,
		SDKPackage:  req.SDKPackage,
		Description: req.Description,
		Enabled:     req.Enabled,
		Menus:       []*config.Menu{},
		Methods:     []*config.Method{},
	}

	// Adicionar ao config
	cfg.Menus = append(cfg.Menus, newMenu)

	// Salvar config atualizado
	err = cfg.SaveConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar config: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Menu criado com sucesso",
		"menu":    newMenu,
		"config":  cfg,
	})
}

// deleteMenuRequest representa a requisição para deletar um menu/submenu
type deleteMenuRequest struct {
	ID string `json:"id" binding:"required"`
}

// deleteMenu remove um menu ou submenu pelo ID
func deleteMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID é obrigatório",
		})
		return
	}

	// Carregar config atual
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	// Procurar e remover o menu/submenu
	removed := removeMenuByID(cfg, id)
	if !removed {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu ou submenu não encontrado",
		})
		return
	}

	// Salvar config atualizado
	err = cfg.SaveConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu removido com sucesso",
		"config":  cfg,
	})
}

// removeMenuByID remove um menu ou submenu pelo ID recursivamente
func removeMenuByID(cfg *config.Config, id string) bool {
	// Procurar nos menus principais
	for i, menu := range cfg.Menus {
		if menu.ID == id {
			cfg.Menus = append(cfg.Menus[:i], cfg.Menus[i+1:]...)
			return true
		}
		// Procurar recursivamente nos submenus
		if menu.Menus != nil && removeMenuFromMenus(&menu.Menus, id) {
			return true
		}
	}
	return false
}

// removeMenuFromMenus remove um menu de uma lista de menus recursivamente
func removeMenuFromMenus(menus *[]*config.Menu, id string) bool {
	if menus == nil {
		return false
	}
	for i, menu := range *menus {
		if menu.ID == id {
			// Remover do slice
			*menus = append((*menus)[:i], (*menus)[i+1:]...)
			return true
		}
		// Procurar recursivamente nos submenus
		if menu.Menus != nil && removeMenuFromMenus(&menu.Menus, id) {
			return true
		}
	}
	return false
}

// moveElementRequest representa a requisição para mover um elemento
type moveElementRequest struct {
	ElementID   string `json:"element_id" binding:"required"`
	TargetID    string `json:"target_id"`                       // ID do destino (menu ou submenu), vazio para raiz
	TargetType  string `json:"target_type"`                     // "menu", "submenu", "root"
	ElementType string `json:"element_type" binding:"required"` // "menu" ou "submenu"
}

// moveElement move um elemento para outro local
func moveElement(c *gin.Context) {
	var req moveElementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	// Encontrar o elemento a ser movido
	var elementToMove *config.Menu
	var parentMenus *[]*config.Menu
	var elementIndex int

	if req.ElementType == "menu" {
		// Procurar nos menus principais
		for i, menu := range cfg.Menus {
			if menu.ID == req.ElementID {
				elementToMove = menu
				parentMenus = &cfg.Menus
				elementIndex = i
				break
			}
		}
		// Se não encontrou, procurar em submenus
		if elementToMove == nil {
			elementToMove, parentMenus, elementIndex = findMenuInSubmenus(cfg.Menus, req.ElementID)
		}
	} else if req.ElementType == "submenu" {
		elementToMove, parentMenus, elementIndex = findMenuInSubmenus(cfg.Menus, req.ElementID)
	}

	if elementToMove == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Elemento não encontrado",
		})
		return
	}

	// Remover do local original
	if parentMenus != nil && elementIndex >= 0 && elementIndex < len(*parentMenus) {
		*parentMenus = append((*parentMenus)[:elementIndex], (*parentMenus)[elementIndex+1:]...)
	}

	// Adicionar ao destino e atualizar parent_menu_id
	if req.TargetType == "root" || req.TargetID == "" {
		// Mover para a raiz - limpar parent_menu_id e atualizar recursivamente todos os submenus filhos
		updateParentMenuIDRecursive(elementToMove, "")
		cfg.Menus = append(cfg.Menus, elementToMove)
	} else {
		// Encontrar o destino
		targetMenu := findMenuByID(cfg.Menus, req.TargetID)
		if targetMenu == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Destino não encontrado",
			})
			return
		}
		if targetMenu.Menus == nil {
			targetMenu.Menus = []*config.Menu{}
		}
		// Atualizar parent_menu_id do menu movido para o ID do menu pai e atualizar recursivamente todos os submenus filhos
		updateParentMenuIDRecursive(elementToMove, targetMenu.ID)
		targetMenu.Menus = append(targetMenu.Menus, elementToMove)
	}

	// Salvar config atualizado
	err = cfg.SaveConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Elemento movido com sucesso",
		"config":  cfg,
	})
}

// findMenuByID encontra um menu pelo ID recursivamente
func findMenuByID(menus []*config.Menu, id string) *config.Menu {
	for _, menu := range menus {
		if menu.ID == id {
			return menu
		}
		if menu.Menus != nil {
			if found := findMenuByID(menu.Menus, id); found != nil {
				return found
			}
		}
	}
	return nil
}

// findMenuInSubmenus encontra um menu em submenus e retorna o menu, o slice pai e o índice
func findMenuInSubmenus(menus []*config.Menu, id string) (*config.Menu, *[]*config.Menu, int) {
	for _, menu := range menus {
		if menu.Menus != nil {
			for i, submenu := range menu.Menus {
				if submenu.ID == id {
					return submenu, &menu.Menus, i
				}
				// Procurar recursivamente nos submenus aninhados
				if submenu.Menus != nil {
					if found, parent, idx := findMenuInSubmenus(submenu.Menus, id); found != nil {
						return found, parent, idx
					}
				}
			}
		}
	}
	return nil, nil, -1
}

// updateParentMenuIDRecursive atualiza recursivamente o parent_menu_id de um menu e todos os seus submenus
func updateParentMenuIDRecursive(menu *config.Menu, parentID string) {
	if menu == nil {
		return
	}
	// Atualizar o parent_menu_id do menu atual
	menu.ParentMenuID = parentID
	// Atualizar recursivamente todos os submenus filhos
	if menu.Menus != nil {
		for _, submenu := range menu.Menus {
			updateParentMenuIDRecursive(submenu, menu.ID)
		}
	}
}

// updateMenuRequest representa a requisição para atualizar um menu/submenu
type updateMenuRequest struct {
	Name             string   `json:"name,omitempty"`
	Enabled          *bool    `json:"enabled,omitempty"`
	Description      string   `json:"description,omitempty"`
	LongDescription  string   `json:"long_description,omitempty"`
	SDKPackage       string   `json:"sdk_package,omitempty"`
	CliGroup         string   `json:"cli_group,omitempty"`
	Alias            []string `json:"alias,omitempty"`
	ServiceInterface string   `json:"service_interface,omitempty"`
	SDKFile          string   `json:"sdk_file,omitempty"`
	CustomFile       string   `json:"custom_file,omitempty"`
	IsGroup          *bool    `json:"is_group,omitempty"`
}

// updateMenu atualiza um menu ou submenu
func updateMenu(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID é obrigatório",
		})
		return
	}

	var req updateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	// Debug: log do request recebido
	fmt.Printf("UpdateMenu - ID: %s, Name: %s, Enabled: %v, IsGroup: %v\n",
		id, req.Name, req.Enabled, req.IsGroup)

	// Carregar config atual

	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	// Encontrar o menu/submenu
	menu := findMenuByID(cfg.Menus, id)
	if menu == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu ou submenu não encontrado",
		})
		return
	}

	// Atualizar campos - sempre atualizar, mesmo que sejam valores vazios/false
	if req.Name != "" {
		menu.Name = req.Name
	}
	// Enabled sempre atualizar se fornecido
	if req.Enabled != nil {
		menu.Enabled = *req.Enabled
	}
	// Sempre atualizar strings (podem ser vazias)
	menu.Description = req.Description
	menu.LongDescription = req.LongDescription
	menu.SDKPackage = req.SDKPackage
	menu.CliGroup = req.CliGroup
	if req.Alias != nil {
		menu.Alias = req.Alias
	}
	menu.ServiceInterface = req.ServiceInterface
	menu.SDKFile = req.SDKFile
	menu.CustomFile = req.CustomFile
	// IsGroup sempre atualizar se fornecido (mesmo que seja false)
	if req.IsGroup != nil {
		menu.IsGroup = *req.IsGroup
	}

	// Salvar config atualizado
	err = cfg.SaveConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar config: %v", err),
		})
		return
	}

	// Log para debug - verificar se o campo foi atualizado
	fmt.Printf("Menu após atualização - ID: %s, Name: %s, Enabled: %v, IsGroup: %v\n",
		menu.ID, menu.Name, menu.Enabled, menu.IsGroup)

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu atualizado com sucesso",
		"menu":    menu,
		"config":  cfg,
	})
}

// updateMethodRequest representa a requisição para atualizar um method
type updateMethodRequest struct {
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	LongDescription string               `json:"long_description,omitempty"`
	Comments        string               `json:"comments,omitempty"`
	Parameters      []config.Parameter   `json:"parameters,omitempty"`
	Returns         []config.Parameter   `json:"returns,omitempty"`
	Confirmation    *config.Confirmation `json:"confirmation,omitempty"`
	IsService       *bool                `json:"is_service,omitempty"`
	ServiceImport   string               `json:"service_import,omitempty"`
	SDKFile         string               `json:"sdk_file,omitempty"`
	CustomFile      string               `json:"custom_file,omitempty"`
}

// updateMethod atualiza um method em um submenu
func updateMethod(c *gin.Context) {
	menuID := c.Param("menuId")
	submenuID := c.Param("submenuId")
	methodIndexStr := c.Param("methodIndex")

	if menuID == "" || submenuID == "" || methodIndexStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "IDs e índice são obrigatórios",
		})
		return
	}

	methodIndex, err := strconv.Atoi(methodIndexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Índice do método inválido",
		})
		return
	}

	var req updateMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	// Carregar config atual
	cfg, err := config.LoadConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao carregar config: %v", err),
		})
		return
	}

	// Encontrar o menu
	menu := findMenuByID(cfg.Menus, menuID)
	if menu == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Menu não encontrado",
		})
		return
	}

	// Se submenuID é igual a menuID, o method está diretamente no menu
	var targetMenu *config.Menu
	if submenuID == menuID {
		targetMenu = menu
	} else {
		// Caso contrário, procurar no submenu
		if menu.Menus != nil {
			for _, sm := range menu.Menus {
				if sm.ID == submenuID {
					targetMenu = sm
					break
				}
			}
		}
	}

	if targetMenu == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Submenu não encontrado",
		})
		return
	}

	if targetMenu.Methods == nil || methodIndex < 0 || methodIndex >= len(targetMenu.Methods) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Método não encontrado",
		})
		return
	}

	method := targetMenu.Methods[methodIndex]

	// Atualizar campos
	if req.Name != "" {
		method.Name = req.Name
	}
	method.Description = req.Description
	method.LongDescription = req.LongDescription
	method.Comments = req.Comments
	if req.Parameters != nil {
		method.Parameters = req.Parameters
	}
	if req.Returns != nil {
		method.Returns = req.Returns
	}
	if req.Confirmation != nil {
		method.Confirmation = req.Confirmation
	}
	if req.IsService != nil {
		method.IsService = *req.IsService
	}
	method.ServiceImport = req.ServiceImport
	method.SDKFile = req.SDKFile
	method.CustomFile = req.CustomFile

	// Salvar config atualizado
	err = cfg.SaveConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao salvar config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Método atualizado com sucesso",
		"method":  method,
		"config":  cfg,
	})
}

// regenerateConfig recria o config.json chamando genconfig.Run()
// Usa os/exec para evitar ciclo de importação
func regenerateConfig(c *gin.Context) {
	// Obter o diretório base do projeto
	baseDir, err := getBaseDir()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao obter diretório base: %v", err),
		})
		return
	}

	// Executar o comando generate-config via os/exec
	// Isso evita o ciclo de importação
	cmd := exec.Command("go", "run", ".", "generate-config")
	cmd.Dir = baseDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao executar generate-config: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config.json recriado com sucesso",
	})
}
