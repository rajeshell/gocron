package gocron

import (
	"testing"
	"time"
	"fmt"
)

func TestMonthlyScheduling(t *testing.T) {
	// Set default timezone to UTC for tests
	loc := time.UTC
	ChangeLoc(loc)

	tests := []struct {
		name           string
		job            *Job
		expectedNextRun time.Time
		expectedError  error
	}{
		{
			name: "Run on specific day of month",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.DayOfMonth(20)
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 20, 0, 0, 0, 0, loc),
		},
		{
			name: "Run on first day of month",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.FirstDayOfMonth()
				return j
			}(),
			expectedNextRun: time.Date(2024, 4, 1, 0, 0, 0, 0, loc),
		},
		{
			name: "Run on last day of month",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.LastDayOfMonth()
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 31, 0, 0, 0, 0, loc),
		},
		{
			name: "Run on specific day with time",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.DayOfMonth(20)
				j.At("15:30")
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 20, 15, 30, 0, 0, loc),
		},
		{
			name: "Invalid day of month",
			job: func() *Job {
				j := NewJob(1)
				j.Loc(loc)
				j.Month() // Explicitly set the unit to months
				j.DayOfMonth(32)
				return j
			}(),
			expectedError: fmt.Errorf("invalid day of month: 32, must be between 1 and 31"),
		},
		{
			name: "Handle February last day",
			job: func() *Job {
				now := time.Date(2024, 1, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.LastDayOfMonth()
				return j
			}(),
			expectedNextRun: time.Date(2024, 1, 31, 0, 0, 0, 0, loc),
		},
		{
			name: "Handle February last day in leap year",
			job: func() *Job {
				now := time.Date(2024, 2, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.LastDayOfMonth()
				return j
			}(),
			expectedNextRun: time.Date(2024, 2, 29, 0, 0, 0, 0, loc),
		},
		{
			name: "Handle day 31 in months with 30 days",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.DayOfMonth(31)
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 31, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.scheduleNextRun()
			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tt.job.nextRun.Equal(tt.expectedNextRun) {
				t.Errorf("expected next run %v, got %v", tt.expectedNextRun, tt.job.nextRun)
			}
		})
	}
}

func TestMonthlySchedulingWithTimezone(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	ChangeLoc(loc)

	tests := []struct {
		name           string
		job            *Job
		expectedNextRun time.Time
	}{
		{
			name: "Run on specific day with timezone",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.DayOfMonth(20)
				j.At("15:30")
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 20, 15, 30, 0, 0, loc),
		},
		{
			name: "Run on last day with timezone",
			job: func() *Job {
				now := time.Date(2024, 3, 15, 10, 0, 0, 0, loc)
				j := NewJob(1)
				j.lastRun = now
				j.Loc(loc)
				j.LastDayOfMonth()
				j.At("15:30")
				return j
			}(),
			expectedNextRun: time.Date(2024, 3, 31, 15, 30, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.job.scheduleNextRun()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if !tt.job.nextRun.Equal(tt.expectedNextRun) {
				t.Errorf("expected next run %v, got %v", tt.expectedNextRun, tt.job.nextRun)
			}
		})
	}
} 