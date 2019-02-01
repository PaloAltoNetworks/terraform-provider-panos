package exp

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a BGP
// export rule.
type Entry struct {
    Name string
    Enable bool
    UsedBy []string
    MatchAsPathRegex string
    MatchCommunityRegex string
    MatchExtendedCommunityRegex string
    MatchMed string
    MatchRouteTable string // 8.0+
    MatchAddressPrefix map[string] bool
    MatchNextHop []string
    MatchFromPeer []string
    Action string
    LocalPreference string
    Med string
    NextHop string
    Origin string
    AsPathLimit int
    AsPathType string
    AsPathValue string
    CommunityType string
    CommunityValue string
    ExtendedCommunityType string
    ExtendedCommunityValue string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Enable = s.Enable
    o.UsedBy = s.UsedBy
    o.MatchAsPathRegex = s.MatchAsPathRegex
    o.MatchCommunityRegex = s.MatchCommunityRegex
    o.MatchExtendedCommunityRegex = s.MatchExtendedCommunityRegex
    s.MatchMed = o.MatchMed
    o.MatchRouteTable = s.MatchRouteTable
    o.MatchAddressPrefix = s.MatchAddressPrefix
    o.MatchNextHop = s.MatchNextHop
    o.MatchFromPeer = s.MatchFromPeer
    o.Action = s.Action
    o.LocalPreference = s.LocalPreference
    o.Med = s.Med
    o.NextHop = s.NextHop
    o.Origin = s.Origin
    o.AsPathLimit = s.AsPathLimit
    o.AsPathType = s.AsPathType
    o.AsPathValue = s.AsPathValue
    o.CommunityType = s.CommunityType
    o.CommunityValue = s.CommunityValue
    o.ExtendedCommunityType = s.ExtendedCommunityType
    o.ExtendedCommunityValue = s.ExtendedCommunityValue
}

/** Structs / functions for this namespace. **/

type normalizer interface {
    Normalize() Entry
}

