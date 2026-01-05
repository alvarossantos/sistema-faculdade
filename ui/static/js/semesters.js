const API_URL = '/api/semesters';

// --- LISTAGEM ---
async function loadSemesters() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Erro ao buscar semestres');

        const semesters = await response.json();
        const tbody = document.querySelector('#semestersTable tbody');

        if (!tbody) return;

        tbody.innerHTML = '';

        if (!semesters || semesters.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="text-center py-4">Nenhum semestre cadastrado.</td></tr>';
            return;
        }

        semesters.forEach(s => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td class="fw-bold">${s.year}</td>
                <td>${s.period}º</td>
                <td><span class="badge bg-info text-dark">${s.year}.${s.period}</span></td>
                <td class="text-end">
                    <button onclick="deleteSemester(${s.id})" class="btn btn-sm btn-danger action-btn" title="Excluir">
                        <i class="bi bi-trash-fill"></i>
                    </button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
        const tbody = document.querySelector('#semestersTable tbody');
        if (tbody) tbody.innerHTML = '<tr><td colspan="4" class="text-center text-danger">Erro ao carregar dados.</td></tr>';
    }
}

// --- CRIAR ---
async function initForm() {
    const form = document.getElementById('semesterForm');
    if (!form) return;

    // Define ano atual como padrão
    document.getElementById('year').value = new Date().getFullYear();

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const data = {
            year: parseInt(document.getElementById('year').value),
            period: parseInt(document.getElementById('period').value)
        };

        try {
            const response = await fetch(API_URL, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                const errText = await response.text();
                // O backend retorna erro 500 se já existir (devia ser 409, mas vamos tratar genericamente)
                if (errText.includes("já existe")) {
                    Swal.fire('Atenção', 'Este semestre já está cadastrado.', 'warning');
                    return;
                }
                throw new Error('Erro ao criar');
            }

            await Swal.fire({
                title: 'Sucesso!',
                text: 'Semestre criado.',
                icon: 'success',
                timer: 1500,
                showConfirmButton: false
            });

            window.location.href = 'semesters.html';
        } catch (error) {
            console.error(error);
            Swal.fire('Erro', 'Falha ao salvar semestre.', 'error');
        }
    });
}

// --- EXCLUIR ---
async function deleteSemester(id) {
    const result = await Swal.fire({
        title: 'Excluir Semestre?',
        text: "Esta ação não pode ser desfeita.",
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
                Swal.fire('Excluído!', 'Semestre removido.', 'success');
                loadSemesters();
            } else {
                const errText = await response.text();
                if (errText.includes("ofertas vinculadas")) {
                    Swal.fire('Bloqueado', 'Não é possível excluir: existem ofertas vinculadas a este semestre.', 'error');
                } else {
                    throw new Error('Erro ao excluir');
                }
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao processar solicitação.', 'error');
        }
    }
}