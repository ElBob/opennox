package legacy

import "C"

//export nox_xxx_gamedataGetFloat_419D40
func nox_xxx_gamedataGetFloat_419D40(k *C.char) C.double {
	key := GoString(k)
	val := C.double(GetServer().S().Balance.Float(key))
	return val
}

//export nox_xxx_gamedataGetFloatTable_419D70
func nox_xxx_gamedataGetFloatTable_419D70(k *C.char, i int) C.double {
	key := GoString(k)
	val := C.double(GetServer().S().Balance.FloatInd(key, i))
	return val
}

