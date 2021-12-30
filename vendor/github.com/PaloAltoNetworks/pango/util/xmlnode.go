package util

import (
	"encoding/xml"
	"strings"
)

const (
	entryPrefix = "entry[@name='"
	entrySuffix = "']"
)

// XmlNode is a generic XML node.
type XmlNode struct {
	XMLName    xml.Name
	Attributes []xml.Attr `xml:",any,attr"`
	Text       []byte     `xml:",innerxml"`
	Nodes      []XmlNode  `xml:",any"`
}

// FindXmlNodeInTree finds a given path in the specified XmlNode tree.
func FindXmlNodeInTree(path []string, elm *XmlNode) *XmlNode {
	if len(path) == 0 {
		return elm
	}

	if elm == nil {
		return elm
	}

	tag := path[0]
	path = path[1:]
	var name string
	if strings.HasPrefix(tag, entryPrefix) {
		name = strings.TrimSuffix(strings.TrimPrefix(tag, entryPrefix), entrySuffix)
		tag = "entry"
	}

	for _, x := range elm.Nodes {
		if x.XMLName.Local == tag {
			if name == "" {
				return FindXmlNodeInTree(path, &x)
			} else {
				for _, atr := range x.Attributes {
					if atr.Name.Local == "name" && atr.Value == name {
						return FindXmlNodeInTree(path, &x)
					}
				}
			}
		}
	}

	return nil
}
