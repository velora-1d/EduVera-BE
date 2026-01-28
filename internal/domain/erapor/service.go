package erapor

import (
	"context"

	"prabogo/internal/domain/erapor/engine"
	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

// Service adalah domain service untuk E-Rapor
type Service struct {
	db outbound_port.ERaporDatabasePort
}

// NewService membuat instance baru Service
func NewService(db outbound_port.ERaporDatabasePort) *Service {
	return &Service{db: db}
}

// ==========================================
// SUBJECT OPERATIONS
// ==========================================

// CreateSubject membuat mata pelajaran baru dengan konfigurasi grading
func (s *Service) CreateSubject(ctx context.Context, input *model.SubjectInput) (*model.Subject, error) {
	// Apply default grading config if empty
	if len(input.GradingConfig.Components) == 0 {
		if input.Type == model.SubjectTypeFormalK13 {
			input.GradingConfig = model.DefaultK13GradingConfig()
		} else {
			input.GradingConfig = model.DefaultMerdekaGradingConfig()
		}
	}
	return s.db.CreateSubject(ctx, input)
}

// UpdateSubject mengupdate mata pelajaran
func (s *Service) UpdateSubject(ctx context.Context, id string, input *model.SubjectInput) (*model.Subject, error) {
	return s.db.UpdateSubject(ctx, id, input)
}

// GetSubjectByID mengambil mata pelajaran berdasarkan ID
func (s *Service) GetSubjectByID(ctx context.Context, id string) (*model.Subject, error) {
	return s.db.GetSubjectByID(ctx, id)
}

// GetSubjectsByTenant mengambil semua mata pelajaran tenant
func (s *Service) GetSubjectsByTenant(ctx context.Context, tenantID string) ([]model.Subject, error) {
	return s.db.GetSubjectsByTenant(ctx, tenantID)
}

// DeleteSubject menghapus mata pelajaran (soft delete)
func (s *Service) DeleteSubject(ctx context.Context, id string) error {
	return s.db.DeleteSubject(ctx, id)
}

// ==========================================
// GRADE OPERATIONS
// ==========================================

// SaveGrade menyimpan nilai siswa dengan kalkulasi predicate otomatis
func (s *Service) SaveGrade(ctx context.Context, input *model.StudentGradeInput) (*model.StudentGrade, error) {
	// Get subject to apply grading rules
	subject, err := s.db.GetSubjectByID(ctx, input.SubjectID)
	if err != nil {
		return nil, err
	}

	// Auto-calculate using Validator Engine
	// Use Subject Type (e.g., FORMAL_MERDEKA) to determining validation strategy
	validator := engine.GetValidator(subject.Type)

	// Map component slice to map for validator
	compMap := make(map[string]float64)
	for i, val := range input.ComponentScores {
		if i < len(subject.GradingConfig.Components) {
			compMap[subject.GradingConfig.Components[i].Name] = val
		}
	}

	result := validator.CalculateGrade(compMap, subject.GradingConfig)

	// Apply calculated values if not manually overridden
	if input.ScoreNumeric == 0 && result.ScoreNumeric > 0 {
		input.ScoreNumeric = result.ScoreNumeric
	}
	if input.ScorePredicate == "" {
		input.ScorePredicate = result.ScorePredicate
	}
	if input.DescriptionHigh == "" {
		input.DescriptionHigh = result.DescriptionHigh
	}
	if input.DescriptionLow == "" {
		input.DescriptionLow = result.DescriptionLow
	}

	return s.db.SaveGrade(ctx, input)
}

// BatchSaveGrades menyimpan banyak nilai sekaligus
func (s *Service) BatchSaveGrades(ctx context.Context, input *model.BatchGradeInput) ([]model.StudentGrade, error) {
	// Get subject for predicate calculation
	subject, err := s.db.GetSubjectByID(ctx, input.SubjectID)
	if err != nil {
		return nil, err
	}

	// Auto-calculate predicates using Validator Engine
	validator := engine.GetValidator(subject.Type)

	for i := range input.Grades {
		// Map component slice to map
		compMap := make(map[string]float64)
		for j, val := range input.Grades[i].ComponentScores {
			if j < len(subject.GradingConfig.Components) {
				compMap[subject.GradingConfig.Components[j].Name] = val
			}
		}

		result := validator.CalculateGrade(compMap, subject.GradingConfig)

		if input.Grades[i].ScoreNumeric == 0 && result.ScoreNumeric > 0 {
			input.Grades[i].ScoreNumeric = result.ScoreNumeric
		}
		if input.Grades[i].ScorePredicate == "" {
			input.Grades[i].ScorePredicate = result.ScorePredicate
		}
	}

	return s.db.BatchSaveGrades(ctx, input)
}

// GetGradesByStudent mengambil semua nilai siswa di semester tertentu
func (s *Service) GetGradesByStudent(ctx context.Context, studentID, semesterID string) ([]model.StudentGrade, error) {
	return s.db.GetGradesByStudent(ctx, studentID, semesterID)
}

// GetGradesBySubject mengambil semua nilai untuk mata pelajaran tertentu
func (s *Service) GetGradesBySubject(ctx context.Context, subjectID, semesterID string) ([]model.StudentGrade, error) {
	return s.db.GetGradesBySubject(ctx, subjectID, semesterID)
}

// GetStudentRapor mengambil data rapor lengkap untuk siswa
func (s *Service) GetStudentRapor(ctx context.Context, studentID, semesterID string) (*model.RaporData, error) {
	return s.db.GetStudentRapor(ctx, studentID, semesterID)
}

// GetGradeStats mengambil statistik nilai
func (s *Service) GetGradeStats(ctx context.Context, tenantID, semesterID string) (map[string]interface{}, error) {
	return s.db.GetGradeStats(ctx, tenantID, semesterID)
}

// ==========================================
// UTILITY FUNCTIONS
// ==========================================

// calculatePredicate menghitung predikat berdasarkan nilai dan config
func (s *Service) calculatePredicate(score float64, config model.GradingConfig) string {
	scoreInt := int(score)
	for _, rule := range config.PredicateRules {
		if scoreInt >= rule.MinScore && scoreInt <= rule.MaxScore {
			return rule.Predicate
		}
	}
	// Default fallback
	if scoreInt >= 90 {
		return "A"
	} else if scoreInt >= 80 {
		return "B"
	} else if scoreInt >= 70 {
		return "C"
	}
	return "D"
}

// CalculateFinalScore menghitung nilai akhir dari component scores
func (s *Service) CalculateFinalScore(componentScores []float64, config model.GradingConfig) float64 {
	if len(componentScores) == 0 || len(config.Components) == 0 {
		return 0
	}

	var total float64
	var totalWeight int

	for i, component := range config.Components {
		if i < len(componentScores) {
			total += componentScores[i] * float64(component.Weight)
			totalWeight += component.Weight
		}
	}

	if totalWeight == 0 {
		return 0
	}

	return total / float64(totalWeight)
}

// GenerateDescription menghasilkan deskripsi capaian kompetensi
func (s *Service) GenerateDescription(score float64, subjectName string) (high, low string) {
	if score >= 90 {
		high = "Ananda menunjukkan penguasaan yang sangat baik pada materi " + subjectName
		low = "Tetap pertahankan dan tingkatkan terus prestasi belajar"
	} else if score >= 80 {
		high = "Ananda menunjukkan pemahaman yang baik terhadap konsep-konsep " + subjectName
		low = "Perlu latihan lebih untuk mencapai hasil yang maksimal"
	} else if score >= 70 {
		high = "Ananda memahami konsep dasar " + subjectName
		low = "Perlu bimbingan tambahan untuk penguatan materi"
	} else {
		high = "Ananda sudah berusaha dalam mempelajari " + subjectName
		low = "Perlu pendampingan intensif untuk peningkatan pemahaman"
	}
	return high, low
}

// GenerateRapor generates a persistent snapshot of the rapor
func (s *Service) GenerateRapor(ctx context.Context, tenantID, studentID, semesterID string, catatanWali string) (*model.Rapor, error) {
	// 1. Get or Create Rapor Periode
	periode, err := s.db.GetOrCreateRaporPeriode(tenantID, semesterID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch Calculated Grades (Dynamic)
	raporData, err := s.db.GetStudentRapor(ctx, studentID, semesterID)
	if err != nil {
		return nil, err
	}

	// 3. Create Rapor Header
	raporHeader := &model.Rapor{
		TenantID:         tenantID,
		PeriodeID:        periode.ID,
		SantriID:         studentID,
		Status:           "Draft",
		CatatanWaliKelas: catatanWali,
	}
	if err := s.db.CreateRapor(raporHeader); err != nil {
		return nil, err
	}

	// 4. Create Rapor Details (Nilai)
	for _, grade := range raporData.Grades {
		nilai := &model.RaporNilai{
			RaporID:    raporHeader.ID,
			Kategori:   "Akademik", // Default category
			Jenis:      grade.SubjectName,
			Nilai:      grade.ScorePredicate, // Use predicate (A, B, C) for rapor
			Keterangan: grade.DescriptionHigh + ". " + grade.DescriptionLow,
		}
		// If subject has specific category requirement, can be adjusted here.
		if err := s.db.CreateRaporNilai(nilai); err != nil {
			// Continue or fail? Ideally transaction.
			// For now log error but continue
			continue
		}
	}

	return raporHeader, nil
}
