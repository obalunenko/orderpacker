package packer

type Packer struct {
	boxes []uint
}

var defaultBoxes = []uint{
	250,
	500,
	1000,
	2000,
	5000,
}

type PackerOption func(*Packer)

func WithBoxes(boxes []uint) PackerOption {
	return func(p *Packer) {
		p.boxes = boxes
	}
}

func WithDefaultBoxes() PackerOption {
	return func(p *Packer) {
		p.boxes = defaultBoxes
	}
}

func NewPacker(opts ...PackerOption) *Packer {
	p := &Packer{}

	if len(opts) == 0 {
		opts = []PackerOption{WithDefaultBoxes()}
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
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
