package usecase

import (
	_ "embed"
	"mk-oapigen-go/entity"
	"testing"
)

var TestStringType1 = entity.Schema{
	Name: "TestStringType1",
	Type: "string",
}

var TestStruct1 = entity.StructSchema{
	Name: "TestStruct1",
	Type: "object",
	PropertyList: []entity.PropertySchema{
		{
			Name: "TestProperty1",
			Type: "Date",
			Tag: &entity.PropertyTag{
				Name:    "test_prop1",
				Type:    "json",
				Options: []string{},
			},
		},
		{
			Name: "TestProperty2",
			Type: "TestStruct2_1",
			Tag: &entity.PropertyTag{
				Name:    "test_prop2",
				Type:    "json",
				Options: []string{"omitempty"},
			},
			StructRef: &TestObject2,
		}, {
			Name: "TestProperty3",
			Type: "TestArray3",
			Tag: &entity.PropertyTag{
				Name:    "test_prop2",
				Type:    "json",
				Options: []string{"omitempty"},
			},
			ArrayRef: &TestArray3,
		},
	},
}

var TestObject2 = entity.StructSchema{
	Name: "TestStruct2_1",
	Type: "object",
	PropertyList: []entity.PropertySchema{
		{
			Name: "TestProperty2_1",
			Type: "string",
			Tag: &entity.PropertyTag{
				Name:    "test_prop2_1",
				Type:    "json",
				Options: []string{},
			},
		},
		{
			Name: "TestProperty2_2",
			Type: "string",
			Tag: &entity.PropertyTag{
				Name:    "test_prop2_1",
				Type:    "json",
				Options: []string{},
			},
		},
	},
}

var TestArray1 = entity.ArraySchema{
	Name:     "TestArray1",
	ItemType: "int64",
}

var TestArray2 = entity.ArraySchema{
	Name:      "TestArray2",
	ItemType:  "TestObject2",
	StructRef: &TestObject2,
}

var TestArray3 = entity.ArraySchema{
	Name:     "TestArray3",
	ItemType: "array",
	ArrayRef: &TestArray1,
}

var TestRequest1 = entity.Schema{
	Name:     "TestRequest1",
	Type:     "array",
	ArrayRef: &TestArray1,
}

var TestRequest2 = entity.Schema{
	Name:     "TestRequest2",
	Type:     "array",
	ArrayRef: &TestArray2,
}
var TestRequest3 = entity.Schema{
	Name:      "TestRequest3",
	Type:      "object",
	StructRef: &TestStruct1,
}

func TestPrintType(t *testing.T) {
	gen := NewTypeGenerator("entity")
	gen.PrintType([]entity.Schema{
		TestStringType1,
		TestRequest1,
		TestRequest2,
		TestRequest3,
	})
}
