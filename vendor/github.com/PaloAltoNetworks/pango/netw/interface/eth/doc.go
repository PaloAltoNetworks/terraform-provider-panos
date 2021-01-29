/*
Package eth is the client.Network.EthernetInterface namespace.

The Set() and Edit() functions in these namespaces takes a vsys parameter.  These
functions do not force you to specify a vsys to import the interface into, however
it should be noted that interfaces must be imported into a vsys in order for PAN-OS
to be able to use that interface.

Interfaces with a Mode of "ha" or "aggregate-group" will not be imported, as
is proper for these types of interfaces.

Normalized object:  Entry
*/
package eth
