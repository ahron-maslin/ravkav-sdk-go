package card

import (
	"bytes"
	"fmt"
	"github.com/derkinderfietsen/ravkav-sdk-go/commands"
	"github.com/derkinderfietsen/ravkav-sdk-go/contracts"
)

type Populate struct {
	card *Card
}

// Populates card with records
func (c *Card) Populate() error {
	err := c.selectApplication()
	if err != nil {
		return err
	}

	for _, lfiItem := range commands.APDUs {
		var command contracts.Command = commands.NewRead(lfiItem.Name, lfiItem.Command)

		res, err := c.reader.Transmit(command.Command())
		if err != nil {
			return err
		}

		record, err := NewRecord(res, command)
		if err != nil {
			return err
		}

		if record == nil { // empty record
			continue
		}

		if lfiItem.Type == commands.Meta {
			c.records.meta[lfiItem.Name] = record
		} else {
			c.records.records[lfiItem.Name] = append(c.records.records[lfiItem.Name], record)
		}
	}
	return nil
}

func (c *Card) selectApplication() error {
	// Selecting application
	var command contracts.Command = commands.NewRaw(commands.ApplicationAPDU.Name, commands.ApplicationAPDU.Command)
	res, err := c.reader.Transmit(command.Command())
	if err != nil {
		return err
	} else if len(res) != 2 {
		return fmt.Errorf("wrong application length")
	}

	// Getting record number
	command = commands.NewRaw(commands.SelectAPDU.Name, append(commands.SelectAPDU.Command, 0x00))
	res, err = c.reader.Transmit(command.Command())
	if err != nil {
		return err
	} else if !bytes.Equal([]byte{res[0]}, []byte{0x6c}) {
		return fmt.Errorf("wrong application record response")
	}

	// Getting application record
	command = commands.NewRaw(commands.SelectAPDU.Name, append(commands.SelectAPDU.Command, res[1]))
	res, err = c.reader.Transmit(command.Command())
	if err != nil {
		return err
	}

	record, err := NewRecord(res, command)
	if err != nil {
		return err
	}

	c.records.meta[commands.ApplicationAPDU.Name] = record
	return nil
}
