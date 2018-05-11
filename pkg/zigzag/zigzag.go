package zigzag

// ZigZag32 encodes an int32 to a uint32
func ZigZag32(i int32) uint32 {
	ui := uint32(i << 1)
	if i < 0 {
		ui = ^ui
	}

	return ui
}

// ZagZig32 decodes an int32 from a uint32
func ZagZig32(ui uint32) int32 {
	i := int32(ui >> 1)
	if ui&1 != 0 {
		i = ^i
	}

	return i
}

// ZigZag encodes an int64 to a uint64
func ZigZag(i int64) uint64 {
	ui := uint64(i << 1)
	if i < 0 {
		ui = ^ui
	}

	return ui
}

// ZagZig decodes an int64 from a uint64
func ZagZig(ui uint64) int64 {
	i := int64(ui >> 1)
	if ui&1 != 0 {
		i = ^i
	}

	return i
}
