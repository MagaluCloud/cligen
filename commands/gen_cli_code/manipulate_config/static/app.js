let configData = null;

// Sistema de undo/redo
let history = []; // Array de estados salvos
let currentHistoryIndex = -1; // Índice atual no histórico
const MAX_HISTORY_SIZE = 50; // Limite máximo de histórico

// Carregar configuração ao iniciar
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    setupKeyboardShortcuts();
    loadConfig();
});

function setupEventListeners() {
    document.getElementById('saveBtn').addEventListener('click', saveConfig);
    document.getElementById('reloadBtn').addEventListener('click', loadConfig);
    document.getElementById('regenerateBtn').addEventListener('click', regenerateConfig);
    document.getElementById('createMenuBtn').addEventListener('click', openCreateMenuModal);
    document.getElementById('cancelCreateMenuBtn').addEventListener('click', closeCreateMenuModal);
    document.getElementById('submitCreateMenuBtn').addEventListener('click', submitCreateMenu);
    document.getElementById('undoBtn').addEventListener('click', undo);
    document.getElementById('redoBtn').addEventListener('click', redo);

    // Fechar modal de criar menu ao clicar no X
    const createModal = document.getElementById('createMenuModal');
    const closeCreateBtn = createModal.querySelector('.close');
    if (closeCreateBtn) {
        closeCreateBtn.addEventListener('click', closeCreateMenuModal);
    }

    // Fechar modal de criar menu ao clicar fora dele
    window.addEventListener('click', (event) => {
        if (event.target === createModal) {
            closeCreateMenuModal();
        }
    });

    // Event listeners para modal de edição
    const editModal = document.getElementById('editModal');
    if (editModal) {
        const closeEditBtn = editModal.querySelector('#closeEditModal');
        const cancelEditBtn = document.getElementById('cancelEditBtn');
        const submitEditBtn = document.getElementById('submitEditBtn');

        if (closeEditBtn) {
            closeEditBtn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                closeEditModal();
            });
        }

        if (cancelEditBtn) {
            cancelEditBtn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                closeEditModal();
            });
        }

        if (submitEditBtn) {
            submitEditBtn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                submitEdit();
            });
        }

        // Fechar modal de edição ao clicar fora dele
        editModal.addEventListener('click', (event) => {
            if (event.target === editModal) {
                closeEditModal();
            }
        });
    }
}

async function loadConfig() {
    const loading = document.getElementById('loading');
    const configTree = document.getElementById('configTree');

    loading.style.display = 'block';
    configTree.style.display = 'none';
    hideStatus();

    try {
        const response = await fetch('/api/config');
        if (!response.ok) {
            throw new Error('Erro ao carregar configuração');
        }

        configData = await response.json();
        // Limpar histórico e salvar estado inicial
        history = [];
        currentHistoryIndex = -1;
        saveState();
        renderConfig(configData);

        loading.style.display = 'none';
        configTree.style.display = 'block';
    } catch (error) {
        loading.textContent = `Erro: ${error.message}`;
        loading.style.color = '#dc3545';
        showStatus('Erro ao carregar configuração', 'error');
    }
}

// Funções de undo/redo
function saveState() {
    if (!configData) return;

    // Criar uma cópia profunda do estado atual
    const stateCopy = JSON.parse(JSON.stringify(configData));

    // Se estamos no meio do histórico (não no final), remover estados futuros
    if (currentHistoryIndex < history.length - 1) {
        history = history.slice(0, currentHistoryIndex + 1);
    }

    // Adicionar novo estado ao histórico
    history.push(stateCopy);
    currentHistoryIndex = history.length - 1;

    // Limitar tamanho do histórico
    if (history.length > MAX_HISTORY_SIZE) {
        history.shift();
        currentHistoryIndex--;
    }

    updateUndoRedoButtons();
}

function undo() {
    if (currentHistoryIndex <= 0) {
        showStatus('Não há mais ações para desfazer', 'info');
        return;
    }

    currentHistoryIndex--;
    configData = JSON.parse(JSON.stringify(history[currentHistoryIndex]));
    renderConfig(configData);
    updateUndoRedoButtons();
    showStatus('Ação desfeita', 'success');
}

function redo() {
    if (currentHistoryIndex >= history.length - 1) {
        showStatus('Não há mais ações para refazer', 'info');
        return;
    }

    currentHistoryIndex++;
    configData = JSON.parse(JSON.stringify(history[currentHistoryIndex]));
    renderConfig(configData);
    updateUndoRedoButtons();
    showStatus('Ação refeita', 'success');
}

function updateUndoRedoButtons() {
    const undoBtn = document.getElementById('undoBtn');
    const redoBtn = document.getElementById('redoBtn');

    if (undoBtn) {
        undoBtn.disabled = currentHistoryIndex <= 0;
        undoBtn.title = currentHistoryIndex <= 0 ? 'Não há ações para desfazer' : 'Desfazer (Ctrl+Z)';
    }

    if (redoBtn) {
        redoBtn.disabled = currentHistoryIndex >= history.length - 1;
        redoBtn.title = currentHistoryIndex >= history.length - 1 ? 'Não há ações para refazer' : 'Refazer (Ctrl+Y)';
    }
}

// Configurar atalhos de teclado para undo/redo
function setupKeyboardShortcuts() {
    document.addEventListener('keydown', (e) => {
        // Ignorar se estiver digitando em um input ou textarea
        if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
            // Permitir Ctrl+Z apenas para desfazer texto, mas não para undo global
            // Permitir Ctrl+Y para redo global mesmo em inputs
            if ((e.ctrlKey || e.metaKey) && e.key === 'y') {
                e.preventDefault();
                redo();
            }
            return;
        }

        // Ctrl+Z ou Cmd+Z para undo
        if ((e.ctrlKey || e.metaKey) && e.key === 'z' && !e.shiftKey) {
            e.preventDefault();
            undo();
        }
        // Ctrl+Y ou Ctrl+Shift+Z para redo
        if ((e.ctrlKey || e.metaKey) && (e.key === 'y' || (e.key === 'z' && e.shiftKey))) {
            e.preventDefault();
            redo();
        }
    });
}

function renderConfig(config) {
    const container = document.getElementById('menusContainer');
    container.innerHTML = '';

    if (!config.menus || config.menus.length === 0) {
        container.innerHTML = '<div class="empty-message">Nenhum menu encontrado</div>';
        return;
    }

    config.menus.forEach((menu, menuIndex) => {
        const menuBlock = createMenuBlock(menu, menuIndex);
        container.appendChild(menuBlock);
    });

    setupDragAndDrop();
    setupExpandCollapse();
}

