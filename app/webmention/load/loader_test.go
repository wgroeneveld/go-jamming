package load

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

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
