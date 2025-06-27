package release

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionToName(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr string
	}{
		{
			name: "empty",
			args: args{
				version: "",
			},
			wantErr: "release: empty version",
		},
		{
			name: "nightly",
			args: args{
				version: "2021-07-01.1b88404.master",
			},
			want: "Emacs.2021-07-01.1b88404.master",
		},
		{
			name: "nightly with variant",
			args: args{
				version: "2021-07-01.1b88404.master-1",
			},
			want: "Emacs.2021-07-01.1b88404.master-1",
		},
		{
			name: "pretest",
			args: args{
				version: "30.0.93-pretest",
			},
			want: "Emacs-30.0.93-pretest",
		},
		{
			name: "pretest with variant",
			args: args{
				version: "30.0.93-pretest-1",
			},
			want: "Emacs-30.0.93-pretest-1",
		},
		{
			name: "stable",
			args: args{
				version: "27.2",
			},
			want: "Emacs-27.2",
		},
		{
			name: "stable with letter",
			args: args{
				version: "23.3b",
			},
			want: "Emacs-23.3b",
		},
		{
			name: "stable with variant",
			args: args{
				version: "23.3-1",
			},
			want: "Emacs-23.3-1",
		},
		{
			name: "stable with letter and variant",
			args: args{
				version: "23.3b-1",
			},
			want: "Emacs-23.3b-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VersionToName(tt.args.version)

			assert.Equal(t, tt.want, got)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGitRefToStableVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr string
	}{
		{
			name: "empty",
			args: args{
				version: "",
			},
			wantErr: "release: git ref is not stable tagged release: \"\"",
		},
		{
			name: "master",
			args: args{
				version: "master",
			},
			wantErr: "release: git ref is not stable tagged release: " +
				"\"master\"",
		},
		{
			name: "feature",
			args: args{
				version: "feature/native-comp",
			},
			wantErr: "release: git ref is not stable tagged release: " +
				"\"feature/native-comp\"",
		},
		{
			name: "stable",
			args: args{
				version: "emacs-27.2",
			},
			want: "27.2",
		},
		{
			name: "stable with letter",
			args: args{
				version: "emacs-23.3b",
			},
			want: "23.3b",
		},
		{
			name: "future stable",
			args: args{
				version: "emacs-239.33",
			},
			want: "239.33",
		},
		{
			name: "future stable with letter",
			args: args{
				version: "emacs-239.33c",
			},
			want: "239.33c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GitRefToStableVersion(tt.args.version)

			assert.Equal(t, tt.want, got)

			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