function createMenuBlock(menu, menuIndex) {
    const menuDiv = document.createElement('div');
    menuDiv.className = 'menu-block collapsed';
    menuDiv.draggable = true;
    menuDiv.dataset.type = 'menu';
    menuDiv.dataset.index = menuIndex;
    menuDiv.dataset.id = menu.id || '';

    const header = document.createElement('div');
    header.className = 'menu-header';
    // Header não deve ser arrastável
    header.draggable = false;

    const expandBtn = document.createElement('span');
    expandBtn.className = 'expand-icon';
    expandBtn.innerHTML = '▶';
    expandBtn.style.cursor = 'pointer';
    expandBtn.draggable = false;
    expandBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        toggleExpand(menuDiv);
        return false;
    };

    header.onclick = (e) => {
        // Só expandir se não clicou no ícone (que já tem seu próprio handler)
        if (!e.target.classList.contains('expand-icon')) {
            e.stopPropagation();
            e.preventDefault();
            e.stopImmediatePropagation();
            toggleExpand(menuDiv);
            return false;
        }
    };

    const headerText = document.createElement('span');
    headerText.className = 'header-text';
    headerText.textContent = `Menu: ${menu.cli_name}`;

    const editBtn = document.createElement('button');
    editBtn.className = 'edit-btn';
    editBtn.innerHTML = '✏️';
    editBtn.title = 'Editar menu';
    editBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        openEditModal('menu', menu, menuIndex);
        return false;
    };

    const deleteBtn = document.createElement('button');
    deleteBtn.className = 'delete-btn';
    deleteBtn.innerHTML = '🗑️';
    deleteBtn.title = 'Remover menu';
    deleteBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        if (confirm(`Tem certeza que deseja remover o menu "${menu.cli_name}"?`)) {
            deleteMenu(menu.id);
        }
        return false;
    };

    header.appendChild(expandBtn);
    header.appendChild(headerText);
    header.appendChild(editBtn);
    header.appendChild(deleteBtn);
    menuDiv.appendChild(header);

    const content = document.createElement('div');
    content.className = 'menu-content';
    // CSS cuida do display inicial (collapsed)

    // Processar submenus primeiro
    if (menu.menus && menu.menus.length > 0) {
        menu.menus.forEach((submenu, submenuIndex) => {
            const submenuBlock = createSubmenuBlock(submenu, menuIndex, submenuIndex);
            content.appendChild(submenuBlock);
        });
    }

    // Processar métodos diretamente no menu (menus de nível superior podem ter methods)
    if (menu.methods && menu.methods.length > 0) {
        const methodsContainer = document.createElement('div');
        methodsContainer.className = 'methods-container';

        menu.methods.forEach((method, methodIndex) => {
            const methodItem = createMethodItem(method, menuIndex, null, methodIndex, null);
            methodsContainer.appendChild(methodItem);
        });

        content.appendChild(methodsContainer);
    }

    menuDiv.appendChild(content);
    return menuDiv;
}

function createSubmenuBlock(submenu, menuIndex, submenuIndex, parentSubmenuIndex = null) {
    const submenuDiv = document.createElement('div');
    submenuDiv.className = 'submenu-block collapsed';
    submenuDiv.draggable = true;
    submenuDiv.dataset.type = 'submenu';
    submenuDiv.dataset.menuIndex = menuIndex;
    submenuDiv.dataset.index = submenuIndex;
    submenuDiv.dataset.id = submenu.id || '';
    if (parentSubmenuIndex !== null) {
        submenuDiv.dataset.parentSubmenuIndex = parentSubmenuIndex;
    }

    const header = document.createElement('div');
    header.className = 'submenu-header';
    // Header não deve ser arrastável
    header.draggable = false;

    const expandBtn = document.createElement('span');
    expandBtn.className = 'expand-icon';
    expandBtn.innerHTML = '▶';
    expandBtn.style.cursor = 'pointer';
    expandBtn.draggable = false;
    expandBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        toggleExpand(submenuDiv);
        return false;
    };

    header.onclick = (e) => {
        // Só expandir se não clicou no ícone (que já tem seu próprio handler)
        if (!e.target.classList.contains('expand-icon')) {
            e.stopPropagation();
            e.preventDefault();
            e.stopImmediatePropagation();
            toggleExpand(submenuDiv);
            return false;
        }
    };

    const headerText = document.createElement('span');
    headerText.className = 'header-text';
    headerText.textContent = `SubMenu: ${submenu.cli_name}`;

    const editBtn = document.createElement('button');
    editBtn.className = 'edit-btn';
    editBtn.innerHTML = '✏️';
    editBtn.title = 'Editar submenu';
    editBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        openEditModal('submenu', submenu, menuIndex, submenuIndex);
        return false;
    };

    const promoteBtn = document.createElement('button');
    promoteBtn.className = 'promote-btn';
    promoteBtn.innerHTML = '⬆️';
    promoteBtn.title = 'Promover para o mesmo nível do parent';
    promoteBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        promoteSubmenuToMenu(submenu.id);
        return false;
    };

    const deleteBtn = document.createElement('button');
    deleteBtn.className = 'delete-btn';
    deleteBtn.innerHTML = '🗑️';
    deleteBtn.title = 'Remover submenu';
    deleteBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        if (confirm(`Tem certeza que deseja remover o submenu "${submenu.cli_name}"?`)) {
            deleteMenu(submenu.id);
        }
        return false;
    };

    header.appendChild(expandBtn);
    header.appendChild(headerText);
    header.appendChild(editBtn);
    header.appendChild(promoteBtn);
    header.appendChild(deleteBtn);
    submenuDiv.appendChild(header);

    const content = document.createElement('div');
    content.className = 'submenu-content';
    // CSS cuida do display inicial (collapsed)

    // Processar submenus aninhados primeiro
    if (submenu.menus && submenu.menus.length > 0) {
        submenu.menus.forEach((nestedSubmenu, nestedIndex) => {
            const nestedSubmenuBlock = createSubmenuBlock(nestedSubmenu, menuIndex, nestedIndex, submenuIndex);
            content.appendChild(nestedSubmenuBlock);
        });
    }

    // Processar métodos
    if (submenu.methods && submenu.methods.length > 0) {
        const methodsContainer = document.createElement('div');
        methodsContainer.className = 'methods-container';

        submenu.methods.forEach((method, methodIndex) => {
            const methodItem = createMethodItem(method, menuIndex, submenuIndex, methodIndex, parentSubmenuIndex);
            methodsContainer.appendChild(methodItem);
        });

        content.appendChild(methodsContainer);
    }

    submenuDiv.appendChild(content);
    return submenuDiv;
}

