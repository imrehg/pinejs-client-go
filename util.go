package pinejs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"

	"github.com/bitly/go-simplejson"
)

// oDataEncodeVals URL Encode values and separate with commas.
func oDataEncodeVals(strs []string) string {
	encoded := make([]string, len(strs))

	for i, str := range strs {
		encoded[i] = url.QueryEscape(str)
	}

	return strings.Join(encoded, ",")
}

// encodeQuery encodes query values, working around a net/url issue whereby keys
// get encoded as well as values. We only want values encoded, otherwise OData
// dies.
func encodeQuery(query map[string][]string) string {
	if len(query) == 0 {
		return ""
	}

	var tuples []string
	for key, vals := range query {
		if len(vals) == 0 {
			continue
		}

		tuple := key + "=" + oDataEncodeVals(vals)
		tuples = append(tuples, tuple)
	}

	return "?" + strings.Join(tuples, "&")
}

// ptrType extracts the pointer type of a provided interface, e.g. *Foo -> Foo.
func ptrType(v interface{}) (reflect.Type, error) {
	if ty := reflect.TypeOf(v); ty.Kind() != reflect.Ptr {
		return nil, errors.New("not a pointer")
	} else {
		return ty.Elem(), nil
	}
}

func assertPointerToStruct(v interface{}) error {
	if el, err := ptrType(v); err != nil {
		return err
	} else if el.Kind() != reflect.Struct {
		return errors.New("not a pointer to a struct")
	}

	return nil
}

func assertPointerToSliceStructs(v interface{}) error {
	if el, err := ptrType(v); err != nil {
		return err
	} else if el.Kind() != reflect.Slice {
		return errors.New("not a pointer to a slice")
	} else if el.Elem().Kind() != reflect.Struct {
		return errors.New("not a pointer to a slice of structs")
	}

	return nil
}

func toJsonReader(v interface{}) (io.Reader, error) {
	if v == nil {
		return nil, nil
	}

	if buf, err := json.Marshal(v); err != nil {
		return nil, err
	} else {
		return bytes.NewReader(buf), nil
	}
}

// Some functionality that is strangely lacking from simplejson...

type jsonNodeType int

const (
	jsonObject jsonNodeType = iota
	jsonArray
	jsonValue // Anything else.
)

func getJsonNodeType(j *simplejson.Json) jsonNodeType {
	// TODO: Reuse returned values.
	if _, err := j.Map(); err == nil {
		return jsonObject
	} else if _, err := j.Array(); err == nil {
		return jsonArray
	} else {
		return jsonValue
	}
}

func getJsonFieldNames(j *simplejson.Json) (ret []string) {
	if obj, err := j.Map(); err != nil {
		// Caller should have checked.
		panic(err)
	} else {
		for name, _ := range obj {
			ret = append(ret, name)
		}
	}

	return
}

func getJsonFields(j *simplejson.Json) map[string]*simplejson.Json {
	ret := make(map[string]*simplejson.Json)

	for _, name := range getJsonFieldNames(j) {
		ret[name] = j.Get(name)
	}

	return ret
}

func getJsonArray(j *simplejson.Json) (ret []*simplejson.Json) {
	if arr, err := j.Array(); err != nil {
		// Caller should have checked.
		panic(err)
	} else {
		// TODO: This sucks. Don't want to remarshal just to use returned data
		// though.
		for i := 0; i < len(arr); i++ {
			ret = append(ret, j.GetIndex(i))
		}
	}

	return
}
