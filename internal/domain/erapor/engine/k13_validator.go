package engine

import (
	"prabogo/internal/model"
)

type K13Validator struct{}

func (v *K13Validator) CalculateGrade(components map[string]float64, config model.GradingConfig) ValidationResult {
	var totalScore float64
	var result ValidationResult
	var totalWeight int

	// 1. Calculate weighted average
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
		result.ScoreNumeric = 0
	}

	// 2. Determine Predicate based on KKM (K13 specific)
	// If KKM is 75:
	// < 75 : D
	// 75 - 82 : C
	// 83 - 90 : B
	// 91 - 100 : A
	// This formula is often dynamic based on KKM value

	kkm := float64(config.KKMValue)
	if kkm <= 0 {
		kkm = 75 // Default fallback
	}

	// Predicate logic (Standard K13 Interval)
	// Interval = (100 - KKM) / 3
	interval := (100 - kkm) / 3

	// A: > (100 - interval)
	// B: > (100 - 2*interval)
	// C: >= KKM
	// D: < KKM

	if result.ScoreNumeric < kkm {
		result.ScorePredicate = "D"
		result.DescriptionLow = "Perlu bimbingan dan remedial"
	} else if result.ScoreNumeric < (100 - 2*interval) {
		result.ScorePredicate = "C"
		result.DescriptionHigh = "Cukup menguasai materi"
	} else if result.ScoreNumeric < (100 - interval) {
		result.ScorePredicate = "B"
		result.DescriptionHigh = "Baik dalam menguasai materi"
	} else {
		result.ScorePredicate = "A"
		result.DescriptionHigh = "Sangat baik dalam menguasai materi"
	}

	result.IsValid = true
	return result
}

func (v *K13Validator) ValidateComponents(components map[string]float64, config model.GradingConfig) error {
	// Ensure all required components are present? Or allow partial?
	// K13 usually strictly requires Pengetahuan (Knowledge) & Keterampilan (Skill)
	// But configuration might vary per school.
	// We just validate against the config passed.
	for _, comp := range config.Components {
		if _, exists := components[comp.Name]; !exists {
			// For K13 we might want to be strict, but for now allow partial updates (zeros)
			// return fmt.Errorf("missing component: %s", comp.Name)
		}
	}
	return nil
}