function createMethodItem(method, menuIndex, submenuIndex, methodIndex, parentSubmenuIndex = null) {
    const methodDiv = document.createElement('div');
    methodDiv.className = 'method-item';
    methodDiv.draggable = true;
    methodDiv.dataset.type = 'method';
    methodDiv.dataset.menuIndex = menuIndex;
    methodDiv.dataset.submenuIndex = submenuIndex;
    methodDiv.dataset.index = methodIndex;
    if (parentSubmenuIndex !== null) {
        methodDiv.dataset.parentSubmenuIndex = parentSubmenuIndex;
    }

    const methodName = document.createElement('span');
    methodName.textContent = `Method: ${method.name}`;

    const editBtn = document.createElement('button');
    editBtn.className = 'edit-btn-small';
    editBtn.innerHTML = '✏️';
    editBtn.title = 'Editar método';
    editBtn.onclick = (e) => {
        e.stopPropagation();
        e.preventDefault();
        e.stopImmediatePropagation();
        // Encontrar menu e submenu pelos índices
        if (configData && configData.menus && configData.menus[menuIndex]) {
            const menu = configData.menus[menuIndex];
            // Se submenuIndex é null, o method está diretamente no menu
            if (submenuIndex === null || submenuIndex === undefined) {
                // Method está diretamente no menu de nível superior
                // Usar o próprio menu como "submenu" para compatibilidade com o backend
                openEditModal('method', method, menuIndex, null, methodIndex, menu.id, menu.id);
            } else {
                const submenu = menu?.menus?.[submenuIndex];
                if (submenu) {
                    openEditModal('method', method, menuIndex, submenuIndex, methodIndex, menu.id, submenu.id);
                }
            }
        }
        return false;
    };

    methodDiv.appendChild(methodName);
    methodDiv.appendChild(editBtn);
    return methodDiv;
}

function toggleExpand(element) {
    if (!element) {
        console.error('toggleExpand: element is null');
        return;
    }

    const content = element.querySelector('.menu-content, .submenu-content');
    const expandIcon = element.querySelector('.expand-icon');

    if (!content) {
        console.error('toggleExpand: content not found', element);
        return;
    }

    if (!expandIcon) {
        console.error('toggleExpand: expandIcon not found', element);
        return;
    }

    // Verificar se está collapsed pela classe
    const isCollapsed = element.classList.contains('collapsed');

    console.log('toggleExpand:', {
        element: element.className,
        isCollapsed: isCollapsed,
        hasContent: !!content,
        hasIcon: !!expandIcon
    });

    if (isCollapsed) {
        // Expandir
        element.classList.remove('collapsed');
        expandIcon.innerHTML = '▼';
        console.log('Expanded');
    } else {
        // Colapsar
        element.classList.add('collapsed');
        expandIcon.innerHTML = '▶';
        console.log('Collapsed');
    }
}

function setupExpandCollapse() {
    // Expand/collapse já está configurado nos botões criados
}

function setupDragAndDrop() {
    // Adicionar listeners apenas para elementos que devem ser arrastáveis
    // Headers não devem ser arrastáveis
    document.querySelectorAll('.menu-block, .submenu-block, .method-item').forEach(element => {
        // Não adicionar drag listeners em headers
        if (element.classList.contains('menu-header') || element.classList.contains('submenu-header')) {
            return;
        }

        element.addEventListener('dragstart', handleDragStart);
        element.addEventListener('dragend', handleDragEnd);
        element.addEventListener('dragover', handleDragOver);
        element.addEventListener('drop', handleDrop);
        element.addEventListener('dragenter', handleDragEnter);
        element.addEventListener('dragleave', handleDragLeave);
    });

    // Permitir drop nos containers também
    document.querySelectorAll('.menu-content, .submenu-content, .menus-container, .methods-container').forEach(container => {
        container.addEventListener('dragover', handleDragOver);
        container.addEventListener('drop', handleContainerDrop);
        container.addEventListener('dragenter', handleDragEnter);
        container.addEventListener('dragleave', handleDragLeave);
    });
}

let draggedElement = null;
let draggedData = null;

function handleDragStart(e) {
    draggedElement = this;
    draggedData = {
        type: this.dataset.type,
        id: this.dataset.id || '',
        menuIndex: this.dataset.menuIndex ? parseInt(this.dataset.menuIndex) : null,
        submenuIndex: this.dataset.submenuIndex ? parseInt(this.dataset.submenuIndex) : null,
        index: this.dataset.index ? parseInt(this.dataset.index) : null,
        parentSubmenuIndex: this.dataset.parentSubmenuIndex ? parseInt(this.dataset.parentSubmenuIndex) : null
    };

    this.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.dropEffect = 'move';
    // Usar text/plain para compatibilidade
    e.dataTransfer.setData('text/plain', JSON.stringify(draggedData));

    // Criar uma imagem fantasma personalizada
    const dragImage = this.cloneNode(true);
    dragImage.style.opacity = '0.5';
    dragImage.style.position = 'absolute';
    dragImage.style.top = '-1000px';
    document.body.appendChild(dragImage);
    e.dataTransfer.setDragImage(dragImage, 0, 0);
    setTimeout(() => document.body.removeChild(dragImage), 0);
}

function handleDragEnd(e) {
    this.classList.remove('dragging');
    document.querySelectorAll('.drag-over').forEach(el => {
        el.classList.remove('drag-over');
    });
    draggedElement = null;
    draggedData = null;
}

function handleDragEnter(e) {
    e.preventDefault();
    e.stopPropagation();

    // Não destacar se for o próprio elemento ou seu pai direto
    if (this === draggedElement || this.contains(draggedElement)) {
        return;
    }

    // Verificar se é um drop válido
    const isValidDrop = isValidDropTarget(this, draggedData);
    if (isValidDrop) {
        this.classList.add('drag-over');
    }
}

function handleDragLeave(e) {
    // Só remove se realmente saiu do elemento (não apenas de um filho)
    if (!this.contains(e.relatedTarget)) {
        this.classList.remove('drag-over');
    }
}

function handleDragOver(e) {
    e.preventDefault();
    e.stopPropagation();

    // Verificar se é um drop válido
    const isValidDrop = isValidDropTarget(this, draggedData);
    e.dataTransfer.dropEffect = isValidDrop ? 'move' : 'none';

    return false;
}

function isValidDropTarget(target, draggedData) {
    if (!draggedData) return false;

    // Não permitir drop no próprio elemento ou em seus filhos
    if (draggedElement && (target === draggedElement || target.contains(draggedElement))) {
        return false;
    }

    const targetType = target.dataset.type;
    const isContainer = target.classList.contains('menu-content') ||
        target.classList.contains('submenu-content') ||
        target.classList.contains('menus-container') ||
        target.classList.contains('methods-container');

    // Menus podem ser soltos em outros menus ou no container raiz
    if (draggedData.type === 'menu') {
        return targetType === 'menu' || target.id === 'menusContainer' ||
            target.classList.contains('menu-content');
    }

    // Submenus podem ser soltos em menus, outros submenus, containers ou na raiz
    if (draggedData.type === 'submenu') {
        // Verificar se não está tentando mover para dentro de si mesmo
        if (targetType === 'submenu' && draggedData.id === target.dataset.id) {
            return false;
        }
        // Verificar se não está tentando mover para dentro de um de seus próprios filhos
        if (draggedElement) {
            const targetSubmenuBlock = target.closest('.submenu-block');
            if (targetSubmenuBlock && draggedElement.contains(targetSubmenuBlock)) {
                return false;
            }
        }
        return targetType === 'menu' || targetType === 'submenu' ||
            target.id === 'menusContainer' ||
            target.classList.contains('menu-content') ||
            target.classList.contains('submenu-content');
    }

    // Métodos podem ser soltos em outros métodos ou no container de métodos
    if (draggedData.type === 'method') {
        return targetType === 'method' || target.classList.contains('methods-container');
    }

    return false;
}

