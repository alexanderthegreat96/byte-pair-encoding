# Pure Go Byte-Pair Encoding (BPE) Tokenizer

A lightweight, pure-Go implementation of **Byte-Pair Encoding (BPE)** tokenization inspired by modern LLM subword tokenizers. Built from scratch with no external dependencies or CGO bindings.

---

## Features

- **Pure Go**: No external dependencies, CGO, or heavy build toolchains required.
- **Byte-Level Foundation**: Initializes from a full ASCII/byte range (`0-255`) to handle any arbitrary string or raw data stream accurately.
- **Dynamic Training**: Iteratively counts byte-pair frequencies and merges the most common pairs up to a custom target vocabulary size.
- **Roundtrip Support**: Encode text into token ID slices and decode token IDs back into their original string payloads seamlessly.

---

## Getting Started

### Installation
Clone or add the file into your Go project:

```bash
go mod init bpe-tokenizer
```

### Usage Example

```go
package main

import "fmt"

func main() {
	// Initialize the vocabulary with 256 base bytes
	v := NewVocabulary()

	// Train on a sample corpus to generate custom merge rules
	corpus := "lower lower lower lowest lowest lowest byte-pair encoding"
	v.Train(corpus, 275)

	// Encode a string into token IDs
	text := "lower"
	tokens := v.Encode(text)
	fmt.Printf("Encoded %q to tokens: %v\n", text, tokens)

	// Decode token IDs back to string
	decoded := v.Decode(tokens)
	fmt.Printf("Decoded back to: %q\n", decoded)
}
```

---

## Running Tests

Execute the unit tests using the standard Go test command:

```bash
go test -v
```

---

## License

This project is open source and available under the MIT License.