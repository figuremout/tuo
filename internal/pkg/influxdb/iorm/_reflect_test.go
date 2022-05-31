package main

import (
	"fmt"
	"go/ast"
	"reflect"
)

type Metric struct {
	Measurement string
	Tags        interface{}
	Fields      interface{}
}

type CPU struct {
	Metric
	Tags   CPUTags
	Fields CPUFields
}

type CPUTags struct {
	CPU  string
	Host string
}

type CPUFields struct {
	Percent int
}

func test(dest interface{}) {
	modelValue := reflect.Indirect(reflect.ValueOf(dest))
	modelType := modelValue.Type()

	for i := 0; i < modelValue.NumField(); i++ {
		v := modelValue.Field(i) // Value
		if ast.IsExported(modelType.Field(i).Name) {
			fmt.Printf("%d: %s %s = %#v\n", i,
				modelType.Field(i).Name, v.Type(), v.Interface())

			if modelType.Field(i).Name == "Tags" {
				tagValue := reflect.Indirect(reflect.ValueOf(v.Interface()))
				tagType := v.Type()
				for i := 0; i < tagValue.NumField(); i++ {
					v := tagValue.Field(i) // Value
					if tagType.Field(i).IsExported() {
						fmt.Printf("%d: %s %s = %#v\n", i,
							tagType.Field(i).Name, v.Type(), v.String())
					}
				}
			}
		}
	}

	for i := 0; i < modelType.NumField(); i++ {
		f := modelType.Field(i) // StructField
		if !f.Anonymous && ast.IsExported(f.Name) {
			fmt.Println(f.Name)
		}
	}
}

func main() {
	cpu := CPU{
		Metric: Metric{Measurement: "cpu"},
		Tags:   CPUTags{CPU: "cpu-total", Host: "local"},
		Fields: CPUFields{Percent: 50},
	}
	test(cpu)
}