function handleDrop(e) {
    e.preventDefault();
    e.stopPropagation();

    if (!draggedElement || !draggedData) return false;

    const targetType = this.dataset.type;
    const targetID = this.dataset.id || '';

    if (draggedData.type === 'menu' || draggedData.type === 'submenu') {
        moveElementToBackend(draggedData.id, draggedData.type, targetID, targetType);
    } else if (draggedData.type === 'method') {
        // Reordenar métodos localmente
        const targetMethodIndex = this.dataset.index ? parseInt(this.dataset.index) : null;
        const targetMenuIndex = this.dataset.menuIndex ? parseInt(this.dataset.menuIndex) : null;
        const targetSubmenuIndex = this.dataset.submenuIndex !== undefined ? (this.dataset.submenuIndex ? parseInt(this.dataset.submenuIndex) : null) : null;

        if (targetMethodIndex !== null && targetMenuIndex !== null &&
            draggedData.menuIndex !== null && draggedData.index !== null) {
            // Verificar se é o mesmo menu/submenu (reordenação dentro do mesmo container)
            const sameMenu = draggedData.menuIndex === targetMenuIndex;
            const sameSubmenu = (draggedData.submenuIndex === null || draggedData.submenuIndex === undefined) === (targetSubmenuIndex === null || targetSubmenuIndex === undefined) &&
                (draggedData.submenuIndex === targetSubmenuIndex || (draggedData.submenuIndex === null && targetSubmenuIndex === null));

            if (sameMenu && sameSubmenu) {
                // Salvar estado antes de modificar
                saveState();
                reorderMethods(draggedData.menuIndex, draggedData.submenuIndex, draggedData.index, targetMethodIndex);
                showStatus('Método reordenado localmente! (Clique em Salvar para persistir)', 'success');
                renderConfig(configData);
            }
        }
    }

    this.classList.remove('drag-over');
    return false;
}

function handleContainerDrop(e) {
    e.preventDefault();
    e.stopPropagation();

    if (!draggedElement || !draggedData) return false;

    const container = this;
    const isMenuContent = container.classList.contains('menu-content');
    const isSubmenuContent = container.classList.contains('submenu-content');
    const isMenusContainer = container.id === 'menusContainer';
    const isMethodsContainer = container.classList.contains('methods-container');

    if ((draggedData.type === 'menu' || draggedData.type === 'submenu')) {
        let targetID = '';
        let targetType = 'root';

        if (isMenusContainer) {
            // Mover para a raiz
            targetType = 'root';
            targetID = '';
        } else if (isMenuContent) {
            // Mover para dentro de um menu
            const menuBlock = container.closest('.menu-block');
            if (menuBlock) {
                targetID = menuBlock.dataset.id || '';
                targetType = 'menu';
            }
        } else if (isSubmenuContent) {
            // Mover para dentro de um submenu (suporta submenus aninhados)
            const submenuBlock = container.closest('.submenu-block');
            if (submenuBlock) {
                targetID = submenuBlock.dataset.id || '';
                targetType = 'submenu';
            }
        }

        // Permitir mover mesmo se targetID for igual (pode ser reordenação)
        // Mas verificar se não está tentando mover para dentro de si mesmo
        if (targetID !== draggedData.id || targetType === 'root') {
            moveElementToBackend(draggedData.id, draggedData.type, targetID, targetType);
        }
    } else if (isMethodsContainer && draggedData.type === 'method') {
        // Reordenar métodos localmente - encontrar o container (menu ou submenu)
        const submenuBlock = container.closest('.submenu-block');
        const menuBlock = container.closest('.menu-block');

        let targetMenuIndex = null;
        let targetSubmenuIndex = null;

        if (submenuBlock) {
            // Methods estão em um submenu
            targetMenuIndex = parseInt(submenuBlock.dataset.menuIndex);
            targetSubmenuIndex = parseInt(submenuBlock.dataset.index);
        } else if (menuBlock) {
            // Methods estão diretamente no menu de nível superior
            targetMenuIndex = parseInt(menuBlock.dataset.index);
            targetSubmenuIndex = null;
        }

        if (targetMenuIndex !== null && draggedData.menuIndex !== null && draggedData.index !== null) {
            // Verificar se é o mesmo container (menu ou submenu)
            const sameMenu = draggedData.menuIndex === targetMenuIndex;
            const sameSubmenu = (draggedData.submenuIndex === null || draggedData.submenuIndex === undefined) === (targetSubmenuIndex === null || targetSubmenuIndex === undefined) &&
                (draggedData.submenuIndex === targetSubmenuIndex || (draggedData.submenuIndex === null && targetSubmenuIndex === null));

            if (sameMenu && sameSubmenu) {
                // Encontrar o índice do método alvo baseado na posição do drop
                const methods = container.querySelectorAll('.method-item');
                let targetIndex = methods.length;

                // Tentar encontrar a posição exata baseada na posição do mouse
                for (let i = 0; i < methods.length; i++) {
                    const rect = methods[i].getBoundingClientRect();
                    if (e.clientY < rect.top + rect.height / 2) {
                        targetIndex = i;
                        break;
                    }
                }

                // Ajustar índice se necessário
                if (draggedData.index < targetIndex) {
                    targetIndex--;
                }

                // Salvar estado antes de modificar
                saveState();
                reorderMethods(draggedData.menuIndex, draggedData.submenuIndex, draggedData.index, targetIndex);
                showStatus('Método reordenado localmente! (Clique em Salvar para persistir)', 'success');
                renderConfig(configData);
            }
        }
    }

    this.classList.remove('drag-over');
    return false;
}

function moveMenuToMenu(fromIndex, toMenuIndex, toSubmenuIndex) {
    if (fromIndex === toMenuIndex) return;

    const menu = configData.menus[fromIndex];
    if (!menu) return;

    // Ajustar índices se necessário
    let adjustedFromIndex = fromIndex;
    let adjustedToIndex = toMenuIndex;
    if (fromIndex < toMenuIndex) {
        adjustedToIndex = toMenuIndex - 1;
    }

    // Remover do local original
    const [movedMenu] = configData.menus.splice(adjustedFromIndex, 1);

    if (toSubmenuIndex !== null && toSubmenuIndex !== undefined) {
        // Mover para dentro de um submenu
        const targetMenu = configData.menus[adjustedToIndex];
        if (targetMenu && targetMenu.menus && targetMenu.menus[toSubmenuIndex]) {
            const targetSubmenu = targetMenu.menus[toSubmenuIndex];
            if (!targetSubmenu.menus) {
                targetSubmenu.menus = [];
            }
            targetSubmenu.menus.push(movedMenu);
        }
    } else if (adjustedToIndex !== null && adjustedToIndex !== undefined && adjustedToIndex >= 0) {
        // Mover para dentro de um menu (transformar em submenu)
        const targetMenu = configData.menus[adjustedToIndex];
        if (targetMenu) {
            if (!targetMenu.menus) {
                targetMenu.menus = [];
            }
            targetMenu.menus.push(movedMenu);
        }
    } else {
        // Mover para o nível raiz (não deveria acontecer aqui, mas por segurança)
        configData.menus.push(movedMenu);
    }
}

