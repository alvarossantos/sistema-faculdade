// ui/js/main.js

const API_URL = '/api/students';

// --- LISTAGEM (index.html) ---
async function loadStudents() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Erro ao buscar dados');
        
        const students = await response.json();
        const tbody = document.querySelector('#studentsTable tbody');
        
        if (!tbody) return;

        tbody.innerHTML = '';

        if (!students || students.length === 0) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center py-4">Nenhum aluno encontrado.</td></tr>';
            return;
        }

        students.forEach(s => {
            const tr = document.createElement('tr');
            
            if (!s.active) {
                tr.classList.add('table-secondary', 'text-muted');
                tr.style.opacity = "0.75";
            }

            let actionButton;
            if (s.active) {
                actionButton = `
                <button onclick="deleteStudent(${s.id})" class="btn btn-sm btn-danger action-btn" title="Inativar">
                    <i class="bi bi-trash-fill"></i>
                </button>`;
            } else {
                actionButton = `
                <button onclick="activateStudent(${s.id})" class="btn btn-sm btn-success action-btn" title="Reativar">
                    <i class="bi bi-arrow-counterclockwise"></i>
                </button>`;
            }

            // Exibimos o course_name na tabela (Visualização)
            tr.innerHTML = `
                <td class="fw-bold">#${s.registration_number}</td>
                <td class="fw-bold">${s.name}</td>
                <td>${s.email || '-'}</td>
                <td><span class="badge bg-info text-dark">${s.course_name}</span></td>
                <td>${s.active ? '<span class="badge bg-success">Ativo</span>' : '<span class="badge bg-secondary">Inativo</span>'}</td>
                <td class="text-end">
                    <a href="students_form.html?id=${s.id}" class="btn btn-sm btn-warning action-btn" title="Editar">
                        <i class="bi bi-pencil-fill"></i>
                    </a>
                    ${actionButton}
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error(error);
        if(document.querySelector('#studentsTable')) {
            Swal.fire('Erro', 'Não foi possível carregar os alunos.', 'error');
        }
    }
}

// --- FUNÇÕES DE DELETE/ATIVAR (Mantidas) ---
async function deleteStudent(id) {
    const result = await Swal.fire({
        title: 'Inativar Aluno?',
        text: "O aluno não terá mais acesso ao sistema.",
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
                Swal.fire('Inativado!', 'O aluno foi inativado.', 'success');
                loadStudents();
            } else {
                throw new Error('Erro ao deletar');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao inativar.', 'error');
        }
    }
}

async function activateStudent(id) {
    const result = await Swal.fire({
        title: 'Reativar Aluno?',
        text: "O aluno voltará a ter status Ativo.",
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
                Swal.fire('Sucesso!', 'Aluno reativado.', 'success');
                loadStudents();
            } else {
                throw new Error('Erro ao ativar');
            }
        } catch (error) {
            Swal.fire('Erro!', 'Falha ao reativar o aluno.', 'error');
        }
    }
}

// --- LÓGICA DO CURSO (SELECT) ---

// Função para buscar lista de cursos no Backend
async function loadCoursesSelect() {
    try {
        const response = await fetch('/api/courses');
        if (!response.ok) throw new Error('Erro ao buscar cursos');
        
        const courses = await response.json();
        const select = document.getElementById('course_id');
        
        // Limpa e adiciona opção padrão
        select.innerHTML = '<option value="" selected disabled>Selecione um curso...</option>';

        courses.forEach(c => {
            const option = document.createElement('option');
            option.value = c.id; // Valor que vai pro banco (ID)
            option.textContent = `ID: ${c.id} - ${c.name}`; // Texto visível
            select.appendChild(option);
        });
    } catch (error) {
        console.error(error);
        const select = document.getElementById('course_id');
        if(select) select.innerHTML = '<option value="">Erro ao carregar cursos</option>';
    }
}


// --- FORMULÁRIO (form.html) ---

async function initForm() {
    const form = document.getElementById('studentForm');
    if (!form) return;

    // 1. Carrega o Select de cursos PRIMEIRO
    await loadCoursesSelect();

    const params = new URLSearchParams(window.location.search);
    const id = params.get('id');
    const title = document.getElementById('formTitle');

    if (id) {
        if(title) title.innerText = 'Editar Aluno';
        document.getElementById('studentId').value = id;
        // 2. Só depois carrega os dados do aluno
        await loadStudentData(id);
    }

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveStudent(id);
    });
}

async function loadStudentData(id) {
    try {
        const response = await fetch(`${API_URL}/${id}`);
        if (!response.ok) throw new Error('Erro na resposta da API');
        
        const s = await response.json();
        
        // Preenche os campos
        document.getElementById('name').value = s.name || '';
        document.getElementById('email').value = s.email || '';
        document.getElementById('cpf').value = s.cpf || '';
        document.getElementById('registration_number').value = s.registration_number || '';
        document.getElementById('gender').value = s.gender || 'M';
        
        // CORREÇÃO: Usa o ID do curso para selecionar a opção correta no Select
        if(document.getElementById('course_id')) {
            document.getElementById('course_id').value = s.course_id; 
        }

        if (s.date_birth) {
            try {
                const dateObj = new Date(s.date_birth);
                if (!isNaN(dateObj)) {
                    document.getElementById('date_birth').value = dateObj.toISOString().split('T')[0];
                }
            } catch (e) { console.warn("Erro data:", e); }
        }

    } catch (error) {
        console.error("Erro loadStudentData:", error);
        Swal.fire('Erro', 'Erro ao carregar dados.', 'error');
    }
}

// Função auxiliar para limpar erros visuais antigos
function clearErrors() {
    const inputs = document.querySelectorAll('.form-control, .form-select');
    inputs.forEach(input => {
        input.classList.remove('is-invalid');
    });
}

// Função auxiliar para mostrar erro em um campo específico
function showFieldError(fieldId, message) {
    const input = document.getElementById(fieldId);
    const errorDiv = document.getElementById('error-' + fieldId);
    
    if (input) {
        input.classList.add('is-invalid'); // Borda vermelha do Bootstrap
        if (errorDiv) {
            errorDiv.innerText = message; // Texto do erro
        }
    }
}

async function saveStudent(id) {
    // 1. Limpa erros anteriores
    clearErrors();

    // ... (sua lógica de captura de dados data = { ... } continua igual) ...
    const courseIdInput = document.getElementById('course_id');
    const data = {
        name: document.getElementById('name').value,
        email: document.getElementById('email').value || null,
        cpf: document.getElementById('cpf').value,
        registration_number: document.getElementById('registration_number').value,
        gender: document.getElementById('gender').value,
        course_id: parseInt(courseIdInput.value),
        date_birth: new Date(document.getElementById('date_birth').value).toISOString()
    };

    const method = id ? 'PUT' : 'POST';
    const url = id ? `${API_URL}/${id}` : API_URL;

    try {
        const response = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        // SE DER ERRO (Status diferente de 200-299)
        if (!response.ok) {
            const errText = await response.text(); // Pega o texto: "CPF já cadastrado"
            
            // Lógica para identificar qual campo deu erro baseado no texto do Go
            // As mensagens devem bater com as do student_repository.go
            
            let fieldErrorFound = false;

            if (errText.includes("CPF")) {
                showFieldError('cpf', errText); // Pinta o CPF
                fieldErrorFound = true;
            } 
            else if (errText.includes("email")) {
                showFieldError('email', errText); // Pinta o Email
                fieldErrorFound = true;
            } 
            else if (errText.includes("matrícula")) {
                showFieldError('registration_number', errText); // Pinta a Matrícula
                fieldErrorFound = true;
            }

            // Se achou um campo específico, não precisa de popup gigante, 
            // mas podemos mostrar um aviso pequeno "toast"
            if (fieldErrorFound) {
                Swal.fire({
                    toast: true,
                    position: 'top-end',
                    icon: 'warning',
                    title: 'Verifique os campos em vermelho.',
                    showConfirmButton: false,
                    timer: 3000
                });
            } else {
                // Erro genérico ou desconhecido
                throw new Error(errText);
            }
            return; // Para a execução aqui
        }

        // SUCESSO
        await Swal.fire({
            title: 'Sucesso!',
            text: 'Dados salvos.',
            icon: 'success',
            timer: 1500,
            showConfirmButton: false
        });

        window.location.href = 'students.html';

    } catch (error) {
        console.error(error);
        Swal.fire('Erro', 'Ocorreu um erro: ' + error.message, 'error');
    }
}