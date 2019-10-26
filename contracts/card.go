package contracts

type Card interface {
	Populate() error
	GetRecords() map[string][]Record
	GetRecord(name string) []Record
	GetMeta() map[string]Record
	Normalize() error
	Output() CardOutput
}

type CardOutput interface {
	JSON() (string, error)
	GetMeta() map[string]map[string]interface{}
	GetRecords() map[string][]map[string]interface{}
}
