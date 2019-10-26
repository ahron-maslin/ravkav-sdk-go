package parsers

import "fmt"

type StreamReader struct {
	src string
	cur int
}

func NewStreamReader(src string) *StreamReader {
	return &StreamReader{src: src}
}

func (s *StreamReader) Read(bits int) (string, error) {
	defer func() {
		s.cur += bits
	}()
	if len(s.src) < s.cur+bits {
		return "", fmt.Errorf("bits out of range: %d-%d", s.cur, s.cur+bits)
	}
	hexVal, err := BinToHex(s.src[s.cur : s.cur+bits])
	if err != nil {
		return "", fmt.Errorf("hex formatting failed for bits range %d-%d | %s", s.cur, s.cur+bits, err)
	}
	return hexVal, nil
}

func (s *StreamReader) BitConditionRead(validator int, bit uint, bits int) (string, error) {
	if s.IsBitOn(validator, bit) {
		readBits, err := s.Read(bits)
		if err != nil {
			return "", fmt.Errorf("error reading bit condition bits")
		}
		return readBits, nil
	}
	return "0", nil
}

func (s *StreamReader) IsBitOn(validator int, bit uint) bool {
	return ((1 << bit) & validator) != 0
}
func (s *StreamReader) SkipBits(bits int) {
	s.cur += bits
}

func (s *StreamReader) BitsLeft() int {
	return len(s.src) - s.cur
}
