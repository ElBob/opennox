package legacy

/*
#include "defs.h"
#include "GAME2.h"

extern uint32_t dword_5d4594_832484;
extern uint32_t dword_5d4594_832492;
extern uint32_t dword_5d4594_832496;
extern uint32_t dword_5d4594_832500;
extern uint32_t dword_5d4594_832504;
extern uint32_t dword_5d4594_832508;
extern uint32_t dword_5d4594_832512;
extern uint32_t dword_5d4594_832516;
extern uint32_t dword_5d4594_832520;
extern uint32_t dword_5d4594_832524;
extern uint32_t dword_5d4594_832528;
extern uint32_t dword_5d4594_832532;
extern uint32_t dword_5d4594_832536;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client"
)

// Sub_44F300 is the Go implementation of sub_44F300 (quest briefing screen drawing); set by the main package.
var Sub_44F300 func(a1, a2 int) int

//export sub_44F300
func sub_44F300(a1, a2 int) int { return Sub_44F300(a1, a2) }

func Sub_44E110() {
	C.sub_44E110()
}

// Get_dword_5d4594_832484 returns the font pointer (the C global holds a 32-bit pointer).
func Get_dword_5d4594_832484() unsafe.Pointer {
	return unsafe.Pointer(uintptr(C.dword_5d4594_832484))
}

func briefDrawable(v C.uint) *client.Drawable {
	return (*client.Drawable)(unsafe.Pointer(uintptr(v)))
}

func Get_dword_5d4594_832492() *client.Drawable { return briefDrawable(C.dword_5d4594_832492) }
func Get_dword_5d4594_832496() *client.Drawable { return briefDrawable(C.dword_5d4594_832496) }
func Get_dword_5d4594_832500() *client.Drawable { return briefDrawable(C.dword_5d4594_832500) }
func Get_dword_5d4594_832504() *client.Drawable { return briefDrawable(C.dword_5d4594_832504) }
func Get_dword_5d4594_832508() *client.Drawable { return briefDrawable(C.dword_5d4594_832508) }
func Get_dword_5d4594_832512() *client.Drawable { return briefDrawable(C.dword_5d4594_832512) }
func Get_dword_5d4594_832516() *client.Drawable { return briefDrawable(C.dword_5d4594_832516) }
func Get_dword_5d4594_832520() *client.Drawable { return briefDrawable(C.dword_5d4594_832520) }
func Get_dword_5d4594_832524() *client.Drawable { return briefDrawable(C.dword_5d4594_832524) }
func Get_dword_5d4594_832528() *client.Drawable { return briefDrawable(C.dword_5d4594_832528) }
func Get_dword_5d4594_832532() *client.Drawable { return briefDrawable(C.dword_5d4594_832532) }
func Get_dword_5d4594_832536() *client.Drawable { return briefDrawable(C.dword_5d4594_832536) }
