package normalizers

import (
	"fmt"
	"github.com/ybaruchel/ravkav-sdk-go/card/parsers"
	"github.com/ybaruchel/ravkav-sdk-go/card/parsers/dictionaries"
	"github.com/ybaruchel/ravkav-sdk-go/contracts"
	"strconv"
	"time"
)

type environmentNormalizer struct{}

func NewEnvironmentNormalizer(_ contracts.CardOutput) contracts.Normalizer {
	return &environmentNormalizer{}
}

func (n *environmentNormalizer) Normalize(record contracts.Record, recordIndex int) (map[string]interface{}, error) {
	stream := parsers.NewStreamReader(record.Binary())

	envApplicationVersionNumber, err := stream.Read(3)
	envCountryId, err := stream.Read(12)
	envIssuerId, err := stream.Read(8)
	envApplicationNo, err := stream.Read(26)
	envIssuingDate, err := stream.Read(14)
	envEndDate, err := stream.Read(14)
	envPayMethod, err := stream.Read(3)
	holderBirthDate, err := stream.Read(32)
	holderCompany, err := stream.Read(14)
	holderCompanyId, err := stream.Read(30)
	holderNumber, err := stream.Read(30)
	holderProf1Code, err := stream.Read(6)
	holderProf1Date, err := stream.Read(14)
	holderProf2Code, err := stream.Read(6)
	holderProf2Date, err := stream.Read(14)
	holderLanguage, err := stream.Read(2)
	holderRFU, err := stream.Read(4)
	if err != nil {
		return nil, fmt.Errorf("error reading from stream | %s", err)
	}

	return map[string]interface{}{
		"env": map[string]string{
			"applicationVersionNumber": parsers.ParseEn1545Number(envApplicationVersionNumber),
			"country":                  n.country(envCountryId),
			"issuer":                   parsers.ParseEn1545Number(envIssuerId),
			"issuerName":               parsers.Operator(envIssuerId),
			"applicationNo":            parsers.ParseEn1545Number(envApplicationNo),
			"issuingDate":              parsers.ParseEn1545Date(envIssuingDate),
			"endDate":                  parsers.ParseEn1545Date(envEndDate),
			"payMethod":                parsers.ParseEn1545Number(envPayMethod),
		},
		"holder": map[string]interface{}{
			"birthDate":   parsers.ParseConcatenatedDate(holderBirthDate),
			"company":     n.company(holderCompany),
			"companyId":   parsers.ParseEn1545Number(holderCompanyId),
			"number":      parsers.ParseEn1545Number(holderNumber),
			"isAnonymous": n.isAnonymous(holderNumber),
			"profiles": []map[string]string{
				{
					"code":    holderProf1Code,
					"name":    n.profile(holderProf1Code),
					"date":    parsers.ParseEn1545Date(holderProf1Date),
					"isValid": n.holderProfileIsValid(parsers.ParseEn1545Date(holderProf1Date)),
				},
				{
					"code":    holderProf2Code,
					"name":    n.profile(holderProf2Code),
					"date":    parsers.ParseEn1545Date(holderProf2Date),
					"isValid": n.holderProfileIsValid(parsers.ParseEn1545Date(holderProf2Date)),
				},
			},
			"language": parsers.ParseEn1545Number(holderLanguage),
			"RFU":      parsers.ParseEn1545Number(holderRFU),
		},
	}, nil
}

func (n *environmentNormalizer) holderProfileIsValid(date string) string {
	profileExpirationDate, err := strconv.Atoi(date)
	if err != nil {
		return "0"
	}
	var secondsInADay = (24 * time.Hour).Seconds()
	if profileExpirationDate < int(secondsInADay) {
		return "0"
	}
	return "1"
}

func (n *environmentNormalizer) isAnonymous(holderNumber string) string {
	if parsers.ParseEn1545Number(holderNumber) == "0" {
		return "1"
	}
	return "0"
}

func (n *environmentNormalizer) country(code string) string {
	if country, ok := dictionaries.CountryCodes[code]; ok {
		return country
	}
	return ""
}

func (n *environmentNormalizer) profile(code string) string {
	if profile, ok := dictionaries.RavkavProfileCodes[code]; ok {
		return profile
	}
	return ""
}

func (n *environmentNormalizer) company(hexString string) string {
	companyID := parsers.ParseEn1545Number(hexString)
	companyIDInt, err := strconv.Atoi(companyID)
	if err != nil {
		return ""
	}
	company := dictionaries.RavkavCompanies[companyIDInt]
	if company != "" {
		return company
	}
	return companyID
}
