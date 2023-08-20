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
	return 0
}

func (c *Command) setCmd() {
}

func (c *Command) getCmd() Cmd {
	return 0
}

func (c *Command) setData(data string) {
}

func (c *Command) getData() string {
	return ""
}

func (c *Command) setDb(db string) {
}

func (c *Command) getDb() string {
	return ""
}

func (c *Command) setIsQuery(isQuery bool) {
}

func (c *Command) getIsQuery() bool {
	return true
}

func (cr *CommandResponse) encode(data types.IoBuffer) {
}

func (cr *CommandResponse) setData(data string) {
}

func (cr *CommandResponse) getData() string {
	return ""
}
