package dmgbuild

import "fmt"

type arrageOrder string

//nolint:golint
var (
	NameOrder           arrageOrder = "name"
	DateModifiedOrder   arrageOrder = "date-modified"
	DateCreatedOrder    arrageOrder = "date-created"
	DateAddedOrder      arrageOrder = "date-added"
	DateLastOpenedOrder arrageOrder = "date-last-opened"
	SizeOrder           arrageOrder = "size"
	KindOrder           arrageOrder = "kind"
	LabelOrder          arrageOrder = "label"
)

type labelPosition string

//nolint:golint
var (
	LabelBottom labelPosition = "bottom"
	LabelRight  labelPosition = "right"
)

type IconView struct {
	ArrangeBy     arrageOrder
	GridOffsetX   int
	GridOffsetY   int
	GridSpacing   float32
	ScrollPosX    float32
	ScrollPosY    float32
	LabelPosition labelPosition
	IconSize      float32
	TextSize      float32
}

func NewIconView() IconView {
	return IconView{
		GridOffsetX:   0,
		GridOffsetY:   0,
		GridSpacing:   100,
		ScrollPosX:    0.0,
		ScrollPosY:    0.0,
		LabelPosition: LabelBottom,
		IconSize:      128,
		TextSize:      16,
	}
}

func (s *IconView) Render() []string {
	r := []string{}

	if s.ArrangeBy != "" {
		r = append(r, "arrange_by = "+pyStr(string(s.ArrangeBy))+"\n")
	}
	if s.GridOffsetX > 0 || s.GridOffsetY > 0 {
		r = append(r, fmt.Sprintf(
			"grid_offset = (%d, %d)\n",
			s.GridOffsetX, s.GridOffsetY,
		))
	}
	if s.GridSpacing > 0 {
		r = append(r, fmt.Sprintf("grid_spacing = %.2f\n", s.GridSpacing))
	}
	if s.ScrollPosX > 0 || s.ScrollPosY > 0 {
		r = append(r, fmt.Sprintf(
			"scroll_position = (%.2f, %.2f)\n",
			s.ScrollPosX, s.ScrollPosY,
		))
	}
	if s.LabelPosition != "" {
		r = append(r, "label_position = "+pyStr(string(s.LabelPosition))+"\n")
	}

	if s.IconSize > 0 {
		r = append(r, fmt.Sprintf("icon_size = %.2f\n", s.IconSize))
	}
	if s.TextSize > 0 {
		r = append(r, fmt.Sprintf("text_size = %.2f\n", s.TextSize))
	}

	return r
}
