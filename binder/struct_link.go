package binder

import (
	"reflect"
)

func invokeLinInkMethodIfExist(f interface{}) {
	v := reflect.ValueOf(f)
	meth := v.MethodByName("LinkIn")
	if meth.Kind() == reflect.Func {
		meth.Call([]reflect.Value{})
	}
}

func invokeLinkOutMethodIfExist(f interface{}) {
	v := reflect.ValueOf(f)
	meth := v.MethodByName("LinkOut")
	if meth.Kind() == reflect.Func {
		meth.Call([]reflect.Value{})
	}
}

func recopyFieldsIn(src, dst reflect.Value, orig interface{}) func() {
	return func() {
		defer invokeLinInkMethodIfExist(orig)
		// fmt.Printf("innnnnnnnnnn %+#v\n", src)
		for i := 0; i < src.NumField(); i++ {
			sf := src.Field(i)
			df := dst.Field(i)

			if sf.Kind() == reflect.Ptr {
				sf, _ = UnwindValue(sf, false, true)
				// sft, _ := UnwindType(sf.Type(), false)
				// _ = sft
				if sf.IsNil() {
					continue
				}
			}

			// fill nil pointer
			if df.Kind() == reflect.Ptr && df.IsNil() {
				FillValue(df)
			}

			CopyValue(sf, df)
		}
	}
}

func recopyFieldsOut(src, dst reflect.Value, orig interface{}) func() {
	return func() {
		defer invokeLinkOutMethodIfExist(orig)
		for i := 0; i < src.NumField(); i++ {
			sf := src.Field(i)
			df := dst.Field(i)

			if df.Kind() == reflect.Ptr {
				df, _ = UnwindValue(df, false, true) // returned nil Value
				// dft, _ := UnwindType(df.Type(), false)
				// _ = dft
				if df.IsNil() {
					// todo(ii64): having nil
					continue
				}
			}

			// fill nil pointer
			if sf.Kind() == reflect.Ptr && sf.IsNil() {
				FillValue(sf)
			}

			CopyValue(df, sf)
		}
	}
}

func restructToPtr(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Struct {
		fields := []reflect.StructField{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			t := f.Type
			reTag(f, &f.Tag)

			var pc int
			t, pc = UnwindType(t, false)
			if t.Kind() == reflect.Struct {
				t = restructToPtr(t)
			}
			t = WindType(t, pc)
			f.Type = reflect.PtrTo(t)
			fields = append(fields, f)
		}
		t = reflect.StructOf(fields)
	}
	return t
}

// Link creates an intermediate struct of original struct fields
func Link(f interface{}) (interface{}, func(), func()) {
	v := reflect.ValueOf(f)
	if f == nil || v.Type().Kind() == reflect.Ptr {
		if f == nil || v.IsNil() {
			panic("Link to nil pointer is not possible")
		}
		v = v.Elem()
	}
	fields := []reflect.StructField{}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		t := f.Type
		reTag(f, &f.Tag)

		var pc = 0
		t, pc = UnwindType(t, false)
		if t.Kind() == reflect.Struct {
			t = restructToPtr(t)
		}
		t = WindType(t, pc)

		f.Type = reflect.PtrTo(t)

		fields = append(fields, f)
	}

	st := reflect.StructOf(fields)
	linker := reflect.New(st)

	// FillValue(linker)
	el := linker.Elem()
	// linker is *struct
	// recopyFieldsIn   linker -> component (actual config struct)
	// recopyFieldsOut  linker <- component (actual config struct)
	in := recopyFieldsIn(el, v, f)
	out := recopyFieldsOut(el, v, f)
	defer out()
	return linker.Interface(), in, out
}
