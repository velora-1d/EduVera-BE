package engine

import (
	"prabogo/internal/model"
)

// GenericValidator serves as a fallback for unknown curriculums or custom/pesantren usage
type GenericValidator struct{}

func (v *GenericValidator) CalculateGrade(components map[string]float64, config model.GradingConfig) ValidationResult {
	var totalScore float64
	var result ValidationResult
	var totalWeight int

	for _, comp := range config.Components {
		score, exists := components[comp.Name]
		if exists {
			totalScore += score * float64(comp.Weight)
			totalWeight += comp.Weight
		}
	}

	if totalWeight > 0 {
		result.ScoreNumeric = totalScore / float64(totalWeight)
	} else {
		result.ScoreNumeric = 0 // Or just average if weight is 0? Assume weight is mandatory for now.
	}

	// Simple 10-point scale or arbitrary mapping if no config
	if len(config.PredicateRules) > 0 {
		for _, rule := range config.PredicateRules {
			if result.ScoreNumeric >= float64(rule.MinScore) && result.ScoreNumeric <= float64(rule.MaxScore) {
				result.ScorePredicate = rule.Predicate
				result.DescriptionHigh = rule.Label
				break
			}
		}
	} else {
		// Fallback standard grading
		if result.ScoreNumeric >= 85 {
			result.ScorePredicate = "A"
		} else if result.ScoreNumeric >= 70 {
			result.ScorePredicate = "B"
		} else if result.ScoreNumeric >= 60 {
			result.ScorePredicate = "C"
		} else {
			result.ScorePredicate = "D"
		}
	}

	result.IsValid = true
	return result
}

func (v *GenericValidator) ValidateComponents(components map[string]float64, config model.GradingConfig) error {
	return nil
}
