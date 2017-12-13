package chunk

func zigzag(i int32) uint32 {
	ui := uint32(i << 1)
	if i < 0 {
		ui = ^ui
	}

	return ui
}

func zagzig(ui uint32) int32 {
	i := int32(ui >> 1)
	if ui&1 != 0 {
		i = ^i
	}

	return i
}
