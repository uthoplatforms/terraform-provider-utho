package api

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/uthoplatforms/terraform-provider-utho/helper"

	"net/http"
)

type TargetGroups struct {
	Status       string        `json:"status"`
	Message      string        `json:"message"`
	Targetgroups []TargetGroup `json:"targetgroups"`
}
type TargetGroup struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Port                string   `json:"port"`
	Protocol            string   `json:"protocol"`
	HealthCheckPath     string   `json:"health_check_path"`
	HealthCheckInterval string   `json:"health_check_interval"`
	HealthCheckProtocol string   `json:"health_check_protocol"`
	HealthCheckTimeout  string   `json:"health_check_timeout"`
	HealthyThreshold    string   `json:"healthy_threshold"`
	UnhealthyThreshold  string   `json:"unhealthy_threshold"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Targets             []Target `json:"targets"`
}
type Target struct {
	Lbid                string `json:"lbid"`
	IP                  string `json:"ip"`
	Cloudid             string `json:"cloudid"`
	Status              string `json:"status"`
	ScalingGroupid      string `json:"scaling_groupid"`
	KubernetesClusterid string `json:"kubernetes_clusterid"`
	BackendPort         string `json:"backend_port"`
	BackendProtocol     string `json:"backend_protocol"`
	TargetgroupID       string `json:"targetgroup_id"`
	FrontendID          string `json:"frontend_id"`
	ID                  string `json:"id"`
}

type CreateTargetGroupArgs struct {
	Name                string `json:"name"`
	Protocol            string `json:"protocol"`
	Port                string `json:"port"`
	HealthCheckPath     string `json:"health_check_path"`
	HealthCheckProtocol string `json:"health_check_protocol"`
	HealthCheckInterval string `json:"health_check_interval"`
	HealthCheckTimeout  string `json:"health_check_timeout"`
	HealthyThreshold    string `json:"healthy_threshold"`
	UnhealthyThreshold  string `json:"unhealthy_threshold"`
	Targets             []TargetArgs
}
type TargetArgs struct {
	BackendProtocol string `json:"backend_protocol"`
	BackendPort     string `json:"backend_port"`
	IP              string `json:"ip"`
}

type UpdateTargetGroupArgs struct {
	Name                string `json:"name"`
	Protocol            string `json:"protocol"`
	Port                string `json:"port"`
	HealthCheckPath     string `json:"health_check_path"`
	HealthCheckProtocol string `json:"health_check_protocol"`
	HealthCheckInterval string `json:"health_check_interval"`
	HealthCheckTimeout  string `json:"health_check_timeout"`
	HealthyThreshold    string `json:"healthy_threshold"`
	UnhealthyThreshold  string `json:"unhealthy_threshold"`
	Targets             []TargetArgs
}

type TargetGroupResponse struct {
	ID      int    `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (c *Client) CreateTargetGroup(ctx context.Context, createTargetGroupArgs CreateTargetGroupArgs) (TargetGroupResponse, error) {
	uri := BASE_URL + "targetgroup"

	jsonPayload, err := json.Marshal(createTargetGroupArgs)
	if err != nil {
		return TargetGroupResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return TargetGroupResponse{}, err
	}

	defer resp.Body.Close()

	var targetGroup TargetGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&targetGroup); err != nil {
		return TargetGroupResponse{}, err
	}
	if targetGroup.Status != "success" {
		return TargetGroupResponse{}, errors.New(targetGroup.Message)
	}

	for _, target := range createTargetGroupArgs.Targets {
		err = c.CreateTarget(ctx, strconv.Itoa(targetGroup.ID), target)
		if err != nil {
			return TargetGroupResponse{}, err
		}
	}

	return targetGroup, nil
}

func (c *Client) GetTargetGroup(ctx context.Context, id string) (TargetGroup, error) {
	uri := BASE_URL + "targetgroup"

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return TargetGroup{}, err
	}
	defer resp.Body.Close()

	var targetGroups TargetGroups
	if err := json.NewDecoder(resp.Body).Decode(&targetGroups); err != nil {
		return TargetGroup{}, err
	}

	if targetGroups.Status != "success" {
		return TargetGroup{}, errors.New(targetGroups.Message)
	}
	if len(targetGroups.Targetgroups) == 0 {
		return TargetGroup{}, errors.New("Target Group not found")
	}

	targetGroup := TargetGroup{}
	for _, t := range targetGroups.Targetgroups {
		if t.ID == id {
			targetGroup = t
		}
	}
	if len(targetGroup.ID) == 0 {
		return TargetGroup{}, errors.New("Target Group not found")
	}

	return targetGroup, nil
}

func (c *Client) UpdateTargetGroup(ctx context.Context, targetGroupId string, updateTargetGroupArgs UpdateTargetGroupArgs) (TargetGroupResponse, error) {
	uri := BASE_URL + "targetgroup/" + targetGroupId

	jsonPayload, err := json.Marshal(updateTargetGroupArgs)
	if err != nil {
		return TargetGroupResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPut, uri, jsonPayload, c.token)
	if err != nil {
		return TargetGroupResponse{}, err
	}

	defer resp.Body.Close()

	var targetGroup TargetGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&targetGroup); err != nil {
		return TargetGroupResponse{}, err
	}
	if targetGroup.Status != "success" {
		return TargetGroupResponse{}, errors.New(targetGroup.Message)
	}

	return targetGroup, nil
}

func (c *Client) DeleteTargetGroup(ctx context.Context, targetGroupId, targetGroupName string) error {
	uri := BASE_URL + "targetgroup/" + targetGroupId + "?name=" + targetGroupName

	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, nil, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var targetGroup TargetGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&targetGroup); err != nil {
		return err
	}
	if targetGroup.Status != "success" {
		return errors.New(targetGroup.Message)
	}

	return nil
}

func (c *Client) CreateTarget(ctx context.Context, targetGroupId string, target TargetArgs) error {
	uri := BASE_URL + "targetgroup/" + targetGroupId + "/target"

	jsonPayload, err := json.Marshal(target)
	if err != nil {
		return err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var targetGroup TargetGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&targetGroup); err != nil {
		return err
	}
	if targetGroup.Status != "success" {
		return errors.New(targetGroup.Message)
	}

	return nil
}
