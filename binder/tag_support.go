package binder

import (
	"reflect"
	"strconv"
	"strings"
)

func SplitTagValue(value string) []string {
	return strings.Split(value, ",")
}

func reTag(f reflect.StructField, tags *reflect.StructTag) {
	a := ParseTag(string(*tags))
	if v, ok := a["mapstructure"]; !ok || v == "" {
		a["mapstructure"] = f.Name
		if overrideTagValue, ok := a[TagMapDefault]; ok && overrideTagValue != "" {
			a["mapstructure"] = overrideTagValue
		}
	}
	// fmt.Printf("%+#v %+#v\n", TagMapDefault, a)
	*tags = JoinTag(a)
}

func JoinTag(tag map[string]string) reflect.StructTag {
	i := 0
	var r string
	for k, v := range tag {
		r = r + k + ":" + strconv.Quote(v)
		if i+1 < len(tag) {
			r = r + " "
		}
		i++
	}
	return reflect.StructTag(r)
}
func ParseTag(tag string) (value map[string]string) {
	value = map[string]string{}
	i := 0
	prev := 0
	for tag != "" {
		prev = 0
		i = 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && (i+1 >= len(tag) || tag[i+1] != '"') {
			i++
		}
		for prev < len(tag) && tag[prev] == ' ' {
			prev++
		}
		tmp := tag[prev:i]
		qvalue := ""
		tag = tag[i:]
		prev = 0
		i = 0
		if len(tag) != 0 && tag[0] == ':' {
			prev = 1
			i = 2
			for i < len(tag) && tag[i] != '"' {
				if tag[i] == '\\' {
					i++
				}
				i++
			}
			if i >= len(tag) {
				qvalue = tag[prev:]
			} else {
				qvalue = string(tag[prev : i+1])
				tag = tag[i+1:]
			}

			var err error
			if qvalue, err = strconv.Unquote(qvalue); err != nil {
				break
			}
		}
		value[tmp] = qvalue
	}
	return
}
