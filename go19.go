// +build go1.9, !go1.12

package breeze

import (
	"reflect"
)

func rangeMap(buf *Buffer, v reflect.Value) (err error) {
	ks := v.MapKeys()
	for _, k := range ks {
		err = writeReflectValue(buf, k, true)
		if err != nil {
			return err
		}
		err = writeReflectValue(buf, v.MapIndex(k), true)
		if err != nil {
			return err
		}
	}
	return err
}

func rangePackedMap(buf *Buffer, v reflect.Value) {
	ks := v.MapKeys()
	var err error
	first := true
	for _, k := range ks {
		if first {
			writeType(buf, k)
			writeType(buf, v.MapIndex(k))
			first = false
		}
		err = writeReflectValue(buf, k, false)
		if err != nil {
			panic(err)
		}
		err = writeReflectValue(buf, v.MapIndex(k), false)
		if err != nil {
			panic(err)
		}
	}
}
