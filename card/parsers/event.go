package parsers

import (
	"strconv"
)

const (
	actionCodeUnknown               string = "-1"
	actionCodeAction                string = "0"
	actionCodeEntry                 string = "1"
	actionExit                      string = "2"
	actionPassage                   string = "3"
	actionOnBoardControl            string = "4"
	actionTest                      string = "5"
	actionTransferTrip              string = "6"
	actionExitStation               string = "7"
	actionOtherUsage                string = "8"
	actionContractCancellation      string = "9"
	actionContractLoadAndValidation string = "12"
	actionContractLoad              string = "13"
	actionCardIssuance              string = "14"
	actionCardInvalidation          string = "15"
	transportPublicTransportation   string = "0"
	transportBus                    string = "1"
	transportInterurbanBus          string = "2"
	transportIntercity              string = "3"
	transportLightRail              string = "4"
	transportTrain                  string = "5"
	transportTransportMean          string = "15"
)

type Event struct{}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) Action(code string) string {
	var actions map[string]string = map[string]string{
		actionCodeUnknown:               "UNKNOWN",
		actionCodeAction:                "ACTION",
		actionCodeEntry:                 "ENTRY",
		actionExit:                      "EXIT",
		actionPassage:                   "PASSAGE",
		actionOnBoardControl:            "ON_BOARD_CONTROL",
		actionTest:                      "TEST",
		actionTransferTrip:              "TRANSFER_TRIP",
		actionExitStation:               "EXIT_STATION",
		actionOtherUsage:                "OTHER_USAGE",
		actionContractCancellation:      "CONTRACT_CANCELLATION",
		actionContractLoadAndValidation: "CONTRACT_LOAD_AND_VALIDATION",
		actionContractLoad:              "CONTRACT_LOAD",
		actionCardIssuance:              "CARD_ISSUANCE",
		actionCardInvalidation:          "CARD_INVALIDATION",
	}
	if val, ok := actions[code]; ok {
		return val
	}
	return ""
}

func (e *Event) TimeRealDate(secondsSince1997 string) int64 {
	secondsSince1997Int, err := strconv.Atoi(secondsSince1997)
	if err != nil || secondsSince1997Int == 0 {
		return 0
	}
	return int64(secondsSince1997Int)*1000 + int64(852076800000)
}

func (e *Event) Transport(code string) string {
	var transports map[string]string = map[string]string{
		transportPublicTransportation: "PUBLIC_TRANSPORTATION",
		transportBus:                  "BUS",
		transportInterurbanBus:        "INTERURBAN_BUS",
		transportIntercity:            "INTERCITY",
		transportLightRail:            "LIGHT_RAIL",
		transportTrain:                "TRAIN",
		transportTransportMean:        "TRANSPORT_MEAN",
	}
	if val, ok := transports[code]; ok {
		return val
	}
	return ""
}
