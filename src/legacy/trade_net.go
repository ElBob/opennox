package legacy

/*
#include <string.h>
#include "defs.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "server__system__trade.h"

// The original code passes the item pointer as a float.
static int nox_xxx_tradeP2PAddOffer2_50F820_trade_obj(int a1, int a2, int item) {
	float f;
	memcpy(&f, &item, sizeof(f));
	return nox_xxx_tradeP2PAddOffer2_50F820_trade(a1, a2, f);
}

// The original code passes the item count as a pointer.
static float* sub_510640_trade_cnt(int a1, int a2, int a3, int cnt) {
	return sub_510640_trade(a1, a2, a3, (float*)cnt);
}
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func sdecodeObj(obj *server.Object) C.int {
	return C.int(uintptr(unsafe.Pointer(asObjectC(obj))))
}

func sdecodeTrade(t *server.TradeSession) C.int {
	return C.int(uintptr(unsafe.Pointer(t)))
}

func sdecodeWStr(data []byte) (*wchar2_t, func()) {
	buf, free := alloc.Make([]byte{}, len(data)+2)
	copy(buf, data)
	return (*wchar2_t)(unsafe.Pointer(&buf[0])), free
}

func Sub_506870(a1 int, u *server.Object, text []byte) {
	if text == nil {
		C.sub_506870(C.int(a1), sdecodeObj(u), nil)
		return
	}
	p, free := sdecodeWStr(text)
	defer free()
	C.sub_506870(C.int(a1), sdecodeObj(u), p)
}

func Sub_506C90(a1 int, u *server.Object, text []byte) {
	if text == nil {
		C.sub_506C90(C.int(a1), sdecodeObj(u), nil)
		return
	}
	p, free := sdecodeWStr(text)
	defer free()
	C.sub_506C90(C.int(a1), sdecodeObj(u), p)
}

func Sub_50F3A0(t *server.TradeSession) {
	C.sub_50F3A0((*C.uint32_t)(unsafe.Pointer(t)))
}

func Nox_xxx_tradeP2PAddOffer2_50F820_trade(t *server.TradeSession, u, item *server.Object) bool {
	return C.nox_xxx_tradeP2PAddOffer2_50F820_trade_obj(sdecodeTrade(t), sdecodeObj(u), sdecodeObj(item)) == 1
}

func Nox_xxx_tradeP2PAddOfferMB_50FE20(t *server.TradeSession, code int) {
	C.nox_xxx_tradeP2PAddOfferMB_50FE20(sdecodeTrade(t), C.int(code))
}

func Nox_xxx_tradeAccept_50F5A0(t *server.TradeSession, u *server.Object) {
	C.nox_xxx_tradeAccept_50F5A0(sdecodeTrade(t), sdecodeObj(u))
}

func Nox_xxx_shopExit_50F4C0(t *server.TradeSession) {
	C.nox_xxx_shopExit_50F4C0((*C.uint32_t)(unsafe.Pointer(t)))
}

func Nox_xxx_servShopStart_50EF10_trade(u, shop *server.Object) {
	C.nox_xxx_servShopStart_50EF10_trade(sdecodeObj(u), sdecodeObj(shop))
}

func Sub_5100C0_trade(u *server.Object, t *server.TradeSession, code int) {
	C.sub_5100C0_trade(sdecodeObj(u), (*C.uint32_t)(unsafe.Pointer(t)), C.int(code))
}

func Sub_510640_trade(u *server.Object, t *server.TradeSession, typ int, cnt int) {
	C.sub_510640_trade_cnt(sdecodeObj(u), sdecodeTrade(t), C.int(typ), C.int(cnt))
}

func Sub_510BE0_trade(u *server.Object, t *server.TradeSession, code int) {
	C.sub_510BE0_trade((*C.int)(unsafe.Pointer(asObjectC(u))), sdecodeTrade(t), (*C.uint32_t)(unsafe.Pointer(uintptr(code))))
}

func Sub_510D10(u *server.Object, t *server.TradeSession, code int, cnt int) {
	C.sub_510D10((*C.int)(unsafe.Pointer(asObjectC(u))), sdecodeTrade(t), C.int(code), C.uint(cnt))
}

func Sub_510AE0(u *server.Object, t *server.TradeSession, code int) {
	C.sub_510AE0((*C.int)(unsafe.Pointer(asObjectC(u))), sdecodeTrade(t), (*C.uint32_t)(unsafe.Pointer(uintptr(code))))
}

func Sub_5109C0_trade(u *server.Object, t *server.TradeSession, code int) {
	C.sub_5109C0_trade((*C.int)(unsafe.Pointer(asObjectC(u))), sdecodeTrade(t), (*C.uint32_t)(unsafe.Pointer(uintptr(code))))
}

func Sub_5108D0(u *server.Object, t *server.TradeSession, code int) {
	C.sub_5108D0(sdecodeObj(u), sdecodeTrade(t), C.int(code))
}
