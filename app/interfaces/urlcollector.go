package interfaces

import (
	"time"
)

type UrlCollectorHandlerFactory interface {
	New() UrlCollectorHandler
}

type UrlCollectorHandler interface {
	Handle(date string) (string, error)
}

type UrlCollectorFactory interface {
	New(UrlCollectorHandler) UrlCollector
}

type UrlCollector interface {
	GetUrlList(from, to time.Time) ([]string, error)
}
