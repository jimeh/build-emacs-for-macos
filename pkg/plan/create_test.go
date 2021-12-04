package plan

import (
	"testing"

	"github.com/jimeh/build-emacs-for-macos/pkg/release"
	"github.com/stretchr/testify/assert"
)

func Test_parseGitRef(t *testing.T) {
	t.Parallel()

	type args struct {
		ref string
	}
	type want struct {
		version string
		channel release.Channel
		err     string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "master",
			args: args{ref: "master"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "emacs-28",
			args: args{ref: "emacs-28"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "emacs-27",
			args: args{ref: "emacs-27"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "emacs-26",
			args: args{ref: "emacs-26"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "emacs-24",
			args: args{ref: "emacs-24"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "feature/native-comp",
			args: args{ref: "feature/native-comp"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "feature/pgtk",
			args: args{ref: "feature/pgtk"},
			want: want{version: "", channel: release.Nightly, err: ""},
		},
		{
			name: "emacs-19.34",
			args: args{ref: "emacs-19.34"},
			want: want{version: "19.34", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-20.4",
			args: args{ref: "emacs-20.4"},
			want: want{version: "20.4", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-22.3",
			args: args{ref: "emacs-22.3"},
			want: want{version: "22.3", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-23.4",
			args: args{ref: "emacs-23.4"},
			want: want{version: "23.4", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-24.0.97",
			args: args{ref: "emacs-24.0.97"},
			want: want{version: "24.0.97", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.2",
			args: args{ref: "emacs-24.2"},
			want: want{version: "24.2", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-24.2.90",
			args: args{ref: "emacs-24.2.90"},
			want: want{version: "24.2.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.2.93",
			args: args{ref: "emacs-24.2.93"},
			want: want{version: "24.2.93", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.3",
			args: args{ref: "emacs-24.3"},
			want: want{version: "24.3", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-24.3-rc1",
			args: args{ref: "emacs-24.3-rc1"},
			want: want{version: "24.3-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-24.3.90",
			args: args{ref: "emacs-24.3.90"},
			want: want{version: "24.3.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.3.94",
			args: args{ref: "emacs-24.3.94"},
			want: want{version: "24.3.94", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.4",
			args: args{ref: "emacs-24.4"},
			want: want{version: "24.4", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-24.4-rc1",
			args: args{ref: "emacs-24.4-rc1"},
			want: want{version: "24.4-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-24.4.90",
			args: args{ref: "emacs-24.4.90"},
			want: want{version: "24.4.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.4.91",
			args: args{ref: "emacs-24.4.91"},
			want: want{version: "24.4.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-24.5",
			args: args{ref: "emacs-24.5"},
			want: want{version: "24.5", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-24.5-rc1",
			args: args{ref: "emacs-24.5-rc1"},
			want: want{version: "24.5-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-24.5-rc3",
			args: args{ref: "emacs-24.5-rc3"},
			want: want{version: "24.5-rc3", channel: release.RC, err: ""},
		},
		{
			name: "emacs-24.5-rc3-fixed",
			args: args{ref: "emacs-24.5-rc3-fixed"},
			want: want{version: "24.5-rc3-fixed", channel: release.RC, err: ""},
		},
		{
			name: "emacs-25.0.90",
			args: args{ref: "emacs-25.0.90"},
			want: want{version: "25.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-25.0.95",
			args: args{ref: "emacs-25.0.95"},
			want: want{version: "25.0.95", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-25.1",
			args: args{ref: "emacs-25.1"},
			want: want{version: "25.1", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-25.1-rc1",
			args: args{ref: "emacs-25.1-rc1"},
			want: want{version: "25.1-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-25.1-rc2",
			args: args{ref: "emacs-25.1-rc2"},
			want: want{version: "25.1-rc2", channel: release.RC, err: ""},
		},
		{
			name: "emacs-25.1.90",
			args: args{ref: "emacs-25.1.90"},
			want: want{version: "25.1.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-25.1.91",
			args: args{ref: "emacs-25.1.91"},
			want: want{version: "25.1.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-25.2",
			args: args{ref: "emacs-25.2"},
			want: want{version: "25.2", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-25.2-rc1",
			args: args{ref: "emacs-25.2-rc1"},
			want: want{version: "25.2-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-25.2-rc2",
			args: args{ref: "emacs-25.2-rc2"},
			want: want{version: "25.2-rc2", channel: release.RC, err: ""},
		},
		{
			name: "emacs-26.0.90",
			args: args{ref: "emacs-26.0.90"},
			want: want{version: "26.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-26.0.91",
			args: args{ref: "emacs-26.0.91"},
			want: want{version: "26.0.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-26.1",
			args: args{ref: "emacs-26.1"},
			want: want{version: "26.1", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-26.1-rc1",
			args: args{ref: "emacs-26.1-rc1"},
			want: want{version: "26.1-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-26.1.90",
			args: args{ref: "emacs-26.1.90"},
			want: want{version: "26.1.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-26.1.92",
			args: args{ref: "emacs-26.1.92"},
			want: want{version: "26.1.92", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-26.2",
			args: args{ref: "emacs-26.2"},
			want: want{version: "26.2", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-26.2.90",
			args: args{ref: "emacs-26.2.90"},
			want: want{version: "26.2.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-26.3",
			args: args{ref: "emacs-26.3"},
			want: want{version: "26.3", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-26.3-rc1",
			args: args{ref: "emacs-26.3-rc1"},
			want: want{version: "26.3-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-27.0.90",
			args: args{ref: "emacs-27.0.90"},
			want: want{version: "27.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-27.0.91",
			args: args{ref: "emacs-27.0.91"},
			want: want{version: "27.0.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-27.1",
			args: args{ref: "emacs-27.1"},
			want: want{version: "27.1", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-27.1-rc1",
			args: args{ref: "emacs-27.1-rc1"},
			want: want{version: "27.1-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-27.1-rc2",
			args: args{ref: "emacs-27.1-rc2"},
			want: want{version: "27.1-rc2", channel: release.RC, err: ""},
		},
		{
			name: "emacs-27.1.90",
			args: args{ref: "emacs-27.1.90"},
			want: want{version: "27.1.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-27.1.91",
			args: args{ref: "emacs-27.1.91"},
			want: want{version: "27.1.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-27.2",
			args: args{ref: "emacs-27.2"},
			want: want{version: "27.2", channel: release.Stable, err: ""},
		},
		{
			name: "emacs-27.2-rc1",
			args: args{ref: "emacs-27.2-rc1"},
			want: want{version: "27.2-rc1", channel: release.RC, err: ""},
		},
		{
			name: "emacs-27.2-rc2",
			args: args{ref: "emacs-27.2-rc2"},
			want: want{version: "27.2-rc2", channel: release.RC, err: ""},
		},
		{
			name: "emacs-28.0.90",
			args: args{ref: "emacs-28.0.90"},
			want: want{version: "28.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-21.0.100",
			args: args{ref: "emacs-pretest-21.0.100"},
			want: want{version: "21.0.100", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-21.0.106",
			args: args{ref: "emacs-pretest-21.0.106"},
			want: want{version: "21.0.106", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-21.0.90",
			args: args{ref: "emacs-pretest-21.0.90"},
			want: want{version: "21.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-21.0.99",
			args: args{ref: "emacs-pretest-21.0.99"},
			want: want{version: "21.0.99", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-22.0.90",
			args: args{ref: "emacs-pretest-22.0.90"},
			want: want{version: "22.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-22.0.99",
			args: args{ref: "emacs-pretest-22.0.99"},
			want: want{version: "22.0.99", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-22.0.990",
			args: args{ref: "emacs-pretest-22.0.990"},
			want: want{version: "22.0.990", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-22.1.90",
			args: args{ref: "emacs-pretest-22.1.90"},
			want: want{version: "22.1.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-22.2.90",
			args: args{ref: "emacs-pretest-22.2.90"},
			want: want{version: "22.2.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.0.90",
			args: args{ref: "emacs-pretest-23.0.90"},
			want: want{version: "23.0.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.1.90",
			args: args{ref: "emacs-pretest-23.1.90"},
			want: want{version: "23.1.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.2.90",
			args: args{ref: "emacs-pretest-23.2.90"},
			want: want{version: "23.2.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.2.91",
			args: args{ref: "emacs-pretest-23.2.91"},
			want: want{version: "23.2.91", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.2.93",
			args: args{ref: "emacs-pretest-23.2.93"},
			want: want{version: "23.2.93", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.2.93.1",
			args: args{ref: "emacs-pretest-23.2.93.1"},
			want: want{version: "23.2.93.1", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-23.3.90",
			args: args{ref: "emacs-pretest-23.3.90"},
			want: want{version: "23.3.90", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-24.0.05",
			args: args{ref: "emacs-pretest-24.0.05"},
			want: want{version: "24.0.05", channel: release.Pretest, err: ""},
		},
		{
			name: "emacs-pretest-24.0.90",
			args: args{ref: "emacs-pretest-24.0.90"},
			want: want{version: "24.0.90", channel: release.Pretest, err: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotChannel, err := parseGitRef(tt.args.ref)

			assert.Equal(t, tt.want.version, got)
			assert.Equal(t, tt.want.channel, gotChannel)

			if tt.want.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want.err)
			}
		})
	}
}
