package logs

import "fmt"

type LogMessage struct {
	level  LogLevel
	format string
	args   []interface{}
}

func (m *LogMessage) String() string {
	return "[" + m.level.String() + "] " + fmt.Sprintf(m.format, m.args...)
}
