package namespace

import (
	"bytes"
	"encoding/xml"
)

// UnpackageXmlInto wraps XML content into a throw-away wrapper for further
// unmarshaling.  This is basically for retrieving sub-object config from
// a parent's raw XML field mapping.
func UnpackageXmlInto(b []byte, res interface{}) error {
	var buf bytes.Buffer
	buf.Grow(len(b) + 7)
	buf.Write([]byte("<a>"))
	buf.Write(b)
	buf.Write([]byte("</a>"))

	data := buf.Bytes()
	return xml.Unmarshal(data, res)
}
