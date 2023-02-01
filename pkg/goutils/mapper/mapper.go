package mapper

import (
	"errors"
	"reflect"
	"time"
)

type Mapper interface {
	Map(source, dest interface{}) error
}

func MapperFactory() Mapper {
	return newMapper()
}

func newMapper() Mapper {
	return &mapper{
		supportedCustomDataTypes: []reflect.Type{
			reflect.TypeOf(time.Time{}),
		},
	}
}

type mapper struct {
	supportedCustomDataTypes []reflect.Type
}

func (mapper *mapper) Map(source, dest interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New(mapperPointerTypeErrorDescription)
	}
	sourceValue := reflect.ValueOf(source)
	destValue := reflect.ValueOf(dest).Elem()

	return mapper.mapValues(sourceValue, destValue)
}

func (mapper *mapper) mapValues(sourceVal, destVal reflect.Value) error {
	destType := destVal.Type()
	destKind := destType.Kind()

	if mapper.isSupportedCustomDataType(destType) && destType == sourceVal.Type() {
		destVal.Set(sourceVal)
		return nil
	}

	if destKind == reflect.Struct {
		if sourceVal.Type().Kind() == reflect.Ptr {
			if sourceVal.IsNil() {
				// If source is nil, it maps to an empty struct
				sourceVal = reflect.New(sourceVal.Type().Elem())
			}
			sourceVal = sourceVal.Elem()
		}

		for i := 0; i < destVal.NumField(); i++ {
			if err := mapper.mapField(sourceVal, destVal, i); err != nil {
				return err
			}
		}
	} else if destType == sourceVal.Type() {
		destVal.Set(sourceVal)
	} else if destType.Kind() == reflect.Ptr {
		if mapper.valueIsNil(sourceVal) {
			return nil
		}

		val := reflect.New(destType.Elem())
		if err := mapper.mapValues(sourceVal, val.Elem()); err != nil {
			return err
		}

		destVal.Set(val)
	} else if destType.Kind() == reflect.Slice {
		return mapper.mapSlice(sourceVal, destVal)
	} else {
		return errors.New(mapperKindNotSupportedErrorDescription)
	}
	return nil
}

func (mapper *mapper) mapSlice(sourceVal, destVal reflect.Value) error {
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)
	for j := 0; j < length; j++ {
		val := reflect.New(destType.Elem()).Elem()
		if err := mapper.mapValues(sourceVal.Index(j), val); err != nil {
			return err
		}

		target.Index(j).Set(val)
	}

	if length == 0 {
		if err := mapper.verifyArrayTypesAreCompatible(sourceVal, destVal); err != nil {
			return err
		}
	}
	destVal.Set(target)

	return nil
}

func (mapper *mapper) verifyArrayTypesAreCompatible(sourceVal, destVal reflect.Value) error {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	return mapper.mapValues(dummySource, dummyDest.Elem())
}

func (mapper *mapper) mapField(source, destVal reflect.Value, i int) error {
	destType := destVal.Type()
	fieldName := destType.Field(i).Name
	defer func() {
		if r := recover(); r != nil {
			// logger.Info(fmt.Sprintf("Error mapping field: %s. DestType: %v. SourceType: %v. Error: %v", fieldName, destType, source.Type(), r))
		}
	}()

	destField := destVal.Field(i)
	if destType.Field(i).Anonymous && mapper.isSupportedCustomDataType(destType.Field(i).Type) {
		source := source.FieldByName(fieldName)
		return mapper.mapValues(source, destField)
	} else {
		if mapper.valueIsContainedInNilEmbeddedType(source, fieldName) {
			return nil
		}
		sourceField := source.FieldByName(fieldName)

		if (sourceField == reflect.Value{}) {
			if destField.Kind() == reflect.Struct {
				return mapper.mapValues(source, destField)
			} else {
				for i := 0; i < source.NumField(); i++ {
					if source.Field(i).Kind() != reflect.Struct {
						continue
					}
					if sourceField = source.Field(i).FieldByName(fieldName); (sourceField != reflect.Value{}) {
						break
					}
				}
			}
		}
		return mapper.mapValues(sourceField, destField)
	}
}

func (mapper *mapper) isSupportedCustomDataType(dataType reflect.Type) bool {
	for _, t := range mapper.supportedCustomDataTypes {
		if dataType == t {
			return true
		}
	}
	return false
}

func (mapper *mapper) valueIsNil(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Ptr && value.IsNil()
}

func (mapper *mapper) valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if mapper.valueIsNil(parentField) {
			return true
		}
	}
	return false
}
