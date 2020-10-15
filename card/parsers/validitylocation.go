package parsers

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ValidityLocation struct {
	locations []map[string]string
	stream    *StreamReader
	handler   func() error
}

func NewValidityLocation(stream *StreamReader) *ValidityLocation {
	return &ValidityLocation{stream: stream}
}

func (v *ValidityLocation) Parse(spatialType int64) error {
	var spatialTypeHandlers map[int64]func() error = map[int64]func() error{
		0:  v.Zones,
		1:  v.FareCode,
		2:  v.LinesList,
		3:  v.Ride,
		4:  v.RideAndRunType,
		5:  v.RideAndRunId,
		6:  v.RideRunAndSeat,
		7:  v.RideZones,
		8:  v.Parking,
		9:  v.PredefinedContractVL,
		10: v.RouteSystemList,
		11: v.FareCodeExtension,
		14: v.RfuKnownSize,
	}
	validityLocationsHandler, ok := spatialTypeHandlers[spatialType]
	if !ok {
		return fmt.Errorf("missing handler for spatial type: %d", spatialType)
	}
	return validityLocationsHandler()
}

func (v *ValidityLocation) Handle(spatialType int) error {
	return v.handler()
}

func (v *ValidityLocation) Locations() []map[string]string {
	return v.locations
}

func (v *ValidityLocation) PredefinedContractVL() error {
	ettb, err := v.stream.Read(3)
	if err != nil {
		return fmt.Errorf("error parsing predefined contract vl [ettb] | %s", err)
	}
	shareCode, err := v.stream.Read(11)
	if err != nil {
		return fmt.Errorf("error parsing predefined contract vl [shareCode] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType": "PredefinedContractVL",
		"ettb":         strconv.FormatInt(Hex2Int64(ettb), 10),
		"shareCode":    strconv.FormatInt(Hex2Int64(shareCode), 10),
	})
	return nil
}

func (v *ValidityLocation) LinesList() error {
	numberOfLinesHex, err := v.stream.Read(4)
	if err != nil {
		return fmt.Errorf("error parsing lines list [numberOfLinesHex] | %s", err)
	}
	numberOfLines := ParseEn1545Number(numberOfLinesHex)
	numberOfLinesInt, err := strconv.Atoi(numberOfLines)
	if err != nil {
		return fmt.Errorf("error parsing lines list [numberOfLinesInt] | %s", err)
	}

	var spatialLines []string
	for i := 0; i < numberOfLinesInt; i++ {
		spatialLine, err := v.stream.Read(16)
		if err != nil {
			return fmt.Errorf("error parsing lines list [spatialLine][%d] | %s", i, err)
		}
		spatialLines = append(spatialLines, spatialLine)
	}

	spatialLinesJson, err := json.Marshal(spatialLines)
	if err != nil {
		return fmt.Errorf("error parsing lines list [spatialLinesJson] | %s", err)
	}

	v.locations = append(v.locations, map[string]string{
		"validityType": "LinesList",
		"spatialLines": string(spatialLinesJson),
	})
	return nil
}

func (v *ValidityLocation) Zones() error {
	spatialRoutesSystem, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing zones [spatialRoutesSystem] | %s", err)
	}
	spatialZones, err := v.stream.Read(12)
	if err != nil {
		return fmt.Errorf("error parsing zones [spatialZones] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":        "Zones",
		"spatialRoutesSystem": spatialRoutesSystem,
		"spatialZones":        spatialZones,
	})
	return nil
}

func (v *ValidityLocation) FareCode() error {
	spatialRoutesSystem, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing fare code [spatialRoutesSystem] | %s", err)
	}
	spatialFareCode, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing fare code [spatialZones] | %s", err)
	}

	var fareCodes map[string]string = map[string]string{
		"3": "ABOVE_TO_SPECIFIED_INCLUDED",
		"0": "ANY",
		"1": "ONLY_SPECIFIED",
		"4": "PREFERRED_UNLESS_OTHER_SPECIFIED",
		"2": "UP_TO_SPECIFIED_INCLUDED",
	}

	v.locations = append(v.locations, map[string]string{
		"validityType":        "FareCode",
		"spatialRoutesSystem": spatialRoutesSystem,
		"spatialFareCode":     spatialFareCode,
		"spatialFareCodeName": fareCodes[spatialFareCode],
	})
	return nil
}

func (v *ValidityLocation) Ride() error {
	spatialLine, err := v.stream.Read(16)
	if err != nil {
		return fmt.Errorf("error parsing ride [spatialLine] | %s", err)
	}
	spatialStationOrigin, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride [spatialStationOrigin] | %s", err)
	}
	spatialStationDestination, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride [spatialStationDestination] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":              "Ride",
		"spatialLine":               spatialLine,
		"spatialStationOrigin":      spatialStationOrigin,
		"spatialStationDestination": spatialStationDestination,
	})
	return nil
}

func (v *ValidityLocation) RideAndRunType() error {
	spatialLine, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing ride and run type [spatialLine] | %s", err)
	}
	spatialStationOrigin, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride and run type [spatialStationOrigin] | %s", err)
	}
	spatialStationDestination, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride and run type [spatialStationDestination] | %s", err)
	}
	spatialRunType, err := v.stream.Read(4)
	if err != nil {
		return fmt.Errorf("error parsing ride and run type [spatialRunType] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":              "RideAndRunType",
		"spatialLine":               spatialLine,
		"spatialStationOrigin":      spatialStationOrigin,
		"spatialStationDestination": spatialStationDestination,
		"spatialRunType":            spatialRunType,
	})
	return nil
}

