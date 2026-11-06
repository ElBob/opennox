package opennox

import (
	"image"
	"unsafe"

	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// drawSetColor is nox_client_drawSetColor_434460.
func (c *Client) drawSetColor(v uint32) {
	c.r.Data().SetColor2(noxcolor.RGBA5551(v))
}

// teamColor32 is nox_xxx_materialGetTeamColor_418D50.
func (c *Client) teamColor32(t *server.Team) uint32 {
	return noxcolor.ToRGBA5551Color(c.srv.Teams.GetTeamColor(t)).Color32()
}

// minimapPolyByte2 is BYTE2(v->field_0[32]) on a polygon check result.
func minimapPolyByte2(p *legacy.Nox_player_polygon_check_data) int {
	return int(p.Field_0[32]>>16) & 0xff
}

func (c *Client) nox_xxx_drawMinimapAndLines_4738E0() int {
	result := 1
	if legacy.Get_nox_client_gui_flag_1556112() != 1 {
		if memmap.Uint8(0x5D4594, 1096424)&1 != 0 {
			v1 := c.Objs.ByNetCodeDynamic(legacy.ClientPlayerNetCode())
			c.nox_xxx_drawMinimap4Sprite_4725C0(v1)
		}
		result = c.nox_xxx_drawMessageLines_445530()
	}
	return result
}

// nox_xxx_drawMessageLines_445530 draws the last three chat/status message lines.
func (c *Client) nox_xxx_drawMessageLines_445530() int {
	v13 := 0
	vp := c.Viewport() // nox_draw_getViewport_437250
	width := legacy.Get_nox_win_width()
	v1 := 3*vp.Size.Y/4 + vp.Screen.Min.Y - 15
	v2 := int(legacy.Get_dword_5d4594_825736())
	i := v2
	result := 0
	fnt := c.r.GetFonts().AsFont(nil)
	for ; ; v2 = i {
		v3 := 161 * v2
		result = int(c.srv.Frame())
		v5 := uintptr(4 * v3)
		if memmap.Uint32(0x5D4594, 824440+v5) < c.srv.Frame() {
			break
		}
		c.r.Data().SetTextColor(noxcolor.RGBA5551(legacy.Get_nox_color_black_2650656()))
		v6 := memmap.PtrOff(0x5D4594, 823804+v5)
		str := legacy.GoWStringP(v6)
		v14 := c.r.GetStringSizeWrapped(fnt, str, 0).X // nox_xxx_drawGetStringSize_43F840
		v7 := width - v14
		v9 := (width - v14) / 2
		for off := uintptr(107848); off < 107880; off += 8 {
			c.r.DrawString(fnt, str, image.Pt(v9+int(memmap.Int32(0x587000, off)), v1+int(memmap.Int32(0x587000, off+4))))
		}
		v10 := v13
		if v13 != 0 {
			c.r.Data().SetTextColor(noxcolor.RGBA5551(memmap.Uint32(0x5D4594, 2597996)))
		} else {
			c.r.Data().SetTextColor(noxcolor.RGBA5551(legacy.Get_nox_color_white_2523948()))
		}
		c.r.DrawString(fnt, str, image.Pt(v7/2, v1))
		v11 := -4 - c.r.FontHeight(fnt)
		result = i
		v1 += v11
		if i != 0 {
			i--
		} else {
			i = 2
		}
		v13++
		if v10+1 >= 3 {
			break
		}
	}
	return result
}

// nox_xxx_polygonGetIdxA_421790 returns the index of a polygon containing the point, or 0.
func (c *Client) nox_xxx_polygonGetIdxA_421790(a1 *[2]int32, a2 int) int {
	v2 := legacy.Nox_xxx_polygonIsPlayerInPolygon_4217B0(unsafe.Pointer(a1), a2)
	if v2 != nil {
		return int(v2.Field_0[20])
	}
	return 0
}

func (c *Client) sub_472540(a1 *client.Drawable) int {
	var a1a [2]int32
	if a1 == c.ClientPlayerUnit() { // *(0x852978, 8)
		p := c.Viewport().World.Max // nox_xxx_getSomeCoods_435670
		a1a[0] = int32(p.X)
		a1a[1] = int32(p.Y)
	} else {
		a1a[0] = int32(a1.PosVec.X)
		a1a[1] = int32(a1.PosVec.Y)
	}
	v2 := c.nox_xxx_polygonGetIdxA_421790(&a1a, int(memmap.Int32(0x5D4594, 1096312)))
	if v2 != 0 {
		*memmap.PtrUint32(0x5D4594, 1096312) = uint32(v2)
	} else {
		v2 = int(memmap.Int32(0x5D4594, 1096312))
	}
	if v2 != 0 {
		return int(*(*uint8)(unsafe.Add(legacy.Nox_xxx_polygonGetByIdx_4214A0(v2), 130)))
	}
	return 1
}

func (c *Client) nox_xxx_drawMinimap4Sprite_4725C0(a1 *client.Drawable) {
	if !c.ClientPlayerUnit().HasEnchant(2) { // nox_client_drawable_testBuff_4356C0(*(0x852978, 8), 2)
		sub_437260()
		*memmap.PtrUint32(0x5D4594, 1096316) = uint32(c.sub_472540(a1))
		c.nox_xxx_cliDrawMinimap_472600(a1, int(memmap.Int32(0x5D4594, 1096316)))
		c.r.SetRectFullScreen() // sub_437290
	}
}

// minimapFirstMonster is nox_xxx_minimapFirstMonster_50AAE0; returns the monster position pointer.
func (c *Client) minimapFirstMonster() *types.Pointf {
	v0 := c.srv.Objs.First()
	legacy.Set_dword_5d4594_1599696(v0)
	if v0 == nil {
		return nil
	}
	for uint32(v0.ObjClass)&2 == 0 {
		v0 = v0.Next()
		legacy.Set_dword_5d4594_1599696(v0)
		if v0 == nil {
			return nil
		}
	}
	return &v0.PosVec
}

// minimapNextMonster is nox_xxx_minimapNextMonster_50AB10.
func (c *Client) minimapNextMonster() *types.Pointf {
	v0 := legacy.Get_dword_5d4594_1599696().Next()
	legacy.Set_dword_5d4594_1599696(v0)
	if v0 == nil {
		return nil
	}
	for uint32(v0.ObjClass)&2 == 0 {
		v0 = v0.Next()
		legacy.Set_dword_5d4594_1599696(v0)
		if v0 == nil {
			return nil
		}
	}
	return &v0.PosVec
}

func (c *Client) nox_xxx_cliDrawMinimap_472600(a1 *client.Drawable, a2 int) {
	_ = a1 // unused in the original
	vp := c.Viewport()
	winW := legacy.Get_nox_win_width()
	winH := legacy.Get_nox_win_height()
	if memmap.Uint8(0x5D4594, 1096300) == 0 {
		*memmap.PtrUint8(0x5D4594, 1096300) = byte(c.srv.Walls.DefIndByName("InvisibleWallSet"))
		*memmap.PtrUint8(0x5D4594, 1096301) = byte(c.srv.Walls.DefIndByName("InvisibleBlockingWallSet"))
	}
	c.r.Data().SetAlphaEnabled(false)
	legacy.Nox_xxx_wndDraw_49F7F0()
	v3 := winW / 6
	v4 := winH - winW/6
	yTop := v4 / 2
	c.r.setClipRects(image.Pt(0, 0), image.Pt(winW, winH)) // nox_client_copyRect_49F6F0
	v5 := vp.Screen.Min.X                                  // *(uint32_t*)v2
	// The original compares the unsigned value: "*(uint32_t*)v2 <= 0" is true only for 0.
	if uint32(v5) == 0 {
		c.r.DrawRectFilledAlpha(0, v4/2, v3, v3)
	} else {
		c.drawSetColor(legacy.Get_nox_color_black_2650656())
		if v5 >= v3 {
			c.r.DrawRectFilledOpaque(0, v4/2, v3, v3, c.r.Data().Color2())
		} else {
			c.r.DrawRectFilledOpaque(0, v4/2, v5, v3, c.r.Data().Color2())
			c.r.DrawRectFilledAlpha(v5, v4/2, v3-v5, v3)
		}
	}
	c.r.Data().SetAlphaEnabled(true)
	c.drawSetColor(legacy.Get_nox_color_black_2650656())
	c.r.Data().SetAlpha(0x5A)
	c.nox_client_drawRectLines_473510(-1, yTop-1, v3+2, v3+2)
	c.r.Data().SetAlpha(0x3C)
	c.nox_client_drawRectLines_473510(-2, yTop-2, v3+4, v3+4)
	c.r.Data().SetAlpha(0x28)
	c.nox_client_drawRectLines_473510(-3, yTop-3, v3+6, v3+6)
	c.r.Data().SetAlphaEnabled(false)
	c.r.setClipRects(image.Pt(0, yTop), image.Pt(v3, v3))
	v6 := int(legacy.Get_nox_xxx_minimap_587000_149232())
	v7 := v3 * v6 / 100
	v84 := c.Viewport().World.Max // nox_xxx_getSomeCoods_435670
	v8 := v84.X - v7/2
	v9 := v84.Y - v7/2
	var xLeft image.Point
	xLeft.X = v84.X - v7/2
	xLeft.Y = v9
	v81 := v8 / 23
	v77 := (v8 + v7) / 23
	v78 := (v9 + v7) / 23
	v74 := 23 * (v8 / 23)
	v80 := 23 * (v9 / 23)
	v10 := v9 / 23
	for i := v9 / 23; i <= v78; i++ {
		v11 := c.srv.Walls.IndexByY(v10) // sub_4106A0
		for v11 != nil {
			// The original stores the tile as (signed) char, so tiles >= 128 never match.
			v12 := int8(v11.Tile1)
			if int(v12) == int(memmap.Uint8(0x5D4594, 1096300)) {
				goto LABEL_37
			}
			if int(v12) == int(memmap.Uint8(0x5D4594, 1096301)) {
				goto LABEL_37
			}
			if uint8(v11.Field8) != 0 && int(uint8(v11.Field8)) != a2 {
				goto LABEL_37
			}
			{
				v13 := int(v11.X5)
				if v13 < v81 {
					goto LABEL_37
				}
				if v13 > v77 {
					break
				}
				v14 := v74 + 23*(v13-v81)
				if uint8(v11.Flags4)&0x10 != 0 {
					// Field32 holds a pointer to a drawable here.
					v15 := legacy.AsDrawableP(unsafe.Pointer(uintptr(v11.Field32)))
					if v15 == nil {
						goto LABEL_37
					}
					v69 := 0
					for off := uintptr(149244); ; off += 8 {
						dx := int(memmap.Int32(0x587000, off-4))
						dy := int(memmap.Int32(0x587000, off))
						v17 := c.srv.Walls.GetWallAtGrid(image.Pt(dx+int(v11.X5), dy+int(v11.Y6)))
						if v17 != nil {
							if v17.Field12 != 0 {
								if v69 < 4 {
									c.minimapDoorLine(v15, xLeft, yTop)
									v8 = xLeft.X
								}
								goto LABEL_37
							}
						} else {
							v18 := c.srv.Walls.GetWallAtGrid(image.Pt(dx+int(v11.X5), int(v11.Y6)))
							var v19 *server.Wall
							if (v18 != nil && v18.Field12 != 0) ||
								func() bool {
									v19 = c.srv.Walls.GetWallAtGrid(image.Pt(int(v11.X5), dy+int(v11.Y6)))
									return v19 != nil && v19.Field12 != 0
								}() {
								if v69 < 4 {
									c.minimapDoorLine(v15, xLeft, yTop)
									v8 = xLeft.X
								}
								goto LABEL_37
							}
						}
						v69++
						if off+8 >= 149276 {
							goto LABEL_37
						}
					}
				}
				if noxflags.HasGame(noxflags.GameNotQuest) || v11.Field12 != 0 { // 0x10000
					v26 := int(legacy.Get_nox_xxx_minimap_587000_149232())
					v25 := v26
					var v76 image.Point
					v76.X = 100 * (v14 - v8) / v26
					v76.Y = yTop + 100*(v80-v9)/v26
					v27 := uint8(v11.Flags4)
					var v28 int8
					if v27&4 != 0 {
						v28 = int8(*(*uint8)(unsafe.Add(v11.Data, 21)))
					}
					if v27&4 == 0 || (v28 != 3 && v28 != 2) {
						if v27&0x20 == 0 {
							c.sub_4730D0(&v76, v11.Dir0, 2300/v25)
						}
					}
				}
			}
		LABEL_37:
			v11 = v11.NextByY24
			v9 = xLeft.Y
		}
		v10 = i + 1
		v80 += 23
	}
	if noxflags.HasEngine(noxflags.EngineShowAI) { // NOX_ENGINE_FLAG_ENABLE_SHOW_AI
		pts := c.srv.AI.Paths.Points()
		v29 := len(pts) // sub_50CB00
		if v29 >= 2 {
			c.drawSetColor(legacy.Get_dword_8531A0_2572())
			for n := 0; n < v29-1; n++ {
				p0 := pts[n]
				p1 := pts[n+1]
				v32 := int(legacy.Get_nox_xxx_minimap_587000_149232())
				xLeft.X = int(int32(100*(uint64(int64(p0.X))-uint64(int64(v8))))) / v32
				c.r.AddPoint(image.Pt(xLeft.X, yTop+int(int32(100*(uint64(int64(p0.Y))-uint64(int64(v9)))))/v32))
				v34 := int(legacy.Get_nox_xxx_minimap_587000_149232())
				xLeft.X = int(int32(100*(uint64(int64(p1.X))-uint64(int64(v8))))) / v34
				c.r.AddPoint(image.Pt(xLeft.X, yTop+int(int32(100*(uint64(int64(p1.Y))-uint64(int64(v9)))))/v34))
				c.r.DrawLineFromPoints(c.r.Data().Color2())
			}
		}
		for j := c.minimapFirstMonster(); j != nil; j = c.minimapNextMonster() {
			c.drawSetColor(memmap.Uint32(0x85B3FC, 940))
			v37 := int(legacy.Get_nox_xxx_minimap_587000_149232())
			xLeft.X = int(int32(100*(uint64(int64(j.X))-uint64(int64(v8))))) / v37
			c.nox_xxx_minimapDrawPoint_473570(xLeft.X,
				yTop+int(int32(100*(uint64(int64(j.Y))-uint64(int64(v9)))))/v37)
		}
	}
	v73 := 0
	if memmap.Uint32(0x5D4594, 1096304) == 0 {
		*memmap.PtrUint32(0x5D4594, 1096304) = uint32(c.Things.IndByID("Crown"))
		*memmap.PtrUint32(0x5D4594, 1096308) = uint32(c.Things.IndByID("GameBall"))
	}
	v70 := nox_xxx_objGetTeamByNetCode_418C80(legacy.ClientPlayerNetCode())
	if v70 != nil && v70.Has() {
		v73 = 1
	}
	player := c.ClientPlayerUnit()
	for k := c.Objs.FirstMinimapList(); k != nil; k = k.Nox_xxx_cliNextMinimapObj_459EC0(k) {
		v41 := legacy.Nox_xxx_polygonIsPlayerInPolygon_4217B0(unsafe.Pointer(&k.PosVec), 0)
		if v41 != nil {
			if minimapPolyByte2(v41) != a2 {
				continue
			}
		} else if a2 != 1 {
			continue
		}
		v42 := int(legacy.Get_nox_xxx_minimap_587000_149232())
		xLeft.X = 100 * (k.PosVec.X - v8) / v42
		xLeft.Y = yTop + 100*(k.PosVec.Y-v9)/v42
		var v43 uint32
		if uint32(k.ObjClass)&0x400000 == 0 {
			v43 = memmap.Uint32(0x85B3FC, 940)
		} else {
			v43 = legacy.Get_nox_color_blue_2650684()
			if uint8(k.ObjSubClass)&8 == 0 {
				v43 = memmap.Uint32(0x85B3FC, 940)
			}
		}
		c.drawSetColor(v43)
		v44 := int(k.TypeIDVal)
		if v44 == int(memmap.Int32(0x5D4594, 1096304)) { // Crown
			if c.srv.Teams.First() != nil {
				// nop
			} else if v45 := c.Objs.FirstPlayerList(); v45 == nil {
				// nop
			} else {
				found := false
				for !v45.HasEnchant(30) {
					v45 = v45.Field104()
					if v45 == nil {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			// LABEL_64
			c.drawSetColor(legacy.Get_dword_8531A0_2572())
			v46 := nox_xxx_objGetTeamByNetCode_418C80(int(k.NetCode32))
			if v46 != nil {
				v47 := c.srv.Teams.ByID(v46.ID)
				if v47 != nil {
					c.drawSetColor(c.teamColor32(v47))
				}
			}
			c.sub_473420(&xLeft)
			continue
		}
		if v44 == int(memmap.Int32(0x5D4594, 1096308)) { // GameBall
			v49 := nox_xxx_objGetTeamByNetCode_418C80(int(k.NetCode32))
			if v49 != nil && v49.Has() {
				v51 := c.srv.Teams.ByID(v49.ID)
				if v51 != nil {
					c.drawSetColor(c.teamColor32(v51))
				}
			} else {
				c.drawSetColor(legacy.Get_nox_color_white_2523948())
			}
			c.nox_video_drawCircleRad3_4734F0(&xLeft) // nox_video_drawCircleRad3_4734F0
			continue
		}
		v53 := uint32(k.ObjClass)
		if v53&0x10000000 != 0 {
			if uint32(k.ObjFlags)&0x1000000 != 0 {
				c.drawSetColor(legacy.Get_nox_color_white_2523948())
				v54 := legacy.Sub_4B9470(*(*unsafe.Pointer)(unsafe.Add(k.C(), 436)))
				v55 := c.srv.Teams.ByID(server.TeamID(v54))
				if v55 != nil {
					c.drawSetColor(c.teamColor32(v55))
					c.sub_4733B0(&xLeft)
					continue
				}
				c.sub_4733B0(&xLeft)
				continue
			}
		} else {
			if v53&4 == 0 {
				c.nox_xxx_minimapDrawPoint_473570(xLeft.X, xLeft.Y)
				continue
			}
			if !noxflags.HasGame(noxflags.GameModeCTF) { // 32
				if k == player {
					c.sub_4735C0(xLeft.X, xLeft.Y)
				} else {
					c.nox_xxx_minimapDrawPoint_473570(xLeft.X, xLeft.Y)
				}
				continue
			}
			v56 := c.srv.Players.ByID(int(k.NetCode32))
			if v56 != nil {
				v81 := v56.WeaponEquip & 1 // *((uint32_t*)v56 + 1) & 1
				if v81 != 0 {
					c.drawSetColor(legacy.Get_nox_color_white_2523948())
					v57 := nox_xxx_objGetTeamByNetCode_418C80(int(k.NetCode32))
					if v57 != nil {
						var v55 *server.Team
						if v57.ID == 1 {
							v55 = c.srv.Teams.ByID(2)
						} else {
							v55 = c.srv.Teams.ByID(1)
						}
						if v55 != nil {
							c.drawSetColor(c.teamColor32(v55))
							c.sub_4733B0(&xLeft)
							continue
						}
					}
					c.sub_4733B0(&xLeft)
					continue
				}
			}
		}
	}
	v79 := legacy.Get_dword_8531A0_2572()
	for l := c.Objs.FirstPlayerList(); l != nil; l = l.Field104() {
		v60 := l.HasEnchant(30)
		v61 := int(l.NetCode32)
		v62 := nox_xxx_objGetTeamByNetCode_418C80(v61)
		v75 := c.srv.Players.ByID(int(l.NetCode32))
		var v76 int
		if v70 != nil && v62 != nil && v73 != 0 {
			v76 = bool2int(v70.SameAs(v62)) // nox_xxx_servCompareTeams_419150
			if v76 != 0 {
				goto LABEL_103
			}
		} else {
			v76 = 0
		}
		if !(v60 || l == player || legacy.Get_dword_8531A0_2576().Field3680&1 != 0) {
			continue
		}
	LABEL_103:
		{
			v63 := legacy.Nox_xxx_polygonIsPlayerInPolygon_4217B0(unsafe.Pointer(&l.PosVec), 0)
			if (v63 == nil || minimapPolyByte2(v63) == a2) && v75 != nil && (v75.Field3680&1) != 1 {
				v64 := int(legacy.Get_nox_xxx_minimap_587000_149232())
				xLeft.X = 100 * (l.PosVec.X - v8) / v64
				xLeft.Y = yTop + 100*(l.PosVec.Y-v9)/v64
				if l == player || v76 != 0 {
					c.drawSetColor(v79)
				} else {
					c.drawSetColor(memmap.Uint32(0x85B3FC, 940))
				}
				if v60 {
					if v62 != nil {
						v65 := c.srv.Teams.ByID(v62.ID)
						if v65 != nil {
							c.drawSetColor(c.teamColor32(v65))
						}
					}
					c.sub_473420(&xLeft)
				} else {
					c.nox_xxx_minimapDrawPoint_473570(xLeft.X, xLeft.Y)
				}
			}
		}
	}
	legacy.Sub_49F860() // the original returns this call's result
}

// minimapDoorLine draws the line for a door wall (shared by two branches of the wall loop).
func (c *Client) minimapDoorLine(v15 *client.Drawable, xLeft image.Point, yTop int) {
	v20 := int(legacy.Get_nox_xxx_minimap_587000_149232())
	v21 := uintptr(8 * int(v15.Field_74_4)) // 8 * *(unsigned char*)(v15 + 299)
	v22 := 100 * (v15.PosVec.X - xLeft.X) / v20
	v85 := 100 * (v15.PosVec.Y - xLeft.Y) / v20
	v23 := 100 * int(memmap.Int32(0x587000, 196184+v21)) / v20
	v24 := 100 * int(memmap.Int32(0x587000, 196188+v21)) / v20
	c.drawSetColor(memmap.Uint32(0x85B3FC, 940))
	c.r.AddPoint(image.Pt(v22, yTop+v85))
	c.r.AddPointRel(image.Pt(v23, v24))
	c.r.DrawLineFromPoints(c.r.Data().Color2())
}

func (c *Client) sub_4730D0(a1 *image.Point, a2 uint8, a3 int) int {
	result := 0
	if legacy.Get_nox_xxx_minimap_587000_149232() <= 2000 {
		v4 := memmap.Uint32(0x85B3FC, 956)
		result = int(a2)
		v5 := a3 / 2
		x, y := a1.X, a1.Y
		l := func(a, b, cc, d int, col uint32) int { return c.sub_473380(a, b, cc, d, col) }
		gray := func() uint32 { return memmap.Uint32(0x85B3FC, 956) }
		switch a2 {
		case 0:
			result = l(x, a3+y, x+a3, y, gray())
		case 1:
			result = l(x, y, x+a3, y+a3, gray())
		case 2:
			l(x, a3+y, x+a3, y, gray())
			result = l(x, y, x+a3, y+a3, v4)
		case 3:
			l(x, y, x+v5, y+v5, gray())
			result = l(x, a3+y, a3+x, y, v4)
		case 4:
			l(x, y, x+a3, y+a3, gray())
			result = l(v5+x, v5+y, a3+x, y, v4)
		case 5:
			l(x, a3+y, x+a3, y, gray())
			result = l(v5+x, v5+y, a3+x, y+a3, v4)
		case 6:
			l(x, y, x+a3, y+a3, gray())
			result = l(x, a3+y, v5+x, y+v5, v4)
		case 7:
			l(x, y, x+v5, y+v5, gray())
			result = l(v5+x, v5+y, a3+x, y, v4)
		case 8:
			l(v5+x, v5+y, x+a3, y+a3, gray())
			result = l(v5+x, v5+y, a3+x, y, v4)
		case 9:
			l(v5+x, v5+y, x+a3, y+a3, gray())
			result = l(v5+x, v5+y, x, y+a3, v4)
		case 10:
			l(x, y, x+v5, y+v5, gray())
			result = l(x, a3+y, v5+x, y+v5, v4)
		default:
			return result
		}
	} else {
		c.drawSetColor(memmap.Uint32(0x85B3FC, 956))
		c.r.DrawPixel(image.Pt(a1.X, a1.Y), c.r.Data().Color2()) // nox_client_drawPixel_49EFA0
	}
	return result
}

func (c *Client) sub_473380(a1, a2, a3, a4 int, a5 uint32) int {
	c.drawSetColor(a5)
	c.r.AddPoint(image.Pt(a1, a2))
	c.r.AddPoint(image.Pt(a3, a4))
	return bool2int(c.r.DrawLineFromPoints(c.r.Data().Color2()))
}

func (c *Client) sub_4733B0(a1 *image.Point) int {
	v1 := a1.Y + 4
	v2 := a1.X - 2
	c.r.AddPoint(image.Pt(v2, v1))
	v1 -= 8
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v2 += 4
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 += 4
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	c.r.AddPoint(image.Pt(v2-4, v1))
	return c.drawLine()
}

// drawLine is nox_client_drawLineFromPoints_49E4B0.
func (c *Client) drawLine() int {
	return bool2int(c.r.DrawLineFromPoints(c.r.Data().Color2()))
}

func (c *Client) sub_473420(a1 *image.Point) int {
	v1 := a1.Y + 6
	v2 := a1.X - 4
	c.r.AddPoint(image.Pt(v2, v1))
	v1 -= 12
	v2 -= 2
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 += 6
	v2 += 4
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 -= 6
	v2 += 2
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 += 6
	v2 += 2
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 -= 6
	v2 += 4
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	v1 += 12
	v2 -= 2
	c.r.AddPoint(image.Pt(v2, v1))
	c.drawLine()
	c.r.AddPoint(image.Pt(v2, v1))
	c.r.AddPoint(image.Pt(v2-8, v1))
	return c.drawLine()
}

func (c *Client) nox_video_drawCircleRad3_4734F0(a1 *image.Point) {
	c.r.DrawCircle(a1.X, a1.Y, 3, c.r.Data().Color2()) // nox_video_drawCircle_4B0B90
}

func (c *Client) nox_client_drawRectLines_473510(a1, a2, a3, a4 int) int {
	c.r.AddPoint(image.Pt(a1, a2))
	v4 := a1 + a3 - 1
	c.r.AddPoint(image.Pt(v4, a2))
	c.drawLine()
	c.r.AddPoint(image.Pt(v4, a2))
	v5 := a4 - 1 + a2
	c.r.AddPoint(image.Pt(v4, v5))
	c.drawLine()
	c.r.AddPoint(image.Pt(v4, v5))
	c.r.AddPoint(image.Pt(a1, v5))
	return c.drawLine()
}

func (c *Client) nox_xxx_minimapDrawPoint_473570(xLeft, yTop int) {
	z := legacy.Get_nox_xxx_minimap_587000_149232()
	if z > 1200 {
		c.r.DrawPoint(image.Pt(xLeft, yTop), bool2int(z < 1750)+2, c.r.Data().Color2()) // nox_xxx_drawPointMB_499B70
	} else {
		c.r.DrawPoint(image.Pt(xLeft, yTop), 4, c.r.Data().Color2())
	}
}

func (c *Client) sub_4735C0(xLeft, yTop int) {
	z := legacy.Get_nox_xxx_minimap_587000_149232()
	if z > 1200 {
		c.r.DrawPoint(image.Pt(xLeft, yTop), bool2int(z < 1750)+4, c.r.Data().Color2())
	} else {
		c.r.DrawPoint(image.Pt(xLeft, yTop), 6, c.r.Data().Color2())
	}
}
