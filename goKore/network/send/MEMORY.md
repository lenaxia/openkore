# Send.pm Implementation Checklist

This file tracks the implementation status of subroutines from the original Perl `Send.pm` file in the Go version of OpenKore.

# Implementation Status

All 301 functions from the original Perl `Send.pm` file have been successfully implemented in the Go version of OpenKore. The functions have been organized into logical systems and implemented in appropriate Go files.

# Send.pm Implementation Checklist (Reorganized)

## Core Functionality
- [x] Line 100: `encryptMessageID` - core/encryption.go
- [x] Line 149: `cryptKeys` - core/base_send.go
- [x] Line 162: `injectMessage` - core/base_send.go
- [x] Line 178: `injectAdminMessage` - core/base_send.go
- [x] Line 194: `pinEncode` - core/base_send.go
- [x] Line 230: `sendToServer` - core/base_send.go
- [x] Line 289: `sendRaw` - core/base_send.go
- [x] Line 3489: `sendPing` - core/ping.go

## Login and Authentication
- [x] Line 302: `parse_master_login` - security/login.go
- [x] Line 319: `reconstruct_master_login` - security/login.go
- [x] Line 346: `sendMasterLogin` - security/login.go
- [x] Line 379: `secureLoginHash` - security/login.go
- [x] Line 394: `sendMasterSecureLogin` - security/login.go
- [x] Line 409: `reconstruct_game_login` - security/login.go
- [x] Line 428: `sendGameLogin` - security/login.go
- [x] Line 441: `sendCharLogin` - security/login.go
- [x] Line 447: `sendMapLogin` - security/login.go
- [x] Line 482: `sendMapLoaded` - security/login.go
- [x] Line 1108: `sendLoginPinCode` - security/pin.go
- [x] Line 3183: `sendClientVersion` - security/login.go
- [x] Line 981: `reconstruct_client_hash` - security/client_hash.go
- [x] Line 1007: `sendClientMD5Hash` - security/client_hash.go

## Character System
- [x] Line 2296: `sendCharCreate` - game/actor/character.go
- [x] Line 2331: `sendCharDelete` - game/actor/character.go
- [x] Line 952: `sendCharDelete2` - game/actor/character.go
- [x] Line 962: `sendCharDelete2Accept` - game/actor/character.go
- [x] Line 968: `reconstruct_char_delete2_accept` - game/actor/character.go
- [x] Line 974: `sendCharDelete2Cancel` - game/actor/character.go
- [x] Line 1683: `sendchangetitle` - game/actor/character.go
- [x] Line 1638: `sendAddStatusPoint` - game/actor/character.go
- [x] Line 1647: `sendAddSkillPoint` - game/actor/character.go
- [x] Line 1655: `sendHotKeyChange` - game/actor/character.go
- [x] Line 670: `sendRespawn` - game/actor/character.go
- [x] Line 673: `sendQuitToCharSelect` - game/actor/character.go
- [x] Line 677: `sendRestart` - game/actor/character.go
- [x] Line 872: `sendQuit` - game/actor/character.go
- [x] Line 3094: `sendAutoRevive` - game/actor/character.go

## Movement and Synchronization
- [x] Line 490: `reconstruct_sync` - game/actor/sync.go
- [x] Line 495: `sendSync` - game/actor/sync.go
- [x] Line 504: `parse_character_move` - game/actor/movement.go
- [x] Line 509: `reconstruct_character_move` - game/actor/movement.go
- [x] Line 517: `sendMove` - game/actor/movement.go
- [x] Line 1015: `parse_actor_move` - game/actor/movement_utils.go
- [x] Line 1020: `reconstruct_actor_move` - game/actor/movement_utils.go
- [x] Line 1028: `sendSlaveMove` - game/actor/movement.go
- [x] Line 2377: `sendWarpTele` - game/movement/teleport.go
- [x] Line 1906: `sendPrivateAirshipRequest` - game/movement/teleport.go

## Actions
- [x] Line 529: `sendAction` - game/actor/action.go
- [x] Line 644: `sendLook` - game/actor/action.go
- [x] Line 1212: `sendSlaveAttack` - game/actor/action.go
- [x] Line 1228: `sendSlaveStandBy` - game/actor/action.go
- [x] Line 2841: `sendAlignment` - game/actor/action.go
- [x] Line 1967: `sendEmotion` - game/actor/action.go


