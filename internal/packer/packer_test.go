package packer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacker_PackOrder(t *testing.T) {
	type args struct {
		items int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "1 - 1x250",
			args: args{
				items: 1,
			},
			want: []int{250},
		},
		{
			name: "251 - 1x500",
			args: args{
				items: 251,
			},
			want: []int{500},
		},
		{
			name: "501 - 1x500, 1x250",
			args: args{
				items: 501,
			},
			want: []int{500, 250},
		},
		{
			name: "12001  - 2x5000, 1x2000, 1x250",
			args: args{
				items: 1,
			},
			want: []int{5000, 5000, 2000, 250},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewDefaultPacker()

			got := p.PackOrder(tt.args.items)

			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
