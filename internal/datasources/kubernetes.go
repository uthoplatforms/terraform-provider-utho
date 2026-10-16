package datasources

import (
	"context"
	"fmt"

	"github.com/uthoplatforms/terraform-provider-utho/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type KubeconfigDataSource struct{ client *client.Client }

type KubeconfigModel struct {
	ID        types.String `tfsdk:"id"`
	ClusterID types.String `tfsdk:"cluster_id"`
	RawConfig types.String `tfsdk:"raw_config"`
}

func NewKubeconfigDataSource() datasource.DataSource { return &KubeconfigDataSource{} }

func (d *KubeconfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_config"
}

func (d *KubeconfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch the kubeconfig for a Utho Kubernetes cluster.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"cluster_id": schema.StringAttribute{Required: true, Description: "Kubernetes cluster ID."},
			"raw_config": schema.StringAttribute{Computed: true, Sensitive: true, Description: "Raw kubeconfig YAML content."},
		},
	}
}

func (d *KubeconfigDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", "Expected *client.Client")
		return
	}
	d.client = c
}

func (d *KubeconfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state KubeconfigModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	kubeconfig, err := d.client.GetKubeconfig(state.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error fetching kubeconfig", fmt.Sprintf("%s", err))
		return
	}

	state.ID = types.StringValue(state.ClusterID.ValueString())
	state.RawConfig = types.StringValue(kubeconfig)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
