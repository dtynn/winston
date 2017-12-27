package storage

import (
	"bytes"
	"testing"
)

func TestRange(t *testing.T) {
	t.Run("Prefix", func(t *testing.T) {
		cases := []struct {
			prefix []byte
			end    []byte
		}{
			{
				nil,
				nil,
			},
			{
				[]byte{},
				nil,
			},
			{
				[]byte{'_', 'a', '_'},
				[]byte{'_', 'a', '_' + 1},
			},
			{
				[]byte{0xff, 0xff, 0xff},
				[]byte{0xff, 0xff, 0xff, 0x00},
			},
		}

		for i, c := range cases {
			end := PrefixEnd(c.prefix)
			if !bytes.Equal(c.end, end) {
				t.Errorf("#%d expected %v, got %v", i+1, c.end, end)
			}

			if !KeyInRange(c.prefix, c.prefix, end) {
				t.Errorf("#%d not in key range", i+1)
			}

			if len(end) == len(c.prefix) {
				for j := 0; j <= 0xff; j++ {
					key := append(c.prefix, byte(j))
					if !KeyInRange(key, c.prefix, end) {
						t.Fatalf("#%d not in key range for append byte %d", i+1, j)
					}
				}
			} else {
				if KeyInRange(append(c.prefix, 0), c.prefix, end) {
					t.Fatalf("#%d should not in key range for bytes full of 0xff", i+1)
				}
			}
		}
	})

	t.Run("Range", func(t *testing.T) {
		cases := []struct {
			start    []byte
			end      []byte
			key      []byte
			expected bool
		}{
			{
				nil,
				nil,
				[]byte{'a'},
				true,
			},
			{
				[]byte{'a'},
				nil,
				[]byte{'a'},
				true,
			},
			{
				[]byte{'b'},
				nil,
				[]byte{'a'},
				false,
			},
			{
				nil,
				[]byte{'a'},
				[]byte{'a'},
				false,
			},
			{
				[]byte{'b'},
				[]byte{'b'},
				[]byte{'b'},
				false,
			},
		}

		for i, c := range cases {
			got := KeyInRange(c.key, c.start, c.end)
			if got != c.expected {
				t.Errorf("#%d expected %v, got %v", i+1, c.expected, got)
			}
		}
	})
}
