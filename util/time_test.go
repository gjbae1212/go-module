package util

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	assert := assert.New(t)

	tt := time.Now()
	ss := TimeToString(tt)
	compare1 := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", tt.Year(), tt.Month(), tt.Day(), tt.Hour(), tt.Minute(), tt.Second())
	assert.Equal(compare1, ss)

	tt2 := StringToTime(ss)
	assert.Equal(tt.Year(), tt2.Year())
	assert.Equal(tt.Month(), tt2.Month())
	assert.Equal(tt.Day(), tt2.Day())
	assert.Equal(tt.Hour(), tt2.Hour())
	assert.Equal(tt.Minute(), tt2.Minute())
	assert.Equal(tt.Second(), tt2.Second())

	ss2 := TimeToDailyStringFormat(tt)
	compare2 := fmt.Sprintf("%d%02d%02d", tt.Year(), tt.Month(), tt.Day())
	assert.Equal(compare2, ss2)

	tt3 := DailyStringToTime(ss2)
	assert.Equal(tt.Year(), tt3.Year())
	assert.Equal(tt.Month(), tt3.Month())
	assert.Equal(tt.Day(), tt3.Day())
	assert.Equal(0, tt3.Hour())
	assert.Equal(0, tt3.Minute())
	assert.Equal(0, tt3.Second())

	ss3 := TimeToHourlyStringFormat(tt)
	compare3 := fmt.Sprintf("%d%02d%02d%02d", tt.Year(), tt.Month(), tt.Day(), tt.Hour())
	assert.Equal(compare3, ss3)

	tt4 := HourlyStringToTime(ss3)
	assert.Equal(tt.Year(), tt4.Year())
	assert.Equal(tt.Month(), tt4.Month())
	assert.Equal(tt.Day(), tt4.Day())
	assert.Equal(tt.Hour(), tt4.Hour())
	assert.Equal(0, tt4.Minute())
	assert.Equal(0, tt4.Second())
}
