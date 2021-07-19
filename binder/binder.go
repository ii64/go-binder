package binder

import (
	"fmt"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

var (
	registered                       = map[string]*item{}
	RegisterCmdArgs   RegisterFunc   = defaultRegisterCmdArgsFlagStd
	mappedConf                       = MappedConfiguration{}
	mappedConfUpdater                = map[string]*funcUpdater{}
	LoadConfig        LoadConfigFunc = defaultLoadConfig
	SaveConfig        SaveConfigFunc = defaultSaveConfig
	SaveOnClose                      = false

	TagMapDefault = ""
)

var (
	ErrHandlerNotSet = errors.New("handler not set")
)

type MappedConfiguration map[string]interface{}

// Passed *MappedConfiguration
type LoadConfigFunc func(mc *MappedConfiguration) error
type SaveConfigFunc func(mc *MappedConfiguration) error

type RegisterFunc func(parent string, fieldType reflect.StructField, fieldValue reflect.Value)

type funcUpdater struct {
	In  func()
	Out func()
}

type item struct {
	bindEnvArgs bool
	bindConf    bool
	val         interface{}
}

func BindArgs(st interface{}, key string) (err error) {
	mid, in, out := link(st)
	mappedConfUpdater[key] = &funcUpdater{
		In: in, Out: out,
	}
	registered[key] = &item{
		bindEnvArgs: true,
		val:         mid,
	}
	return
}

func BindConf(st interface{}, key string) (err error) {
	mid, in, out := link(st)
	mappedConfUpdater[key] = &funcUpdater{
		In: in, Out: out,
	}
	registered[key] = &item{
		bindConf: true,
		val:      mid,
	}
	return
}

func BindArgsConf(st interface{}, key string) (err error) {
	mid, in, out := link(st)
	mappedConfUpdater[key] = &funcUpdater{
		In: in, Out: out,
	}
	registered[key] = &item{
		bindEnvArgs: true,
		bindConf:    true,
		val:         mid,
	}
	return
}

func Init() (err error) {
	// order priority:
	// 1. defaulf value         - runtime
	// 2. args value            - tag "args"
	// 3. env value             - tag "env"
	// 4. configuration value   - from conf file

	for k, s := range registered {
		if s.bindEnvArgs {
			defer addBindArgs(k, s)
		}
		if s.bindConf {
			addBindConf(k, s)
		}
	}
	var t MappedConfiguration
	defer func() {
		setBackMap(&mappedConf, t)
		loadReupdate()
	}()
	if err = LoadConfig(&t); err != nil {
		return
	}
	return
}

func Save() (err error) {
	if SaveConfig == nil {
		return ErrHandlerNotSet
	}
	saveReupdate()
	if err = SaveConfig(&mappedConf); err != nil {
		err = errors.Wrap(err, "binder.Save")
		return
	}
	return
}

func In() {
	// spew.Dump(mappedConf)
	loadReupdate()
}

func Out() {
	saveReupdate()
}

func Close() (err error) {
	if SaveOnClose {
		if err = Save(); err != nil {
			panic(err)
		}
	}
	return
}

func setBackMap(dst *MappedConfiguration, val MappedConfiguration) {
	ds := reflect.ValueOf(dst)
	if ds.Kind() == reflect.Ptr {
		ds = ds.Elem()
	}
	iface := ds.Interface().(MappedConfiguration)
	for k, orig := range iface {
		src := val[k]
		// spew.Dump(orig)

		if src == nil {
			continue
		}

		iorig := reflect.ValueOf(orig)
		// if iorig.Kind() == reflect.Ptr {
		// 	iorig = iorig.Elem()
		// }
		ival := reflect.ValueOf(src)
		if ival.Kind() == reflect.Ptr {
			ival = ival.Elem()
		}

		// spew.Dump(iorig)
		ifaceToStruct(ival, iorig)
		// spew.Dump(iorig)

		// fmt.Printf(">> SET BACK ival %+#v\n", ival)
		// fmt.Printf(">> SET BACK iorg %+#v\n", iorig)
	}
}

func ifaceToStruct(ival reflect.Value, iorig reflect.Value) {
	if iorig.Kind() != reflect.Struct && (iorig.Kind() != reflect.Ptr && !iorig.IsNil() && iorig.Elem().Kind() == reflect.Struct) {
		conv := ival.Convert(iorig.Type())
		iorig.Set(conv)
		return
	}

	// check if ival is map[string]interface{}
	vk := ival.Type().Key().Kind()
	vv := ival.Type().Elem().Kind()
	if vk != reflect.String && vv != reflect.Interface {
		panic(fmt.Sprintf("convert is not supported : has key %s and val %s", vk, vv))
	}

	// fmt.Printf("\nsdsds %s %s\n\n", vk, vv)
	// to := iorig.Type()
	// for i := 0; i < iorig.NumField(); i++ {

	// }
	// if ival.Kind() != reflect.Map {
	// 	conv := ival.Convert(iorig.Type())
	// 	iorig.Set(conv)
	// 	return
	// }

	// with the help of mapstructure! :D
	if err := mapstructure.Decode(ival.Interface(), iorig.Interface()); err != nil {
		panic(errors.Wrap(err, "mapstructure error"))
	}
	// spew.Dump(iorig.Interface())
}

func addBindArgs(key string, s *item) {
	v := reflect.ValueOf(s.val)
	_ = v
	instFields(key, v, wrapperUnwind(wrapperOSEnv(RegisterCmdArgs)))
}

func addBindConf(key string, s *item) {
	v := reflect.ValueOf(s.val)
	_ = v
	instFields(key, v, wrapperUnwind(wrapperOSEnv(defaultRegisterConf)))
	mappedConf[key] = s.val
}

//
func instFields(parent string, t reflect.Value, eachFieldF RegisterFunc) {
	if eachFieldF == nil {
		return
	}

	if t.Kind() == reflect.Ptr {
		instFields(parent, t.Elem(), eachFieldF)
		return
	}
	tt := t.Type()
	// fmt.Printf("%+#v %+#v\n", t, t.Kind() == reflect.Ptr)
	if prefix := tt.Name(); prefix != "" {
		if parent == "" {
			parent = prefix
		}
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		c := tt.Field(i)
		// fmt.Printf("got val %s %+#v\n", f.Kind(), f.Interface())

		k := f
		for k.Kind() == reflect.Ptr {
			if t := k.Elem(); t.Kind() != reflect.Ptr {
				// fmt.Printf("?> %+#v\n", t.Interface())
				if t.Kind() == reflect.Struct {
					instSub(parent, c, f, eachFieldF)
				}
				break
			} else {
				if t.IsNil() {
					nt := reflect.New(t.Type().Elem())
					t.Set(nt)
				}
				k = t
				// fmt.Printf(">> %+#v\n", k.Interface())
			}
		}

		// if f.Kind() == reflect.Struct {
		// 	instSub(parent, c, f, eachFieldF)
		// 	continue
		// }

		// if f.Kind() == reflect.Ptr {
		// 	newVal := f
		// 	if f.IsNil() {
		// 		newVal = reflect.New(f.Type().Elem())
		// 		f.Set(newVal)
		// 	}

		// 	var n reflect.Value
		// 	unwind := f.Type().Elem()
		// 	// fmt.Printf("uw %s\n", unwind)

		// 	// unwind = unwind.Elem()
		// 	for unwind.Kind() == reflect.Ptr && newVal.Elem().IsNil() {
		// 		// fmt.Printf("unwind %s\n", unwind)
		// 		unwind = unwind.Elem()
		// 		// fmt.Printf("unwind el %s\n", unwind)
		// 		n = reflect.New(unwind)
		// 		newVal.Elem().Set(n)
		// 		newVal = n
		// 	}

		// 	if f.Type().Elem().Kind() == reflect.Struct {
		// 		instSub(parent, c, f, eachFieldF)
		// 		continue
		// 	}
		// }
		eachFieldF(parent, c, f)
	}
}

func instSub(parent string, sf reflect.StructField, v reflect.Value, fc RegisterFunc) {
	if v.Kind() == reflect.Ptr {
		// initialize nil pointer
		next := v.Elem()
		if next.Kind() == reflect.Invalid { // nil
			// v.SetPointer(unsafe.Pointer(reflect.New(v.Type()).Pointer()))
			// val := reflect.New(sf.Type.Elem())
			val := reflect.New(v.Type().Elem())
			// fmt.Printf("%s [%s] %+#v %s\n", parent, sf.Name, val, sf.Type)

			v.Set(val)
			next = v.Elem()
		}
		instSub(parent, sf, next, fc)
		return
	}
	// fmt.Printf("%s %s %+#v\n", v.Kind(), sf.Type, v.Kind() == reflect.Ptr)
	if v.Kind() != reflect.Struct {
		fc(parent, sf, v)
		return
	}
	tt := v.Type()
	bindName := sf.Tag.Get("bind")
	bindNameParsed := SplitTagValue(bindName)
	if len(bindNameParsed) > 0 && bindNameParsed[0] != "" {
		bindName = bindNameParsed[0]
		if parent != "" {
			parent = parent + "." + bindName
		} else {
			parent = bindName
		}
	} else if prefix := v.Type().Name(); prefix != "" {
		if parent != "" {
			parent = parent + "." + prefix
		} else {
			parent = prefix
		}
	} else if prefix = sf.Name; prefix != "" {
		if parent != "" {
			parent = parent + "." + prefix
		} else {
			parent = prefix
		}
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		c := tt.Field(i)

		k := f
		for k.Kind() == reflect.Ptr {
			if t := k.Elem(); t.Kind() != reflect.Ptr {
				break
			} else {
				if t.IsNil() {
					nt := reflect.New(t.Type().Elem())
					t.Set(nt)
				}
				k = t
			}
		}

		// fmt.Printf("fc %+#v %+#v %+#v\n", f.Interface(), f.Elem(), k.Interface())
		fc(parent, c, f)
	}
}

func loadReupdate() {
	for _, v := range mappedConfUpdater {
		v.In()
	}
}

func saveReupdate() {
	for _, v := range mappedConfUpdater {
		v.Out()
	}
}

func wrapperUnwind(f RegisterFunc) RegisterFunc {
	return func(parent string, fieldType reflect.StructField, fieldValue reflect.Value) {
		defer f(parent, fieldType, fieldValue)

		for fieldValue.Kind() == reflect.Ptr {
			if t := fieldValue.Elem(); t.Kind() != reflect.Ptr {
				break
			} else {
				fieldValue = t
			}
		}

	}
}

func wrapperOSEnv(f RegisterFunc) RegisterFunc {
	return func(parent string, fieldType reflect.StructField, fieldValue reflect.Value) {
		defer f(parent, fieldType, fieldValue)

		environName := fieldType.Tag.Get("environ")
		envName := fieldType.Tag.Get("env")

		environNameParsed := SplitTagValue(environName)
		envNameParsed := SplitTagValue(envName)

		if len(envNameParsed) > 0 {
			envName = envNameParsed[0]
		}

		if len(environNameParsed) > 0 && environNameParsed[0] != "" { // dedicated name
			environName = environNameParsed[0]
			envName = environName
		} else if envName == "" {
			// don't lookup any
			return
		} else {
			if parent != "" {
				envName = parent + "." + envName
			}
		}

		// fmt.Printf("checking env %+#v\n", envName)
		val, ok := os.LookupEnv(envName)
		if ok {
			// fmt.Printf("header exists!! %+#v\n", envName)
			fieldValue.Elem().Set(reflect.ValueOf(val))
		}
	}
}
