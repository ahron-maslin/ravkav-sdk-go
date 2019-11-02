package main

import (
	"go-ravkav/card"
	"go-ravkav/contracts"
	"go-ravkav/reader"
)

func NewReader() {
	reader.NewReader()
}

func NewCard(reader contracts.Reader) {
	card.NewByReader(reader)
}
