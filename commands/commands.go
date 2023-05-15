package commands

import (
	"github.com/ahron-maslin/ravkav-sdk-go/contracts"
)

type Command struct {
	name    string
	kind    string
	command []byte
}

var calypsoCLA byte = 0x94

func NewRead(name string, command []byte) contracts.Command {
	var cmd []byte = getCommandPrefix()

	var prefix byte = 0xb2
	var suffix byte = 0x1d
	cmd = append(cmd, prefix)
	cmd = append(cmd, command...)
	cmd = append(cmd, suffix)

	return &Command{
		name:    name,
		kind:    "read",
		command: cmd,
	}
}

func NewRaw(name string, command []byte) contracts.Command {
	return &Command{
		name:    name,
		kind:    "raw",
		command: command,
	}
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Kind() string {
	return c.kind
}

func (c *Command) Command() []byte {
	return c.command
}

func getCommandPrefix() []byte {
	return []byte{calypsoCLA}
}
