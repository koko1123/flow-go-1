package request

import (
	"fmt"
	"math"
	"strconv"
)

const sealed = "sealed"
const final = "final"

// Special height values
const SealedHeight = math.MaxUint64 - 1
const FinalHeight = math.MaxUint64 - 2
const EmptyHeight = math.MaxUint64 - 3

type Height uint64

func (h *Height) Parse(raw string) error {
	if raw == "" { // allow empty
		*h = EmptyHeight
		return nil
	}

	if raw == sealed {
		*h = SealedHeight
		return nil
	}
	if raw == final {
		*h = FinalHeight
		return nil
	}

	height, err := strconv.ParseUint(raw, 0, 64)
	if err != nil {
		return fmt.Errorf("invalid height format")
	}

	*h = Height(height)
	return nil
}

func (h Height) Flow() uint64 {
	return uint64(h)
}

type Heights []Height

func (h *Heights) Parse(raw []string) error {
	var height Height
	heights := make([]Height, 0)
	for _, r := range raw {
		err := height.Parse(r)
		if err != nil {
			return err
		}
		// don't include empty heights
		if height == EmptyHeight {
			continue
		}

		heights = append(heights, height)
	}

	*h = heights
	return nil
}

func (h Heights) Flow() []uint64 {
	heights := make([]uint64, len(h))
	for i, he := range h {
		heights[i] = he.Flow()
	}
	return heights
}
