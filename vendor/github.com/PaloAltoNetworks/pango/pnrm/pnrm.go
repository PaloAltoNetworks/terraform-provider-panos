/*
Package pnrm is the client.Panorama namespace.
*/
package pnrm


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/pnrm/dg"
    "github.com/PaloAltoNetworks/pango/pnrm/template"
    "github.com/PaloAltoNetworks/pango/pnrm/template/stack"
    "github.com/PaloAltoNetworks/pango/pnrm/template/variable"
)


// Pnrm is the panorama.DeviceGroup namespace.
type Pnrm struct {
    DeviceGroup *dg.Dg
    Template *template.Template
    TemplateStack *stack.Stack
    TemplateVariable *variable.Variable
}

// Initialize is invoked on panorama.Initialize().
func (c *Pnrm) Initialize(i util.XapiClient) {
    c.DeviceGroup = &dg.Dg{}
    c.DeviceGroup.Initialize(i)

    c.Template = &template.Template{}
    c.Template.Initialize(i)

    c.TemplateStack = &stack.Stack{}
    c.TemplateStack.Initialize(i)

    c.TemplateVariable = &variable.Variable{}
    c.TemplateVariable.Initialize(i)
}
