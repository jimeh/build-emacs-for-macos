package dmgbuild

import (
	"fmt"
	"sort"
	"strings"
)

type listColumn string

//nolint:golint
var (
	NameColumn           listColumn = "name"
	DateModifiedColumn   listColumn = "date-modified"
	DateCreatedColumn    listColumn = "date-created"
	DateAddedColumn      listColumn = "date-added"
	DateLastOpenedColumn listColumn = "date-last-opened"
	SizeColumn           listColumn = "size"
	KindColumn           listColumn = "kind"
	LabelColumn          listColumn = "label"
	VersionColumn        listColumn = "version"
	CommentsColumn       listColumn = "comments"
)

type direction string

//nolint:golint
var (
	Ascending  direction = "ascending"
	Descending direction = "descending"
)

type ListView struct {
	SortBy               listColumn
	ScrollPosX           int
	ScrollPosY           int
	IconSize             float32
	TextSize             float32
	UseRelativeDates     bool
	CalculateAllSizes    bool
	Columns              []listColumn
	ColumnWidths         map[listColumn]int
	ColumnSortDirections map[listColumn]direction
}

func NewListView() ListView {
	return ListView{
		SortBy:           NameColumn,
		IconSize:         16,
		TextSize:         12,
		UseRelativeDates: true,
		Columns: []listColumn{
			NameColumn,
			DateModifiedColumn,
			SizeColumn,
			KindColumn,
			DateAddedColumn,
		},
		ColumnWidths: map[listColumn]int{
			(NameColumn):           300,
			(DateModifiedColumn):   181,
			(DateCreatedColumn):    181,
			(DateAddedColumn):      181,
			(DateLastOpenedColumn): 181,
			(SizeColumn):           97,
			(KindColumn):           115,
			(LabelColumn):          100,
			(VersionColumn):        75,
			(CommentsColumn):       300,
		},
		ColumnSortDirections: map[listColumn]direction{
			(NameColumn):           Ascending,
			(DateModifiedColumn):   Descending,
			(DateCreatedColumn):    Descending,
			(DateAddedColumn):      Descending,
			(DateLastOpenedColumn): Descending,
			(SizeColumn):           Descending,
			(KindColumn):           Ascending,
			(LabelColumn):          Ascending,
			(VersionColumn):        Ascending,
			(CommentsColumn):       Ascending,
		},
	}
}

func (s *ListView) Render() []string {
	r := []string{}

	if s.SortBy != "" {
		r = append(r, "list_sort_by = "+pyStr(string(s.SortBy))+"\n")
	}
	if s.ScrollPosX > 0 || s.ScrollPosY > 0 {
		r = append(r, fmt.Sprintf(
			"list_scroll_position = (%d, %d)\n",
			s.ScrollPosX, s.ScrollPosY,
		))
	}
	if s.IconSize > 0 {
		r = append(r, fmt.Sprintf("list_icon_size = %.2f\n", s.IconSize))
	}
	if s.TextSize > 0 {
		r = append(r, fmt.Sprintf("list_text_size = %.2f\n", s.TextSize))
	}
	r = append(r, "list_use_relative_dates = "+pyBool(s.UseRelativeDates)+"\n")
	r = append(
		r, "list_calculate_all_sizes = "+pyBool(s.CalculateAllSizes)+"\n",
	)

	if len(s.Columns) > 0 {
		var cols []string
		for _, col := range s.Columns {
			cols = append(cols, pyStr(string(col)))
		}
		r = append(r,
			"list_columns = [\n    "+strings.Join(cols, ",\n    ")+"\n]\n",
		)
	}

	if len(s.ColumnWidths) > 0 {
		var cols []string
		for col, w := range s.ColumnWidths {
			cols = append(cols, fmt.Sprintf(
				"%s: %d", pyStr(string(col)), w,
			))
		}
		sort.SliceStable(cols, func(i, j int) bool {
			return cols[i] < cols[j]
		})
		r = append(r,
			"list_column_widths = {\n    "+
				strings.Join(cols, ",\n    ")+
				"\n}\n",
		)
	}

	if len(s.ColumnSortDirections) > 0 {
		var cols []string
		for col, direction := range s.ColumnSortDirections {
			cols = append(cols, fmt.Sprintf(
				"%s: %s", pyStr(string(col)), pyStr(string(direction)),
			))
		}
		sort.SliceStable(cols, func(i, j int) bool {
			return cols[i] < cols[j]
		})
		r = append(r,
			"list_column_sort_directions = {\n    "+
				strings.Join(cols, ",\n    ")+
				"\n}\n",
		)
	}

	return r
}
