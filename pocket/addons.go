package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	rawCacheFile = "cache.json"
	dumpFile     = "dump.json"
)

func (this *RetrieveResponse) dump(filename string) {
	pretty, _ := json.MarshalIndent(this, "", "	")
	ioutil.WriteFile(filename, pretty, 0666)
}

func (this *Api) RetrieveAllArticles() (items []Item) {
	request := Request{
		State:       "all",
		ContentType: "article",
		DetailType:  "complete",
		Sort:        "newest",
	}
	return this.CachedRetrieve(request)
}

/* Cached retrieve for testing */
func (this *Api) CachedRetrieve(request Request) (items []Item) {
	request.Api = *this

	data, err := ioutil.ReadFile(rawCacheFile)
	if err != nil {
		fmt.Println("Requsting API")
		data = this.doRequest(request)
		ioutil.WriteFile(rawCacheFile, data, 0666)
	} else {
		fmt.Println("Loaded from cache")
	}

	response := this.parseResponse(bytes.NewReader(data))
	response.dump(dumpFile)
	items = response.getItems()
	return
}
