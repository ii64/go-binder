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
	if v, ok := a[TagName]; !ok || v == "" {
		a[TagName] = f.Name
		if overrideTagValue, ok := a[TagName]; ok && overrideTagValue != "" {
			a[TagName] = overrideTagValue
		}
	}
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
	prevKey := []string{}
	for tag != "" {
		prev = 0
		i = 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' {
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
		if tag != "" && tag[0] == ':' {
			prev = 1
			i = 1
			del := tag[i]
			if del == '\'' || del == '`' {
				panic("unexpected quote found")
			}
			if del != '"' {
				del = ' '
			} else {
				i++
			}
			for i < len(tag) && tag[i] != del {
				if tag[i] == '\\' {
					i++
				}
				i++
			}
			if i >= len(tag) {
				qvalue = tag[prev:]
			} else if del != ' ' {
				qvalue = string(tag[prev : i+1])
				tag = tag[i+1:]
				var err error
				if qvalue, err = strconv.Unquote(qvalue); err != nil {
					break
				}
			} else {
				qvalue = string(tag[prev:i])
				tag = tag[i:]
			}
			prevKey = append(prevKey, tmp)
			for _, pk := range prevKey {
				if _, ok := value[pk]; !ok {
					value[pk] = qvalue
				}
			}
			prevKey = []string{}
		}
		prevKey = append(prevKey, tmp)
		if len(tag) < 1 {
			for _, pk := range prevKey {
				if _, ok := value[pk]; !ok {
					value[pk] = qvalue
				}
			}
			prevKey = []string{}
		}
	}
	return
}
