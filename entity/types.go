package entity

type Property interface {
	GetName() string
	GetType() string
}

type PropertySchema struct {
	Name      string
	Type      string
	StructRef *StructSchema
	ArrayRef  *ArraySchema
	Tag       *PropertyTag
}

func (p PropertySchema) GetName() string { return p.Name }
func (p PropertySchema) GetType() string { return p.Type }

type PropertyTag struct {
	Name    string
	Type    string
	Options []string
}

type StructSchema struct {
	Name         string
	Type         string
	PropertyList []PropertySchema
}

type ArraySchema struct {
	Name      string
	ItemType  string
	StructRef *StructSchema
	ArrayRef  *ArraySchema
}

type Schema struct {
	Name      string
	Type      string
	StructRef *StructSchema
	ArrayRef  *ArraySchema
}

func (p Schema) GetName() string { return p.Name }
func (p Schema) GetType() string { return p.Type }
