// ui/js/main.js

document.addEventListener('DOMContentLoaded', () => {
    console.log("UniSystem Dashboard Carregado");

    // Aqui futuramente você pode colocar:
    // 1. Lógica para carregar estatísticas (Ex: Total de Alunos no card)
    // 2. Verificação de Login (Se o usuário tem token de acesso)
    
    animateCards();
});

// Efeito simples de entrada sequencial nos cards (Opcional, mas fica bonito)
function animateCards() {
    const cards = document.querySelectorAll('.dashboard-card');
    cards.forEach((card, index) => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            card.style.transition = 'all 0.5s ease';
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, 100 * (index + 1)); // Delay progressivo
    });
}