type container_v1 struct {
    Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        UsedBy: util.MemToStr(o.Answer.UsedBy),
    }

    if o.Answer.Match != nil {
        ans.MatchMed = o.Answer.Match.MatchMed
        ans.MatchNextHop = util.MemToStr(o.Answer.Match.MatchNextHop)
        ans.MatchFromPeer = util.MemToStr(o.Answer.Match.MatchFromPeer)

        if o.Answer.Match.MatchAsPathRegex != nil {
            ans.MatchAsPathRegex = o.Answer.Match.MatchAsPathRegex.Regex
        }

        if o.Answer.Match.MatchCommunityRegex != nil {
            ans.MatchCommunityRegex = o.Answer.Match.MatchCommunityRegex.Regex
        }

        if o.Answer.Match.MatchExtendedCommunityRegex != nil {
            ans.MatchExtendedCommunityRegex = o.Answer.Match.MatchExtendedCommunityRegex.Regex
        }

        if o.Answer.Match.MatchAddressPrefix != nil {
            m := make(map[string] bool)
            for _, v := range o.Answer.Match.MatchAddressPrefix.Entry {
                m[v.Name] = util.AsBool(v.Exact)
            }
            ans.MatchAddressPrefix = m
        }
    }

    if o.Answer.Action != nil {
        if o.Answer.Action.Deny != nil {
            ans.Action = ActionDeny
        } else if o.Answer.Action.Allow != nil {
            ans.Action = ActionAllow

            if o.Answer.Action.Allow.Update != nil {
                ans.LocalPreference = o.Answer.Action.Allow.Update.LocalPreference
                ans.Med = o.Answer.Action.Allow.Update.Med
                ans.NextHop = o.Answer.Action.Allow.Update.NextHop
                ans.Origin = o.Answer.Action.Allow.Update.Origin
                ans.AsPathLimit = o.Answer.Action.Allow.Update.AsPathLimit

                if o.Answer.Action.Allow.Update.AsPath != nil {
                    if o.Answer.Action.Allow.Update.AsPath.None != nil {
                        ans.AsPathType = AsPathTypeNone
                    } else if o.Answer.Action.Allow.Update.AsPath.Remove != nil {
                        ans.AsPathType = AsPathTypeRemove
                    } else if o.Answer.Action.Allow.Update.AsPath.Prepend != "" {
                        ans.AsPathType = AsPathTypePrepend
                        ans.AsPathValue = o.Answer.Action.Allow.Update.AsPath.Prepend
                    } else if o.Answer.Action.Allow.Update.AsPath.RemoveAndPrepend != "" {
                        ans.AsPathType = AsPathTypeRemoveAndPrepend
                        ans.AsPathValue = o.Answer.Action.Allow.Update.AsPath.RemoveAndPrepend
                    }
                }

                if o.Answer.Action.Allow.Update.Community != nil {
                    if o.Answer.Action.Allow.Update.Community.None != nil {
                        ans.CommunityType = CommunityTypeNone
                    } else if o.Answer.Action.Allow.Update.Community.RemoveAll != nil {
                        ans.CommunityType = CommunityTypeRemoveAll
                    } else if o.Answer.Action.Allow.Update.Community.RemoveRegex != "" {
                        ans.CommunityType = CommunityTypeRemoveRegex
                        ans.CommunityValue = o.Answer.Action.Allow.Update.Community.RemoveRegex
                    } else if o.Answer.Action.Allow.Update.Community.Append != nil {
                        ans.CommunityType = CommunityTypeAppend
                        ans.CommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.Community.Append)
                    } else if o.Answer.Action.Allow.Update.Community.Overwrite != nil {
                        ans.CommunityType = CommunityTypeOverwrite
                        ans.CommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.Community.Overwrite)
                    }
                }

                if o.Answer.Action.Allow.Update.ExtendedCommunity != nil {
                    if o.Answer.Action.Allow.Update.ExtendedCommunity.None != nil {
                        ans.ExtendedCommunityType = CommunityTypeNone
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveAll != nil {
                        ans.ExtendedCommunityType = CommunityTypeRemoveAll
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveRegex != "" {
                        ans.ExtendedCommunityType = CommunityTypeRemoveRegex
                        ans.ExtendedCommunityValue = o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveRegex
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.Append != nil {
                        ans.ExtendedCommunityType = CommunityTypeAppend
                        ans.ExtendedCommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.ExtendedCommunity.Append)
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.Overwrite != nil {
                        ans.ExtendedCommunityType = CommunityTypeOverwrite
                        ans.ExtendedCommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.ExtendedCommunity.Overwrite)
                    }
                }
            }
        }
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        UsedBy: util.MemToStr(o.Answer.UsedBy),
    }

    if o.Answer.Match != nil {
        ans.MatchMed = o.Answer.Match.MatchMed
        ans.MatchRouteTable = o.Answer.Match.MatchRouteTable
        ans.MatchNextHop = util.MemToStr(o.Answer.Match.MatchNextHop)
        ans.MatchFromPeer = util.MemToStr(o.Answer.Match.MatchFromPeer)

        if o.Answer.Match.MatchAsPathRegex != nil {
            ans.MatchAsPathRegex = o.Answer.Match.MatchAsPathRegex.Regex
        }

        if o.Answer.Match.MatchCommunityRegex != nil {
            ans.MatchCommunityRegex = o.Answer.Match.MatchCommunityRegex.Regex
        }

        if o.Answer.Match.MatchExtendedCommunityRegex != nil {
            ans.MatchExtendedCommunityRegex = o.Answer.Match.MatchExtendedCommunityRegex.Regex
        }

        if o.Answer.Match.MatchAddressPrefix != nil {
            m := make(map[string] bool)
            for _, v := range o.Answer.Match.MatchAddressPrefix.Entry {
                m[v.Name] = util.AsBool(v.Exact)
            }
            ans.MatchAddressPrefix = m
        }
    }

    if o.Answer.Action != nil {
        if o.Answer.Action.Deny != nil {
            ans.Action = ActionDeny
        } else if o.Answer.Action.Allow != nil {
            ans.Action = ActionAllow

            if o.Answer.Action.Allow.Update != nil {
                ans.LocalPreference = o.Answer.Action.Allow.Update.LocalPreference
                ans.Med = o.Answer.Action.Allow.Update.Med
                ans.NextHop = o.Answer.Action.Allow.Update.NextHop
                ans.Origin = o.Answer.Action.Allow.Update.Origin
                ans.AsPathLimit = o.Answer.Action.Allow.Update.AsPathLimit

                if o.Answer.Action.Allow.Update.AsPath != nil {
                    if o.Answer.Action.Allow.Update.AsPath.None != nil {
                        ans.AsPathType = AsPathTypeNone
                    } else if o.Answer.Action.Allow.Update.AsPath.Remove != nil {
                        ans.AsPathType = AsPathTypeRemove
                    } else if o.Answer.Action.Allow.Update.AsPath.Prepend != "" {
                        ans.AsPathType = AsPathTypePrepend
                        ans.AsPathValue = o.Answer.Action.Allow.Update.AsPath.Prepend
                    } else if o.Answer.Action.Allow.Update.AsPath.RemoveAndPrepend != "" {
                        ans.AsPathType = AsPathTypeRemoveAndPrepend
                        ans.AsPathValue = o.Answer.Action.Allow.Update.AsPath.RemoveAndPrepend
                    }
                }

                if o.Answer.Action.Allow.Update.Community != nil {
                    if o.Answer.Action.Allow.Update.Community.None != nil {
                        ans.CommunityType = CommunityTypeNone
                    } else if o.Answer.Action.Allow.Update.Community.RemoveAll != nil {
                        ans.CommunityType = CommunityTypeRemoveAll
                    } else if o.Answer.Action.Allow.Update.Community.RemoveRegex != "" {
                        ans.CommunityType = CommunityTypeRemoveRegex
                        ans.CommunityValue = o.Answer.Action.Allow.Update.Community.RemoveRegex
                    } else if o.Answer.Action.Allow.Update.Community.Append != nil {
                        ans.CommunityType = CommunityTypeAppend
                        ans.CommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.Community.Append)
                    } else if o.Answer.Action.Allow.Update.Community.Overwrite != nil {
                        ans.CommunityType = CommunityTypeOverwrite
                        ans.CommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.Community.Overwrite)
                    }
                }

                if o.Answer.Action.Allow.Update.ExtendedCommunity != nil {
                    if o.Answer.Action.Allow.Update.ExtendedCommunity.None != nil {
                        ans.ExtendedCommunityType = CommunityTypeNone
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveAll != nil {
                        ans.ExtendedCommunityType = CommunityTypeRemoveAll
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveRegex != "" {
                        ans.ExtendedCommunityType = CommunityTypeRemoveRegex
                        ans.ExtendedCommunityValue = o.Answer.Action.Allow.Update.ExtendedCommunity.RemoveRegex
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.Append != nil {
                        ans.ExtendedCommunityType = CommunityTypeAppend
                        ans.ExtendedCommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.ExtendedCommunity.Append)
                    } else if o.Answer.Action.Allow.Update.ExtendedCommunity.Overwrite != nil {
                        ans.ExtendedCommunityType = CommunityTypeOverwrite
                        ans.ExtendedCommunityValue = util.MemToOneStr(o.Answer.Action.Allow.Update.ExtendedCommunity.Overwrite)
                    }
                }
            }
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    UsedBy *util.MemberType `xml:"used-by"`
    Match *match_v1 `xml:"match"`
    Action *action `xml:"action"`
}

type match_v1 struct {
    MatchAsPathRegex *regex `xml:"as-path"`
    MatchCommunityRegex *regex `xml:"community"`
    MatchExtendedCommunityRegex *regex `xml:"extended-community"`
    MatchMed string `xml:"med,omitempty"`
    MatchAddressPrefix *addPre `xml:"address-prefix"`
    MatchNextHop *util.MemberType `xml:"nexthop"`
    MatchFromPeer *util.MemberType `xml:"from-peer"`
}

type addPre struct {
    Entry []apEntry `xml:"entry"`
}

type apEntry struct {
    Name string `xml:"name,attr"`
    Exact string `xml:"exact"`
}

type regex struct {
    Regex string `xml:"regex,omitempty"`
}

type action struct {
    Deny *string `xml:"deny"`
    Allow *allow `xml:"allow"`
}

type allow struct {
    Update *update `xml:"update"`
}

type update struct {
    LocalPreference string `xml:"local-preference,omitempty"`
    Med string `xml:"med,omitempty"`
    NextHop string `xml:"nexthop,omitempty"`
    Origin string `xml:"origin,omitempty"`
    AsPathLimit int `xml:"as-path-limit,omitempty"`
    AsPath *asPath `xml:"as-path"`
    Community *allowCom `xml:"community"`
    ExtendedCommunity *allowCom `xml:"extended-community"`
}

type asPath struct {
    None *string `xml:"none"`
    Remove *string `xml:"remove"`
    Prepend string `xml:"prepend,omitempty"`
    RemoveAndPrepend string `xml:"remove-and-prepend,omitempty"`
}

type allowCom struct {
    None *string `xml:"none"`
    RemoveAll *string `xml:"remove-all"`
    RemoveRegex string `xml:"remove-regex,omitempty"`
    Append *util.MemberType `xml:"append"`
    Overwrite *util.MemberType `xml:"overwrite"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        UsedBy: util.StrToMem(e.UsedBy),
    }
    s := ""

    if e.MatchAsPathRegex != "" || e.MatchCommunityRegex != "" || e.MatchExtendedCommunityRegex != "" || e.MatchMed != "" || len(e.MatchAddressPrefix) > 0 || len(e.MatchNextHop) > 0 || len(e.MatchFromPeer) > 0 {
        ans.Match = &match_v1{
            MatchMed: e.MatchMed,
            MatchNextHop: util.StrToMem(e.MatchNextHop),
            MatchFromPeer: util.StrToMem(e.MatchFromPeer),
        }

        if e.MatchAsPathRegex != "" {
            ans.Match.MatchAsPathRegex = &regex{
                Regex: e.MatchAsPathRegex,
            }
        }

        if e.MatchCommunityRegex != "" {
            ans.Match.MatchCommunityRegex = &regex{
                Regex: e.MatchCommunityRegex,
            }
        }

        if e.MatchExtendedCommunityRegex != "" {
            ans.Match.MatchExtendedCommunityRegex = &regex{
                Regex: e.MatchExtendedCommunityRegex,
            }
        }

        if len(e.MatchAddressPrefix) > 0 {
            apList := make([]apEntry, 0, len(e.MatchAddressPrefix))
            for k, v := range e.MatchAddressPrefix {
                apList = append(apList, apEntry{
                    Name: k,
                    Exact: util.YesNo(v),
                })
            }
            ans.Match.MatchAddressPrefix = &addPre{apList}
        }
    }

    switch e.Action {
    case ActionDeny:
        ans.Action = &action{
            Deny: &s,
        }
    case ActionAllow:
        ans.Action = &action{
            Allow: &allow{
            },
        }

        if e.LocalPreference != "" || e.Med != "" || e.NextHop != "" || e.Origin != "" || e.AsPathLimit != 0 || e.AsPathType != "" || e.CommunityType != "" || e.ExtendedCommunityType != "" {
            u := update{
                LocalPreference: e.LocalPreference,
                Med: e.Med,
                NextHop: e.NextHop,
                Origin: e.Origin,
                AsPathLimit: e.AsPathLimit,
            }

            switch e.AsPathType {
            case AsPathTypeNone:
                u.AsPath = &asPath{
                    None: &s,
                }
            case AsPathTypeRemove:
                u.AsPath = &asPath{
                    Remove: &s,
                }
            case AsPathTypePrepend:
                u.AsPath = &asPath{
                    Prepend: e.AsPathValue,
                }
            case AsPathTypeRemoveAndPrepend:
                u.AsPath = &asPath{
                    RemoveAndPrepend: e.AsPathValue,
                }
            }

            switch e.CommunityType {
            case CommunityTypeNone:
                u.Community = &allowCom{
                    None: &s,
                }
            case CommunityTypeRemoveAll:
                u.Community = &allowCom{
                    RemoveAll: &s,
                }
            case CommunityTypeRemoveRegex:
                u.Community = &allowCom{
                    RemoveRegex: e.CommunityValue,
                }
            case CommunityTypeAppend:
                u.Community = &allowCom{
                    Append: util.OneStrToMem(e.CommunityValue),
                }
            case CommunityTypeOverwrite:
                u.Community = &allowCom{
                    Overwrite: util.OneStrToMem(e.CommunityValue),
                }
            }

            switch e.ExtendedCommunityType {
            case CommunityTypeNone:
                u.ExtendedCommunity = &allowCom{
                    None: &s,
                }
            case CommunityTypeRemoveAll:
                u.ExtendedCommunity = &allowCom{
                    RemoveAll: &s,
                }
            case CommunityTypeRemoveRegex:
                u.ExtendedCommunity = &allowCom{
                    RemoveRegex: e.ExtendedCommunityValue,
                }
            case CommunityTypeAppend:
                u.ExtendedCommunity = &allowCom{
                    Append: util.OneStrToMem(e.ExtendedCommunityValue),
                }
            case CommunityTypeOverwrite:
                u.ExtendedCommunity = &allowCom{
                    Overwrite: util.OneStrToMem(e.ExtendedCommunityValue),
                }
            }

            ans.Action.Allow.Update = &u
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    UsedBy *util.MemberType `xml:"used-by"`
    Match *match_v2 `xml:"match"`
    Action *action `xml:"action"`
}

type match_v2 struct {
    MatchAsPathRegex *regex `xml:"as-path"`
    MatchCommunityRegex *regex `xml:"community"`
    MatchExtendedCommunityRegex *regex `xml:"extended-community"`
    MatchMed string `xml:"med,omitempty"`
    MatchRouteTable string `xml:"route-table,omitempty"`
    MatchAddressPrefix *addPre `xml:"address-prefix"`
    MatchNextHop *util.MemberType `xml:"nexthop"`
    MatchFromPeer *util.MemberType `xml:"from-peer"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        UsedBy: util.StrToMem(e.UsedBy),
    }
    s := ""

    if e.MatchAsPathRegex != "" || e.MatchCommunityRegex != "" || e.MatchExtendedCommunityRegex != "" || e.MatchMed != "" || e.MatchRouteTable != "" || len(e.MatchAddressPrefix) > 0 || len(e.MatchNextHop) > 0 || len(e.MatchFromPeer) > 0 {
        ans.Match = &match_v2{
            MatchMed: e.MatchMed,
            MatchRouteTable: e.MatchRouteTable,
            MatchNextHop: util.StrToMem(e.MatchNextHop),
            MatchFromPeer: util.StrToMem(e.MatchFromPeer),
        }

        if e.MatchAsPathRegex != "" {
            ans.Match.MatchAsPathRegex = &regex{
                Regex: e.MatchAsPathRegex,
            }
        }

        if e.MatchCommunityRegex != "" {
            ans.Match.MatchCommunityRegex = &regex{
                Regex: e.MatchCommunityRegex,
            }
        }

        if e.MatchExtendedCommunityRegex != "" {
            ans.Match.MatchExtendedCommunityRegex = &regex{
                Regex: e.MatchExtendedCommunityRegex,
            }
        }

        if len(e.MatchAddressPrefix) > 0 {
            apList := make([]apEntry, 0, len(e.MatchAddressPrefix))
            for k, v := range e.MatchAddressPrefix {
                apList = append(apList, apEntry{
                    Name: k,
                    Exact: util.YesNo(v),
                })
            }
            ans.Match.MatchAddressPrefix = &addPre{apList}
        }
    }

    switch e.Action {
    case ActionDeny:
        ans.Action = &action{
            Deny: &s,
        }
    case ActionAllow:
        ans.Action = &action{
            Allow: &allow{
            },
        }

        if e.LocalPreference != "" || e.Med != "" || e.NextHop != "" || e.Origin != "" || e.AsPathLimit != 0 || e.AsPathType != "" || e.CommunityType != "" || e.ExtendedCommunityType != "" {
            u := update{
                LocalPreference: e.LocalPreference,
                Med: e.Med,
                NextHop: e.NextHop,
                Origin: e.Origin,
                AsPathLimit: e.AsPathLimit,
            }

            switch e.AsPathType {
            case AsPathTypeNone:
                u.AsPath = &asPath{
                    None: &s,
                }
            case AsPathTypeRemove:
                u.AsPath = &asPath{
                    Remove: &s,
                }
            case AsPathTypePrepend:
                u.AsPath = &asPath{
                    Prepend: e.AsPathValue,
                }
            case AsPathTypeRemoveAndPrepend:
                u.AsPath = &asPath{
                    RemoveAndPrepend: e.AsPathValue,
                }
            }

            switch e.CommunityType {
            case CommunityTypeNone:
                u.Community = &allowCom{
                    None: &s,
                }
            case CommunityTypeRemoveAll:
                u.Community = &allowCom{
                    RemoveAll: &s,
                }
            case CommunityTypeRemoveRegex:
                u.Community = &allowCom{
                    RemoveRegex: e.CommunityValue,
                }
            case CommunityTypeAppend:
                u.Community = &allowCom{
                    Append: util.OneStrToMem(e.CommunityValue),
                }
            case CommunityTypeOverwrite:
                u.Community = &allowCom{
                    Overwrite: util.OneStrToMem(e.CommunityValue),
                }
            }

            switch e.ExtendedCommunityType {
            case CommunityTypeNone:
                u.ExtendedCommunity = &allowCom{
                    None: &s,
                }
            case CommunityTypeRemoveAll:
                u.ExtendedCommunity = &allowCom{
                    RemoveAll: &s,
                }
            case CommunityTypeRemoveRegex:
                u.ExtendedCommunity = &allowCom{
                    RemoveRegex: e.ExtendedCommunityValue,
                }
            case CommunityTypeAppend:
                u.ExtendedCommunity = &allowCom{
                    Append: util.OneStrToMem(e.ExtendedCommunityValue),
                }
            case CommunityTypeOverwrite:
                u.ExtendedCommunity = &allowCom{
                    Overwrite: util.OneStrToMem(e.ExtendedCommunityValue),
                }
            }

            ans.Action.Allow.Update = &u
        }
    }

    return ans
}
