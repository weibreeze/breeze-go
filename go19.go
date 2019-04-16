// +build go1.9, !go1.12

package breeze

import (
	"reflect"
)

func rangeMap(buf *Buffer, v reflect.Value) (err error) {
	ks := v.MapKeys()
	for _, k := range ks {
		err = WriteValue(buf, k)
		if err != nil {
			return err
		}
		err = WriteValue(buf, v.MapIndex(k))
		if err != nil {
			return err
		}
	}
	return err
}
