package plugin

import (
	"encoding/xml"
)

type GetPlugins struct {
	XMLName xml.Name `xml:"show"`
	Cmd     string   `xml:"plugins>packages"`
}

type PackageListing struct {
	Answer []PackageInfo `xml:"result>plugins>entry"`
}

func (o *PackageListing) Listing() []Info {
	ans := make([]Info, 0, len(o.Answer))
	for _, x := range o.Answer {
		ans = append(ans, Info{
			Name:           x.Name,
			Version:        x.Version,
			ReleaseDate:    x.ReleaseDate,
			ReleaseNoteUrl: x.RelNote.ReleaseNoteUrl,
			PackageFile:    x.PackageFile,
			Size:           x.Size,
			Platform:       x.Platform,
			Installed:      x.Installed,
			Downloaded:     x.Downloaded,
		})
	}

	return ans
}

type PackageInfo struct {
	Name        string      `xml:"name"`
	Version     string      `xml:"version"`
	ReleaseDate string      `xml:"release-date"`
	RelNote     ReleaseNote `xml:"release-note-url"`
	PackageFile string      `xml:"pkg-file"`
	Size        string      `xml:"size"`
	Platform    string      `xml:"platform"`
	Installed   string      `xml:"installed"`
	Downloaded  string      `xml:"downloaded"`
}

type ReleaseNote struct {
	ReleaseNoteUrl string `xml:",cdata"`
}

// Info is normalized information on plugin packages available to PAN-OS.
type Info struct {
	Name           string
	Version        string
	ReleaseDate    string
	ReleaseNoteUrl string
	PackageFile    string
	Size           string
	Platform       string
	Installed      string
	Downloaded     string
}
