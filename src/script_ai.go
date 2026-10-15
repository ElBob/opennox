package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func (s *Server) scriptMoveTo(obj *server.Object, wp *server.Waypoint) { // nox_server_scriptMoveTo_5123C0
	if obj.Flags().Has(object.FlagDead) {
		return
	}
	if obj.Class().Has(object.ClassMonster) {
		ud := obj.UpdateDataMonster()
		obj.ClearActionStack()
		obj.MonsterPushAction(ai.ACTION_REPORT, 8)
		if wp.PointsCnt != 0 {
			obj.MonsterPushAction(ai.ACTION_ROAM, wp.C(), 0, uint32(uint8(ud.Field333)))
		}
		obj.MonsterPushAction(ai.ACTION_FAR_MOVE_TO, wp.PosVec, 0)
		return
	}
	moverTyp := legacy.Get_dword_5d4594_2386836()
	if int(obj.TypeInd) == moverTyp {
		s.moverGoTo(obj, wp)
		return
	}
	for it := s.Objs.First(); it != nil; it = it.Next() {
		if int(it.TypeInd) == moverTyp && it.UpdateDataMover().Field_8 == obj.Extent {
			s.moverGoTo(it, wp)
		}
	}
}

func (s *Server) moverGoTo(obj *server.Object, wp *server.Waypoint) { // nox_xxx_moverGoTo_5124C0
	ud := obj.UpdateDataMover()
	legacy.Nox_xxx_objectSetOn_4E75B0(obj)
	obj.VelVec = types.Pointf{}
	ud.Field_0 = 0
	ud.Field_2 = int32(wp.Index)
	s.Objs.AddToUpdatable(obj)
}

func scriptMonsterRoam(obj *server.Object) { // nox_xxx_scriptMonsterRoam_512930
	if !obj.Class().Has(object.ClassMonster) || obj.Flags().Has(object.FlagDead) {
		return
	}
	ud := obj.UpdateDataMonster()
	obj.ClearActionStack()
	obj.MonsterPushAction(ai.ACTION_REPORT, 10)
	obj.MonsterPushAction(ai.ACTION_ROAM, 0, 0, uint32(uint8(ud.Field333)))
}

func scriptMonsterGoHome(obj *server.Object) { // nox_server_gotoHome
	if !obj.Class().Has(object.ClassMonster) || obj.Flags().Has(object.FlagDead) {
		return
	}
	ud := obj.UpdateDataMonster()
	obj.ClearActionStack()
	obj.MonsterPushAction(ai.ACTION_REPORT, 37)
	cos, sin := server.SinCosDir(byte(ud.Direction94))
	// The original code computes it with double precision.
	pos := types.Pointf{
		X: float32(float64(cos)*10.0 + float64(obj.PosVec.X)),
		Y: float32(float64(sin)*10.0 + float64(obj.PosVec.Y)),
	}
	obj.MonsterPushAction(ai.ACTION_FACE_LOCATION, pos)
	obj.MonsterPushAction(ai.ACTION_MOVE_TO_HOME, ud.Pos95, 0)
}
