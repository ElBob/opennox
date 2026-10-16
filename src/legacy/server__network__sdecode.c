#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_1.h"
#include "GAME4_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "MixPatch.h"
#include "client__gui__guicon.h"
#include "client__gui__window.h"
#include "client__network__cdecode.h"
#include "client__system__ctrlevnt.h"
#include "client__system__parsecmd.h"
#include "common__magic__speltree.h"
#include "common__net_list.h"
#include "common__strman.h"
#include "common__system__team.h"
#include "defs.h"
#include "operators.h"

#include "common/fs/nox_fs.h"
#include "server__system__trade.h"

//----- (0051BAD0) --------------------------------------------------------
int nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(int a1, unsigned char* data, int dsz, nox_playerInfo* v8p, nox_object_t* unitp, void* v10p) {
	int v8 = v8p;
	int unit = unitp;
	int* v10 = v10p;
	int v85 = v10;
	int v7;                          // edx
	char* v11 = 0;                       // ecx
	int v14;                         // ecx
	int v15;                         // edi
	int v16;                         // eax
	int v17;                         // edi
	int v18;                         // eax
	int v19;                         // edi
	int v20;                         // eax
	unsigned short v21;              // dx
	int v22;                         // edi
	int v23;                         // ebx
	int v24;                         // edi
	int v26;                         // edi
	int v27;                         // eax
	int v28;                         // edi
	int v29;                         // eax
	void (*v30)(int, int, uint32_t); // ecx
	int v31;                         // edi
	int v32;                         // eax
	int v34;                         // edi
	int v35;                         // eax
	int* v36;                        // edi
	int v37;                         // ebx
	uint32_t* v38;                   // eax
	int v39;                         // ecx
	int v40;                         // ebx
	int v57;                         // eax
	int v58;                         // edi
	int v59;                         // eax
	int v60;                         // ebx
	int v61;                         // edi
	unsigned char v62;               // al
	uint32_t* v63;                   // eax
	int v64;                         // eax
	int v65;                         // eax
	int v66;                         // eax
	int v67;                         // edi
	int* v68;                        // edi
	int v69;                         // eax
	int v70;                         // edi
	int v71;                         // eax
	uint32_t* v72;                   // eax
	int v73;                         // eax
	int v74;                         // eax
	int v75;                         // eax
	int v76;                         // eax
	int v77;                         // eax
	int v78;                         // eax
	unsigned char v79;               // al
	int v80;                         // eax
	int v81;                         // ecx
	unsigned int v82;                // [esp-8h] [ebp-65Ch]
	float v83;                       // [esp+0h] [ebp-654h]
	float v84;                       // [esp+4h] [ebp-650h]
	char v88[4];                     // [esp+24h] [ebp-630h]
	char v90[6];                     // [esp+2Ch] [ebp-628h]
	float2 v94;                      // [esp+44h] [ebp-610h]

	int op = data[0];
	switch (op) {
	case 0xEEu: // MSG_VOTE
		switch (data[1]) {
		case 0u:;
			int c1 = 0;
			if (nox_common_gameFlags_check_40A5C0(4096)) {
				c1 = 3;
			}
			sub_506870(c1, unit, (wchar2_t*)data + 1);
			return 52;
		case 1u:
			sub_506870(1, unit, (wchar2_t*)data + 1);
			return 52;
		case 2u:;
			int c2 = 0;
			if (nox_common_gameFlags_check_40A5C0(4096)) {
				c2 = 3;
			}
			sub_506C90(c2, unit, (wchar2_t*)data + 1);
			return 52;
		case 3u:
			sub_506C90(1, unit, (wchar2_t*)data + 1);
			return 52;
		case 4u:
			sub_506870(2, unit, 0);
			return 2;
		case 5u:
			sub_506C90(2, unit, 0);
			return 2;
		default:
			return -1;
		}
	case 0xC9u: // MSG_TRADE
		switch (data[1]) {
		case 0xEu:
			if (v10[70]) {
				sub_50F3A0((uint32_t*)v10[70]);
			}
			return 2;
		case 0xFu:
			v67 = nox_xxx_packetDynamicUnitCode_578B40(*((unsigned short*)data + 1));
			if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_ENABLE_NET_DEBUG)) {
				nox_xxx_netTestHighBit_578B70(*((unsigned short*)data + 1));
			}
			v68 = (int*)nox_xxx_equipedItemByCode_4F7920(unit, v67);
			if (v68) {
				v69 = *(uint32_t*)(v85 + 280);
				if (v69) {
					if (nox_xxx_tradeP2PAddOffer2_50F820_trade(v69, unit, *(float*)&v68) == 1) {
						sub_4ED0C0(unit, v68);
					}
				}
			}
			v10 = (int*)v85;
			return 4;
		case 0x10u:
			v70 = nox_xxx_packetDynamicUnitCode_578B40(*((unsigned short*)data + 1));
			if (nox_common_getEngineFlag(NOX_ENGINE_FLAG_ENABLE_NET_DEBUG)) {
				nox_xxx_netTestHighBit_578B70(*((unsigned short*)data + 1));
			}
			v71 = *(uint32_t*)(v85 + 280);
			if (v71) {
				nox_xxx_tradeP2PAddOfferMB_50FE20(v71, v70);
			}
			v10 = (int*)v85;
			return 4;
		case 0x11u:
			v66 = v10[70];
			if (v66) {
				nox_xxx_tradeAccept_50F5A0(v66, unit);
			}
			return 2;
		case 0x12u:
			if (v10[70]) {
				nox_xxx_shopExit_50F4C0((uint32_t*)v10[70]);
			}
			return 2;
		case 0x15u:
			if (!nox_xxx_gameGet_4DB1B0() && !(*(uint8_t*)(v10[69] + 3680) & 3)) {
				v64 = nox_xxx_packetDynamicUnitCode_578B40(*((unsigned short*)data + 1));
				v65 = nox_server_getObjectFromNetCode_4ECCB0(v64);
				if (v65) {
					if (*(uint8_t*)(v65 + 12) & 8) {
						nox_xxx_servShopStart_50EF10_trade(unit, v65);
					}
				}
			}
			return 4;
		case 0x16u:
			v72 = (uint32_t*)v10[70];
			if (v72) {
				sub_5100C0_trade(unit, v72, *((unsigned short*)data + 1));
			}
			return 4;
		case 0x17u:
			v73 = v10[70];
			if (v73) {
				sub_510640_trade(unit, v73, *((unsigned short*)data + 1), (float*)data[4]);
			}
			return 5;
		case 0x18u:
			v74 = v10[70];
			if (v74) {
				sub_510BE0_trade((int*)unit, v74, (uint32_t*)*((unsigned short*)data + 1));
			}
			return 4;
		case 0x19u:
			v76 = v10[70];
			if (v76) {
				sub_510D10((int*)unit, v76, *((unsigned short*)data + 1), data[4]);
			}
			return 5;
		case 0x1Au:
			v75 = v10[70];
			if (v75) {
				sub_510AE0((int*)unit, v75, (uint32_t*)*((unsigned short*)data + 1));
			}
			return 4;
		case 0x1Cu:
			v77 = v10[70];
			if (v77) {
				sub_5109C0_trade((int*)unit, v77, (uint32_t*)*((unsigned short*)data + 1));
			}
			return 4;
		case 0x1Eu:
			v78 = v10[70];
			if (v78) {
				sub_5108D0(unit, v78, *((unsigned short*)data + 1));
			}
			return 4;
		default:
			return -1;
		}
	case 0xF0u: // MSG_GAUNTLET
		v79 = data[1];
		if (v79 == 3) {
			v80 = v10[69];
			v81 = *(uint32_t*)(v80 + 2056);
			if (v81) {
				v7 = *(uint32_t*)(v81 + 16);
				if ((v7 & 0x8000) != 0) {
					v10[137] = 0;
					nox_xxx_playerRespawn_4F7EF0(*(uint32_t*)(v80 + 2056));
				}
			}
		} else {
			if (v79 != 27) {
				return -1;
			}
			sub_4DD0B0(unit);
		}
		return 2;
	default:
		return -1;
	}
}
