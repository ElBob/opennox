package legacy

import "C"
import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

//export sub_5007E0
func sub_5007E0(pattern *C.char) {
	GetServer().S().QuestStatus.Reset(GoString(pattern))
}

//export sub_500A60
func sub_500A60() C.int {
	if err := GetServer().S().QuestStatus.Save(cryptfile.Global(), noxflags.HasGame(noxflags.GameModeCoop)); err != nil {
		return 0
	}
	return 1
}

//export sub_500B70
func sub_500B70() C.int {
	if err := GetServer().S().QuestStatus.Load(cryptfile.Global()); err != nil {
		return 0
	}
	return 1
}
