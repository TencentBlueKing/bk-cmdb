// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// DateHistogramAggregation is a multi-bucket aggregation similar to the
// histogram except it can only be applied on date values.
// See: https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-aggregations-bucket-datehistogram-aggregation.html
type DateHistogramAggregation struct {
	field           string
	script          *Script
	missing         interface{}
	subAggregations map[string]Aggregation
	meta            map[string]interface{}

	interval          string
	fixedInterval     string
	calendarInterval  string
	order             string
	orderAsc          bool
	minDocCount       *int64
	extendedBoundsMin interface{}
	extendedBoundsMax interface{}
	timeZone          string
	format            string
	offset            string
	keyed             *bool
}

// NewDateHistogramAggregation creates a new DateHistogramAggregation.
func NewDateHistogramAggregation() *DateHistogramAggregation {
	return &DateHistogramAggregation{
		subAggregations: make(map[string]Aggregation),
	}
}

// Field on which the aggregation is processed.
func (a *DateHistogramAggregation) Field(field string) *DateHistogramAggregation {
	a.field = field
	return a
}

func (a *DateHistogramAggregation) Script(script *Script) *DateHistogramAggregation {
	a.script = script
	return a
}

// Missing configures the value to use when documents miss a value.
func (a *DateHistogramAggregation) Missing(missing interface{}) *DateHistogramAggregation {
	a.missing = missing
	return a
}

func (a *DateHistogramAggregation) SubAggregation(name string, subAggregation Aggregation) *DateHistogramAggregation {
	a.subAggregations[name] = subAggregation
	return a
}

// Meta sets the meta data to be included in the aggregation response.
func (a *DateHistogramAggregation) Meta(metaData map[string]interface{}) *DateHistogramAggregation {
	a.meta = metaData
	return a
}

// Interval by which the aggregation gets processed. This field
// will be replaced by the two FixedInterval and CalendarInterval
// fields (see below).
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.4/search-aggregations-bucket-datehistogram-aggregation.html
//
// Deprecated: This field will be removed in the future.
func (a *DateHistogramAggregation) Interval(interval string) *DateHistogramAggregation {
	a.interval = interval
	return a
}

// FixedInterval by which the aggregation gets processed.
//
// Allowed values are: "year", "1y", "quarter", "1q", "month", "1M",
// "week", "1w", "day", "1d", "hour", "1h", "minute", "1m", "second",
// or "1s". It also supports time settings like "1.5h".
//
// These units are not calendar-aware and are simply multiples of
// fixed, SI units. This is mutually exclusive with CalendarInterval.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.4/search-aggregations-bucket-datehistogram-aggregation.html
func (a *DateHistogramAggregation) FixedInterval(fixedInterval string) *DateHistogramAggregation {
	a.fixedInterval = fixedInterval
	return a
}

// CalendarInterval by which the aggregation gets processed.
//
// Allowed values are: "year" ("1y", "y"), "quarter" ("1q", "q"),
// "month" ("1M", "M"), "week" ("1w", "w"), "day" ("d", "1d")
//
// These units are calendar-aware, meaning they respect leap
// additions, variable days per month etc. This is mutually
// exclusive with FixedInterval.
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.4/search-aggregations-bucket-datehistogram-aggregation.html
func (a *DateHistogramAggregation) CalendarInterval(calendarInterval string) *DateHistogramAggregation {
	a.calendarInterval = calendarInterval
	return a
}

