#include "client__shell__noxworld.h"

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3.h"
#include "GAME3_2.h"
#include "GAME5_2.h"
#include "MixPatch.h"
#include "client__gui__window.h"
#include "common__magic__speltree.h"
#include "common__random.h"
#include "common__strman.h"
#include "defs.h"
#include "input_common.h"
#include "operators.h"

#include "client__gui__gamewin__gamewin.h"
#include "client__shell__mainmenu.h"
#include "client__shell__selchar.h"

extern uint32_t dword_5d4594_815096;
extern uint32_t dword_5d4594_815052;
extern uint32_t dword_5d4594_815100;
extern uint32_t dword_5d4594_814548;
extern uint32_t dword_5d4594_814992;
extern uint32_t nox_client_connError_814552;
extern uint32_t dword_5d4594_815056;
extern uint32_t dword_5d4594_814996;
extern uint32_t dword_5d4594_814988;
extern uint32_t dword_587000_87408;
extern uint64_t qword_5d4594_814956;
extern uint64_t qword_5d4594_815068;
extern uint32_t dword_5d4594_815044;
extern uint32_t nox_wol_server_result_cnt_815088;
extern void* dword_5d4594_814624;
extern uint32_t dword_5d4594_815032;
extern uint32_t dword_5d4594_815020;
extern uint32_t dword_5d4594_815024;
extern uint32_t dword_5d4594_815028;
extern uint32_t dword_5d4594_814984;
extern uint32_t dword_5d4594_815016;
extern uint32_t nox_game_createOrJoin_815048;
extern uint32_t dword_587000_87412;
extern uint32_t dword_5d4594_815000;
extern uint32_t nox_wol_wnd_gameList_815012;
extern nox_window* dword_5d4594_815004;
extern nox_window* nox_wol_wnd_world_814980;

int dword_5d4594_815104 = 0;
nox_gui_animation* nox_wnd_xxx_815040 = 0;

//----- (004383A0) --------------------------------------------------------
int sub_4383A0() {
	wchar2_t* v0; // eax
	int result;  // eax

	if (!wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_814992)) {
		nox_window_set_hidden(*(int*)&dword_5d4594_814992, 1);
		sub_489870();
	}
	nox_window_set_hidden(*(int*)&nox_wol_wnd_gameList_815012, 0);
	nox_window_set_hidden(*(int*)&dword_5d4594_814984, 1);
	nox_window_set_hidden(*(int*)&dword_5d4594_814988, 1);
	sub_46AD20(*(uint32_t**)&nox_wol_wnd_world_814980, 10006, 10007, 1);
	sub_46ACE0(*(uint32_t**)&nox_wol_wnd_world_814980, 10047, 10051, 0);
	v0 = nox_strman_loadString_40F1D0("ListJoinServer", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 418);
	nox_window_call_field_94(*(int*)&dword_5d4594_814996, 16385, (int)v0, 0);
	dword_587000_87408 = 1;
	return 0;
}

