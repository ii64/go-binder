package binder

import (
	"flag"
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

func defaultRegisterCmdArgsFlagStd(parent string, fieldType reflect.StructField, fieldValue reflect.Value) (err error) {
	// fmt.Printf("ARGS %s: %s [%s] has value %+#v  [%s]\n", parent, fieldType.Name, fieldType.Type, fieldValue, fieldType.Tag)
	argxName := fieldType.Tag.Get("argx")
	argName := fieldType.Tag.Get("arg")

	argxNameParsed := SplitTagValue(argxName)
	argNameParsed := SplitTagValue(argName)

	if len(argNameParsed) > 0 {
		argName = argNameParsed[0]
	}

	if len(argxNameParsed) > 0 && argxNameParsed[0] != "" {
		argxName = argxNameParsed[0]
		if argxName == "-" {
			return // skip flag binding if argx:"-"
		}
		argName = argxName
	} else if argName == "-" {
		return // skip flag binding if arg:"-"
	} else if argName == "" {
		if parent != "" {
			argName = parent + "." + fieldType.Name
		} else {
			argName = fieldType.Name
		}
	} else {
		if parent != "" {
			argName = parent + "." + argName
		} else {
			argName = fieldType.Name
		}
	}
	argUsage := fieldType.Tag.Get("usage")

	var valueKind = fieldValue // has value
	if fieldValue.Kind() == reflect.Pointer {
		valueKind = fieldValue.Elem()
	}
	value := fieldValue.Interface() // has reference
	switch valueKind.Kind() {
	case reflect.Bool:
		ref := (*bool)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.BoolVar(ref, argName, *ref, argUsage)

	case reflect.Int:
		ref := (*int)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.IntVar(ref, argName, *ref, argUsage)
	case reflect.Int64:
		ref := (*int64)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.Int64Var(ref, argName, *ref, argUsage)

	case reflect.Uint:
		ref := (*uint)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.UintVar(ref, argName, *ref, argUsage)
	case reflect.Uint64:
		ref := (*uint64)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.Uint64Var(ref, argName, *ref, argUsage)

	case reflect.Float64:
		ref := (*float64)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.Float64Var(ref, argName, *ref, argUsage)

	case reflect.String:
		ref := (*string)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
		flag.StringVar(ref, argName, *ref, argUsage)

	case reflect.Slice:
		flag.Var(newDefaultSliceVar(valueKind), argName, argUsage)
	case reflect.Array:
		flag.Var(newDefaultArrayVar(valueKind), argName, argUsage)

	default:
		switch valueKind.Interface().(type) {
		case time.Duration:
			ref := (*time.Duration)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
			flag.DurationVar(ref, argName, *ref, argUsage)
		case func():
			ref := (*func(string) error)((*[2]unsafe.Pointer)(unsafe.Pointer(&value))[1])
			flag.Func(argName, argUsage, *ref)
		default:
			switch ref := value.(type) {
			case flag.Value:
				flag.Var(ref, argName, argUsage)
			}
		}
	}

	return
}

func defaultLoadConfig(mc *MappedConfiguration) error {
	return nil
}

func defaultSaveConfig(mc *MappedConfiguration) error {
	return nil
}

type defaultSliceVar struct {
	value     reflect.Value
	container reflect.Type
}

func newDefaultSliceVar(value reflect.Value) *defaultSliceVar {
	container := value.Type()
	return &defaultSliceVar{value: value, container: container}
}

func (s *defaultSliceVar) Set(value string) error {
	val := convertStringToType(value, s.container.Elem(), nil)
	s.value.Set(reflect.Append(s.value, val))
	return nil
}

func (s *defaultSliceVar) String() string {
	// TODO: tbd.
	return ""
}

type defaultArrayVar struct {
	value     reflect.Value
	container reflect.Type
	cur       int
}

func newDefaultArrayVar(value reflect.Value) *defaultArrayVar {
	container := value.Type()
	return &defaultArrayVar{value: value, container: container}
}

func (s *defaultArrayVar) Set(value string) error {
	if s.cur > int(s.container.Len()) {
		return fmt.Errorf("array set exceeded container size")
	}
	val := convertStringToType(value, s.container.Elem(), nil)
	s.value.Index(s.cur).Set(val)
	s.cur++
	return nil
}

func (s *defaultArrayVar) String() string {
	// TODO: tbd.
	return ""
}
