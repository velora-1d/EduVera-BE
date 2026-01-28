package model

import "time"

// ==========================================
// E-RAPOR MODELS - Flexible Grading System
// ==========================================

// GradingComponent represents a single grading component (e.g., Sumatif, Formatif)
type GradingComponent struct {
	Name   string `json:"name"`   // "Sumatif", "Formatif", "Praktik"
	Weight int    `json:"weight"` // Percentage (0-100)
}

// GradingConfig represents the grading configuration for a subject
// Stored as JSONB in PostgreSQL for flexibility
type GradingConfig struct {
	UseKKM         bool               `json:"use_kkm"`         // K13 style with KKM threshold
	KKMValue       int                `json:"kkm_value"`       // Kriteria Ketuntasan Minimal (e.g., 75)
	UseDescriptive bool               `json:"use_descriptive"` // Kurikulum Merdeka style with descriptions
	Components     []GradingComponent `json:"components"`      // Weighted components
	PredicateRules []PredicateRule    `json:"predicate_rules"` // Rules for converting score to predicate
}

// PredicateRule defines how to convert numeric score to predicate
type PredicateRule struct {
	MinScore  int    `json:"min_score"` // Minimum score for this predicate
	MaxScore  int    `json:"max_score"` // Maximum score for this predicate
	Predicate string `json:"predicate"` // A, B, C, D
	Label     string `json:"label"`     // "Sangat Baik", "Baik", etc.
}

// Subject represents a school subject with its grading configuration
type Subject struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	Name          string        `json:"name"`           // "Matematika", "Bahasa Indonesia"
	Code          string        `json:"code"`           // "MATH-01"
	Type          string        `json:"type"`           // FORMAL_MERDEKA, FORMAL_K13, PESANTREN_KITAB
	GradingConfig GradingConfig `json:"grading_config"` // JSONB
	IsActive      bool          `json:"is_active"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// SubjectInput for creating/updating subjects
type SubjectInput struct {
	TenantID      string        `json:"tenant_id"`
	Name          string        `json:"name"`
	Code          string        `json:"code"`
	Type          string        `json:"type"`
	GradingConfig GradingConfig `json:"grading_config"`
}

// StudentGrade represents a student's grade for a subject in a semester
type StudentGrade struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	StudentID       string    `json:"student_id"`
	StudentName     string    `json:"student_name,omitempty"` // Joined from students table
	SubjectID       string    `json:"subject_id"`
	SubjectName     string    `json:"subject_name,omitempty"` // Joined from subjects table
	SemesterID      string    `json:"semester_id"`            // "2025-2026-1" (Year-Semester)
	ScoreNumeric    float64   `json:"score_numeric"`          // Final calculated score (0-100)
	ScorePredicate  string    `json:"score_predicate"`        // A/B/C/D
	DescriptionHigh string    `json:"description_high"`       // Kompetensi tertinggi
	DescriptionLow  string    `json:"description_low"`        // Kompetensi perlu ditingkatkan
	ComponentScores []float64 `json:"component_scores"`       // Scores per component (JSONB)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// StudentGradeInput for saving grades
type StudentGradeInput struct {
	TenantID        string    `json:"tenant_id"`
	StudentID       string    `json:"student_id"`
	SubjectID       string    `json:"subject_id"`
	SemesterID      string    `json:"semester_id"`
	ScoreNumeric    float64   `json:"score_numeric"`
	ScorePredicate  string    `json:"score_predicate"`
	DescriptionHigh string    `json:"description_high"`
	DescriptionLow  string    `json:"description_low"`
	ComponentScores []float64 `json:"component_scores"`
}

// BatchGradeInput for bulk saving grades
type BatchGradeInput struct {
	TenantID   string              `json:"tenant_id"`
	SemesterID string              `json:"semester_id"`
	SubjectID  string              `json:"subject_id"`
	Grades     []StudentGradeInput `json:"grades"`
}

// RaporData represents complete rapor data for a student
type RaporData struct {
	StudentID       string                `json:"student_id"`
	StudentName     string                `json:"student_name"`
	StudentNISN     string                `json:"student_nisn"`
	ClassName       string                `json:"class_name"`
	SemesterID      string                `json:"semester_id"`
	SemesterLabel   string                `json:"semester_label"` // "Semester 1 TA 2025/2026"
	Grades          []StudentGrade        `json:"grades"`
	Attendance      AttendanceData        `json:"attendance"`
	Extracurricular []ExtracurricularData `json:"extracurricular"`
	TeacherNotes    string                `json:"teacher_notes"`
}

// AttendanceData for rapor
type AttendanceData struct {
	Sakit int `json:"sakit"`
	Izin  int `json:"izin"`
	Alpha int `json:"alpha"`
}

// ExtracurricularData for rapor
type ExtracurricularData struct {
	Name        string `json:"name"`
	Predicate   string `json:"predicate"` // A/B/C
	Description string `json:"description"`
}

// Subject Types
const (
	SubjectTypeFormalMerdeka    = "FORMAL_MERDEKA"
	SubjectTypeFormalK13        = "FORMAL_K13"
	SubjectTypePesantrenKitab   = "PESANTREN_KITAB"
	SubjectTypePesantrenTahfidz = "PESANTREN_TAHFIDZ"
)

// Default grading config for Kurikulum Merdeka
func DefaultMerdekaGradingConfig() GradingConfig {
	return GradingConfig{
		UseKKM:         false,
		UseDescriptive: true,
		Components: []GradingComponent{
			{Name: "Sumatif", Weight: 60},
			{Name: "Formatif", Weight: 40},
		},
		PredicateRules: []PredicateRule{
			{MinScore: 90, MaxScore: 100, Predicate: "A", Label: "Sangat Baik"},
			{MinScore: 80, MaxScore: 89, Predicate: "B", Label: "Baik"},
			{MinScore: 70, MaxScore: 79, Predicate: "C", Label: "Cukup"},
			{MinScore: 0, MaxScore: 69, Predicate: "D", Label: "Perlu Bimbingan"},
		},
	}
}

// Default grading config for K13
func DefaultK13GradingConfig() GradingConfig {
	return GradingConfig{
		UseKKM:         true,
		KKMValue:       75,
		UseDescriptive: false,
		Components: []GradingComponent{
			{Name: "Pengetahuan", Weight: 50},
			{Name: "Keterampilan", Weight: 50},
		},
		PredicateRules: []PredicateRule{
			{MinScore: 90, MaxScore: 100, Predicate: "A", Label: "Sangat Baik"},
			{MinScore: 80, MaxScore: 89, Predicate: "B", Label: "Baik"},
			{MinScore: 70, MaxScore: 79, Predicate: "C", Label: "Cukup"},
			{MinScore: 0, MaxScore: 69, Predicate: "D", Label: "Kurang"},
		},
	}
}