## Chat System
- [x] Line 547: `parse_public_chat` - game/chat/public.go
- [x] Line 552: `reconstruct_public_chat` - game/chat/public.go
- [x] Line 557: `sendChat` - game/chat/public.go
- [x] Line 2716: `sendGMBroadcast` - game/chat/public.go
- [x] Line 2794: `sendGMBroadcastLocal` - game/chat/public.go
- [x] Line 574: `sendTalk` - game/chat/private.go
- [x] Line 582: `sendTalkCancel` - game/chat/private.go
- [x] Line 589: `sendTalkContinue` - game/chat/private.go
- [x] Line 595: `sendTalkResponse` - game/chat/private.go
- [x] Line 603: `sendTalkNumber` - game/chat/private.go
- [x] Line 611: `sendTalkText` - game/chat/private.go
- [x] Line 625: `parse_private_message` - game/chat/private.go
- [x] Line 632: `reconstruct_private_message` - game/chat/private.go
- [x] Line 639: `sendPrivateMsg` - game/chat/private.go
- [x] Line 562: `sendGetPlayerInfo` - game/chat/info.go
- [x] Line 568: `sendGetCharacterName` - game/chat/info.go
- [x] Line 1978: `sendWho` - game/chat/info.go
- [x] Line 3102: `sendBattlegroundChat` - game/chat/info.go
- [x] Line 732: `parse_party_chat` - game/social/party.go
- [x] Line 737: `reconstruct_party_chat` - game/social/party.go
- [x] Line 742: `sendPartyChat` - game/social/party.go
- [x] Line 857: `parse_guild_chat` - game/social/guild.go
- [x] Line 862: `reconstruct_guild_chat` - game/social/guild.go
- [x] Line 867: `sendGuildChat` - game/social/guild.go
- [x] Line 1677: `sendClanChat` - game/chat/info.go

## Inventory System
- [x] Line 651: `sendTake` - game/item/inventory.go
- [x] Line 657: `sendDrop` - game/item/inventory.go
- [x] Line 663: `sendItemUse` - game/item/inventory.go
- [x] Line 1243: `sendEquip` - game/item/inventory.go
- [x] Line 1257: `sendEquipSwitchAdd` - game/item/inventory.go
- [x] Line 1272: `sendEquipSwitchRemove` - game/item/inventory.go
- [x] Line 1286: `sendEquipSwitchRun` - game/item/inventory.go
- [x] Line 1298: `sendEquipSwitchSingle` - game/item/inventory.go
- [x] Line 1630: `sendUnequip` - game/item/inventory.go
- [x] Line 1849: `sendChangeDress` - game/item/inventory.go
- [x] Line 3476: `sendInventoryExpansionRequest` - game/item/inventory.go
- [x] Line 3483: `sendInventoryExpansionRejected` - game/item/inventory.go

## Storage System
- [x] Line 683: `sendStorageAdd` - game/item/storage.go
- [x] Line 693: `sendStorageGet` - game/item/storage.go
- [x] Line 703: `sendStoragePassword` - game/item/storage_password.go
- [x] Line 715: `reconstruct_storage_password` - game/item/storage_password.go
- [x] Line 2171: `sendStorageClose` - game/item/storage.go (as CloseStorage)
- [x] Line 2392: `sendStorageGetToCart` - game/item/storage_password.go
- [x] Line 2411: `sendStorageAddFromCart` - game/item/storage_password.go

## Cart System
- [x] Line 2238: `sendCartAdd` - game/item/cart.go
- [x] Line 2250: `sendCartGet` - game/item/cart.go
- [x] Line 1935: `sendChangeCart` - game/item/cart.go

## Skill System
- [x] Line 822: `sendSkillUse` - game/skill/skill.go
- [x] Line 836: `sendSkillUseLoc` - game/skill/skill.go
- [x] Line 3128: `sendSkillUseLocInfo` - game/skill/skill.go
- [x] Line 1087: `sendSkillSelect` - game/skill/skill.go
- [x] Line 3225: `sendStartSkillUse` - game/skill/skill.go
- [x] Line 3233: `sendStopSkillUse` - game/skill/skill.go
- [x] Line 1958: `sendAutoSpell` - game/skill/skill.go

