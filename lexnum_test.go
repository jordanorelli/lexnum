package lexnum

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

var lexNumTests = []struct {
	in  int
	out string
}{
	{0, "0"},
	{1, "x1"},
	{9, "x9"},
	{10, "xx210"},
	{99, "xx299"},
	{100, "xx3100"},
	{12345, "xx512345"},
	{123456789, "xx9123456789"},
	{1234567890, "xxx2101234567890"},
	{-1, "o8"},
	{-2, "o7"},
	{-9, "o0"},
	{-10, "oo789"},
	{-11, "oo788"},
	{-123, "oo6876"},
	{-123456789, "oo0876543210"},
	{-1234567890, "ooo7898765432109"},
}

func TestLexnum(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	e := NewEncoder('x', 'o')
	for _, test := range lexNumTests {
		s := e.EncodeInt(test.in)
		t.Logf("%d -> %s", test.in, test.out)
		if s != test.out {
			t.Errorf("encode input: %v, output: %v, expected output: %v", test.in, s, test.out)
		}
		n, err := e.DecodeInt(s)
		if err != nil {
			t.Errorf("decode error: %v", err)
			continue
		}
		if n != test.in {
			t.Errorf("decode input: %v, output: %v, expected output: %v", s, n, test.in)
		}
	}

	runsize := 8
	for runz := 0; runz < 4; runz += 1 {
		nums, stringz := make([]int, runsize), make([]string, runsize)
		for i := 0; i < runsize; i++ {
			nums[i] = rand.Int()
			if rand.Int()%2 == 0 {
				nums[i] = -nums[i]
			}
			stringz[i] = e.EncodeInt(nums[i])
		}
		sort.Strings(stringz)
		sort.Ints(nums)
		t.Logf("stringz sorted: %v", stringz)
		t.Logf("nums sorted: %v", nums)
		for i := 0; i < runsize; i++ {
			n, err := e.DecodeInt(stringz[i])
			if err != nil {
				t.Errorf("unable to decode our own input: %v", stringz[i])
				continue
			}
			if n != nums[i] {
				t.Errorf("sorting is broken")
			}
		}
	}
}
