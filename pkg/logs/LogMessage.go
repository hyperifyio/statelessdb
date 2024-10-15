package logs

import "fmt"

type LogMessage struct {
	context string
	level   LogLevel
	format  string
	args    []interface{}
	depth   int
}

func (m *LogMessage) String() string {
	return "[" + m.context + "] [" + m.level.String() + "] [" + fmt.Sprintf("%03d", m.depth) + "] " + fmt.Sprintf(m.format, m.args...)
}
