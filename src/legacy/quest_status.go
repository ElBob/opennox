package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_1570272;
*/
import "C"
import (
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/server"
)

func init() {
	server.SetQuestStatusHead(unsafe.Pointer(&C.dword_5d4594_1570272))
}

//export sub_5007E0
func sub_5007E0(pattern *C.char) {
	GetServer().S().QuestStatus.Reset(GoString(pattern))
}

//export sub_500A60
func sub_500A60() C.int {
	return C.int(GetServer().S().QuestStatus.Save(cryptfile.Global(), noxflags.HasGame(noxflags.GameModeCoop)))
}

//export sub_500B70
func sub_500B70() C.int {
	return C.int(GetServer().S().QuestStatus.Load(cryptfile.Global()))
}
