# Feistel Cipher

After watching [Computerphile's video](https://www.youtube.com/watch?v=FGhj3CGxl8I) on Feistel Networks (aka Feistel Chains or Feistel Ciphers) I wanted to build my own.

## Example

```
input := []byte("ABCXYZ")

// filler HMAC secrets
keys := [][]byte{
  {50, 100, 150},
  {125},
  {123, 234},
}

// Encrypt
result, err := Transform(input, len(keys), func(data []byte, round int) (result []byte) {
  h := sha256.Sum256(append(data, keys[round]...))
  return h[:]
})

// Decrypt
result, err = Transform(result, len(keys), func(data []byte, round int) (result []byte) {
  h := sha256.Sum256(append(data, keys[len(keys)-round-1]...))
  return h[:]
})
```

## Quick Explaination

A feistel cipher splits it's input in equal parts usually named `left` and `right`. `right` is run through a round function `F`, xored with `left`, that result is stored in `left`, finally `left` and `right` are swapped. This is called a Round. The process runs for multiple rounds and different keys should be used in each Round by `F`.

<p align="center">
 <img style="background-color: white;" src="https://upload.wikimedia.org/wikipedia/commons/f/fa/Feistel_cipher_diagram_en.svg" />
</p>

`F` can be just about anything. Cryptographic hash function, another cipher, any process that returns a non-zero length of bytes to xor against `left`.

Because the encryption process is the same as the decryption process, you just need to Transform the encrypted data to decrypt it. The only gotcha is that when decrypting you should use the keys in the reverse order of the encryption process.

This library expects the input to contain an even number of bytes. If supplied with an odd length input, then the returned byte slice will be nil and the err will be `ErrUnbalancedInput`.

You provide the number of rounds and the round function. If a round count less than 2 is provided, then the returned byte slice will be nil and the err will be `ErrInvalidRoundCount`. If the round function returns a `nil` or empty slice, then the returned byte slice will be nil and the err will be `ErrFReturnedEmpty`. If the round function itself is `nil`, then the returned byte slice will be nil and the err will be `ErrInvalidF`.

The xor process used in this library will xor one []byte with a second []byte, repeating the second []byte if it's shorter than the first. If the second []byte is longer than the first then a portion of the second goes unused.