function moveSubmenuToMenu(fromMenuIndex, fromSubmenuIndex, toMenuIndex) {
    const sourceMenu = configData.menus[fromMenuIndex];
    if (!sourceMenu || !sourceMenu.menus || fromSubmenuIndex < 0 || fromSubmenuIndex >= sourceMenu.menus.length) return;

    const submenu = sourceMenu.menus[fromSubmenuIndex];

    // Remover do menu original
    sourceMenu.menus.splice(fromSubmenuIndex, 1);

    // Adicionar ao menu destino
    const targetMenu = configData.menus[toMenuIndex];
    if (targetMenu) {
        if (!targetMenu.menus) {
            targetMenu.menus = [];
        }
        targetMenu.menus.push(submenu);
    }
}

function moveSubmenuToSubmenu(fromMenuIndex, fromSubmenuIndex, toMenuIndex, toSubmenuIndex) {
    const sourceMenu = configData.menus[fromMenuIndex];
    if (!sourceMenu || !sourceMenu.menus || !sourceMenu.menus[fromSubmenuIndex]) return;

    const submenu = sourceMenu.menus[fromSubmenuIndex];
    sourceMenu.menus.splice(fromSubmenuIndex, 1);

    const targetMenu = configData.menus[toMenuIndex];
    if (targetMenu && targetMenu.menus && targetMenu.menus[toSubmenuIndex]) {
        const targetSubmenu = targetMenu.menus[toSubmenuIndex];
        if (!targetSubmenu.menus) {
            targetSubmenu.menus = [];
        }
        targetSubmenu.menus.push(submenu);
    }
}

function reorderMenus(fromIndex, toIndex) {
    if (fromIndex === toIndex) return;
    const menus = configData.menus;
    const [moved] = menus.splice(fromIndex, 1);
    menus.splice(toIndex, 0, moved);
}

function reorderSubmenus(menuIndex, fromIndex, toIndex) {
    const menu = configData.menus[menuIndex];
    if (menu && menu.menus) {
        const [moved] = menu.menus.splice(fromIndex, 1);
        menu.menus.splice(toIndex, 0, moved);
    }
}

function reorderMethods(menuIndex, submenuIndex, fromIndex, toIndex) {
    const menu = configData.menus[menuIndex];
    if (!menu) return;

    // Se submenuIndex é null, o method está diretamente no menu
    if (submenuIndex === null || submenuIndex === undefined) {
        if (menu.methods) {
            const [moved] = menu.methods.splice(fromIndex, 1);
            menu.methods.splice(toIndex, 0, moved);
        }
    } else {
        // Caso contrário, procurar no submenu
        if (menu.menus && menu.menus[submenuIndex]) {
            const submenu = menu.menus[submenuIndex];
            if (submenu && submenu.methods) {
                const [moved] = submenu.methods.splice(fromIndex, 1);
                submenu.methods.splice(toIndex, 0, moved);
            }
        }
    }
}

async function saveConfig() {
    if (!configData) {
        showStatus('Nenhuma configuração carregada', 'error');
        return;
    }

    try {
        const response = await fetch('/api/config/save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(configData),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Erro ao salvar configuração');
        }

        const result = await response.json();
        showStatus(`Config salvo com sucesso em ${result.path}`, 'success');
    } catch (error) {
        showStatus(`Erro ao salvar: ${error.message}`, 'error');
    }
}

async function regenerateConfig() {
    if (!confirm('Tem certeza que deseja recriar o config.json? Isso irá sobrescrever o arquivo atual.')) {
        return;
    }

    const loading = document.getElementById('loading');
    const configTree = document.getElementById('configTree');

    loading.style.display = 'block';
    configTree.style.display = 'none';
    hideStatus();

    try {
        const response = await fetch('/api/config/regenerate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Erro ao recriar config.json');
        }

        const result = await response.json();
        showStatus(result.message || 'Config.json recriado com sucesso', 'success');

        // Recarregar o config após a regeneração
        await loadConfig();
    } catch (error) {
        loading.style.display = 'none';
        configTree.style.display = 'block';
        showStatus(`Erro ao recriar config.json: ${error.message}`, 'error');
    }
}

function showStatus(message, type) {
    const status = document.getElementById('status');
    status.textContent = message;
    status.className = `status ${type}`;

    if (type === 'success' || type === 'info') {
        setTimeout(() => hideStatus(), 3000);
    }
}

function hideStatus() {
    const status = document.getElementById('status');
    status.className = 'status';
}

function openCreateMenuModal() {
    const modal = document.getElementById('createMenuModal');
    modal.style.display = 'block';
    // Limpar formulário
    document.getElementById('createMenuForm').reset();
    document.getElementById('menuEnabled').checked = true;
}

function closeCreateMenuModal() {
    const modal = document.getElementById('createMenuModal');
    modal.style.display = 'none';
}

