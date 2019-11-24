package formatters

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/SladeThe/common-go/log/terminal"
)

const (
	colorModeAuto = iota
	colorModeForceEnable
	colorModeForceDisable

	ColorModeEnvName = "S1_LOGRUS_COLOR_MODE"

	DefaultTimeFormat = "2006-01-02 15:04:05.000 MST"
	DefaultTimeColor  = -1

	DefaultCallerColor  = terminal.ForegroundColorBrightBlue
	DefaultMessageColor = -1
)

var (
	functionReplacementByRegexp = map[*regexp.Regexp]string{
		regexp.MustCompile("\\.\\(\\*([a-zA-Z0-9_]+)\\)\\."): ".$1.", // Remove dereference sign from type
		regexp.MustCompile("\\.func([0-9]+)$"):               "",     // Remove anonymous function name
	}

	colorModeByEnv = map[string]int{
		"auto":  colorModeAuto,
		"force": colorModeForceEnable,
		"off":   colorModeForceDisable,
	}

	colorMode = colorModeAuto
)

func init() {
	if colorModeEnv := strings.TrimSpace(os.Getenv(ColorModeEnvName)); len(colorModeEnv) > 0 {
		colorMode = colorModeByEnv[strings.ToLower(colorModeEnv)]
	}
}

var _ logrus.Formatter = &Logrus{}

type Logrus struct {
	// On Windows combine it with "github.com/mattn/go-colorable".
	// For example:
	//	if runtime.GOOS == "windows" {
	//		log.SetOutput(colorable.NewColorableStdout())
	//	} else {
	//		log.SetOutput(os.Stdout)
	//	}
	EnableColors bool

	// TODO The feature is under development.
	// Set this to true to enable advanced coloring and message parsing.
	// Disable, if you have any performance issues.
	AdvancedColors bool

	// Set this to true, if you do not want to print timestamp,
	// e.g. it is added by external log system.
	SkipTime bool

	// Default is DefaultTimeFormat.
	TimeFormat string

	// Use values less than zero to disable timestamp coloring.
	// Use a valid terminal color to set up coloring.
	// Any other value (for example, 0) means DefaultTimeColor.
	// Ignored if coloring is not enabled for this formatter.
	TimeColor int

	// Use values less than zero to disable caller coloring.
	// Use a valid terminal color to set up coloring.
	// Any other value (for example, 0) means DefaultCallerColor.
	// Ignored if coloring is not enabled for this formatter.
	CallerColor int

	// Use values less than zero to disable caller coloring.
	// Use a valid terminal color to set up coloring.
	// Any other value (for example, 0) means DefaultMessageColor.
	// Ignored if coloring is not enabled for this formatter.
	MessageColor int
}

func (formatter *Logrus) Format(entry *logrus.Entry) ([]byte, error) {
	buf := formatter.buffer(entry)

	if !formatter.SkipTime {
		timeText := entry.Time.Format(formatter.timeFormat())
		buf.WriteString(formatter.colorize(timeText, formatter.TimeColor, DefaultTimeColor))
		buf.WriteByte(' ')
	}

	formatter.printLevel(buf, entry)

	hasCaller := entry.HasCaller()

	if hasCaller && len(entry.Caller.Function) > 0 {
		functionName := filepath.Base(entry.Caller.Function)
		for compiledRegexp, replacement := range functionReplacementByRegexp {
			functionName = compiledRegexp.ReplaceAllString(functionName, replacement)
		}
		buf.WriteString(formatter.colorize(functionName, formatter.CallerColor, DefaultCallerColor))
		buf.WriteString(": ")
	}

	formatter.printMessage(buf, entry)

	if hasCaller && len(entry.Caller.File) > 0 {
		fileText := formatter.colorize(filepath.Base(entry.Caller.File), formatter.CallerColor, DefaultCallerColor)
		lineText := formatter.colorize(strconv.Itoa(entry.Caller.Line), formatter.CallerColor, DefaultCallerColor)
		buf.WriteString(fmt.Sprintf(" (%s:%s)", fileText, lineText))
	}

	buf.WriteByte('\n')

	return buf.Bytes(), nil
}

func (formatter *Logrus) buffer(entry *logrus.Entry) *bytes.Buffer {
	if entry.Buffer == nil {
		return &bytes.Buffer{}
	} else {
		return entry.Buffer
	}
}

func (formatter *Logrus) timeFormat() string {
	if len(formatter.TimeFormat) > 0 {
		return formatter.TimeFormat
	} else {
		return DefaultTimeFormat
	}
}

func (formatter *Logrus) printLevel(buf *bytes.Buffer, entry *logrus.Entry) {
	var levelText string

	switch entry.Level {
	case logrus.TraceLevel:
		levelText = "[TRACE]"
	case logrus.DebugLevel:
		levelText = "[DEBUG]"
	case logrus.InfoLevel:
		levelText = "[INFO ]"
	case logrus.WarnLevel:
		levelText = "[WARN ]"
	case logrus.ErrorLevel:
		levelText = "[ERROR]"
	case logrus.FatalLevel:
		levelText = "[FATAL]"
	case logrus.PanicLevel:
		levelText = "[PANIC]"
	}

	if len(levelText) > 0 {
		buf.WriteString(formatter.colorize2(levelText, levelColors(entry), [2]int{0, 0})) // TODO override colors
		buf.WriteByte(' ')
	}
}

func levelColors(entry *logrus.Entry) [2]int {
	switch entry.Level {
	case logrus.TraceLevel:
		return [2]int{terminal.ForegroundColorGray, -1}
	case logrus.DebugLevel:
		return [2]int{terminal.ForegroundColorBrightCyan, -1}
	case logrus.WarnLevel:
		return [2]int{terminal.ForegroundColorBrightYellow, -1}
	case logrus.ErrorLevel:
		return [2]int{terminal.ForegroundColorRed, -1}
	case logrus.FatalLevel, logrus.PanicLevel:
		return [2]int{terminal.ForegroundColorWhite, terminal.BackgroundColorRed}
	default:
		return [2]int{terminal.ForegroundColorCyan, -1}
	}
}

func (formatter *Logrus) printMessage(buf *bytes.Buffer, entry *logrus.Entry) {
	if formatter.AdvancedColors {
		buf.WriteString(formatter.colorize(entry.Message, formatter.MessageColor, DefaultMessageColor)) // TODO
	} else {
		buf.WriteString(formatter.colorize(entry.Message, formatter.MessageColor, DefaultMessageColor))
	}
}

func (formatter *Logrus) areColorsEnabled() bool {
	switch colorMode {
	case colorModeForceEnable:
		return true
	case colorModeForceDisable:
		return false
	default:
		return formatter.EnableColors
	}
}

func (formatter *Logrus) colorize(text string, color, defaultColor int) string {
	if formatter.areColorsEnabled() && color >= 0 {
		if !terminal.IsValidColor(color) {
			color = defaultColor
		}

		if terminal.IsValidColor(color) {
			return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, text)
		}
	}

	return text
}

func (formatter *Logrus) colorize2(text string, colors, defaultColors [2]int) string {
	if colors[0] < 0 {
		return formatter.colorize(text, colors[1], defaultColors[1])
	}

	if colors[1] < 0 {
		return formatter.colorize(text, colors[0], defaultColors[0])
	}

	if formatter.areColorsEnabled() {
		for i, color := range colors {
			if !terminal.IsValidColor(color) {
				colors[i] = defaultColors[i]
			}
		}

		return fmt.Sprintf("\x1b[%d;%dm%s\x1b[0m", colors[0], colors[1], text)
	}

	return text
}
