package collector

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"url-collector/daterange"
	"url-collector/interfaces"
)

const defaultConcurrentRequests = 5

var concurrentRequests = defaultConcurrentRequests

type BaseUrlCollectorFactory struct {
}

type BaseUrlCollector struct {
	accessChannel    chan struct{}
	collectorHandler interfaces.UrlCollectorHandler
}

func (bucf BaseUrlCollectorFactory) New(collectorHandler interfaces.UrlCollectorHandler) interfaces.UrlCollector {
	val, convErr := strconv.Atoi(os.Getenv("CONCURRENT_REQUESTS"))
	if convErr == nil && val > 0 {
		concurrentRequests = val
	}

	collector := BaseUrlCollector{
		accessChannel:    make(chan struct{}, concurrentRequests),
		collectorHandler: collectorHandler,
	}
	return &collector
}

func handleCollection(date string) error {
	time.Sleep(1 * time.Second)
	fmt.Println(date)
	return nil
}

func (buc BaseUrlCollector) getPictureUrlLink(wg *sync.WaitGroup, goroutineAccessChannel *chan struct{}, date string, urlChannel *chan string, errorChannel *chan error) {
	defer wg.Done()
	*goroutineAccessChannel <- struct{}{}
	defer func() {
		<-*goroutineAccessChannel
	}()

	url, err := buc.collectorHandler.Handle(date)
	if err != nil {
		*errorChannel <- err
		return
	}
	*urlChannel <- url
}

func collectPictureUrls(urlChannel *chan string, urls *[]string, doneChannel *chan struct{}, errorChannel *chan error, err *error) {
	for {
		select {
		case url := <-*urlChannel:
			*urls = append(*urls, url)
		case e := <-*errorChannel:
			if *err == nil {
				*err = e
			}
		case <-*doneChannel:
			return
		}
	}
}

func (buc BaseUrlCollector) GetUrlList(from, to time.Time) ([]string, error) {
	urlChannel := make(chan string)
	collectorDoneChannel := make(chan struct{})
	errorChannel := make(chan error)
	var collectionError error
	var urlList []string
	var wgRequests sync.WaitGroup

	go collectPictureUrls(&urlChannel, &urlList, &collectorDoneChannel, &errorChannel, &collectionError)

	currentDate := from
	for currentDate.Unix() <= to.Unix() {
		wgRequests.Add(1)
		go buc.getPictureUrlLink(&wgRequests, &buc.accessChannel, daterange.TimeToDayString(currentDate), &urlChannel, &errorChannel)
		currentDate = currentDate.Add(24 * time.Hour)
	}
	wgRequests.Wait()
	collectorDoneChannel <- struct{}{}
	return urlList, collectionError
}
