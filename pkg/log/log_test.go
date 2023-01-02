package log

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func runLoggers() {
	Debug.Print("foo")
	Info.Print("bar")
	Warn.Print("baz")
	Error.Print("qux")
}

var logName = "go-pkg-log-test.log"
var logFile = filepath.Join(os.TempDir(), logName)

func rmLogFile(t *testing.T) {
	err := file.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(logFile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDebugLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "debug")

	setupLoggers()
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)

	want := fmt.Sprintf(`%[1]s [INFO] LOGLEVEL set to debug
[DEBUG] %[1]s log_test.go:12: foo
[INFO] %[1]s log_test.go:13: bar
[WARNING] %[1]s log_test.go:14: baz
[ERROR] %[1]s log_test.go:15: qux
`, timeNow)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}

func TestInfoLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)

	setupLoggers()
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)

	want := fmt.Sprintf(
		"%[1]s [INFO] bar\n%[1]s [WARNING] baz\n%[1]s [ERROR] qux\n", timeNow,
	)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}

func TestWarningLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "warning")

	setupLoggers()
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)

	want := fmt.Sprintf(
		`%[1]s [INFO] LOGLEVEL set to warning
%[1]s [WARNING] baz
%[1]s [ERROR] qux
`, timeNow,
	)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}

func TestErrorLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "error")

	setupLoggers()
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)

	want := fmt.Sprintf(
		"%[1]s [INFO] LOGLEVEL set to error\n%[1]s [ERROR] qux\n", timeNow,
	)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}
