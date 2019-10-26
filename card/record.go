package card

import (
	"bytes"
	"go-ravkav/card/parsers"
	"go-ravkav/contracts"
)

// Record Single record
type Record struct {
	ByteVal          []byte
	BinVal           string
	HexVal           string
	normalizer       contracts.Normalizer
	NormalizedValues map[string]interface{}
}

var emptyResponse []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 144, 0}

func NewRecord(response []byte, command contracts.Command) (contracts.Record, error) {
	binVal, hexVal := parsers.ParseCardResponse(response)

	if bytes.Equal(response, emptyResponse) { // Empty record
		return nil, nil
	}

	return &Record{
		ByteVal:          response,
		BinVal:           binVal,
		HexVal:           hexVal,
		NormalizedValues: make(map[string]interface{}),
	}, nil
}

func (r *Record) Normalized() map[string]interface{} {
	return r.NormalizedValues
}

func (r *Record) NormalizedVal(key string) interface{} {
	return r.NormalizedValues[key]
}

func (r *Record) SetNormalized(values map[string]interface{}) {
	r.NormalizedValues = values
}

func (r *Record) SetNormalizedValue(key string, value string) {
	r.NormalizedValues[key] = value
}

func (r *Record) Bytes() []byte {
	return r.ByteVal
}

func (r *Record) Binary() string {
	return r.BinVal
}

func (r *Record) Hex() string {
	return r.HexVal
}
