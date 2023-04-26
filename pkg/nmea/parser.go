package nmea

import (
	"errors"
	"fmt"
	"io"
)

// Tokenize reads a sentence from in.
func Tokenize(in io.Reader) (Sentence, error) {
	result := []string{}

	b := make([]byte, 50)
	var checksum byte

	var mode int
	const modeFindStart = 0
	const modeReadField = 1

	for {
		b = b[:0]
		d, err := readUntil(in, &b)
		if err != nil {
			if errors.Is(err, tokenError) {
				// out of buffer while looking for delimiter
				continue
			}
			return nil, err
		}

		if mode == modeFindStart {
			if d == '$' {
				// start found
				mode = modeReadField
				checksum = 0
			}
		} else {
			if d == ',' || d == '*' {
				// field found
				result = append(result, string(b))
				for _, c := range b {
					checksum ^= c
				}
				if d == ',' {
					checksum ^= d
				}
			}
			if d == '*' {
				// last field, next 2 bytes contain checksum
				b := []byte{0, 0}
				_, err := in.Read(b)
				if err != nil {
					return nil, err
				}
				cs := (b[0]-'0')<<4 + b[1] - '0'
				if checksum != cs {
					return nil, fmt.Errorf("checksum expected %02x, got: %02x", cs, checksum)
				}

				return result, nil
			}
		}
	}
}

// ReadUntil reads from in until a ',' or '*' delimiter and returns the result and the delimiter.
func readUntil(in io.Reader, result *[]byte) (byte, error) {
	b := []byte{0}
	for i := 0; i < cap(*result); i++ {
		_, err := in.Read(b)
		if err != nil {
			return 0, err
		}

		if b[0] == ',' || b[0] == '$' || b[0] == '*' {
			return b[0], nil
		}
		*result = append(*result, b[0])
	}
	return 0, tokenError
}

type tokenErr string

func (te tokenErr) Error() string {
	return string(te)
}

const tokenError = tokenErr("result is full")
