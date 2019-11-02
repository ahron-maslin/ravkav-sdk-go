package normalizers

import (
	"fmt"
	"math"
	"ravkav-sdk-go/card/parsers"
	"ravkav-sdk-go/contracts"
	"strconv"
)

type counterNormalizer struct{}

func NewCounterNormalizer(_ contracts.CardOutput) contracts.Normalizer {
	return &counterNormalizer{}
}

func (n *counterNormalizer) Normalize(record contracts.Record, recordIndex int) (map[string]interface{}, error) {
	recordBin := record.Binary()
	stream := parsers.NewStreamReader(recordBin)

	maxCounters := int(math.Floor(float64(len(recordBin[:len(recordBin)-16]) / 24)))
	var counters map[string]interface{} = make(map[string]interface{})
	for i := 1; i < maxCounters; i++ {
		nCnt, err := stream.Read(24)
		if err != nil {
			return nil, fmt.Errorf("error reading counter from stream | %s", err)
		}
		nVal, err := strconv.Atoi(parsers.ParseEn1545Number(nCnt))
		if err != nil {
			return nil, fmt.Errorf("error normalizing counter[%d] | %s", i, err)
		}
		counters[strconv.Itoa(i)] = nVal
	}

	return counters, nil
}
