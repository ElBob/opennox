#include <math.h>
#include <stdio.h>
#include <string.h>

#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_2.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "client__video__draw_common.h"
#include "common__magic__speltree.h"
#include "common__random.h"
#include "common__strman.h"
#include "common__system__team.h"
#include "operators.h"
#include "server__gamemech__explevel.h"
#include "server__magic__plyrspel.h"
#include "server__script__activator.h"
#include "server__script__builtin.h"
#include "server__script__internal.h"
#include "server__script__script.h"

// TODO: move somewhere else
int nox_xxx_netSendChat_528AC0(nox_object_t* a1, wchar2_t* a2, wchar2_t a3);

extern unsigned int dword_5d4594_2386836;

unsigned int dword_5d4594_1599628 = 0;
int dword_5d4594_2386848 = 0;
unsigned int dword_5d4594_2386852 = 0;

//----- (00512E80) --------------------------------------------------------
int sub_512E80(wchar2_t* a1) {
	int v1; // eax

	v1 = dword_5d4594_1599628;
	if (*(int*)&dword_5d4594_1599628 < 1024) {
		*getMemU32Ptr(0x973F18, 26664 + 4 * dword_5d4594_1599628) = a1;
		dword_5d4594_1599628 = ++v1;
	}
	return v1 - 1;
}

//----- (00513B00) --------------------------------------------------------
void nox_xxx_playerCanCarryItem_513B00(nox_object_t* a1p, nox_object_t* a2p) {
	int a1 = a1p;
	int a2 = a2p;
	uint32_t* v2; // ebp
	int v3;       // edi
	int v4;       // esi
	int v5;       // eax
	int v6;       // eax
	float2 v7;    // [esp+0h] [ebp-8h]

	if (!*getMemU32Ptr(0x5D4594, 2386856)) {
		*getMemU32Ptr(0x5D4594, 2386856) = nox_xxx_getNameId_4E3AA0("Glyph");
	}
	if (sub_467B00(*(unsigned short*)(a2 + 4), 1) - *(int*)&dword_5d4594_2386848 <= 0) {
		v2 = 0;
		v3 = 999999;
		v4 = nox_xxx_inventoryGetFirst_4E7980(a1);
		if (v4) {
			do {
				if (!(*(uint8_t*)(v4 + 8) & 0x10)) {
					v5 = *(uint32_t*)(v4 + 16);
					if (!(v5 & 0x100) && *(unsigned short*)(v4 + 4) != *getMemU32Ptr(0x5D4594, 2386856) &&
						!nox_xxx_ItemIsDroppable_53EBF0(v4)) {
						v6 = nox_xxx_shopGetItemCost_50E3D0(1, 0, *(float*)&v4);
						if (v6 < v3) {
							v3 = v6;
							v2 = (uint32_t*)v4;
						}
					}
				}
				v4 = nox_xxx_inventoryGetNext_4E7990(v4);
			} while (v4);
			if (v2) {
				sub_4ED970(50.0, (float2*)(a1 + 56), &v7);
				nox_xxx_drop_4ED790(a1, v2, &v7);
				if (!dword_5d4594_2386852) {
					nox_xxx_netPriMsgToPlayer_4DA2C0(a1, "pickup.c:CarryingTooMuch", 0);
					dword_5d4594_2386852 = 1;
				}
			}
		}
	}
}

//----- (00514BE0) --------------------------------------------------------
int nox_script_SetQuestInt_514BE0() {
	int v0; // esi
	int v1; // eax

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	nox_xxx_journalQuestSet_500540(nox_script_getString_512E40(v0), v1);
	return 0;
}

//----- (00514C10) --------------------------------------------------------
int nox_script_SetQuestFloat_514C10() {
	int v0; // esi
	int v1; // eax

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	nox_xxx_journalQuestSetBool_5006B0(nox_script_getString_512E40(v0), v1);
	return 0;
}

//----- (00514C40) --------------------------------------------------------
int nox_script_GetQuestInt_514C40() {
	int v0; // eax
	int v1; // eax

	v0 = nox_script_pop();
	v1 = sub_500750(nox_script_getString_512E40(v0));
	nox_script_push(v1);
	return 0;
}

//----- (00514C60) --------------------------------------------------------
int nox_script_GetQuestFloat_514C60() {
	int v0; // eax
	int v2; // [esp+0h] [ebp-4h]

	v0 = nox_script_pop();
	*(float*)&v2 = sub_500770((char*)nox_script_getString_512E40(v0));
	nox_script_push(v2);
	return 0;
}

