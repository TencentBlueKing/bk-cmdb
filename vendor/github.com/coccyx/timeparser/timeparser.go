// Package timeparser implements Splunk's relative time syntax, or attempts to guess as
// best as possible if not relative time.
package timeparser

import (
	"fmt"
	godate "github.com/joyt/godate"
	"math"
	"regexp"
	"strconv"
	"time"
)

var (
	unitsre   string         = "(seconds|second|secs|sec|minutes|minute|min|hours|hour|hrs|hr|days|day|weeks|week|w[0-6]|months|month|mon|quarters|quarter|qtrs|qtr|years|year|yrs|yr|s|h|m|d|w|y|w|q)"
	reltimere string         = "(?i)(?P<plusminus>[+-]*)(?P<num>\\d{1,})(?P<unit>" + unitsre + "{1})(([\\@](?P<snapunit>" + unitsre + "{1})((?P<snapplusminus>[+-])(?P<snaprelnum>\\d+)(?P<snaprelunit>" + unitsre + "{1}))*)*)"
	re        *regexp.Regexp = regexp.MustCompile(reltimere)
	now                      = time.Now
	loc       *time.Location

	secs  = map[string]bool{"s": true, "sec": true, "secs": true, "second": true, "seconds": true}
	mins  = map[string]bool{"m": true, "min": true, "minute": true, "minutes": true}
	hours = map[string]bool{"h": true, "hr": true, "hrs": true, "hour": true, "hours": true}
	days  = map[string]bool{"d": true, "day": true, "days": true}
	weeks = map[string]bool{"w": true, "week": true, "weeks": true,
		"w0": true, "w1": true, "w2": true, "w3": true, "w4": true, "w5": true, "w6": true}
	months   = map[string]bool{"mon": true, "month": true, "months": true}
	quarters = map[string]bool{"q": true, "qtr": true, "qtrs": true, "quarter": true, "quarters": true}
	years    = map[string]bool{"y": true, "yr": true, "yrs": true, "year": true, "years": true}
)

func init() {
	loc, _ = time.LoadLocation("Local")
}

// TimeParser returns a parsed time based on the current time in the local time zone
func TimeParser(ts string) (time.Time, error) {
	return TimeParserNowInLocation(ts, now, loc)
}

// TimeParser returns a parsed time based on now returned from the passed callback in the local time zone
func TimeParserNow(ts string, now func() time.Time) (time.Time, error) {
	return TimeParserNowInLocation(ts, now, loc)
}

// TimeParser returns a parsed time based on the current time in the passed time zone
func TimeParserInLocation(ts string, loc *time.Location) (time.Time, error) {
	return TimeParserNowInLocation(ts, now, loc)
}

// TimeParser returns a parsed time based on now returned from the passed callback in the passed time zone
func TimeParserNowInLocation(ts string, now func() time.Time, loc *time.Location) (time.Time, error) {
	if ts == "now" {
		return now(), nil
	} else {
		if ts[:1] == "+" || ts[:1] == "-" {
			ret := now()

			match := re.FindStringSubmatch(ts)
			results := make(map[string]string)
			for i, name := range re.SubexpNames() {
				if i != 0 {
					results[name] = match[i]
				}
			}

			// Handle first part of the time string
			if len(results["plusminus"]) != 0 && len(results["num"]) != 0 && len(results["unit"]) != 0 {
				timeParserTimeMath(results["plusminus"], results["num"], results["unit"], &ret)

				snapunit := results["snapunit"]
				if len(snapunit) > 0 {
					switch {
					case secs[snapunit]:
						ret = time.Date(ret.Year(), ret.Month(), ret.Day(), ret.Hour(), ret.Minute(), ret.Second(), 0, loc)
					case mins[snapunit]:
						ret = time.Date(ret.Year(), ret.Month(), ret.Day(), ret.Hour(), ret.Minute(), 0, 0, loc)
					case hours[snapunit]:
						ret = time.Date(ret.Year(), ret.Month(), ret.Day(), ret.Hour(), 0, 0, 0, loc)
					case days[snapunit]:
						ret = time.Date(ret.Year(), ret.Month(), ret.Day(), 0, 0, 0, 0, loc)
					case weeks[snapunit]:
						// Only w[0-6] have length of 2
						if len(snapunit) == 2 {
							wdnum, err := strconv.Atoi(snapunit[1:2])
							if err != nil {
								return ret, err
							}
							wd := int(ret.Weekday())

							if wdnum <= wd {
								ret = ret.Add(time.Duration(24*(wdnum-wd)) * time.Hour)
								ret = time.Date(ret.Year(), ret.Month(), ret.Day(), 0, 0, 0, 0, loc)
							} else {
								ret = ret.Add(time.Duration(24*7*-1) * time.Hour)
								ret = ret.Add(time.Duration(24*-1*wd) * time.Hour)
								ret = ret.Add(time.Duration(24*wdnum) * time.Hour)
								ret = time.Date(ret.Year(), ret.Month(), ret.Day(), 0, 0, 0, 0, loc)
							}
						}
					case months[snapunit]:
						ret = time.Date(ret.Year(), ret.Month(), 1, 0, 0, 0, 0, loc)
					case quarters[snapunit]:
						tmonth := int(math.Floor(float64(ret.Month()/3)) * 3)
						ret = time.Date(ret.Year(), time.Month(tmonth), 1, 0, 0, 0, 0, loc)
					case years[snapunit]:
						ret = time.Date(ret.Year(), 1, 1, 0, 0, 0, 0, loc)
					}

					if len(results["snapplusminus"]) != 0 && len(results["snaprelnum"]) != 0 && len(results["snaprelunit"]) != 0 {
						timeParserTimeMath(results["snapplusminus"], results["snaprelnum"], results["snaprelunit"], &ret)
					}
				}
				return ret, nil
			}
		} else { // We're not a relative time, so try our best to interpret the date passed
			return godate.ParseInLocation(ts, loc)
		}
	}
	return now(), fmt.Errorf("Got to the end but didn't return")
}

func timeParserTimeMath(plusminus string, numstr string, unit string, ret *time.Time) {
	num, _ := strconv.Atoi(numstr)
	if plusminus == "-" {
		num *= -1
	}

	switch {
	case secs[unit]:
		*ret = ret.Add(time.Duration(num) * time.Second)
	case mins[unit]:
		*ret = ret.Add(time.Duration(num) * time.Minute)
	case hours[unit]:
		*ret = ret.Add(time.Duration(num) * time.Hour)
	case days[unit]:
		*ret = ret.AddDate(0, 0, num)
	case weeks[unit]:
		*ret = ret.AddDate(0, 0, num*7)
	case months[unit]:
		*ret = ret.AddDate(0, num, 0)
	case quarters[unit]:
		*ret = ret.AddDate(0, num*3, 0)
	case years[unit]:
		*ret = ret.AddDate(num, 0, 0)
	}
}
