package cmdutils

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
)

var colors = []func(string) string{
	glg.Gray,
	glg.Cyan,
	glg.Green,
	glg.Orange,
	glg.Purple,
	glg.Brown,
	glg.Yellow,
}

var (
	// ERR_EXT is the error extension tacked onto scopes
	ERR_EXT = "-err"

	// LN_FMT is the formatting of each log
	LN_FMT = "%-35s %s"

	lastUsedColor = 0
)

// Logger is the logging interface
type Logger interface {
	KV(string, interface{}) Logger
	Log(string) Logger
	Error(string) Logger
}

/*

	Underlying logging state

*/
type logger struct {
	glg    *glg.Glg
	scope  string
	buffer map[string]interface{}
}

// Error prints an error to stderr
func (lg *logger) Error(s string) Logger {

	fields := lg.checkBuffer()
	if len(lg.scope) > 0 {
		lg.glg.CustomLogf(lg.scope+ERR_EXT, LN_FMT, s, fields)
		return lg
	}

	lg.glg.Errorf(LN_FMT, s, fields)
	return lg

}

// Log prints a string to log output
func (lg *logger) Log(s string) Logger {

	fields := lg.checkBuffer()
	if len(lg.scope) > 0 {
		lg.glg.CustomLogf(lg.scope, LN_FMT, s, fields)
		return lg
	}

	lg.glg.Logf(LN_FMT, s, fields)
	return lg

}

// KV attaches a key value pair to a log statement
func (lg *logger) KV(k string, v interface{}) Logger {

	lg.buffer[k] = v
	return lg

}

/*

	Clear and return the buffer as a JSON string

*/
func (lg *logger) checkBuffer() string {

	if len(lg.buffer) > 0 {
		jsn, err := json.Marshal(lg.buffer)
		lg.buffer = make(map[string]interface{})
		if err != nil {
			return fmt.Sprintf("formatErr: %v", err)
		}
		return string(jsn)
	}

	return ""

}

/*

	Set the new log scope levels

*/
func (lg *logger) setScopeLevels() {

	errLevel := lg.scope + ERR_EXT
	lg.glg.AddErrLevel(errLevel)
	lg.glg.SetLevelColor(errLevel, glg.Red)
	lg.glg.AddStdLevel(lg.scope)

	// set a color for the output
	lg.glg.SetLevelColor(lg.scope, lg.nextColor())

}

/*

	Grab the next color in colors list after the last used one

*/
func (lg *logger) nextColor() func(string) string {

	if len(colors) == lastUsedColor+1 {
		lastUsedColor = 0
	} else {
		lastUsedColor = lastUsedColor + 1
	}

	return colors[lastUsedColor]

}

// NewLogger returns a new Logger interface
func NewLogger(scope string) Logger {

	g := glg.New()

	l := logger{
		glg:    g,
		scope:  scope,
		buffer: make(map[string](interface{})),
	}

	l.setScopeLevels()
	return &l

}
