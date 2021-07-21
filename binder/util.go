package binder

import (
	"reflect"
)

func UnwindType(v reflect.Type, usingLastPtr bool) (r reflect.Type, pc int) {
	r = v
	for r.Kind() == reflect.Ptr {
		pc++
		t := r.Elem()
		if t.Kind() != reflect.Ptr {
			if !usingLastPtr {
				r = t
			}
			break
		} else {
			r = t
		}
	}
	return
}

func WindType(v reflect.Type, pc int) (r reflect.Type) {
	r = v
	for i := 0; i < pc; i++ {
		r = reflect.PtrTo(r)
	}
	return
}

func FillValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// new pointer type
			nc := reflect.New(v.Type().Elem())
			v.Set(nc)
		}
		v = FillValue(v.Elem())
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			t := f.Type()
			if t.Kind() == reflect.Ptr {
				FillValue(f)
			}
		}
	}
	return v
}

func UnwindValue(v reflect.Value, fillNil bool, usingLastPtr bool) (r reflect.Value, pc int) {
	for v.Kind() == reflect.Ptr {
		pc++
		if v.IsNil() {
			if fillNil {
				FillValue(v)
			} else {
				r = v
				return
			}
		}
		t := v.Elem()
		if t.Kind() != reflect.Ptr {
			if !usingLastPtr {
				v = t
			}
			break
		} else {
			v = t
		}
	}
	r = v
	return
}

func CopyValue(src, dst reflect.Value) {
	var srcPc, dstPc int
	src, srcPc = UnwindValue(src, true, false)
	dst, dstPc = UnwindValue(dst, true, false)
	_, _ = srcPc, dstPc
	if src.Kind() == dst.Kind() {
		if src.Kind() == reflect.Struct {
			for i := 0; i < src.NumField(); i++ {
				sf := src.Field(i)
				sf, _ = UnwindValue(sf, true, false)
				df := dst.Field(i)
				df, _ = UnwindValue(df, true, false)

				if src.Kind() == reflect.Struct {
					CopyValue(sf, df)
					continue
				}

				df.Set(sf)
			}
		} else {
			dst.Set(src)
		}
	}
}
