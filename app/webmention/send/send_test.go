package send

import (
	"brainbaking.com/go-jamming/app/mf"
	"brainbaking.com/go-jamming/common"
	"brainbaking.com/go-jamming/mocks"
	"brainbaking.com/go-jamming/rest"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"
)

func TestSinceForDomain(t *testing.T) {
	cases := []struct {
		label        string
		sinceInParam string
		sinceInFile  string
		expected     time.Time
	}{
		{
			"Is since parameter if provided",
			"2021-03-09T15:51:43.732Z",
			"",
			time.Date(2021, time.March, 9, 15, 51, 43, 732, time.UTC),
		},
		{
			"Is file contents if since parameter is empty and file is not",
			"",
			"2021-03-09T15:51:43.732Z",
			time.Date(2021, time.March, 9, 15, 51, 43, 732, time.UTC),
		},
		{
			"Is empty time if both parameter and file are not present",
			"",
			"",
			time.Time{},
		},
	}

	snder := Sender{
		Conf: &common.Config{
			DataPath: "testdata",
		},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			os.MkdirAll("testdata", os.ModePerm)
			defer os.RemoveAll("testdata")

			if tc.sinceInFile != "" {
				ioutil.WriteFile("testdata/domain-since.txt", []byte(tc.sinceInFile), os.ModePerm)
			}

			actual := snder.sinceForDomain("domain", tc.sinceInParam)
			assert.Equal(t, tc.expected.Year(), actual.Year())
			assert.Equal(t, tc.expected.Month(), actual.Month())
			assert.Equal(t, tc.expected.Day(), actual.Day())
			assert.Equal(t, tc.expected.Hour(), actual.Hour())
			assert.Equal(t, tc.expected.Minute(), actual.Minute())
			assert.Equal(t, tc.expected.Second(), actual.Second())
		})
	}
}

func TestSendMentionAsWebmention(t *testing.T) {
	passedFormValues := url.Values{}
	snder := Sender{
		RestClient: &mocks.RestClientMock{
			PostFormFunc: func(endpoint string, formValues url.Values) error {
				passedFormValues = formValues
				return nil
			},
		},
	}

	sendMentionAsWebmention(&snder, mf.Mention{
		Source: "mysource",
		Target: "mytarget",
	}, "someendpoint")

	assert.Equal(t, "mysource", passedFormValues.Get("source"))
	assert.Equal(t, "mytarget", passedFormValues.Get("target"))
}

// Stress test for opening HTTP connections en masse.
// Works great for up to 1000 runs. 10k hits: "http: Accept error: accept tcp [::]:6666: accept: too many open files in system; retrying in 10ms"
// Crashed even GoLand and the open Spotify client...
// The rate limiter fixes this, and in reality, we never send out 10k links anyway.
func TestSendMentionIntegrationStressTest(t *testing.T) {
	snder := Sender{
		Conf:       common.Configure(),
		RestClient: &rest.HttpClient{},
	}
	defer os.RemoveAll("data")

	runs := 100
	responses := make(chan bool, runs)

	http.HandleFunc("/pingback", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("pingbacked stuff."))
		responses <- true
	})
	http.HandleFunc("/target", func(writer http.ResponseWriter, request *http.Request) {
		target := `<html>
							<head>
								<link rel="pingback" href="http://localhost:6666/pingback" />
							</head>
							<body>sup!</body>
						</html>
						`
		writer.WriteHeader(200)
		writer.Write([]byte(target))
	})

	go func() {
		fmt.Println("Serving stub at 6666...")
		http.ListenAndServe(":6666", nil)
		fmt.Println("Stub stopped?")
	}()

	fmt.Println("Bootstrapping runs...")
	for i := 0; i < runs; i++ {
		snder.sendMention(mf.Mention{
			Source: "http://localhost:6666/source",
			Target: "http://localhost:6666/target",
		})
	}
	fmt.Println("Asserting...")
	for i := 0; i < runs; i++ {
		pingbacked := <-responses
		assert.True(t, pingbacked)
	}
}

func TestSendIntegrationTestCanSendBothWebmentionsAndPingbacks(t *testing.T) {
	posted := map[string]interface{}{}
	var lock = sync.Mutex{}
	defer os.RemoveAll("data")

	snder := Sender{
		Conf: common.Configure(),
		RestClient: &mocks.RestClientMock{
			GetBodyFunc: mocks.RelPathGetBodyFunc(t, "./../../../mocks/"),
			PostFunc: func(url string, contentType string, body string) error {
				lock.Lock()
				defer lock.Unlock()
				posted[url] = body
				return nil
			},
			PostFormFunc: func(endpoint string, formValues url.Values) error {
				lock.Lock()
				defer lock.Unlock()
				posted[endpoint] = formValues
				return nil
			},
		},
	}

	snder.Send("brainbaking.com", "2021-03-16T16:00:00.000Z")
	assert.Equal(t, 3, len(posted))

	wmPost1 := posted["http://aaronpk.example/webmention-endpoint-header"].(url.Values)
	assert.Equal(t, "https://brainbaking.com/notes/2021/03/16h17m07s14/", wmPost1.Get("source"))
	assert.Equal(t, "https://brainbaking.com/link-discover-test-multiple.html", wmPost1.Get("target"))

	wmPost2 := posted["http://aaronpk.example/pingback-endpoint-body"].(string)
	expectedPost2 := `<?xml version="1.0" encoding="UTF-8"?>
<methodCall>
	<methodName>pingback.ping</methodName>
	<params>
		<param>
			<value><string>https://brainbaking.com/notes/2021/03/16h17m07s14/</string></value>
		</param>
		<param>
			<value><string>https://brainbaking.com/pingback-discover-test-single.html</string></value>
		</param>
	</params>
</methodCall>`
	assert.Equal(t, expectedPost2, wmPost2)

	wmPost3 := posted["http://aaronpk.example/webmention-endpoint-body"].(url.Values)
	assert.Equal(t, "https://brainbaking.com/notes/2021/03/16h17m07s14/", wmPost3.Get("source"))
	assert.Equal(t, "https://brainbaking.com/link-discover-test-single.html", wmPost3.Get("target"))
}
