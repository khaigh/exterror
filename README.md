exterror
=========

Extending GO's basic error handling by adding additional context and formatting to make debugging easier

Installation
------------

	go get github.com/khaigh/exterror


Basic Usage
-----------

	package main
	
	import (
	"os"
		"strconv"
	
	"github.com/khaigh/exterror"
	)
	
	func main() {
		err := innerFunc()
		if err != nil {
			os.Exit(1)
		}
	}
	
	func innerFunc() *exterror.Error {
		return innerInnerFunc()
	}
	
	func innerInnerFunc() *exterror.Error {
		_, err := strconv.Atoi("invalid number")
		return exterror.New(1234, "We have a problem, Houston", err).AndLog()
	}


Sample Log Output
-------------

	2014/07/10 16:44:26 ERROR:main.go:main.innerInnerFunc:23: We have a problem, Houston
		Error Number: 1234
		Trace (most recent call last):
			main.innerInnerFunc(0x4332b3)
				/home/vagrant/go/src/github.com/khaigh/sample/main.go:23 +0x6c
			main.innerFunc(0x400d47)
				/home/vagrant/go/src/github.com/khaigh/sample/main.go:18 +0x1e
			main.main()
				/home/vagrant/go/src/github.com/khaigh/sample/main.go:11 +0x1e
			
		Parent Error:	strconv.ParseInt: parsing "invalid number": invalid syntax
	exit status 1


Advanced Usage
--------------

### Debug Messages

In some cases, you want to add additional debugging information into the error log:

	return exterror.New(1234, "We have a problem, Houston", err).WithDebugMsg("debug message").AndLog()

Or dump local context to aid in debugging:

	numvar := 2
	strvar := "hello"
	return exterror.New(1234, "We have a problem, Houston", err).WithDebugField("numvar", numvar).WithDebugField("strvar", strval).AndLog()

### Error Formatting

If you would like to use a different error format, override the template using the text/template package with your own:

	var templateStr = "Error at {{.Filename}}:{{.CallingMethod}}:{{.Line}}: {{.EndUserMsg}}"
	tmpl, err := template.New("simpleTmpl").Parse(templateStr)
	if err != nil {
		t.Error("Error creating template")
	}
	return exterror.New(1234, "We have a problem, Houston", err).WithErrTemplate(tmpl).AndLog()

### Error Id
The error id makes it easy to do a find.  Line numbers are insufficent when people are making changes to code and you don't know which exact version of the app is generating which stack trace.
I use this shell command to generate them:

#!/bin/sh
od -vAn -N4 -tu4 < /dev/urandom | tr -d " \n"

Add this command to your favorite editor to generate unique ids.

