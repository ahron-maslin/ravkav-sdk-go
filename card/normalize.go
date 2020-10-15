package card

import (
	"fmt"
	"github.com/derkinderfietsen/ravkav-sdk-go/card/normalizers"
	"github.com/derkinderfietsen/ravkav-sdk-go/contracts"
)

var normalizersArr map[string]func(_ contracts.CardOutput) contracts.Normalizer = map[string]func(_ contracts.CardOutput) contracts.Normalizer{
	"application":    normalizers.NewApplicationNormalizer,
	"environment":    normalizers.NewEnvironmentNormalizer,
	"contracts":      normalizers.NewContractNormalizer,
	"counters":       normalizers.NewCounterNormalizer,
	"event_logs":     normalizers.NewEventNormalizer,
	"special_events": normalizers.NewEventNormalizer,
}

// Normalizes card fields
func (c *Card) Normalize() error {
	for metaName, record := range (*c).GetMeta() {
		err := c.normalizeRecord(c, record, metaName, 0)
		if err != nil {
			return err
		}
	}

	for recordName, records := range (*c).GetRecords() {
		for index, record := range records {
			err := c.normalizeRecord(c, record, recordName, index)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Card) normalizeRecord(card *Card, record contracts.Record, name string, recordIndex int) error {
	normalizer, err := c.newNormalizer(name)
	if err != nil {
		return fmt.Errorf("error normalizing record %s[%d] | %s", name, recordIndex, err)
	}
	normalizedRecord, err := normalizer.Normalize(record, recordIndex)
	if err != nil {
		return fmt.Errorf("error normalizing record %s[%d] | %s", name, recordIndex, err)
	}
	record.SetNormalized(normalizedRecord)
	return nil
}

func (c *Card) newNormalizer(nodeName string) (contracts.Normalizer, error) {
	if normalizer, ok := normalizersArr[nodeName]; ok {
		return normalizer(c.Output()), nil
	}

	return nil, fmt.Errorf("normalizer for %s not found", nodeName)
}
