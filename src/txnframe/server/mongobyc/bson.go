package mongobyc

// #include <stdlib.h>
// #include "mongo.h"
import "C"

import "unsafe"

type bson struct {
	data string
}

func (b *bson) ToDocument() (*C.bson_t, error) {

	var err C.bson_error_t
	innerData := C.CString(b.data)
	defer C.free(unsafe.Pointer(innerData))
	bson := C.bson_new_from_json((*C.uchar)(unsafe.Pointer(innerData)), -1, &err)
	if nil == bson {
		return nil, TransformError(err)
	}
	return bson, nil
}
