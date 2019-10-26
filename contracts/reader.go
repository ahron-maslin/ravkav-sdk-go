package contracts

type Reader interface {
	ListReaders() ([]string, error)
	Connect(reader string) error
	Disconnect() error
	Transmit([]byte) ([]byte, error)
}