//----- (00439370) --------------------------------------------------------
void nox_client_gui_serverInfoBlock_4394D0(int a1);
void nox_client_gui_serverInfoBlockCheckExp_439370(int2* a1, int a2) {
	if ((*(uint8_t*)(a2 + 164) & 0x10) && (0)) {
		wchar2_t* v3 = nox_strman_loadString_40F1D0("GeneralPrint:InformExpansion", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 688);
		nox_xxx_dialogMsgBoxCreate_449A10(nox_wol_wnd_world_814980, 0, (int)v3, 33, 0, 0);
		sub_44A360(0);
		sub_44A4B0();
		return;
	}
	sub_439450(a1->field_0, a1->field_4, a1);
	sub_46B120(dword_5d4594_815000, 0);
	nox_xxx_wndShowModalMB_46A8C0(dword_5d4594_815000);
	sub_46C690(dword_5d4594_815000);
	nox_window_setPos_46A9B0(dword_5d4594_815000, a1->field_0, a1->field_4);
	nox_client_gui_serverInfoBlock_4394D0(a2);
	dword_5d4594_815056 = 1;
	*getMemU16Ptr(0x5D4594, 814604) = *(uint16_t*)(a2 + 109);
	if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_DISABLE_GRAPHICS_RENDERING)) {
		nox_xxx_wnd_46ABB0(*getMemIntPtr(0x5D4594, 815008), 0);
	}
}
//----- (004394D0) --------------------------------------------------------
void nox_client_gui_serverInfoBlock_4394D0(int a1) {
	int v1;            // ebp
	short v9;          // ax
	wchar2_t* v10;      // eax
	wchar2_t* v14;      // eax
	int v15;           // ebx
	int v16;           // esi
	int v17;           // edi
	int v18;           // ebp
	int v19;           // eax
	wchar2_t* v20;      // eax
	int v21;           // esi
	wchar2_t* v22;      // eax
	int v23;           // ebx
	int v24;           // edi
	int v25;           // ebp
	wchar2_t* v27;      // eax
	int v28;           // edi
	wchar2_t* v29;      // eax
	int v30;           // esi
	int v31;           // ebx
	wchar2_t* v32;      // eax
	wchar2_t* v33;      // [esp-8h] [ebp-14Ch]
	unsigned char v34; // [esp+10h] [ebp-134h]
	unsigned char v35; // [esp+10h] [ebp-134h]
	int v36;           // [esp+14h] [ebp-130h]

	char buf[256] = {0};
	wchar2_t wbuf[512] = {0};

	v36 = a1 + 111;
	v1 = *(uint8_t*)(a1 + 102) & 0x7F;
	nox_window_call_field_94(dword_5d4594_815004, 16399, 0, 0);
	wchar2_t* nameTitle = nox_strman_loadString_40F1D0("Name", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 749);
	nox_window_call_field_94(dword_5d4594_815004, 16397, nameTitle, 14);
	if (*(uint8_t*)(a1 + 120)) {
		strncpy(buf, (const char*)(a1 + 120), sizeof(buf)-1);
	} else {
		nox_sprintAddrPort_43BC80(a1 + 12, *(uint16_t*)(a1 + 109), buf);
	}
	buf[sizeof(buf)-1] = 0;
	nox_swprintf(wbuf, L"%S", buf);
	nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	if (0) {
		nox_window_call_field_94(dword_5d4594_815004, 16397, getMemAt(0x587000, 89332), -1);
		wchar2_t* hostTitle = nox_strman_loadString_40F1D0("GameHost", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 766);
		nox_window_call_field_94(dword_5d4594_815004, 16397, hostTitle, 14);
		sub_439CC0(*(uint32_t*)(a1 + 32), buf);
		nox_swprintf(wbuf, L"%S", buf);
		nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, getMemAt(0x587000, 89396), -1);
	wchar2_t* pingTitle = nox_strman_loadString_40F1D0("Ping", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 774);
	nox_window_call_field_94(dword_5d4594_815004, 16397, pingTitle, 14);
	if (*(uint32_t*)(a1 + 96) == 9999) {
		nox_swprintf(wbuf, L"--");
	} else {
		nox_swprintf(wbuf, L"%d", *(uint32_t*)(a1 + 96));
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	nox_window_call_field_94(dword_5d4594_815004, 16397, getMemAt(0x587000, 89464), -1);
	wchar2_t* modeTitle = nox_strman_loadString_40F1D0("GameType", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 787);
	nox_window_call_field_94(dword_5d4594_815004, 16397, modeTitle, 14);
	wchar2_t* gameMode = nox_gui_wol_gameModeString_43BCB0(*(uint16_t*)(a1 + 163));
	nox_window_call_field_94(dword_5d4594_815004, 16397, gameMode, -1);
	if (*(uint8_t*)(a1 + 164) & 0x10) { // quest
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89520), -1);
		wchar2_t* stageTitle = nox_strman_loadString_40F1D0("Stage", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 794);
		nox_window_call_field_94(dword_5d4594_815004, 16397, stageTitle, 14);
		nox_swprintf(wbuf, L"%d", *(unsigned short*)(a1 + 165));
		nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89580), -1);
	wchar2_t* mapTitle = nox_strman_loadString_40F1D0("Map", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 801);
	nox_window_call_field_94(dword_5d4594_815004, 16397, mapTitle, 14);
	strncpy(buf, (const char*)(a1 + 111), sizeof(buf)-1);
	buf[sizeof(buf)-1] = 0;
	nox_swprintf(wbuf, L"%S", buf);
	nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89636), -1);
	v9 = *(uint16_t*)(a1 + 163);
	if (v9 & 0xC000u) {
		if (v9 & 0x4000) {
			v33 = nox_strman_loadString_40F1D0("Individual", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 811);
		} else {
			v33 = nox_strman_loadString_40F1D0("Clan", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 815);
		}
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v33, 6);
		v10 = nox_strman_loadString_40F1D0("Ladder", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 817);
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v10, 6);
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89788), -1);
	wchar2_t* playersTitle = nox_strman_loadString_40F1D0("Occupancy", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 822);
	nox_window_call_field_94(dword_5d4594_815004, 16397, playersTitle, 14);
	nox_swprintf(wbuf, L"%d/%d\n", *(unsigned char*)(a1 + 103), *(unsigned char*)(a1 + 104));
	nox_window_call_field_94(dword_5d4594_815004, 16397, wbuf, -1);
	if ((*(uint8_t*)(a1 + 164) & 0x20) == 0) {
		return;
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89860), -1);
	wchar2_t* resTitle = nox_strman_loadString_40F1D0("Resolution", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 831);
	nox_window_call_field_94(dword_5d4594_815004, 16397, resTitle, 14);
	int rstr = get_video_mode_string(v1);
	nox_window_call_field_94(dword_5d4594_815004, 16397, rstr, -1);
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 89916), -1);
	v14 = nox_strman_loadString_40F1D0("DisabledSpells", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 836);
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v14, 14);
	v15 = 0;
	v16 = 1;
	v34 = 0;
	v17 = 1;
	v18 = 136;
	do {
		if (v16 == 0x80000000) {
			v16 = 1;
			++v34;
		} else {
			v16 *= 2;
		}
		if (nox_xxx_spellIsValid_424B50(v17) && !(v16 & *(uint32_t*)(v36 + 4 * v34 + 24)) &&
			nox_xxx_spellFlags_424A70(v17) & 0x7000000) {
			v15 = 1;
			v19 = nox_xxx_spellTitle_424930(v17);
			nox_window_call_field_94(dword_5d4594_815004, 16397, v19, 4);
		}
		++v17;
		--v18;
	} while (v18);
	if (!v15) {
		v20 = nox_strman_loadString_40F1D0("None", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 866);
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v20, 4);
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 90024), -1);
	v21 = 0;
	v22 = nox_strman_loadString_40F1D0("DisabledWeapons", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 872);
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v22, 14);
	v23 = 1;
	v35 = 0;
	v24 = 1;
	v25 = 27;
	do {
		if (!((unsigned char)v23 & *(uint8_t*)(v35 + v36 + 44))) {
			int v26 = sub_4159F0(v24);
			if (v26) {
				nox_window_call_field_94(dword_5d4594_815004, 16397, v26, -1);
				v21 = 1;
			}
		}
		if (v23 == 128) {
			v23 = 1;
			++v35;
		} else {
			v23 *= 2;
		}
		v24 *= 2;
		--v25;
	} while (v25);
	if (!v21) {
		v27 = nox_strman_loadString_40F1D0("None", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 899);
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v27, 4);
	}
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)getMemAt(0x587000, 90132), -1);
	v28 = 0;
	v29 = nox_strman_loadString_40F1D0("DisabledArmor", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 905);
	nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v29, 14);
	v30 = 1;
	v31 = 26;
	do {
		if (!(v30 & *(uint32_t*)(a1 + 159))) {
			int result = sub_415E80(v30);
			if (result) {
				v28 = 1;
				nox_window_call_field_94(dword_5d4594_815004, 16397, result, -1);
			}
		}
		v30 *= 2;
		--v31;
	} while (v31);
	if (!v28) {
		v32 = nox_strman_loadString_40F1D0("None", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 920);
		nox_window_call_field_94(dword_5d4594_815004, 16397, (int)v32, 4);
	}
}

