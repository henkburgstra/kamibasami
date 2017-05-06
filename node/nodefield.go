package node

type FieldType int

const (
	String FieldType = iota
	Int
	Date
	Unknown
)

func (t FieldType) String() string {
	switch t {
	case String:
		return "String"
	case Int:
		return "Int"
	case Date:
		return "Date"
	case Unknown:
		fallthrough
	default:
		return "Unknown"
	}
}

func NewFieldType(t string) FieldType {
	switch t {
	case "String":
		return String
	case "Int":
		return Int
	case "Date":
		return Date
	case "Unknown":
		fallthrough
	default:
		return Unknown
	}
}

type Field struct {
	Name        string
	Type        FieldType
	Size        int
	Default     interface{}
	Description string
	Tooltip     string
}