func (v *ValidityLocation) RideAndRunId() error {
	spatialLine, err := v.stream.Read(16)
	if err != nil {
		return fmt.Errorf("error parsing ride and run id [spatialLine] | %s", err)
	}
	spatialStationOrigin, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride and run id [spatialStationOrigin] | %s", err)
	}
	spatialStationDestination, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride and run id [spatialStationDestination] | %s", err)
	}
	spatialRunId, err := v.stream.Read(12)
	if err != nil {
		return fmt.Errorf("error parsing ride and run id [spatialRunId] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":              "RideAndRunId",
		"spatialLine":               spatialLine,
		"spatialStationOrigin":      spatialStationOrigin,
		"spatialStationDestination": spatialStationDestination,
		"spatialRunId":              spatialRunId,
	})
	return nil
}

func (v *ValidityLocation) RideRunAndSeat() error {
	spatialLine, err := v.stream.Read(16)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialLine] | %s", err)
	}
	spatialStationOrigin, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialStationOrigin] | %s", err)
	}
	spatialStationDestination, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialStationDestination] | %s", err)
	}
	spatialRunId, err := v.stream.Read(12)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialRunId] | %s", err)
	}
	spatialVehicleCoach, err := v.stream.Read(4)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialVehicleCoach] | %s", err)
	}
	spatialSeat, err := v.stream.Read(7)
	if err != nil {
		return fmt.Errorf("error parsing ride run and seat [spatialSeat] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":              "RideRunAndSeat",
		"spatialLine":               spatialLine,
		"spatialStationOrigin":      spatialStationOrigin,
		"spatialStationDestination": spatialStationDestination,
		"spatialRunId":              spatialRunId,
		"spatialVehicleCoach":       spatialVehicleCoach,
		"spatialSeat":               spatialSeat,
	})
	return nil
}

func (v *ValidityLocation) RideZones() error {
	spatialRoutesSystemFrom, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing ride zones [spatialRoutesSystemFrom] | %s", err)
	}
	spatialZonesFrom, err := v.stream.Read(12)
	if err != nil {
		return fmt.Errorf("error parsing ride zones [spatialZonesFrom] | %s", err)
	}
	spatialRoutesSystemTo, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing ride zones [spatialRoutesSystemTo] | %s", err)
	}
	spatialZonesTo, err := v.stream.Read(12)
	if err != nil {
		return fmt.Errorf("error parsing ride zones [spatialZonesTo] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":            "RideZones",
		"spatialRoutesSystemFrom": spatialRoutesSystemFrom,
		"spatialZonesFrom":        spatialZonesFrom,
		"spatialRoutesSystemTo":   spatialRoutesSystemTo,
		"spatialZonesTo":          spatialZonesTo,
	})
	return nil
}

func (v *ValidityLocation) Parking() error {
	skipBits, err := v.stream.Read(6)
	if err != nil {
		return fmt.Errorf("error parsing parking [skipBits] | %s", err)
	}
	skipBitsInt, err := strconv.Atoi(skipBits)
	if err != nil {
		return fmt.Errorf("error parsing parking [skipBitsInt] | %s", err)
	}
	v.stream.SkipBits(skipBitsInt + 12)
	return nil
}

func (v *ValidityLocation) RouteSystemList() error {
	routeSystemsNumber, err := v.stream.Read(4)
	if err != nil {
		return fmt.Errorf("error parsing route system list [routeSystemsNumber] | %s", err)
	}
	routeSystemsNumberInt, err := strconv.Atoi(routeSystemsNumber)
	if err != nil {
		return fmt.Errorf("error parsing route system list [routeSystemsNumberInt] | %s", err)
	}

	var routesSystems []string
	for i := 0; i < routeSystemsNumberInt; i++ {
		routeSystem, err := v.stream.Read(10)
		if err != nil {
			return fmt.Errorf("error parsing route system list [routeSystem][%d] | %s", i, err)
		}
		routesSystems = append(routesSystems, routeSystem)
	}

	routesSystemsJson, err := json.Marshal(routesSystems)
	if err != nil {
		return fmt.Errorf("error parsing route system list [routesSystemsJson] | %s", err)
	}

	v.locations = append(v.locations, map[string]string{
		"validityType":  "RouteSystemList",
		"routesSystems": string(routesSystemsJson),
	})
	return nil
}

func (v *ValidityLocation) FareCodeExtension() error {
	spatialRoutesSystem, err := v.stream.Read(10)
	if err != nil {
		return fmt.Errorf("error parsing fare code extension [spatialRoutesSystem] | %s", err)
	}
	fareRestrictionCode, err := v.stream.Read(3)
	if err != nil {
		return fmt.Errorf("error parsing fare code extension [fareRestrictionCode] | %s", err)
	}
	spatialFareCode, err := v.stream.Read(8)
	if err != nil {
		return fmt.Errorf("error parsing fare code extension [spatialFareCode] | %s", err)
	}
	v.locations = append(v.locations, map[string]string{
		"validityType":        "FareCodeExtension",
		"spatialRoutesSystem": spatialRoutesSystem,
		"fareRestrictionCode": fareRestrictionCode,
		"spatialFareCode":     spatialFareCode,
	})
	return nil
}

func (v *ValidityLocation) RfuKnownSize() error {
	skipBits, err := v.stream.Read(6)
	if err != nil {
		return fmt.Errorf("error parsing rfu known size [skipBits] | %s", err)
	}
	skipBitsInt, err := strconv.Atoi(skipBits)
	if err != nil {
		return fmt.Errorf("error parsing rfu known size [skipBitsInt] | %s", err)
	}
	v.stream.SkipBits(skipBitsInt + 12)
	return nil
}
