package output

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/benthosdev/benthos/v4/internal/component/metrics"
	"github.com/benthosdev/benthos/v4/internal/log"
	"github.com/benthosdev/benthos/v4/internal/manager/mock"
	"github.com/benthosdev/benthos/v4/internal/message"
	"github.com/benthosdev/benthos/v4/internal/old/processor"
)

func TestFallbackOutputBasic(t *testing.T) {
	dir := t.TempDir()

	outOne, outTwo, outThree := NewConfig(), NewConfig(), NewConfig()
	outOne.Type, outTwo.Type, outThree.Type = TypeHTTPClient, TypeFile, TypeFile
	outOne.HTTPClient.URL = "http://localhost:11111111/badurl"
	outOne.HTTPClient.NumRetries = 1
	outOne.HTTPClient.Retry = "1ms"
	outTwo.File.Path = filepath.Join(dir, "two", `bar-${!count("fallbacktofoo")}-${!count("fallbacktobar")}.txt`)
	outTwo.File.Codec = "all-bytes"
	outThree.File.Path = "/dev/null"

	procOne, procTwo, procThree := processor.NewConfig(), processor.NewConfig(), processor.NewConfig()
	procOne.Type, procTwo.Type, procThree.Type = processor.TypeBloblang, processor.TypeBloblang, processor.TypeBloblang
	procOne.Bloblang = `root = "this-should-never-appear %v".format(count("fallbacktofoo")) + content()`
	procTwo.Bloblang = `root = "two-" + content()`
	procThree.Bloblang = `root = "this-should-never-appear %v".format(count("fallbacktobar")) + content()`

	outOne.Processors = append(outOne.Processors, procOne)
	outTwo.Processors = append(outTwo.Processors, procTwo)
	outThree.Processors = append(outThree.Processors, procThree)

	conf := NewConfig()
	conf.Type = TypeFallback
	conf.Fallback = append(conf.Fallback, outOne, outTwo, outThree)

	s, err := New(conf, mock.NewManager(), log.Noop(), metrics.Noop())
	require.NoError(t, err)

	sendChan := make(chan message.Transaction)
	resChan := make(chan error)
	require.NoError(t, s.Consume(sendChan))

	t.Cleanup(func() {
		s.CloseAsync()
		require.NoError(t, s.WaitForClose(time.Second))
	})

	inputs := []string{
		"first", "second", "third", "fourth",
	}
	expFiles := map[string]string{
		"./two/bar-2-1.txt": "two-first",
		"./two/bar-4-2.txt": "two-second",
		"./two/bar-6-3.txt": "two-third",
		"./two/bar-8-4.txt": "two-fourth",
	}

	for _, input := range inputs {
		testMsg := message.QuickBatch([][]byte{[]byte(input)})
		select {
		case sendChan <- message.NewTransaction(testMsg, resChan):
		case <-time.After(time.Second * 2):
			t.Fatal("Action timed out")
		}

		select {
		case res := <-resChan:
			if res != nil {
				t.Fatal(res)
			}
		case <-time.After(time.Second * 2):
			t.Fatal("Action timed out")
		}
	}

	for k, exp := range expFiles {
		k = filepath.Join(dir, k)
		fileBytes, err := os.ReadFile(k)
		if err != nil {
			t.Errorf("Expected file '%v' could not be read: %v", k, err)
			continue
		}
		if act := string(fileBytes); exp != act {
			t.Errorf("Wrong contents for file '%v': %v != %v", k, act, exp)
		}
	}
}