package engine

import (
	"prabogo/internal/model"
)

type MerdekaValidator struct{}

func (v *MerdekaValidator) CalculateGrade(components map[string]float64, config model.GradingConfig) ValidationResult {
	var totalScore float64
	var result ValidationResult
	var totalWeight int

	// 1. Calculate weighted average (Sumatif + Formatif)
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

	// 2. Merdeka Logic: Focus on Descriptions (Capaian Pembelajaran)
	// Predicate is secondary but still used for simplified view
	// Usually strict ranges not dependent on KKM (KKM is technically removed in Merdeka)

	if result.ScoreNumeric >= 91 {
		result.ScorePredicate = "A"
		result.DescriptionHigh = "Menunjukkan penguasaan yang sangat baik dalam " // To be appended with TP (Tujuan Pembelajaran)
	} else if result.ScoreNumeric >= 81 {
		result.ScorePredicate = "B"
		result.DescriptionHigh = "Menunjukkan penguasaan yang baik dalam "
	} else if result.ScoreNumeric >= 71 {
		result.ScorePredicate = "C"
		result.DescriptionHigh = "Menunjukkan penguasaan yang cukup dalam "
	} else {
		result.ScorePredicate = "D"
		result.DescriptionLow = "Perlu bimbingan dalam "
	}

	// Note: In real Merdeka, descriptions are dynamic based on WHICH TP has high/low scores.
	// This simplified engine assumes the caller will append specific TP names later
	// or we can extend the input to include TP-level scores.

	result.IsValid = true
	return result
}

func (v *MerdekaValidator) ValidateComponents(components map[string]float64, config model.GradingConfig) error {
	return nil
}
