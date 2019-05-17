package date

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	formatyyyyMMdd   = "__06.01._2"
	formatyyyyMdd    = "__06.1._2"
	formatyyyyMMMdd  = "__06.Jan._2"
	formatyyyyMMMMdd = "__06.January._2"
	formatMMddyyyy   = "01._2.__06"
	formatMddyyyy    = "1._2.__06"
	formatMMMddyyyy  = "Jan._2.__06"
	formatMMMMddyyyy = "January._2.__06"
	formatddMMyyyy   = "_2.01.__06"
	formatddMyyyy    = "_2.1.__06"
	formatddMMMyyyy  = "_2.Jan.__06"
	formatddMMMMyyyy = "_2.January.__06"

	formatyyyyMM   = "__06.01"
	formatyyyyM    = "__06.1"
	formatyyyyMMM  = "__06.Jan"
	formatyyyyMMMM = "__06.January"
	formatMMyyyy   = "01.__06"
	formatMyyyy    = "1.__06"
	formatMMMyyyy  = "Jan.__06"
	formatMMMMyyyy = "January.__06"

	formatMMdd   = "01._2" // Always chosen over ddMM if both are possible.
	formatMdd    = "1._2"  // Always chosen over ddMM if both are possible.
	formatMMMdd  = "Jan._2"
	formatMMMMdd = "January._2"
	formatddMM   = "_2.01"
	formatddM    = "_2.1"
	formatddMMM  = "_2.Jan"
	formatddMMMM = "_2.January"

	formatHHmmss  = "15:04:05"
	formatHHmm    = "15:04"
	formathhmmssa = "03:04:05pm"
	formathhmma   = "03:04pm"
	formathmmssa  = "3:04:05pm"
	formathmma    = "3:04pm"

	formatzzz  = "MST"
	formatZZZ  = "-0700"
	formatZZZZ = "GMT-07:00"
	// TODO: Handle Z format timezones outside of standard.

	formatEEE  = "Mon"
	formatEEEE = "Monday"
)

