const API_URL = '/api/departments';

async function loadDepartments() {
    try {
        const response = await fetch(API_URL);
        const departments = await response.json();
        const tbody = document.querySelector('#departmentsTable tbody');
        
        if (!tbody) return;
        tbody.innerHTML = '';

        if (!departments || departments.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="text-center py-4">Nenhum departamento encontrado.</td></tr>';
            return;
        }

        departments.forEach(d => {
            // Formata a data simples
            const date = new Date(d.created_at).toLocaleDateString('pt-BR');
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>#${d.id}</td>
                <td class="fw-bold">${d.name}</td>
                <td><span class="badge bg-secondary">${d.abbreviation}</span></td>
                <td>${date}</td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
    }
}

async function initForm() {
    const form = document.getElementById('deptForm');
    if (!form) return;

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const data = {
            name: document.getElementById('name').value,
            abbreviation: document.getElementById('abbreviation').value
        };

        try {
            const response = await fetch(API_URL, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            if (response.ok) {
                await Swal.fire('Sucesso!', 'Departamento criado.', 'success');
                window.location.href = 'departments.html';
            } else {
                throw new Error('Erro ao criar');
            }
        } catch (error) {
            Swal.fire('Erro', 'Falha ao salvar departamento.', 'error');
        }
    });
}