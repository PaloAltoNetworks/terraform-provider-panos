package group

import (
	"encoding/xml"
)

/** Structs / functions for this namespace. **/

/*
type normalizer interface {
    Normalize() []map[string] string
}
*/

type pmContainer_v1 struct {
	Group cGroup `xml:"result>port-maps>entry"`
}

func (o *pmContainer_v1) Normalize() []map[string]string {
	if len(o.Group.Clusters) == 0 {
		return nil
	}

	lenAns := 0
	for i := range o.Group.Clusters {
		lenAns += len(o.Group.Clusters[i].Services)
	}

	ans := make([]map[string]string, 0, lenAns)
	for i := range o.Group.Clusters {
		cName := o.Group.Clusters[i].ClusterName
		for j := range o.Group.Clusters[i].Services {
			ans = append(ans, map[string]string{
				"cluster-name":     cName,
				"namespace":        o.Group.Clusters[i].Services[j].Namespace,
				"service":          o.Group.Clusters[i].Services[j].Service,
				"external-lb-port": o.Group.Clusters[i].Services[j].ExternalLbPort,
				"internal-lb-port": o.Group.Clusters[i].Services[j].InternalLbPort,
				"target-port":      o.Group.Clusters[i].Services[j].TargetPort,
				"target-protocol":  o.Group.Clusters[i].Services[j].TargetProtocol,
				"dup-np":           o.Group.Clusters[i].Services[j].DupNp,
			})
		}
	}

	return ans
}

type cGroup struct {
	GroupName string    `xml:"name,attr"`
	Clusters  []cluster `xml:"gke-clusters>entry"`
}

type cluster struct {
	ClusterName string    `xml:"name,attr"`
	Services    []service `xml:"service>entry"`
}

type service struct {
	XMLName        xml.Name `xml:"entry"`
	Namespace      string   `xml:"namespace"`
	Service        string   `xml:"name,attr"`
	ExternalLbPort string   `xml:"port"`
	InternalLbPort string   `xml:"ilb-port"`
	TargetPort     string   `xml:"target-port"`
	TargetProtocol string   `xml:"protocol"`
	DupNp          string   `xml:"dup-np"`
}

type pmReq_v1 struct {
	XMLName xml.Name `xml:"show"`
	Group   string   `xml:"plugins>gcp>port-map>gke>name"`
	Cmd     string   `xml:"plugins>gcp>port-map>gke>service"`
}
