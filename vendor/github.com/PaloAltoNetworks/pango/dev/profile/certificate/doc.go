/*
Package certificate is the client.Device.CertificateProfile namespace.

For Panorama, there are three possibilities:
- local to Panorama
- in /config/shared
- inside a template

To manage certificates on Panorama, leave "tmpl" and "ts" params empty, then
either leave `dg` as an empty string (for certs in /config/panorama) or specying
`dg="shared"` (for certs in /config/shared).

To manage objects in a template, specify the template name and the vsys (if
unspecified, defaults to "shared").

Normalized object:  Entry
*/
package certificate
