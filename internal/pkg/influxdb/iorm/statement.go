package iorm

import (
	"fmt"
	"strings"

	"github.com/githubzjm/tuo/internal/pkg/influxdb/iorm/schema"
	log "github.com/sirupsen/logrus"
)

// Only for query, because write use point instead of flux statement
type Statement struct {
	Flux strings.Builder
}

func NewStatement() *Statement {
	return &Statement{
		Flux: strings.Builder{},
	}
}

func (s *Statement) Build() string {
	return s.Flux.String()
}

func (s *Statement) From(bucket string) *Statement {
	s.Flux.WriteString(fmt.Sprintf(`from(bucket: "%s")`, bucket))
	return s
}

// start is necessary, stop default is "now()"
func (s *Statement) Range(start, stop string) *Statement {
	if start == "" {
		return s
	}
	var res string
	if stop == "" {
		res = fmt.Sprintf(`|> range(start: %s)`, start)
	} else {
		res = fmt.Sprintf(`|> range(start: %s, stop: %s)`, start, stop)
	}
	s.Flux.WriteString(res)
	return s
}

// can filter only one _field and _value pair
// rule can be >, <, ==
func (s *Statement) Filter(query interface{}, rule string) *Statement {
	point := schema.Parse(query)

	measurement := point.Name()
	measurementFilter := fmt.Sprintf(`r._measurement == "%s"`, measurement)

	tagsFilters := make([]string, 10)
	for _, tag := range point.TagList() {
		tagsFilters = append(tagsFilters, fmt.Sprintf(`r.%s == "%s"`, tag.Key, tag.Value))
	}

	var fieldFilter, valueFilter string
	for _, field := range point.FieldList() {
		fieldFilter = fmt.Sprintf(`r._field == "%s"`, field.Key)
		if rule != "" {
			valueFilter = fmt.Sprintf(`r._value %s %s`, rule, field.Value)
		}
	}

	filters := make([]string, 20)
	filters = append(filters, measurementFilter)
	filters = append(filters, tagsFilters...)
	if fieldFilter != "" {
		filters = append(filters, fieldFilter)
	}
	if valueFilter != "" {
		filters = append(filters, valueFilter)
	}
	res := fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(filters, "and"))
	s.Flux.WriteString(res)
	return s
}

func (s *Statement) Keep(args ...interface{}) *Statement {
	columns := make([]string, 10)
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			columns = append(columns, arg)
		case []string:
			columns = append(columns, arg...)
		default:
			log.Errorf("unsupported select args %v", args)
			return s
		}
	}
	for i, column := range columns {
		columns[i] = fmt.Sprintf(`"%s"`, column)
	}
	res := fmt.Sprintf(`|> keep(columns: [%s])`, strings.Join(columns, ", "))
	s.Flux.WriteString(res)
	return s
}

func (s *Statement) Last() *Statement {
	s.Flux.WriteString(`|> last()`)
	return s
}
