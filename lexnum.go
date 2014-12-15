package lexnum

import (
	"fmt"
	"strconv"
)

type Encoder struct {
	pos rune
	neg rune
}

func NewEncoder(pos rune, neg rune) *Encoder {
	if pos < neg {
		panic("positive lexnum rune must be of higher rank than negative lexnum rune")
	}
	return &Encoder{pos: pos, neg: neg}
}

func (l Encoder) EncodeInt(i int) string {
	if i == 0 {
		return "0"
	}
	if i > 0 {
		return l.encodePos(i)
	}
	return l.encodeNeg(i)
}

func (l Encoder) encodePos(i int) string {
	s := strconv.Itoa(i)
	if len(s) == 1 {
		return fmt.Sprintf("%c%s", l.pos, s)
	}
	return fmt.Sprintf("%c%s%s", l.pos, l.encodePos(len(s)), s)
}

func (l Encoder) encodeNeg(i int) string {
	if i < 0 {
		i = -i
	}
	runes := []rune(strconv.Itoa(i))
	for i := range runes {
		runes[i] = l.flip(runes[i])
	}
	if len(runes) == 1 {
		return fmt.Sprintf("%c%s", l.neg, string(runes))
	}
	return fmt.Sprintf("%c%s%s", l.neg, l.encodeNeg(len(runes)), string(runes))
}

func (l Encoder) flip(r rune) rune {
	switch r {
	case '0':
		return '9'
	case '1':
		return '8'
	case '2':
		return '7'
	case '3':
		return '6'
	case '4':
		return '5'
	case '5':
		return '4'
	case '6':
		return '3'
	case '7':
		return '2'
	case '8':
		return '1'
	case '9':
		return '0'
	default:
		panic(fmt.Sprintf("can't flip illegal rune %c", r))
	}
}

func (l Encoder) flipInPlace(runes []rune) {
	for i := range runes {
		runes[i] = l.flip(runes[i])
	}
}

func (l Encoder) isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func (l Encoder) prefixCount(runes []rune) int {
	i := 0
	for _, r := range runes {
		if r == l.neg || r == l.pos {
			i++
		} else {
			break
		}
	}
	return i
}

func (l Encoder) DecodeInt(s string) (int, error) {
	if s == "" {
		return 0, fmt.Errorf("illegal Lexnum decode of empty string")
	}
	runes := []rune(s)
	if len(runes) == 1 {
		if runes[0] == '0' {
			return 0, nil
		}
		return 0, fmt.Errorf("illegal Lexnum decode of non-zero unit string: %s", s)
	}
	switch runes[0] {
	case l.neg:
		return l.decodeNeg(runes)
	case l.pos:
		return l.decodePos(runes)
	default:
		return 0, fmt.Errorf("illegal Lexnum decode of string without %c or %c as initial rune: %s", l.neg, l.pos, s)
	}
}

func (l Encoder) decodePos(runes []rune) (int, error) {
	return l._decodePos(runes, 1, l.prefixCount(runes))
}

func (l Encoder) _decodePos(runes []rune, size int, index int) (int, error) {
	n, err := strconv.ParseInt(string(runes[index:index+size]), 10, 64)
	if err != nil {
		return 0, err
	}
	if index+size > len(runes) {
		return 0, fmt.Errorf("illegal Lexnum decode of abnormally long string %s", string(runes))
	}
	if index+size == len(runes) {
		return int(n), nil
	}
	return l._decodePos(runes, int(n), index+size)

}

func (l Encoder) decodeNeg(runes []rune) (int, error) {
	p := l.prefixCount(runes)
	l.flipInPlace(runes[p:len(runes)])
	n, err := l._decodePos(runes, 1, p)
	if err != nil {
		return 0, err
	}
	return -n, nil
}
