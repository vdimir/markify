package testutil

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MustReadData open and read file checking errors
func MustReadData(t *testing.T, path string) []byte {
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return data
}

// GetTempFolder return tmp forder and cleanup func
func GetTempFolder(t require.TestingT, prefix string) (tmpPath string, cleanup func()) {
	tmpPath, err := ioutil.TempDir("", prefix)
	require.NoError(t, err)

	cleanup = func() {
		defer func() {
			assert.NoError(t, os.RemoveAll(tmpPath))
		}()
	}
	return tmpPath, cleanup
}

// ChooseRandomUnusedPort returns random free port
func ChooseRandomUnusedPort() (port uint16) {
	for i := 0; i < 10; i++ {
		port = 40000 + uint16(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port)); err == nil {
			_ = ln.Close()
			time.Sleep(time.Millisecond * 10)
			return port
		}
	}
	return 0
}

// WaitForHTTPSServerStart wait up to 3 second to server start
func WaitForHTTPSServerStart(host string, port uint16) error {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	for i := 0; i < 300; i++ {
		time.Sleep(time.Millisecond * 10)
		conn, _ := net.DialTimeout("tcp", hostPort, time.Millisecond*10)
		if conn != nil {
			_ = conn.Close()
			return nil
		}
	}
	return errors.Errorf("cannot dial %s", hostPort)
}
