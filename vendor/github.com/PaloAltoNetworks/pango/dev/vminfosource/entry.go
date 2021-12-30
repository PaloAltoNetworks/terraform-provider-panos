package vminfosource

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// VM information source.
type Entry struct {
	Name          string
	AwsVpc        *AwsVpc
	Esxi          *Esxi
	Vcenter       *Vcenter
	GoogleCompute *GoogleCompute // PAN-OS 8.1
}

type AwsVpc struct {
	Description     string
	Disabled        bool
	Source          string
	AccessKeyId     string
	SecretAccessKey string // encrypted
	UpdateInterval  int
	EnableTimeout   bool
	Timeout         int
	VpcId           string
}

type Esxi struct {
	Description    string
	Port           int
	Disabled       bool
	EnableTimeout  bool
	Timeout        int
	Source         string
	Username       string
	Password       string // encrypted
	UpdateInterval int
}

type Vcenter struct {
	Description    string
	Port           int
	Disabled       bool
	EnableTimeout  bool
	Timeout        int
	Source         string
	Username       string
	Password       string // encrypted
	UpdateInterval int
}

type GoogleCompute struct {
	Description              string
	Disabled                 bool
	AuthType                 string
	ServiceAccountCredential string // encrypted
	ProjectId                string
	ZoneName                 string
	UpdateInterval           int
	EnableTimeout            bool
	Timeout                  int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	if s.AwsVpc == nil {
		o.AwsVpc = nil
	} else {
		s.AwsVpc = &AwsVpc{
			Description:     s.AwsVpc.Description,
			Disabled:        s.AwsVpc.Disabled,
			Source:          s.AwsVpc.Source,
			AccessKeyId:     s.AwsVpc.AccessKeyId,
			SecretAccessKey: s.AwsVpc.SecretAccessKey,
			UpdateInterval:  s.AwsVpc.UpdateInterval,
			EnableTimeout:   s.AwsVpc.EnableTimeout,
			Timeout:         s.AwsVpc.Timeout,
			VpcId:           s.AwsVpc.VpcId,
		}
	}
	if s.Esxi == nil {
		o.Esxi = nil
	} else {
		s.Esxi = &Esxi{
			Description:    s.Esxi.Description,
			Port:           s.Esxi.Port,
			Disabled:       s.Esxi.Disabled,
			EnableTimeout:  s.Esxi.EnableTimeout,
			Timeout:        s.Esxi.Timeout,
			Source:         s.Esxi.Source,
			Username:       s.Esxi.Username,
			Password:       s.Esxi.Password,
			UpdateInterval: s.Esxi.UpdateInterval,
		}
	}
	if s.Vcenter == nil {
		o.Vcenter = nil
	} else {
		s.Vcenter = &Vcenter{
			Description:    s.Vcenter.Description,
			Port:           s.Vcenter.Port,
			Disabled:       s.Vcenter.Disabled,
			EnableTimeout:  s.Vcenter.EnableTimeout,
			Timeout:        s.Vcenter.Timeout,
			Source:         s.Vcenter.Source,
			Username:       s.Vcenter.Username,
			Password:       s.Vcenter.Password,
			UpdateInterval: s.Vcenter.UpdateInterval,
		}
	}
	if s.GoogleCompute == nil {
		o.GoogleCompute = nil
	} else {
		o.GoogleCompute = &GoogleCompute{
			Description:              s.GoogleCompute.Description,
			Disabled:                 s.GoogleCompute.Disabled,
			AuthType:                 s.GoogleCompute.AuthType,
			ServiceAccountCredential: s.GoogleCompute.ServiceAccountCredential,
			ProjectId:                s.GoogleCompute.ProjectId,
			ZoneName:                 s.GoogleCompute.ZoneName,
			UpdateInterval:           s.GoogleCompute.UpdateInterval,
			EnableTimeout:            s.GoogleCompute.EnableTimeout,
			Timeout:                  s.GoogleCompute.Timeout,
		}
	}
}

