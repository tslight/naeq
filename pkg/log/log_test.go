package log

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
[DEBUG] %[1]s log_test.go:15: foo
[INFO] %[1]s log_test.go:16: bar
[WARNING] %[1]s log_test.go:17: baz
[ERROR] %[1]s log_test.go:18: qux
`, timeNow)

	if got != want {
		t.Fatalf("\n%s\nNOT THE SAME AS:\n%s", got, want)
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
		t.Fatalf("\n%s\nNOT THE SAME AS:\n%s", got, want)
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
		t.Fatalf("\n%s\nNOT THE SAME AS:\n%s", got, want)
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
		t.Fatalf("\n%s\nNOT THE SAME AS:\n%s", got, want)
	}
}

func TestLogRequest(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ip := req.RemoteAddr
	setupLoggers()
	timeNow := time.Now().Format("2006/01/02 15:04:05")
	Request(req)

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := strings.TrimSpace(string(f))

	want := fmt.Sprintf(
		"%[1]s [INFO] GET to example.com/ from %[2]s\n%[1]s [INFO] AGENT:", timeNow, ip,
	)

	if got != want {
		t.Fatalf("\n%s\nNOT THE SAME AS:\n%s", got, want)
	}
}

func TestTraceLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "trace")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// trace is alias for debug
	if !strings.Contains(got, "[DEBUG]") {
		t.Fatalf("Expected [DEBUG] output for trace level, got:\n%s", got)
	}
}

func TestInformationLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "information")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// information is alias for info
	if !strings.Contains(got, "[WARNING]") || strings.Contains(got, "[DEBUG]") {
		t.Fatalf("Expected INFO level output, got:\n%s", got)
	}
}

func TestErrLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "err")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// err is alias for error
	if !strings.Contains(got, "[ERROR]") || strings.Contains(got, "[WARNING]") {
		t.Fatalf("Expected ERROR level output, got:\n%s", got)
	}
}

func TestInvalidLogLevel(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("LOGLEVEL", "invalid_level")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// Invalid level should warn and fall back to INFO
	if !strings.Contains(got, "INVALID") || !strings.Contains(got, "[INFO]") {
		t.Fatalf("Expected INVALID warning and INFO output, got:\n%s", got)
	}
}

func TestDebugEnvVar(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("DEBUG", "true")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// DEBUG=true should enable debug output
	if !strings.Contains(got, "[DEBUG]") {
		t.Fatalf("Expected [DEBUG] output when DEBUG=true, got:\n%s", got)
	}
}

func TestDebugEnvVarCaseInsensitive(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("DEBUG", "True")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// DEBUG=True (uppercase) should also work
	if !strings.Contains(got, "[DEBUG]") {
		t.Fatalf("Expected [DEBUG] output when DEBUG=True, got:\n%s", got)
	}
}

func TestDebugEnvVarFalse(t *testing.T) {
	t.Setenv("LOGFILE", logFile)
	t.Setenv("DEBUG", "false")
	t.Setenv("LOGLEVEL", "info")

	setupLoggers()
	runLoggers()

	f, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	rmLogFile(t)

	got := string(f)
	// DEBUG=false should NOT enable debug output
	if strings.Contains(got, "[DEBUG]") {
		t.Fatalf("Did not expect [DEBUG] output when DEBUG=false, got:\n%s", got)
	}
}
