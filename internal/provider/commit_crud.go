package provider

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"

	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type commitReq struct {
	XMLName xml.Name `xml:"commit"`
	action  string
}

func (o commitReq) Action() string {
	return o.action
}

func (o commitReq) Element() any {
	return o
}

func (o *CommitAction) InvokeCustom(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {

	cmd := &xmlapi.Commit{
		Command: &commitReq{},
		Target:  o.client.GetTarget(),
	}

	var commitResp xmlapi.JobResponse

	_, _, err := o.client.Communicate(ctx, cmd, false, &commitResp)
	if err != nil {
		resp.Diagnostics.AddError("Failed to schedule a commit", err.Error())
		return
	}

	err = o.client.WaitForJob(ctx, commitResp.Id, 2*time.Second, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for commit task to finish", err.Error())
		return
	}
}
