package severity

type Severity int

//nolint:varcheck,deadcode,unused
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
