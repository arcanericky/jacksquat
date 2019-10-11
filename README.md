# Jack Squat

[![Build Status](https://travis-ci.com/arcanericky/jacksquat.svg?branch=master)](https://travis-ci.com/arcanericky/jacksquat)
[![codecov](https://codecov.io/gh/arcanericky/jacksquat/branch/master/graph/badge.svg)](https://codecov.io/gh/arcanericky/jacksquat)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

For when you need a login shell that can't do Jack Squat.

### Quick Start

Copy `jacksquat` to a place where it can be used as a login shell. For Debian, `/usr/sbin/` is a good location.

```
$ cp jacksquat /usr/sbin/
```

Create the locked user, configured with the `jacksquat` shell.

```
$ useradd -m -s /usr/sbin/jacksquat jack
$ passwd jack
New password: 
Retype new password: 
passwd: password updated successfully
```

Validate the user "can't do jack squat".

```
$ ssh jack@localhost
Password:
```

If a log entry and/or a notice is desired when `jacksquat` is executed, add an `/etc/jacksquat.conf` file in json format containing a Go template of your log line (`logtemplate`) and notice (`noticetemplate`). Both fields are optional as is the configuration file. Logging is done with the `jacksquat` tag with a priority of `syslog.LOG_NOTICE|syslog.LOG_AUTH`. The following is an example and shows all of the available fields.

```
echo '{"logtemplate":"login by {{.UserName}} (UID: {{.UserID}}) on {{.TTYName}}","noticetemplate":"Welcome {{.UserName}}. This is a captive account."}' > /etc/jacksquat.conf
```

You may want to use a [tool such as `jq`](https://stedolan.github.io/jq/) to validate the json is parsable.

```
cat /etc/jacksquat.conf | jq '.'
{
  "logtemplate": "login by {{.UserName}} (UID: {{.UserID}}) on {{.TTYName}}",
  "noticetemplate": "Welcome {{.UserName}}. This is a captive account."
}
```

### Inspiration

I need a user login that could remain logged in yet "can't do jack squat". The logging, use of Go templates in a configuration file, and testing are most certainly overkill.