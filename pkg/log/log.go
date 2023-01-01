package log

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func setLogLevel(outMW io.Writer) {
	isDebug, isDebugPresent := os.LookupEnv("DEBUG")
	if isDebugPresent && strings.ToLower(isDebug) == "true" {
		Debug.SetOutput(outMW)
		Debug.Println("LOGLEVEL set to [DEBUG]")
		return
	}
	level, levelPresent := os.LookupEnv("LOGLEVEL")
	if levelPresent {
		switch strings.ToLower(level) {
		case "1", "debug", "trace":
			Debug.SetOutput(outMW)
			Debug.Println("LOGLEVEL set to " + level)
		case "2", "info", "information":
			Debug.SetOutput(io.Discard)
		case "3", "warn", "warning":
			Debug.SetOutput(io.Discard)
			Info.SetOutput(io.Discard)
		case "4", "err", "error":
			Debug.SetOutput(io.Discard)
			Info.SetOutput(io.Discard)
			Warn.SetOutput(io.Discard)
		default:
			Warn.Printf("LOGLEVEL=%s is INVALID! Using default [INFO].", level)
		}
	}
}

func init() {
	var file *os.File
	var err error

	envFile, envFilePresent := os.LookupEnv("LOGFILE")
	if envFilePresent {
		log.Println("Logging to " + envFile)
		file, err = os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	} else {
		t := time.Now().Format("2006-01-02")
		file, err = os.OpenFile(t+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	}

	if err != nil {
		log.Fatal(err)
	}

	outMW := io.MultiWriter(os.Stdout, file)
	errMW := io.MultiWriter(os.Stderr, file)

	Debug = log.New(io.Discard, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(outMW, "[INFO] ", log.LstdFlags)
	Warn = log.New(outMW, "[WARNING] ", log.LstdFlags)
	Error = log.New(errMW, "[ERROR] ", log.LstdFlags)

	setLogLevel(outMW)
}

func Request(r *http.Request) {
	Debug.Printf("REQUEST:\n%#v", r)
	Debug.Printf("REQUEST.URL:\n%#v", r.URL)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		Error.Fatal(err)
	}
	// Replace the body with a new reader after reading from the original
	r.Body = io.NopCloser(bytes.NewBuffer(b))
	var ip string
	ip = r.Header.Get("Client-Ip")
	if ip == "" {
		ip = r.RemoteAddr
	}
	Info.Printf(
		"%s to %s%s from %s\n", r.Method, r.Host, r.URL.Path, ip,
	)
	agent := r.Header.Get("User-Agent")
	Info.Print("AGENT: ", agent)
	if r.URL.RawQuery != "" {
		Info.Print("QUERY: ", r.URL.RawQuery)
	}
	bstr := string(b)
	if bstr != "" {
		Info.Print("BODY: ", bstr)
	}
}
