package main

import (
    "github.com/PaloAltoNetworks/pango"
    "github.com/PaloAltoNetworks/pango/objs/addr"

    "github.com/hashicorp/terraform/helper/schema"
)


func resourceAddressObject() *schema.Resource {
    return &schema.Resource{
        Create: createAddressObject,
        Read: readAddressObject,
        Update: updateAddressObject,
        Delete: deleteAddressObject,

        Schema: map[string] *schema.Schema{
            "name": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
                Description: "The address object's name",
            },
            "vsys": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Default: "vsys1",
                Description: "The vsys to put this address object in",
            },
            "type": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
                Default: "ip-netmask",
                Description: "The type of address object (ip-netmask, ip-range, fqdn)",
            },
            "value": &schema.Schema{
                Type: schema.TypeString,
                Required: true,
            },
            "description": &schema.Schema{
                Type: schema.TypeString,
                Optional: true,
            },
            "tag": &schema.Schema{
                Type: schema.TypeList,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
                Description: "Administrative tags for the address object",
            },
        },
    }
}

func parseAddressObject(d *schema.ResourceData) (string, addr.Entry) {
    vsys := d.Get("vsys").(string)
    o := addr.Entry{
        Name: d.Get("name").(string),
        Value: d.Get("value").(string),
        Type: d.Get("type").(string),
        Description: d.Get("description").(string),
        Tag: asStringList(d, "tag"),
    }

    return vsys, o
}

func saveDataAddressObject(d *schema.ResourceData, o addr.Entry) {
    d.SetId(o.Name)
    d.Set("value", o.Value)
    d.Set("type", o.Type)
    d.Set("description", o.Description)
    d.Set("tag", o.Tag)
}

func createAddressObject(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys, o := parseAddressObject(d)

    list, err := fw.Objects.Address.GetList(vsys)
    if err != nil {
        return err
    }

    for i := range list {
        if list[i] == o.Name {
            d.SetId(o.Name)
            return nil
        }
    }

    if err = fw.Objects.Address.Set(vsys, o); err != nil {
        d.SetId(o.Name)
    }

    return nil
}

func readAddressObject(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    name := d.Get("name").(string)
    vsys := d.Get("vsys").(string)

    o, err := fw.Objects.Address.Get(vsys, name)
    if err != nil {
        d.SetId("")
        return nil
    }

    saveDataAddressObject(d, o)
    return nil
}

func updateAddressObject(d *schema.ResourceData, meta interface{}) error {
    var err error
    fw := meta.(*pango.Firewall)
    vsys, o := parseAddressObject(d)

    lo, err := fw.Objects.Address.Get(vsys, o.Name)
    if err == nil {
        lo.Copy(o)
        err = fw.Objects.Address.Edit(vsys, lo)
    } else {
        err = fw.Objects.Address.Set(vsys, o)
    }

    if err == nil {
        saveDataAddressObject(d, o)
    }
    return err
}

func deleteAddressObject(d *schema.ResourceData, meta interface{}) error {
    fw := meta.(*pango.Firewall)
    vsys := d.Get("vsys").(string)
    name := d.Get("name").(string)

    _ = fw.Objects.Address.Delete(vsys, name)
    d.SetId("")
    return nil
}
