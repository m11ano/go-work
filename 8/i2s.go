package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {

	outValue := reflect.ValueOf(out)
	if outValue.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("out [%T] is not pointer", out)
	}

	return i2sDo(data, outValue.Elem())

}

func i2sDo(data interface{}, out reflect.Value) error {

	value := reflect.ValueOf(data)
	valueType := value.Type()
	switch valueType.Kind() {
	case reflect.Slice:
		if out.Kind() != reflect.Slice {
			return fmt.Errorf("out field [%T] is not slice", out)
		}
		for i := 0; i < value.Len(); i++ {
			v := value.Index(i)
			if v.CanInterface() {
				if out.Len()-1 < i {
					out.Set(reflect.Append(out, reflect.New(out.Type().Elem()).Elem()))
				}
				outItem := out.Index(i)
				if outItem.IsValid() && outItem.CanSet() {
					err := i2sDo(v.Interface(), outItem)
					if err != nil {
						return err
					}
				}
			}
		}
	case reflect.Map:
		if out.Kind() != reflect.Struct {
			return fmt.Errorf("out field [%T] is not struct", out)
		}
		for _, mI := range value.MapKeys() {
			mV := value.MapIndex(mI)
			mIS := mI.String()
			if mV.CanInterface() {
				mVI := mV.Interface()

				outField := out.FieldByName(mIS)
				if outField.IsValid() && outField.CanSet() {
					err := i2sDo(mVI, outField)
					if err != nil {
						return err
					}
				}

			}
		}
	case reflect.Int:
		switch out.Kind() {
		case reflect.Int:
			out.SetInt(value.Int())
		case reflect.Float32, reflect.Float64:
			out.SetFloat(float64(value.Int()))
		default:
			return fmt.Errorf("out field [%T] is not int/float", out)
		}
	case reflect.Float32, reflect.Float64:
		switch out.Kind() {
		case reflect.Int:
			out.SetInt(int64(value.Float()))
		case reflect.Float32, reflect.Float64:
			out.SetFloat(value.Float())
		default:
			return fmt.Errorf("out field [%T] is not int/float", out)
		}
	case reflect.String:
		switch out.Kind() {
		case reflect.String:
			out.SetString(value.String())
		default:
			return fmt.Errorf("out field [%T] is not string", out)
		}

	case reflect.Bool:
		switch out.Kind() {
		case reflect.Bool:
			out.SetBool(value.Bool())
		default:
			return fmt.Errorf("out field [%T] is not bool", out)
		}
	}

	return nil
}
