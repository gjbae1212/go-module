package util

import "time"

var (
	maxTime       = time.Unix(0, (1<<63)-1)
	utcLayout     = "2006-01-02 15:04:05"
	yearlyLayout  = "2006"
	monthlyLayout = "200601"
	dailyLayout   = "20060102"
	hourlyLayout  = "2006010215"
)

func TimestampByMaxTime() int64 {
	return maxTime.UnixNano()
}

func StringToTime(s string) time.Time {
	t, err := time.Parse(utcLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToString(t time.Time) string {
	return t.Format(utcLayout)
}

func YearlyStringToTime(s string) time.Time {
	t, err := time.Parse(yearlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToYearlyStringFormat(t time.Time) string {
	return t.Format(yearlyLayout)
}

func MonthlyStringToTime(s string) time.Time {
	t, err := time.Parse(monthlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToMonthlyStringFormat(t time.Time) string {
	return t.Format(monthlyLayout)
}

func DailyStringToTime(s string) time.Time {
	t, err := time.Parse(dailyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToDailyStringFormat(t time.Time) string {
	return t.Format(dailyLayout)
}

func HourlyStringToTime(s string) time.Time {
	t, err := time.Parse(hourlyLayout, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func TimeToHourlyStringFormat(t time.Time) string {
	return t.Format(hourlyLayout)
}
