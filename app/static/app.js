class MenuApp {
    constructor() {
        this.cart = [];
        this.menuContainer = document.getElementById('menu-container');
        this.filters = document.querySelectorAll('.filter-btn');
        this.currentCategory = 'all';
        this.init();
    }

    init() {
        this.loadMenu();
        this.setupFilters();
        this.setupAuthModals(); 
        this.setupNavigation();
    }

    setupNavigation() {
        document.querySelectorAll('.nav-link').forEach(link => {
            link.addEventListener('click', (e) => {
                e.preventDefault();
                const section = e.target.dataset.section;
                this.showSection(section);
            });
        });
    }

    showSection(sectionName) {
        // Скрыть все разделы
        document.querySelectorAll('main, section').forEach(section => {
            section.classList.add('hidden');
        });

        // Показать выбранный раздел
        if (sectionName === 'menu') {
            document.querySelector('main').classList.remove('hidden');
        } else {
            document.getElementById(sectionName + '-section').classList.remove('hidden');
        }

        // Обновить активную ссылку в навигации
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        document.querySelector(`[data-section="${sectionName}"]`).classList.add('active');
    }


    async loadMenu() {
    try {
        this.showSkeleton(6); // Показываем 6 скелетонов
        
        const response = await fetch('/api/menu');
        if (!response.ok) throw new Error('Ошибка загрузки меню');
        
        this.products = await response.json(); // Сохраняем продукты в свойство
        this.renderMenu(this.products);
    } catch (error) {
        console.error('Error:', error);
        this.menuContainer.innerHTML = `
            <div class="error">
                <p>😕 Не удалось загрузить меню</p>
                <button onclick="location.reload()">Попробовать снова</button>
            </div>
        `;
    }
}

    setupAuthModals() {
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                this.closeModals();
            }
        });
    }

    showSkeleton(count) {
        this.menuContainer.innerHTML = Array(count).fill(`
            <div class="product-card skeleton"></div>
        `).join('');
    }

    renderMenu(products) {
        this.menuContainer.innerHTML = products.map(product => `
            <div class="product-card" data-category="${product.category}">
            <div class="product-image" style="background-image: url('/static/images/${product.id}.png')">
                ${product.type === 1 ? '<div class="badge-new">НОВИНКА!</div>' : ''}
            </div>
            <h3>${product.name}</h3>
            <div class="price">${product.price} ₽</div>
            <div class="description">Количество: ${product.count} шт.</div>
            <button class="add-btn" data-id="${product.id}">
                Добавить в корзину
            </button>
        </div>
        </div>
        `).join('');

        this.setupAddButtons();
    }

    setupAddButtons() {
        document.querySelectorAll('.add-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const productId = e.target.dataset.id;
                this.addToCart(productId);
            });
        });
    }

    setupFilters() {
        this.filters.forEach(btn => {
            btn.addEventListener('click', (e) => {
                this.filters.forEach(b => b.classList.remove('active'));
                e.target.classList.add('active');
                
                this.currentCategory = e.target.textContent.toLowerCase();
                this.applyFilter();
            });
        });
    }

    applyFilter() {
        const cards = document.querySelectorAll('.product-card');
        
        cards.forEach(card => {
            if (this.currentCategory === 'все' || this.currentCategory === 'all') {
                card.style.display = 'block';
            } else {
                const category = card.dataset.category;
                const shouldShow = this.getCategoryName(category) === this.currentCategory;
                card.style.display = shouldShow ? 'block' : 'none';
            }
        });
    }

    getCategoryName(categoryId) {
        const categories = {
            0: 'бургеры',
            1: 'закуски', 
            2: 'напитки',
            3: 'десерты',
        };
        return categories[categoryId] || '';
    }

    addToCart(productId) {
    // Анимация добавления
    const btn = document.querySelector(`[data-id="${productId}"]`);
    btn.textContent = '✅ Добавлено!';
    setTimeout(() => {
        btn.textContent = 'Добавить в корзину';
    }, 2000);

    // Находим продукт в сохраненном массиве
    const product = this.products.find(p => p.id == productId);
    if (!product) {
        console.error('Product not found:', productId);
        return;
    }
    
    // Проверяем, есть ли уже такой продукт в корзине
    const existingItem = this.cart.find(item => item.id == productId);
    
    if (existingItem) {
        // Увеличиваем количество
        existingItem.quantity = (existingItem.quantity || 1) + 1;
    } else {
        // Добавляем новый продукт с количеством 1
        this.cart.push({
            ...product,
            quantity: 1
        });
    }
    
    this.updateCartUI();
    this.showNotification(`Добавлено: ${product.name}`);
}

    setupAuthModals() {
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                this.closeModals();
            }
        });
    }

    showLoginModal() {
        document.getElementById('loginModal').style.display = 'block';
    }

    showRegisterModal() {
        document.getElementById('registerModal').style.display = 'block';
    }

    closeModals() {
        document.querySelectorAll('.modal').forEach(modal => {
            modal.style.display = 'none';
        });
    }

    async handleLogin(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = {
            username: formData.get('username'),
            password: formData.get('password')
        };

        try {
            const response = await fetch('/api/auth/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                const { token } = await response.json();
                localStorage.setItem('token', token);
                this.closeModals();
                this.updateAuthUI();
            } else {
                alert('Ошибка входа');
            }
        } catch (error) {
            console.error('Login error:', error);
        }
    }

    async handleRegister(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const data = {
            username: formData.get('username'),
            email: formData.get('email'),
            password: formData.get('password')
        };

        try {
            const response = await fetch('/api/auth/register', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });
            
            if (response.ok) {
                alert('Регистрация успешна! Теперь войдите.');
                this.closeModals();
                this.showLoginModal();
            } else {
                alert('Ошибка регистрации');
            }
        } catch (error) {
            console.error('Register error:', error);
        }
    }

    updateAuthUI() {
        const token = localStorage.getItem('token');
        if (token) {
            // Показываем email пользователя и кнопку выхода
        } else {
            // Показываем кнопки входа/регистрации
        }
    }

    updateCartUI() {
    const cartItemsContainer = document.getElementById('cart-items');
    const totalPriceElement = document.getElementById('total-price');
    
    if (!cartItemsContainer || !totalPriceElement) return;
    
    // Очищаем контейнер
    cartItemsContainer.innerHTML = '';
    
    if (this.cart.length === 0) {
        cartItemsContainer.innerHTML = '<p class="empty-cart">Корзина пуста</p>';
        totalPriceElement.textContent = '0';
        return;
    }
    
    // Добавляем товары с новой структурой
    this.cart.forEach(item => {
        const itemElement = document.createElement('div');
        itemElement.className = 'cart-item';
        itemElement.innerHTML = `
            <div class="cart-item-content">
                <div class="cart-item-details">
                    <div class="cart-item-title">${item.name}</div>
                    <div class="cart-item-price">${item.price} ₽ × ${item.quantity} = ${item.price * item.quantity} ₽</div>
                    <div class="cart-item-quantity">
                        <button class="quantity-btn" onclick="app.changeQuantity(${item.id}, -1)">-</button>
                        <span class="quantity-number">${item.quantity}</span>
                        <button class="quantity-btn" onclick="app.changeQuantity(${item.id}, 1)">+</button>
                    </div>
                </div>
            </div>
            <button class="remove-from-cart" onclick="app.removeItemCompletely(${item.id})">×</button>
        `;
        cartItemsContainer.appendChild(itemElement);
    });
    
    // Обновляем общую сумму
    const total = this.cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    totalPriceElement.textContent = total;
}

