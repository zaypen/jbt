package util

import "reflect"

func In(v interface{}, in interface{}) bool {
	val := reflect.Indirect(reflect.ValueOf(in))
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if v == val.Index(i).Interface() {
				return true
			}
		}
	}
	return false
}
