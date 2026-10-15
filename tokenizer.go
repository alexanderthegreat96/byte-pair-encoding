package main

// Byte-Pair Encoding Implementation in the Go Programming Language
// Based on Meta.ai and OpenAI's algorithm

// Algorithm implementation is split in 2 phases:
// A -> Vocabulary Initialization
// B -> Vocabulary Training -> finding the best pairs

// I could use 2 types of data -> a slice of rune or a slice of bytes
// I chose to go with bytes, as bytes are a lot more accurate and closer
// to the metal and I don't want runes to hide explicit behavior
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

func (v *Vocabulary) GetTokenAtId(tokenId int) []byte {
	if tokenId > len(v.data) {
		return nil
	}

	return v.data[tokenId]
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

func (v *Vocabulary) TrainFromBytes(inputStream []byte, maxVocabularySize int) {
	integerSlice := make([]int, 0, len(inputStream))

	for _, b := range inputStream {
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

func (v *Vocabulary) Train(inputData string, maxVocabularySize int) {
	if len(inputData) < 1 {
		return
	}

	stringToBytes := []byte(inputData)
	integerSlice := make([]int, 0, len(stringToBytes))

	for _, b := range stringToBytes {
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
