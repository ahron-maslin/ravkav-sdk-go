package parsers

import (
	"strconv"
	"time"
)

func ParseEn1545Number(hexString string) string {
	return strconv.Itoa(Hex2Int64(hexString))
}

func ParseEn1545Date(hexString string) string {
	const epoch int64 = 852076800000
	return strconv.Itoa(Hex2Int64(hexString)*1000*60*60*24 + epoch)
}

func ParseConcatenatedDate(date string) string {
	if len(date) != 8 {
		return ""
	}
	timeDate, err := time.Parse("2006-01-02", date[0:4]+"-"+date[4:6]+"-"+date[6:8])
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(timeDate.UnixNano() / 1000000)) // To millisecond date
}
