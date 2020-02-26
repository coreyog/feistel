package feistel

import (
	"crypto/sha256"
	"reflect"
	"testing"
)

func TestTransform(t *testing.T) {
	input := []byte("ABCXYZ")

	t.Logf(" input = %x\n", input)
	t.Logf("string = %s\n", string(input))

	keys := [][]byte{
		{50, 100, 150},
		{125},
		{123, 234},
	}

	result, err := Transform(input, len(keys), func(data []byte, round int) (result []byte) {
		h := sha256.Sum256(append(data, keys[round]...))
		return h[:]
	})

	if err != nil {
		panic(err)
	}

	t.Logf("result = %x\n", result)
	t.Logf("string = %s\n", string(result))

	if reflect.DeepEqual(input, result) {
		t.Log("expected: input != result, actual: input == result")
		t.Fail()
	}

	result, err = Transform(result, len(keys), func(data []byte, round int) (result []byte) {
		h := sha256.Sum256(append(data, keys[len(keys)-round-1]...))
		return h[:]
	})

	if err != nil {
		panic(err)
	}

	t.Logf("result = %x\n", result)
	t.Logf("string = %s\n", string(result))

	if !reflect.DeepEqual(input, result) {
		t.Log("expected: input == result, actual: input != result")
		t.Fail()
	}
}

func TestUnbalanced(t *testing.T) {
	msg := []byte{1, 2, 3}

	ciphertext, err := Transform(msg, 3, func(data []byte, i int) []byte {
		return data
	})

	if ciphertext != nil {
		t.Log("expected: ciphertext == nil, actual, ciphertext != nil")
		t.Fail()
	}

	if err != ErrUnbalancedInput {
		t.Log("expected: err == ErrUnbalancedInput, actual: err != ErrUnbalancedInput")
		t.Fail()
	}
}

func TestBadRoundCount(t *testing.T) {
	msg := []byte{1, 2, 3, 4}

	ciphertext, err := Transform(msg, -1, func(data []byte, i int) []byte {
		return data
	})

	if ciphertext != nil {
		t.Log("expected: ciphertext == nil, actual, ciphertext != nil")
		t.Fail()
	}

	if err != ErrInvalidRoundCount {
		t.Log("expected: err == ErrInvalidRoundCount, actual: err != ErrInvalidRoundCount")
		t.Fail()
	}
}

func TestInvalidF(t *testing.T) {
	msg := []byte{1, 2, 3, 4}

	ciphertext, err := Transform(msg, 3, nil)

	if ciphertext != nil {
		t.Log("expected: ciphertext == nil, actual, ciphertext != nil")
		t.Fail()
	}

	if err != ErrInvalidF {
		t.Log("expected: err == ErrInvalidF, actual: err != ErrInvalidF")
		t.Fail()
	}
}

func TestEmptyF(t *testing.T) {
	msg := []byte{1, 2, 3, 4}

	// nil
	ciphertext, err := Transform(msg, 3, func(data []byte, i int) []byte {
		return nil
	})

	if ciphertext != nil {
		t.Log("expected: ciphertext == nil, actual, ciphertext != nil")
		t.Fail()
	}

	if err != ErrFReturnedEmpty {
		t.Log("expected: err == ErrFReturnedEmpty, actual: err != ErrFReturnedEmpty")
		t.Fail()
	}

	// empty
	ciphertext, err = Transform(msg, 3, func(data []byte, i int) []byte {
		return []byte{}
	})

	if ciphertext != nil {
		t.Log("expected: ciphertext == nil, actual, ciphertext != nil")
		t.Fail()
	}

	if err != ErrFReturnedEmpty {
		t.Log("expected: err == ErrFReturnedEmpty, actual: err != ErrFReturnedEmpty")
		t.Fail()
	}
}
