package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_toAPIResponse(t *testing.T) {
	type args struct {
		boxes []uint
	}
	tests := []struct {
		name string
		args args
		want PackResponse
	}{
		{
			name: "[500, 500, 500]",
			args: args{
				boxes: []uint{500, 500, 500},
			},
			want: PackResponse{
				Packs: []Pack{
					{
						Box:      500,
						Quantity: 3,
					},
				},
			},
		},
		{
			name: "[500, 2000, 500]",
			args: args{
				boxes: []uint{500, 2000, 500},
			},
			want: PackResponse{
				Packs: []Pack{
					{
						Box:      2000,
						Quantity: 1,
					},
					{
						Box:      500,
						Quantity: 2,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toAPIResponse(tt.args.boxes)

			assert.Equal(t, tt.want, got)
		})
	}
}
