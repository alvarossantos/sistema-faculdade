// ui/js/teachers.js

const API_URL = '/api/teachers';

// --- LISTAGEM (teachers.html) ---

async function loadTeachers() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Erro ao buscar professores');
        
        const teachers = await response.json();
        const tbody = document.querySelector('#teachersTable tbody');
        
        if (!tbody) return; // Proteção se estivermos em outra página

        tbody.innerHTML = '';

        if (!teachers || teachers.length === 0) {
            tbody.innerHTML = '<tr><td colspan="7" class="text-center py-4">Nenhum professor encontrado.</td></tr>';
            return;
        }

        teachers.forEach(t => {
            const tr = document.createElement('tr');
            
            // Visual: Linha cinza se inativo
            if (!t.active) {
                tr.classList.add('table-secondary', 'text-muted');
                tr.style.opacity = "0.75";
            }

            // Botão Dinâmico: Lixeira (se ativo) ou Reativar (se inativo)
            let actionButton;
            if (t.active) {
                actionButton = `
                <button onclick="deleteTeacher(${t.id})" class="btn btn-sm btn-danger action-btn" title="Inativar">
                    <i class="bi bi-trash-fill"></i>
                </button>`;
            } else {
                actionButton = `
                <button onclick="activateTeacher(${t.id})" class="btn btn-sm btn-success action-btn" title="Reativar">
                    <i class="bi bi-arrow-counterclockwise"></i>
                </button>`;
            }

            // Dentro do loop teachers.forEach(t => { ... })
            // Certifique-se de que a linha que exibe o departamento use "department_name"

            tr.innerHTML = `
                <td class="fw-bold text-secondary">#${t.id}</td>
                <td class="fw-bold">${t.name}</td>
                <td>${t.email}</td>
                <td>${t.telephone}</td>
                
                <td><span class="badge bg-primary">${t.department_name || 'Sem Depto'}</span></td>
                
                <td>${t.active ? '<span class="badge bg-success">Ativo</span>' : '<span class="badge bg-secondary">Inativo</span>'}</td>
                <td class="text-end">
                    <a href="teacher_form.html?id=${t.id}" class="btn btn-sm btn-warning action-btn" title="Editar">
                        <i class="bi bi-pencil-fill"></i>
                    </a>
                    ${actionButton}
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
        if(document.querySelector('#teachersTable')) {
            Swal.fire('Erro', 'Não foi possível carregar os professores.', 'error');
        }
    }
}

// --- AÇÕES (DELETE / ACTIVATE) ---

async function deleteTeacher(id) {
    const result = await Swal.fire({
        title: 'Inativar Professor?',
        text: "O professor será marcado como inativo.",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#e74a3b',
        cancelButtonColor: '#858796',
        confirmButtonText: 'Sim, inativar'
    });

    if (result.isConfirmed) {
        try {
            const response = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
            if (response.ok) {
                Swal.fire('Inativado!', 'Professor inativado com sucesso.', 'success');
                loadTeachers();
            } else {
                throw new Error('Erro ao inativar');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao processar solicitação.', 'error');
        }
    }
}

async function activateTeacher(id) {
    const result = await Swal.fire({
        title: 'Reativar Professor?',
        text: "O cadastro voltará a ficar ativo.",
        icon: 'question',
        showCancelButton: true,
        confirmButtonColor: '#1cc88a',
        cancelButtonColor: '#858796',
        confirmButtonText: 'Sim, reativar'
    });

    if (result.isConfirmed) {
        try {
            const response = await fetch(`${API_URL}/${id}/activate`, { method: 'PATCH' });
            if (response.ok) {
                Swal.fire('Sucesso!', 'Professor reativado.', 'success');
                loadTeachers();
            } else {
                throw new Error('Erro ao ativar');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao reativar.', 'error');
        }
    }
}

// --- FORMULÁRIO (teacher_form.html) ---

async function initForm() {
    const form = document.getElementById('teacherForm');
    if (!form) return;

    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');
    const title = document.getElementById('formTitle');

    if (id) {
        if(title) title.innerText = 'Editar Professor';
        document.getElementById('teacherId').value = id;
        await loadTeacherData(id);
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveTeacher(id);
    });
}

async function loadTeacherData(id) {
    try {
        const response = await fetch(`${API_URL}/${id}`);
        if (!response.ok) throw new Error('Erro na requisição');
        
        const t = await response.json();
        
        // Preenche campos
        document.getElementById('name').value = t.name || '';
        document.getElementById('email').value = t.email || '';
        document.getElementById('cpf').value = t.cpf || '';
        document.getElementById('telephone').value = t.telephone || '';
        document.getElementById('department_id').value = t.department_id || '';

        // Formata data de contratação
        if (t.date_contract) {
            const dateObj = new Date(t.date_contract);
            if (!isNaN(dateObj)) {
                document.getElementById('date_contract').value = dateObj.toISOString().split('T')[0];
            }
        }

    } catch (error) {
        console.error("Erro loadTeacherData:", error);
        Swal.fire('Erro', 'Erro ao carregar dados do professor.', 'error');
    }
}

async function saveTeacher(id) {
    const data = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value,
        cpf: document.getElementById('cpf').value,
        telephone: document.getElementById('telephone').value,
        department_id: parseInt(document.getElementById('department_id').value),
        date_contract: new Date(document.getElementById('date_contract').value).toISOString()
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `${API_URL}/${id}` : API_URL;

    try {
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const errText = await response.text();
            throw new Error(errText || 'Erro na requisição');
        }

        await Swal.fire({
            title: 'Sucesso!',
            text: 'Professor salvo com sucesso.',
            icon: 'success',
            timer: 1500,
            showConfirmButton: false
        });

        window.location.href = 'teachers.html';
    } catch (error) {
        console.error(error);
        Swal.fire('Erro', 'Ocorreu um erro ao salvar: ' + error.message, 'error');
    }
}