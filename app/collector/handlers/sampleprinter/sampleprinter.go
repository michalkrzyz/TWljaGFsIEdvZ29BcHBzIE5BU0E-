package sampleprinter

import (
	"fmt"
	"time"

	"url-collector/interfaces"
)

type SamplePrinterUrlCollectorHandlerFactory struct {
}

type SamplePrinterUrlCollectorHandler struct {
}

func (spucf SamplePrinterUrlCollectorHandlerFactory) New() interfaces.UrlCollectorHandler {
	return SamplePrinterUrlCollectorHandler{}
}

func (spuc SamplePrinterUrlCollectorHandler) Handle(date string) (string, error) {
	time.Sleep(1 * time.Second)
	fmt.Println(date)
	return date, nil
}
