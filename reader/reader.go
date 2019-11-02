package reader

import (
	"errors"
	"github.com/ebfe/scard"
	"ravkav-sdk-go/contracts"
)

type Reader struct {
	card    *scard.Card
	context *scard.Context
}

func NewReader() contracts.Reader {
	return &Reader{}
}

func (r *Reader) ListReaders() ([]string, error) {
	// Establish a PC/SC context
	context, err := scard.EstablishContext()
	if err != nil {
		return nil, err
	}
	r.context = context

	// List available readers
	readers, err := r.context.ListReaders()
	if err != nil {
		return nil, err
	}

	if len(readers) == 0 {
		return nil, errors.New("error selecting reader")
	}

	return readers, nil
}

func (r *Reader) Connect(reader string) error {
	// Connect to the card
	card, err := r.context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		return err
	}
	r.card = card

	return nil
}

func (r *Reader) Disconnect() error {
	err := r.card.Disconnect(scard.LeaveCard)
	if err != nil {
		return err
	}

	err = r.context.Release()
	if err != nil {
		return err
	}
	return nil
}

func (r *Reader) Transmit(command []byte) ([]byte, error) {
	return r.card.Transmit(command)
}
