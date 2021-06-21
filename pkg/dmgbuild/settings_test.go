package dmgbuild

import (
	"bytes"
	"strings"
	"testing"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSettings_Write(t *testing.T) {
	test := []struct {
		name         string
		entitlements *Settings
		want         string
	}{
		{
			name:         "empty",
			entitlements: &Settings{},
			want: undent.String(`
                # -*- coding: utf-8 -*-
                from __future__ import unicode_literals
                show_status_bar = False
                show_tab_view = False
                show_toolbar = False
                show_pathbar = False
                show_sidebar = False
                show_icon_preview = False
                show_item_info = False
                include_icon_view_settings = False
                include_list_view_settings = False
                list_use_relative_dates = False
                list_calculate_all_sizes = False`,
			),
		},
		{
			name: "full",
			entitlements: &Settings{
				Filename:         "/builds/Emacs.2021-05-25.f4dc646.master.dmg",
				VolumeName:       "Emacs.2021-05-25.f4dc646.master",
				Format:           UDBZFormat,
				CompressionLevel: 8,
				Size:             "100m",
				Files: []*File{
					{
						Path: "/builds/Emacs.app",
						PosX: 200,
						PosY: 200,
					},
					{
						Path:          "/builds/README.rtf",
						PosX:          200,
						PosY:          300,
						HideExtension: true,
					},
					{
						Path:   "/builds/hide-me.png",
						Hidden: true,
					},
				},
				Symlinks: []*Symlink{
					{
						Name:   "Applications",
						Target: "/Applications",
						PosX:   400,
						PosY:   400,
					},
					{
						Name:   "QuickLook",
						Target: "/Library/QuickLook",
						PosX:   500,
						PosY:   400,
						Hidden: true,
					},
					{
						Name:          "System",
						Target:        "/System",
						HideExtension: true,
					},
				},
				Icon:      "/opt/misc/assets/volIcon.icns",
				BadgeIcon: "/builds/Emacs.app/Contents/Resources/Icon.icns",
				Window: Window{
					PoxX:                    200,
					PosY:                    250,
					Width:                   680,
					Height:                  446,
					Background:              "/opt/misc/assets/bg.tif",
					ShowStatusBar:           true,
					ShowTabView:             true,
					ShowToolbar:             true,
					ShowPathbar:             true,
					ShowSidebar:             true,
					SidebarWidth:            165,
					DefaultView:             list,
					ShowIconPreview:         true,
					ShowItemInfo:            true,
					IncludeIconViewSettings: true,
					IncludeListViewSettings: true,
				},
				IconView: IconView{
					ArrangeBy:     NameOrder,
					GridOffsetX:   42,
					GridOffsetY:   43,
					GridSpacing:   44.5,
					ScrollPosX:    4.5,
					ScrollPosY:    5.5,
					LabelPosition: LabelBottom,
					IconSize:      160,
					TextSize:      15,
				},
				ListView: ListView{
					SortBy:            NameColumn,
					ScrollPosX:        7,
					ScrollPosY:        8,
					IconSize:          16,
					TextSize:          12,
					UseRelativeDates:  true,
					CalculateAllSizes: true,
					Columns: []listColumn{
						NameColumn,
						DateModifiedColumn,
						DateCreatedColumn,
						DateAddedColumn,
						DateLastOpenedColumn,
						SizeColumn,
						KindColumn,
						LabelColumn,
						VersionColumn,
						CommentsColumn,
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
				},
				License: License{
					DefaultLanguage: LocaleEnUS,
					Licenses: map[locale]string{
						//nolint:lll
						(LocaleEnGB): undent.String(`
                            {\rtf1\ansi\ansicpg1252\cocoartf1504\cocoasubrtf820
                             {\fonttbl\f0\fnil\fcharset0 Helvetica-Bold;\f1\fnil\fcharset0 Helvetica;}
                             {\colortbl;\red255\green255\blue255;\red0\green0\blue0;}
                             {\*\expandedcolortbl;;\cssrgb\c0\c0\c0;}
                             \paperw11905\paperh16837\margl1133\margr1133\margb1133\margt1133
                             \deftab720
                             \pard\pardeftab720\sa160\partightenfactor0
                             \f0\b\fs60 \cf2 \expnd0\expndtw0\kerning0
                             \up0 \nosupersub \ulnone \outl0\strokewidth0 \strokec2 Test License\
                             \pard\pardeftab720\sa160\partightenfactor0
                             \fs36 \cf2 \strokec2 What is this?\
                             \pard\pardeftab720\sa160\partightenfactor0
                             \f1\b0\fs22 \cf2 \strokec2 This is the English license. It says what you are allowed to do with this software.\
                             \
                             }`,
						),
						//nolint:lll
						(LocaleSe): undent.String(`
                            {\rtf1\ansi\ansicpg1252\cocoartf1504\cocoasubrtf820
                             {\fonttbl\f0\fnil\fcharset0 Helvetica-Bold;\f1\fnil\fcharset0 Helvetica;}
                             {\colortbl;\red255\green255\blue255;\red0\green0\blue0;}
                             {\*\expandedcolortbl;;\cssrgb\c0\c0\c0;}
                             \paperw11905\paperh16837\margl1133\margr1133\margb1133\margt1133
                             \deftab720
                             \pard\pardeftab720\sa160\partightenfactor0
                             \f0\b\fs60 \cf2 \expnd0\expndtw0\kerning0
                             \up0 \nosupersub \ulnone \outl0\strokewidth0 \strokec2 Test License\
                             \pard\pardeftab720\sa160\partightenfactor0
                             \fs36 \cf2 \strokec2 What is this?\
                             \pard\pardeftab720\sa160\partightenfactor0
                             \f1\b0\fs22 \cf2 \strokec2 Detta är den engelska licensen. Det står vad du får göra med den här programvaran.\
                             \
                             }`,
						),
					},
					Buttons: map[locale]Buttons{
						(LocaleEnGB): {
							LanguageName: "English",
							Agree:        "Agree",
							Disagree:     "Disagree",
							Print:        "Print",
							Save:         "Save",
							Message: "If you agree with the terms of this " +
								"license, press \"Agree\" to install the " +
								"software.  If you do not agree, press " +
								"\"Disagree\".",
						},
						(LocaleSe): {
							LanguageName: "Svenska",
							Agree:        "Godkänn",
							Disagree:     "Håller inte med",
							Print:        "Skriv ut",
							Save:         "Spara",
							Message: "Om du godkänner villkoren i denna " +
								"licens, tryck på \"Godkänn\" för att " +
								"installera programvaran. Om du inte håller " +
								"med, tryck på \"Håller inte med\".",
						},
					},
				},
			},
			//nolint:lll
			want: undent.String(`
                # -*- coding: utf-8 -*-
                from __future__ import unicode_literals
                filename = "/builds/Emacs.2021-05-25.f4dc646.master.dmg"
                volume_name = "Emacs.2021-05-25.f4dc646.master"
                format = "UDBZ"
                compression_level = 8
                size = "100m"
                files = [
                    "/builds/Emacs.app",
                    "/builds/README.rtf",
                    "/builds/hide-me.png"
                ]
                symlinks = {
                    "Applications": "/Applications",
                    "QuickLook": "/Library/QuickLook",
                    "System": "/System"
                }
                hide = [
                    "hide-me.png",
                    "QuickLook"
                ]
                hide_extensions = [
                    "README.rtf",
                    "System"
                ]
                icon_locations = {
                    "Emacs.app": (200, 200),
                    "README.rtf": (200, 300),
                    "Applications": (400, 400),
                    "QuickLook": (500, 400)
                }
                icon = "/opt/misc/assets/volIcon.icns"
                badge_icon = "/builds/Emacs.app/Contents/Resources/Icon.icns"
                background = "/opt/misc/assets/bg.tif"
                show_status_bar = True
                show_tab_view = True
                show_toolbar = True
                show_pathbar = True
                show_sidebar = True
                sidebar_width = 165
                default_view = "list-view"
                window_rect = ((200, 250), (680, 446))
                show_icon_preview = True
                show_item_info = True
                include_icon_view_settings = True
                include_list_view_settings = True
                arrange_by = "name"
                grid_offset = (42, 43)
                grid_spacing = 44.50
                scroll_position = (4.50, 5.50)
                label_position = "bottom"
                icon_size = 160.00
                text_size = 15.00
                list_sort_by = "name"
                list_scroll_position = (7, 8)
                list_icon_size = 16.00
                list_text_size = 12.00
                list_use_relative_dates = True
                list_calculate_all_sizes = True
                list_columns = [
                    "name",
                    "date-modified",
                    "date-created",
                    "date-added",
                    "date-last-opened",
                    "size",
                    "kind",
                    "label",
                    "version",
                    "comments"
                ]
                list_column_widths = {
                    "comments": 300,
                    "date-added": 181,
                    "date-created": 181,
                    "date-last-opened": 181,
                    "date-modified": 181,
                    "kind": 115,
                    "label": 100,
                    "name": 300,
                    "size": 97,
                    "version": 75
                }
                list_column_sort_directions = {
                    "comments": "ascending",
                    "date-added": "descending",
                    "date-created": "descending",
                    "date-last-opened": "descending",
                    "date-modified": "descending",
                    "kind": "ascending",
                    "label": "ascending",
                    "name": "ascending",
                    "size": "descending",
                    "version": "ascending"
                }
                license = {
                    "default-language": "en_US",
                    "licenses": {
                        "en_GB": """{\\rtf1\\ansi\\ansicpg1252\\cocoartf1504\\cocoasubrtf820
                 {\\fonttbl\\f0\\fnil\\fcharset0 Helvetica-Bold;\\f1\\fnil\\fcharset0 Helvetica;}
                 {\\colortbl;\\red255\\green255\\blue255;\\red0\\green0\\blue0;}
                 {\\*\\expandedcolortbl;;\\cssrgb\\c0\\c0\\c0;}
                 \\paperw11905\\paperh16837\\margl1133\\margr1133\\margb1133\\margt1133
                 \\deftab720
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\f0\\b\\fs60 \\cf2 \\expnd0\\expndtw0\\kerning0
                 \\up0 \\nosupersub \\ulnone \\outl0\\strokewidth0 \\strokec2 Test License\\
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\fs36 \\cf2 \\strokec2 What is this?\\
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\f1\\b0\\fs22 \\cf2 \\strokec2 This is the English license. It says what you are allowed to do with this software.\\
                 \\
                 }""",
                        "se": """{\\rtf1\\ansi\\ansicpg1252\\cocoartf1504\\cocoasubrtf820
                 {\\fonttbl\\f0\\fnil\\fcharset0 Helvetica-Bold;\\f1\\fnil\\fcharset0 Helvetica;}
                 {\\colortbl;\\red255\\green255\\blue255;\\red0\\green0\\blue0;}
                 {\\*\\expandedcolortbl;;\\cssrgb\\c0\\c0\\c0;}
                 \\paperw11905\\paperh16837\\margl1133\\margr1133\\margb1133\\margt1133
                 \\deftab720
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\f0\\b\\fs60 \\cf2 \\expnd0\\expndtw0\\kerning0
                 \\up0 \\nosupersub \\ulnone \\outl0\\strokewidth0 \\strokec2 Test License\\
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\fs36 \\cf2 \\strokec2 What is this?\\
                 \\pard\\pardeftab720\\sa160\\partightenfactor0
                 \\f1\\b0\\fs22 \\cf2 \\strokec2 Detta är den engelska licensen. Det står vad du får göra med den här programvaran.\\
                 \\
                 }"""
                    },
                    "buttons": {
                        "en_GB": (
                            "English",
                            "Agree",
                            "Disagree",
                            "Print",
                            "Save",
                            "If you agree with the terms of this license, press \"Agree\" to install the software.  If you do not agree, press \"Disagree\"."
                        ),
                        "se": (
                            "Svenska",
                            "Godkänn",
                            "Håller inte med",
                            "Skriv ut",
                            "Spara",
                            "Om du godkänner villkoren i denna licens, tryck på \"Godkänn\" för att installera programvaran. Om du inte håller med, tryck på \"Håller inte med\"."
                        )
                    }
                }`,
			),
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := tt.entitlements.Write(&buf)
			require.NoError(t, err)

			assert.Equal(t, tt.want, strings.TrimSpace(buf.String()))
		})
	}
}

func Test_pyStr(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "empty",
			s:    "",
			want: `""`,
		},
		{
			name: "regular string",
			s:    "foo-bar nope :)",
			want: `"foo-bar nope :)"`,
		},
		{
			name: "with single quotes",
			s:    "john's lost 'flip-flop'",
			want: `"john's lost 'flip-flop'"`,
		},
		{
			name: "with double quotes",
			s:    `john has lost a "flip-flop"`,
			want: `"john has lost a \"flip-flop\""`,
		},
		{
			name: "with backslashes",
			s:    `C:\path\to\file.txt`,
			want: `"C:\\path\\to\\file.txt"`,
		},
		{
			name: "with line-feed",
			s:    "hello\nworld",
			want: `"hello\nworld"`,
		},
		{
			name: "with carriage return",
			s:    "hello\rworld",
			want: `"hello\rworld"`,
		},
		{
			name: "with backslashes, single and double quotes",
			s:    `john's "lost" C:\path\to\file.txt`,
			want: `"john's \"lost\" C:\\path\\to\\file.txt"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pyStr(tt.s)

			assert.Equal(t, tt.want, got)
		})
	}
}
