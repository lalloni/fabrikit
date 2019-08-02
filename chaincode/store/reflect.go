package store

import (
	"reflect"
)

func FieldClear(name string) MutatorFunc {
	return func(v interface{}) {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}
		f := vv.FieldByName(name)
		f.Set(reflect.Zero(f.Type()))
	}
}

func FieldSetter(name string) SetterFunc {
	return func(v interface{}, w interface{}) {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}
		vv.FieldByName(name).Set(reflect.ValueOf(w))
	}
}

func FieldGetter(name string) GetterFunc {
	return func(v interface{}) interface{} {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}
		return vv.FieldByName(name).Interface()
	}
}

func MapEnumerator() EnumeratorFunc {
	return func(v interface{}) []Item {
		items := []Item{}
		mv := reflect.ValueOf(v)
		vs := mv.MapKeys()
		for _, v := range vs {
			items = append(items, NewItem(v.String(), mv.MapIndex(v).Interface()))
		}
		return items
	}
}

func MapCollector() CollectorFunc {
	return func(v interface{}, item Item) {
		reflect.ValueOf(v).SetMapIndex(reflect.ValueOf(item.Identifier), reflect.ValueOf(item.Value))
	}
}

func ValueCreator(ref interface{}) CreatorFunc {
	t := reflect.TypeOf(ref).Elem()
	return func() interface{} {
		return reflect.New(t).Interface()
	}
}
