package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_ReleaseURL(t *testing.T) {
	type fields struct {
		Type   Type
		Source string
	}
	type args struct {
		releaseName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "empty name",
			fields: fields{Type: GitHub, Source: "foo/bar"},
			args:   args{releaseName: ""},
			want:   "",
		},
		{
			name:   "GitHub, foo/bar, v1.0.0",
			fields: fields{Type: GitHub, Source: "foo/bar"},
			args:   args{releaseName: "v1.0.0"},
			want:   "https://github.com/foo/bar/releases/tag/v1.0.0",
		},
		{
			name:   "Not GitHub, foo/bar, v1.0.0",
			fields: fields{Type: Type("oops"), Source: "foo/bar"},
			args:   args{releaseName: "v1.0.0"},
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Repository{
				Type:   tt.fields.Type,
				Source: tt.fields.Source,
			}

			got := repo.ReleaseURL(tt.args.releaseName)

			assert.Equal(t, tt.want, got)
		})
	}
}