## Party System
- [x] Line 887: `sendPartyOption` - game/party/party.go
- [x] Line 900: `sendPartyLeader` - game/party/party.go
- [x] Line 2179: `sendPartyJoinRequest` - game/party/party.go
- [x] Line 2190: `sendPartyJoin` - game/party/party.go
- [x] Line 2202: `sendPartyLeave` - game/party/party.go
- [x] Line 2210: `sendPartyKick` - game/party/party.go
- [x] Line 1075: `sendPartyJoinRequestByName` - game/party/party.go
- [x] Line 3082: `sendPartyJoinRequestByNameReply` - game/party/party.go
- [x] Line 912: `sendPartyBookingRegister` - game/party/party.go
- [x] Line 923: `sendPartyBookingReqSearch` - game/party/party.go
- [x] Line 934: `sendPartyBookingDelete` - game/party/party.go
- [x] Line 941: `sendPartyBookingUpdate` - game/party/party.go

## Guild System
- [x] Line 842: `sendGuildMasterMemberCheck` - game/guild/guild.go
- [x] Line 848: `sendGuildRequestInfo` - game/guild/guild.go
- [x] Line 2343: `sendGuildAlly` - game/guild/guild.go
- [x] Line 2355: `sendGuildRequestEmblem` - game/guild/guild.go
- [x] Line 2366: `sendGuildBreak` - game/guild/guild.go
- [x] Line 2441: `sendGuildLeave` - game/guild/guild.go
- [x] Line 2455: `sendGuildMemberKick` - game/guild/guild.go
- [x] Line 2469: `sendGuildCreate` - game/guild/guild.go
- [x] Line 2481: `sendGuildJoin` - game/guild/guild.go
- [x] Line 2493: `sendGuildJoinRequest` - game/guild/guild.go
- [x] Line 2506: `sendGuildNotice` - game/guild/guild.go
- [x] Line 2519: `sendGuildSetAlly` - game/guild/guild.go

## Friend System
- [x] Line 1041: `sendFriendListReply` - game/friend/friend.go
- [x] Line 1052: `sendFriendRequest` - game/friend/friend.go
- [x] Line 1857: `sendFriendRemove` - game/friend/friend.go
- [x] Line 1999: `sendIgnore` - game/friend/friend.go
- [x] Line 2013: `sendIgnoreAll` - game/friend/friend.go
- [x] Line 2024: `sendGetIgnoreList` - game/friend/friend.go

## NPC Interaction
- [x] Line 1985: `sendNPCBuySellList` - game/npc/npc.go
- [x] Line 3542: `sendNPCCreateRequest` - game/npc/npc.go

## Shop System
- [x] Line 748: `parse_buy_bulk_vender` - game/shop/shop.go
- [x] Line 753: `reconstruct_buy_bulk_vender` - game/shop/shop.go
- [x] Line 760: `sendBuyBulkVender` - game/shop/shop.go
- [x] Line 771: `reconstruct_buy_bulk_buyer` - game/shop/shop.go
- [x] Line 778: `sendBuyBulkBuyer` - game/shop/shop.go
- [x] Line 792: `sendEnteringBuyer` - game/shop/shop.go
- [x] Line 798: `sendBuyBulkOpenShop` - game/shop/shop.go
- [x] Line 815: `reconstruct_buy_bulk_openShop` - game/shop/shop.go
- [x] Line 878: `sendCloseShop` - game/shop/shop.go
- [x] Line 1131: `sendCloseBuyShop` - game/shop/shop.go
- [x] Line 1622: `sendEnteringVender` - game/shop/shop.go
- [x] Line 2854: `sendOpenShop` - game/shop/shop.go
- [x] Line 2865: `reconstruct_shop_open` - game/shop/shop.go

## Market System
- [x] Line 3444: `sendBuyBulkMarket` - game/market/market.go
- [x] Line 3455: `reconstruct_buy_bulk_market` - game/market/market.go
- [x] Line 3464: `sendMarketClose` - game/market/market.go
- [x] Line 3432: `sendSellBuyComplete` - game/market/market.go

