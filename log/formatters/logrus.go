package formatters

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"
)

var functionReplacementByRegexp = map[*regexp.Regexp]string{
	regexp.MustCompile("\\.\\(\\*([a-zA-Z0-9_]+)\\)\\."): ".$1.", // Remove dereference sign from type
	regexp.MustCompile("\\.func([0-9]+)$"):               "",     // Remove anonymous function name
}

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 36
	colorGray   = 37
)

var _ logrus.Formatter = &Logrus{}

type Logrus struct {
	// On Windows combine it with "github.com/mattn/go-colorable".
	// For example:
	//	if runtime.GOOS == "windows" {
	//		log.SetOutput(colorable.NewColorableStdout())
	//	} else {
	//		log.SetOutput(os.Stdout)
	//	}
	// TODO env var to force enable/disable/auto
	EnableColors bool
}

func (formatter *Logrus) Format(entry *logrus.Entry) ([]byte, error) {
	var buf *bytes.Buffer
	if entry.Buffer == nil {
		buf = &bytes.Buffer{}
	} else {
		buf = entry.Buffer
	}

	timeFormat := "2006-01-02 15:04:05.999"
	timeString := entry.Time.Format(timeFormat)

	for len(timeString) < len(timeFormat) {
		timeString += "0"
	}

	buf.WriteString(timeString)

	var printLevel func(levelText string)

	if formatter.EnableColors {
		printLevel = func(levelText string) {
			var levelColor int

			switch entry.Level {
			case logrus.DebugLevel, logrus.TraceLevel:
				levelColor = colorGray
			case logrus.WarnLevel:
				levelColor = colorYellow
			case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
				levelColor = colorRed
			default:
				levelColor = colorBlue
			}

			_, _ = fmt.Fprintf(buf, "\x1b[%dm%s\x1b[0m", levelColor, levelText)
		}
	} else {
		printLevel = func(levelText string) {
			buf.WriteString(levelText)
		}
	}

	switch entry.Level {
	case logrus.TraceLevel:
		printLevel(" [TRACE] ")
	case logrus.DebugLevel:
		printLevel(" [DEBUG] ")
	case logrus.InfoLevel:
		printLevel(" [INFO ] ")
	case logrus.WarnLevel:
		printLevel(" [WARN ] ")
	case logrus.ErrorLevel:
		printLevel(" [ERROR] ")
	case logrus.FatalLevel:
		printLevel(" [FATAL] ")
	case logrus.PanicLevel:
		printLevel(" [PANIC] ")
	}

	if entry.HasCaller() {
		functionName := filepath.Base(entry.Caller.Function)
		for compiledRegexp, replacement := range functionReplacementByRegexp {
			functionName = compiledRegexp.ReplaceAllString(functionName, replacement)
		}
		buf.WriteString(functionName)
		buf.WriteString(": ")
	}

	buf.WriteString(entry.Message)

	if entry.HasCaller() {
		buf.WriteString(fmt.Sprintf(" (%s:%d)", filepath.Base(entry.Caller.File), entry.Caller.Line))
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
