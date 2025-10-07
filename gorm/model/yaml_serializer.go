package model

import (
	"context"
	"errors"
	"reflect"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm/schema"
)

type YamlMap map[string]interface{}

func (ym *YamlMap) Scan(ctx context.Context, field *schema.Field,
	dst reflect.Value, dbValue interface{}) error {
	if dbValue == nil {
		*ym = nil
		return nil
	}
	bytes, ok := dbValue.([]byte)
	if !ok {
		return errors.New("failed to unmarshal YAML value: source data not bytes")
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(bytes, &m); err != nil {
		return err
	}

	*ym = m
	return nil
}

func (ym YamlMap) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	if ym == nil {
		return nil, nil
	}
	return yaml.Marshal(ym)
}