//----- (00514C90) --------------------------------------------------------
int nox_script_ResetQuestStatus_514C90() {
	int v0; // eax

	v0 = nox_script_pop();
	sub_5007E0(nox_script_getString_512E40(v0));
	return 0;
}

//----- (00515550) --------------------------------------------------------
int nox_script_JournalDelete_515550() {
	int v0; // esi
	int v1; // eax
	int v2; // eax

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	if (v1) {
		v2 = nox_server_scriptValToObjectPtr_511B60(v1);
		if (v2) {
			nox_xxx_comJournalEntryRemove_427630(v2, nox_script_getString_512E40(v0));
			return 0;
		}
	} else {
		nox_xxx_comRemoveEntryAll_427680(nox_script_getString_512E40(v0));
	}
	return 0;
}

//----- (005155A0) --------------------------------------------------------
int nox_script_JournalEdit_5155A0() {
	short v0; // si
	int v1;   // edi
	int v2;   // eax
	int v3;   // eax

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	v2 = nox_script_pop();
	if (v2) {
		v3 = nox_server_scriptValToObjectPtr_511B60(v2);
		if (v3) {
			nox_xxx_comJournalEntryUpdate_427720(v3, nox_script_getString_512E40(v1), v0);
			return 0;
		}
	} else {
		nox_xxx_comUpdateEntryAll_427770(nox_script_getString_512E40(v1), v0);
	}
	return 0;
}

//----- (00515C40) --------------------------------------------------------
int nox_script_SetRoamFlag_515C40(int a1) {
	char v1; // bl
	int v2;  // eax
	int v3;  // eax
	char v5; // [esp+1h] [ebp-1h]

	v5 = HIBYTE(a1);
	v1 = nox_script_pop();
	v2 = nox_script_pop();
	v3 = nox_server_scriptValToObjectPtr_511B60(v2);
	if (v3) {
		v5 = v1;
		sub_515C80(v3, &v5);
	}
	return 0;
}

//----- (00515CB0) --------------------------------------------------------
int nox_script_SetRoamFlagGroup_515CB0() {
	char v0;           // bl
	int v1;            // eax
	unsigned char* v2; // eax
	char v4;           // [esp+1h] [ebp-1h]

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	v4 = v0;
	v2 = (unsigned char*)nox_server_scriptGetGroup_57C0A0(v1);
	nox_server_scriptExecuteFnForEachGroupObj_502670(v2, 0, sub_515C80, (int)&v4);
	return 0;
}

//----- (00516190) --------------------------------------------------------
int nox_script_GiveExp_516190() {
	int v0; // esi
	int v1; // eax
	int v2; // eax

	v0 = nox_script_pop();
	v1 = nox_script_pop();
	v2 = nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		nox_xxx_plyrGiveExp_4EF3A0_exp_level(v2, *(float*)&v0);
	}
	return 0;
}

void nox_script_StartupScreen_516600_A() {
	int i;  // esi
	int v1; // esi
	int v2; // edi

	for (i = nox_xxx_getFirstPlayerUnit_4DA7C0(); i; i = nox_xxx_getNextPlayerUnit_4DA7F0(i)) {
		if (*(uint8_t*)(*(uint32_t*)(*(uint32_t*)(i + 748) + 276) + 2064) == 31) {
			break;
		}
	}
	sub_4277B0(i, 0xEu);
	v1 = nox_xxx_inventoryGetFirst_4E7980(i);
	if (v1) {
		do {
			v2 = nox_xxx_inventoryGetNext_4E7990(v1);
			if (*(uint8_t*)(v1 + 8) & 0x40) {
				nox_xxx_delayedDeleteObject_4E5CC0(v1);
			}
			v1 = v2;
		} while (v2);
	}
	*getMemU32Ptr(0x5D4594, 2386832) = 1;
}

//----- (005166A0) --------------------------------------------------------
int nox_script_IsTalking_5166A0() {
	char* v0; // eax
	int v1;   // ecx

	v0 = nox_common_playerInfoFromNum_417090(31);
	v1 = 0;
	if (v0 && *(uint32_t*)(*(uint32_t*)(*((uint32_t*)v0 + 514) + 748) + 284)) {
		v1 = 1;
	}
	nox_script_push(v1);
	return 0;
}

//----- (005166E0) --------------------------------------------------------
int nox_script_PlayerIsTrading_5166E0() {
	char* v0; // eax
	int v1;   // ecx

	v0 = nox_common_playerInfoFromNum_417090(31);
	v1 = 0;
	if (v0 && *(uint32_t*)(*(uint32_t*)(*((uint32_t*)v0 + 514) + 748) + 280)) {
		v1 = 1;
	}
	nox_script_push(v1);
	return 0;
}

