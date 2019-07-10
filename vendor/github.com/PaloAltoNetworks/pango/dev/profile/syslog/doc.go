/*
Package syslog is the client.Object.SyslogServerProfile namespace.

For Panorama, there are two possibilities:  managing this object on Panorama
itself or inside of a Template.

To manage objects save on Panorama, leave "tmpl" and "ts" params empty and
set "dg" to "shared" (which is also the default).

To manage objects in a template, specify the template name and the vsys (if
unspecified, defaults to "shared").

Normalized object:  Entry
*/
package syslog
