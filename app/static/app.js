class MenuApp {
    constructor() {
        this.menuContainer = document.getElementById('menu-container');
        this.filters = document.querySelectorAll('.filter-btn');
        this.currentCategory = 'all';
        this.init();
    }

    init() {
        this.loadMenu();
        this.setupFilters();
    }

    async loadMenu() {
        try {
            this.showSkeleton(6); // Показываем 6 скелетонов
            
            const response = await fetch('/api/menu');
            if (!response.ok) throw new Error('Ошибка загрузки меню');
            
            const products = await response.json();
            this.renderMenu(products);
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

    showSkeleton(count) {
        this.menuContainer.innerHTML = Array(count).fill(`
            <div class="product-card skeleton"></div>
        `).join('');
    }

    renderMenu(products) {
        this.menuContainer.innerHTML = products.map(product => `
            <div class="product-card" data-category="${product.category}">
            <div class="product-image" style="background-image: url('/static/images/${product.id}.png')"></div>
            <h3>${product.name}</h3>
            <div class="price">${product.price} ₽</div>
            <div class="description">Колличество: ${product.count} шт.</div>
            <button class="add-btn" data-id="${product.id}">
                Добавить в корзину
            </button>
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
            3: 'десерты'
        };
        return categories[categoryId] || '';
    }

    addToCart(productId) {
        // Здесь будет логика добавления в корзину
        console.log('Adding to cart:', productId);
        
        // Временная анимация
        const btn = document.querySelector(`[data-id="${productId}"]`);
        btn.textContent = '✅ Добавлено!';
        setTimeout(() => {
            btn.textContent = 'Добавить в корзину';
        }, 2000);
    }
}

// Запуск приложения
document.addEventListener('DOMContentLoaded', () => {
    new MenuApp();
});