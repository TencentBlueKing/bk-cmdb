package timeparser

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeParser(t *testing.T) {
	_ = fmt.Sprintf("foo") // Just to avoid annoying warnings

	loc, _ := time.LoadLocation("Local")
	n := time.Date(2001, 10, 20, 12, 0, 0, 100000, loc)

	now := func() time.Time {
		return n
	}

	// Test Now
	tn, _ := TimeParserNow("now", now)
	assert.Equal(t, n, tn)

	// Test -1h
	x := n.Add(time.Duration(1) * time.Hour * -1)
	tn, _ = TimeParserNow("-1h", now)
	assert.Equal(t, x, tn)

	// Test -10s
	x = n.Add(time.Duration(10) * time.Second * -1)
	tn, _ = TimeParserNow("-10s", now)
	assert.Equal(t, x, tn)

	// Test -59m
	x = n.Add(time.Duration(59) * time.Minute * -1)
	tn, _ = TimeParserNow("-59m", now)
	assert.Equal(t, x, tn)

	// Test +1day
	x = n.AddDate(0, 0, 1)
	tn, _ = TimeParserNow("+1day", now)
	assert.Equal(t, x, tn)

	// Test Snapto sec
	x = n.Add(time.Duration(1) * time.Second * -1)
	x = time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), 0, loc)
	tn, _ = TimeParserNow("-1s@s", now)
	assert.Equal(t, x, tn)

	// Test Snapto min
	x = n.Add(time.Duration(1) * time.Minute * -1)
	x = time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), 0, loc)
	tn, _ = TimeParserNow("-1s@m", now)
	assert.Equal(t, x, tn)

	// Test Snapto hour
	x = n.Add(time.Duration(1) * time.Hour * -1)
	x = time.Date(x.Year(), x.Month(), x.Day(), x.Hour(), x.Minute(), x.Second(), 0, loc)
	tn, _ = TimeParserNow("-1s@h", now)
	assert.Equal(t, x, tn)

	// Test Snapto day
	x = time.Date(2001, 10, 20, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@d", now)
	assert.Equal(t, x, tn)

	// Test Snapto month
	x = time.Date(2001, 10, 1, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@mon", now)
	assert.Equal(t, x, tn)

	// Test Snapto quarter
	x = time.Date(2001, 9, 1, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@qtr", now)
	assert.Equal(t, x, tn)

	// Test Snapto year
	x = time.Date(2001, 1, 1, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@year", now)
	assert.Equal(t, x, tn)

	// Test Snapto year with math
	x = time.Date(2001, 2, 1, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@year+1mon", now)
	assert.Equal(t, x, tn)

	// Test Snapto weekday
	x = time.Date(2001, 10, 14, 0, 0, 0, 0, loc)
	tn, _ = TimeParserNow("-1h@w0", now)
	assert.Equal(t, x, tn)
}
