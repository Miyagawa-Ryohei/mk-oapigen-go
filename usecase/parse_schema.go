package usecase

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
	"mk-oapigen-go/entity"
	"strings"
)

type SchemaBuilder struct {
}

func (b SchemaBuilder) BuildPropertyTag(name string, schemaRef *openapi3.SchemaRef, structTagType string) *entity.PropertyTag {
	options := []string{name}
	if schemaRef.Value.Nullable {
		options = append(options, "omitempty")
	}
	return &entity.PropertyTag{
		Name:    name,
		Type:    structTagType,
		Options: options,
	}
}

func (b SchemaBuilder) BuildPropertySchema(name string, schemaRef *openapi3.SchemaRef) entity.PropertySchema {
	itemType := b.detectType(schemaRef)
	if itemType == "array" || itemType == "object" {
		itemType = b.convNameFromRef(schemaRef.Ref)
	}
	if schemaRef.Value.Nullable {
		itemType = "*" + itemType
	}

	property := entity.PropertySchema{
		Name:      strcase.ToCamel(name),
		Type:      itemType,
		StructRef: b.BuildObjectSchema(schemaRef),
		ArrayRef:  b.BuildArraySchema(schemaRef),
		Tag:       b.BuildPropertyTag(name, schemaRef, "json"),
	}
	return property
}

func (b SchemaBuilder) BuildObjectSchema(schemaRef *openapi3.SchemaRef) *entity.StructSchema {
	if schemaRef.Value.Type != "object" {
		return nil
	}
	properties := []entity.PropertySchema{}
	for name, p := range schemaRef.Value.Properties {
		properties = append(properties, b.BuildPropertySchema(name, p))
	}
	ret := &entity.StructSchema{
		Name:         b.convNameFromRef(schemaRef.Ref),
		Type:         b.detectType(schemaRef),
		PropertyList: properties,
	}

	return ret
}

func (b SchemaBuilder) BuildPropertySchemaFromParameter(parameter *openapi3.ParameterRef) entity.PropertySchema {
	schema := parameter.Value.Schema
	itemType := b.detectType(schema)
	if itemType == "array" || itemType == "object" {
		itemType = b.convNameFromRef(schema.Ref)
	}
	if schema.Value.Nullable {
		itemType = "*" + itemType
	}

	property := entity.PropertySchema{
		Name:      strcase.ToCamel(parameter.Value.Name),
		Type:      itemType,
		StructRef: b.BuildObjectSchema(schema),
		ArrayRef:  b.BuildArraySchema(schema),
		Tag:       b.BuildPropertyTag(parameter.Value.Name, schema, parameter.Value.In),
	}
	return property
}

func (b SchemaBuilder) BuildObjectSchemaFromParam(name string, parameters *openapi3.Parameters) *entity.StructSchema {
	if len(*parameters) == 0 {
		return nil
	}
	properties := []entity.PropertySchema{}
	for _, p := range *parameters {
		properties = append(properties, b.BuildPropertySchemaFromParameter(p))
	}
	ret := &entity.StructSchema{
		Name:         name,
		Type:         name,
		PropertyList: properties,
	}

	return ret
}

func (b SchemaBuilder) BuildArraySchema(schemaRef *openapi3.SchemaRef) *entity.ArraySchema {
	if schemaRef == nil {
		return nil
	}
	if schemaRef.Value.Type != "array" {
		return nil
	}
	itemType := b.detectType(schemaRef.Value.Items)
	if itemType == "array" || itemType == "object" {
		itemType = b.convNameFromRef(schemaRef.Value.Items.Ref)
	}

	ret := &entity.ArraySchema{
		Name:      b.convNameFromRef(schemaRef.Ref),
		ItemType:  itemType,
		StructRef: b.BuildObjectSchema(schemaRef.Value.Items),
		ArrayRef:  b.BuildArraySchema(schemaRef.Value.Items),
	}

	return ret
}

func (b SchemaBuilder) BuildSchemaFromRequestBody(bodyRef *openapi3.RequestBodyRef) *entity.Schema {
	for _, cont := range bodyRef.Value.Content {
		return &entity.Schema{
			Name:      b.convNameFromRef(cont.Schema.Ref),
			Type:      b.detectType(cont.Schema),
			StructRef: b.BuildObjectSchema(cont.Schema),
			ArrayRef:  b.BuildArraySchema(cont.Schema),
		}
	}
	return nil
}

func (b SchemaBuilder) BuildSchemaFromParameters(params *openapi3.Parameters, methodName string) entity.Schema {
	return entity.Schema{
		Name:      methodName + "Param",
		Type:      "object",
		StructRef: b.BuildObjectSchemaFromParam(methodName+"Param", params),
	}
}

func (b SchemaBuilder) BuildSchemaFromRespRef(respRef *openapi3.ResponseRef) *entity.Schema {
	for _, cont := range respRef.Value.Content {
		return &entity.Schema{
			Name:      b.convNameFromRef(cont.Schema.Ref),
			Type:      b.detectType(cont.Schema),
			StructRef: b.BuildObjectSchema(cont.Schema),
			ArrayRef:  b.BuildArraySchema(cont.Schema),
		}
	}
	return nil
}

func (b SchemaBuilder) detectType(schema *openapi3.SchemaRef) string {
	switch schema.Value.Type {
	case "object":
	case "array":
		return b.convNameFromRef(schema.Ref)
	case "string":
		if schema.Value.Format == "date-time" || schema.Value.Format == "date" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if schema.Value.Format == "int64" || schema.Value.Format == "int32" {
			return schema.Value.Format
		}
		return "int64"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	}
	return schema.Value.Type
}

func (b SchemaBuilder) convNameFromRef(ref string) string {
	split := strings.Split(ref, "/")
	return split[len(split)-1]
}
