// ui/static/js/main.js

document.addEventListener('DOMContentLoaded', () => {
    console.log("UniSystem Dashboard Iniciado");

    // 1. Se os elementos de contagem existem, carrega as estatísticas
    if (document.getElementById('count-students')) {
        loadDashboardStats();
    }

    // 2. Inicia a animação de entrada dos cards
    animateCards();
});

async function loadDashboardStats() {
    try {
        const response = await fetch('/api/dashboard/stats');

        if (!response.ok) {
            console.error("Erro na resposta do servidor:", response.status);
            return;
        }

        const stats = await response.json();

        // Anima os números de 0 até o valor final
        animateValue("count-students", 0, stats.students, 1000);
        animateValue("count-teachers", 0, stats.teachers, 1000);
        animateValue("count-courses", 0, stats.courses, 1000);
        animateValue("count-disciplines", 0, stats.disciplines, 1000);
        animateValue("count-semesters", 0, stats.semesters, 1000);
        animateValue("count-depts", 0, stats.departments, 1000);

    } catch (error) {
        console.error("Erro ao carregar stats:", error);
    }
}

// Função matemática para o efeito de contagem
function animateValue(id, start, end, duration) {
    const obj = document.getElementById(id);
    if (!obj) return;

    // Se o valor for 0 ou nulo, apenas exibe 0 sem animar
    if (!end || end === 0) {
        obj.innerHTML = "0";
        return;
    }

    let startTimestamp = null;
    const step = (timestamp) => {
        if (!startTimestamp) startTimestamp = timestamp;
        const progress = Math.min((timestamp - startTimestamp) / duration, 1);
        obj.innerHTML = Math.floor(progress * (end - start) + start);
        if (progress < 1) {
            window.requestAnimationFrame(step);
        } else {
            obj.innerHTML = end; // Garante que termina no número exato
        }
    };
    window.requestAnimationFrame(step);
}

// Efeito visual (CSS) para os cards aparecerem suavemente
function animateCards() {
    const cards = document.querySelectorAll('.dashboard-card');
    cards.forEach((card, index) => {
        // Estado inicial (invisível e deslocado)
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';

        // Animação com delay progressivo
        setTimeout(() => {
            card.style.transition = 'all 0.5s ease';
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, 100 * (index + 1));
    });
}