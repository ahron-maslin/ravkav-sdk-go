package commands

const (
	Meta   string = "meta"
	Record string = "record"
)

type APDU struct {
	Name         string
	Command      []byte
	Type         string
	Transmission string
}

// App name (1TIC.ICA)
var ApplicationAPDU APDU = APDU{Name: "application", Command: []byte{0x00, 0xa4, 0x04, 0x00, 0x08, 0x31, 0x54, 0x49, 0x43, 0x2e, 0x49, 0x43, 0x41}, Type: Meta}
var SelectAPDU APDU = APDU{Name: "select", Command: []byte{0x00, 0xc0, 0x00, 0x00}, Type: Meta}

// List containing all APDUs to be executed on card
var APDUs []APDU = []APDU{
	{Name: "environment", Command: []byte{0x01, 0x3c}, Type: Meta},
	{Name: "counters", Command: []byte{0x01, 0xcc}, Type: Meta},

	{Name: "event_logs", Command: []byte{0x01, 0x44}, Type: Record},
	{Name: "event_logs", Command: []byte{0x02, 0x44}, Type: Record},
	{Name: "event_logs", Command: []byte{0x03, 0x44}, Type: Record},
	{Name: "event_logs", Command: []byte{0x04, 0x44}, Type: Record},
	{Name: "event_logs", Command: []byte{0x05, 0x44}, Type: Record},
	{Name: "event_logs", Command: []byte{0x06, 0x44}, Type: Record},

	{Name: "special_events", Command: []byte{0x01, 0xec}, Type: Record},
	{Name: "special_events", Command: []byte{0x02, 0xec}, Type: Record},
	{Name: "special_events", Command: []byte{0x03, 0xec}, Type: Record},
	{Name: "special_events", Command: []byte{0x04, 0xec}, Type: Record},

	{Name: "contracts", Command: []byte{0x01, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x02, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x04, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x05, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x03, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x06, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x07, 0x4c}, Type: Record},
	{Name: "contracts", Command: []byte{0x08, 0x4c}, Type: Record},
}
