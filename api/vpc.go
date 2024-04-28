package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoplatforms/terraform-provider-utho/helper"

	"net/http"
)

type Vpc struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	Total     int    `json:"total"`
	Available int    `json:"available"`
	Network   string `json:"network"`
	Dcslug    string `json:"dcslug"`
	IsDefault string `json:"is_default"`
	Id        string `json:"id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

type VpcRequest struct {
	Dcslug  string `json:"dcslug"`
	Name    string `json:"name"`
	Planid  string `json:"planid"`
	Network string `json:"network"`
	Size    string `json:"size"`
}

func (c *Client) CreateVpc(ctx context.Context, vpcRequest VpcRequest) (Vpc, error) {
	uri := BASE_URL + "vpc/create"

	jsonPayload, err := json.Marshal(vpcRequest)
	if err != nil {
		return Vpc{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return Vpc{}, err
	}

	defer resp.Body.Close()

	var vpc Vpc
	if err := json.NewDecoder(resp.Body).Decode(&vpc); err != nil {
		return Vpc{}, err
	}
	if vpc.Status != "success" {
		return Vpc{}, errors.New(vpc.Message)
	}
	return vpc, nil
}

func (c *Client) GetVpc(ctx context.Context, id string) (Vpc, error) {
	uri := BASE_URL + "vpc/" + id

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Vpc{}, err
	}
	defer resp.Body.Close()

	var vpc Vpc
	if err := json.NewDecoder(resp.Body).Decode(&vpc); err != nil {
		return Vpc{}, err
	}

	if vpc.Status == "error" {
		return Vpc{}, errors.New(vpc.Message)
	}
	return vpc, nil
}

// ///////////////////////////////////////////////////////////////////
func (c *Client) DeleteVpc(ctx context.Context, id string) error {
	uri := BASE_URL + "vpc/" + id + "/destroy"

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var vpc Vpc
	if err := json.NewDecoder(resp.Body).Decode(&vpc); err != nil {
		return err
	}
	if vpc.Status != "success" {
		return errors.New(vpc.Message)
	}
	return nil
}