//----- (00439E70) --------------------------------------------------------
int nox_xxx_windowMultiplayerSub_439E70(int a1, unsigned int a2, int* a3, int a4) {
	int v4;             // eax
	int v5;             // esi
	int* v6;            // eax
	int* v7;            // eax
	int v8;             // esi
	int v9;             // esi
	int v10;            // eax
	int result;         // eax
	const wchar2_t* v12; // eax
	int v13;            // eax
	int v14;            // eax
	wchar2_t* v15;       // eax
	int v16;            // eax
	uint32_t* v17;      // eax
	uint32_t* v18;      // eax
	int v19;            // esi
	int v20;            // esi
	wchar2_t* v21;       // eax
	int v22;            // eax
	wchar2_t* v23;       // eax
	char* v24;          // eax
	char* v25;          // eax
	char* v26;          // eax
	int v27;            // eax
	uint16_t v28;       // [esp-8h] [ebp-A0h]
	int v30;            // [esp+Ch] [ebp-8Ch]
	int v31;            // [esp+10h] [ebp-88h]
	char buf[24];        // [esp+14h] [ebp-84h]
	char v34[36];       // [esp+2Ch] [ebp-6Ch]
	char v35[72];       // [esp+50h] [ebp-48h]

	if (a2 > 0x4010) {
		if (a2 == 16403 || a2 == 16412) {
			nox_window_call_field_94(*(int*)&dword_5d4594_815016, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815020, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815028, a2, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815032, a2, (int)a3, 0);
		}
		return 0;
	}
	switch (a2) {
	case 0x4010u:
		if (nox_xxx_wndGetID_46B0A0(a3) == 10061) {
			nox_point mpos = nox_client_getMousePos_4309F0();
			dword_5d4594_814624 = sub_4A28C0(a4);
			nox_client_gui_serverInfoBlockCheckExp_439370(&mpos, dword_5d4594_814624);
			return 0;
		}
		return 0;
	case 0x17u:
		return 1;
	case 0x4000u:
		v27 = nox_xxx_wndGetID_46B0A0(a3);
		if (v27 >= 10043 && v27 <= 10044) {
			nox_window_call_field_94(*(int*)&dword_5d4594_815016, 0x4000, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815020, 0x4000, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, 0x4000, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815028, 0x4000, (int)a3, 0);
			nox_window_call_field_94(*(int*)&dword_5d4594_815032, 0x4000, (int)a3, 0);
			return 0;
		}
		return 0;
	}
	if (a2 != 16391) {
		return 0;
	}
	v4 = nox_xxx_wndGetID_46B0A0(a3);
	v5 = v4;
	if (v4 != 10043 && v4 != 10044 && v4 != 10035 && v4 != 10036) {
		nox_xxx_clientPlaySoundSpecial_452D80(766, 100);
	}
	if (v5 >= 10070) {
		nox_point mpos = nox_client_getMousePos_4309F0();
		v30 = mpos.x - 216;
		v31 = mpos.y - 27;
		v6 = (int*)sub_4A0020();
		if (sub_4A25C0(&v30, v6) >= 2) {
			v7 = (int*)sub_4A0020();
			*getMemU32Ptr(0x5D4594, 815036) = sub_4A2610(*(int*)&nox_wol_wnd_world_814980, &v30, v7);
		} else {
			dword_5d4594_814624 = sub_4A0490(v5 - 10070);
			nox_client_gui_serverInfoBlockCheckExp_439370(&mpos, dword_5d4594_814624);
		}
		nox_xxx_cursorSetTooltip_4776B0((wchar2_t*)getMemAt(0x5D4594, 815112));
	}
	if (v5 > 10006) {
		switch (v5) {
		case 10007:
			if (nox_game_createOrJoin_815048) {
				sub_4373A0();
			}
			return 0;
		case 10010:
			sub_4373A0();
			return 0;
		case 10047:
		case 10048:
		case 10049:
		case 10050:
		case 10051:
			sub_4379C0();
			nox_wol_servers_sortBtnHandler_4A0290(v5);
			sub_4A0390();
			return 0;
		case 10054:
			dword_587000_87412 = 0;
			sub_43A810();
			nox_client_refreshServerList_4378B0();
			if (nox_common_gameFlags_check_40A5C0(0x2000000)) {
				v24 = sub_4A7EF0();
				sub_439D90(*(uint32_t*)v24 + 216, *((uint32_t*)v24 + 1) + 27);
			} else {
				sub_439D90(0x198u, 0xEFu);
			}
			result = 0;
			break;
		case 10055:
			dword_587000_87412 = 1;
			sub_43A810();
			nox_client_refreshServerList_4378B0();
			if (!nox_common_gameFlags_check_40A5C0(0x2000000)) {
				sub_439D90(0x198u, 0xEFu);
				result = 0;
			} else {
				v25 = sub_4A7EF0();
				sub_439D90(*(uint32_t*)v25 + 216, *((uint32_t*)v25 + 1) + 27);
				return 0;
			}
			break;
		case 10056:
			dword_587000_87412 = 2;
			sub_43A810();
			nox_client_refreshServerList_4378B0();
			if (!nox_common_gameFlags_check_40A5C0(0x2000000)) {
				sub_439D90(0x198u, 0xEFu);
			} else {
				v24 = sub_4A7EF0();
				sub_439D90(*(uint32_t*)v24 + 216, *((uint32_t*)v24 + 1) + 27);
			}
			result = 0;
			break;
		case 10057:
			dword_587000_87412 = 3;
			sub_43A810();
			nox_client_refreshServerList_4378B0();
			if (nox_common_gameFlags_check_40A5C0(0x2000000)) {
				v26 = sub_4A7EF0();
				sub_439D90(*(uint32_t*)v26 + 216, *((uint32_t*)v26 + 1) + 27);
			} else {
				sub_439D90(0x198u, 0xEFu);
			}
			result = 0;
			break;
		default:
			return 0;
		}
		return result;
	}
	if (v5 == 10006) {
		nox_game_createOrJoin_815048 = 0;
		nox_client_refreshServerList_4378B0();
		return 0;
	}
	if (v5 > 10002) {
		v19 = v5 - 10003;
		if (v5 != 10003) {
			v20 = v19 - 1;
			if (!v20) {
				nox_game_createOrJoin_815048 = 0;
				sub_4383A0();
				nox_client_refreshServerList_4378B0();
				return 0;
			}
			if (v20 == 1) {
				nox_game_createOrJoin_815048 = 0;
				nox_window_set_hidden(*(int*)&nox_wol_wnd_gameList_815012, 1);
				nox_window_set_hidden(*(int*)&dword_5d4594_814984, 1);
				nox_window_set_hidden(*(int*)&dword_5d4594_814988, 1);
				nox_window_set_hidden(*(int*)&dword_5d4594_814992, 0);
				sub_46AD20(*(uint32_t**)&nox_wol_wnd_world_814980, 10006, 10007, 0);
				sub_46ACE0(*(uint32_t**)&nox_wol_wnd_world_814980, 10047, 10051, 1);
				v21 = nox_strman_loadString_40F1D0("FilterMsg", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 2186);
				nox_window_call_field_94(*(int*)&dword_5d4594_814996, 16385, (int)v21, 0);
				return 0;
			}
			return 0;
		}
		// host quest game
		LOBYTE(v22) = 1; // has expansion
		if (!v22) {
			v23 = nox_strman_loadString_40F1D0("GeneralPrint:InformExpansion", 0,
											   "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 2013);
			nox_xxx_dialogMsgBoxCreate_449A10(*(int*)&nox_wol_wnd_world_814980, 0, (int)v23, 33, 0, 0);
			sub_44A360(0);
			sub_44A4B0();
			return 0;
		}
		nox_game_createOrJoin_815048 = 1;
		nox_common_gameFlags_unset_40A540(0x10000);
		nox_client_setMousePos_430B10(408, 239);
		sub_4375C0(0);
		if (!wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_814992)) {
			sub_489870();
		}
		nox_xxx_setQuest_4D6F60(1);
		nox_xxx_cliShowHideTubes_470AA0(1);
		if (0) {
			if (*(int*)&dword_587000_87412 != -1) {
				return 0;
			}
			v16 = nox_common_randomInt_415FA0(10054, 10057);
			v17 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_world_814980, v16);
			nox_window_call_field_94(*(int*)&nox_wol_wnd_world_814980, 16391, (int)v17, 0);
			return 0;
		}
		dword_587000_87412 = 0;
		v18 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_world_814980, 10020);
		nox_window_call_field_93((int)v18, 5, 15663512, 0);
		return 0;
	}
	if (v5 == 10002) {
		// host regular game
		if (nox_common_gameFlags_check_40A5C0(0x1000000)) {
			return 0;
		}
		nox_game_createOrJoin_815048 = 1;
		nox_xxx_setGameFlags_40A4D0(0x10000);
		nox_xxx_cliShowHideTubes_470AA0(0);
		nox_client_setMousePos_430B10(408, 239);
		sub_4375C0(0);
		if (!wndIsShown_nox_xxx_wndIsShown_46ACC0(*(int*)&dword_5d4594_814992)) {
			sub_489870();
		}
		if (nox_common_gameFlags_check_40A5C0(0x2000000) || sub_4D6F30()) {
			return 0;
		}
		if (0) {
			if (*(int*)&dword_587000_87412 != -1) {
				return 0;
			}
			v16 = nox_common_randomInt_415FA0(10054, 10057);
			v17 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_world_814980, v16);
			nox_window_call_field_94(*(int*)&nox_wol_wnd_world_814980, 16391, (int)v17, 0);
			return 0;
		}
		dword_587000_87412 = 0;
		v18 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&nox_wol_wnd_world_814980, 10020);
		nox_window_call_field_93((int)v18, 5, 15663512, 0);
		return 0;
	}
	v8 = v5 - 4001;
	if (!v8) {
		if (sub_43AF80() == 6) {
			v12 = (const wchar2_t*)sub_449E60(4);
			wchar2_t* wbuf = &buf[4];
			nox_wcsncpy(wbuf, v12, 9u);
			wbuf[8] = 0;
			v28 = nox_client_getServerPort_43B320();
			v13 = nox_client_getServerAddr_43B300();
			sub_5550D0(v13, v28, buf);
			sub_43AF90(3);
			*(uint64_t*)&qword_5d4594_814956 = nox_platform_get_ticks() + 20000;
			sub_449EA0(0);
			return 0;
		}
		if (sub_43AF80() == 10) {
			v14 = sub_449E60(4);
			nox_sprintf(v35, "%S", v14);
			v15 = nox_strman_loadString_40F1D0("Finding", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 2209);
			sub_449E30((int)v15);
			sub_43AF90(11);
			sub_449EA0(0);
			return 0;
		}
		if (sub_43AF80() == 1) {
			sub_43A920();
			nox_game_showGameSel_4379F0();
			return 0;
		}
		if (dword_5d4594_815104) {
			*(uint64_t*)&qword_5d4594_815068 = nox_platform_get_ticks() + 1000;
			return 0;
		}
		return 0;
	}
	v9 = v8 - 1;
	if (!v9) {
		sub_43A920();
		nox_game_showGameSel_4379F0();
		return 0;
	}
	if (v9 != 5999 || dword_5d4594_815044) {
		return 0;
	}
	v10 = sub_43B340();
	if (v10 & 0x1000) {
		sub_4D6F80(1);
		nox_xxx_cliShowHideTubes_470AA0(1);
	}
	sub_43B460();
	if (sub_4D6F70()) {
		if (nox_client_countPlayerFiles04_4DC7D0()) {
			sub_4A7A70(1);
			nox_game_showSelChar_4A4DB0();
			dword_5d4594_815044 = 1;
			sub_4A2890();
			return 0;
		}
	} else if (nox_client_countPlayerFiles02_4DC630()) {
		sub_4A7A70(1);
		nox_game_showSelChar_4A4DB0();
		dword_5d4594_815044 = 1;
		sub_4A2890();
		return 0;
	}
	sub_4A7A70(0);
	nox_game_showSelClass_4A4840();
	dword_5d4594_815044 = 1;
	sub_4A2890();
	return 0;
}
//----- (0043A810) --------------------------------------------------------
void sub_43A810() {
	int v0;      // edi
	int* v1;     // esi
	wchar2_t* v2; // [esp-8h] [ebp-Ch]

	nox_window_set_hidden(*(int*)&dword_5d4594_814984, 0);
	nox_window_set_hidden(*(int*)&dword_5d4594_814988, 1);
	nox_xxx_wndSetIcon_46AE60(*(int*)&dword_5d4594_814984, *getMemU32Ptr(0x5D4594, 814900 + 4 * dword_587000_87412));
	if (nox_game_createOrJoin_815048 == 1) {
		v2 = nox_strman_loadString_40F1D0("CreateMsg", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 377);
	} else {
		v2 = nox_strman_loadString_40F1D0("JoinServer", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 382);
	}
	nox_window_call_field_94(*(int*)&dword_5d4594_814996, 16385, (int)v2, 0);
	if (0) {
		sub_46ACE0(*(uint32_t**)&dword_5d4594_814984, 10600, 10611, 1);
		v0 = 0;
		v1 = getMemIntPtr(0x587000, 87560);
		do {
			sub_46ACE0(*(uint32_t**)&dword_5d4594_814984, *v1, v1[1], v0 != dword_587000_87412);
			v1 += 2;
			++v0;
		} while ((int)v1 < (int)getMemAt(0x587000, 87592));
	}
	sub_49FDB0(*(int*)&dword_587000_87412);
}

