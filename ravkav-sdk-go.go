package ravkavsdkgo

import (
	"github.com/derkinderfietsen/ravkav-sdk-go/card"
	"github.com/derkinderfietsen/ravkav-sdk-go/contracts"
	"github.com/derkinderfietsen/ravkav-sdk-go/reader"
)

func main() {}

// NewReader creates a new reader instance used for connecting to
// the card and transmit commands for getting card's content.
func NewReader() contracts.Reader {
	return reader.NewReader()
}

// NewCard creates a new card instance, it receives a reader instance
// for flexibility purposes and for tests.
func NewCard(reader contracts.Reader) contracts.Card {
	return card.NewByReader(reader)
}
