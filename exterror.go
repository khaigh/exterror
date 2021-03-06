// Copyright (c) Ken Haigh @khaigh All Rights Reserved.
// This source code is release under the MIT license, http://opensource.org/licenses/MIT

// The exterror package extends go's basic error handling by adding additional context and formatting to make
// debugging easier.
package exterror

import (
	"bytes"
	"log"
	"path"
	"runtime"
	"strings"
	"text/template"
)

//
type Error struct {
	Id            int64
	Filename      string
	CallingMethod string
	Line          int
	EndUserMsg    string
	DebugMsg      string
	DebugFields   map[string]interface{}
	ParentErr     error
	StackTrace    string
	Template      *template.Template
}

const templateStr = "ERROR:{{.Filename}}:{{.CallingMethod}}:{{.Line}}: {{.EndUserMsg}}" +
	"\n\tError Number: {{.Id}}" +
	"{{if (gt (len .DebugMsg) 0)}}" +
	"\n\tDebug Message: {{.DebugMsg}}" +
	"{{end}}" +
	"{{range $key, $value := .DebugFields}}" +
	"\n\t{{$key}}: {{printf \"%+v\" $value}}" +
	"{{end}}" +
	"\n\tTrace (most recent call last):" +
	"{{range $line := (trimleft .StackTrace | splitLines) }}" +
	"\n\t{{concat \"\t\" $line}}" +
	"{{end}}" +
	"{{if (gt (len (error .ParentErr)) 0)}}" +
	"\n\tParent Error:" +
	"{{range $line := (error .ParentErr | trimleft | splitLines) }}" +
	"{{concat \"\t\" $line}}\n" +
	"{{end}}" +
	"{{end}}"

var defaultTemplate = initDefaultTemplate()

func initDefaultTemplate() *template.Template {
	funcMap := template.FuncMap{
		"trimleft":   func(a string) string { return strings.TrimLeft(a, " ") },
		"splitLines": func(para string) []string { return strings.Split(para, "\n") },
		"concat":     func(a string, b string) string { return a + b },
		"error": func(e error) string {
			if e != nil {
				return e.Error()
			}
			return ""
		},
	}
	tmpl, err := template.New("defaultErrTmpl").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// New creates a pointer to an extended error. The id is a unique number
// generated by the caller to make finding a error easier. The endUserMsg provides a user appropriate error
// message. The parentErr allows for chaining of errors to find the source error.
func New(id int64, endUserMsg string, parentErr error) *Error {
	e := new(Error)
	e.Id = id
	e.EndUserMsg = endUserMsg
	e.ParentErr = parentErr
	e.DebugFields = make(map[string]interface{})
	e.Line, e.Filename, e.CallingMethod = getCallingFuncInfo()
	e.StackTrace = getStackTrace()
	return e
}

func getCallingFuncInfo() (lineNum int, fileName string, callingMethod string) {
	const skipTwoFunctions = 2 // this function plus parent
	pc, file, line, ok := runtime.Caller(skipTwoFunctions)
	if !ok {
		return 0, "", ""
	}
	f := runtime.FuncForPC(pc)
	return line, path.Base(file), f.Name()
}

func getStackTrace() string {
	const skipTwoFrames = 5 // the goroutine line plus the frame for this method and the new func
	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)
	return strings.Join(strings.Split(string(buf[:n]), "\n")[skipTwoFrames:], "\n")
}

// AndLog is provided for convenience to write the error to the log upon creation.
// For example, exterror.New(1234, "user message", err).AndLog()
func (e *Error) AndLog() *Error {
	log.Print(e)
	return e
}

// WithDebugField allows the caller to include extra debugging information into an error message.
// For example, exterror.New(1234, "user message", err).WithDebugField("testvar", a).AndLog()
func (e *Error) WithDebugField(key string, value interface{}) *Error {
	e.DebugFields[key] = value
	return e
}

// WithDebugMsg allows the caller to set a debug message into an error message
// For example, exterror.New(1234, "user message", err).WithDebugMsg("example msg").AndLog()
func (e *Error) WithDebugMsg(msg string) *Error {
	e.DebugMsg = msg
	return e
}

// WithTemplate allows the caller to change the formatting of an error message
// For example, exterror.New(1234, "user message", err).WithTemplate(mysimpletmpl).AndLog()
func (e *Error) WithTemplate(tmpl *template.Template) *Error {
	e.Template = tmpl
	return e
}

// Error implements the error interface and returns a formatted error string for an extended error.
func (e *Error) Error() string {
	var tmpl *template.Template
	var buf bytes.Buffer

	if e.Template != nil {
		tmpl = e.Template
	} else {
		tmpl = defaultTemplate
	}
	err := tmpl.Execute(&buf, e)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
