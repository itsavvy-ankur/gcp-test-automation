package stepImpl

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	. "github.com/getgauge-contrib/gauge-go/testsuit"
	"google.golang.org/api/iterator"
)

var vowels map[rune]bool

var _ = gauge.Step("Publish message <message> to pubsub topic <topic>", func(message string, topic string) {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "asitech-dev"
	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	for topics := client.Topics(ctx); ; {
		topic, err := topics.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}

		log.Printf("Topic ID : %s", topic.ID())
	}

})

var _ = gauge.Step("Sample message from file <file>", func(content string) {
	gauge.WriteMessage("Content : %s", content)
	
})

var _ = gauge.Step("Vowels in English language are <vowels>.", func(vowelString string) {
	vowels = make(map[rune]bool, 0)
	for _, ch := range vowelString {
		vowels[ch] = true
	}
})

var _ = gauge.Step("Almost all words have vowels <table>", func(tbl *m.Table) {
	for _, row := range tbl.Rows {
		word := row.Cells[0]
		expectedCount, err := strconv.Atoi(row.Cells[1])
		if err != nil {
			T.Fail(fmt.Errorf("Failed to parse string %s to integer", row.Cells[1]))
		}
		actualCount := countVowels(word)
		if actualCount != expectedCount {
			T.Fail(fmt.Errorf("Vowel count in word %s - got: %d, want: %d", word, actualCount, expectedCount))
		}
	}
})

var _ = gauge.Step("The word <word> has <expectedCount> vowels.", func(word string, expected string) {
	actualCount := countVowels(word)
	expectedCount, err := strconv.Atoi(expected)
	if err != nil {
		T.Fail(fmt.Errorf("Failed to parse string %s to integer", expected))
	}
	if actualCount != expectedCount {
		T.Fail(fmt.Errorf("got: %d, want: %d", actualCount, expectedCount))
	}
})

func countVowels(word string) int {
	vowelCount := 0
	for _, ch := range word {
		if _, ok := vowels[ch]; ok {
			vowelCount++
		}
	}
	return vowelCount
}
