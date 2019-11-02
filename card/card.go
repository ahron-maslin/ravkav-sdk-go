package card

import (
	"github.com/ybaruchel/ravkav-sdk-go/contracts"
)

type Records struct {
	meta    map[string]contracts.Record
	records map[string][]contracts.Record
}

type Card struct {
	reader  contracts.Reader
	records *Records
	output  contracts.CardOutput
}

func NewByReader(ravkavReader contracts.Reader) contracts.Card {
	var records Records = Records{meta: make(map[string]contracts.Record), records: make(map[string][]contracts.Record)}
	return &Card{reader: ravkavReader, records: &records}
}

func (c *Card) GetRecords() map[string][]contracts.Record {
	return c.records.records
}

func (c *Card) GetRecord(name string) []contracts.Record {
	return c.records.records[name]
}

func (c *Card) GetMeta() map[string]contracts.Record {
	return c.records.meta
}

func (o *Output) GetMeta() map[string]map[string]interface{} {
	return o.Meta
}

func (o *Output) GetRecords() map[string][]map[string]interface{} {
	return o.Records
}
