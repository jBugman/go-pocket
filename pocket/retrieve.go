package pocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	rawCacheFile = "cache.json"
	dumpFile     = "dump.json"
)

func (this *Api) Retrieve() (items []Item) {
	request := Request{
		Api:         *this,
		State:       "all",
		ContentType: "article",
		Count:       3,
		DetailType:  "complete",
	}

	data := this.doRequest(request)
	response := this.parseResponse(bytes.NewReader(data))
	items = response.getItems()
	return
}

// Cached retrieve for testing
func (this *Api) TestRetrieve() (items []Item) {
	request := Request{
		Api:         *this,
		State:       "all",
		ContentType: "article",
		Count:       3,
		DetailType:  "complete",
	}

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

func (this *Api) doRequest(request Request) []byte {
	requestBody, _ := json.Marshal(request)
	response, _ := http.Post(retrieveUrl, "application/json", bytes.NewReader(requestBody))
	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return bodyBytes
}

func (this *Api) parseResponse(body io.Reader) (response RetrieveResponse) {
	jsonParser := json.NewDecoder(body)
	_ = jsonParser.Decode(&response)
	return response
}

type Request struct {
	Api
	State string `json:"state,omitempty"`
	// unread = only return unread items (default)
	// archive = only return archived items
	// all = return both unread and archived items
	Favorite int `json:"favorite,omitempty"`
	// 0 = only return un-favorited items
	// 1 = only return favorited items
	Tag string `json:"tag,omitempty"`
	// tag_name = only return items tagged with tag_name
	// _untagged_ = only return untagged items
	ContentType string `json:"contentType,omitempty"`
	// article = only return articles
	// video = only return videos or articles with embedded videos
	// image = only return images
	Sort string `json:"sort,omitempty"`
	// newest = return items in order of newest to oldest
	// oldest = return items in order of oldest to newest
	// title = return items in order of title alphabetically
	// site = return items in order of url alphabetically
	DetailType string `json:"detailType,omitempty"`
	// simple = only return the titles and urls of each item
	// complete = return all data about each item, including tags, images, authors, videos and more
	Search string `json:"search,omitempty"`
	// Only return items whose title or url contain the search string
	Domain string `json:"domain,omitempty"`
	// Only return items from a particular domain
	/* since	timestamp		Only return items modified since the given since unix timestamp */
	Count int `json:"count,omitempty"`
	// Only return count number of items
	Offset int `json:"offset,omitempty"`
	// Used only with count; start returning from offset position of results
}

type RetrieveResponse struct {
	Status   int              `json:"status"`
	Complete int              `json:"complete"`
	List     InternalItemList `json:"list"`
	Since    int              `json:"since"`
}

func (this *RetrieveResponse) getItems() []Item {
	return this.List.Values()
}

func (this *RetrieveResponse) dump(filename string) {
	pretty, _ := json.MarshalIndent(this, "", "	")
	ioutil.WriteFile(filename, pretty, 0666)
}

type InternalItemList map[string]Item

func (this *InternalItemList) Values() (result []Item) {
	result = make([]Item, 0, len(*this))

	for _, value := range *this {
		result = append(result, value)
	}
	return
}

type Item struct {
	Id string `json:"item_id"`
	// A unique identifier matching the saved item. This id must be used to perform any actions through the v3/modify endpoint.
	ResolvedId string `json:"resolved_id"`
	// A unique identifier similar to the item_id but is unique to the actual url of the saved item. The resolved_id identifies unique urls. For example a direct link to a New York Times article and a link that redirects (ex a shortened bit.ly url) to the same article will share the same resolved_id. If this value is 0, it means that Pocket has not processed the item. Normally this happens within seconds but is possible you may request the item before it has been resolved.
	GivenUrl string `json:"given_url"`
	// The actual url that was saved with the item. This url should be used if the user wants to view the item.
	Url string `json:"resolved_url"`
	// The final url of the item. For example if the item was a shortened bit.ly link, this will be the actual article the url linked to.
	GivenTitle string `json:"given_title"`
	// The title that was saved along with the item.
	Title string `json:"resolved_title"`
	// The title that Pocket found for the item when it was parsed
	Favorite int `json:"favorite"`
	// 0 or 1 - 1 If the item is favorited
	Status int `json:"status"`
	// 0, 1, 2 - 1 if the item is archived - 2 if the item should be deleted
	Excerpt string `json:"excerpt"`
	// The first few lines of the item (articles only)
	IsArticle int `json:"is_article"`
	// 0 or 1 - 1 if the item is an article
	HasImage int `json:"has_image"`
	// 0, 1, or 2 - 1 if the item has images in it - 2 if the item is an image
	HasVideo int `json:"has_video"`
	// 0, 1, or 2 - 1 if the item has videos in it - 2 if the item is a video
	WordCount int `json:"word_count"`
	// How many words are in the article
	Tags map[string]Tag `json:"tags,omitempty"`
	// A JSON object of the user tags associated with the item
	Authors map[string]Author `json:"authors,omitempty"`
	// A JSON object listing all of the authors associated with the item
	Images map[string]interface{} `json:"images,omitempty"`
	// A JSON object listing all of the images associated with the item
	Videos map[string]interface{} `json:"videos,omitempty"`
	// A JSON object listing all of the videos associated with the item
}

func (this *Item) GetTags() []string {
	tags := (*this).Tags
	result := make([]string, 0, len(tags))

	for _, value := range tags {
		result = append(result, value.Name)
	}
	return result
}

type Tag struct {
	ItemId string `json:"item_id"`
	Name   string `json:"tag"`
}

type Author struct {
	ItemId string `json:"item_id"`
	Id     string `json:"author_id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
}
