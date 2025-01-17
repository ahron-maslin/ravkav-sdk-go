package normalizers

import (
	"encoding/hex"
	"strconv"

	"github.com/ahron-maslin/ravkav-sdk-go/card/parsers"
	"github.com/ahron-maslin/ravkav-sdk-go/contracts"
)

type applicationNormalizer struct{}

func NewApplicationNormalizer(_ contracts.CardOutput) contracts.Normalizer {
	return &applicationNormalizer{}
}

func (n *applicationNormalizer) Normalize(record contracts.Record, recordIndex int) (map[string]interface{}, error) {
	bytes := record.Bytes()
	bytes = bytes[19 : bytes[18]+19]

	var bytesLong []byte = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	for i := range bytes {
		bytesLong[(8-len(bytes))+i] = bytes[i]
	}

	return map[string]interface{}{
		"cardNumber": strconv.FormatInt(parsers.Hex2Int64(hex.EncodeToString(bytesLong)), 10),
	}, nil
}
