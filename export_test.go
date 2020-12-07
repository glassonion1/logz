package logz

import (
	"time"

	"github.com/glassonion1/logz/internal/logger"
)

func SetNow(t time.Time) {
	logger.NowFunc = func() time.Time {
		return t
	}
}
