package log

import "testing"

func TestLogger(t *testing.T) {
	SetLevel(InfoLevel)
	SetReportCaller(true)

	Debugf("should not be seen")

	WithContext(SetDebugModeContext(nil, true)).Debugf("should be seen")
}
