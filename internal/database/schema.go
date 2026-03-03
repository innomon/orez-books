package database

type FieldType string

const (
	FieldTypeData         FieldType = "Data"
	FieldTypeSelect       FieldType = "Select"
	FieldTypeLink         FieldType = "Link"
	FieldTypeDate         FieldType = "Date"
	FieldTypeDatetime     FieldType = "Datetime"
	FieldTypeTable        FieldType = "Table"
	FieldTypeAutoComplete FieldType = "AutoComplete"
	FieldTypeCheck        FieldType = "Check"
	FieldTypeAttachImage  FieldType = "AttachImage"
	FieldTypeDynamicLink  FieldType = "DynamicLink"
	FieldTypeInt          FieldType = "Int"
	FieldTypeFloat        FieldType = "Float"
	FieldTypeCurrency     FieldType = "Currency"
	FieldTypeText         FieldType = "Text"
	FieldTypeColor        FieldType = "Color"
	FieldTypeButton       FieldType = "Button"
	FieldTypeAttachment   FieldType = "Attachment"
)

type SelectOption struct {
	Value interface{} `json:"value"`
	Label interface{} `json:"label"`
}

type Field struct {
	Fieldname   string         `json:"fieldname"`
	Fieldtype   FieldType      `json:"fieldtype"`
	Label       string         `json:"label"`
	SchemaName  string         `json:"schemaName,omitempty"`
	Required    bool           `json:"required,omitempty"`
	Hidden      bool           `json:"hidden,omitempty"`
	Invisible   bool           `json:"invisible,omitempty"`
	ReadOnly    bool           `json:"readOnly,omitempty"`
	Description string         `json:"description,omitempty"`
	Default     interface{}    `json:"default,omitempty"`
	Placeholder string         `json:"placeholder,omitempty"`
	GroupBy     string         `json:"groupBy,omitempty"`
	Meta        bool           `json:"meta,omitempty"`
	Filter      bool           `json:"filter,omitempty"`
	Computed    bool           `json:"computed,omitempty"`
	Section     string         `json:"section,omitempty"`
	Tab         string         `json:"tab,omitempty"`
	Abstract    interface{}    `json:"abstract,omitempty"`
	IsCustom    bool           `json:"isCustom,omitempty"`
	Bold        bool           `json:"bold,omitempty"`
	SubLabel    string         `json:"sub_label,omitempty"`
	Options     []SelectOption `json:"options,omitempty"`
	Target      string         `json:"target,omitempty"`
	References  string         `json:"references,omitempty"`
	MinValue    float64        `json:"minvalue,omitempty"`
	MaxValue    float64        `json:"maxvalue,omitempty"`
	Rows        int            `json:"rows,omitempty"`
}

type Schema struct {
	Name            string   `json:"name"`
	Label           string   `json:"label"`
	Fields          []Field  `json:"fields"`
	IsTree          bool     `json:"isTree,omitempty"`
	Extends         string   `json:"extends,omitempty"`
	IsChild         bool     `json:"isChild,omitempty"`
	IsSingle        bool     `json:"isSingle,omitempty"`
	IsAbstract      bool     `json:"isAbstract,omitempty"`
	TableFields     []string `json:"tableFields,omitempty"`
	IsSubmittable   bool     `json:"isSubmittable,omitempty"`
	KeywordFields   []string `json:"keywordFields,omitempty"`
	QuickEditFields []string `json:"quickEditFields,omitempty"`
	LinkDisplayField string   `json:"linkDisplayField,omitempty"`
	Create          bool     `json:"create,omitempty"`
	Naming          string   `json:"naming,omitempty"`
	TitleField      string   `json:"titleField,omitempty"`
	RemoveFields    []string `json:"removeFields,omitempty"`
}

type SchemaMap map[string]*Schema
