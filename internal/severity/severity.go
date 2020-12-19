package severity

import "context"

// Severity is type of severity that extended int
type Severity int

// Enum for severity
const (
	Default   Severity = iota * 100 // 0
	Debug                           // 100
	Info                            // 200
	Notice                          // 300
	Warning                         // 400
	Error                           // 500
	Critical                        // 600
	Alert                           // 700
	Emergency                       // 800
)

var severityMap = map[Severity]string{
	Default:   "DEFAULT",
	Debug:     "DEBUG",
	Info:      "INFO",
	Notice:    "NOTICE",
	Warning:   "WARNING",
	Error:     "ERROR",
	Critical:  "CRITICAL",
	Alert:     "ALERT",
	Emergency: "EMERGENCY",
}

// String returns text representation for the severity
func (s Severity) String() string {
	return severityMap[s]
}

// ContextSeverity is list of severity that managed by context
type ContextSeverity struct {
	logged []Severity
}

// Add adds a list of severity
func (c *ContextSeverity) Add(s Severity) {
	c.logged = append(c.logged, s)
}

// Max returns max from list of severity
func (c *ContextSeverity) Max() Severity {
	max := Default
	for _, s := range c.logged {
		if s > max {
			max = s
		}
	}
	return max
}

type contextKey struct{}

var contextSeverityKey = &contextKey{}

// SetContextSeverity sets the ContextSeverity instance to context
func SetContextSeverity(ctx context.Context, cs *ContextSeverity) context.Context {
	return context.WithValue(ctx, contextSeverityKey, cs)
}

// GetContextSeverity gets the ContextSeverity instance from context
func GetContextSeverity(ctx context.Context) *ContextSeverity {
	v, _ := ctx.Value(contextSeverityKey).(*ContextSeverity)
	return v
}
