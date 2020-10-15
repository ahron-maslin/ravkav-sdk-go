package parsers

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const (
	minutesInDay      int = 1440
	minutesInHalfHour int = 30
	minutesInMonth    int = 43200
	minutesInWeek     int = 10080
)

const (
	AccessTypeParkingAndTransport              string = "3"
	AccessTypeParkingOnly                      string = "2"
	AccessTypeTransportOnly                    string = "1"
	CounterTypeDateAndRemainingJourneys        string = "1"
	CounterTypeMonetaryAmount                  string = "3"
	CounterTypeNotUsed                         string = "0"
	CounterTypeNumberOfTokens                  string = "2"
	DurationTypeDays                           int    = 2
	DurationTypeHalfHours                      int    = 3
	DurationTypeMonths                         int    = 0
	DurationTypeWeeks                          int    = 1
	EttAFreeCertificate                        string = "4"
	EttAParkAndRideSeasonPass                  string = "5"
	EttASeasonPass                             string = "2"
	EttASingleOrMultiRide                      string = "1"
	EttASingleOrMultiRideExtension             string = "7"
	EttAStoredValue                            string = "6"
	EttATransferTicket                         string = "3"
	RestrictCodeNoRestriction                  string = "0"
	RestrictCodeQuaysOnly                      string = "2"
	RestrictCodeRestrictDuration5MinutesUnits  string = "16"
	RestrictCodeStationAreaOnly                string = "1"
	RestrictDurationEndOfExploitation          string = "62"
	RestrictDurationEndOfService               string = "63"
	RestrictDurationNoRestriction              string = "0"
	RestrictTimeNotValidEarlyMorningHours      string = "1"
	RestrictTimeNotValidPeakHours              string = "2"
	RestrictTimeNotValidWeekendDaysAndHolidays string = "4"
	RestrictTimeNotValidWorkDays               string = "3"
	RestrictTimeNoRestriction                  string = "0"
	RestrictTimeValidOnlySchoolDays            string = "5"
)

type Contract struct{}

func NewContract() *Contract {
	return &Contract{}
}

func (c *Contract) CalculateValidityEndDate(validityStartDate, validityEndDate, validityDurationType, validityDurationValue string) string {
	if validityEndDate != "" && validityEndDate != "0" {
		return validityEndDate
	}
	if validityDurationType == "" || validityDurationValue == "" {
		return ""
	}
	validityDurationTypeInt, err := strconv.Atoi(validityDurationType)
	if err != nil {
		return ""
	}
	validityDurationValueInt, err := strconv.Atoi(validityDurationValue)
	if err != nil {
		return ""
	}
	validityDurationInMinutes, err := c.getValidityDurationInMinutes(validityDurationTypeInt, validityDurationValueInt)
	if validityDurationInMinutes == 0 || validityDurationInMinutes%minutesInDay != 0 || err != nil {
		return ""
	}
	var validityDurationInDays int = validityDurationInMinutes / minutesInDay
	validityStartDateInt, err := strconv.Atoi(validityStartDate)
	if err != nil {
		return ""
	}
	calculatedTime := time.Unix(int64(validityStartDateInt/1000), 0)
	calculatedTime = calculatedTime.AddDate(0, 0, validityDurationInDays-1)
	return strconv.Itoa(int(calculatedTime.UnixNano() / 1000000))
}

func (c *Contract) getValidityDurationInMinutes(validityDurationType, validityDurationValue int) (int, error) {
	switch validityDurationType {
	case DurationTypeMonths:
		return minutesInMonth * validityDurationValue, nil
	case DurationTypeWeeks:
		return validityDurationValue * minutesInWeek, nil
	case DurationTypeDays:
		return validityDurationValue * minutesInDay, nil
	case DurationTypeHalfHours:
		return validityDurationValue * minutesInHalfHour, nil
	default:
		return 0, fmt.Errorf("invalid validityDurationType: %d", validityDurationType)
	}
}

func (c *Contract) InterchangeTimeInMinutes(restrictCode, restrictDuration string) string {
	if restrictDuration == "" || restrictDuration == RestrictDurationNoRestriction || restrictDuration == RestrictDurationEndOfExploitation || restrictDuration == RestrictDurationEndOfService {
		return ""
	}
	restrictDurationInt, err := strconv.Atoi(restrictDuration)
	if err != nil {
		return ""
	}
	if restrictCode != RestrictCodeRestrictDuration5MinutesUnits {
		return strconv.Itoa(restrictDurationInt * minutesInHalfHour)
	}
	return strconv.Itoa(restrictDurationInt * 5)
}

func (c *Contract) CounterValue(counterType string, counter interface{}) string {
	rawCounter := counter.(int)
	switch counterType {
	case CounterTypeDateAndRemainingJourneys:
		return strconv.Itoa(rawCounter & 1023)
	case CounterTypeNumberOfTokens, CounterTypeMonetaryAmount:
		return strconv.Itoa(rawCounter)
	default:
		return ""
	}
}

