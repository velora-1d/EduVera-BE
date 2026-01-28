package dashboard

import (
	"context"
	"prabogo/internal/model"
)

type mockService struct{}

func NewMockService() Service {
	return &mockService{}
}

func (s *mockService) GetStats(ctx context.Context, tenantID string) (*model.DashboardStats, error) {
	// Mock data based on typical pesantren values for demo
	return &model.DashboardStats{
		TotalSantri:      156,
		TotalAsrama:      4,
		TotalUstadz:      12,
		TotalPengurus:    8,
		AttendanceRate:   95.5,
		ActiveViolations: 3,
		CashBalance:      45000000,
		IncomeMonth:      15000000,
		ExpenseMonth:     8500000,
	}, nil
}
