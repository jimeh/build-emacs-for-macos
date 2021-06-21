package dmgbuild

import "fmt"

type view string

//nolint:golint
var (
	Icon      view = "icon-view"
	list      view = "list-view"
	Column    view = "column-view"
	Coverflow view = "coverflow"
)

type Window struct {
	PoxX                    int
	PosY                    int
	Width                   int
	Height                  int
	Background              string
	ShowStatusBar           bool
	ShowTabView             bool
	ShowToolbar             bool
	ShowPathbar             bool
	ShowSidebar             bool
	SidebarWidth            int
	DefaultView             view
	ShowIconPreview         bool
	ShowItemInfo            bool
	IncludeIconViewSettings bool
	IncludeListViewSettings bool
}

func NewWindow() Window {
	return Window{
		PoxX:        100,
		PosY:        150,
		Width:       640,
		Height:      280,
		Background:  "builtin-arrow",
		DefaultView: Icon,
	}
}

func (s *Window) Render() []string {
	r := []string{}

	if s.Background != "" {
		r = append(r, "background = "+pyStr(s.Background)+"\n")
	}

	r = append(r, "show_status_bar = "+pyBool(s.ShowStatusBar)+"\n")
	r = append(r, "show_tab_view = "+pyBool(s.ShowTabView)+"\n")
	r = append(r, "show_toolbar = "+pyBool(s.ShowToolbar)+"\n")
	r = append(r, "show_pathbar = "+pyBool(s.ShowPathbar)+"\n")
	r = append(r, "show_sidebar = "+pyBool(s.ShowSidebar)+"\n")

	if s.SidebarWidth > 0 {
		r = append(r, fmt.Sprintf(
			"sidebar_width = %d\n", s.SidebarWidth,
		))
	}
	if s.DefaultView != "" {
		r = append(r, "default_view = "+pyStr(string(s.DefaultView))+"\n")
	}
	if s.Width > 0 && s.Height > 0 {
		r = append(r, fmt.Sprintf(
			"window_rect = ((%d, %d), (%d, %d))\n",
			s.PoxX, s.PosY, s.Width, s.Height,
		))
	}

	r = append(r, "show_icon_preview = "+pyBool(s.ShowIconPreview)+"\n")
	r = append(r, "show_item_info = "+pyBool(s.ShowIconPreview)+"\n")
	r = append(
		r, "include_icon_view_settings = "+pyBool(s.ShowIconPreview)+"\n",
	)
	r = append(
		r, "include_list_view_settings = "+pyBool(s.ShowIconPreview)+"\n",
	)

	return r
}
