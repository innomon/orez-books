package database

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

//go:embed schemas/app/*.json schemas/core/*.json schemas/meta/*.json schemas/regional/**/*.json
var schemaFS embed.FS

type SchemaLoader struct {
	SchemaMap SchemaMap
}

func NewSchemaLoader() *SchemaLoader {
	return &SchemaLoader{
		SchemaMap: make(SchemaMap),
	}
}

func (l *SchemaLoader) LoadSchemas(countryCode string) error {
	coreSchemas, err := l.loadDirectory("schemas/core")
	if err != nil {
		return err
	}

	appSchemas, err := l.loadDirectory("schemas/app")
	if err != nil {
		return err
	}

	metaSchemas, err := l.loadDirectory("schemas/meta")
	if err != nil {
		return err
	}

	// Basic regional schema loading
	regionalPath := fmt.Sprintf("schemas/regional/%s", countryCode)
	regionalSchemas, _ := l.loadDirectory(regionalPath)

	// Combine app and regional
	for name, schema := range regionalSchemas {
		if appSchema, ok := appSchemas[name]; ok {
			appSchemas[name] = l.combineSchemas(appSchema, schema)
		} else {
			appSchemas[name] = schema
		}
	}

	// Merge all into SchemaMap
	for name, schema := range appSchemas {
		l.SchemaMap[name] = schema
	}
	for name, schema := range coreSchemas {
		l.SchemaMap[name] = schema
	}

	// Handle inheritance (Abstract schemas)
	l.resolveInheritance()

	// Add Meta Fields
	l.addMetaFields(metaSchemas)

	// Finalize fields
	l.finalizeFields()

	return nil
}

func (l *SchemaLoader) loadDirectory(path string) (SchemaMap, error) {
	entries, err := schemaFS.ReadDir(path)
	if err != nil {
		return nil, nil // Directory might not exist for regional
	}

	schemas := make(SchemaMap)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		data, err := schemaFS.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}

		var schema Schema
		if err := json.Unmarshal(data, &schema); err != nil {
			return nil, fmt.Errorf("error unmarshaling %s: %w", entry.Name(), err)
		}
		schemas[schema.Name] = &schema
	}
	return schemas, nil
}

func (l *SchemaLoader) combineSchemas(base *Schema, override *Schema) *Schema {
	// Deep copy base and override
	// For simplicity in this port, we merge fields by fieldname
	fieldMap := make(map[string]Field)
	for _, f := range base.Fields {
		fieldMap[f.Fieldname] = f
	}
	for _, f := range override.Fields {
		fieldMap[f.Fieldname] = f
	}

	newFields := make([]Field, 0, len(fieldMap))
	// Maintain some order or just append
	for _, f := range fieldMap {
		newFields = append(newFields, f)
	}

	base.Fields = newFields
	// Merge other properties if needed
	return base
}

func (l *SchemaLoader) resolveInheritance() {
	// Port getAbstractCombinedSchemas logic
	for name, schema := range l.SchemaMap {
		if schema.Extends != "" {
			if parent, ok := l.SchemaMap[schema.Extends]; ok {
				l.SchemaMap[name] = l.combineSchemas(parent, schema)
			}
		}
	}

	// Remove abstract schemas that are not extended (or just all abstract schemas)
	for name, schema := range l.SchemaMap {
		if schema.IsAbstract {
			delete(l.SchemaMap, name)
		}
	}
}

func (l *SchemaLoader) addMetaFields(metaSchemas SchemaMap) {
	base := metaSchemas["base"]
	tree := l.combineSchemas(metaSchemas["tree"], base)
	child := metaSchemas["child"]
	submittable := l.combineSchemas(metaSchemas["submittable"], base)
	// Simplified submittableTree
	submittableTree := l.combineSchemas(tree, metaSchemas["submittable"])

	for _, schema := range l.SchemaMap {
		if schema.IsSingle {
			continue
		}

		var metaFields []Field
		if schema.IsTree && schema.IsSubmittable {
			metaFields = submittableTree.Fields
		} else if schema.IsTree {
			metaFields = tree.Fields
		} else if schema.IsSubmittable {
			metaFields = submittable.Fields
		} else if schema.IsChild {
			metaFields = child.Fields
		} else {
			metaFields = base.Fields
		}

		schema.Fields = append(schema.Fields, metaFields...)
	}
}

func (l *SchemaLoader) finalizeFields() {
	for _, schema := range l.SchemaMap {
		// Add name field if missing and not single
		if !schema.IsSingle {
			hasName := false
			for _, f := range schema.Fields {
				if f.Fieldname == "name" {
					hasName = true
					break
				}
			}
			if !hasName {
				nameField := Field{
					Fieldname: "name",
					Label:     "ID",
					Fieldtype: FieldTypeData,
					Required:  true,
					ReadOnly:  true,
				}
				schema.Fields = append([]Field{nameField}, schema.Fields...)
			}
		}

		// Set SchemaName on fields
		for i := range schema.Fields {
			schema.Fields[i].SchemaName = schema.Name
		}

		// Set default TitleField
		if schema.TitleField == "" {
			schema.TitleField = "name"
		}
	}
}
