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

type ObjectStoragePlan struct {
	Pricing []Pricing `json:"pricing"`
	Rcode   string    `json:"rcode"`
}
type Pricing struct {
	ID             string `json:"id"`
	UUID           string `json:"uuid"`
	Type           string `json:"type"`
	Slug           string `json:"slug"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Disk           string `json:"disk"`
	RAM            string `json:"ram"`
	CPU            string `json:"cpu"`
	Bandwidth      string `json:"bandwidth"`
	IsFeatured     string `json:"is_featured"`
	DedicatedVcore string `json:"dedicated_vcore"`
	Price          string `json:"price"`
	Monthly        string `json:"monthly"`
}

func (c *Client) GetObjectStoragePlan(ctx context.Context) (ObjectStoragePlan, error) {
	uri := BASE_URL + "pricing/objectstorage"

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return ObjectStoragePlan{}, err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return ObjectStoragePlan{}, errors.New("unexpected http error code received for geting Object Storage Plan data status code :" + strconv.Itoa(resp.StatusCode) + " body" + string(body))
	}
	defer resp.Body.Close()

	var objectStoragePlan ObjectStoragePlan
	if err := json.NewDecoder(resp.Body).Decode(&objectStoragePlan); err != nil {
		return ObjectStoragePlan{}, err
	}

	return objectStoragePlan, nil
}
