// Copyright (c) Ken Haigh @khaigh All Rights Reserved.
// This source code is release under the MIT license, http://opensource.org/licenses/MIT

package exterror

import (
	"fmt"
	"runtime"
	"testing"
	"text/template"
)

func TestNew(t *testing.T) {
	_, _, line, _ := runtime.Caller(0)
	err := New(1234, "test error #1", nil)
	if err.Id != 1234 {
		t.Error("Id does not match")
	}
	if err.Filename != "exterror_test.go" {
		t.Error("Filename does not match")
	}
	if err.CallingMethod != "github.com/khaigh/exterror.TestNew" {
		t.Error("CallingMethod does not match")
	}
	if err.Line != (line + 1) {
		t.Error("Line does not match")
	}
	if err.EndUserMsg != "test error #1" {
		t.Error("EndUserMsg does not match")
	}
}

func TestAndLog(t *testing.T) {
	err := New(5678, "test error #2", nil).AndLog()
	if err.Id != 5678 {
		t.Error("Id does not match")
	}
}

func TestWithDebugMsg(t *testing.T) {
	err := New(11111111, "test error #3", nil).WithDebugMsg("debug message").AndLog()
	if err.Id != 11111111 {
		t.Error("Id does not match")
	}
	if err.DebugMsg != "debug message" {
		t.Error("DebugMsg does not match")
	}
}

func TestWithDebugFields(t *testing.T) {
	var testnum = 5
	var teststr = "test"
	err := New(4, "test error #4", nil).WithDebugField("testnum", testnum).WithDebugField("teststr", teststr).AndLog()
	if err.Id != 4 {
		t.Error("Id does not match")
	}
	if err.DebugFields["testnum"] != testnum {
		t.Error("DebugFields does not match")
	}
	if err.DebugFields["teststr"] != teststr {
		t.Error("DebugFields does not match")
	}
}

func TestParent(t *testing.T) {
	parentErr := New(1, "parent error", nil)
	err := New(2, "test error #5", parentErr).AndLog()
	if err.Id != 2 {
		t.Error("Id does not match")
	}
	if err.ParentErr.Error() == "" {
		t.Error("ParentErr does not match")
	}
}

func TestTemplate(t *testing.T) {
	var templateStr = "Error at {{.Filename}}:{{.CallingMethod}}:{{.Line}}: {{.EndUserMsg}}"
	tmpl, err := template.New("simpleTmpl").Parse(templateStr)
	if err != nil {
		t.Error("Error creating template")
	}
	_, _, line, _ := runtime.Caller(0)
	eerr := New(6, "test error #6", nil).WithErrTemplate(tmpl).AndLog()
	if eerr.Error() != fmt.Sprintf("Error at exterror_test.go:github.com/khaigh/exterror.TestTemplate:%d: test error #6", line+1) {
		t.Error("ErrTemplate does not match")
	}
}

func TestBadTemplate(t *testing.T) {
	funcMap := template.FuncMap{
		"badfunc": func(a []string) string { return a[1] },
	}
	var templateStr = "Error at {{badfunc .DebugMsg}}"
	tmpl, err := template.New("simpleTmpl").Funcs(funcMap).Parse(templateStr)
	if err != nil {
		t.Error("Error creating template")
	}
	defer func() {
		err := recover()
		if err == nil {
			t.Error("Did not throw an exception")
		}
	}()
	eerr := New(6, "test error #7", nil).WithErrTemplate(tmpl)
	eerr.Error()
	t.Error("Unreachable code")
}
