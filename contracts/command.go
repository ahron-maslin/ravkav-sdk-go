package contracts

type Command interface {
	Name() string
	Kind() string
	Command() []byte
}
