// +build go1.12

package breeze

import "reflect"

func rangeMap(buf *Buffer, v reflect.Value) (err error) {
	iter := v.MapRange()
	for iter.Next() {
		err = writeReflectValue(buf, iter.Key(), true)
		if err != nil {
			return err
		}
		err = writeReflectValue(buf, iter.Value(), true)
		if err != nil {
			return err
		}
	}
	return nil
}

func rangePackedMap(buf *Buffer, v reflect.Value) {
	iter := v.MapRange()
	var err error
	first := true
	for iter.Next() {
		if first {
			writeType(buf, iter.Key())
			writeType(buf, iter.Value())
			first = false
		}
		err = writeReflectValue(buf, iter.Key(), false)
		if err != nil {
			panic(err)
		}
		err = writeReflectValue(buf, iter.Value(), false)
		if err != nil {
			panic(err)
		}
	}
}
