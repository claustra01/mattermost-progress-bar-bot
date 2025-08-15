package date_test

import (
	"testing"
	"time"

	"github.com/claustra01/typetalk-progress-bar-bot/pkg/date"
)

func TestGetRemainingDays(t *testing.T) {
	tests := []struct {
		now    time.Time
		expect int
	}{
		{time.Date(2024, 7, 22, 6, 0, 0, 0, time.Now().Location()), 230},
		{time.Date(2025, 3, 8, 6, 0, 0, 0, time.Now().Location()), 1},
		{time.Date(2025, 3, 9, 6, 0, 0, 0, time.Now().Location()), 0},
	}

	for _, tt := range tests {
		result := date.GetRemainingDays(tt.now)
		if result != tt.expect {
			t.Errorf("GetRemainingDays(%v) = %v; want %v", tt.now, result, tt.expect)
		}
	}
}
