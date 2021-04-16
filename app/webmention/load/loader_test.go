package load

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

// stress tests to see what concurrent disk access is like. Runs fine, even with 5000 runs and 100 files.
// this means worker pools do not have to be implemented in FromDisk().
// However, if runs := 10000, some results are empty. At other times, even ioutil.ReadDir() panics...
// The rate limiter should catch this, combined with a domain read lock in the caller.
// Furthermore, a run of 1 and files of 50k breaks the OS without using a semaphore to limit the nr. of open files!
func TestFromDiskStressTest(t *testing.T) {
	runs := 100
	files := 100

	os.MkdirAll("testdata/somedomain", os.ModePerm)
	defer os.RemoveAll("testdata")

	for i := 0; i < files; i++ {
		json := `{"author":{"name":"Jef Klakveld","picture":"https://brainbaking.com/img/avatar.jpg"},"name":"I much prefer Sonic Mania's Lock On to Belgium's t...","content":"I much prefer Sonic Mania’s Lock On to Belgium’s third Lock Down. Sigh. At least 16-bit 2D platformers make me smile: https://jefklakscodex.com/articles/reviews/sonic-mania/\n\n\n\nEnclosed Toot image","published":"2021-03-25T10:45:00","url":"https://brainbaking.com/notes/2021/03/25h10m45s09/","type":"mention","source":"https://brainbaking.com/notes/2021/03/25h10m45s09/","target":"https://jefklakscodex.com/articles/reviews/sonic-mania/"}`
		ioutil.WriteFile(fmt.Sprintf("testdata/somedomain/%d.json", i), []byte(json), os.ModePerm)
	}

	amounts := make(chan int, runs)
	for i := 0; i < runs; i++ {
		go func(nr int) {
			data := FromDisk("somedomain", "testdata")
			itms := len(data.Data)

			fmt.Printf("From disk #%d - found %d items\n", nr, itms)
			amounts <- itms
		}(i)
	}

	fmt.Println("Asserting...")
	for i := 0; i < runs; i++ {
		actual := <-amounts
		assert.Equal(t, files, actual)
	}
}

func TestFromDiskReturnsAllJsonFilesFromDiskWrappedInResult(t *testing.T) {
	os.MkdirAll("testdata/somedomain", os.ModePerm)
	defer os.RemoveAll("testdata")

	json1 := `{"author":{"name":"Wouter Groeneveld","picture":"https://brainbaking.com/img/avatar.jpg"},"name":"I much prefer Sonic Mania's Lock On to Belgium's t...","content":"I much prefer Sonic Mania’s Lock On to Belgium’s third Lock Down. Sigh. At least 16-bit 2D platformers make me smile: https://jefklakscodex.com/articles/reviews/sonic-mania/\n\n\n\nEnclosed Toot image","published":"2021-03-25T10:45:00","url":"https://brainbaking.com/notes/2021/03/25h10m45s09/","type":"mention","source":"https://brainbaking.com/notes/2021/03/25h10m45s09/","target":"https://jefklakscodex.com/articles/reviews/sonic-mania/"}`
	json2 := `{"author":{"name":"Jef Klakveld","picture":"https://brainbaking.com/img/avatar.jpg"},"name":"I much prefer Sonic Mania's Lock On to Belgium's t...","content":"I much prefer Sonic Mania’s Lock On to Belgium’s third Lock Down. Sigh. At least 16-bit 2D platformers make me smile: https://jefklakscodex.com/articles/reviews/sonic-mania/\n\n\n\nEnclosed Toot image","published":"2021-03-25T10:45:00","url":"https://brainbaking.com/notes/2021/03/25h10m45s09/","type":"mention","source":"https://brainbaking.com/notes/2021/03/25h10m45s09/","target":"https://jefklakscodex.com/articles/reviews/sonic-mania/"}`
	ioutil.WriteFile("testdata/somedomain/testjson1.json", []byte(json1), os.ModePerm)
	ioutil.WriteFile("testdata/somedomain/testjson2.json", []byte(json2), os.ModePerm)

	result := FromDisk("somedomain", "testdata")
	sort.SliceStable(result.Data, func(i, j int) bool {
		comp := strings.Compare(result.Data[i].Author.Name, result.Data[j].Author.Name)
		if comp > 0 {
			return false
		}
		return true
	})

	assert.Equal(t, "success", result.Status)
	assert.Equal(t, "Jef Klakveld", result.Data[0].Author.Name)
	assert.Equal(t, "Wouter Groeneveld", result.Data[1].Author.Name)
}
