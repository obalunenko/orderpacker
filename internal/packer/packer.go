package packer

type Packer struct {
	Boxes []int
}

var boxes = []int{
	250,
	500,
	1000,
	2000,
	5000,
}

func NewPacker(boxes []int) *Packer {
	return &Packer{Boxes: boxes}
}

func NewDefaultPacker() *Packer {
	return &Packer{Boxes: boxes}
}

func (p Packer) PackOrder(items int) []int {
	return []int{
		items,
	}
}
