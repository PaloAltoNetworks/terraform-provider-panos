package tags

// These are the color constants you can use in Entry.SetColor().  Note that
// each version of PANOS has added colors, so if you are looking for maximum
// compatibility, only use the first 16 colors (17 including None).
const (
	None = iota
	Red
	Green
	Blue
	Yellow
	Copper
	Orange
	Purple
	Gray
	LightGreen
	Cyan
	LightGray
	BlueGray
	Lime
	Black
	Gold
	Brown
	Olive
	_
	Maroon
	RedOrange
	YellowOrange
	ForestGreen
	TurquoiseBlue
	AzureBlue
	CeruleanBlue
	MidnightBlue
	MediumBlue
	CobaltBlue
	VioletBlue
	BlueViolet
	MediumViolet
	MediumRose
	Lavender
	Orchid
	Thistle
	Peach
	Salmon
	Magenta
	RedViolet
	Mahogany
	BurntSienna
	Chestnut
)

const (
	singular = "tag"
	plural   = "tags"
)
