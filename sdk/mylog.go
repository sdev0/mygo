package sdk

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

type MyLog struct {
	mu         sync.Mutex // ensures atomic writes; protects the following fields
	prefix     string     // prefix on each line to identify the logger (but see Lmsgprefix)
	flag       int        // properties
	out        io.Writer  // destination for output
	buf        []byte     // for accumulating text to write
	headerlen  int
	showHeader bool
}

func NewMyLog(out io.Writer, prefix string, flag int) *MyLog {
	return &MyLog{out: out, prefix: prefix, flag: flag, headerlen: 40, showHeader: true}
}

// SetOutput sets the output destination for the logger.
func (l *MyLog) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

var std = NewMyLog(os.Stderr, "", LstdFlags)

// Default returns the standard logger used by the package-level output functions.
func Default() *MyLog { return std }

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank and Lmsgprefix is unset),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided),
//   * l.prefix (if it's not blank and Lmsgprefix is set).
func (l *MyLog) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	if l.flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.prefix...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, '[')
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		lgth := len(*buf)
		for lgth <= l.headerlen {
			*buf = append(*buf, ' ')
			lgth++
		}
		*buf = append(*buf, "] "...)
	}
	if l.flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// MyLog. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *MyLog) Output(calldepth int, s string) (string, string, error) {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	header := l.buf
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return string(header), string(l.buf[len(header):]), err
}

// Printf calls res, _ := l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *MyLog) Printf(format string, v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf(format, v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}

// Print calls res, _ := l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *MyLog) Print(v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprint(v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}

// Println calls res, _ := l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *MyLog) Println(v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintln(v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Println(res)
}

// Flags returns the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *MyLog) Flags() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flag
}

// SetFlags sets the output flags for the logger.
// The flag bits are Ldate, Ltime, and so on.
func (l *MyLog) SetFlags(flag int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.flag = flag
}

// Prefix returns the output prefix for the logger.
func (l *MyLog) Prefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the logger.
func (l *MyLog) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

// Writer returns the output destination for the logger.
func (l *MyLog) Writer() io.Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out
}

// set head length
func (l *MyLog) SetHeaderLen(length int) {
	l.headerlen = length
}
// set whether show header
func (l *MyLog) SetShowHeader(show bool) {
	l.showHeader = show
}
func (l *MyLog) Llog(v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintln(v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}
func (l *MyLog) Llogf(format string, v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf(format, v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}
func (l *MyLog) LogInfo(v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf("[INFO] %v\n", v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}
func (l *MyLog) LogInfof(format string, v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf("[INFO] "+format, v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}
func (l *MyLog) LogErr(v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf("[ERROR] %v\n", v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}
func (l *MyLog) LogErrf(format string, v ...interface{}) {
	head, res, _ := l.Output(2, fmt.Sprintf("[ERROR] "+format, v...))
	if l.showHeader {
		fmt.Print(head)
	}
	fmt.Print(res)
}

// 初始化log，启动时以每1小时为单位创建
//  @param  logpath [string]
//  @param  flag [int]
//	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
//	Ltime                         // the time in the local time zone: 01:23:23
//	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
//	Llongfile                     // full file name and line number: /a/b/c/d.go:23
//	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
//	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
//	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
//	LstdFlags     = Ldate | Ltime // initial values for the standard logger
//  @return [*MyLog]
func InitLog(logpath string, flag int) *MyLog {
	err := os.MkdirAll(logpath, os.ModePerm)
	CheckError(err, true, "[ERROR] create log path failed. path:", logpath)
	logfile, err := os.OpenFile(logpath+"/logger."+GetTimeStrByKey(T_DATE|T_HOUR)+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	CheckError(err, true, "[ERROR] create log file failed.")

	logger := NewMyLog(logfile, "", flag)
	return logger
}
