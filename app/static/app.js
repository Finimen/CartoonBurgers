class MenuApp {
    constructor() {
        this.cart = [];
        this.menuContainer = document.getElementById('menu-container');
        this.filters = document.querySelectorAll('.filter-btn');
        this.currentCategory = 'all';
        this.init();
        this.updateAuthUI();
        this.map = null;
        this.mapInitialized = false;
        this.setupOrderForm();
    }

    init() {
        this.loadMenu();
        this.setupFilters();
        this.setupAuthModals(); 
        this.setupNavigation();
    }

    setupOrderForm() {
        const checkoutForm = document.getElementById('checkout-form');
        const paymentMethod = document.getElementById('payment-method');
        const cardFields = document.getElementById('card-fields');
        
        if (!checkoutForm || !paymentMethod || !cardFields) {
            console.error('Не найдены элементы формы заказа');
            return;
        }

        paymentMethod.addEventListener('change', () => {
            if (paymentMethod.value === 'card') {
                cardFields.classList.remove('hidden');
            } else {
                cardFields.classList.add('hidden');
            }
        });

        this.setupInputMasks();

        checkoutForm.addEventListener('submit', (e) => {
            e.preventDefault();
            
            if (!this.validatePhone()) {
                return;
            }
            
            if (paymentMethod.value === 'card' && !this.validateCard()) {
                return;
            }
            
            this.processOrder();
        });
    }

    validatePhone() {
        const phoneInput = document.getElementById('delivery-phone');
        if (!phoneInput) return false;

        const phoneValue = phoneInput.value.trim();
        const cleanPhone = phoneValue.replace(/[^\d+]/g, '');
        
        const phoneRegex = /^(\+7|8)(\d{10})$/;
        
        if (!phoneRegex.test(cleanPhone)) {
            this.showNotification('Пожалуйста, введите корректный номер телефона России');
            phoneInput.focus();
            return false;
        }
        
        return true;
    }

    validateCard() {
        const cardNumber = document.getElementById('card-number');
        const cardExpiry = document.getElementById('card-expiry');
        const cardCvv = document.getElementById('card-cvv');
        
        if (!cardNumber || !cardExpiry || !cardCvv) return false;

        const cardRegex = /^\d{16}$/;
        if (!cardRegex.test(cardNumber.value.replace(/\s/g, ''))) {
            this.showNotification('Введите корректный номер карты (16 цифр)');
            cardNumber.focus();
            return false;
        }
        
        const expiryRegex = /^\d{2}\/\d{2}$/;
        if (!expiryRegex.test(cardExpiry.value)) {
            this.showNotification('Введите корректную дату (ММ/ГГ)');
            cardExpiry.focus();
            return false;
        }
        
        const cvvRegex = /^\d{3}$/;
        if (!cvvRegex.test(cardCvv.value)) {
            this.showNotification('Введите корректный CVV (3 цифры)');
            cardCvv.focus();
            return false;
        }
        
        return true;
    }

    setupInputMasks() {
        const cardNumber = document.getElementById('card-number');
        if (cardNumber) {
            cardNumber.addEventListener('input', (e) => {
                let value = e.target.value.replace(/\D/g, '');
                if (value.length > 16) value = value.slice(0, 16);
                
                value = value.replace(/(\d{4})/g, '$1 ').trim();
                e.target.value = value;
            });
        }

        const cardExpiry = document.getElementById('card-expiry');
        if (cardExpiry) {
            cardExpiry.addEventListener('input', (e) => {
                let value = e.target.value.replace(/\D/g, '');
                if (value.length > 4) value = value.slice(0, 4);
                
                if (value.length > 2) {
                    value = value.slice(0, 2) + '/' + value.slice(2);
                }
                e.target.value = value;
            });
        }

        const cardCvv = document.getElementById('card-cvv');
        if (cardCvv) {
            cardCvv.addEventListener('input', (e) => {
                e.target.value = e.target.value.replace(/\D/g, '').slice(0, 3);
            });
        }

        const phoneInput = document.getElementById('delivery-phone');
        if (phoneInput) {
            phoneInput.addEventListener('input', (e) => {
                let value = e.target.value.replace(/\D/g, '');
                
                if (value.length > 11) value = value.slice(0, 11);
                
                if (value.startsWith('7')) {
                    value = '+7' + value.slice(1);
                } else if (value.startsWith('8')) {
                    value = '+7' + value.slice(1);
                } else if (!value.startsWith('+7') && value.length > 0) {
                    value = '+7' + value;
                }
                
                if (value.length > 2) {
                    value = value.replace(/^(\+7)(\d{0,3})/, '$1 ($2');
                }
                if (value.length > 7) {
                    value = value.replace(/^(\+7\s?\(\d{3})\)?(\d{0,3})/, '$1) $2');
                }
                if (value.length > 11) {
                    value = value.replace(/^(\+7\s?\(\d{3}\)\s?\d{3})-?(\d{0,2})/, '$1-$2');
                }
                if (value.length > 14) {
                    value = value.replace(/^(\+7\s?\(\d{3}\)\s?\d{3}-\d{2})-?(\d{0,2})/, '$1-$2');
                }
                
                e.target.value = value;
            });
        }
    }

    initMap() {
    if (this.mapInitialized) return;

    try {
        const mapContainer = document.getElementById('map');
        if (!mapContainer) {
            console.error('Map container not found');
            return;
        }

        this.map = L.map('map').setView([55.7558, 37.6173], 10); // Москва

        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(this.map);

        this.map.on('click', (e) => {
            const latlng = e.latlng;
            
            if (this.marker) {
                this.map.removeLayer(this.marker);
            }
            
            this.marker = L.marker(latlng, {
                icon: L.divIcon({
                    className: 'custom-marker',
                    html: '<div style="background:#8b5cf6;width:20px;height:20px;border-radius:50%;border:3px solid white;"></div>',
                    iconSize: [20, 20]
                })
            }).addTo(this.map);
            
            this.getAddressFromCoordinates(latlng.lat, latlng.lng);
        });

        this.mapInitialized = true;
        
        setTimeout(() => {
            if (this.map) this.map.invalidateSize();
        }, 100);

    } catch (error) {
        console.error('Error initializing map:', error);
    }
}

    async getAddressFromCoordinates(lat, lng) {
        try {
            const response = await fetch(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&zoom=18&addressdetails=1`);
            const data = await response.json();
            
            if (data && data.display_name) {
                document.getElementById('delivery-address').value = data.display_name;
            }
        } catch (error) {
            console.error('Ошибка получения адреса:', error);
            document.getElementById('delivery-address').value = `Широта: ${lat.toFixed(6)}, Долгота: ${lng.toFixed(6)}`;
        }
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
        document.querySelectorAll('main, section').forEach(section => {
            section.classList.add('hidden');
        });

        if (sectionName === 'menu') {
            document.querySelector('main').classList.remove('hidden');
        } else if (sectionName === 'cart') {
            document.getElementById('cart-section').classList.remove('hidden');
            setTimeout(() => {
                if (!this.mapInitialized) {
                    this.initMap();
                } else {
                    setTimeout(() => this.map.invalidateSize(), 50);
                }
            }, 100);
        } else if (sectionName === 'profile') {
            if (!this.checkAuth()) {
                this.showLoginModal();
                return;
            }
            document.getElementById('profile-section').classList.remove('hidden');
            this.loadProfileData();
        } else {
            document.getElementById(sectionName + '-section').classList.remove('hidden');
        }

        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
        });
        document.querySelector(`[data-section="${sectionName}"]`).classList.add('active');
    }

    processOrder() {
        if (this.cart.length === 0) {
            this.showNotification('Корзина пуста');
            return;
        }
        
        const address = document.getElementById('delivery-address').value;
        const phone = document.getElementById('delivery-phone').value;
        
        this.showNotification('Заказ успешно оформлен!');
        this.cart = [];
        this.updateCartUI();
    }


    async loadMenu() {
    try {
        this.showSkeleton(6);
        
        const response = await fetch('/api/menu');
        if (!response.ok) throw new Error('Ошибка загрузки меню');
        
        this.products = await response.json(); 
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
    const btn = document.querySelector(`[data-id="${productId}"]`);
    btn.textContent = '✅ Добавлено!';
    setTimeout(() => {
        btn.textContent = 'Добавить в корзину';
    }, 2000);

    const product = this.products.find(p => p.id == productId);
    if (!product) {
        console.error('Product not found:', productId);
        return;
    }
    
    const existingItem = this.cart.find(item => item.id == productId);
    
    if (existingItem) {
        existingItem.quantity = (existingItem.quantity || 1) + 1;
    } else {
        this.cart.push({
            ...product,
            quantity: 1
        });
    }
    
    this.updateCartUI();
    this.showNotification(`Добавлено: ${product.name}`);
}

showSection(sectionName) {
    document.querySelectorAll('main, section').forEach(section => {
        section.classList.add('hidden');
    });

    if (sectionName === 'menu') {
        document.querySelector('main').classList.remove('hidden');
    } else if (sectionName === 'profile') {
        if (!this.checkAuth()) {
            this.showLoginModal();
            return;
        }
        document.getElementById('profile-section').classList.remove('hidden');
        this.loadProfileData();
    } else {
        document.getElementById(sectionName + '-section').classList.remove('hidden');
    }

    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
    });
    document.querySelector(`[data-section="${sectionName}"]`).classList.add('active');
}

checkAuth() {
    const token = localStorage.getItem('token');
    return !!token;
}

async loadProfileData() {
    const token = localStorage.getItem('token');
    const userInfo = document.getElementById('user-info');
    const authMessage = document.getElementById('auth-required-message');
    
    if (!token) {
        userInfo.classList.add('hidden');
        authMessage.classList.remove('hidden');
        return;
    }

    try {
        const response = await fetch('/api/profile', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (response.ok) {
            const profileData = await response.json();
            this.renderProfile(profileData);
            userInfo.classList.remove('hidden');
            authMessage.classList.add('hidden');
        } else if (response.status === 401) {
            localStorage.removeItem('token');
            this.updateAuthUI();
            userInfo.classList.add('hidden');
            authMessage.classList.remove('hidden');
        }
    } catch (error) {
        console.error('Ошибка загрузки профиля:', error);
        userInfo.classList.add('hidden');
        authMessage.classList.remove('hidden');
    }
}

    renderProfile(data) {
        const usernameElement = document.getElementById('user-username');
        const emailElement = document.getElementById('user-email');
        const bonusElement = document.getElementById('user-bonus');
        
        // Анимированное появление данных
        setTimeout(() => {
            usernameElement.textContent = data.username || 'Не указан';
            usernameElement.style.animation = 'typewriter 2s steps(20) forwards';
        }, 300);
        
        setTimeout(() => {
            emailElement.textContent = data.email || 'Не указан';
            emailElement.style.opacity = '0';
            emailElement.style.transform = 'translateX(-20px)';
            setTimeout(() => {
                emailElement.style.transition = 'all 0.5s ease';
                emailElement.style.opacity = '1';
                emailElement.style.transform = 'translateX(0)';
            }, 100);
        }, 800);
        
        setTimeout(() => {
            bonusElement.textContent = data.bonuses || 0;
            bonusElement.style.animation = 'countUp 1s ease-out forwards';
        }, 1200);
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
            this.showNotification('Успешный вход!');
        } else {
            alert('Ошибка входа: неверные учетные данные');
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Ошибка сети');
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
    const authButtons = document.querySelector('.auth-buttons');
    
    if (token) {
        authButtons.innerHTML = `
            <button class="auth-btn" onclick="app.logout()">Выйти</button>
        `;
    } else {
        authButtons.innerHTML = `
            <button class="auth-btn" onclick="app.showLoginModal()">Войти</button>
            <button class="auth-btn primary" onclick="app.showRegisterModal()">Регистрация</button>
        `;
    }
}

logout() {
    localStorage.removeItem('token');
    this.updateAuthUI();
    this.showNotification('Вы вышли из системы');
    this.showSection('menu');
}

    updateCartUI() {
    const cartItemsContainer = document.getElementById('cart-items');
    const totalPriceElement = document.getElementById('total-price');
    
    if (!cartItemsContainer || !totalPriceElement) return;
    
    cartItemsContainer.innerHTML = '';
    
    if (this.cart.length === 0) {
        cartItemsContainer.innerHTML = '<p class="empty-cart">Корзина пуста</p>';
        totalPriceElement.textContent = '0';
        return;
    }
    
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
    
    const total = this.cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);
    totalPriceElement.textContent = total;
}

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

removeItemCompletely(productId) {
    const itemIndex = this.cart.findIndex(item => item.id == productId);
    
    if (itemIndex !== -1) {
        this.cart.splice(itemIndex, 1);
        this.updateCartUI();
        this.showNotification('Товар удален из корзины');
    }
}

showNotification(message) {
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
    
    setTimeout(() => {
        notification.remove();
    }, 3000);
}
}

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

class ParticlesBackground {
    constructor() {
        this.canvas = document.getElementById('particles-canvas');
        this.ctx = this.canvas.getContext('2d');
        this.particles = [];
        this.mouse = { x: null, y: null, radius: 100 };
        
        this.init();
        this.animate();
    }

    init() {
        this.resize();
        window.addEventListener('resize', () => this.resize());

        this.createParticles();

        window.addEventListener('mousemove', (e) => {
            this.mouse.x = e.x;
            this.mouse.y = e.y;
        });

        window.addEventListener('mouseout', () => {
            this.mouse.x = null;
            this.mouse.y = null;
        });
    }

    resize() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
    }

    createParticles() {
        const particleCount = Math.min(100, Math.floor(window.innerWidth / 10));
        
        for (let i = 0; i < particleCount; i++) {
            this.particles.push({
                x: Math.random() * this.canvas.width,
                y: Math.random() * this.canvas.height,
                size: Math.random() * 2 + 1,
                speedX: Math.random() * 0.5 - 0.25,
                speedY: Math.random() * 0.5 - 0.25,
                color: `rgba(139, 92, 246, ${Math.random() * 0.3 + 0.1})`
            });
        }
    }

    animate() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        
        for (let particle of this.particles) {
            this.updateParticle(particle);
            this.drawParticle(particle);
        }

        this.connectParticles();

        requestAnimationFrame(() => this.animate());
    }

    updateParticle(particle) {
        particle.x += particle.speedX;
        particle.y += particle.speedY;

        if (particle.x > this.canvas.width || particle.x < 0) {
            particle.speedX = -particle.speedX;
        }
        if (particle.y > this.canvas.height || particle.y < 0) {
            particle.speedY = -particle.speedY;
        }

        if (this.mouse.x && this.mouse.y) {
            const dx = particle.x - this.mouse.x;
            const dy = particle.y - this.mouse.y;
            const distance = Math.sqrt(dx * dx + dy * dy);

            if (distance < this.mouse.radius) {
                const angle = Math.atan2(dy, dx);
                const force = (this.mouse.radius - distance) / this.mouse.radius;
                particle.x += Math.cos(angle) * force * 2;
                particle.y += Math.sin(angle) * force * 2;
            }
        }
    }

    drawParticle(particle) {
        this.ctx.beginPath();
        this.ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
        this.ctx.fillStyle = particle.color;
        this.ctx.fill();
    }

    connectParticles() {
        const maxDistance = 150;
        
        for (let i = 0; i < this.particles.length; i++) {
            for (let j = i + 1; j < this.particles.length; j++) {
                const dx = this.particles[i].x - this.particles[j].x;
                const dy = this.particles[i].y - this.particles[j].y;
                const distance = Math.sqrt(dx * dx + dy * dy);

                if (distance < maxDistance) {
                    const opacity = 1 - distance / maxDistance;
                    this.ctx.beginPath();
                    this.ctx.strokeStyle = `rgba(139, 92, 246, ${opacity * 0.2})`;
                    this.ctx.lineWidth = 1;
                    this.ctx.moveTo(this.particles[i].x, this.particles[i].y);
                    this.ctx.lineTo(this.particles[j].x, this.particles[j].y);
                    this.ctx.stroke();
                }
            }
        }
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.app = new MenuApp();
    new ParticlesBackground(); 
});