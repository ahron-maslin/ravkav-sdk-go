package contracts

type Normalizer interface {
	Normalize(record Record, recordIndex int) (map[string]interface{}, error)
}
