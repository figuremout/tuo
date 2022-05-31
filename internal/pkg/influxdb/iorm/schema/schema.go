package schema

import (
	"reflect"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Model struct {
	Start       time.Time `iorm:"_start"`
	Stop        time.Time `iorm:"_stop"`
	Time        time.Time `iorm:"_time"`
	Measurement string    `iorm:"_measurement"`
	Tags        interface{}
	Fields      interface{}
}

func Parse(dest interface{}) *write.Point {
	modelValue := reflect.Indirect(reflect.ValueOf(dest))

	measurementValue := modelValue.FieldByName("Measurement")
	timeValue := modelValue.FieldByName("Time")
	tagsValue := modelValue.FieldByName("Tags")
	fieldsValue := modelValue.FieldByName("Fields")

	if measurementValue.IsZero() || timeValue.IsZero() || tagsValue.IsZero() || fieldsValue.IsZero() {
		return nil
	}

	measurement := measurementValue.String()

	Itime := timeValue.Interface()
	ts, ok := Itime.(time.Time)
	if !ok {
		return nil
	}

	tags := parseTags(tagsValue)

	fields := parseFields(fieldsValue)

	return influxdb2.NewPoint(measurement, tags, fields, ts)
}

func parseTags(v reflect.Value) map[string]string { // v is tags
	res := make(map[string]string)
	tagsValue := reflect.Indirect(reflect.ValueOf(v.Interface()))
	tagsType := v.Type()
	for i := 0; i < tagsValue.NumField(); i++ {
		v := tagsValue.Field(i)
		t := tagsType.Field(i)
		if t.IsExported() && !v.IsZero() { // only handle exported and assigned explicitly fields
			if tag, ok := t.Tag.Lookup("iorm"); ok {
				res[tag] = v.String()
			} else {
				res[strings.ToLower(tagsType.Field(i).Name)] = v.String()
			}
			// fmt.Printf("%d: %s %s = %#v\n", i,
			// 	tagType.Field(i).Name, v.Type(), v.String())
		}
	}
	return res
}

func parseFields(v reflect.Value) map[string]interface{} {
	res := make(map[string]interface{})
	fieldsValue := reflect.Indirect(reflect.ValueOf(v.Interface()))
	fieldsType := v.Type()
	for i := 0; i < fieldsValue.NumField(); i++ {
		v := fieldsValue.Field(i)
		t := fieldsType.Field(i)
		if t.IsExported() && !v.IsZero() {
			if tag, ok := t.Tag.Lookup("iorm"); ok {
				res[tag] = v.Interface()
			} else {
				res[strings.ToLower(fieldsType.Field(i).Name)] = v.Interface()
			}
		}
	}
	return res
}

// load src into dest struct
func Load(src map[string]interface{}, dest interface{}) {
	reflectValue := reflect.ValueOf(dest)
	reflectType := reflectValue.Type()

	if reflectValue.Kind() != reflect.Ptr { // dest is not a ptr
		return
	}

	fieldTagNameMap := make(map[string]string)
	for i := 0; i < reflectValue.NumField(); i++ {
		v := reflectValue.Field(i)
		t := reflectType.Field(i)
		if t.IsExported() && !v.IsZero() {
			if tag, ok := t.Tag.Lookup("iorm"); ok {
				fieldTagNameMap[tag] = t.Name
			}
		}
	}

	value := reflectValue.Elem() // the value of dest
	for k, v := range src {
		var name string
		if n, ok := fieldTagNameMap[k]; ok { // key is defined by tag
			name = n
		} else { // key is just the lower case of field name
			name = strings.ToUpper(k[:1]) + k[1:] // upper the first letter
		}

		value.FieldByName(name).Set(reflect.ValueOf(v).Elem()) // assign the field of dest
	}
}
