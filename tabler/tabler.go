package tabler

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
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
	dataFrame := make([][]string, len(data)+1)
	dataFrame[0] = headers
	for i, row := range data {
		dataFrame[i+1] = row.PrepRow()
	}
	return dataFrame
}

func RenderTable(data [][]string, headers []string) error {
	/*table := tablewriter.NewWriter(os.Stdout)
	table.Configure(func(cfg *tablewriter.Config) {
		cfg.MaxWidth = 120
		cfg.Row.Formatting.AutoWrap = 2
	})*/

	colorCfg := renderer.ColorizedConfig{
		Header: renderer.Tint{
			FG: renderer.Colors{color.FgGreen, color.Bold}, // Green bold headers
		},
		Column: renderer.Tint{
			FG: renderer.Colors{color.FgCyan}, // Default cyan for rows
			Columns: []renderer.Tint{
				{FG: renderer.Colors{color.FgMagenta}}, // Magenta for column 0
				{},                                     // Inherit default (cyan)
				{FG: renderer.Colors{color.FgHiRed}},   // High-intensity red for column 2
			},
		},
		Footer: renderer.Tint{
			FG: renderer.Colors{color.FgYellow, color.Bold}, // Yellow bold footer
			Columns: []renderer.Tint{
				{},                                      // Inherit default
				{FG: renderer.Colors{color.FgHiYellow}}, // High-intensity yellow for column 1
				{},                                      // Inherit default
			},
		},
		Border:    renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White borders
		Separator: renderer.Tint{FG: renderer.Colors{color.FgWhite}}, // White separators
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRenderer(renderer.NewColorized(colorCfg)),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal}, // Wrap long content
				Alignment:    tw.CellAlignment{Global: tw.AlignLeft},     // Left-align rows
				ColMaxWidths: tw.CellWidth{Global: 40},
			},
			Footer: tw.CellConfig{
				Alignment: tw.CellAlignment{Global: tw.AlignRight},
			},
		}),
	)
	table.Header(headers)
	table.Bulk(data)
	err := table.Render()
	if err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}
	return nil
}
