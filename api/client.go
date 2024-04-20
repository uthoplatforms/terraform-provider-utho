package api

const BASE_URL = "https://api.utho.com/v2/"

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
	}
}