function openEditModal(type, element, menuIndex, submenuIndex = null, methodIndex = null, menuId = null, submenuId = null) {
    const modal = document.getElementById('editModal');
    const title = document.getElementById('editModalTitle');
    const formFields = document.getElementById('editFormFields');

    // Limpar campos anteriores
    formFields.innerHTML = '';

    // Configurar campos ocultos
    document.getElementById('editElementType').value = type;

    if (type === 'menu' || type === 'submenu') {
        document.getElementById('editElementId').value = element.id || '';
        title.textContent = `Editar ${type === 'menu' ? 'Menu' : 'SubMenu'}: ${element.cli_name}`;

        // Campos para Menu/SubMenu
        addFormField(formFields, 'cli_name', 'Nome do comando', element.cli_name || '', true);
        // Sempre criar campo enabled, mesmo se não existir no elemento
        const enabledValue = element.hasOwnProperty('enabled') ? element.enabled : true;
        addFormField(formFields, 'enabled', 'Habilitado', enabledValue, false, 'checkbox');
        addFormField(formFields, 'description', 'Descrição', element.description || '', false, 'textarea');
        addFormField(formFields, 'long_description', 'Descrição Longa', element.long_description || '', false, 'textarea');
        addFormField(formFields, 'sdk_package', 'SDK Package', element.sdk_package || '', false);
        addFormField(formFields, 'cli_group', 'CLI Group', element.cli_group || '', false);
        addFormField(formFields, 'service_interface', 'Service Interface', element.service_interface || '', false);
        addFormField(formFields, 'sdk_file', 'SDK File', element.sdk_file || '', false);
        addFormField(formFields, 'custom_file', 'Custom File', element.custom_file || '', false);
        // Sempre criar campo is_group, mesmo se não existir no elemento
        const isGroupValue = element.hasOwnProperty('is_group') ? element.is_group : false;
        addFormField(formFields, 'is_group', 'Is Group', isGroupValue, false, 'checkbox');

        // Alias como campo de texto (separado por vírgulas)
        const aliasValue = element.alias ? element.alias.join(', ') : '';
        addFormField(formFields, 'alias', 'Alias (separado por vírgulas)', aliasValue, false);
    } else if (type === 'method') {
        document.getElementById('editMenuId').value = menuId || '';
        document.getElementById('editSubmenuId').value = submenuId || '';
        document.getElementById('editMethodIndex').value = methodIndex || '';
        title.textContent = `Editar Method: ${element.name}`;

        // Campos para Method
        addFormField(formFields, 'name', 'Nome', element.name || '', true);
        addFormField(formFields, 'description', 'Descrição', element.description || '', false, 'textarea');
        addFormField(formFields, 'long_description', 'Descrição Longa', element.long_description || '', false, 'textarea');
        addFormField(formFields, 'comments', 'Comentários', element.comments || '', false, 'textarea');
        addFormField(formFields, 'service_import', 'Service Import', element.service_import || '', false);
        addFormField(formFields, 'sdk_file', 'SDK File', element.sdk_file || '', false);
        addFormField(formFields, 'custom_file', 'Custom File', element.custom_file || '', false);
        addFormField(formFields, 'is_service', 'Is Service', element.is_service || false, false, 'checkbox');

        // Parameters e Returns como JSON (para edição avançada)
        addFormField(formFields, 'parameters', 'Parameters (JSON)', JSON.stringify(element.parameters || [], null, 2), false, 'textarea');
        addFormField(formFields, 'returns', 'Returns (JSON)', JSON.stringify(element.returns || [], null, 2), false, 'textarea');
        if (element.confirmation) {
            addFormField(formFields, 'confirmation', 'Confirmation (JSON)', JSON.stringify(element.confirmation, null, 2), false, 'textarea');
        }
    }

    modal.style.display = 'block';
}

function addFormField(container, name, label, value, required = false, type = 'text') {
    const formGroup = document.createElement('div');
    formGroup.className = 'form-group';

    const labelEl = document.createElement('label');
    labelEl.setAttribute('for', `edit_${name}`);

    let input;
    if (type === 'textarea') {
        input = document.createElement('textarea');
        input.rows = name.includes('JSON') ? 6 : 3;
        labelEl.textContent = label + (required ? ' *' : '');
    } else if (type === 'checkbox') {
        input = document.createElement('input');
        input.type = 'checkbox';
        input.id = `edit_${name}`;  // IMPORTANTE: definir o ID antes de retornar
        input.name = name;
        input.checked = value;
        labelEl.style.display = 'flex';
        labelEl.style.alignItems = 'center';
        labelEl.appendChild(input);
        labelEl.appendChild(document.createTextNode(' ' + label + (required ? ' *' : '')));
        formGroup.appendChild(labelEl);
        container.appendChild(formGroup);
        return;
    } else {
        input = document.createElement('input');
        input.type = type;
        labelEl.textContent = label + (required ? ' *' : '');
    }

    input.id = `edit_${name}`;
    input.name = name;
    if (type !== 'checkbox') {
        input.value = value || '';
    }
    if (required) {
        input.required = true;
    }

    formGroup.appendChild(labelEl);
    formGroup.appendChild(input);
    container.appendChild(formGroup);
}

function closeEditModal() {
    const modal = document.getElementById('editModal');
    if (modal) {
        modal.style.display = 'none';
    }
}

// Função auxiliar para encontrar menu/submenu por ID
function findMenuByIdInConfig(menus, id) {
    for (const menu of menus) {
        if (menu.id === id) {
            return menu;
        }
        if (menu.menus) {
            const found = findMenuByIdInConfig(menu.menus, id);
            if (found) return found;
        }
    }
    return null;
}

// Função auxiliar para encontrar método por menuId, submenuId e index
function findMethodInConfig(menuId, submenuId, methodIndex) {
    const menu = findMenuByIdInConfig(configData.menus, menuId);
    if (!menu) return null;

    // Se submenuId é igual a menuId ou null, o method está diretamente no menu
    if (!submenuId || submenuId === menuId) {
        if (!menu.methods) return null;
        if (methodIndex < 0 || methodIndex >= menu.methods.length) return null;
        return { menu, submenu: menu, method: menu.methods[methodIndex] };
    }

    // Caso contrário, procurar no submenu
    if (!menu.menus) return null;
    const submenu = menu.menus.find(sm => sm.id === submenuId);
    if (!submenu || !submenu.methods) return null;

    if (methodIndex < 0 || methodIndex >= submenu.methods.length) return null;

    return { menu, submenu, method: submenu.methods[methodIndex] };
}

