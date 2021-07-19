package binder

import (
	"flag"
	"reflect"
)

func defaultRegisterCmdArgsFlagStd(parent string, fieldType reflect.StructField, fieldValue reflect.Value) {
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
		argName = argxName
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

	// unwind types until last pointer
	// for fieldValue.Kind() == reflect.Ptr {
	// 	if t := fieldValue.Elem(); t.Kind() != reflect.Ptr {
	// 		break
	// 	} else {
	// 		fieldValue = t
	// 		// fmt.Printf("unwind --> %+#v\n", fieldValue.Interface())
	// 	}
	// }

	switch val := fieldValue.Interface().(type) {
	case *bool:
		flag.BoolVar(val, argName, *val, argUsage)
	case *int:
		flag.IntVar(val, argName, *val, argUsage)
	case *int64:
		flag.Int64Var(val, argName, *val, argUsage)
	case *float64:
		flag.Float64Var(val, argName, *val, argUsage)
	case *string:
		flag.StringVar(val, argName, *val, argUsage)

	}
}

func defaultRegisterConf(parent string, fieldType reflect.StructField, fieldValue reflect.Value) {
	// fmt.Printf("CONF %s: %s [%s] has value %+#v  [%s]\n", parent, fieldType.Name, fieldType.Type, fieldValue, fieldType.Tag)
	// fmt.Printf("json tag:%q\n", fieldType.Tag.Get("json"))
}

func defaultLoadConfig(mc *MappedConfiguration) error {
	return nil
}

func defaultSaveConfig(mc *MappedConfiguration) error {
	return nil
}
