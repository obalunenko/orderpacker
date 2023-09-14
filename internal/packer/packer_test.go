package packer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPacker_PackOrder(t *testing.T) {
	type args struct {
		items uint
	}
	tests := []struct {
		name string
		args args
		want []uint
	}{
		{
			name: "1 - 1x250",
			args: args{
				items: 1,
			},
			want: []uint{250},
		},
		{
			name: "251 - 1x500",
			args: args{
				items: 251,
			},
			want: []uint{500},
		},
		{
			name: "501 - 1x500, 1x250",
			args: args{
				items: 501,
			},
			want: []uint{500, 250},
		},
		{
			name: "12001  - 2x5000, 1x2000, 1x250",
			args: args{
				items: 12001,
			},
			want: []uint{5000, 5000, 2000, 250},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPacker(WithDefaultBoxes())
			require.NoError(t, err)

			got := p.PackOrder(tt.args.items)

			compareSlices(t, tt.want, got)
		})
	}
}

func compareSlices(t *testing.T, expected, actual []uint) {
	bexp, err := json.Marshal(expected)
	require.NoError(t, err)

	bact, err := json.Marshal(actual)
	require.NoError(t, err)

	assert.Equal(t, string(bexp), string(bact))
}

func TestNewPacker(t *testing.T) {
	type args struct {
		opts []PackerOption
	}
	tests := []struct {
		name    string
		args    args
		want    *Packer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "default boxes",
			args: args{
				opts: []PackerOption{},
			},
			want: &Packer{
				boxes: DefaultBoxes,
			},
			wantErr: assert.NoError,
		},
		{
			name: "custom boxes",
			args: args{
				opts: []PackerOption{
					WithBoxes([]uint{32, 1, 2, 2, 4, 16, 8, 16}),
				},
			},
			want: &Packer{
				boxes: []uint{1, 2, 4, 8, 16, 32},
			},
			wantErr: assert.NoError,
		},
		{
			name: "custom boxes empty - error",
			args: args{
				opts: []PackerOption{
					WithBoxes([]uint{}),
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPacker(tt.args.opts...)
			if !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
