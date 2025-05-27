package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

var InspectionModeError = "Resources are unavailable when the provider is in inspection mode.  Resources are only available in API mode."

type DiagnosticsError struct {
	message     string
	diagnostics diag.Diagnostics
}

func NewDiagnosticsError(message string, diags diag.Diagnostics) *DiagnosticsError {
	return &DiagnosticsError{
		diagnostics: diags.Errors(),
	}
}

func (o *DiagnosticsError) Diagnostics() diag.Diagnostics {
	return o.diagnostics
}

func (o *DiagnosticsError) Error() string {
	var summaries []string
	for _, elt := range o.diagnostics {
		summaries = append(summaries, elt.Summary())
	}
	return fmt.Sprintf("%s: %s", o.message, strings.Join(summaries, ", "))
}
