
package webmention

import (
	"testing"
	"os"
	"errors"

	"github.com/wgroeneveld/go-jamming/common"
	"github.com/wgroeneveld/go-jamming/mocks"
)

var conf = &common.Config{
	AllowedWebmentionSources: []string {
		"jefklakscodex.com",
	},
	DataPath: "testdata",
}


func TestConvertWebmentionToPath(t *testing.T) {
	wm := webmention{
		source: "https://brainbaking.com",
		target: "https://jefklakscodex.com/articles",
	}

	result := wm.asPath(conf)
	if result != "testdata/jefklakscodex.com/99be66594fdfcf482545fead8e7e4948.json" {
		t.Fatalf("md5 hash check failed, got " + result)
	}
}

func writeSomethingTo(filename string) {
	file, _ := os.Create(filename)
	file.WriteString("lolz")
	defer file.Close()	
}

func TestReceiveTargetDoesNotExistAnymoreDeletesPossiblyOlderWebmention(t *testing.T) {
	os.MkdirAll("testdata/jefklakscodex.com", os.ModePerm)
	defer os.RemoveAll("testdata")

	wm := webmention{
		source: "https://brainbaking.com",
		target: "https://jefklakscodex.com/articles",
	}
	filename := wm.asPath(conf)
	writeSomethingTo(filename)

	client := &mocks.RestClientMock{
		GetBodyFunc: func(url string) (string, error) {
			return "", errors.New("whoops")
		},
	}	
	receiver := &receiver {
		conf: conf,
		restClient: client,
	}

	receiver.receive(wm)
  	if _, err := os.Stat(filename); err == nil {
  		t.Fatalf("Expected possibly older webmention to be deleted, but it wasn't!")
  	}
}

func TestProcessSourceBodyAbortsIfNoMentionOfTargetFoundInSourceHtml(t *testing.T) {
	os.MkdirAll("testdata/jefklakscodex.com", os.ModePerm)
	defer os.RemoveAll("testdata")

	wm := webmention{
		source: "https://brainbaking.com",
		target: "https://jefklakscodex.com/articles",
	}
	filename := wm.asPath(conf)

	receiver := &receiver {
		conf: conf,
	}

	receiver.processSourceBody("<html>my nice body</html>", wm)
  	if _, err := os.Stat(filename); err == nil {
  		t.Fatalf("Expected no file to be created!")
  	}
}

