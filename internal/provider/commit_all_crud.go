package provider

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"

	"github.com/PaloAltoNetworks/pango/xmlapi"
)

type commitAllReq struct {
	XMLName xml.Name `xml:"commit-all"`
}

func (o commitAllReq) Action() string {
	return "all"
}

func (o commitAllReq) Element() any {
	return o
}

func (o *CommitAllAction) InvokeCustom(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {

	cmd := &xmlapi.Commit{
		Command: &commitAllReq{},
		Target:  o.client.GetTarget(),
	}

	var commitResp xmlapi.JobResponse

	_, _, err := o.client.Communicate(ctx, cmd, false, &commitResp)
	if err != nil {
		resp.Diagnostics.AddError("Failed to schedule a commit-all (push to devices)", err.Error())
		return
	}

	err = o.client.WaitForJob(ctx, commitResp.Id, 2*time.Second, nil)
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for commit-all task to finish", err.Error())
		return
	}
}
