//  Copyright (c) Ken Haigh (http://kenhaigh.com) All Rights Reserved.

package eerr

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type Error struct {
	Num           int64
	Filename      string
	CallingMethod string
	Line          int
	EndUserMsg    string
	DebugMsg      string
	DebugFields   map[string]interface{}
	Err           error
	StackTrace    string
}

func New(num int64, endUserMsg string, parentErr error) *Error {
	e := new(Error)
	e.Num = num
	e.EndUserMsg = endUserMsg
	e.Err = parentErr
	e.DebugFields = make(map[string]interface{})

	pc, file, line, ok := runtime.Caller(1)

	if ok {
		e.Line = line
		components := strings.Split(file, "/")
		e.Filename = components[(len(components) - 1)]
		f := runtime.FuncForPC(pc)
		e.CallingMethod = f.Name()
	}

	const size = 1 << 12
	buf := make([]byte, size)
	n := runtime.Stack(buf, false)

	e.StackTrace = string(buf[:n])

	log.Print(e)
	return e
}

func (e *Error) AddDebugField(key string, value interface{}) {
	e.DebugFields[key] = value
}

func (e *Error) Error() string {
	parentError := "nil"
	if e.Err != nil {
		parentError = prependToLines(e.Err.Error(), "-- ")
	}
	debugFieldStrings := make([]string, 0, len(e.DebugFields))
	for k, v := range e.DebugFields {
		str := fmt.Sprintf("\n-- DebugField[%s]: %+v", k, v)
		debugFieldStrings = append(debugFieldStrings, str)
	}
	dbgMsg := ""
	if len(e.DebugMsg) > 0 {
		dbgMsg = "\n-- DebugMsg: " + e.DebugMsg
	}

	return fmt.Sprintln(
		"\n\n-- Error",
		e.Num,
		e.Filename,
		e.CallingMethod,
		"line:", e.Line,
		"\n-- EndUserMsg: ", e.EndUserMsg,
		dbgMsg,
		strings.Join(debugFieldStrings, ""),
		"\n-- StackTrace:",
		strings.TrimLeft(prependToLines(e.StackTrace, "-- "), " "),
		"\n-- ParentError:", parentError,
	)
}

func prependToLines(para, prefix string) string {
	lines := strings.Split(para, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}
