-- E-Rapor Module Tables Migration
-- Created: 2026-01-28
-- Description: Create tables for Subject, StudentGrade, and related E-Rapor entities

-- =============================================
-- SUBJECTS TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    grade_level INTEGER NOT NULL, -- Tingkat kelas (10, 11, 12)
    curriculum VARCHAR(50) NOT NULL, -- kurikulum_merdeka, k13
    is_muatan_lokal BOOLEAN DEFAULT FALSE,
    grading_config JSONB DEFAULT '{}', -- Konfigurasi penilaian polymorphic
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for subjects
CREATE INDEX IF NOT EXISTS idx_subjects_tenant_id ON subjects(tenant_id);
CREATE INDEX IF NOT EXISTS idx_subjects_grade_level ON subjects(grade_level);
CREATE INDEX IF NOT EXISTS idx_subjects_curriculum ON subjects(curriculum);
CREATE UNIQUE INDEX IF NOT EXISTS idx_subjects_code_tenant ON subjects(tenant_id, code);

-- =============================================
-- STUDENT_GRADES TABLE
-- =============================================
CREATE TABLE IF NOT EXISTS student_grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    student_id UUID NOT NULL, -- Reference to students table
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    semester INTEGER NOT NULL CHECK (semester IN (1, 2)),
    academic_year VARCHAR(9) NOT NULL, -- Format: 2025/2026
    
    -- Nilai Pengetahuan
    nilai_tugas DECIMAL(5,2),
    nilai_ulangan_harian DECIMAL(5,2),
    nilai_uts DECIMAL(5,2),
    nilai_uas DECIMAL(5,2),
    nilai_pengetahuan DECIMAL(5,2), -- Calculated
    
    -- Nilai Keterampilan
    nilai_praktik DECIMAL(5,2),
    nilai_proyek DECIMAL(5,2),
    nilai_portofolio DECIMAL(5,2),
    nilai_keterampilan DECIMAL(5,2), -- Calculated
    
    -- Final
    nilai_akhir DECIMAL(5,2),
    predikat CHAR(1), -- A, B, C, D
    
    -- Deskripsi (Kurikulum Merdeka)
    deskripsi_pengetahuan TEXT,
    deskripsi_keterampilan TEXT,
    catatan_guru TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for student_grades
CREATE INDEX IF NOT EXISTS idx_student_grades_tenant_id ON student_grades(tenant_id);
CREATE INDEX IF NOT EXISTS idx_student_grades_student_id ON student_grades(student_id);
CREATE INDEX IF NOT EXISTS idx_student_grades_subject_id ON student_grades(subject_id);
CREATE INDEX IF NOT EXISTS idx_student_grades_semester ON student_grades(semester);
CREATE INDEX IF NOT EXISTS idx_student_grades_academic_year ON student_grades(academic_year);

-- Unique constraint: one grade per student per subject per semester per year
CREATE UNIQUE INDEX IF NOT EXISTS idx_student_grades_unique 
    ON student_grades(student_id, subject_id, semester, academic_year);

-- =============================================
-- TRIGGERS
-- =============================================

-- Subjects
DROP TRIGGER IF EXISTS update_subjects_updated_at ON subjects;
CREATE TRIGGER update_subjects_updated_at
    BEFORE UPDATE ON subjects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Student Grades
DROP TRIGGER IF EXISTS update_student_grades_updated_at ON student_grades;
CREATE TRIGGER update_student_grades_updated_at
    BEFORE UPDATE ON student_grades
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
