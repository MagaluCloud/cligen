package cmdutils

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type Item struct {
	Title string
	Value string
}

type Selector interface {
	Run() ([]string, error)
}

type multiSelectModel struct {
	allItems     []Item
	filtered     []Item
	selected     map[string]bool
	filter       string
	cursor       int
	width        int
	height       int
	scrollOffset int
	selectText   string
}

type selectorImpl struct {
	model multiSelectModel
}

func MultiSelect(options []Item, selectText string) ([]string, error) {
	return NewMultiSelect(options, selectText).Run()
}

func NewMultiSelect(options []Item, selectText string) Selector {
	m := multiSelectModel{
		allItems:     options,
		filtered:     options,
		selected:     map[string]bool{},
		cursor:       0,
		width:        0,
		height:       0,
		scrollOffset: 0,
		selectText:   selectText,
	}

	return &selectorImpl{model: m}
}

func (s *selectorImpl) Run() ([]string, error) {
	p := tea.NewProgram(&s.model, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := final.(*multiSelectModel)

	var result []string
	for v, ok := range m.selected {
		if ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (m *multiSelectModel) Init() tea.Cmd {
	m.applyFilter()
	return nil
}

func (m *multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		// navegação
		switch key {
		case "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
			if m.cursor < m.scrollOffset {
				m.scrollOffset--
			}
			return m, nil

		case "down":
			m.cursor++
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			visible := m.height - 4
			if m.cursor >= m.scrollOffset+visible {
				m.scrollOffset++
			}
			return m, nil

		default:
			if len(msg.String()) == 1 {
				m.cursor = 0
				m.scrollOffset = 0
			}
		}

		// ações
		switch key {
		case "enter": // seleciona/desseleciona item
			if len(m.filtered) > 0 {
				item := m.filtered[m.cursor]
				m.selected[item.Value] = !m.selected[item.Value]
			}
			return m, nil

		case "tab": // confirmar seleção
			return m, tea.Quit

		case "right": // selecionar todos
			for _, it := range m.filtered {
				m.selected[it.Value] = true
			}
			return m, nil

		case "left": // desseleciona todos
			for _, it := range m.filtered {
				m.selected[it.Value] = false
			}
			return m, nil

		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.applyFilter()
			}
			return m, nil

		case "esc":
			m.filter = ""
			m.applyFilter()
			return m, nil

		case "ctrl+c":
			return m, tea.Quit
		}

		// filtro ao digitar
		if len(key) == 1 {
			m.filter += key
			m.applyFilter()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = max(msg.Height-5, 3)
		return m, nil
	}

	return m, nil
}

func (m *multiSelectModel) View() string {
	if m.width == 0 || m.height == 0 {
		// Ainda não recebemos WindowSizeMsg
		return "Carregando..."
	}

	visible := m.height - 4
	start := m.scrollOffset
	end := min(start+visible, len(m.filtered))

	var b strings.Builder

	b.WriteString(color.CyanString(m.selectText) + m.filter + "\n")

	magenta := color.New(color.FgMagenta).SprintFunc()
	boldMagenta := color.New(color.FgMagenta, color.Bold).SprintFunc()

	if len(m.filtered) == 0 {
		b.WriteString("Nenhum resultado.\n")
		return b.String()
	}

	for i := start; i < end; i++ {
		cursor := " "
		if i == m.cursor {
			cursor = boldMagenta(">")
		}

		check := color.RedString("✗")
		if m.selected[m.filtered[i].Value] {
			check = color.GreenString("✓")
		}

		b.WriteString(fmt.Sprintf("%s [%s ] %s\n", cursor, check, m.filtered[i].Title))
	}

	info := "\n" +
		magenta("enter: ") + boldMagenta("select") +
		magenta(" | tab: ") + boldMagenta("confirm") +
		magenta(" | left: ") + boldMagenta("none") +
		magenta(" | right: ") + boldMagenta("all") +
		magenta(" | type: ") + boldMagenta("to filter") + "\n"

	b.WriteString(info)

	return b.String()
}

func (m *multiSelectModel) applyFilter() {
	if m.filter == "" {
		m.filtered = m.allItems
		return
	}

	lower := strings.ToLower(m.filter)
	var result []Item

	for _, it := range m.allItems {
		if strings.Contains(strings.ToLower(it.Title), lower) {
			result = append(result, it)
		}
	}

	m.filtered = result

	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}
