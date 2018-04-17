package pango

/*
Connect opens a connection to the PAN-OS client, then uses the "model" info
to return a pointer to either a Firewall or Panorama struct.

The Initialize function is invoked as part of this discovery, so there is no
need to Initialize() the Client connection prior to invoking this.
*/
func Connect(c Client) (interface{}, error) {
    var err error

    logg := c.Logging
    c.Logging = LogQuiet

    if err = c.Initialize(); err != nil {
        return nil, err
    }

    model := c.SystemInfo["model"]
    if model == "Panorama" || model[:2] == "M-" {
        pano := &Panorama{Client: c}
        pano.Logging = logg
        if err = pano.Initialize(); err != nil {
            return nil, err
        }
        return pano, nil
    } else {
        fw := &Firewall{Client: c}
        fw.Logging = logg
        if err = fw.Initialize(); err != nil {
            return nil, err
        }
        return fw, nil
    }
}