var (
	separatorRegex       = regexp.MustCompile("[^a-zA-Z0-9]+")
	timezoneRegex        = regexp.MustCompile("[A-Z]{2,4}T")
	timezoneNumericRegex = regexp.MustCompile("((GMT)|Z)?[-+][0-9]{2}:?[0-9]{2}")

	// TODO: Nov and Mar may conflict with some timezones.
	months          = []string{"january", "februrary", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december"}
	monthsShort     = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
	daysOfWeekShort = []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
	daysOfWeek      = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	standardDateFormats = []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RubyDate,
		time.UnixDate,
		time.Kitchen,
		time.ANSIC,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
	dateFormatsWithShortMonth = []string{
		formatyyyyMMMdd,
		formatMMMddyyyy,
		formatddMMMyyyy,
		formatMMMdd,
		formatddMMM,
		formatMMMyyyy,
		formatyyyyMMM,
	}
	dateFormatsWithLongMonth = []string{
		formatyyyyMMMMdd,
		formatMMMMddyyyy,
		formatddMMMMyyyy,
		formatMMMMdd,
		formatddMMMM,
		formatMMMMyyyy,
		formatyyyyMMMM,
	}
	dateFormatsNumeric = []string{
		formatMddyyyy,
		formatMMddyyyy,
		formatddMMyyyy,
		formatddMyyyy,
		formatyyyyMMdd,
		formatyyyyMdd,
		formatMMdd,
		formatMdd,
		formatddMM,
		formatddM,
		formatMMyyyy,
		formatMyyyy,
		formatyyyyMM,
		formatyyyyM,
	}
	dateFormatsNumericYearFirst = []string{
		formatyyyyMM,
		formatyyyyM,
		formatyyyyMMdd,
		formatyyyyMdd,
		formatMddyyyy,
		formatMMddyyyy,
		formatddMMyyyy,
		formatddMyyyy,
		formatMMyyyy,
		formatMyyyy,
		formatMMdd,
		formatMdd,
		formatddMM,
		formatddM,
	}
	timeFormats = []string{
		formatHHmmss,
		formatHHmm,
		formathhmmssa,
		formathhmma,
		formathmmssa,
		formathmma,
	}
	timeZoneFormats = []string{
		formatzzz,
		formatZZZ,
		formatZZZZ,
	}
)

func init() {
	timezoneRegex.Longest()
	timezoneNumericRegex.Longest()
	separatorRegex.Longest()
}

// Parse attempts to parse this string as a timestamp, returning an error
// if it cannot. Example inputs: "July 9 1977", "07/9/1977 5:03pm".
// Assumes UTC if a timezone is not provided.
func Parse(s string) (time.Time, error) {
	d, _, err := ParseAndGetLayout(s)
	return d, err
}

// MustParse is like Parse except it panics if the string is not parseable.
func MustParse(s string) time.Time {
	d, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return d
}

// ParseInLocation is like Parse except it uses the given location when parsing the date.
func ParseInLocation(s string, loc *time.Location) (time.Time, error) {
	_, l, err := ParseAndGetLayout(s)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.ParseInLocation(l, s, loc)
	return t, err
}

// ParseAndGetLayout attempts to parse this string as a timestamp
// and if successful, returns the timestamp and the layout of the
// string.
func ParseAndGetLayout(date string) (time.Time, string, error) {
	if len(strings.TrimSpace(date)) == 0 {
		return time.Time{}, "", errors.New("Empty string cannot be parsed to date")
	}
	// Check standard date formats first.
	for _, f := range standardDateFormats {
		if t, err := time.Parse(f, date); err == nil {
			return t, f, nil
		}
	}
	s := strings.ToLower(date)
	layout := &bytes.Buffer{}
	prefix := getPrefix(s)
	layout.WriteString(prefix)
	s = strings.TrimPrefix(s, prefix)

	// Check for day of week.
	for _, d := range daysOfWeek {
		if strings.HasPrefix(s, d) {
			s = strings.TrimPrefix(s, d)
			layout.WriteString(formatEEEE)
		}
	}
	for _, d := range daysOfWeekShort {
		if strings.HasPrefix(s, d) {
			s = strings.TrimPrefix(s, d)
			layout.WriteString(formatEEE)
		}
	}

	// Get rid of prefix and suffix.
	prefix = getPrefix(s)
	separators := separatorRegex.FindAllStringSubmatch(s, -1)
	if len(prefix) > 0 {
		s = strings.TrimPrefix(s, prefix)
		layout.WriteString(prefix)
		separators = separators[1:]
	}
	var suffix string
	if len(separators) > 0 && strings.HasSuffix(s, separators[len(separators)-1][0]) {
		suffix = separators[len(separators)-1][0]
		s = strings.TrimSuffix(s, suffix)
		separators = separators[:len(separators)-1]
	}

	// Narrow down formats needed to check.
	// TODO: Make more efficient by checking fewer formats variations.
	var formats []string
	containsTime := containsTime(date)
	containsTimezone := containsTimezone(date)
	var onlyTime bool
	if containsLongMonth(s) {
		formats = dateFormatsWithLongMonth
	} else if containsShortMonth(s) {
		formats = dateFormatsWithShortMonth
	} else if (len(separators) <= 3 && containsTime && containsTimezone) || (len(separators) <= 2 && containsTime) {
		if containsTimezone {
			formats = getCombinations(timeFormats, timeZoneFormats, false)
		} else {
			formats = timeFormats
		}
		onlyTime = true
	} else if len(separators) == 0 {
		// If the date is all munged together, assume year is first rather than month.
		formats = dateFormatsNumericYearFirst
	} else {
		formats = dateFormatsNumeric
	}
	if containsTimezone {
		// time.Parse only accepts uppercase timezones names, so need to convert back.
		if tz := timezoneRegex.FindStringSubmatch(date); len(tz) > 0 {
			s = strings.Replace(s, strings.ToLower(tz[0]), tz[0], -1)
		}
	}

	// Check possible formats.
	for _, f := range formats {
		variations := getVariations(f, containsTime && !onlyTime, containsTimezone && !onlyTime)
		var correct string
		for _, v := range variations {
			if strings.Contains(v, ":") {
				v = strings.Replace(v, ":", ".", -1)
			}
			l := formatWithSeparators(v, separators)
			if _, err := time.Parse(l, s); err == nil {
				correct = l
				break
			}
		}
		if len(correct) > 0 {
			layout.WriteString(correct)
			break
		}
	}
	layout.WriteString(suffix)

	// Return Date and format.
	date = strings.Replace(date, "AM", "am", -1)
	date = strings.Replace(date, "PM", "pm", -1)
	t, err := time.Parse(layout.String(), date)
	if err != nil {
		return time.Time{}, "", err
	}
	return t, layout.String(), err
}

// Layout returns the layout of this date, appropriate for use
// with the Go time package, for example "Jan 02 2006".
// See http://golang.org/pkg/time/ for more examples.
func Layout(s string) string {
	_, l, err := ParseAndGetLayout(s)
	if err == nil {
		return l
	}
	return ""
}

// LayoutUnicode returns the layout of this date according to the
// Unicode standard: http://www.unicode.org/reports/tr35/tr35-19.html#Date_Format_Patterns
// NOT TESTED.
func LayoutUnicode(s string) string {
	_, l, err := ParseAndGetLayout(s)
	if err == nil {
		return ConvertGoLayoutToUnicode(l)
	}
	return ""
}

// ConvertGoLayoutToUnicode converts the given time layout string
// to on using the Unicode standard prescribed in
// http://www.unicode.org/reports/tr35/tr35-19.html#Date_Format_Patterns
// NOT TESTED.
func ConvertGoLayoutToUnicode(layout string) string {
	// Year
	layout = strings.Replace(layout, "20", "yy", -1)
	layout = strings.Replace(layout, "06", "yy", -1)

	// Month
	layout = strings.Replace(layout, "January", "MMMM", -1)
	layout = strings.Replace(layout, "Jan", "MMM", -1)
	layout = strings.Replace(layout, "01", "MM", -1)
	layout = strings.Replace(layout, "1", "M", -1)

	// Day
	layout = strings.Replace(layout, "02", "dd", -1)
	layout = strings.Replace(layout, "_2", "d", -1)

	// Weekday
	layout = strings.Replace(layout, "Mon", "EEE", -1)
	layout = strings.Replace(layout, "Monday", "EEEE", -1)

	// Hour
	layout = strings.Replace(layout, "03", "hh", -1)
	layout = strings.Replace(layout, "3", "h", -1)
	layout = strings.Replace(layout, "15", "HH", -1)
	layout = strings.Replace(layout, "PM", "a", -1)

	// Minute
	layout = strings.Replace(layout, "04", "mm", -1)

	// Second
	layout = strings.Replace(layout, "05", "ss", -1)

	// Timezone
	layout = strings.Replace(layout, "MST", "zzz", -1)
	layout = strings.Replace(layout, "-0700", "ZZZ", -1)
	layout = strings.Replace(layout, "GMT-07:00", "ZZZZ", -1)

	return layout
}

// Private Methods

func getVariations(f string, includeTime, includeTimezone bool) []string {
	var v []string
	if strings.Contains(f, "__") {
		v = []string{strings.Replace(f, "__", "20", 1), strings.Replace(f, "__", "", 1)}
	} else {
		v = []string{f}
	}
	l := len(v)
	for i := 0; i < l; i++ {
		if strings.Contains(v[i], "_") {
			v = append(v, strings.Replace(v[i], "_", "0", -1))
		}
	}
	if includeTime {
		if includeTimezone {
			times := getCombinations(timeFormats, timeZoneFormats, false)
			v = getCombinations(v, times, true)
		} else {
			v = getCombinations(v, timeFormats, true)
		}
	}
	return v
}

func getCombinations(a, b []string, switchOrder bool) []string {
	var res []string
	for _, s := range a {
		for _, s2 := range b {
			res = append(res, fmt.Sprintf("%s.%s", s, s2))
			if switchOrder {
				res = append(res, fmt.Sprintf("%s.%s", s2, s))
			}
		}
	}
	return res
}

func getPrefix(s string) string {
	firstMatch := separatorRegex.FindStringSubmatchIndex(s)
	if len(firstMatch) == 0 {
		return ""
	}
	if firstMatch[0] == 0 {
		return s[:firstMatch[1]]
	}
	return ""
}

func containsShortMonth(s string) bool {
	for _, m := range monthsShort {
		if strings.Contains(s, m) {
			return true
		}
	}
	return false
}

func containsLongMonth(s string) bool {
	for _, m := range months {
		if strings.Contains(s, m) {
			return true
		}
	}
	return false
}

func containsTime(s string) bool {
	return strings.Contains(s, ":")
}

func containsTimezone(s string) bool {
	return (timezoneRegex.MatchString(s) || timezoneNumericRegex.MatchString(s)) &&
		!(strings.Contains(s, "SEPT") || strings.Contains(s, "OCT") || strings.Contains(s, "SAT"))
}

func formatWithSeparators(f string, sep [][]string) string {
	if len(sep) == 0 {
		return strings.Replace(f, ".", "", -1)
	}
	for i := 0; i < len(sep); i++ {
		s := sep[i][0]
		if s == "." {
			s = "．"
		}
		f = strings.Replace(f, ".", s, 1)
	}
	return strings.Replace(f, "．", ".", -1)
}
