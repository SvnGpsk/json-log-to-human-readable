package cmd

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"
)

// QuarkusLogMessage Quarkus Standard Log message type
type QuarkusLogMessage struct {
	Timestamp  string    `json:"timestamp"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Exception  Exception `json:"exception,omitempty"`
	LoggerName string    `json:"loggerName"`
	Tracing    Tracing   `json:"mdc"`
}

// SpringBootLogMessage Spring Boot Log message type
type SpringBootLogMessage struct {
	Timestamp  string `json:"@timestamp"`
	Level      string `json:"level"`
	Message    string `json:"message"`
	Exception  string `json:"stack_trace,omitempty"`
	LoggerName string `json:"logger_name"`
}

// GoZapLogMessage Uber Zap log message type
type GoZapLogMessage struct {
	Level      string  `json:"level"`
	Timestamp  float64 `json:"ts"`
	Logger     string  `json:"logger"`
	Message    string  `json:"msg"`
	Controller string  `json:"controller,omitempty"`
	Request    string  `json:"request,omitempty"`
	Error      string  `json:"error,omitempty"`
	Stacktrace string  `json:"stacktrace,omitempty"`
}

// DotNetLogMessage type
type DotNetLogMessage struct {
	Timestamp  string `json:"Timestamp"`
	Level      string `json:"LogLevel"`
	Message    string `json:"Message"`
	LoggerName string `json:"Category"`
}

// CommonLogMessage interface
type CommonLogMessage interface {
	transform(w io.Writer)
}

// Exception for Java Log message
type Exception struct {
	RefID         int      `json:"refId"`
	ExceptionType string   `json:"exceptionType"`
	Message       string   `json:"message"`
	CausedBy      CausedBy `json:"causedBy"`
	Frames        *[]Frame `json:"frames"`
}

// CausedBy Exception caused by filed
type CausedBy struct {
	Exception *Exception `json:"exception,omitempty"`
}

// Frame for Uber zap log message
type Frame struct {
	Class  string `json:"class"`
	Method string `json:"method"`
	Line   int    `json:"line"`
}

// Tracing Log message in Java based logging
type Tracing struct {
	TraceID string `json:"traceId"`
	SpanID  string `json:"spanId"`
	Sampled string `json:"sampled"`
}

func (lm *QuarkusLogMessage) transform(w io.Writer) {
	// log contains a tracing message
	if lm.Tracing != (Tracing{}) {
		if strings.EqualFold(lm.Level, "Trace") {
			fmt.Fprintf(w, "\033[36m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Debug") {
			fmt.Fprintf(w, "\033[35m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Info") {
			fmt.Fprintf(w, "\033[32m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Warn") {
			fmt.Fprintf(w, "\033[33m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Error") {
			fmt.Fprintf(w, "\033[31m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Severe") {
			fmt.Fprintf(w, "\033[91m%v \033[39m%v\t\033[36mtraceId=%v \033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.Tracing.TraceID, lm.LoggerName, lm.Message)
		}
	} else {
		if strings.EqualFold(lm.Level, "Trace") {
			fmt.Fprintf(w, "\033[36m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Debug") {
			fmt.Fprintf(w, "\033[35m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Info") {
			fmt.Fprintf(w, "\033[32m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Warn") {
			fmt.Fprintf(w, "\033[33m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Error") {
			fmt.Fprintf(w, "\033[31m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
		if strings.EqualFold(lm.Level, "Severe") {
			fmt.Fprintf(w, "\033[91m%v \033[39m%v\t\033[37m%v\t\033[35m%v\n\033[39m", lm.Level, lm.Timestamp, lm.LoggerName, lm.Message)
		}
	}

	// log message contains an error error
	if lm.Exception != (Exception{}) {
		lm.Exception.transform(w)
	}
}

func (ex *Exception) transform(w io.Writer) {
	fmt.Fprintf(w, "\033[31mCaused by: %v. %s:\n\033[31m", ex.ExceptionType, ex.Message)

	for _, frame := range *ex.Frames {
		fmt.Fprintf(w, "\t \033[31mat %s(%s:%v)\n\033[39m", frame.Method, frame.Class, frame.Line)
	}

	if ex.CausedBy != (CausedBy{}) {
		ex.CausedBy.Exception.transform(w)
	}
}

func (alm *SpringBootLogMessage) transform(w io.Writer) {
	fmt.Fprintf(w, "%v %v\t%v\t%v\n", alm.Level, alm.Timestamp, alm.LoggerName, alm.Message)
	// log message contains an error error
	if alm.Exception != "" {
		fmt.Fprintf(w, "Exception: %s", alm.Exception)
	}
}

func (glm *GoZapLogMessage) transform(w io.Writer) {
	sec, dec := math.Modf(glm.Timestamp)
	timestamp := time.Unix(int64(sec), int64(dec*(1e9))).UTC()
	fmt.Fprintf(w, "%v %v\t%v\tmsg: %v\tcontroller: %v\trequest: %v\n", glm.Level, timestamp, glm.Logger, glm.Message, glm.Controller, glm.Request)
	// log message contains an error error
	if glm.Error != "" {
		fmt.Fprintf(w, "\033[31merror: %s\033[39m", glm.Error)
		fmt.Fprintf(w, "\033[31mstacktrace: %s\n\033[39m", glm.Stacktrace)
	}
}

func (dnlm *DotNetLogMessage) transform(w io.Writer) {
	fmt.Fprintf(w, "%v %v\t%v\t%v\n", dnlm.Level, dnlm.Timestamp, dnlm.LoggerName, dnlm.Message)
}
