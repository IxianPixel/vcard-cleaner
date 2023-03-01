package main

import (
	"github.com/emersion/go-vcard"
	"io"
	"log"
	"os"
	"strings"
)

func parseIcloudCard() ([]vcard.Card, int) {
	f, err := os.Open("cards.vcf")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	personalCards := make([]vcard.Card, 0)

	outlookCounter := 0
	personalCounter := 0
	dec := vcard.NewDecoder(f)
	for {
		card, err := dec.Decode()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		note := card.PreferredValue(vcard.FieldNote)

		if strings.Contains(note, "This contact is read-only") {
			outlookCounter++
		} else {
			personalCounter++
			personalCards = append(personalCards, card)
		}
	}

	return personalCards, outlookCounter
}

func writeOutputCard(cards []vcard.Card) {
	f, err := os.Create("output.vcf")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	enc := vcard.NewEncoder(f)

	for _, card := range cards {
		enc.Encode(card)
	}
}

func main() {
	personalCards, outlookCards := parseIcloudCard()
	writeOutputCard(personalCards)

	log.Printf("Personal: %d, Outlook: %d (%d)", len(personalCards), outlookCards, (len(personalCards) + outlookCards))
}