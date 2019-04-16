// +build go1.12

package breeze

import "reflect"

func rangeMap(buf *Buffer, v reflect.Value) (err error) {
	iter := v.MapRange()
	for iter.Next() {
		err = WriteValue(buf, iter.Key())
		if err != nil {
			return err
		}
		err = WriteValue(buf, iter.Value())
		if err != nil {
			return err
		}
	}
	return nil
}
