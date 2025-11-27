# 🎓 UniSystem - Sistema de Gestão Acadêmica

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![JavaScript](https://img.shields.io/badge/JavaScript-ES6+-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)
![Bootstrap](https://img.shields.io/badge/Bootstrap-5-7952B3?style=for-the-badge&logo=bootstrap&logoColor=white)
![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)

> **UniSystem** é uma plataforma Full-Stack desenvolvida para simplificar a administração de instituições de ensino. O projeto une a performance do **Go** no backend com a leveza do **Vanilla JavaScript** no frontend, criando uma solução rápida, moderna e sem dependências pesadas de frameworks SPA.

---

## 🚀 Visão Geral e Recursos

O sistema foi projetado seguindo o padrão de arquitetura REST, separando claramente as responsabilidades entre a interface do usuário e a lógica de negócios.

### ✨ Funcionalidades Principais
* **Dashboard Administrativo:** Painel visual com cartões animados para acesso rápido aos módulos.
* **Gestão Completa (CRUD):** Criação, Leitura, Atualização e Exclusão para:
    * 👨‍🎓 **Alunos:** Controle de matrículas, dados pessoais e associação a cursos.
    * 👨‍🏫 **Professores:** Cadastro detalhado com vínculo a departamentos.
    * 📚 **Cursos:** Definição de grade curricular e duração.
    * 🏢 **Departamentos:** Organização estrutural da instituição.
* **Soft Delete:** Implementação de inativação lógica (os dados não são perdidos, apenas arquivados), permitindo reativação futura.
* **Feedback Visual:** Utilização de *SweetAlert2* para confirmações e alertas elegantes, e validação de formulários com feedback em tempo real (ex: conflito de CPF).
* **Relacionamentos:** Integridade referencial entre Alunos/Cursos e Professores/Departamentos.

---

## 🛠️ Arquitetura e Tecnologias

### Stack Tecnológica
* **Backend:** Go (Golang) purista, utilizando apenas a biblioteca padrão `net/http` para roteamento e `lib/pq` para conexão com banco.
* **Frontend:** HTML5, CSS3 (Bootstrap 5) e JavaScript (Vanilla ES6+).
* **Banco de Dados:** PostgreSQL.
* **Design:** Responsivo e adaptável a dispositivos móveis.

### 📂 Estrutura de Arquivos
O projeto segue uma organização limpa e modular, facilitando a manutenção e escalabilidade:

```text
sistema-faculdade/
├── cmd/
│   └── api/
│       ├── main.go           # Ponto de entrada (Entrypoint), carrega envs e sobe o servidor
│       ├── routes.go         # Definição de rotas e servidor de arquivos estáticos
│       ├── students.go       # Handlers HTTP para Alunos
│       ├── teachers.go       # Handlers HTTP para Professores
│       ├── courses.go        # Handlers HTTP para Cursos
│       └── departments.go    # Handlers HTTP para Departamentos
├── internal/
│   ├── data/                 # Camada de Persistência (Repositories/SQL)
│   │   ├── student_repository.go
│   │   ├── teacher_repository.go
│   │   └── ...
│   └── models/               # Estruturas de Dados (Structs Go)
│       ├── student.go
│       └── ...
├── ui/                       # Frontend (Servido estaticamente pelo Go)
│   ├── css/                  # Estilos globais
│   ├── js/                   # Lógica do cliente (Fetch API, DOM manipulation)
│   │   ├── main.js           # Lógica do Dashboard
│   │   ├── students.js       # Lógica específica de Alunos
│   │   └── ...
│   ├── index.html            # Dashboard Principal
│   ├── students.html         # Listagem de Alunos
│   └── ...
├── .env                      # Variáveis de ambiente (Configuração do Banco)
└── README.md                 # Documentação do projeto
````

-----

## 🌐 Endpoints da API

A API segue o padrão RESTful. Abaixo os principais recursos disponíveis:

| Método | Endpoint | Descrição |
| :--- | :--- | :--- |
| **Alunos** | | |
| `GET` | `/api/students` | Lista todos os alunos (com paginação/filtros futuros). |
| `POST` | `/api/students` | Cria um novo aluno. Valida CPF e E-mail únicos. |
| `PUT` | `/api/students/{id}` | Atualiza dados do aluno. |
| `DELETE` | `/api/students/{id}` | Inativa o aluno (Soft Delete). |
| `PATCH` | `/api/students/{id}/activate` | Reativa um aluno inativo. |
| **Outros** | | |
| `GET` | `/api/courses` | Lista cursos para preencher dropdowns. |
| `GET` | `/api/departments` | Lista departamentos disponíveis. |

> *Nota: Endpoints similares existem para Professores, Cursos e Departamentos.*

-----

## 🏁 Guia de Instalação e Execução

### 1\. Pré-requisitos

  * [Go](https://go.dev/dl/) (v1.22+)
  * [PostgreSQL](https://www.postgresql.org/download/)

### 2\. Configuração do Banco de Dados

Crie um banco de dados no PostgreSQL e execute o script de criação das tabelas:

```sql
CREATE DATABASE unisystem_db;

-- Execute as tabelas (students, teachers, courses, departments)
-- Exemplo tabela Courses:
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    total_credits_required INT,
    duration_semesters INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- (Repita para as outras tabelas conforme os models do projeto)
```

### 3\. Configuração de Ambiente

Na raiz do projeto, crie um arquivo `.env` com a string de conexão do seu banco:

```env
DB_DSN=postgres://seu_usuario:sua_senha@localhost:5432/unisystem_db?sslmode=disable
```

### 4\. Executando a Aplicação

Navegue até a pasta `cmd/api` e inicie o servidor:

```bash
cd cmd/api
go run .
```

O terminal exibirá: `Servidor pronto! Conectado ao banco.`

### 5\. Acessando

Abra seu navegador e vá para:
👉 **http://localhost:8080**

-----

## 🔮 Roadmap Futuro

  * [ ] Implementação de Login/Auth (JWT).
  * [ ] Relatórios em PDF de matriculas.
  * [ ] Dashboard com gráficos (Chart.js) consumindo dados reais.
  * [ ] Paginação nas tabelas de listagem.

## 🤝 Contribuição

Contribuições são bem-vindas\! Sinta-se à vontade para abrir uma **Issue** para discutir novas features ou enviar um **Pull Request**.

-----

Desenvolvido usando Go e JavaScript.

```