package ravkav

import (
	"ravkav-sdk-go/card"
	"ravkav-sdk-go/contracts"
	"ravkav-sdk-go/reader"
)

func main() {}

func NewReader() contracts.Reader {
	return reader.NewReader()
}

func NewCard(reader contracts.Reader) contracts.Card {
	return card.NewByReader(reader)
}
