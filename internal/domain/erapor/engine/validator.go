package engine

import (
	"prabogo/internal/model"
)

// ValidationResult represents the output of a grading calculation
type ValidationResult struct {
	ScoreNumeric    float64
	ScorePredicate  string
	DescriptionHigh string // For Merdeka: "Shows mastery in..."
	DescriptionLow  string // For Merdeka: "Needs improvement in..."
	IsValid         bool
	Errors          []string
}

// CurriculumValidator defines the strategy interface for different curriculums
type CurriculumValidator interface {
	// CalculateGrade computes the final grade based on components
	CalculateGrade(components map[string]float64, config model.GradingConfig) ValidationResult

	// ValidateComponents checks if the input components match the required config
	ValidateComponents(components map[string]float64, config model.GradingConfig) error
}

// ValidatorFactory creates the appropriate validator based on curriculum code
func GetValidator(curriculumCode string) CurriculumValidator {
	switch curriculumCode {
	case model.SubjectTypeFormalK13:
		return &K13Validator{}
	case model.SubjectTypeFormalMerdeka:
		return &MerdekaValidator{}
	default:
		// Default validation (generic) or fallback
		return &GenericValidator{}
	}
}
