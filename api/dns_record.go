package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoterraform/terraform-provider-utho/helper"

	"net/http"
)

type DnsRecordRequest struct {
	Domain   string `json:"domain"`
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	Value    string `json:"value"`
	TTL      string `json:"ttl"`
	Porttype string `json:"porttype"`
	Port     string `json:"port"`
	Priority string `json:"priority"`
	Wight    string `json:"wight"`
}

type DnsRecordResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) CreateDnsRecord(ctx context.Context, dnsRecordRequest DnsRecordRequest) (DnsRecordResponse, error) {
	uri := BASE_URL + "dns/" + dnsRecordRequest.Domain + "/record/add"

	jsonPayload, err := json.Marshal(dnsRecordRequest)
	if err != nil {
		return DnsRecordResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return DnsRecordResponse{}, err
	}

	defer resp.Body.Close()

	var dnsRecord DnsRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&dnsRecord); err != nil {
		return DnsRecordResponse{}, err
	}
	if dnsRecord.Status != "success" {
		return DnsRecordResponse{}, errors.New(dnsRecord.Message)
	}

	return dnsRecord, nil
}

func (c *Client) DeleteDnsRecord(ctx context.Context, domain, recordId string) error {
	uri := BASE_URL + "dns/" + domain + "/record/" + recordId + "/delete"

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var firewall DnsRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&firewall); err != nil {
		return err
	}
	if firewall.Status != "success" {
		return errors.New(firewall.Message)
	}

	return nil
}
