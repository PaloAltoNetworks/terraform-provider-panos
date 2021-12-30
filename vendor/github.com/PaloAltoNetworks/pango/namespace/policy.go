package namespace

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PaloAltoNetworks/pango/audit"
	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

/*
Policy is a namespace struct for config that is not imported into a vsys.

This struct contains additional operational state functions relevant for
policy rules.
*/
type Policy struct {
	Standard
}

// HitCount gets the rule hit count for the given rules.
//
// If the rules param is nil, then the hit count for all rules is returned.
func (n *Policy) HitCount(base, vsys string, rules []string) ([]util.HitCount, error) {
	if !n.Client.Versioning().Gte(version.Number{8, 1, 0, ""}) {
		return nil, fmt.Errorf("rule hit count requires PAN-OS 8.1+")
	}

	req := util.NewHitCountRequest(base, vsys, rules)
	var resp util.HitCountResponse
	if _, err := n.Client.Op(req, "", nil, &resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}

// SetAuditComment sets an audit comment for the given rule.
func (n *Policy) SetAuditComment(pather Pather, rule, comment string) error {
	if rule == "" {
		return fmt.Errorf("rule must be specified")
	}
	path, err := pather([]string{rule})
	if err != nil {
		return err
	}

	n.Client.LogOp("(op) set audit comment for %q: %s", rule, comment)

	req := audit.SetComment{
		Xpath:   util.AsXpath(path),
		Comment: comment,
	}

	_, err = n.Client.Op(req, "", nil, nil)
	return err
}

// CurrentAuditComment gets the uncommitted audit comment for the given rule.
func (n *Policy) CurrentAuditComment(pather Pather, rule string) (string, error) {
	if rule == "" {
		return "", fmt.Errorf("rule must be specified")
	}
	path, err := pather([]string{rule})
	if err != nil {
		return "", err
	}

	n.Client.LogOp("(op) getting current audit comment for %q", rule)

	req := audit.GetComment{
		Xpath: util.AsXpath(path),
	}
	var resp audit.UncommittedComment

	_, err = n.Client.Op(req, "", nil, &resp)
	if err != nil {
		return "", err
	}

	return resp.Comment, nil
}

// AuditCommentHistory retrieves a chunk of historical audit comment logs.
func (n *Policy) AuditCommentHistory(pather Pather, rule, direction string, nlogs, skip int) ([]audit.Comment, error) {
	if rule == "" {
		return nil, fmt.Errorf("rule must be specified")
	}
	path, err := pather([]string{rule})
	if err != nil {
		return nil, err
	} else if len(path) != 6 && len(path) != 9 {
		return nil, fmt.Errorf("Invalid path length %d != (6, 9)", len(path))
	}

	var vsysDg string
	switch len(path) {
	case 6:
		vsysDg = "shared"
	case 9:
		tokens := strings.Split(path[4], "'")
		if len(tokens) != 3 {
			return nil, fmt.Errorf("vsys/dg retrieval not possible: %s", path[4])
		}
		vsysDg = tokens[1]
	}
	base := path[len(path)-4]
	rType := path[len(path)-3]
	query := strings.Join([]string{
		"(subtype eq audit-comment)",
		fmt.Sprintf("(path contains '\\'%s\\'')", rule),   // Name.
		fmt.Sprintf("(path contains '%s')", rType),        // Rule type.
		fmt.Sprintf("(path contains %s)", base),           // Rulebase.
		fmt.Sprintf("(path contains '\\'%s\\'')", vsysDg), // Vsys or device group.
	}, " and ")

	n.Client.LogLog("(log) retrieving %s audit comment history: %s", rType, rule)

	extras := url.Values{}
	extras.Set("uniq", "yes")

	var job util.JobResponse
	if _, err := n.Client.Log("config", "", query, direction, nlogs, skip, extras, &job); err != nil {
		return nil, err
	}

	var ans audit.CommentHistory
	if _, err = n.Client.WaitForLogs(job.Id, 500*time.Millisecond, 0, &ans); err != nil {
		return nil, err
	}

	if len(ans.Comments) != 0 {
		if clock, err := n.Client.Clock(); err == nil {
			for i := range ans.Comments {
				ans.Comments[i].SetTime(clock)
			}
		}
	}

	return ans.Comments, nil
}
