package model

// TierFeatures defines which features are available for each tier
var TierFeatures = map[string]map[string]bool{
	TierBasic: {
		"payment_gateway":   false, // Manual payment only
		"wa_auto_notif":     false, // No auto WhatsApp
		"api_access":        false, // No API access
		"multi_bank":        false, // Single bank only
		"bulk_spp_generate": false, // No bulk SPP generation
	},
	TierPremium: {
		"payment_gateway":   true, // Moota/Tripay integration
		"wa_auto_notif":     true, // Auto WhatsApp notifications
		"api_access":        true, // API access for integrations
		"multi_bank":        true, // Multiple bank accounts
		"bulk_spp_generate": true, // Bulk SPP generation
	},
}

// Feature constants for easy reference
const (
	FeaturePaymentGateway  = "payment_gateway"
	FeatureWAAutoNotif     = "wa_auto_notif"
	FeatureAPIAccess       = "api_access"
	FeatureMultiBank       = "multi_bank"
	FeatureBulkSPPGenerate = "bulk_spp_generate"
)

// HasFeature checks if a given tier has access to a specific feature
func HasFeature(tier, feature string) bool {
	if tier == "" {
		tier = TierBasic // Default to basic
	}
	if tierFeatures, ok := TierFeatures[tier]; ok {
		return tierFeatures[feature]
	}
	return false
}

// GetTierDisplayName returns Indonesian display name for tier
func GetTierDisplayName(tier string) string {
	switch tier {
	case TierPremium:
		return "Premium"
	default:
		return "Basic"
	}
}
