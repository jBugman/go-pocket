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

func (this *Api) RetrieveAllArticles() (items []Item, err error) {
	request := Request{
		State:       "all",
		ContentType: "article",
		DetailType:  "complete",
		Sort:        "newest",
	}
	return this.CachedRetrieve(request)
}

/* Cached retrieve for testing */
func (this *Api) CachedRetrieve(request Request) (items []Item, err error) {
	request.Api = *this

	data, err := ioutil.ReadFile(rawCacheFile)
	if err != nil {
		fmt.Println("Requsting API")
		data, err = this.doRequest(request)
		if err != nil {
			return nil, err
		}
		ioutil.WriteFile(rawCacheFile, data, 0666)
	} else {
		fmt.Println("Loaded from cache")
	}

	response, err := this.parseResponse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	response.dump(dumpFile)
	items = response.getItems()
	return
}
