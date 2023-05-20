package lzw

import "testing"

// ================================================ test ================================================ //

func TestCompress12bitArray(t *testing.T) {
    input := []uint16{0x0000, 0x0111, 0x0222, 0x0333, 0x0444, 0x0555, 0x0666, 0x0777, 0x0888, 0x0999, 0x0aaa, 0x0bbb, 0x0ccc, 0x0ddd, 0x0eee, 0x0fff}
    expected := []uint16{0x0001, 0x1122, 0x2333, 0x4445, 0x5566, 0x6777, 0x8889, 0x99aa, 0xabbb, 0xcccd, 0xddee, 0xefff}
    output := compress12bitArray(input)

    if len(output) != len(expected) {
        t.Errorf("Something went wrong\nexpected: %+v\noutput: %+v", expected, output)
    }

    for i := 0; i < len(expected); i++ {
        if expected[i] != output[i] {
            t.Errorf("Something went wrong\nexpected: %+v\noutput: %+v", expected, output)
        }
    }
}
