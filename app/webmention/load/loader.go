package load

import (
	"brainbaking.com/go-jamming/app/mf"
	"encoding/json"
	"io/ioutil"
	"path"
)

func FromDisk(domain string, dataPath string) mf.IndiewebDataResult {
	// assume that params have already been validated.
	loadPath := path.Join(dataPath, domain)
	info, _ := ioutil.ReadDir(loadPath)
	amountOfFiles := len(info)
	results := make(chan *mf.IndiewebData, amountOfFiles)

	for _, file := range info {
		fileName := file.Name()
		go func() {
			data, _ := ioutil.ReadFile(path.Join(loadPath, fileName))
			indiewebData := &mf.IndiewebData{}
			json.Unmarshal(data, indiewebData)
			results <- indiewebData
		}()
	}

	indiewebResults := gather(amountOfFiles, results)
	return mf.WrapResult(indiewebResults)
}

func gather(amount int, results chan *mf.IndiewebData) []*mf.IndiewebData {
	var indiewebResults []*mf.IndiewebData
	for i := 0; i < amount; i++ {
		result := <-results
		// json marshal errors are ignored in the above scatter func.Highly unlikely, but still.
		if result.Url != "" {
			indiewebResults = append(indiewebResults, result)
		}
	}
	return indiewebResults
}
