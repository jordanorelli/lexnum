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
	{1, "=1"},
	{9, "=9"},
	{10, "==210"},
	{99, "==299"},
	{100, "==3100"},
	{12345, "==512345"},
	{123456789, "==9123456789"},
	{1234567890, "===2101234567890"},
	{-1, "-8"},
	{-2, "-7"},
	{-9, "-0"},
	{-10, "--789"},
	{-11, "--788"},
	{-123, "--6876"},
	{-123456789, "--0876543210"},
	{-1234567890, "---7898765432109"},
}

func TestLexnum(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	e := NewEncoder('=', '-')
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

	// -150 to 150 test
	nums, stringz := make([]int, 0, 301), make([]string, 0, 301)
	for x := -150; x <= 150; x += 1 {
		nums = append(nums, x)
		stringz = append(stringz, e.EncodeInt(x))
	}
	sort.Strings(stringz)
	sort.Ints(nums)
	for i := 0; i < len(nums); i++ {
		n, err := e.DecodeInt(stringz[i])
		if err != nil {
			t.Errorf("unable to decode our own input: %v", stringz[i])
			continue
		}
		if n != nums[i] {
			t.Errorf("sorting is broken in range test")
			t.Log(stringz, "\n", nums)
			break
		}
	}

	// random test
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
