package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"os"
	"os/user"
	"text/template"
	"time"
)

const (
	unknown = "UNKNOWN"
)

type configValues struct {
	LogTemplate    string `json:"logtemplate"`
	NoticeTemplate string `json:"noticetemplate"`
}

func thisUser() (string, string) {
	userName := unknown
	uID := unknown

	if c, err := user.Current(); err == nil {
		userName = c.Username
		uID = c.Uid
	}

	return userName, uID
}

func thisTTYName() string {
	tty := unknown

	dest, _ := os.Readlink("/proc/self/fd/0")

	if len(dest) > 0 {
		tty = dest
	}

	return tty
}

func log(data string) {
	if lw, err := syslog.New(syslog.LOG_NOTICE|syslog.LOG_AUTH, "jacksquat"); err == nil {
		fmt.Fprintf(lw, data)
		lw.Close()
	}
}

func getConfigFromReader(configReader io.Reader) configValues {
	var byteValue []byte
	var err error
	var config configValues

	if byteValue, err = ioutil.ReadAll(configReader); err != nil {
		log("config file could not be read")
	} else if json.Unmarshal(byteValue, &config) != nil {
		log("config file could not be parsed")
	}

	return config
}

func getConfig(configFile string) configValues {
	var config configValues

	if reader, err := os.Open(configFile); err == nil {
		defer reader.Close()
		config = getConfigFromReader(reader)
	}

	return config
}

func captureLoginWithConfig(config configValues, duration time.Duration, exitCheck func() bool) {
	if len(config.LogTemplate) > 0 {
		loginID, uID := thisUser()
		loginInfo := struct {
			UserName string
			UserID   string
			TTYName  string
		}{loginID, uID, thisTTYName()}
		if t, err := template.New("sysloginfo").Parse(config.LogTemplate); err == nil {
			var buf bytes.Buffer
			if t.Execute(&buf, loginInfo) == nil {
				log(buf.String())
			}
		}

		if t, err := template.New("noticeinfo").Parse(config.NoticeTemplate); err == nil {
			var buf bytes.Buffer
			if t.Execute(&buf, loginInfo) == nil {
				fmt.Println(buf.String())
			}
		}
	}

	for exitCheck() {
		time.Sleep(duration)
	}
}

func captureLogin(configFile string, duration time.Duration, exitCheck func() bool) {
	captureLoginWithConfig(getConfig(configFile), duration, exitCheck)
}
