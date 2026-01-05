const API_URL = '/api/disciplines';

// --- LISTAGEM (disciplines.html) ---

async function loadDisciplines() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Erro ao buscar disciplinas');

        const disciplines = await response.json();
        const tbody = document.querySelector('#disciplinesTable tbody');

        if (!tbody) return;

        tbody.innerHTML = '';

        if (!disciplines || disciplines.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center py-4">Nenhuma disciplina encontrada.</td></tr>';
            return;
        }

        disciplines.forEach(d => {
            const tr = document.createElement('tr');

            tr.innerHTML = `
                <td><span class="badge bg-secondary">${d.code}</span></td>
                <td class="fw-bold text-wrap" style="max-width: 200px;">${d.name}</td>
                
                <td class="d-none d-md-table-cell">${d.credits}</td>
                <td class="d-none d-md-table-cell">${d.workload_hours}h</td>
                <td class="d-none d-md-table-cell"><span class="badge bg-primary">${d.department_name || 'N/A'}</span></td>
                
                <td class="text-end">
                    <a href="discipline_details.html?id=${d.id}" class="btn btn-sm btn-info action-btn text-white" title="Ver">
                        <i class="bi bi-eye-fill"></i>
                    </a>
                    <a href="discipline_form.html?id=${d.id}" class="btn btn-sm btn-warning action-btn" title="Editar">
                        <i class="bi bi-pencil-fill"></i>
                    </a>
                    <button onclick="deleteDiscipline(${d.id})" class="btn btn-sm btn-danger action-btn" title="Excluir">
                        <i class="bi bi-trash-fill"></i>
                    </button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
        if (document.querySelector('#disciplinesTable')) {
            const tbody = document.querySelector('#disciplinesTable tbody');
            if (tbody) tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger">Erro ao carregar dados.</td></tr>';
        }
    }
}

// --- AÇÕES ---

async function deleteDiscipline(id) {
    const result = await Swal.fire({
        title: 'Tem certeza?',
        text: "Essa ação removerá a disciplina permanentemente.",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonColor: '#e74a3b',
        cancelButtonColor: '#858796',
        confirmButtonText: 'Sim, excluir'
    });

    if (result.isConfirmed) {
        try {
            const response = await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
            if (response.ok) {
                Swal.fire('Excluído!', 'Disciplina removida.', 'success');
                loadDisciplines();
            } else {
                throw new Error('Erro ao excluir');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao excluir disciplina.', 'error');
        }
    }
}

// --- FORMULÁRIO (discipline_form.html) ---

async function loadDepartmentsSelect() {
    try {
        const response = await fetch('/api/departments');
        if (!response.ok) throw new Error('Erro ao buscar departamentos');

        const departments = await response.json();
        const select = document.getElementById('department_id');

        if (!select) return;

        select.innerHTML = '<option value="" selected disabled>Selecione um departamento...</option>';

        departments.forEach(dep => {
            const option = document.createElement('option');
            option.value = dep.id;
            option.textContent = dep.name;
            select.appendChild(option);
        });
    } catch (error) {
        console.error("Erro ao carregar departamentos:", error);
    }
}

async function initForm() {
    const form = document.getElementById('disciplineForm');
    if (!form) return;

    // 1. Carrega o Select de Departamentos primeiro
    await loadDepartmentsSelect();

    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');
    const title = document.getElementById('formTitle');

    // 2. Se for edição, busca dados da disciplina
    if (id) {
        if (title) title.innerText = 'Editar Disciplina';
        document.getElementById('disciplineId').value = id;
        await loadDisciplineData(id);
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveDiscipline(id);
    });
}

async function loadDisciplineData(id) {
    try {
        const response = await fetch(`${API_URL}/${id}`);
        if (!response.ok) throw new Error('Erro na requisição');

        const d = await response.json();

        document.getElementById('name').value = d.name || '';
        document.getElementById('code').value = d.code || '';
        document.getElementById('credits').value = d.credits || '';
        document.getElementById('workload_hours').value = d.workload_hours || '';
        document.getElementById('description').value = d.description || '';

        if (document.getElementById('department_id')) {
            document.getElementById('department_id').value = d.department_id;
        }

    } catch (error) {
        console.error("Erro loadDisciplineData:", error);
        Swal.fire('Erro', 'Erro ao carregar dados da disciplina.', 'error');
    }
}

async function saveDiscipline(id) {
    const data = {
        name: document.getElementById('name').value,
        code: document.getElementById('code').value,
        credits: parseInt(document.getElementById('credits').value),
        workload_hours: parseInt(document.getElementById('workload_hours').value),
        description: document.getElementById('description').value || null,
        department_id: parseInt(document.getElementById('department_id').value)
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

            if (response.status === 409) {
                Swal.fire('Conflito', errText, 'warning');
                return;
            }

            throw new Error(errText || 'Erro na requisição');
        }

        await Swal.fire({
            title: 'Sucesso!',
            text: 'Disciplina salva com sucesso.',
            icon: 'success',
            timer: 1500,
            showConfirmButton: false
        });

        window.location.href = 'disciplines.html';
    } catch (error) {
        console.error(error);
        Swal.fire('Erro', 'Ocorreu um erro ao salvar: ' + error.message, 'error');
    }
}

// --- DETALHES (discipline_details.html) ---

async function initDetails() {
    const nameEl = document.getElementById('detail-name');
    if (!nameEl) return;

    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');

    if (!id) { window.location.href = 'disciplines.html'; return; }

    try {
        const response = await fetch(`${API_URL}/${id}`);
        if (!response.ok) throw new Error('Erro ao buscar disciplina');

        const d = await response.json();

        document.getElementById('detail-name').innerText = d.name;
        document.getElementById('detail-code').innerText = d.code;
        document.getElementById('detail-credits').innerText = d.credits;
        document.getElementById('detail-hours').innerText = d.workload_hours;
        document.getElementById('detail-desc').innerText = d.description || 'Sem descrição cadastrada.';

        document.getElementById('detail-dept').innerText = d.department_name || 'N/A';

        // Botão de Editar na página de detalhes
        const actionsDiv = document.getElementById('action-buttons');
        actionsDiv.innerHTML = `
            <a href="discipline_form.html?id=${d.id}" class="btn btn-warning text-white">
                <i class="bi bi-pencil-fill"></i> Editar
            </a>
        `;

    } catch (error) {
        console.error(error);
        Swal.fire('Erro', 'Erro ao carregar detalhes.', 'error');
    }
}