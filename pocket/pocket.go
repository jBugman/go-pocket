package pocket

const (
	retrieveUrl     = "https://getpocket.com/v3/get"
	tokenRequestUrl = "https://getpocket.com/v3/oauth/request"
	authorizeUrl    = "https://getpocket.com/v3/oauth/authorize"
)

type Api struct {
	ConsumerKey string `json:"consumer_key"`
	AccessToken string `json:"access_token"`
}