async function submitEdit() {
    const elementTypeEl = document.getElementById('editElementType');
    if (!elementTypeEl) {
        showStatus('Erro: elemento não encontrado', 'error');
        return;
    }

    const elementType = elementTypeEl.value;
    if (!elementType) {
        showStatus('Tipo de elemento não definido', 'error');
        return;
    }

    if (!configData) {
        showStatus('Configuração não carregada', 'error');
        return;
    }

    let data = {};

    if (elementType === 'menu' || elementType === 'submenu') {
        const elementId = document.getElementById('editElementId').value;

        // Encontrar o elemento no configData local
        const element = findMenuByIdInConfig(configData.menus, elementId);
        if (!element) {
            showStatus('Elemento não encontrado', 'error');
            return;
        }

        // Coletar dados do formulário
        const cliNameEl = document.getElementById('edit_cli_name');
        const enabledEl = document.getElementById('edit_enabled');
        const isGroupEl = document.getElementById('edit_is_group');

        data.cli_name = cliNameEl ? cliNameEl.value : '';
        if (enabledEl) {
            data.enabled = enabledEl.checked;
        } else {
            data.enabled = true;
        }
        data.description = document.getElementById('edit_description')?.value || '';
        data.long_description = document.getElementById('edit_long_description')?.value || '';
        data.sdk_package = document.getElementById('edit_sdk_package')?.value || '';
        data.cli_group = document.getElementById('edit_cli_group')?.value || '';
        data.service_interface = document.getElementById('edit_service_interface')?.value || '';
        data.sdk_file = document.getElementById('edit_sdk_file')?.value || '';
        data.custom_file = document.getElementById('edit_custom_file')?.value || '';

        if (isGroupEl) {
            data.is_group = Boolean(isGroupEl.checked);
        } else {
            data.is_group = false;
        }

        // Processar alias
        const aliasEl = document.getElementById('edit_alias');
        const aliasStr = aliasEl ? aliasEl.value : '';
        data.alias = aliasStr ? aliasStr.split(',').map(a => a.trim()).filter(a => a) : [];

        if (!data.cli_name || data.cli_name.trim() === '') {
            showStatus('Nome do comando é obrigatório', 'error');
            return;
        }

        // Salvar estado antes de modificar
        saveState();

        // Atualizar localmente no configData
        element.cli_name = data.cli_name;
        element.enabled = data.enabled;
        element.description = data.description;
        element.long_description = data.long_description;
        element.sdk_package = data.sdk_package;
        element.cli_group = data.cli_group;
        element.service_interface = data.service_interface;
        element.sdk_file = data.sdk_file;
        element.custom_file = data.custom_file;
        element.is_group = data.is_group;
        element.alias = data.alias;

        showStatus(`${elementType === 'menu' ? 'Menu' : 'SubMenu'} atualizado localmente! (Clique em Salvar para persistir)`, 'success');
        closeEditModal();

        // Atualizar interface visualmente
        renderConfig(configData);
    } else if (elementType === 'method') {
        const menuId = document.getElementById('editMenuId').value;
        const submenuId = document.getElementById('editSubmenuId').value;
        const methodIndex = parseInt(document.getElementById('editMethodIndex').value);

        // Encontrar o método no configData local
        const found = findMethodInConfig(menuId, submenuId, methodIndex);
        if (!found) {
            showStatus('Método não encontrado', 'error');
            return;
        }

        const { method } = found;

        data.name = document.getElementById('edit_name').value;
        data.description = document.getElementById('edit_description').value || '';
        data.long_description = document.getElementById('edit_long_description').value || '';
        data.comments = document.getElementById('edit_comments').value || '';
        data.service_import = document.getElementById('edit_service_import')?.value || '';
        data.sdk_file = document.getElementById('edit_sdk_file')?.value || '';
        data.custom_file = document.getElementById('edit_custom_file')?.value || '';

        const isServiceEl = document.getElementById('edit_is_service');
        if (isServiceEl) {
            data.is_service = isServiceEl.checked;
        } else {
            data.is_service = false;
        }

        // Processar JSON fields
        try {
            const paramsStr = document.getElementById('edit_parameters').value || '[]';
            data.parameters = JSON.parse(paramsStr);
        } catch (e) {
            showStatus('Erro ao processar Parameters JSON: ' + e.message, 'error');
            return;
        }

        try {
            const returnsStr = document.getElementById('edit_returns').value || '[]';
            data.returns = JSON.parse(returnsStr);
        } catch (e) {
            showStatus('Erro ao processar Returns JSON: ' + e.message, 'error');
            return;
        }

        const confirmationEl = document.getElementById('edit_confirmation');
        if (confirmationEl && confirmationEl.value) {
            try {
                data.confirmation = JSON.parse(confirmationEl.value);
            } catch (e) {
                // Ignorar se não for JSON válido
            }
        }

        if (!data.name || data.name.trim() === '') {
            showStatus('Nome é obrigatório', 'error');
            return;
        }

        // Salvar estado antes de modificar
        saveState();

        // Atualizar localmente no configData
        method.name = data.name;
        method.description = data.description;
        method.long_description = data.long_description;
        method.comments = data.comments;
        method.service_import = data.service_import;
        method.sdk_file = data.sdk_file;
        method.custom_file = data.custom_file;
        method.is_service = data.is_service;
        method.parameters = data.parameters;
        method.returns = data.returns;
        if (data.confirmation !== undefined) {
            method.confirmation = data.confirmation;
        }

        showStatus('Método atualizado localmente! (Clique em Salvar para persistir)', 'success');
        closeEditModal();

        // Atualizar interface visualmente
        renderConfig(configData);
    }
}

// Função auxiliar para remover menu/submenu por ID recursivamente
function removeMenuByIdFromConfig(menus, id) {
    for (let i = 0; i < menus.length; i++) {
        if (menus[i].id === id) {
            menus.splice(i, 1);
            return true;
        }
        if (menus[i].menus) {
            if (removeMenuByIdFromConfig(menus[i].menus, id)) {
                return true;
            }
        }
    }
    return false;
}

async function deleteMenu(id) {
    if (!id) {
        showStatus('ID do menu não encontrado', 'error');
        return;
    }

    if (!configData || !configData.menus) {
        showStatus('Configuração não carregada', 'error');
        return;
    }

    // Salvar estado antes de modificar
    saveState();

    // Remover localmente do configData
    const removed = removeMenuByIdFromConfig(configData.menus, id);

    if (!removed) {
        showStatus('Menu ou submenu não encontrado', 'error');
        return;
    }

    showStatus('Menu removido localmente! (Clique em Salvar para persistir)', 'success');

    // Atualizar interface visualmente
    renderConfig(configData);
}

// Função auxiliar para encontrar e remover elemento por ID, retornando o elemento removido
function findAndRemoveElementById(menus, id) {
    // Procurar nos menus principais
    for (let i = 0; i < menus.length; i++) {
        if (menus[i].id === id) {
            const element = menus[i];
            menus.splice(i, 1);
            return { element, parentMenus: menus, index: i };
        }
        // Procurar recursivamente nos submenus
        if (menus[i].menus) {
            const result = findAndRemoveElementInSubmenus(menus[i].menus, id);
            if (result) return result;
        }
    }
    return null;
}

function findAndRemoveElementInSubmenus(menus, id) {
    for (let i = 0; i < menus.length; i++) {
        if (menus[i].id === id) {
            const element = menus[i];
            menus.splice(i, 1);
            return { element, parentMenus: menus, index: i };
        }
        // Procurar recursivamente nos submenus aninhados
        if (menus[i].menus) {
            const result = findAndRemoveElementInSubmenus(menus[i].menus, id);
            if (result) return result;
        }
    }
    return null;
}

