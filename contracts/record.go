package contracts

type Record interface {
	Normalized() map[string]interface{}
	NormalizedVal(key string) interface{}
	SetNormalized(values map[string]interface{})
	SetNormalizedValue(key string, value string)
	Bytes() []byte
	Binary() string
	Hex() string
}
