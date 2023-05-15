package normalizers

import (
	"fmt"
	"github.com/ahron-maslin/ravkav-sdk-go/card/parsers"
	"github.com/ahron-maslin/ravkav-sdk-go/contracts"
)

type eventNormalizer struct{}

func NewEventNormalizer(_ contracts.CardOutput) contracts.Normalizer {
	return &eventNormalizer{}
}

func (n *eventNormalizer) Normalize(record contracts.Record, recordIndex int) (map[string]interface{}, error) {
	stream := parsers.NewStreamReader(record.Binary())

	versionNumber, err := stream.Read(3)
	serviceProvider, err := stream.Read(8)
	contractPointer, err := stream.Read(4)
	transportType, err := stream.Read(4)
	eventAction, err := stream.Read(4)
	dateTimeStamp, err := stream.Read(30) // Parse as calypso date
	journeyInterchanges, err := stream.Read(1)
	dateTimeFirstStamp, err := stream.Read(30) // Parse as calypso date
	bestContractPriorities, err := n.bestContractPriorities(stream)
	locationRead, err := stream.Read(7)
	location := parsers.Hex2Int64("00" + locationRead)
	if err != nil {
		return nil, fmt.Errorf("error reading event bits | %s", err)
	}
	place, err := stream.BitConditionRead(location, 0, 16)
	line, err := stream.BitConditionRead(location, 1, 16)
	candidateContracts, err := stream.BitConditionRead(location, 2, 8)
	runId10, err := stream.BitConditionRead(location, 3, 10)
	runId2, err := stream.BitConditionRead(location, 3, 2)
	device, err := stream.BitConditionRead(location, 4, 14)
	if stream.IsBitOn(location, 5) {
		stream.SkipBits(4)
	}
	interchangeRights, err := stream.BitConditionRead(location, 6, 8)
	extensionRead, err := stream.Read(3)
	extension := parsers.Hex2Int64("00" + extensionRead)
	if err != nil {
		return nil, fmt.Errorf("error converting event extension to int | %s", err)
	}
	ticketRoutesSystem, err := stream.BitConditionRead(extension, 0, 10)
	ticketAccumulatedFareCode, err := stream.BitConditionRead(extension, 0, 8)
	ticketDebitAmount, err := stream.BitConditionRead(extension, 0, 16)
	passengersNumber, err := stream.BitConditionRead(extension, 1, 5)
	if err != nil {
		return nil, fmt.Errorf("error reading bit conditional event bits | %s", err)
	}

	eventParser := parsers.NewEvent()
	return map[string]interface{}{
		"versionNumber":             parsers.ParseEn1545Number(versionNumber),
		"serviceProvider":           parsers.Operator(serviceProvider),
		"contractPointer":           contractPointer, // Check how can be connected to contract
		"transportType":             eventParser.Transport(parsers.ParseEn1545Number(transportType)),
		"eventAction":               eventParser.Action(parsers.ParseEn1545Number(eventAction)),
		"dateTimeStamp":             eventParser.TimeRealDate(parsers.ParseEn1545Number(dateTimeStamp)),
		"journeyInterchanges":       journeyInterchanges,
		"dateTimeFirstStamp":        eventParser.TimeRealDate(parsers.ParseEn1545Number(dateTimeFirstStamp)),
		"bestContractPriorities":    bestContractPriorities,
		"place":                     parsers.ParseEn1545Number(place),
		"line":                      parsers.ParseEn1545Number(line),
		"candidateContracts":        candidateContracts,
		"runId10":                   runId10,
		"runId2":                    runId2,
		"device":                    parsers.ParseEn1545Number(device),
		"interchangeRights":         interchangeRights,
		"ticketRoutesSystem":        parsers.ParseEn1545Number(ticketRoutesSystem),
		"ticketAccumulatedFareCode": ticketAccumulatedFareCode,
		"ticketDebitAmount":         ticketDebitAmount,
		"passengersNumber":          passengersNumber,
	}, nil
}

func (n *eventNormalizer) bestContractPriorities(stream *parsers.StreamReader) ([]string, error) {
	var bestContractPriorities []string
	for i := 0; i < 8; i++ {
		contractPriority, err := stream.Read(4)
		if err != nil {
			return nil, fmt.Errorf("error reading contract priorities bits | %s", err)
		}
		bestContractPriorities = append(bestContractPriorities, parsers.ParseEn1545Number(contractPriority))
	}
	return bestContractPriorities, nil
}
