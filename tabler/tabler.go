package tabler

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type RowPrepper interface {
	PrepRow() []string
}

// practice with generics
func PrepStructRow[T any](s T) []string {
	v := reflect.ValueOf(s).Elem()
	numFields := v.NumField()
	row := make([]string, numFields)
	// iterate over fields and add type specific logic to format field as string in row
	for i := 0; i < numFields; i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String { // Check if it's a slice of strings
				tags := make([]string, 0, field.Len())
				for j := 0; j < field.Len(); j++ {
					tags = append(tags, field.Index(j).String())
				}
				row[i] = strings.Join(tags, ", ")
			} else {
				row[i] = field.String()
			}
		case reflect.Int:
			row[i] = strconv.Itoa(int(field.Int()))
		case reflect.String:
			row[i] = field.String()
		case reflect.Bool:
			row[i] = strconv.FormatBool(field.Bool())
		case reflect.Float64:
			row[i] = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Float32:
			row[i] = strconv.FormatFloat(field.Float(), 'f', -1, 32)
		default:
			// attempt to cast to string
			row[i] = field.String()
		}
	}
	return row
}

func PrepTable(data []RowPrepper, headers []string) [][]string {
	if len(data) == 0 {
		return [][]string{}
	}
	fmt.Println("data:\n", data)
	dataFrame := make([][]string, len(data) + 1)
	dataFrame[0] = headers
	for i, row := range data {
		dataFrame[i+1] = row.PrepRow()
	}
	return dataFrame
}