func (c *Contract) CounterDate(counterType string, counter interface{}) string {
	rawCounter := counter.(int)
	switch counterType {
	case CounterTypeDateAndRemainingJourneys:
		counterDate, err := c.InvertedDate(strconv.Itoa(rawCounter >> 10))
		if err != nil {
			return ""
		}
		return counterDate
	default:
		return ""
	}
}

func (c *Contract) AccessType(typeNumber string) string {
	var types map[string]string = map[string]string{
		AccessTypeTransportOnly:       "TRANSPORT_ONLY",
		AccessTypeParkingOnly:         "PARKING_ONLY",
		AccessTypeParkingAndTransport: "PARKING_AND_TRANSPORT",
	}
	if typeVal, ok := types[typeNumber]; ok {
		return typeVal
	}
	return ""
}

func (c *Contract) CounterType(typeNumber string) string {
	var types map[string]string = map[string]string{
		CounterTypeNotUsed:                  "NOT_USED",
		CounterTypeDateAndRemainingJourneys: "NUMBER_OF_TOKENS",
		CounterTypeNumberOfTokens:           "DATE_AND_REMAINING_JOURNEYS",
		CounterTypeMonetaryAmount:           "MONETARY_AMOUNT",
	}
	if typeVal, ok := types[typeNumber]; ok {
		return typeVal
	}
	return ""
}

func (c *Contract) RestrictType(typeNumber string) string {
	var types map[string]string = map[string]string{
		RestrictCodeNoRestriction:                 "NO_RESTRICTION",
		RestrictCodeQuaysOnly:                     "QUAYS_ONLY",
		RestrictCodeRestrictDuration5MinutesUnits: "RESTRICT_DURATION_5_MINUTES_UNITS",
		RestrictCodeStationAreaOnly:               "STATION_AREA_ONLY",
	}
	if typeVal, ok := types[typeNumber]; ok {
		return typeVal
	}
	return ""
}

func (c *Contract) RestrictDuration(code string) string {
	var durations map[string]string = map[string]string{
		RestrictDurationEndOfExploitation: "END_OF_EXPLOITATION",
		RestrictDurationEndOfService:      "END_OF_SERVICE",
		RestrictDurationNoRestriction:     "NO_RESTRICTION",
	}
	if durationVal, ok := durations[code]; ok {
		return durationVal
	}
	return ""
}

func (c *Contract) RestrictTime(code string) string {
	var times map[string]string = map[string]string{
		RestrictTimeNotValidEarlyMorningHours:      "NOT_VALID_EARLY_MORNING_HOURS",
		RestrictTimeNotValidPeakHours:              "NOT_VALID_PEAK_HOURS",
		RestrictTimeNotValidWeekendDaysAndHolidays: "NOT_VALID_WEEKEND_DAYS_AND_HOLIDAYS",
		RestrictTimeNotValidWorkDays:               "NOT_VALID_WORK_DAYS",
		RestrictTimeNoRestriction:                  "NO_RESTRICTION",
		RestrictTimeValidOnlySchoolDays:            "VALID_ONLY_SCHOOL_DAYS",
	}
	if durationVal, ok := times[code]; ok {
		return durationVal
	}
	return ""
}

func (c *Contract) Etta(etta string) string {
	var ettas map[string]string = map[string]string{
		EttAFreeCertificate:            "FREE_CERTIFICATE",
		EttAParkAndRideSeasonPass:      "PARK_AND_RIDE_SEASON_PASS",
		EttASeasonPass:                 "SEASON_PASS",
		EttASingleOrMultiRide:          "SINGLE_OR_MULTI_RIDE",
		EttASingleOrMultiRideExtension: "SINGLE_OR_MULTI_RIDE_EXTENSION",
		EttAStoredValue:                "STORED_VALUE",
		EttATransferTicket:             "TRANSFER_TICKET",
	}
	if ettaVal, ok := ettas[etta]; ok {
		return ettaVal
	}
	return ""
}

func (c *Contract) InvertedDate(hexString string) (string, error) {
	if hexString == "" || hexString == "0" {
		return hexString, nil
	}

	byteVal, err := hex.DecodeString(hexString)
	if err != nil {
		return "", fmt.Errorf("error parsing inverted date hex | %s", err)
	}
	binVal := BytesToBin(byteVal)
	binVal = binVal[2:]
	stream := NewStreamReader(binVal)
	hexVal, err := stream.Read(len(binVal))
	if err != nil {
		return "", fmt.Errorf("error reading inverted date | %s", err)
	}
	intVal := Hex2Int64(hexVal)
	xorVal := intVal ^ 0x3fff
	dateVal := 852069600 + xorVal*24*3600
	return strconv.Itoa(dateVal * 1000), nil
}
