package fetch

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

type delay struct {
	minTime time.Duration
	maxTime time.Duration
}

func (d delay) sleep() {
	rf := rand.Float32()
	sleepSec := d.minTime + time.Duration(float32(d.maxTime-d.minTime)*rf)

	<-time.After(sleepSec * time.Second)
}

// Mock emulates fetcher behaviour
type Mock struct {
	pages  map[string][]byte
	delays map[string]delay
}

// NewMock create new Mock
func NewMock() Fetcher {
	return &Mock{
		pages:  map[string][]byte{},
		delays: map[string]delay{},
	}
}

// Fetch waits specified time and return result if presented
func (f *Mock) Fetch(url string) (io.ReadCloser, error) {
	if dl, ok := f.delays[url]; ok {
		dl.sleep()
	}

	if data, ok := f.pages[url]; ok {
		return ioutil.NopCloser(bytes.NewReader(data)), nil
	}
	return nil, errors.Errorf("cannot get data from %q", url)
}

// SetDelay setup delay for fetching data from url
func (f *Mock) SetDelay(url string, minTime time.Duration, maxTime time.Duration) {
	f.delays[url] = delay{minTime, maxTime}
}

// SetData add data that will be fetched by url
func (f *Mock) SetData(url string, data []byte) {
	f.pages[url] = data
}