//----- (00516720) --------------------------------------------------------
int nox_script_MakeFriendly_516720() {
	char* v0; // esi
	int v1;   // eax
	int v2;   // eax
	int v3;   // ecx
	int v4;   // ecx

	v0 = nox_common_playerInfoFromNum_417090(31);
	v1 = nox_script_pop();
	v2 = nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		v3 = *(uint32_t*)(v2 + 12);
		BYTE1(v3) |= 1u;
		*(uint32_t*)(v2 + 12) = v3;
		v4 = *((uint32_t*)v0 + 514);
		if (v4) {
			nox_xxx_unitSetOwner_4EC290(v4, v2);
		}
	}
	return 0;
}

//----- (00516760) --------------------------------------------------------
int nox_script_MakeEnemy_516760() {
	int v0; // eax
	int v1; // eax
	int v2; // edx

	v0 = nox_script_pop();
	v1 = nox_server_scriptValToObjectPtr_511B60(v0);
	if (v1) {
		v2 = *(uint32_t*)(v1 + 12);
		BYTE1(v2) &= 0xFEu;
		*(uint32_t*)(v1 + 12) = v2;
		nox_xxx_unitClearOwner_4EC300(v1);
	}
	return 0;
}

//----- (00516790) --------------------------------------------------------
int nox_script_builtin_516790(void* this) {
	int v1;   // eax
	int v2;   // eax
	void* v4; // [esp+0h] [ebp-4h]

	v4 = this;
	v1 = nox_script_pop();
	v2 = nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		nox_script_push((*(uint32_t*)(v2 + 12) >> 8) & 1);
	} else {
		nox_script_push((int)v4);
	}
	return 0;
}

//----- (005167D0) --------------------------------------------------------
int nox_script_BecomePet_5167D0() {
	char* v0; // esi
	int v1;   // eax
	int v2;   // eax
	int v3;   // ecx

	v0 = nox_common_playerInfoFromNum_417090(31);
	v1 = nox_script_pop();
	v2 = nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		v3 = *((uint32_t*)v0 + 514);
		if (v3) {
			nox_xxx_unitBecomePet_4E7B00(v3, v2);
		}
	}
	return 0;
}

//----- (00516810) --------------------------------------------------------
int nox_script_BecomeEnemy_516810() {
	char* v0;     // esi
	int v1;       // eax
	uint32_t* v2; // eax
	int v3;       // ecx

	v0 = nox_common_playerInfoFromNum_417090(31);
	v1 = nox_script_pop();
	v2 = (uint32_t*)nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		v3 = *((uint32_t*)v0 + 514);
		if (v3) {
			nox_xxx_monsterRemoveMonitors_4E7B60(v3, v2);
		}
	}
	return 0;
}

//----- (00516850) --------------------------------------------------------
int nox_script_builtin_516850(void* this) {
	int v1;   // eax
	int v2;   // eax
	void* v4; // [esp+0h] [ebp-4h]

	v4 = this;
	v1 = nox_script_pop();
	v2 = nox_server_scriptValToObjectPtr_511B60(v1);
	if (v2) {
		nox_script_push((*(uint32_t*)(v2 + 12) & 0xFFu) >> 7);
	} else {
		nox_script_push((int)v4);
	}
	return 0;
}

//----- (00516890) --------------------------------------------------------
int nox_script_OblivionGive_516890() {
	uint32_t* v0; // edi
	int v1;       // esi
	int v2;       // ebx
	uint32_t* v3; // eax
	uint32_t* v4; // eax

	v0 = (uint32_t*)*((uint32_t*)nox_common_playerInfoFromNum_417090(31) + 514);
	v1 = 0;
	v2 = nox_script_pop();
	v3 = (uint32_t*)v0[126];
	if (v3) {
		while (!(v3[2] & 0x1000000) || !(v3[3] & 0x7800000)) {
			v3 = (uint32_t*)v3[124];
			if (!v3) {
				goto LABEL_7;
			}
		}
		v1 = (v3[4] >> 8) & 1;
		nox_xxx_delayedDeleteObject_4E5CC0((int)v3);
	}
LABEL_7:
	v4 = nox_xxx_playerRespawnItem_4EF750((int)v0, *(char**)getMemAt(0x587000, 247336 + 4 * v2), 0, 1, 1);
	if (v1) {
		nox_xxx_playerTryEquip_4F2F70(v0, (int)v4);
	}
	return 0;
}