/** Structs / functions for normalization. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	AwsVpc  *vpc     `xml:"AWS-VPC"`
	Esxi    *esxi    `xml:"VMware-ESXi"`
	Vcenter *vcenter `xml:"VMware-vCenter"`
}

type vpc struct {
	Description     string `xml:"description,omitempty"`
	Disabled        string `xml:"disabled"`
	Source          string `xml:"source"`
	AccessKeyId     string `xml:"access-key-id"`
	SecretAccessKey string `xml:"secret-access-key"`
	UpdateInterval  int    `xml:"update-interval,omitempty"`
	EnableTimeout   string `xml:"vm-info-timeout-enable"`
	Timeout         int    `xml:"vm-info-timeout,omitempty"`
	VpcId           string `xml:"vpc-id"`
}

func (e *vpc) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local vpc
	ans := local{
		UpdateInterval: 60,
		Timeout:        2,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = vpc(ans)
	return nil
}

type esxi struct {
	Description    string `xml:"description,omitempty"`
	Port           int    `xml:"port,omitempty"`
	Disabled       string `xml:"disabled"`
	EnableTimeout  string `xml:"vm-info-timeout-enable"`
	Timeout        int    `xml:"vm-info-timeout,omitempty"`
	Source         string `xml:"source"`
	Username       string `xml:"username"`
	Password       string `xml:"password"`
	UpdateInterval int    `xml:"update-interval,omitempty"`
}

func (e *esxi) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local esxi
	ans := local{
		Port:           443,
		Timeout:        2,
		UpdateInterval: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = esxi(ans)
	return nil
}

type vcenter struct {
	Description    string `xml:"description,omitempty"`
	Port           int    `xml:"port,omitempty"`
	Disabled       string `xml:"disabled"`
	EnableTimeout  string `xml:"vm-info-timeout-enable"`
	Timeout        int    `xml:"vm-info-timeout,omitempty"`
	Source         string `xml:"source"`
	Username       string `xml:"username"`
	Password       string `xml:"password"`
	UpdateInterval int    `xml:"update-interval,omitempty"`
}

func (e *vcenter) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local vcenter
	ans := local{
		Port:           443,
		Timeout:        2,
		UpdateInterval: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = vcenter(ans)
	return nil
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if e.AwsVpc != nil {
		ans.AwsVpc = &vpc{
			Description:     e.AwsVpc.Description,
			Disabled:        util.YesNo(e.AwsVpc.Disabled),
			Source:          e.AwsVpc.Source,
			AccessKeyId:     e.AwsVpc.AccessKeyId,
			SecretAccessKey: e.AwsVpc.SecretAccessKey,
			UpdateInterval:  e.AwsVpc.UpdateInterval,
			EnableTimeout:   util.YesNo(e.AwsVpc.EnableTimeout),
			Timeout:         e.AwsVpc.Timeout,
			VpcId:           e.AwsVpc.VpcId,
		}
	}

	if e.Esxi != nil {
		ans.Esxi = &esxi{
			Description:    e.Esxi.Description,
			Port:           e.Esxi.Port,
			Disabled:       util.YesNo(e.Esxi.Disabled),
			EnableTimeout:  util.YesNo(e.Esxi.EnableTimeout),
			Timeout:        e.Esxi.Timeout,
			Source:         e.Esxi.Source,
			Username:       e.Esxi.Username,
			Password:       e.Esxi.Password,
			UpdateInterval: e.Esxi.UpdateInterval,
		}
	}

	if e.Vcenter != nil {
		ans.Vcenter = &vcenter{
			Description:    e.Vcenter.Description,
			Port:           e.Vcenter.Port,
			Disabled:       util.YesNo(e.Vcenter.Disabled),
			EnableTimeout:  util.YesNo(e.Vcenter.EnableTimeout),
			Timeout:        e.Vcenter.Timeout,
			Source:         e.Vcenter.Source,
			Username:       e.Vcenter.Username,
			Password:       e.Vcenter.Password,
			UpdateInterval: e.Vcenter.UpdateInterval,
		}
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name: e.Name,
	}

	if e.AwsVpc != nil {
		ans.AwsVpc = &AwsVpc{
			Description:     e.AwsVpc.Description,
			Disabled:        util.AsBool(e.AwsVpc.Disabled),
			Source:          e.AwsVpc.Source,
			AccessKeyId:     e.AwsVpc.AccessKeyId,
			SecretAccessKey: e.AwsVpc.SecretAccessKey,
			UpdateInterval:  e.AwsVpc.UpdateInterval,
			EnableTimeout:   util.AsBool(e.AwsVpc.EnableTimeout),
			Timeout:         e.AwsVpc.Timeout,
			VpcId:           e.AwsVpc.VpcId,
		}
	}

	if e.Esxi != nil {
		ans.Esxi = &Esxi{
			Description:    e.Esxi.Description,
			Port:           e.Esxi.Port,
			Disabled:       util.AsBool(e.Esxi.Disabled),
			EnableTimeout:  util.AsBool(e.Esxi.EnableTimeout),
			Timeout:        e.Esxi.Timeout,
			Source:         e.Esxi.Source,
			Username:       e.Esxi.Username,
			Password:       e.Esxi.Password,
			UpdateInterval: e.Esxi.UpdateInterval,
		}
	}

	if e.Vcenter != nil {
		ans.Vcenter = &Vcenter{
			Description:    e.Vcenter.Description,
			Port:           e.Vcenter.Port,
			Disabled:       util.AsBool(e.Vcenter.Disabled),
			EnableTimeout:  util.AsBool(e.Vcenter.EnableTimeout),
			Timeout:        e.Vcenter.Timeout,
			Source:         e.Vcenter.Source,
			Username:       e.Vcenter.Username,
			Password:       e.Vcenter.Password,
			UpdateInterval: e.Vcenter.UpdateInterval,
		}
	}

	return ans
}

// PAN-OS 8.1
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v2 struct {
	XMLName       xml.Name       `xml:"entry"`
	Name          string         `xml:"name,attr"`
	AwsVpc        *vpc           `xml:"AWS-VPC"`
	Esxi          *esxi          `xml:"VMware-ESXi"`
	Vcenter       *vcenter       `xml:"VMware-vCenter"`
	GoogleCompute *googleCompute `xml:"Google-Compute-Engine"`
}

type googleCompute struct {
	Description    string `xml:"description,omitempty"`
	Disabled       string `xml:"disabled"`
	Auth           gcAuth `xml:"service-auth-type"`
	ProjectId      string `xml:"project-id"`
	ZoneName       string `xml:"zone-name"`
	UpdateInterval int    `xml:"update-interval,omitempty"`
	EnableTimeout  string `xml:"vm-info-timeout-enable"`
	Timeout        int    `xml:"vm-info-timeout,omitempty"`
}

type gcAuth struct {
	ServiceInGce   *string               `xml:"service-in-gce"`
	ServiceAccount *gcAuthServiceAccount `xml:"service-account"`
}

type gcAuthServiceAccount struct {
	ServiceAccountCredential string `xml:"service-account-cred"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	if e.AwsVpc != nil {
		ans.AwsVpc = &vpc{
			Description:     e.AwsVpc.Description,
			Disabled:        util.YesNo(e.AwsVpc.Disabled),
			Source:          e.AwsVpc.Source,
			AccessKeyId:     e.AwsVpc.AccessKeyId,
			SecretAccessKey: e.AwsVpc.SecretAccessKey,
			UpdateInterval:  e.AwsVpc.UpdateInterval,
			EnableTimeout:   util.YesNo(e.AwsVpc.EnableTimeout),
			Timeout:         e.AwsVpc.Timeout,
			VpcId:           e.AwsVpc.VpcId,
		}
	}

	if e.Esxi != nil {
		ans.Esxi = &esxi{
			Description:    e.Esxi.Description,
			Port:           e.Esxi.Port,
			Disabled:       util.YesNo(e.Esxi.Disabled),
			EnableTimeout:  util.YesNo(e.Esxi.EnableTimeout),
			Timeout:        e.Esxi.Timeout,
			Source:         e.Esxi.Source,
			Username:       e.Esxi.Username,
			Password:       e.Esxi.Password,
			UpdateInterval: e.Esxi.UpdateInterval,
		}
	}

	if e.Vcenter != nil {
		ans.Vcenter = &vcenter{
			Description:    e.Vcenter.Description,
			Port:           e.Vcenter.Port,
			Disabled:       util.YesNo(e.Vcenter.Disabled),
			EnableTimeout:  util.YesNo(e.Vcenter.EnableTimeout),
			Timeout:        e.Vcenter.Timeout,
			Source:         e.Vcenter.Source,
			Username:       e.Vcenter.Username,
			Password:       e.Vcenter.Password,
			UpdateInterval: e.Vcenter.UpdateInterval,
		}
	}

	if e.GoogleCompute != nil {
		ans.GoogleCompute = &googleCompute{
			Description:    e.GoogleCompute.Description,
			Disabled:       util.YesNo(e.GoogleCompute.Disabled),
			ProjectId:      e.GoogleCompute.ProjectId,
			ZoneName:       e.GoogleCompute.ZoneName,
			UpdateInterval: e.GoogleCompute.UpdateInterval,
			EnableTimeout:  util.YesNo(e.GoogleCompute.EnableTimeout),
			Timeout:        e.GoogleCompute.Timeout,
		}

		switch e.GoogleCompute.AuthType {
		case AuthTypeServiceInGce:
			s := ""
			ans.GoogleCompute.Auth.ServiceInGce = &s
		case AuthTypeServiceAccount:
			ans.GoogleCompute.Auth.ServiceAccount = &gcAuthServiceAccount{
				ServiceAccountCredential: e.GoogleCompute.ServiceAccountCredential,
			}
		}
	}

	return ans
}

func (e *entry_v2) normalize() Entry {
	ans := Entry{
		Name: e.Name,
	}

	if e.AwsVpc != nil {
		ans.AwsVpc = &AwsVpc{
			Description:     e.AwsVpc.Description,
			Disabled:        util.AsBool(e.AwsVpc.Disabled),
			Source:          e.AwsVpc.Source,
			AccessKeyId:     e.AwsVpc.AccessKeyId,
			SecretAccessKey: e.AwsVpc.SecretAccessKey,
			UpdateInterval:  e.AwsVpc.UpdateInterval,
			EnableTimeout:   util.AsBool(e.AwsVpc.EnableTimeout),
			Timeout:         e.AwsVpc.Timeout,
			VpcId:           e.AwsVpc.VpcId,
		}
	}

	if e.Esxi != nil {
		ans.Esxi = &Esxi{
			Description:    e.Esxi.Description,
			Port:           e.Esxi.Port,
			Disabled:       util.AsBool(e.Esxi.Disabled),
			EnableTimeout:  util.AsBool(e.Esxi.EnableTimeout),
			Timeout:        e.Esxi.Timeout,
			Source:         e.Esxi.Source,
			Username:       e.Esxi.Username,
			Password:       e.Esxi.Password,
			UpdateInterval: e.Esxi.UpdateInterval,
		}
	}

	if e.Vcenter != nil {
		ans.Vcenter = &Vcenter{
			Description:    e.Vcenter.Description,
			Port:           e.Vcenter.Port,
			Disabled:       util.AsBool(e.Vcenter.Disabled),
			EnableTimeout:  util.AsBool(e.Vcenter.EnableTimeout),
			Timeout:        e.Vcenter.Timeout,
			Source:         e.Vcenter.Source,
			Username:       e.Vcenter.Username,
			Password:       e.Vcenter.Password,
			UpdateInterval: e.Vcenter.UpdateInterval,
		}
	}

	if e.GoogleCompute != nil {
		ans.GoogleCompute = &GoogleCompute{
			Description:    e.GoogleCompute.Description,
			Disabled:       util.AsBool(e.GoogleCompute.Disabled),
			ProjectId:      e.GoogleCompute.ProjectId,
			ZoneName:       e.GoogleCompute.ZoneName,
			UpdateInterval: e.GoogleCompute.UpdateInterval,
			EnableTimeout:  util.AsBool(e.GoogleCompute.EnableTimeout),
			Timeout:        e.GoogleCompute.Timeout,
		}

		switch {
		case e.GoogleCompute.Auth.ServiceInGce != nil:
			ans.GoogleCompute.AuthType = AuthTypeServiceInGce
		case e.GoogleCompute.Auth.ServiceAccount != nil:
			ans.GoogleCompute.AuthType = AuthTypeServiceAccount
			ans.GoogleCompute.ServiceAccountCredential = e.GoogleCompute.Auth.ServiceAccount.ServiceAccountCredential
		}
	}

	return ans
}