//----- (0043B630) --------------------------------------------------------
uint32_t* sub_43B630() {
	wchar2_t* v0; // eax

	nox_xxx_wnd_46C6E0(dword_5d4594_815000);
	v0 = nox_strman_loadString_40F1D0("AttemptingConn", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 1284);
	return nox_xxx_dialogMsgBoxCreate_449A10(*(int*)&nox_wol_wnd_world_814980, 0, (int)v0, 34, 0, 0);
}

//----- (0043B6E0) --------------------------------------------------------
void sub_43B6E0() {
	wchar2_t* v1; // eax
	wchar2_t* v2; // [esp-10h] [ebp-10h]
	if (nox_wol_wnd_world_814980) {
		v2 = nox_strman_loadString_40F1D0("Kicked", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3331);
		v1 = nox_strman_loadString_40F1D0("Notification", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3331);
		nox_xxx_dialogMsgBoxCreate_449A10(0, (int)v1, (int)v2, 33, 0, 0);
		sub_44A360(1);
		dword_5d4594_815096 = 0;
	} else {
		dword_5d4594_815096 = 1;
	}
}

//----- (0043B750) --------------------------------------------------------
void sub_43B750() {
	wchar2_t* v1; // eax
	wchar2_t* v2; // [esp-10h] [ebp-10h]

	if (nox_wol_wnd_world_814980) {
		v2 = nox_strman_loadString_40F1D0("Timeout", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3351);
		v1 = nox_strman_loadString_40F1D0("Notification", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3351);
		nox_xxx_dialogMsgBoxCreate_449A10(0, (int)v1, (int)v2, 33, 0, 0);
		sub_44A360(1);
		dword_5d4594_815100 = 0;
	} else {
		dword_5d4594_815100 = 1;
	}
}

