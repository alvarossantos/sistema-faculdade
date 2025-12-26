-- =========================================================
-- DATABASE
-- =========================================================
-- CREATE DATABASE sistema_faculdade;

-- =========================================================
-- TABELA DE DEPARTAMENTOS
-- =========================================================
CREATE TABLE departments (
  id SERIAL PRIMARY KEY,
  name VARCHAR(120) UNIQUE NOT NULL,
  abbreviation VARCHAR(20) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- TABELA DE CURSOS
-- =========================================================
CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  name VARCHAR(120) UNIQUE NOT NULL,
  total_credits_required INT NOT NULL CHECK(total_credits_required > 0),
  duration_semesters INT NOT NULL CHECK(duration_semesters > 0),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- TABELA DE PROFESSORES
-- =========================================================
CREATE TABLE teachers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(120) UNIQUE NOT NULL CHECK (email LIKE '%@%'),
  cpf CHAR(11) UNIQUE NOT NULL CHECK (cpf ~ '^[0-9]{11}$'),
  telephone VARCHAR(20) NOT NULL,
  active BOOLEAN DEFAULT TRUE NOT NULL,
  department_id INT REFERENCES departments(id) ON DELETE RESTRICT,
  date_contract DATE CHECK (date_contract <= CURRENT_DATE),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- TABELA DE ALUNOS
-- =========================================================
CREATE TABLE students (
  id SERIAL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(120) UNIQUE,
  gender CHAR(1) CHECK (gender IN ('M','F','O')),
  date_birth DATE NOT NULL CHECK (date_birth <= CURRENT_DATE),
  cpf CHAR(11) UNIQUE NOT NULL CHECK (cpf ~ '^[0-9]{11}$'),
  registration_number VARCHAR(30) UNIQUE NOT NULL,
  active BOOLEAN DEFAULT TRUE NOT NULL,
  course_id INT REFERENCES courses(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- TABELA DE DISCIPLINAS
-- =========================================================
CREATE TABLE disciplines (
  id SERIAL PRIMARY KEY,
  name VARCHAR(120) UNIQUE NOT NULL,
  code VARCHAR(20) UNIQUE NOT NULL,
  credits INT NOT NULL CHECK(credits > 0),
  workload_hours INT NOT NULL CHECK(workload_hours > 0),
  description TEXT,
  department_id INT REFERENCES departments(id),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- TABELA DE SEMESTRES ACADÊMICOS
-- =========================================================
CREATE TABLE academic_semesters (
  id SERIAL PRIMARY KEY,
  year INT NOT NULL,
  period SMALLINT NOT NULL CHECK(period IN (1,2)),
  UNIQUE (year, period)
);

-- =========================================================
-- OFERTA DE DISCIPLINAS
-- =========================================================
CREATE TABLE discipline_offers (
  id SERIAL PRIMARY KEY,
  discipline_id INT REFERENCES disciplines(id) ON DELETE CASCADE,
  semester_id INT REFERENCES academic_semesters(id) ON DELETE CASCADE,
  teacher_id INT REFERENCES teachers(id) ON DELETE RESTRICT,
  schedule VARCHAR(100) NOT NULL,
  UNIQUE(discipline_id, semester_id)
);

-- =========================================================
-- STATUS POSSÍVEL DO ALUNO NA DISCIPLINA
-- =========================================================
CREATE TYPE registration_status AS ENUM (
  'in_progress',
  'approved',
  'failed',
  'take_test'
);

-- =========================================================
-- MATRÍCULAS DO ALUNO EM DISCIPLINAS OFERTADAS
-- =========================================================
CREATE TABLE registrations (
  id SERIAL PRIMARY KEY,
  student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
  offer_id INT NOT NULL REFERENCES discipline_offers(id) ON DELETE CASCADE,
  final_grade DECIMAL(5,2),
  frequency DECIMAL(5,2),
  absences INT DEFAULT 0,
  status registration_status DEFAULT 'in_progress' NOT NULL,
  approved BOOLEAN DEFAULT FALSE NOT NULL,
  UNIQUE(student_id, offer_id),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- LANÇAMENTOS DE NOTAS
-- =========================================================
CREATE TABLE grade_items (
  id SERIAL PRIMARY KEY,
  registration_id INT NOT NULL REFERENCES registrations(id) ON DELETE CASCADE,
  title VARCHAR(100) NOT NULL,
  grade DECIMAL(5,2) NOT NULL CHECK(grade >= 0 AND grade <= 100),
  weight DECIMAL(5,2) NOT NULL CHECK(weight >= 0 AND weight <= 1),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- =========================================================
-- REGISTRO DE FALTAS
-- =========================================================
CREATE TABLE attendance_records (
  id SERIAL PRIMARY KEY,
  registration_id INT NOT NULL REFERENCES registrations(id) ON DELETE CASCADE,
  class_date DATE NOT NULL,
  hours_absent INT NOT NULL CHECK(hours_absent >= 0),
  UNIQUE(registration_id, class_date)
);

-- =========================================================
-- FUNÇÃO: CALCULAR NOTA FINAL
-- =========================================================
CREATE OR REPLACE FUNCTION calculate_final_grade(reg_id INT)
RETURNS DECIMAL AS $$
DECLARE
  result DECIMAL;
BEGIN
  SELECT COALESCE(SUM(grade * weight), 0)
  INTO result
  FROM grade_items
  WHERE registration_id = reg_id;

  RETURN result;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- FUNÇÃO: CALCULAR FREQUÊNCIA E FALTAS
-- =========================================================
CREATE OR REPLACE FUNCTION calculate_attendance(reg_id INT)
RETURNS TABLE(absences INT, frequency DECIMAL) AS $$
DECLARE
  workload INT;
BEGIN
  SELECT COALESCE(SUM(hours_absent), 0)
  INTO absences
  FROM attendance_records
  WHERE registration_id = reg_id;

  SELECT d.workload_hours
  INTO workload
  FROM registrations r
  JOIN discipline_offers o ON o.id = r.offer_id
  JOIN disciplines d ON d.id = o.discipline_id
  WHERE r.id = reg_id;

  IF workload = 0 THEN
    frequency := 0;
  ELSE
    frequency := 100 - ((absences::DECIMAL / workload) * 100);
  END IF;

  RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- FUNÇÃO: ATUALIZAR STATUS DA MATRÍCULA
-- =========================================================
CREATE OR REPLACE FUNCTION update_registration_status()
RETURNS TRIGGER AS $$
DECLARE
  calc_grade DECIMAL;
  abs INT;
  freq DECIMAL;
  reg_id INT;
BEGIN
  -- Determina o ID correto, OLD para DELETE, NEW para INSERT/UPDATE
  IF TG_OP = 'DELETE' THEN
    reg_id := OLD.registration_id;
  ELSE
    reg_id := NEW.registration_id;
  END IF;

  calc_grade := calculate_final_grade(reg_id);
  SELECT absences, frequency INTO abs, freq
  FROM calculate_attendance(reg_id);

  IF TG_OP = 'DELETE' THEN
    RETURN OLD;
  ELSE
    NEW.final_grade := calc_grade;
    NEW.absences := abs;
    NEW.frequency := freq;

    IF freq < 75 THEN
      NEW.status := 'failed';
      NEW.approved := FALSE;
    ELSIF calc_grade >= 60 THEN
      NEW.status := 'approved';
      NEW.approved := TRUE;
    ELSIF calc_grade < 60 AND freq >= 75 THEN
      NEW.status := 'take_test';
      NEW.approved := FALSE;
    ELSE
      NEW.status := 'failed';
      NEW.approved := FALSE;
    END IF;

    RETURN NEW;
  END IF;
END;
$$ LANGUAGE plpgsql;

-- =========================================================
-- TRIGGERS: ATUALIZAR STATUS AUTOMATICAMENTE
-- =========================================================
CREATE TRIGGER trg_update_registration_after_grades
AFTER INSERT OR UPDATE OR DELETE ON grade_items
FOR EACH ROW
EXECUTE PROCEDURE update_registration_status();

CREATE TRIGGER trg_update_registration_after_attendance
AFTER INSERT OR UPDATE OR DELETE ON attendance_records
FOR EACH ROW
EXECUTE PROCEDURE update_registration_status();

-- =========================================================
-- FUNÇÃO PARA ATUALIZAR updated_at
-- =========================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_teachers_modtime
BEFORE UPDATE ON teachers
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_students_modtime
BEFORE UPDATE ON students
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_disciplines_modtime
BEFORE UPDATE ON disciplines
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_registrations_modtime
BEFORE UPDATE ON registrations
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =========================================================
-- DADOS DE TESTE
-- =========================================================
-- Departamentos
INSERT INTO departments (name, abbreviation)
VALUES
('Tecnologia da Informação', 'TI'),
('Computação', 'COMP');

-- Cursos
INSERT INTO courses (name, total_credits_required, duration_semesters)
VALUES
('Sistemas de Informação', 200, 8),
('Ciência da Computação', 220, 8);

-- Professores
INSERT INTO teachers (name, email, cpf, telephone, department_id)
VALUES
('Carlos Souza', 'carlos.souza@facul.com', '12345678901', '11999999999', 1),
('Mariana Lima', 'mariana.lima@facul.com', '98765432100', '21988888888', 2);

-- Alunos
INSERT INTO students (name, date_birth, cpf, registration_number, email, gender, course_id)
VALUES
('João Pereira', '2000-05-10', '11122233344', '2025001', 'joao@aluno.com', 'M', 1),
('Ana Santos', '1999-10-02', '55566677788', '2025002', 'ana@aluno.com', 'F', 1),
('Pedro Alves', '2001-02-20', '99988877766', '2025003', 'pedro@aluno.com', 'M', 2);

-- Semestre
INSERT INTO academic_semesters (year, period)
VALUES (2025, 1);

-- Disciplinas
INSERT INTO disciplines (name, code, credits, workload_hours, description, department_id)
VALUES
('Banco de Dados', 'BD101', 4, 80, 'Fundamentos e modelagem', 1),
('Programação Go', 'GO202', 3, 60, 'Programação moderna com Go', 2),
('Algoritmos', 'ALG303', 4, 80, 'Introdução à lógica', 1);

-- Ofertas
INSERT INTO discipline_offers (discipline_id, semester_id, teacher_id, schedule)
VALUES
(1, 1, 1, 'Seg/Qua 10h'),
(2, 1, 2, 'Ter/Qui 14h'),
(3, 1, 1, 'Seg/Qua 08h');

-- Matrículas
INSERT INTO registrations (student_id, offer_id)
VALUES
(1, 1),
(1, 2),
(2, 1),
(3, 3);
