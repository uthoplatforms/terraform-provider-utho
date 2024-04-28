package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoplatforms/terraform-provider-utho/helper"

	"net/http"
)

type Loadbalancers struct {
	Loadbalancers []Loadbalancer `json:"loadbalancers"`
}

type Loadbalancer struct {
	ID            string `json:"id"`
	Userid        string `json:"userid"`
	IP            string `json:"ip"`
	Name          string `json:"name"`
	Algorithm     string `json:"algorithm"`
	Cookie        string `json:"cookie"`
	Cookiename    string `json:"cookiename"`
	Redirecthttps string `json:"redirecthttps"`
	Type          string `json:"type"`
	Country       string `json:"country"`
	Cc            string `json:"cc"`
	City          string `json:"city"`
	Backendcount  string `json:"backendcount"`
	CreatedAt     string `json:"created_at"`
	Status        string `json:"status"`
	// Backends      string  `json:"backends"`
	Rules []Rules `json:"rules"`
	// Acls          string  `json:"acls"`
	// Routes        string  `json:"routes"`
	// ScalingGroups string  `json:"scaling_groups"`
	// Frontends     string  `json:"frontends"`
}
type Rules struct {
	ID          string `json:"id"`
	Lb          string `json:"lb"`
	SrcProto    string `json:"src_proto"`
	SrcPort     string `json:"src_port"`
	DstProto    string `json:"dst_proto"`
	DstPort     string `json:"dst_port"`
	Timeadded   string `json:"timeadded"`
	Timeupdated string `json:"timeupdated"`
}

type LoadbalancerRequest struct {
	Dcslug string `json:"dcslug"`
	Type   string `json:"type"`
	Name   string `json:"name"`
}
type LoadbalancerResponse struct {
	Status         string `json:"status"`
	Loadbalancerid string `json:"loadbalancerid"`
	Message        string `json:"message"`
}

func (c *Client) CreateLoadbalancer(ctx context.Context, loadbalancerRequest LoadbalancerRequest) (LoadbalancerResponse, error) {
	uri := BASE_URL + "loadbalancer"

	jsonPayload, err := json.Marshal(loadbalancerRequest)
	if err != nil {
		return LoadbalancerResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return LoadbalancerResponse{}, err
	}

	defer resp.Body.Close()

	var loadbalancer LoadbalancerResponse
	if err := json.NewDecoder(resp.Body).Decode(&loadbalancer); err != nil {
		return LoadbalancerResponse{}, err
	}
	if loadbalancer.Status != "success" {
		return LoadbalancerResponse{}, errors.New(loadbalancer.Message)
	}
	return loadbalancer, nil
}

func (c *Client) GetLoadbalancer(ctx context.Context, id string) (Loadbalancer, error) {
	uri := BASE_URL + "loadbalancer/" + id

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Loadbalancer{}, err
	}
	defer resp.Body.Close()

	var loadbalancers Loadbalancers
	if err := json.NewDecoder(resp.Body).Decode(&loadbalancers); err != nil {
		return Loadbalancer{}, err
	}

	if len(loadbalancers.Loadbalancers) == 0 {
		return Loadbalancer{}, errors.New("domain not found")
	}
	return loadbalancers.Loadbalancers[0], nil
}

func (c *Client) DeleteLoadbalancer(ctx context.Context, id string) error {
	uri := BASE_URL + "loadbalancer/" + id

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var loadbalancer LoadbalancerResponse
	if err := json.NewDecoder(resp.Body).Decode(&loadbalancer); err != nil {
		return err
	}
	if loadbalancer.Status != "success" {
		return errors.New(loadbalancer.Message)
	}

	return nil
}
