package node

type FieldType int

const (
	String FieldType = iota
	Int
	Date
)

func (t FieldType) String() string {
	switch t {
	case String:
		return "String"
	case Int:
		return "Int"
	case Date:
		return "Date"
	default:
		return "Unknown"
	}
}

type Field struct {
	Name        string
	Type        FieldType
	Size        int
	Description string
	Tooltip     string
}