// Order specifies the sort order. Valid values for order are:
// "_key", "_count", a sub-aggregation name, or a sub-aggregation name
// with a metric.
func (a *DateHistogramAggregation) Order(order string, asc bool) *DateHistogramAggregation {
	a.order = order
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByCount(asc bool) *DateHistogramAggregation {
	// "order" : { "_count" : "asc" }
	a.order = "_count"
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByCountAsc() *DateHistogramAggregation {
	return a.OrderByCount(true)
}

func (a *DateHistogramAggregation) OrderByCountDesc() *DateHistogramAggregation {
	return a.OrderByCount(false)
}

func (a *DateHistogramAggregation) OrderByKey(asc bool) *DateHistogramAggregation {
	// "order" : { "_key" : "asc" }
	a.order = "_key"
	a.orderAsc = asc
	return a
}

func (a *DateHistogramAggregation) OrderByKeyAsc() *DateHistogramAggregation {
	return a.OrderByKey(true)
}

func (a *DateHistogramAggregation) OrderByKeyDesc() *DateHistogramAggregation {
	return a.OrderByKey(false)
}

// OrderByAggregation creates a bucket ordering strategy which sorts buckets
// based on a single-valued calc get.
func (a *DateHistogramAggregation) OrderByAggregation(aggName string, asc bool) *DateHistogramAggregation {
	// {
	//     "aggs" : {
	//         "genders" : {
	//             "terms" : {
	//                 "field" : "gender",
	//                 "order" : { "avg_height" : "desc" }
	//             },
	//             "aggs" : {
	//                 "avg_height" : { "avg" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = aggName
	a.orderAsc = asc
	return a
}

// OrderByAggregationAndMetric creates a bucket ordering strategy which
// sorts buckets based on a multi-valued calc get.
func (a *DateHistogramAggregation) OrderByAggregationAndMetric(aggName, metric string, asc bool) *DateHistogramAggregation {
	// {
	//     "aggs" : {
	//         "genders" : {
	//             "terms" : {
	//                 "field" : "gender",
	//                 "order" : { "height_stats.avg" : "desc" }
	//             },
	//             "aggs" : {
	//                 "height_stats" : { "stats" : { "field" : "height" } }
	//             }
	//         }
	//     }
	// }
	a.order = aggName + "." + metric
	a.orderAsc = asc
	return a
}

// MinDocCount sets the minimum document count per bucket.
// Buckets with less documents than this min value will not be returned.
func (a *DateHistogramAggregation) MinDocCount(minDocCount int64) *DateHistogramAggregation {
	a.minDocCount = &minDocCount
	return a
}

// TimeZone sets the timezone in which to translate dates before computing buckets.
func (a *DateHistogramAggregation) TimeZone(timeZone string) *DateHistogramAggregation {
	a.timeZone = timeZone
	return a
}

// Format sets the format to use for dates.
func (a *DateHistogramAggregation) Format(format string) *DateHistogramAggregation {
	a.format = format
	return a
}

// Offset sets the offset of time intervals in the histogram, e.g. "+6h".
func (a *DateHistogramAggregation) Offset(offset string) *DateHistogramAggregation {
	a.offset = offset
	return a
}

// ExtendedBounds accepts int, int64, string, or time.Time values.
// In case the lower value in the histogram would be greater than min or the
// upper value would be less than max, empty buckets will be generated.
func (a *DateHistogramAggregation) ExtendedBounds(min, max interface{}) *DateHistogramAggregation {
	a.extendedBoundsMin = min
	a.extendedBoundsMax = max
	return a
}

// ExtendedBoundsMin accepts int, int64, string, or time.Time values.
func (a *DateHistogramAggregation) ExtendedBoundsMin(min interface{}) *DateHistogramAggregation {
	a.extendedBoundsMin = min
	return a
}

// ExtendedBoundsMax accepts int, int64, string, or time.Time values.
func (a *DateHistogramAggregation) ExtendedBoundsMax(max interface{}) *DateHistogramAggregation {
	a.extendedBoundsMax = max
	return a
}

// Keyed specifies whether to return the results with a keyed response (or not).
//
// See https://www.elastic.co/guide/en/elasticsearch/reference/7.0/search-aggregations-bucket-datehistogram-aggregation.html#_keyed_response_3.
func (a *DateHistogramAggregation) Keyed(keyed bool) *DateHistogramAggregation {
	a.keyed = &keyed
	return a
}

func (a *DateHistogramAggregation) Source() (interface{}, error) {
	// Example:
	// {
	//     "aggs" : {
	//         "articles_over_time" : {
	//             "date_histogram" : {
	//                 "field" : "date",
	//                 "fixed_interval" : "month"
	//             }
	//         }
	//     }
	// }
	//
	// This method returns only the { "date_histogram" : { ... } } part.

	source := make(map[string]interface{})
	opts := make(map[string]interface{})
	source["date_histogram"] = opts

	// ValuesSourceAggregationBuilder
	if a.field != "" {
		opts["field"] = a.field
	}
	if a.script != nil {
		src, err := a.script.Source()
		if err != nil {
			return nil, err
		}
		opts["script"] = src
	}
	if a.missing != nil {
		opts["missing"] = a.missing
	}

	if s := a.interval; s != "" {
		opts["interval"] = s
	}
	if s := a.fixedInterval; s != "" {
		opts["fixed_interval"] = s
	}
	if s := a.calendarInterval; s != "" {
		opts["calendar_interval"] = s
	}

	if a.minDocCount != nil {
		opts["min_doc_count"] = *a.minDocCount
	}
	if a.order != "" {
		o := make(map[string]interface{})
		if a.orderAsc {
			o[a.order] = "asc"
		} else {
			o[a.order] = "desc"
		}
		opts["order"] = o
	}
	if a.timeZone != "" {
		opts["time_zone"] = a.timeZone
	}
	if a.offset != "" {
		opts["offset"] = a.offset
	}
	if a.format != "" {
		opts["format"] = a.format
	}
	if a.extendedBoundsMin != nil || a.extendedBoundsMax != nil {
		bounds := make(map[string]interface{})
		if a.extendedBoundsMin != nil {
			bounds["min"] = a.extendedBoundsMin
		}
		if a.extendedBoundsMax != nil {
			bounds["max"] = a.extendedBoundsMax
		}
		opts["extended_bounds"] = bounds
	}
	if a.keyed != nil {
		opts["keyed"] = *a.keyed
	}

	// AggregationBuilder (SubAggregations)
	if len(a.subAggregations) > 0 {
		aggsMap := make(map[string]interface{})
		source["aggregations"] = aggsMap
		for name, aggregate := range a.subAggregations {
			src, err := aggregate.Source()
			if err != nil {
				return nil, err
			}
			aggsMap[name] = src
		}
	}

	// Add Meta data if available
	if len(a.meta) > 0 {
		source["meta"] = a.meta
	}

	return source, nil
}
