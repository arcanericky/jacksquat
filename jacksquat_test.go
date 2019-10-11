package main

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

type testReader struct {
	data       string
	done       bool
	throwError bool
}

func (t *testReader) Read(p []byte) (n int, err error) {
	if t.throwError {
		return 0, io.ErrNoProgress
	}

	if t.done {
		return 0, io.EOF
	}

	for i, b := range []byte(t.data) {
		p[i] = b
	}

	t.done = true
	return len(t.data), nil
}

func TestJackSquat(t *testing.T) {
	// Retrieval of user name
	userName, userID := thisUser()
	if len(userName) == 0 {
		t.Error("Getting user name failed")
	}
	if len(userID) == 0 {
		t.Error("Getting user ID failed")
	}

	// Retrieval of TTY name
	ttyName := thisTTYName()
	if len(ttyName) == 0 {
		t.Error("Getting TTY name failed")
	}

	const configFilename = "jacksquat.conf"

	// Config file does not exist
	os.Remove(configFilename)
	getConfig(configFilename)

	ioutil.WriteFile(configFilename,
		[]byte(`{ "logtemplate": "login by {{.UserName}} (UID: {{.UserID}}) on {{.TTYName}}",
"noticetemplate": "Welcome {{.UserName}}. This is a captive account." }`),
		0644)

	// Config file exists with valid data
	loopCount := 0
	captureLogin(configFilename, 0, func() bool {
		ret := true

		if loopCount > 0 {
			ret = false
		}
		loopCount++

		return ret
	})

	os.Remove(configFilename)
	if loopCount == 0 {
		t.Error("capture loop was not called")
	}

	// Data parse error
	tr := &testReader{"data", false, false}
	if len(getConfigFromReader(tr).LogTemplate) > 0 {
		t.Error("config returned data on data parse error")
	}

	// Reader returns error
	tr.throwError = true
	if len(getConfigFromReader(tr).LogTemplate) > 0 {
		t.Error("config returned data on read error")
	}
}
