package load

import (
	"brainbaking.com/go-jamming/app/mf"
	"io/ioutil"
	"path"
)

// FromDisk assumes that params have already been validated.
func FromDisk(domain string, dataPath string) mf.IndiewebDataResult {
	loadPath := path.Join(dataPath, domain)

	info, _ := ioutil.ReadDir(loadPath)
	amountOfFiles := len(info)
	results := make(chan *mf.IndiewebData, amountOfFiles)

	for _, file := range info {
		go func(fileName string) {
			results <- mf.RequireFromFile(path.Join(loadPath, fileName))
		}(file.Name())
	}

	indiewebResults := gather(amountOfFiles, results)
	return mf.WrapResult(indiewebResults)
}

func gather(amount int, results <-chan *mf.IndiewebData) []*mf.IndiewebData {
	var indiewebResults []*mf.IndiewebData
	for i := 0; i < amount; i++ {
		result := <-results
		if !result.IsEmpty() {
			indiewebResults = append(indiewebResults, result)
		}
	}
	return indiewebResults
}
