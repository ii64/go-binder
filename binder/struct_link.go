package binder

import (
	"reflect"
)

func invokeLinInkMethodIfExist(f interface{}) {
	v := reflect.ValueOf(f)
	meth := v.MethodByName("LinkIn")
	if meth.Kind() == reflect.Func {
		// fmt.Printf("inv linkIn meth: %+#v %+#v\n", v, meth)
		meth.Call([]reflect.Value{})
	}
}

func invokeLinkOutMethodIfExist(f interface{}) {
	v := reflect.ValueOf(f)
	meth := v.MethodByName("LinkOut")
	if meth.Kind() == reflect.Func {
		// fmt.Printf("inv linkOut meth: %+#v %+#v\n", v, meth)
		meth.Call([]reflect.Value{})
	}
}

func recopyFieldsIn(src, dst reflect.Value, orig interface{}) func() {
	return func() {
		defer invokeLinInkMethodIfExist(orig)
		// fmt.Printf("innnnnnnnnnn %+#v\n", src)
		for i := 0; i < src.NumField(); i++ {
			val := src.Field(i)
			d := dst.Field(i)
			// fmt.Printf("SRC VAL [%s]: ", "")
			// spew.Dump(val.Interface())
			// fmt.Printf("DST VAL [%s]: ", "")
			// spew.Dump(d.Interface())
			if val.Type().Kind() == reflect.Struct || (val.Type().Kind() == reflect.Ptr && val.Type().Elem().Kind() == reflect.Struct) {
				restructToPtrInitialize(val, d)
				continue
			}
			d.Set(val.Elem())
			// fmt.Printf("DST VAL after: %+#v\n", d)
		}
	}
}

func recopyFieldsOut(src, dst reflect.Value, orig interface{}) func() {
	return func() {
		defer invokeLinkOutMethodIfExist(orig)
		for i := 0; i < src.NumField(); i++ {
			val := src.Field(i)
			d := dst.Field(i)
			if val.Type().Kind() == reflect.Struct || (val.Type().Kind() == reflect.Ptr && val.Type().Elem().Kind() == reflect.Struct) {
				restructToPtrInitialize(d, val)
				continue
			}
			d.Elem().Set(val)
		}
	}
}

func restructToPtr(t reflect.Type) reflect.Type { // return non-pointer type
	// fmt.Printf("instrospect invoked\n")
	fields := []reflect.StructField{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		reTag(f, &f.Tag)
		if f.Type.Kind() == reflect.Struct || (f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct) {
			f.Type = restructToPtr(f.Type)
		}
		f.Type = reflect.PtrTo(f.Type)
		fields = append(fields, f)
	}
	st := reflect.StructOf(fields)
	// fmt.Printf("got %s\n", st)
	return st
}

func restructToPtrInitialize(orig reflect.Value, v reflect.Value) {
	if orig.Kind() == reflect.Ptr {
		orig = orig.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// fmt.Printf("orig f: %+#v\nv f: %+#v\n", orig.Interface(), v.Interface())
	for i := 0; i < orig.NumField(); i++ {
		f := orig.Field(i)
		flink := v.Field(i)

		if f.Type().Kind() == reflect.Struct {
			restructToPtrInitialize(f, flink)
			continue
		}

		// fmt.Printf("restruct: f\n")
		// spew.Dump(f.Interface())
		// fmt.Printf("restruct: flink\n")
		// spew.Dump(flink.Interface())

		if !f.CanAddr() || f.IsNil() {
			fptr := reflect.New(f.Type().Elem())
			fptr.Elem().Set(flink)
			f.Set(fptr)
			continue
		}
		flink.Set(f.Elem())
	}
}

// link creates an intermediate struct of original struct fields
func link(f interface{}) (interface{}, func(), func()) {
	v := reflect.ValueOf(f)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// fmt.Printf("%+#v\n", v)
	fields := []reflect.StructField{}
	t := v.Type()
	// fmt.Printf("%s\n", t)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		reTag(f, &f.Tag)
		if f.Type.Kind() == reflect.Struct || (f.Type.Kind() == reflect.Ptr && f.Type.Elem().Kind() == reflect.Struct) {
			f.Type = restructToPtr(f.Type)
		}
		f.Type = reflect.PtrTo(f.Type)
		fields = append(fields, f)
	}
	// fmt.Printf("%+#v\n", fields)
	st := reflect.StructOf(fields)
	linker := reflect.New(st)
	el := linker.Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		flink := el.Field(i)

		if f.Type().Kind() == reflect.Struct {
			// fmt.Printf("%s %+#v\n", f.Type(), flink)
			fptr := reflect.New(flink.Type().Elem())
			flink.Set(fptr)
			restructToPtrInitialize(flink, f)
			continue
		}
		fptr := reflect.New(f.Type())
		fptr.Elem().Set(f)
		flink.Set(fptr)
	}
	// fmt.Printf("%+#v\n", linker)
	return linker.Interface(), recopyFieldsIn(el, v, f), recopyFieldsOut(v, el, f)
}
