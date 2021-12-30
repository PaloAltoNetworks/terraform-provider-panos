/*
Package certificate is the client.Device.Certificate namespace.

For Panorama, there are two possibilities:  managing this object on Panorama
itself or inside of a Template.

To manage objects on Panorama, leave "tmpl" and "vsys" params empty.

To manage objects in a template, specify the template name and the vsys (if
unspecified, defaults to "shared").

Configuring things such as "Forward Trust Certificate", "Forward Untrust
Certificate", and "Trusted Root CA" is done from the Device.SslDecrypt
namespace.

Note: PAN-OS 7.1+

Normalized object:  Entry
*/
package certificate
