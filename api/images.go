package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/uthoterraform/terraform-provider-utho/helper"
)

type Images struct {
	Images []Image `json:"images"`
}

type Image struct {
	Distro       string `json:"distro"`
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
	Image        string `json:"image"`
	Cost         int    `json:"cost"`
}

func (c *Client) GetImages(ctx context.Context) (Images, error) {
	uri := BASE_URL + "cloud/images"

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Images{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return Images{}, errors.New("unexpected http error code received for geting images data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var images Images
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return Images{}, err
	}

	return images, nil
}
