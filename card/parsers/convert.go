package parsers

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// Parse for parsing response from transmitions
func ParseCardResponse(byteResponse []byte) (binaryParse string, hexParse string) {
	return BytesToBin(byteResponse), hex.EncodeToString(byteResponse)
}

func BytesToBin(resp []byte) string {
	var response string
	for _, n := range resp {
		response += fmt.Sprintf("%08b", n)
	}
	return response
}

func BinToHex(s string) (string, error) {
	ui, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", ui), err
}

func Hex2Int64(hexStr string) int64 {
	// base 16 for hexadecimal
	result, _ := strconv.ParseUint(hexStr, 16, 64)
	return int64(result)
}
