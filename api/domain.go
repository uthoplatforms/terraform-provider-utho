package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoplatforms/terraform-provider-utho/helper"

	"net/http"
)

type (
	DnsDomains struct {
		Domains []Domain `json:"domains"`
	}
	Domain struct {
		Domain         string   `json:"domain"`
		Status         string   `json:"status"`
		Message        string   `json:"message"`
		Nspoint        string   `json:"nspoint"`
		CreatedAt      string   `json:"created_at"`
		DnsrecordCount string   `json:"dnsrecord_count"`
		Records        []Record `json:"records"`
	}
	Record struct {
		ID       string `json:"id"`
		Hostname string `json:"hostname"`
		Type     string `json:"type"`
		Value    string `json:"value"`
		TTL      string `json:"ttl"`
		Priority string `json:"priority"`
	}
)
type DomainRequest struct {
	Domain string `json:"domain"`
}

func (c *Client) CreateDomain(ctx context.Context, serverRequest DomainRequest) (Domain, error) {
	uri := BASE_URL + "dns/adddomain"

	jsonPayload, err := json.Marshal(serverRequest)
	if err != nil {
		return Domain{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return Domain{}, err
	}

	defer resp.Body.Close()

	var domain Domain
	if err := json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return Domain{}, err
	}
	if domain.Status != "success" {
		return Domain{}, errors.New(domain.Message)
	}
	return domain, nil
}

func (c *Client) GetDomain(ctx context.Context, domainName string) (Domain, error) {
	uri := BASE_URL + "dns/" + domainName

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Domain{}, err
	}
	defer resp.Body.Close()

	var dns_domains DnsDomains
	if err := json.NewDecoder(resp.Body).Decode(&dns_domains); err != nil {
		return Domain{}, err
	}

	if len(dns_domains.Domains) == 0 {
		return Domain{}, errors.New("domain not found")
	}
	return dns_domains.Domains[0], nil
}

func (c *Client) DeleteDomain(ctx context.Context, domainName string) error {
	uri := BASE_URL + "dns/" + domainName + "/delete"

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var domain Domain
	if err := json.NewDecoder(resp.Body).Decode(&domain); err != nil {
		return err
	}
	if domain.Status != "success" {
		return errors.New(domain.Message)
	}
	return nil
}
