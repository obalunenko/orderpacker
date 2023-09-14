package packer

import (
	"fmt"
	"slices"
)

type Packer struct {
	boxes []uint
}

var DefaultBoxes = []uint{
	250,
	500,
	1000,
	2000,
	5000,
}

type PackerOption func(*Packer)

func WithBoxes(boxes []uint) PackerOption {
	return func(p *Packer) {
		slices.Sort(boxes)

		unique := make(map[uint]struct{}, len(boxes))

		isUnique := func(x uint) bool {
			if _, exist := unique[x]; !exist {
				unique[x] = struct{}{}

				return true
			}

			return false
		}

		filtered := boxes[:0]

		for _, b := range boxes {
			if isUnique(b) {
				filtered = append(filtered, b)
			}
		}

		p.boxes = filtered
	}
}

func WithDefaultBoxes() PackerOption {
	return func(p *Packer) {
		p.boxes = DefaultBoxes
	}
}

func NewPacker(opts ...PackerOption) (*Packer, error) {
	var p Packer

	if len(opts) == 0 {
		opts = []PackerOption{WithDefaultBoxes()}
	}

	for _, opt := range opts {
		opt(&p)
	}

	if len(p.boxes) == 0 {
		return nil, fmt.Errorf("boxes list is empty")
	}

	return &p, nil
}

func (p Packer) PackOrder(items uint) []uint {
	if items == 0 {
		return []uint{}
	}

	var result []uint

	for i := len(p.boxes) - 1; i >= 0; i-- {
		box := p.boxes[i]

		if box > items {
			if i == 0 {
				result = append(result, box)

				break
			}

			continue
		}

		if box <= items {
			if i == 0 {
				result = append(result, p.boxes[i+1])

				break
			}

			n := items / box

			if n == 0 {
				result = append(result, box)

				break
			}

			for j := uint(0); j < n; j++ {
				result = append(result, box)
			}

			left := items % box

			if left == 0 {
				break
			}

			items = left
		}
	}

	return result
}
