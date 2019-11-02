package ravkavsdkgo

import (
	"github.com/ybaruchel/ravkav-sdk-go/card"
	"github.com/ybaruchel/ravkav-sdk-go/contracts"
	"github.com/ybaruchel/ravkav-sdk-go/reader"
)

func main() {}

func NewReader() contracts.Reader {
	return reader.NewReader()
}

func NewCard(reader contracts.Reader) contracts.Card {
	return card.NewByReader(reader)
}
