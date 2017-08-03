package cmdutils

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
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

	sentry = false
)

type tagMap map[string]string
type dataMap map[string]interface{}

// Logger is the logging interface
type Logger interface {
	EnableAggregation(string)
	KV(string, interface{}) Logger
	Tag(string, string)
	Log(string) Logger
	Error(error) Logger
	Fatal(error)
}

type LogMeta struct {
	Tags tagMap  `json:"-"`
	Data dataMap `json:"data"`
}

/*

	Underlying logging state

*/
type logger struct {
	glg    *glg.Glg
	scope  string
	buffer LogMeta
}

// NewLogger returns a new Logger interface
func NewLogger(scope string) Logger {

	g := glg.New()

	l := logger{
		glg:    g,
		scope:  scope,
		buffer: LogMeta{Tags: make(tagMap), Data: make(dataMap)},
	}

	l.setScopeLevels()
	return &l

}

// Error prints an error to stderr
func (lg *logger) Error(err error) Logger {

	if sentry {
		raven.CaptureError(err, lg.buffer.Tags)
	}

	s := err.Error()
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

	if sentry {
		raven.Capture(lg.packet(s), lg.buffer.Tags)
	}

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

	lg.buffer.Data[k] = v
	return lg

}

// Tag adds a log tag to be sent to sentry
func (lg *logger) Tag(k, v string) {

	lg.buffer.Tags[k] = v

}

// Fatal logs err and kills the process
func (lg *logger) Fatal(err error) {

	if sentry {
		raven.CapturePanicAndWait(func() {
			panic(err)
		}, lg.buffer.Tags)
	}

	fields := lg.checkBuffer()
	lg.glg.Fatalf(LN_FMT, err.Error(), fields)

}

// EnableAggregation will send logs to sentry
func (lg *logger) EnableAggregation(tokenAndProject string) {

	info := SplitBySpaces(tokenAndProject)
	if len(info) < 2 {
		lg.Fatal(fmt.Errorf("Bad Sentry Data"))
	}

	err := raven.SetDSN(sentryUrl(info[0], info[1]))
	if err != nil {
		lg.Fatal(err)
	}

	sentry = true

}

/*

	Clear and return the buffer as a JSON string

*/
func (lg *logger) checkBuffer() string {

	if len(lg.buffer.Data) > 0 {
		jsn, err := json.Marshal(lg.buffer.Data)
		lg.buffer.Data = make(dataMap)
		if err != nil {
			return fmt.Sprintf("formatErr: %v", err)
		}
		return string(jsn)
	}

	lg.buffer.Tags = make(tagMap)
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

/*

	Create a sentry connect url

*/
func sentryUrl(tk, project string) string {

	return "https://" + tk + "@sentry.io/" + project

}

/*

	Get a new packet for sending info logs to sentry

*/
func (lg *logger) packet(s string) *raven.Packet {

	p := raven.NewPacket(s)
	for k, v := range lg.buffer.Data {
		p.Extra[k] = v
	}
	p.Level = raven.INFO
	p.Logger = lg.scope

	return p

}
