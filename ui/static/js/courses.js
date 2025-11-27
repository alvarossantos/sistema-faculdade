const API_URL = '/api/courses';

async function loadCourses() {
    try {
        const response = await fetch(API_URL);
        const courses = await response.json();
        const tbody = document.querySelector('#coursesTable tbody');
        
        if (!tbody) return;
        tbody.innerHTML = '';

        if (!courses || courses.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="text-center py-4">Nenhum curso encontrado.</td></tr>';
            return;
        }

        courses.forEach(c => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>#${c.id}</td>
                <td class="fw-bold">${c.name}</td>
                <td>${c.total_credits_required}</td>
                <td>${c.duration_semesters} semestres</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
    }
}

async function initForm() {
    const form = document.getElementById('courseForm');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        // Mapeia os inputs para o JSON que o Go espera em models/course.go
        const data = {
            name: document.getElementById('name').value,
            total_credits_required: parseInt(document.getElementById('credits').value),
            duration_semesters: parseInt(document.getElementById('semesters').value)
        };

        try {
            const response = await fetch(API_URL, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                await Swal.fire('Sucesso!', 'Curso criado.', 'success');
                window.location.href = 'courses.html';
            } else {
                throw new Error('Erro ao criar');
            }
        } catch (error) {
            Swal.fire('Erro', 'Falha ao salvar curso.', 'error');
        }
    });
}