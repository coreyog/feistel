package feistel

import "errors"

var (
	ErrUnbalancedInput   = errors.New("unbalanced input")
	ErrFReturnedEmpty    = errors.New("round function returned an empty byte slice")
	ErrInvalidRoundCount = errors.New("the number of rounds must be greater than or equal to 2")
	ErrInvalidF          = errors.New("a round function must be provided")
)

// Transform
func Transform(msg []byte, rounds int, f func([]byte, int) []byte) (transformed []byte, err error) {
	if len(msg)%2 == 1 {
		return nil, ErrUnbalancedInput
	}

	if rounds < 2 {
		return nil, ErrInvalidRoundCount
	}

	if f == nil {
		return nil, ErrInvalidF
	}

	left := msg[:len(msg)/2]
	right := msg[len(msg)/2:]

	for i := 0; i < rounds; i++ {
		// f(R)
		// "f of R"
		fR := f(append(right[:0:0], right...), i)
		if len(fR) == 0 {
			return nil, ErrFReturnedEmpty
		}

		// left = left ^ f(R)
		// "xor f of R"
		xfR := xor(left, fR)

		left, right = right, xfR
	}

	left, right = right, left

	transformed = append(left, right...)

	return transformed, nil
}

// xor xor one byte array with another, if the first array is longer than the second, wrap around the second array
func xor(a []byte, b []byte) (c []byte) {
	c = make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		c[i] = a[i] ^ b[i%len(b)]
	}

	return c
}