//----- (0043B7C0) --------------------------------------------------------
void nox_gui_wol_newServerLine_43B7C0(nox_gui_server_ent_t* srv) {
	wchar2_t wbuf[128];
	char buf[332];

	if (dword_587000_87408) { // LAN games
		nox_window_call_field_94(nox_wol_wnd_gameList_815012, 16397, (int)getMemAt(0x587000, 91164), 4);

		strncpy(buf, srv->server_name, 15);
		buf[15] = 0;
		if (!memcmp(buf, getMemAt(0x5D4594, 815120), 1u)) {
			nox_sprintAddrPort_43BC80(srv->addr, srv->port, buf);
		}
		nox_swprintf(wbuf, L"%S", buf);
		sub_43BC10(wbuf, 100);
		nox_window_call_field_94(dword_5d4594_815016, 16397, wbuf, 4);

		nox_swprintf(wbuf, L"%d/%d", srv->players, srv->max_players);
		nox_window_call_field_94(dword_5d4594_815020, 16397, wbuf, 4);

		wchar2_t* v6 = nox_gui_wol_gameModeString_43BCB0(srv->flags);
		if (srv->flags & 0x1000) { // quest
			nox_swprintf((wchar2_t*)getMemAt(0x5D4594, 814772), L"%s %d", v6, srv->quest_level);
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16397, (int)getMemAt(0x5D4594, 814772), 4);
		} else {
			nox_window_call_field_94(*(int*)&dword_5d4594_815024, 16397, v6, 4);
		}

		if (srv->ping == 9999) {
			nox_swprintf(wbuf, L"--");
		} else {
			nox_itow(srv->ping, wbuf, 10);
		}
		nox_window_call_field_94(dword_5d4594_815028, 16397, wbuf, 4);

		wbuf[0] = 0;
		if (srv->status & 0x20) {
			wchar2_t* v9 = nox_strman_loadString_40F1D0("Noxworld.wnd:private", 0,
													   "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3486);
			nox_wcscat(wbuf, v9);
		}
		if (srv->status & 0x10) {
			if (wbuf[0]) {
				nox_wcscat(wbuf, L"+");
			}
			wchar2_t* v10 = nox_strman_loadString_40F1D0("Noxworld.wnd:closed", 0,
														"C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3494);
			nox_wcscat(wbuf, v10);
		}
		if (!wbuf[0]) {
			wchar2_t* v11;
			if (srv->players < srv->max_players) {
				v11 = nox_strman_loadString_40F1D0("Open", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3505);
			} else {
				v11 = nox_strman_loadString_40F1D0("Full", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3501);
			}
			nox_wcscat(wbuf, v11);
		}
		nox_window_call_field_94(*(int*)&dword_5d4594_815032, 16397, wbuf, 4);
	} else { // WOL games
		int a1 = srv;
		memset(buf, 0, sizeof(buf));
		*(uint32_t*)&buf[8] = 257;
		*(uint32_t*)&buf[16] = nox_wol_wnd_world_814980;
		int v4;
		if (*(int*)&dword_587000_87412 == -1) {
			int v1 = sub_437860(*(short*)(a1 + 44), *(short*)(a1 + 46));
			short v2 = *getMemU16Ptr(0x587000, 87528 + 8 * v1);
			v1 += 10054;
			*(uint16_t*)(a1 + 44) -= v2;
			*(uint16_t*)(a1 + 46) -= *getMemU16Ptr(0x587000, 7098 + 8 * v1);
			*(uint16_t*)(a1 + 44) >>= 1;
			*(uint16_t*)(a1 + 46) >>= 1;
			uint32_t* v3 = nox_xxx_wndGetChildByID_46B0C0(*(uint32_t**)&dword_5d4594_814988, v1);
			v4 = nox_gui_newButtonOrCheckbox_4A91A0((int)v3, 1185, *(short*)(a1 + 44) - 5, *(short*)(a1 + 46) - 5, 10,
													10, buf);
		} else {
			*(uint16_t*)(a1 + 44) -= *getMemU16Ptr(0x587000, 87528 + 8 * dword_587000_87412);
			*(uint16_t*)(a1 + 46) -= *getMemU16Ptr(0x587000, 87530 + 8 * dword_587000_87412);
			v4 = nox_gui_newButtonOrCheckbox_4A91A0(*(int*)&dword_5d4594_814984, 1192, *(short*)(a1 + 44) - 10,
													*(short*)(a1 + 46) - 10, 20, 20, buf);
		}
		*(uint32_t*)(a1 + 28) = v4;
		sub_437320(a1);
		char v13[32];
		if (*(uint8_t*)(a1 + 120)) {
			strncpy(v13, (const char*)(a1 + 120), 15);
			v13[15] = 0;
		} else {
			nox_sprintAddrPort_43BC80(a1 + 12, *(uint16_t*)(a1 + 109), v13);
		}
		if (*(uint32_t*)(a1 + 96) == 9999) {
			nox_swprintf(wbuf, L"%S -- ms", v13);
		} else {
			nox_swprintf(wbuf, L"%S %dms", v13, *(uint32_t*)(a1 + 96));
		}
		nox_xxx_wndWddSetTooltip_46B000((wchar2_t*)(*(uint32_t*)(a1 + 28) + 36), wbuf);
		nox_xxx_wndSetProc_46B2C0(*(uint32_t*)(a1 + 28), nox_xxx_windowMultiplayerSub_439E70);
		uint32_t* result = *(uint32_t**)(a1 + 28);
		*result = *(uint32_t*)(a1 + 36) + 10070;
	}
}

//----- (0043BCB0) --------------------------------------------------------
wchar2_t* nox_gui_wol_gameModeString_43BCB0(short a1) {
	if (a1 & 0x1000) {
		return nox_strman_loadString_40F1D0("Quest", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3534);
	}
	if (a1 & 0x20) {
		return nox_strman_loadString_40F1D0("CTF", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3541);
	}
	if (a1 & 0x400) {
		return nox_strman_loadString_40F1D0("Highlander", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3545);
	}
	if (a1 & 0x10) {
		return nox_strman_loadString_40F1D0("KotR", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3550);
	}
	if (a1 & 0x40) {
		return nox_strman_loadString_40F1D0("Flagball", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3555);
	}
	if ((a1 & 0x80u) == 0) {
		return nox_strman_loadString_40F1D0("Arena", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3563);
	}
	return nox_strman_loadString_40F1D0("Chat", 0, "C:\\NoxPost\\src\\client\\shell\\noxworld.c", 3559);
}
