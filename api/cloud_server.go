package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uthoterraform/terraform-provider-utho/helper"

	"net/http"
)

type CloudServer struct {
	Cloud   []Cloud `json:"cloud"`
	Status  string  `json:"status"`
	Message string  `json:"message"`
}

type Cloud struct {
	Cloudid           string                 `json:"cloudid"`
	Hostname          string                 `json:"hostname"`
	CPU               string                 `json:"cpu"`
	RAM               string                 `json:"ram"`
	ManagedOs         string                 `json:"managed_os"`
	ManagedFull       string                 `json:"managed_full"`
	ManagedOnetime    string                 `json:"managed_onetime"`
	PlanDisksize      int                    `json:"plan_disksize"`
	Disksize          int                    `json:"disksize"`
	Ha                string                 `json:"ha"`
	Status            string                 `json:"status"`
	Iso               string                 `json:"iso"`
	IP                string                 `json:"ip"`
	Billingcycle      string                 `json:"billingcycle"`
	Cost              float64                `json:"cost"`
	Vmcost            float64                `json:"vmcost"`
	Imagecost         int                    `json:"imagecost"`
	Backupcost        int                    `json:"backupcost"`
	Hourlycost        float64                `json:"hourlycost"`
	Cloudhourlycost   float64                `json:"cloudhourlycost"`
	Imagehourlycost   int                    `json:"imagehourlycost"`
	Backuphourlycost  int                    `json:"backuphourlycost"`
	Creditrequired    float64                `json:"creditrequired"`
	Creditreserved    int                    `json:"creditreserved"`
	Nextinvoiceamount float64                `json:"nextinvoiceamount"`
	Nextinvoicehours  string                 `json:"nextinvoicehours"`
	Consolepassword   string                 `json:"consolepassword"`
	Powerstatus       string                 `json:"powerstatus"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
	Nextduedate       string                 `json:"nextduedate"`
	Bandwidth         string                 `json:"bandwidth"`
	BandwidthUsed     int                    `json:"bandwidth_used"`
	BandwidthFree     int                    `json:"bandwidth_free"`
	Features          Features               `json:"features"`
	Image             CloudServerImage       `json:"image"`
	Dclocation        Dclocation             `json:"dclocation"`
	Networks          Networks               `json:"networks"`
	Storages          []Storages             `json:"storages"`
	Snapshots         []Snapshots            `json:"snapshots"`
	Firewalls         []CloudServerFirewalls `json:"firewalls"`
	GpuAvailable      string                 `json:"gpu_available"`
}
type CloudServerFirewalls struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
type Features struct {
	Backups string `json:"backups"`
}
type CloudServerImage struct {
	Name         string `json:"name"`
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
	Image        string `json:"image"`
	Cost         string `json:"cost"`
}
type Dclocation struct {
	Location string `json:"location"`
	Country  string `json:"country"`
	Dc       string `json:"dc"`
	Dccc     string `json:"dccc"`
}
type Public struct {
	V4 []V4 `json:"v4"`
}
type V4 struct {
	IPAddress string `json:"ip_address"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	Type      string `json:"type"`
	Nat       bool   `json:"nat"`
	Primary   string `json:"primary"`
}
type Private struct {
	V4 []PrivateV4 `json:"v4"`
}
type PrivateV4 struct {
	Noip      int    `json:"noip"`
	IPAddress string `json:"ip_address"`
	VpcName   string `json:"vpc_name"`
	Network   string `json:"network"`
	VpcID     string `json:"vpc_id"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	Type      string `json:"type"`
	Primary   string `json:"primary"`
}
type Networks struct {
	Public  Public  `json:"public"`
	Private Private `json:"private"`
}
type Storages struct {
	ID        string `json:"id"`
	Size      int    `json:"size"`
	DiskUsed  string `json:"disk_used"`
	DiskFree  string `json:"disk_free"`
	DiskUsedp string `json:"disk_usedp"`
	CreatedAt string `json:"created_at"`
	Bus       string `json:"bus"`
	Type      string `json:"type"`
}
type Snapshots struct {
	ID        string `json:"id"`
	Size      string `json:"size"`
	CreatedAt string `json:"created_at"`
	Note      string `json:"note"`
	Name      string `json:"name"`
}

type CreateCloudServerArgs struct {
	Dcslug       string          `json:"dcslug"`
	Image        string          `json:"image"`
	Planid       string          `json:"planid"`
	Auth         string          `json:"auth"`
	RootPassword string          `json:"root_password"`
	Firewall     string          `json:"firewall"`
	Enablebackup string          `json:"enablebackup"`
	Support      string          `json:"support"`
	Management   string          `json:"management"`
	Billingcycle string          `json:"billingcycle"`
	Backupid     string          `json:"backupid"`
	Snapshotid   string          `json:"snapshotid"`
	Sshkeys      string          `json:"sshkeys"`
	Cloud        []CloudHostname `json:"cloud"`
}
type CloudHostname struct {
	Hostname string `json:"hostname"`
}

type DeleteCloudServerArgs struct {
	Confirm string `json:"confirm"`
}

type CloudServerResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Cloudid  string `json:"cloudid"`
	Password string `json:"password"`
	Ipv4     string `json:"ipv4"`
}

func (c *Client) CreateCloudServer(ctx context.Context, createCloudServerArgs CreateCloudServerArgs) (CloudServerResponse, error) {
	uri := BASE_URL + "cloud/deploy"

	jsonPayload, err := json.Marshal(createCloudServerArgs)
	if err != nil {
		return CloudServerResponse{}, err
	}

	resp, err := helper.NewUthoRequest(ctx, http.MethodPost, uri, jsonPayload, c.token)
	if err != nil {
		return CloudServerResponse{}, err
	}

	defer resp.Body.Close()

	var cloudServer CloudServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudServer); err != nil {
		return CloudServerResponse{}, err
	}
	if cloudServer.Status != "success" {
		return CloudServerResponse{}, errors.New(cloudServer.Message)
	}
	return cloudServer, nil
}

func (c *Client) GetCloudServer(ctx context.Context, id string) (Cloud, error) {
	uri := BASE_URL + "cloud/" + id

	resp, err := helper.NewUthoRequest(ctx, http.MethodGet, uri, nil, c.token)
	if err != nil {
		return Cloud{}, err
	}
	defer resp.Body.Close()
	var cloudServer CloudServer
	if err := json.NewDecoder(resp.Body).Decode(&cloudServer); err != nil {
		return Cloud{}, err
	}

	if len(cloudServer.Cloud) == 0 {
		return Cloud{}, errors.New("Cloud server not found")
	}
	if cloudServer.Status == "error" {
		return Cloud{}, errors.New(cloudServer.Message)
	}

	return cloudServer.Cloud[0], nil
}

func (c *Client) DeleteCloudServer(ctx context.Context, id string) error {
	uri := BASE_URL + "cloud/" + id + "/destroy"

	deleteCloudServerArgs := DeleteCloudServerArgs{
		Confirm: "I am aware this action will delete data and server permanently",
	}
	jsonPayload, err := json.Marshal(deleteCloudServerArgs)
	if err != nil {
		return err
	}
	resp, err := helper.NewUthoRequest(ctx, http.MethodDelete, uri, jsonPayload, c.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cloudServer CloudServerResponse
	if err := json.NewDecoder(resp.Body).Decode(&cloudServer); err != nil {
		return err
	}
	if cloudServer.Status != "success" {
		return errors.New(cloudServer.Message)
	}

	return nil
}
