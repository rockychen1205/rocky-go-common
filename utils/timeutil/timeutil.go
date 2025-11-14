// Package timeutil provides time-related utility functions
package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseTimeRange parses time range strings like "14-15", "14:30-15:45"
func ParseTimeRange(timeRange string) (time.Time, time.Time, error) {
	parts := strings.Split(timeRange, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid time range format: %s", timeRange)
	}
	
	start, err := parseTime(parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time: %w", err)
	}
	
	end, err := parseTime(parts[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time: %w", err)
	}
	
	// If end time is before start time, assume it's next day
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}
	
	return start, end, nil
}

// parseTime parses time string like "14", "14:30", "14:30:45"
func parseTime(timeStr string) (time.Time, error) {
	now := time.Now()
	timeStr = strings.TrimSpace(timeStr)
	
	parts := strings.Split(timeStr, ":")
	
	switch len(parts) {
	case 1:
		// Hour only: "14"
		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %s", parts[0])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location()), nil
		
	case 2:
		// Hour and minute: "14:30"
		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %s", parts[0])
		}
		minute, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid minute: %s", parts[1])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location()), nil
		
	case 3:
		// Hour, minute and second: "14:30:45"
		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %s", parts[0])
		}
		minute, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid minute: %s", parts[1])
		}
		second, err := strconv.Atoi(parts[2])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid second: %s", parts[2])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location()), nil
		
	default:
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}
}

// FormatDuration formats a duration to human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.0fms", float64(d)/float64(time.Millisecond))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// ParseDate parses date string in various formats
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
		"20060102",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// IsToday checks if a time is today
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday checks if a time is yesterday
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// StartOfDay returns the start of day (00:00:00) for the given time
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of day (23:59:59) for the given time
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}