async function moveElementToBackend(elementID, elementType, targetID, targetType) {
    if (!elementID) {
        showStatus('ID do elemento não encontrado', 'error');
        return;
    }

    if (!configData || !configData.menus) {
        showStatus('Configuração não carregada', 'error');
        return;
    }

    // Verificar se está tentando mover para dentro de si mesmo
    if (targetID === elementID) {
        showStatus('Não é possível mover um elemento para dentro de si mesmo', 'error');
        return;
    }

    // Verificar se está tentando mover para dentro de um de seus próprios filhos
    const elementToMove = findMenuByIdInConfig(configData.menus, elementID);
    if (elementToMove && targetID) {
        const targetMenu = findMenuByIdInConfig(configData.menus, targetID);
        if (targetMenu && isDescendantOf(elementToMove, targetMenu)) {
            showStatus('Não é possível mover um elemento para dentro de seus próprios filhos', 'error');
            return;
        }
    }

    // Salvar estado antes de modificar
    saveState();

    // Encontrar e remover o elemento do local original
    const removed = findAndRemoveElementById(configData.menus, elementID);
    if (!removed) {
        showStatus('Elemento não encontrado', 'error');
        return;
    }

    const { element } = removed;

    // Adicionar ao destino e atualizar parent_menu_id
    if (targetType === 'root' || !targetID) {
        // Mover para a raiz - limpar parent_menu_id e atualizar recursivamente todos os submenus filhos
        updateParentMenuIDRecursive(element, '');
        configData.menus.push(element);
    } else {
        // Encontrar o destino
        const targetMenu = findMenuByIdInConfig(configData.menus, targetID);
        if (!targetMenu) {
            // Se não encontrou, tentar adicionar de volta ao local original
            removed.parentMenus.splice(removed.index, 0, element);
            showStatus('Destino não encontrado', 'error');
            return;
        }
        if (!targetMenu.menus) {
            targetMenu.menus = [];
        }
        // Atualizar parent_menu_id do menu movido para o ID do menu pai e atualizar recursivamente todos os submenus filhos
        updateParentMenuIDRecursive(element, targetMenu.id);
        targetMenu.menus.push(element);
    }

    showStatus('Elemento movido localmente! (Clique em Salvar para persistir)', 'success');

    // Atualizar interface visualmente
    renderConfig(configData);
}

// Função auxiliar para verificar se um menu é descendente de outro
function isDescendantOf(menu, ancestor) {
    if (!menu.menus || menu.menus.length === 0) {
        return false;
    }
    for (const child of menu.menus) {
        if (child.id === ancestor.id) {
            return true;
        }
        if (isDescendantOf(child, ancestor)) {
            return true;
        }
    }
    return false;
}

// Função para atualizar recursivamente o parent_menu_id de um menu e todos os seus submenus
function updateParentMenuIDRecursive(menu, parentID) {
    if (!menu) {
        return;
    }
    // Atualizar o parent_menu_id do menu atual
    menu.parent_menu_id = parentID;
    // Atualizar recursivamente todos os submenus filhos
    if (menu.menus && menu.menus.length > 0) {
        for (const submenu of menu.menus) {
            updateParentMenuIDRecursive(submenu, menu.id);
        }
    }
}

// Função para promover submenu para o mesmo nível do parent
function promoteSubmenuToMenu(submenuID) {
    if (!submenuID) {
        showStatus('ID do submenu não encontrado', 'error');
        return;
    }

    if (!configData || !configData.menus) {
        showStatus('Configuração não carregada', 'error');
        return;
    }

    // Salvar estado antes de modificar
    saveState();

    // Encontrar o submenu
    const submenu = findMenuByIdInConfig(configData.menus, submenuID);
    if (!submenu) {
        showStatus('Submenu não encontrado', 'error');
        return;
    }

    // Obter o parent_menu_id do submenu
    const parentMenuID = submenu.parent_menu_id;

    // Encontrar e remover o submenu do local original
    const removed = findAndRemoveElementById(configData.menus, submenuID);
    if (!removed) {
        showStatus('Submenu não encontrado', 'error');
        return;
    }

    const { element } = removed;

    // Se o submenu tem um parent (não está na raiz)
    if (parentMenuID) {
        // Encontrar o parent do submenu
        const parentMenu = findMenuByIdInConfig(configData.menus, parentMenuID);
        if (parentMenu) {
            // Verificar se o parent tem um parent (não está na raiz)
            if (parentMenu.parent_menu_id) {
                // Encontrar o parent do parent (avô)
                const grandParentMenu = findMenuByIdInConfig(configData.menus, parentMenu.parent_menu_id);
                if (grandParentMenu) {
                    // Mover para dentro do avô (mesmo nível do parent)
                    if (!grandParentMenu.menus) {
                        grandParentMenu.menus = [];
                    }
                    updateParentMenuIDRecursive(element, grandParentMenu.id);
                    grandParentMenu.menus.push(element);
                    showStatus('Submenu promovido para o mesmo nível do parent! (Clique em Salvar para persistir)', 'success');
                } else {
                    // Se não encontrou o avô, mover para a raiz
                    updateParentMenuIDRecursive(element, '');
                    configData.menus.push(element);
                    showStatus('Submenu promovido para menu! (Clique em Salvar para persistir)', 'success');
                }
            } else {
                // O parent está na raiz, então mover para a raiz também
                updateParentMenuIDRecursive(element, '');
                configData.menus.push(element);
                showStatus('Submenu promovido para menu! (Clique em Salvar para persistir)', 'success');
            }
        } else {
            // Se não encontrou o parent, mover para a raiz
            updateParentMenuIDRecursive(element, '');
            configData.menus.push(element);
            showStatus('Submenu promovido para menu! (Clique em Salvar para persistir)', 'success');
        }
    } else {
        // Se o submenu já está na raiz (não tem parent), não há o que promover
        // Adicionar de volta ao local original
        removed.parentMenus.splice(removed.index, 0, element);
        showStatus('Submenu já está no nível raiz', 'info');
        return;
    }

    // Atualizar interface visualmente
    renderConfig(configData);
}

// Função auxiliar para gerar UUID simples
function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

async function submitCreateMenu() {
    const form = document.getElementById('createMenuForm');
    const formData = new FormData(form);

    const menuData = {
        sdk_name: formData.get('sdk_name') || document.getElementById('menuSDKName').value,
        cli_name: formData.get('cli_name') || document.getElementById('menuCliName').value,
        sdk_package: formData.get('sdk_package') || document.getElementById('menuSDKPackage').value,
        description: formData.get('description') || document.getElementById('menuDescription').value,
        enabled: document.getElementById('menuEnabled').checked
    };

    if (!menuData.sdk_name || menuData.sdk_name.trim() === '') {
        showStatus('Nome do SDK é obrigatório', 'error');
        return;
    }

    if (!menuData.cli_name || menuData.cli_name.trim() === '') {
        showStatus('Nome do comando é obrigatório', 'error');
        return;
    }

    if (!configData) {
        showStatus('Configuração não carregada', 'error');
        return;
    }

    // Salvar estado antes de modificar
    saveState();

    // Criar novo menu localmente
    const newMenu = {
        id: generateUUID(),
        sdk_name: menuData.sdk_name,
        cli_name: menuData.cli_name,
        sdk_package: menuData.sdk_package || '',
        description: menuData.description || '',
        enabled: menuData.enabled,
        menus: [],
        methods: []
    };

    // Adicionar ao configData
    if (!configData.menus) {
        configData.menus = [];
    }
    configData.menus.push(newMenu);

    showStatus(`Menu "${menuData.cli_name}" criado localmente! (Clique em Salvar para persistir)`, 'success');
    closeCreateMenuModal();

    // Atualizar interface visualmente
    renderConfig(configData);
}
