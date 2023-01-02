package log

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func runLoggers() {
	Debug.Print("foo")
	Info.Print("bar")
	Warn.Print("baz")
	Error.Print("qux")
}

func TestDebugLogLevel(t *testing.T) {
	logFile := os.TempDir() + "debug.log"
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "debug")
	setupLoggers()

	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(logFile)
	if err != nil {
		t.Fatal(t)
	}

	got := string(f)

	want := fmt.Sprintf(
		`%[1]s [INFO] LOGLEVEL set to debug
[DEBUG] %[1]s log_test.go:11: foo
[INFO] %[1]s log_test.go:12: bar
[WARNING] %[1]s log_test.go:13: baz
[ERROR] %[1]s log_test.go:14: qux
`, timeNow)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}

func TestInfoLogLevel(t *testing.T) {
	logFile := os.TempDir() + "info.log"
	t.Setenv("LOGFILE", logFile)
	setupLoggers()

	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(logFile)
	if err != nil {
		t.Fatal(t)
	}

	got := string(f)

	want := fmt.Sprintf(
		"%[1]s [INFO] bar\n%[1]s [WARNING] baz\n%[1]s [ERROR] qux\n", timeNow,
	)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}

func TestWarningLogLevel(t *testing.T) {
	logFile := os.TempDir() + "warning.log"
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "warning")
	setupLoggers()

	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(logFile)
	if err != nil {
		t.Fatal(t)
	}

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
	logFile := os.TempDir() + "error.log"
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "error")
	setupLoggers()

	timeNow := time.Now().Format("2006/01/02 15:04:05")
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(logFile)
	if err != nil {
		t.Fatal(t)
	}

	got := string(f)

	want := fmt.Sprintf(
		"%[1]s [INFO] LOGLEVEL set to error\n%[1]s [ERROR] qux\n", timeNow,
	)

	if got != want {
		t.Fatalf("%sNOT THE SAME AS:\n%s", got, want)
	}
}
