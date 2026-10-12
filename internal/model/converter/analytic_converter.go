package converter

import (
	"fmt"
	"math"

	"github.com/paricorpdev/dashboard-pay2u-be/internal/model/payload"
	"github.com/paricorpdev/dashboard-pay2u-be/internal/repository"
)

func SummaryToResponse(month int, year int, res *repository.MonthlySummaryResult) *payload.AnalyticSummaryResponse {
	var successRate float64 = 0
	if res.TotalCount > 0 {
		successRate = math.Round((float64(res.SuccessCount)/float64(res.TotalCount))*10000) / 100
	}

	return &payload.AnalyticSummaryResponse{
		Month:                 month,
		Year:                  year,
		TotalVolumeAmount:     res.TotalVolume,
		TotalTransactionCount: res.TotalCount,
		SuccessCount:          res.SuccessCount,
		PendingCount:          res.PendingCount,
		ExpiredCount:          res.ExpiredCount,
		SuccessRate:           successRate,
		TotalAdminFee:         res.TotalAdminFee,
	}
}

func RawDailyTrendsToCompleteMonth(year int, month int, daysInMonth int, raw []repository.RawDailyTrend) []payload.DailyTrendItem {
	dataMap := make(map[string]repository.RawDailyTrend)
	for _, r := range raw {
		dataMap[r.DayStr] = r
	}

	result := make([]payload.DailyTrendItem, 0, daysInMonth)
	for day := 1; day <= daysInMonth; day++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		if item, exists := dataMap[dateStr]; exists {
			result = append(result, payload.DailyTrendItem{
				Date:         dateStr,
				Day:          day,
				TotalAmount:  item.TotalAmount,
				TotalCount:   item.TotalCount,
				SuccessCount: item.SuccessCount,
				PendingCount: item.PendingCount,
				ExpiredCount: item.ExpiredCount,
			})
		} else {
			result = append(result, payload.DailyTrendItem{
				Date:         dateStr,
				Day:          day,
				TotalAmount:  0,
				TotalCount:   0,
				SuccessCount: 0,
				PendingCount: 0,
				ExpiredCount: 0,
			})
		}
	}

	return result
}

func RawPaymentMethodToDistribution(raw []repository.RawPaymentMethodDist, totalCount int64) []payload.PaymentMethodDistributionItem {
	result := make([]payload.PaymentMethodDistributionItem, 0, len(raw))
	for _, item := range raw {
		var pct float64 = 0
		if totalCount > 0 {
			pct = math.Round((float64(item.Count)/float64(totalCount))*10000) / 100
		}
		result = append(result, payload.PaymentMethodDistributionItem{
			PaymentMethod: item.PaymentMethod,
			Count:         item.Count,
			Amount:        item.Amount,
			Percentage:    pct,
		})
	}
	return result
}

func RawStatusToDistribution(raw []repository.RawStatusDist, totalCount int64) []payload.StatusDistributionItem {
	result := make([]payload.StatusDistributionItem, 0, len(raw))
	for _, item := range raw {
		var pct float64 = 0
		if totalCount > 0 {
			pct = math.Round((float64(item.Count)/float64(totalCount))*10000) / 100
		}
		result = append(result, payload.StatusDistributionItem{
			Status:     item.Status,
			Count:      item.Count,
			Amount:     item.Amount,
			Percentage: pct,
		})
	}
	return result
}