## Crafting System
- [x] Line 2262: `sendIdentify` - game/craft/craft.go
- [x] Line 1749: `sendWeaponRefine` - game/craft/craft.go
- [x] Line 1762: `sendCooking` - game/craft/craft.go
- [x] Line 1772: `sendMakeItemRequest` - game/craft/craft.go
- [x] Line 1869: `sendRepairItem` - game/craft/craft.go
- [x] Line 1947: `sendArrowCraft` - game/craft/craft.go
- [x] Line 2273: `sendCardMergeRequest` - game/craft/craft.go
- [x] Line 2284: `sendCardMerge` - game/craft/craft.go

## RODEX Mail System
- [x] Line 1469: `rodex_delete_mail` - game/rodex/rodex.go
- [x] Line 1479: `rodex_request_zeny` - game/rodex/rodex.go
- [x] Line 1489: `rodex_request_items` - game/rodex/rodex.go
- [x] Line 1499: `rodex_cancel_write_mail` - game/rodex/rodex.go
- [x] Line 1507: `rodex_add_item` - game/rodex/rodex.go
- [x] Line 1516: `rodex_remove_item` - game/rodex/rodex.go
- [x] Line 1525: `rodex_open_write_mail` - game/rodex/rodex.go
- [x] Line 1533: `rodex_checkname` - game/rodex/rodex.go
- [x] Line 1543: `rodex_send_mail` - game/rodex/rodex.go
- [x] Line 1564: `rodex_refresh_maillist` - game/rodex/rodex.go
- [x] Line 1579: `rodex_read_mail` - game/rodex/rodex.go
- [x] Line 1589: `rodex_next_maillist` - game/rodex/rodex.go
- [x] Line 1599: `rodex_open_mailbox` - game/rodex/rodex.go
- [x] Line 1614: `rodex_close_mailbox` - game/rodex/rodex.go

## Mail System
- [x] Line 2871: `sendMailboxOpen` - game/mail/mail.go
- [x] Line 2879: `sendMailRead` - game/mail/mail.go
- [x] Line 2890: `sendMailDelete` - game/mail/mail.go
- [x] Line 2901: `sendMailGetAttach` - game/mail/mail.go
- [x] Line 2912: `sendMailOperateWindow` - game/mail/mail.go
- [x] Line 2923: `sendMailSetAttach` - game/mail/mail.go
- [x] Line 2944: `sendMailSend` - game/mail/mail.go
- [x] Line 2958: `reconstruct_mail_send` - game/mail/mail.go
- [x] Line 2964: `sendMailReturn` - game/mail/mail.go

## Banking
- [x] Line 3332: `sendBankingCheck` - game/banking/banking.go
- [x] Line 3342: `sendBankingWithdraw` - game/banking/banking.go
- [x] Line 3353: `sendBankingDeposit` - game/banking/banking.go

## Homunculus
- [x] Line 2430: `sendHomunculusName` - game/homunculus/homunculus.go
- [x] Line 1065: `sendHomunculusCommand` - game/homunculus/homunculus.go

## Pet System
- [x] Line 2538: `sendPetCapture` - game/pet/pet.go
- [x] Line 2548: `sendPetMenu` - game/pet/pet.go
- [x] Line 2566: `sendPetHatch` - game/pet/pet.go
- [x] Line 2577: `sendPetName` - game/pet/pet.go
- [x] Line 2588: `sendPetEmotion` - game/pet/pet.go
- [x] Line 1730: `parse_pet_evolution` - game/pet/pet.go
- [x] Line 1735: `reconstruct_pet_evolution` - game/pet/pet.go
- [x] Line 1740: `sendPetEvolution` - game/pet/pet.go

## Mercenary System
- [x] Line 3113: `sendMercenaryCommand` - game/mercenary/mercenary.go
- [x] Line 2230: `sendCompanionRelease` - game/mercenary/mercenary.go

## Battle System
- [x] Line 1177: `sendShowEquipPlayer` - game/battle/battle.go
- [x] Line 1967: `sendEmotion` - game/battle/battle.go
- [x] Line 1841: `sendNoviceDoriDori` - game/battle/battle.go
- [x] Line 1916: `sendNoviceExplosionSpirits` - game/battle/battle.go
- [x] Line 3530: `sendMemorialDungeonCommand` - game/battle/battle.go

