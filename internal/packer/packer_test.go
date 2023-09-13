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
			p := NewPacker(WithDefaultBoxes())

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

func TestDiv(t *testing.T) {
	var (
		a uint = 10
		b uint = 3
	)

	c := a / b
	d := a % b

	t.Logf("%d / %d = %d", a, b, c)
	t.Logf("%d %% %d = %d", a, b, d)
}
