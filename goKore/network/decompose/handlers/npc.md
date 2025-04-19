**NPC Interaction Handlers:**

- npc_talk (lines 7408-7451)
  - Handles NPC dialog messages (ZC_SAY_DIALOG)
  - Features:
    - Auto-creates TalkNPC task if needed
    - Processes NPC ID and message
    - Removes RO color codes
    - Maintains conversation state
    - Triggers npc_talk hook

- npc_talk_close (lines 7453-7474)
  - Handles NPC dialog close button (ZC_CLOSE_DIALOG)
  - Features:
    - Validates NPC ID
    - Clears talk state
    - Triggers npc_talk_done hook

- npc_talk_continue (lines 7476-7485)
  - Handles NPC dialog continue button (ZC_WAIT_DIALOG)
  - Features:
    - Updates talk state to 'next'
    - Maintains timing information

- npc_talk_number (lines 7487-7497)
  - Handles NPC number input dialog (ZC_OPEN_EDITDLG)
  - Features:
    - Updates talk state to 'number'
    - Maintains timing information

- npc_talk_responses (lines 7499-7557)
  - Handles NPC menu selections (ZC_MENU_LIST)
  - Features:
    - Auto-creates TalkNPC task if needed
    - Processes menu items split by ':'
    - Removes RO color codes
    - Handles special itemID format
    - Adds "Cancel Chat" option
    - Triggers npc_talk_responses hook

- npc_talk_text (lines 7559-7569)
  - Handles NPC text input dialog (ZC_OPEN_EDITDLGSTR)
  - Features:
    - Updates talk state to 'text'
    - Maintains timing information

- npc_store_begin (lines 7571-7581)
  - Initiates NPC shop buy/sell dialog (ZC_SELECT_DEALTYPE)
  - Features:
    - Clears previous talk state
    - Sets up initial shop state
    - Gets NPC name

- npc_store_info (lines 7583-7633)
  - Handles NPC shop item list (ZC_PC_PURCHASE_ITEMLIST)
  - Features:
    - Processes different packet formats
    - Creates item entries with price/type info
    - Handles duplicate items
    - Updates store state
    - Triggers 'store' command if not in buyAuto

- npc_sell_list (lines 7635-7663)
  - Handles NPC sellable items list (ZC_PC_SELL_ITEMLIST)
  - Features:
    - Marks items as sellable/unsellable
    - Displays sellable items info
    - Updates talk state to 'sell'

- buy_result (lines 7671-7697)
  - Handles purchase result notifications (ZC_PC_PURCHASE_RESULT)
  - Features:
    - Processes all possible result codes
    - Provides appropriate success/error messages
    - Handles cases:
      - Success
      - Insufficient zeny
      - Overweight
      - Too many items
      - Invalid items
      - Invalid store