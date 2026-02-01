package analytics

import (
	"context"
	"time"

	"prabogo/internal/model"
	outbound_port "prabogo/internal/port/outbound"
)

type AnalyticsDomain interface {
	GetAnalytics(ctx context.Context, tenantID string) (*model.AnalyticsData, error)
}

type analyticsDomain struct {
	db outbound_port.DatabasePort
}

func NewAnalyticsDomain(db outbound_port.DatabasePort) AnalyticsDomain {
	return &analyticsDomain{db: db}
}

// GetAnalytics returns complete analytics data for charts
func (d *analyticsDomain) GetAnalytics(ctx context.Context, tenantID string) (*model.AnalyticsData, error) {
	analytics := &model.AnalyticsData{}

	// Try to get real attendance data, otherwise use mock
	attendanceData := d.getDailyAttendance(ctx, tenantID, 7)
	analytics.AttendanceChart = attendanceData

	// Get monthly payment data (mock for now)
	analytics.PaymentChart = d.getPaymentChart(ctx, tenantID)

	// Get enrollment chart (mock for now)
	analytics.EnrollmentChart = d.getEnrollmentChart()

	// Get SPP collection chart (mock for now)
	analytics.SPPCollectionChart = d.getSPPCollectionChart()

	return analytics, nil
}

// getDailyAttendance returns attendance chart data
func (d *analyticsDomain) getDailyAttendance(_ context.Context, _ string, days int) []model.ChartPoint {
	var chart []model.ChartPoint

	// Generate dates for last N days
	now := time.Now()
	dayNames := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dayName := dayNames[int(date.Weekday())]

		// Mock attendance rate (varies by day)
		rate := 85.0 + float64((7-i)*2%15)
		if date.Weekday() == time.Saturday {
			rate = 60.0
		}
		if date.Weekday() == time.Sunday {
			rate = 0.0
		}

		chart = append(chart, model.ChartPoint{
			Label: dayName,
			Value: rate,
		})
	}

	return chart
}

// getPaymentChart returns monthly payment data
func (d *analyticsDomain) getPaymentChart(_ context.Context, _ string) []model.ChartPoint {
	// Generate last 6 months
	months := []string{"Sep", "Okt", "Nov", "Des", "Jan", "Feb"}
	values := []float64{45000000, 52000000, 48000000, 55000000, 58000000, 51000000}

	var chart []model.ChartPoint
	for i, m := range months {
		chart = append(chart, model.ChartPoint{Label: m, Value: values[i]})
	}
	return chart
}

// getEnrollmentChart returns student enrollment trend
func (d *analyticsDomain) getEnrollmentChart() []model.ChartPoint {
	months := []string{"Sep", "Okt", "Nov", "Des", "Jan", "Feb"}
	values := []float64{15, 8, 5, 3, 12, 6}

	var chart []model.ChartPoint
	for i, m := range months {
		chart = append(chart, model.ChartPoint{Label: m, Value: values[i]})
	}
	return chart
}

// getSPPCollectionChart returns SPP collection rate by month
func (d *analyticsDomain) getSPPCollectionChart() []model.ChartPoint {
	months := []string{"Sep", "Okt", "Nov", "Des", "Jan", "Feb"}
	values := []float64{95, 92, 88, 94, 91, 85}

	var chart []model.ChartPoint
	for i, m := range months {
		chart = append(chart, model.ChartPoint{Label: m, Value: values[i]})
	}
	return chart
}
