package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid yaml file - loaded config",
			args: args{
				path: filepath.Join("testdata", "config.yaml"),
			},
			want: &Config{
				HTTP: httpConfig{
					Port: "8080",
				},
				Pack: packConfig{
					Boxes: []uint{1, 2, 4, 8, 16, 32},
				},
				Log: logConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "valid json file - loaded config",
			args: args{
				path: filepath.Join("testdata", "config.json"),
			},
			want: &Config{
				HTTP: httpConfig{
					Port: "8080",
				},
				Pack: packConfig{
					Boxes: []uint{1, 2, 4, 8, 16, 32},
				},
				Log: logConfig{
					Level:  "info",
					Format: "json",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "not exist - error",
			args: args{
				path: filepath.Join("testdata", "not_exist.yaml"),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				assert.ErrorIs(t, err, ErrNotExists)

				return assert.Error(t, err, msgAndArgs...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(tt.args.path)
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
