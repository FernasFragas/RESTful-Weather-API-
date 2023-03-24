package chat

import (
	"github.com/jdkato/prose/v2"
	"log"
)

func GetLocationFromChat(userInput string) string {
	doc, err := prose.NewDocument(userInput)
	if err != nil {
		log.Fatal(err.Error())
	}

	var location string

	for _, entity := range doc.Tokens() {
		if entity.Label == "B-GPE" {
			location = entity.Text
			break
		}
	}
	return location
}
