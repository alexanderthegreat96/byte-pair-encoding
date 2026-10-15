package main

import "fmt"

const (
	VOCABULARY_CAPACITY       int    = 256
	SHOULD_TRAIN_MODEL        bool   = false
	TOKENIZER_DATA_STORE_FILE string = "tokenizer_model_data.json"
	TRAINING_DATA_FILE        string = "train.txt"
)

func main() {
	// training phase
	if SHOULD_TRAIN_MODEL {
		vocabulary := NewVocabulary(VOCABULARY_CAPACITY)

		trainData, err := LoadTrainingData(TRAINING_DATA_FILE)
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
		vocabulary.TrainFromBytes(trainData, 2000)

		err = SaveTokenizer(vocabulary, TOKENIZER_DATA_STORE_FILE)
		if err != nil {
			panic(err)
		}

	}

	// this will return an instance of vocabulary that we can do stuff on
	var loadVocab *Vocabulary
	var err error

	loadVocab, err = LoadTokenizer(TOKENIZER_DATA_STORE_FILE)
	if err != nil {
		panic(err)
	}

	// now, we can make use of all of this stuff
	inputData := "This is just some sample stuff I am typing to test the vocabulary"

	encoded := loadVocab.Encode(inputData)
	decoded := loadVocab.Decode(encoded)

	fmt.Printf("Used Data: [%s]\n", inputData)
	fmt.Printf("Encoded Tokens: %v\n", encoded)
	fmt.Printf("Decoded Tokens: %v\n", decoded)

	for _, id := range encoded {
		fmt.Printf("Token ID: %d -> %q\n", id, loadVocab.GetTokenAtId(id))
	}
}
