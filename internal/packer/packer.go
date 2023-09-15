package packer

import (
	"fmt"
	log "log/slog"
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

	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("failed to validate packer: %w", err)
	}

	log.Info("Packer created", "boxes", p.boxes)

	return &p, nil
}

func (p Packer) validate() error {
	if len(p.boxes) == 0 {
		return fmt.Errorf("boxes list is empty")
	}

	// There should be no box with zero volume.
	for _, box := range p.boxes {
		if box == 0 {
			return fmt.Errorf("box with zero volume")
		}
	}

	return nil

}

func (p Packer) PackOrder(items uint) []uint {
	if items == 0 {
		return []uint{}
	}

	if p.boxes[0] == 0 {
		// This should never happen, cause we validate boxes on creation.
		panic(fmt.Errorf("packer has box with zero volume: boxes [%v]", p.boxes))
	}

	// Preallocate memory for the result slice.
	// Make a prediction based on the number of items and the smallest box.
	result := make([]uint, 0, items/p.boxes[0])

	if len(p.boxes) == 1 {
		box := p.boxes[0]

		if items < box {
			return []uint{box}
		}

		n := items / box

		last := items % box
		if last != 0 {
			n++
		}

		for i := uint(0); i < n; i++ {
			result = append(result, box)
		}

		return result
	}

	for i := len(p.boxes) - 1; i >= 0; i-- {
		box := p.boxes[i]

		if box >= items {
			if i == 0 {
				result = append(result, box)

				break
			}

			continue
		}

		if box < items {
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
