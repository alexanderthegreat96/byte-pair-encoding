package main

import (
	"fmt"

	"github.com/alexanderthegreat96/byte-pair-encoding/tokenizer"
)

const (
	SHOULD_TRAIN_MODEL                bool   = false // flip this to train first
	VOCABULARY_CAPACITY               int    = 256   // default to 256 bytes
	TOKENIZER_DATA_STORE_FILE         string = "tokenizer_model_data.json"
	TRAINING_DATA_FILE                string = "train.txt"
	TRAINING_DATA_MAX_VOCABULARY_SIZE int    = 2000 // you should change this in case you have a larger file than what I have
)

func main() {
	// training phase
	if SHOULD_TRAIN_MODEL {
		vocabulary := tokenizer.NewVocabulary(VOCABULARY_CAPACITY)

		// unix tokens for example
		specialTokens := []string{
			"[PAD]",
			"[UNK]",
			"<|endoftext|>",
		}

		for _, token := range specialTokens {
			vocabulary.AddSpecialTokenFromString(token)
		}

		trainData, err := tokenizer.LoadTrainingData(TRAINING_DATA_FILE)
		if err != nil {
			panic(err)
		}

		// the core idea here is that we want the tokenizer to learn the stuff properly
		// which would result in something that is meaningful
		// now, we are training on a lot of text
		// so if we were to give it 3-400, that's not enough
		// probably the sweet spot for this would be between 1 and 3k
		// so I chose 2k
		// if we go above this, it will not produce meaningful results
		vocabulary.TrainFromBytes(trainData, TRAINING_DATA_MAX_VOCABULARY_SIZE)

		err = tokenizer.SaveTokenizer(vocabulary, TOKENIZER_DATA_STORE_FILE)
		if err != nil {
			panic(err)
		}

	}

	// this will return an instance of vocabulary that we can do stuff on
	var loadVocab *tokenizer.Vocabulary
	var err error

	loadVocab, err = tokenizer.LoadTokenizer(TOKENIZER_DATA_STORE_FILE)
	if err != nil {
		panic(err)
	}

	// now, we can make use of all of this stuff
	inputData := "This is just some sample stuff I am typing to test the vocabulary"

	encoded := loadVocab.Encode(inputData)
	decoded := loadVocab.Decode(encoded)
	tokenized := loadVocab.FormatTokens(encoded)

	fmt.Printf("Input Data:\n%s\n", inputData)
	fmt.Printf("Encoded Tokens:\n%v\n", encoded)
	fmt.Printf("Decoded Tokens:\n%v\n", decoded)
	fmt.Printf("Tokenized Data:\n%v\n", tokenized)

	for _, id := range encoded {
		fmt.Printf("Token ID: %d -> %q\n", id, loadVocab.GetTokenAtId(id))
	}
}
