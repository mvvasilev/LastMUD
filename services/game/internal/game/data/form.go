package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
)

type FormFieldType byte

const (
	FormFieldTypeString FormFieldType = iota
	FormFieldTypeNumber
	FormFieldTypeDecimal
	FormFieldTypeYN
)

type FormFieldTypeComponent struct {
	FieldType FormFieldType
}

func (c FormFieldTypeComponent) Type() ecs.ComponentType {
	return TypeFormFieldType
}

type FormFieldValueStringComponent struct {
	Value string
}

func (c FormFieldValueStringComponent) Type() ecs.ComponentType {
	return TypeFormFieldValueString
}

type FormFieldValueNumberComponent struct {
	Value int
}

func (c FormFieldValueNumberComponent) Type() ecs.ComponentType {
	return TypeFormFieldValueNumber
}

type FormFieldValueDecimalComponent struct {
	Value float32
}

func (c FormFieldValueDecimalComponent) Type() ecs.ComponentType {
	return TypeFormFieldValueDecimal
}

type YesNo bool

const (
	Yes YesNo = true
	No  YesNo = false
)

type FormFieldValueYesNoComponent struct {
	Value YesNo
}

func (c FormFieldValueYesNoComponent) Type() ecs.ComponentType {
	return TypeFormFieldValueYesNo
}
