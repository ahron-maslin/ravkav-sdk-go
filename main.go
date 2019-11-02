package main

import (
	"go-ravkav/card"
	"go-ravkav/contracts"
	"go-ravkav/reader"
)

func main() {}

func NewReader() contracts.Reader {
	return reader.NewReader()
}

func NewCard(reader contracts.Reader) contracts.Card {
	return card.NewByReader(reader)
}
