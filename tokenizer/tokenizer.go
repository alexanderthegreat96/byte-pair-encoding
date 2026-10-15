package tokenizer

import (
	"strings"
)

// Byte-Pair Encoding Algorithm in the Go Programming Language
// Based on the original 1994 data compression algorithm by Philip Gage.
//
// The algorithm implementation is split into three phases:
//
//	A. Vocabulary Initialization
//	B. Vocabulary Training (finding the best pairs)
//	C. Vocabulary Usage
//
// This implementation uses a slice of bytes ([]byte) rather than a slice of runes ([]rune).
// Bytes are chosen because they are closer to the metal and avoid rune abstractions
// hiding explicit byte-level behavior.

type Vocabulary struct {
	data      map[int][]byte
	merges    map[[2]int]int
	mergeList []mergeRule
}

type mergeRule struct {
	Left  int
	Right int
	NewId int
}

func NewVocabulary(capacity int) *Vocabulary {
	dataMap := make(map[int][]byte)

	for num := range capacity {
		dataMap[num] = []byte{byte(num)}
	}

	return &Vocabulary{
		data:      dataMap,
		merges:    make(map[[2]int]int),
		mergeList: []mergeRule{},
	}
}

func (v *Vocabulary) AddSpecialTokenFromBytes(token []byte) {
	if len(token) < 1 {
		return
	}

	nextId := len(v.data)
	v.data[nextId] = token
}

func (v *Vocabulary) AddSpecialTokenFromString(token string) {
	if len(token) < 1 {
		return
	}

	nextId := len(v.data)
	v.data[nextId] = []byte(token)
}

func (v *Vocabulary) GetTokenAtId(tokenId int) []byte {
	if tokenId > len(v.data) {
		return nil
	}

	return v.data[tokenId]
}

func (v *Vocabulary) GetData() map[int][]byte {
	return v.data
}

func (v *Vocabulary) GetMerges() map[[2]int]int {
	return v.merges
}

// translate a string input to a slice of ints
func (v *Vocabulary) Encode(input string) []int {
	if len(input) < 1 {
		return nil
	}

	stringToBytes := []byte(input)

	// this is a small optimization
	// i am allocating exactly how many bytes I need
	// for this integer slice based on the length of the bytes slice
	integerSlice := make([]int, 0, len(stringToBytes))

	for _, b := range stringToBytes {
		integerSlice = append(integerSlice, int(b))
	}

	// I am keeping track of the merges basically
	// previously, I had a separate pair logic
	// now, i don't need all of that because the mergePairs handles that
	// and we can store the merges inside the Vocabulary struct
	for pair, replacementId := range v.merges {
		integerSlice = mergePairs(integerSlice, pair, replacementId)
	}

	return integerSlice
}

// can be used to verify if the encoding process produces
// good results
func (v *Vocabulary) Decode(tokens []int) string {
	if len(tokens) == 0 {
		return ""
	}

	var fullBytes []byte
	for _, id := range tokens {
		if tokenBytes, exists := v.data[id]; exists {
			fullBytes = append(fullBytes, tokenBytes...)
		}
	}

	return string(fullBytes)
}

// debugging function which outputs a "tokenized" format
// ex: This is just some sample stuff
// output: Th_is_ _is_ _j_us_t_ s_o_me_ s_am_p_le_ s_tu_ff_
func (v *Vocabulary) FormatTokens(inputTokens []int) string {
	if len(inputTokens) < 1 {
		return ""
	}

	stringBuff := make([]string, 0, len(inputTokens))

	for _, tokenId := range inputTokens {
		tokenBytes := v.GetTokenAtId(tokenId)
		if tokenBytes != nil {
			stringBuff = append(stringBuff, string(tokenBytes))
		} else {
			stringBuff = append(stringBuff, "[UNK]")
		}

	}

	return strings.Join(stringBuff, "_")
}

// this should be utilized inside a real scenario
// where we read a lot of data from a file
// therefore, a string conversion from the byte buffer
// will cause performance issues
// and we don't want that
func (v *Vocabulary) TrainFromBytes(inputStream []byte, maxVocabularySize int) {
	v.trainHelper(inputStream, maxVocabularySize)
}

// this is useful for testing and small set scenarios
func (v *Vocabulary) TrainFromString(inputData string, maxVocabularySize int) {
	if len(inputData) < 1 {
		return
	}

	stringToBytes := []byte(inputData)
	v.trainHelper(stringToBytes, maxVocabularySize)
}

func mergePairs(inputSlice []int, targetPair [2]int, replacementId int) []int {
	newSlice := make([]int, 0, len(inputSlice))

	for i := 0; i < len(inputSlice); {
		if i < len(inputSlice)-1 && inputSlice[i] == targetPair[0] && inputSlice[i+1] == targetPair[1] {
			newSlice = append(newSlice, replacementId)
			i += 2
		} else {
			newSlice = append(newSlice, inputSlice[i])
			i += 1
		}
	}
	return newSlice
}

// this is what does the magic of training
// nothing crazy though, but it definately does it's job
func (v *Vocabulary) trainHelper(inputBytes []byte, maxVocabularySize int) {
	if len(inputBytes) < 1 {
		return
	}

	integerSlice := make([]int, 0, len(inputBytes))

	for _, b := range inputBytes {
		integerSlice = append(integerSlice, int(b))
	}

	for len(v.data) < maxVocabularySize {
		if len(integerSlice) < 12 {
			break
		}

		freqCounter := make(map[[2]int]int)

		var pair [2]int
		for i := 0; i < len(integerSlice)-1; i++ {
			pair[0] = integerSlice[i]
			pair[1] = integerSlice[i+1]
			freqCounter[pair]++
		}

		if len(freqCounter) == 0 {
			break
		}

		var bestPair [2]int
		var maxCount int = -1
		for p, count := range freqCounter {
			if count > maxCount {
				maxCount = count
				bestPair = p
			}
		}

		newTokenId := len(v.data)

		leftBytes := v.data[bestPair[0]]
		rightBytes := v.data[bestPair[1]]

		combinedBytes := append([]byte{}, leftBytes...)
		combinedBytes = append(combinedBytes, rightBytes...)

		v.data[newTokenId] = combinedBytes
		v.merges[bestPair] = newTokenId
		v.mergeList = append(v.mergeList, mergeRule{
			Left:  bestPair[0],
			Right: bestPair[1],
			NewId: newTokenId,
		})

		integerSlice = mergePairs(integerSlice, bestPair, newTokenId)
	}
}
