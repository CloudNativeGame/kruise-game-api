package builder

type BytesStack struct {
	stack []byte
}

func (bs *BytesStack) AppendByte(b byte) {
	bs.stack = append(bs.stack, b)
}

func (bs *BytesStack) PopByte() byte {
	c := bs.stack[len(bs.stack)-1]
	bs.stack = bs.stack[:len(bs.stack)-1]
	return c
}

func (bs *BytesStack) TopByte() byte {
	if len(bs.stack) == 0 {
		return 0
	}
	return bs.stack[len(bs.stack)-1]
}

func (bs *BytesStack) IsEmpty() bool {
	return len(bs.stack) == 0
}