// Добавляем метод для изменения количества
changeQuantity(productId, delta) {
    const item = this.cart.find(item => item.id == productId);
    
    if (item) {
        item.quantity += delta;
        
        if (item.quantity <= 0) {
            this.removeItemCompletely(productId);
        } else {
            this.updateCartUI();
        }
    }
}

// Метод для полного удаления товара
removeItemCompletely(productId) {
    const itemIndex = this.cart.findIndex(item => item.id == productId);
    
    if (itemIndex !== -1) {
        this.cart.splice(itemIndex, 1);
        this.updateCartUI();
        this.showNotification('Товар удален из корзины');
    }
}

showNotification(message) {
    // Создаем уведомление
    const notification = document.createElement('div');
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: #8b5cf6;
        color: white;
        padding: 15px 20px;
        border-radius: 8px;
        z-index: 1001;
        animation: slideIn 0.3s ease;
    `;
    notification.textContent = message;
    
    document.body.appendChild(notification);
    
    // Удаляем через 3 секунды
    setTimeout(() => {
        notification.remove();
    }, 3000);
}
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', () => {
    window.app = new MenuApp();
});

function showLoginModal() {
    if (window.app) {
        window.app.showLoginModal();
    }
}

function showRegisterModal() {
    if (window.app) {
        window.app.showRegisterModal();
    }
}

function handleLogin(event) {
    event.preventDefault();
    if (window.app) {
        window.app.handleLogin(event);
    }
}

function handleRegister(event) {
    event.preventDefault();
    if (window.app) {
        window.app.handleRegister(event);
    }
}

function closeModals() {
    if (window.app) {
        window.app.closeModals();
    }
}