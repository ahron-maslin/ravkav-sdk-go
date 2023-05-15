package normalizers

import (
	"fmt"
	"github.com/ahron-maslin/ravkav-sdk-go/card/parsers"
	"github.com/ahron-maslin/ravkav-sdk-go/contracts"
	"strconv"
)

type contractNormalizer struct {
	cardOutput contracts.CardOutput
}

func NewContractNormalizer(cardOutput contracts.CardOutput) contracts.Normalizer {
	return &contractNormalizer{cardOutput}
}

func (n *contractNormalizer) Normalize(record contracts.Record, recordIndex int) (map[string]interface{}, error) {
	counters := n.cardOutput.GetMeta()["counters"]
	counter := counters[strconv.Itoa(recordIndex+1)]

	stream := parsers.NewStreamReader(record.Binary())

	versionNumber, err := stream.Read(3)
	validityStartDate, err := stream.Read(14)
	issuer, err := stream.Read(8)
	accessType, err := stream.Read(2)
	counterType, err := stream.Read(3)
	etta, err := stream.Read(6)
	saleDate, err := stream.Read(14)
	saleDevice, err := stream.Read(12)
	saleNumberDaily, err := stream.Read(10)
	journeyInterchanges, err := stream.Read(1)
	validityInfoRead, err := stream.Read(9)
	validityInfo := parsers.Hex2Int64("00" + validityInfoRead)
	if err != nil {
		return nil, fmt.Errorf("error reading contract bits | %s", err)
	}
	restrictTimeCode, err := stream.BitConditionRead(validityInfo, 0, 5)
	restrictCode, err := stream.BitConditionRead(validityInfo, 1, 5)
	restrictDuration, err := stream.BitConditionRead(validityInfo, 2, 6)
	validityEndDate, err := stream.BitConditionRead(validityInfo, 3, 14)
	validityDurationType, err := stream.BitConditionRead(validityInfo, 4, 2)
	validityDurationValue, err := stream.BitConditionRead(validityInfo, 4, 6)
	periodJourneysType, err := stream.BitConditionRead(validityInfo, 5, 2)
	periodJourneysValue, err := stream.BitConditionRead(validityInfo, 5, 6)
	customerProfile, err := stream.BitConditionRead(validityInfo, 6, 6)
	passengersNumber, err := stream.BitConditionRead(validityInfo, 7, 5)
	if err != nil {
		return nil, fmt.Errorf("error reading bit conditional contract bits | %s", err)
	}
	if stream.IsBitOn(validityInfo, 8) {
		stream.SkipBits(32)
	}

	var validityLocations = parsers.NewValidityLocation(stream)
	for {
		spatialTypeRead, err := stream.Read(4)
		if err != nil {
			return nil, fmt.Errorf("error reading spatial type from stream | %s", err)
		}
		spatialType := parsers.Hex2Int64(spatialTypeRead)
		if spatialType != 15 {
			err := validityLocations.Parse(spatialType)
			if err != nil {
				return nil, fmt.Errorf("error parsing validity locations | %s", err)
			}
		} else if stream.BitsLeft() < 8 {
			return nil, fmt.Errorf("less than 8 bits left")
		} else {
			stream.SkipBits(stream.BitsLeft() - 8)

			contractParser := parsers.NewContract()
			validityStartDate, err := contractParser.InvertedDate(validityStartDate)
			if err != nil {
				return nil, fmt.Errorf("error parsing validityStartDate | %s", err)
			}
			validityEndDate, err := contractParser.InvertedDate(validityEndDate)
			if err != nil {
				return nil, fmt.Errorf("error parsing validityEndDate | %s", err)
			}
			authenticator, err := stream.Read(8)
			if err != nil {
				return nil, fmt.Errorf("error reading authenticator | %s", err)
			}
			return map[string]interface{}{
				"contractNumber":      recordIndex + 1,
				"versionNumber":       parsers.ParseEn1545Number(versionNumber),
				"issuer":              parsers.Operator(issuer),
				"journeyInterchanges": parsers.ParseEn1545Number(journeyInterchanges),
				"customerProfile":     customerProfile,
				"passengersNumber":    passengersNumber,
				"authenticator":       authenticator,
				"etta": map[string]interface{}{
					"code": parsers.ParseEn1545Number(etta),
					"name": contractParser.Etta(parsers.ParseEn1545Number(etta)),
				},
				"validity": map[string]interface{}{
					"startDate": validityStartDate,
					"endDate":   contractParser.CalculateValidityEndDate(validityStartDate, validityEndDate, validityDurationType, parsers.ParseEn1545Number(validityDurationValue)),
					"locations": validityLocations.Locations(),
				},
				"counter": map[string]interface{}{
					"type":  parsers.ParseEn1545Number(counterType),
					"name":  contractParser.CounterType(parsers.ParseEn1545Number(counterType)),
					"value": contractParser.CounterValue(parsers.ParseEn1545Number(counterType), counter),
					"date":  contractParser.CounterDate(parsers.ParseEn1545Number(counterType), counter),
				},
				"access": map[string]interface{}{
					"type": parsers.ParseEn1545Number(accessType),
					"name": contractParser.AccessType(parsers.ParseEn1545Number(accessType)),
				},
				"sale": map[string]interface{}{
					"date":        parsers.ParseEn1545Date(saleDate),
					"device":      parsers.ParseEn1545Number(saleDevice),
					"numberDaily": parsers.ParseEn1545Number(saleNumberDaily),
				},
				"restrict": map[string]interface{}{
					"code": parsers.ParseEn1545Number(restrictCode),
					"name": contractParser.RestrictType(parsers.ParseEn1545Number(restrictCode)),
					"duration": map[string]interface{}{
						"code": parsers.ParseEn1545Number(restrictDuration),
						"name": contractParser.RestrictDuration(parsers.ParseEn1545Number(restrictDuration)),
					},
					"time": map[string]interface{}{
						"code": parsers.ParseEn1545Number(restrictTimeCode),
						"name": contractParser.RestrictTime(parsers.ParseEn1545Number(restrictTimeCode)),
					},
					"interchangeTimeInMinutes": contractParser.InterchangeTimeInMinutes(parsers.ParseEn1545Number(restrictCode), parsers.ParseEn1545Number(restrictDuration)),
				},
				"period": map[string]interface{}{
					"journeysType":  parsers.ParseEn1545Number(periodJourneysType),
					"journeysValue": periodJourneysValue,
				},
			}, nil
		}
	}
}
