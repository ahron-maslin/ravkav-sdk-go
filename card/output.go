package card

import (
	"encoding/json"
	"go-ravkav/contracts"
)

type Output struct {
	Meta    map[string]map[string]interface{}   `json:"meta"`
	Records map[string][]map[string]interface{} `json:"records"`
}

// Returns card normalized output
func (c *Card) Output() contracts.CardOutput {
	metaMap := make(map[string]map[string]interface{})
	for name, record := range c.GetMeta() {
		metaMap[name] = record.Normalized()
	}

	recordsMap := make(map[string][]map[string]interface{})
	for name, records := range c.GetRecords() {
		for _, record := range records {
			recordsMap[name] = append(recordsMap[name], record.Normalized())
		}
	}

	var output contracts.CardOutput = &Output{Meta: metaMap, Records: recordsMap}
	c.output = output

	return c.output
}

func (o *Output) JSON() (string, error) {
	res, err := json.Marshal(o)
	return string(res), err
}
