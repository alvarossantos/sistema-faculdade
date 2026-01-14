const API_URL = '/api/offers';

// --- LISTAGEM (offers.html) ---

async function loadOffers() {
    try {
        const response = await fetch(API_URL);
        
        if (!response.ok) throw new Error('Erro ao buscar ofertas');
        
        const offers = await response.json();
        const tbody = document.querySelector('#offersTable tbody');
        
        if (!tbody) return;

        tbody.innerHTML = '';

        if (!offers || offers.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center py-4">Nenhuma turma ofertada.</td></tr>';
            return;
        }

        offers.forEach(o => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td><span class="badge bg-info text-dark">${o.semester_label}</span></td>
                <td><span class="badge bg-secondary">${o.class_code}</span></td>
                <td class="fw-bold">${o.discipline_name}</td>
                <td>${o.teacher_name}</td>
                <td class="d-none d-md-table-cell"><small>${o.schedule}</small></td>
                <td class="text-end">
                    <button onclick="deleteOffer(${o.id})" class="btn btn-sm btn-danger action-btn" title="Excluir">
                        <i class="bi bi-trash-fill"></i>
                    </button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
        const tbody = document.querySelector('#offersTable tbody');
        if(tbody) tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger">Erro ao carregar dados.</td></tr>';
    }
}

// --- FORMULÁRIO (offer_form.html) ---

async function initForm() {
    const form = document.getElementById('offerForm');
    if (!form) return;

    // Carrega as dependências (Selects)
    try {
        await Promise.all([
            loadSelect('/api/semesters', 'semester_id', (item) => `${item.year}.${item.period}`),
            loadSelect('/api/disciplines', 'discipline_id', 'name'),
            loadSelect('/api/teachers', 'teacher_id', 'name')
        ]);
    } catch (error) {
        console.error(error);
        Swal.fire('Erro', 'Falha ao carregar listas (Semestres, Disciplinas ou Professores).', 'error');
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        // Mapeia os dados conforme o struct DisciplineOffer em offer.go
        const data = {
            semester_id: parseInt(document.getElementById('semester_id').value),
            discipline_id: parseInt(document.getElementById('discipline_id').value),
            teacher_id: parseInt(document.getElementById('teacher_id').value),
            class_code: document.getElementById('class_code').value,
            schedule: document.getElementById('schedule').value
        };

        try {
            const response = await fetch(API_URL, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            });

            if (!response.ok) {
                const errText = await response.text();
                // Tratamento de erro 409 (se você implementar validação de duplicidade)
                if (errText.includes("existe")) {
                    Swal.fire('Atenção', 'Esta disciplina já foi ofertada neste semestre.', 'warning');
                    return;
                }
                throw new Error('Erro ao salvar');
            }

            await Swal.fire({
                title: 'Sucesso!',
                text: 'Turma criada com sucesso.',
                icon: 'success',
                timer: 1500,
                showConfirmButton: false
            });

            window.location.href = 'offers.html';

        } catch (error) {
            console.error(error);
            Swal.fire('Erro', 'Ocorreu um erro ao salvar.', 'error');
        }
    });
}

// Função auxiliar para preencher selects
async function loadSelect(url, elementId, textProperty) {
    const response = await fetch(url);
    if (!response.ok) throw new Error(`Erro ao carregar ${url}`);
    
    const items = await response.json();
    const select = document.getElementById(elementId);
    
    if (!select) return;

    select.innerHTML = '<option value="" selected disabled>Selecione...</option>';

    if (!items) return;

    items.forEach(item => {
        const option = document.createElement('option');
        option.value = item.id;
        
        // Verifica se a propriedade de texto é uma função ou string
        if (typeof textProperty === 'function') {
            option.textContent = textProperty(item);
        } else {
            option.textContent = item[textProperty];
        }
        select.appendChild(option);
    });
}

// --- DELETE ---

async function deleteOffer(id) {
    const result = await Swal.fire({
        title: 'Excluir Oferta?',
        text: "Isso removerá a turma do sistema.",
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
                Swal.fire('Excluído!', 'Turma removida.', 'success');
                loadOffers();
            } else {
                throw new Error('Erro ao excluir');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao processar solicitação.', 'error');
        }
    }
}