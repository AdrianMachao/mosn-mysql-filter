package mysql

import "mosn.io/mosn/pkg/types"

type Cmd int

type Command struct {
	MySQLCodec
	cmd     Cmd
	data    string
	db      string
	isQuery bool
}

type CommandResponse struct {
	MySQLCodec
	data string
}

const (
	Null  Cmd = -1
	Sleep Cmd = iota - 1
	Quit
	InitDb
	Query
	FieldList
	CreateDb
	DropDb
	Refresh
	Shutdown
	Statistics
	ProcessInfo
	Connect
	ProcessKill
	Debug
	Ping
	Time
	DelayedInsert
	ChangeUser
	Daemon          Cmd = 29
	ResetConnection Cmd = 31
)

func (c *Command) parseCmd(data types.IoBuffer) Cmd {
	cmdBytes := data.Peek(1)
	if cmdBytes == nil {
		return Null
	}

	cmd := Cmd(uint32(cmdBytes[0]))

	data.Drain(1)

	return cmd
}

func (c *Command) setCmd(cmd Cmd) {
	c.cmd = cmd
}

func (c *Command) getCmd() Cmd {
	return c.cmd
}

func (c *Command) setData(data string) {
	c.data = data
}

func (c *Command) getData() string {
	return c.data
}

func (c *Command) setDb(db string) {
	c.db = db
}

func (c *Command) getDb() string {
	return c.db
}

func (c *Command) setIsQuery(isQuery bool) {
	c.isQuery = isQuery
}

func (c *Command) getIsQuery() bool {
	return c.isQuery
}

func (cr *Command) encode(data types.IoBuffer) {
}

func (cr *Command) decode(data types.IoBuffer, seq uint8, length int) DecodeStatus {
	cr.seq = seq
	return cr.parseMessage(data, length)
}

func (cr *Command) parseMessage(data types.IoBuffer, length int) DecodeStatus {
	cmd := cr.parseCmd(data)

	cr.setCmd(cmd)
	if cmd == Null {
		return Failure
	}

	switch cmd {
	case InitDb:
	case CreateDb:
	case DropDb:
		db := data.Peek(length - 1)
		cr.setDb(string(db))
		break
	case Query:
		cr.isQuery = true
	default:
		data.Peek(length - 1)
		break
	}

	return Success
}

func (cr *CommandResponse) parseMessage(data types.IoBuffer, length int) DecodeStatus {
	cmdResp := data.Peek(length - 1)
	if cmdResp == nil {
		return Failure
	}

	cr.data = string(cmdResp)

	return Success
}

// if (BufferHelper::readStringBySize(buffer, len, data_) != DecodeStatus::Success) {
//     ENVOY_LOG(debug, "error when parsing command response");
//     return DecodeStatus::Failure;
//   }
//   return DecodeStatus::Success;

func (cr *CommandResponse) encode(data types.IoBuffer) {
}

func (cr *CommandResponse) setData(data string) {
	cr.data = data
}

func (cr *CommandResponse) getData() string {
	return cr.data
}
