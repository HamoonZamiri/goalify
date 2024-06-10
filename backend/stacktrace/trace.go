package stacktrace

import "fmt"

type TraceLogger interface {
	GetTrace(function string) string
}

type DomainStackTraceLogger struct {
	domain string
}

func NewDomainStackTraceLogger(domain string) *DomainStackTraceLogger {
	return &DomainStackTraceLogger{domain: domain}
}

func (d *DomainStackTraceLogger) GetTrace(function string) string {
	return fmt.Sprintf("domain=%s: %s", d.domain, function)
}