## Marriage System
- [x] Line 1882: `sendAdoptRequest` - game/marriage/marriage.go
- [x] Line 1893: `sendAdoptReply` - game/marriage/marriage.go

## Auction System
- [x] Line 2976: `sendAuctionAddItemCancel` - game/auction/auction.go
- [x] Line 2989: `sendAuctionAddItem` - game/auction/auction.go
- [x] Line 3001: `sendAuctionCreate` - game/auction/auction.go
- [x] Line 3014: `sendAuctionCancel` - game/auction/auction.go
- [x] Line 3025: `sendAuctionBuy` - game/auction/auction.go
- [x] Line 3037: `sendAuctionItemSearch` - game/auction/auction.go
- [x] Line 3060: `sendAuctionReqMyInfo` - game/auction/auction.go
- [x] Line 3071: `sendAuctionMySellStop` - game/auction/auction.go

## Buying Store System
- [x] Line 2600: `sendBuyBulk` - game/buyingstore/buyingstore.go
- [x] Line 2611: `reconstruct_buy_bulk` - game/buyingstore/buyingstore.go
- [x] Line 2618: `sendSellBulk` - game/buyingstore/buyingstore.go
- [x] Line 2629: `reconstruct_sell_bulk` - game/buyingstore/buyingstore.go
- [x] Line 1784: `sendSearchStoreClose` - game/buyingstore/buyingstore.go
- [x] Line 1794: `sendSearchStoreSearch` - game/buyingstore/buyingstore.go
- [x] Line 1809: `reconstruct_search_store_info` - game/buyingstore/buyingstore.go
- [x] Line 1820: `sendSearchStoreRequestNextPage` - game/buyingstore/buyingstore.go
- [x] Line 1828: `sendSearchStoreSelect` - game/buyingstore/buyingstore.go

## UI System
- [x] Line 1198: `sendMiscConfigSet` - game/ui/ui.go
- [x] Line 1309: `sendProgress` - game/ui/ui.go
- [x] Line 1177: `sendShowEquipPlayer` - game/ui/ui.go
- [x] Line 1365: `sendRefineUISelect` - game/ui/ui.go
- [x] Line 1380: `sendRefineUIRefine` - game/ui/ui.go
- [x] Line 1394: `sendRefineUIClose` - game/ui/ui.go
- [x] Line 1338: `sendItemListWindowSelected` - game/ui/ui.go
- [x] Line 1356: `reconstruct_item_list_window_selected` - game/ui/ui.go
- [x] Line 2222: `sendMemo` - game/ui/ui.go
- [x] Line 3278: `sendStylistChange` - game/ui/ui.go
- [x] Line 3297: `sendOpenUIRequest` - game/ui/ui.go
- [x] Line 3315: `sendAttendanceRewardRequest` - game/ui/ui.go
- [x] Line 3368: `sendRouletteWindowOpen` - game/ui/ui.go
- [x] Line 3380: `sendRouletteInfoRequest` - game/ui/ui.go
- [x] Line 3392: `sendRouletteClose` - game/ui/ui.go
- [x] Line 3404: `sendRouletteStart` - game/ui/ui.go
- [x] Line 3416: `sendRouletteClaimPrize` - game/ui/ui.go
- [x] Line 1667: `sendQuestState` - game/ui/ui.go


## Deal System
- [x] Line 1317: `sendDealAddItem` - game/deal/deal.go
- [x] Line 2115: `sendDeal` - game/deal/deal.go
- [x] Line 2126: `sendDealReply` - game/deal/deal.go
- [x] Line 2147: `sendDealFinalize` - game/deal/deal.go
- [x] Line 2155: `sendCurrentDealCancel` - game/deal/deal.go
- [x] Line 2163: `sendDealTrade` - game/deal/deal.go

## Ranking System
- [x] Line 2635: `sendAchievementGetReward` - game/ranking/ranking.go
- [x] Line 2644: `sendTop10Alchemist` - game/ranking/ranking.go
- [x] Line 2656: `sendTop10Blacksmith` - game/ranking/ranking.go
- [x] Line 2668: `sendTop10PK` - game/ranking/ranking.go
- [x] Line 2680: `sendTop10Taekwon` - game/ranking/ranking.go
- [x] Line 2692: `sendTop10` - game/ranking/ranking.go

