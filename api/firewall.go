package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoterraform/terraform-provider-utho/helper"

	"net/http"
)

type Firewalls struct {
	Firewalls []Firewall `json:"firewalls"`
	Status    string     `json:"status"`
	Message   string     `json:"message"`
}

type Firewall struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreatedAt    string `json:"created_at"`
	Rulecount    string `json:"rulecount"`
	Serverscount string `json:"serverscount"`
}

type FirewallRequest struct {
	Name string `json:"name"`
}
type FirewallResponse struct {
	ID      string `json:"firewallid"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) CreateFirewall(ctx context.Context, firewallRequest FirewallRequest) (FirewallResponse, error) {
	uri := BASE_URL + "firewall/create"

	jsonPayload, err := json.Marshal(firewallRequest)
	if err != nil {
		return FirewallResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return FirewallResponse{}, err
	}

	defer resp.Body.Close()

	var firewall FirewallResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewall); err != nil {
		return FirewallResponse{}, err
	}
	if firewall.Status != "success" {
		return FirewallResponse{}, errors.New(firewall.Message)
	}

	return firewall, nil
}

func (c *Client) GetFirewall(ctx context.Context, id string) (Firewall, error) {
	uri := BASE_URL + "firewall/" + id

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Firewall{}, err
	}
	defer resp.Body.Close()
	var firewalls Firewalls
	if err := json.NewDecoder(resp.Body).Decode(&firewalls); err != nil {
		return Firewall{}, err
	}

	if len(firewalls.Firewalls) == 0 {
		return Firewall{}, errors.New("firewall not found")
	}

	return firewalls.Firewalls[0], nil
}

func (c *Client) DeleteFirewall(ctx context.Context, id string) error {
	uri := BASE_URL + "firewall/" + id + "/destroy"

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var firewall FirewallResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewall); err != nil {
		return err
	}
	if firewall.Status != "success" {
		return errors.New(firewall.Message)
	}

	return nil
}
