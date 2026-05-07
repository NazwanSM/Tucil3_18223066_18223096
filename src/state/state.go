package state

type State struct {
	X         int
	Y         int
	NextDigit int
}

func (s State) Key() uint64 {
	const m24 = uint64(1)<<24 - 1
	const m16 = uint64(1)<<16 - 1
	return (uint64(s.NextDigit)&m16)<<48 |
		(uint64(s.Y)&m24)<<24 |
		(uint64(s.X) & m24)
}

func (s State) Equals(other State) bool {
	return s == other 
}