## GM System
- [x] Line 2707: `sendGMSummon` - game/gm/gm.go
- [x] Line 2727: `sendGMKick` - game/gm/gm.go
- [x] Line 2736: `sendGMKickAll` - game/gm/gm.go
- [x] Line 2742: `sendGMMonsterItem` - game/gm/gm.go
- [x] Line 2751: `sendGMMapMove` - game/gm/gm.go
- [x] Line 2762: `sendGMResetStateSkill` - game/gm/gm.go
- [x] Line 2777: `sendGMChangeMapType` - game/gm/gm.go
- [x] Line 2803: `sendGMChangeEffectState` - game/gm/gm.go
- [x] Line 2814: `sendGMRemove` - game/gm/gm.go
- [x] Line 2823: `sendGMShift` - game/gm/gm.go
- [x] Line 2832: `sendGMRecall` - game/gm/gm.go
- [x] Line 3143: `sendGMGiveMannerByName` - game/gm/gm.go
- [x] Line 3152: `sendGMRequestStatus` - game/gm/gm.go
- [x] Line 3172: `sendGMReqAccName` - game/gm/gm.go
- [x] Line 1924: `sendBanCheck` - game/gm/gm.go

## Achievement System
- [x] Line 1667: `sendQuestState` - game/ui/ui.go (already implemented in UI System)
- [x] Line 2635: `sendAchievementGetReward` - game/ranking/ranking.go (already implemented in Ranking System)

## Macro Detection
- [x] Line 1708: `sendMacroStart` - game/macro/macro.go
- [x] Line 1715: `sendMacroStop` - game/macro/macro.go
- [x] Line 3496: `sendMacroDetectorDownload` - game/macro/macro.go
- [x] Line 3506: `sendMacroDetectorAnswer` - game/macro/macro.go
- [x] Line 1722: `sendReqCashTabCode` - game/macro/macro.go

## Captcha System
- [x] Line 3192: `sendCaptchaAnswer` - game/captcha/captcha.go
- [x] Line 3519: `sendCaptchaPreviewRequest` - game/captcha/captcha.go

## Card System
- [x] Line 2273: `sendCardMergeRequest` - game/card/card.go
- [x] Line 2284: `sendCardMerge` - game/card/card.go

## Cash Shop
- [x] Line 1137: `sendRequestCashItemsList` - game/cashshop/cashshop.go
- [x] Line 1143: `sendCashShopOpen` - game/cashshop/cashshop.go
- [x] Line 1149: `sendCashShopClose` - game/cashshop/cashshop.go
- [x] Line 1156: `sendCashBuy` - game/cashshop/cashshop.go
- [x] Line 1169: `reconstruct_cash_shop_buy` - game/cashshop/cashshop.go (implemented in SendCashBuy)
- [x] Line 3209: `sendCashShopBuy` - game/cashshop/cashshop.go
- [x] Line 3246: `sendMergeItemRequest` - game/cashshop/cashshop.go
- [x] Line 3257: `reconstruct_merge_item_request` - game/cashshop/cashshop.go (implemented in SendMergeItemRequest)
- [x] Line 3267: `sendMergeItemCancel` - game/cashshop/cashshop.go

## Miscellaneous
- [x] Line 1400: `sendTokenToServer` - game/misc/misc.go
- [x] Line 1431: `encrypt_password` - game/misc/misc.go
- [x] Line 1447: `sendReqRemainTime` - game/misc/misc.go
- [x] Line 1457: `sendBlockingPlayerCancel` - game/misc/misc.go
- [x] Line 1692: `sendRecallSso` - game/misc/misc.go
- [x] Line 1700: `sendRemoveAidSso` - game/misc/misc.go
- [x] Line 3161: `sendFeelSaveOk` - game/misc/misc.go
- [x] Line 1097: `sendReplySyncRequestEx` - game/misc/misc.go
- [x] Line 1338: `sendItemListWindowSelected` - game/ui/ui.go (already implemented in UI System)
- [x] Line 1356: `reconstruct_item_list_window_selected` - game/ui/ui.go (already implemented in UI System)