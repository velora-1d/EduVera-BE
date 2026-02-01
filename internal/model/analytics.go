package model

import "time"

// AnalyticsData holds chart data for dashboard
type AnalyticsData struct {
	// Attendance over time (last 7 days)
	AttendanceChart []ChartPoint `json:"attendance_chart"`

	// Payment over time (last 6 months)
	PaymentChart []ChartPoint `json:"payment_chart"`

	// Student enrollment trend (last 6 months)
	EnrollmentChart []ChartPoint `json:"enrollment_chart"`

	// SPP collection rate by month
	SPPCollectionChart []ChartPoint `json:"spp_collection_chart"`
}

// ChartPoint represents a single data point in a chart
type ChartPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// AnalyticsFilter for querying analytics data
type AnalyticsFilter struct {
	TenantID  string    `json:"tenant_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Period    string    `json:"period"` // daily, weekly, monthly
}

// MonthlyStats holds monthly aggregated stats
type MonthlyStats struct {
	Month          string  `json:"month"`
	Year           int     `json:"year"`
	TotalIncome    float64 `json:"total_income"`
	TotalExpense   float64 `json:"total_expense"`
	NewStudents    int     `json:"new_students"`
	AttendanceRate float64 `json:"attendance_rate"`
	SPPTotal       float64 `json:"spp_total"`
	SPPCollected   float64 `json:"spp_collected"`
	CollectionRate float64 `json:"collection_rate"`
}

// DailyAttendance holds daily attendance data
type DailyAttendance struct {
	Date       string  `json:"date"`
	Present    int     `json:"present"`
	Absent     int     `json:"absent"`
	Permission int     `json:"permission"`
	Sick       int     `json:"sick"`
	Total      int     `json:"total"`
	Rate       float64 `json:"rate"`
}
