package provider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	rsschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Locationer interface {
	IsValid() error
}

type RuleInfo struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
}

func EncodeLocation(loc Locationer) (string, error) {
	b, err := json.Marshal(loc)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func DecodeLocation(s string, loc Locationer) error {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(b, loc); err != nil {
		return err
	}

	return loc.IsValid()
}

func ProviderParamDescription(desc, defaultValue, envName, jsonName string) string {
	var b strings.Builder

	b.WriteString(desc)

	if defaultValue != "" {
		b.WriteString(fmt.Sprintf(" Default: `%s`.", defaultValue))
	}

	if envName != "" {
		b.WriteString(fmt.Sprintf(" Environment variable: `%s`.", envName))
	}

	if jsonName != "" {
		b.WriteString(fmt.Sprintf(" JSON config file variable: `%s`.", jsonName))
	}

	return b.String()
}

type AncestorType string

const (
	AncestorObjectEntry AncestorType = "object-entry"
	AncestorListEntry   AncestorType = "list-entry"
)

type Ancestor interface {
	AncestorName() string
	EntryName() *string
}

type XpathAncestorError struct {
	name    string
	message string
}

func (o XpathAncestorError) Error() string {
	message := o.message
	message += fmt.Sprintf(": %s", o.name)
	return message
}

func CreateXpathForAttributeWithAncestors(ancestors []Ancestor, attribute string) (string, error) {
	var xpath []string

	createXpathElements := func(attr Ancestor) ([]string, error) {
		elts := []string{"/" + attr.AncestorName()}
		name := attr.EntryName()
		if name != nil {
			elts = append(elts, fmt.Sprintf("/entry[@name=\"%s\"]", *name))

		}

		return elts, nil
	}

	for _, elt := range ancestors {
		xpathElts, err := createXpathElements(elt)
		if err != nil {
			return "", err
		}

		xpath = append(xpath, xpathElts...)
	}

	xpath = append(xpath, "/"+attribute)
	return strings.Join(xpath, ""), nil
}

// TypesObjectToMap converts a Terraform types.Object to a map[string]interface{} for JSON marshaling.
// This allows marshaling partial objects without requiring all struct fields.
// If schemaAttr is provided, validates the object against the schema to catch unknown fields.
func TypesObjectToMap(obj types.Object, schemaAttr ...rsschema.Attribute) (interface{}, error) {
	if obj.IsNull() {
		return nil, nil
	}

	attrs := obj.Attributes()
	result := make(map[string]interface{})

	// Validate against schema if provided
	if len(schemaAttr) > 0 && schemaAttr[0] != nil {
		schema, ok := schemaAttr[0].(rsschema.SingleNestedAttribute)
		if ok {
			// Validate no unknown fields
			for key := range attrs {
				if _, exists := schema.Attributes[key]; !exists {
					return nil, fmt.Errorf("unknown field %q in location object", key)
				}
			}
			// Validate not empty
			if len(attrs) == 0 {
				return nil, fmt.Errorf("location object cannot be empty")
			}
		}
	}

	for key, val := range attrs {
		switch v := val.(type) {
		case types.Object:
			// For nested objects, don't validate (no schema available)
			nested, err := TypesObjectToMap(v)
			if err != nil {
				return nil, err
			}
			result[key] = nested
		case types.String:
			if !v.IsNull() {
				result[key] = v.ValueString()
			}
		case types.Bool:
			if !v.IsNull() {
				result[key] = v.ValueBool()
			}
		case types.Int64:
			if !v.IsNull() {
				result[key] = v.ValueInt64()
			}
		case types.Float64:
			if !v.IsNull() {
				result[key] = v.ValueFloat64()
			}
		case types.List:
			if !v.IsNull() {
				var list []interface{}
				for _, elem := range v.Elements() {
					switch e := elem.(type) {
					case types.Object:
						nested, err := TypesObjectToMap(e)
						if err != nil {
							return nil, err
						}
						list = append(list, nested)
					case types.String:
						if !e.IsNull() {
							list = append(list, e.ValueString())
						}
					default:
						list = append(list, elem)
					}
				}
				result[key] = list
			}
		case types.Map:
			if !v.IsNull() {
				mapResult := make(map[string]interface{})
				for k, mapVal := range v.Elements() {
					switch mv := mapVal.(type) {
					case types.String:
						if !mv.IsNull() {
							mapResult[k] = mv.ValueString()
						}
					default:
						mapResult[k] = mapVal
					}
				}
				result[key] = mapResult
			}
		}
	}

	return result, nil
}

// MapToTypesObject converts a map[string]interface{} to a Terraform types.Object using the provided schema.
// This automatically handles missing fields by creating typed null values.
// Validates that the map doesn't contain unknown fields or is empty.
func MapToTypesObject(data map[string]interface{}, schemaAttr rsschema.Attribute) (types.Object, error) {
	schema, ok := schemaAttr.(rsschema.SingleNestedAttribute)
	if !ok {
		return types.ObjectNull(nil), fmt.Errorf("schema attribute is not a SingleNestedAttribute")
	}

	attrTypes := make(map[string]attr.Type)
	attrValues := make(map[string]attr.Value)

	for name, schemaField := range schema.Attributes {
		attrTypes[name] = schemaField.GetType()

		if val, exists := data[name]; exists && val != nil {
			// Convert JSON value to appropriate terraform type
			switch schemaField.GetType().(type) {
			case basetypes.ObjectType:
				// Recursively handle nested objects
				nestedMap, ok := val.(map[string]interface{})
				if !ok {
					return types.ObjectNull(attrTypes), fmt.Errorf("expected map for nested object %s", name)
				}
				nestedObj, err := MapToTypesObject(nestedMap, schemaField)
				if err != nil {
					return types.ObjectNull(attrTypes), err
				}
				attrValues[name] = nestedObj
			case basetypes.StringType:
				strVal, ok := val.(string)
				if !ok {
					return types.ObjectNull(attrTypes), fmt.Errorf("expected string for field %s", name)
				}
				attrValues[name] = types.StringValue(strVal)
			default:
				// For other types, try to set them as-is
				attrValues[name] = types.StringValue(fmt.Sprintf("%v", val))
			}
		} else {
			// Create typed null for missing fields
			switch schemaField.GetType().(type) {
			case basetypes.ObjectType:
				objType := schemaField.GetType().(basetypes.ObjectType)
				attrValues[name] = types.ObjectNull(objType.AttrTypes)
			default:
				// For non-object types, create a generic null
				attrValues[name] = types.StringNull()
			}
		}
	}

	// Validate no unknown fields in input
	for key := range data {
		if _, exists := schema.Attributes[key]; !exists {
			return types.ObjectNull(attrTypes), fmt.Errorf("unknown field %q in location object", key)
		}
	}

	// Validate not empty
	if len(data) == 0 {
		return types.ObjectNull(attrTypes), fmt.Errorf("location object cannot be empty")
	}

	obj, diags := types.ObjectValue(attrTypes, attrValues)
	if diags.HasError() {
		return types.ObjectNull(attrTypes), fmt.Errorf("failed to create object: %v", diags.Errors())
	}
	return obj, nil
}
