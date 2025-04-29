# Advanced Testing for Network::Receive Functions

This document describes how to use the expanded test harness to verify that the Go implementation of Network::Receive matches the Perl implementation.

## Overview

The test harness has been expanded to test all functions in Network::Receive.pm. This ensures that the Go implementation matches the Perl implementation for all network receive functionality.

## Test Files

The test harness consists of the following files:

1. `generate_receive_tests.sh`: Generates test data for all Network::Receive functions
2. `run_receive_tests.sh`: Runs tests for Network::Receive functions and compares results
3. `receive_function_test_harness.pl`: Perl implementation of the test harness
4. `receive_function_test_harness.go`: Go implementation of the test harness

## Test Data

Test data is stored in JSON files in the `test_data/receive_functions` directory. Each function has its own test data file with the following structure:

```json
{
    "function_name": "function_name",
    "args": {
        "server": "ServerType0",
        "arguments": []
    },
    "expected_output": "",
    "validation": {
        "type": "exact_match",
        "description": "Test function_name function",
        "requirements": [
            "Function should return the expected output"
        ]
    }
}
```

## Generating Test Data

To generate test data for all Network::Receive functions:

```bash
./generate_receive_tests.sh
```

This will create test data files for each function in Network::Receive.pm in the `test_data/receive_functions` directory.

## Running Tests

### Running Tests for All Functions

To run tests for all Network::Receive functions:

```bash
./run_receive_tests.sh test-all
```

This will:
1. Run tests for all functions in the Perl implementation
2. Run tests for all functions in the Go implementation
3. Compare the results between the Perl and Go implementations

### Running Tests for a Specific Function

To run tests for a specific function:

```bash
./run_receive_tests.sh run <function_name> <implementation>
```

For example:

```bash
./run_receive_tests.sh run exp perl
./run_receive_tests.sh run exp go
```

### Comparing Results

To compare the results between Perl and Go implementations for a specific function:

```bash
./run_receive_tests.sh compare <function_name>
```

For example:

```bash
./run_receive_tests.sh compare exp
```

## Test Results

Test results are stored in the `results/receive_functions` directory. Each function has two result files:

1. `perl_<function_name>.txt`: Result from the Perl implementation
2. `go_<function_name>.txt`: Result from the Go implementation

## Implementing Go Functions

To implement a Go function that matches the Perl implementation:

1. Look at the Perl implementation in `src/Network/Receive.pm`
2. Implement the equivalent function in Go in `goKore/network/implementation/network/receive`
3. Run the test to verify that the Go implementation matches the Perl implementation

## Troubleshooting

If a test fails, check the following:

1. Make sure the function is implemented in the Go code
2. Check that the function signature matches the Perl implementation
3. Check that the function behavior matches the Perl implementation
4. Look at the error messages in the result files

## Complete List of Functions

The test harness covers all functions in Network::Receive.pm, including:

- exp
- parse
- queryLoginPinCode
- queryAndSaveLoginPinCode
- changeToInGameState
- received_characters_blockSize
- received_characters_unpackString
- received_characters_slots_info
- received_characters
- sync_received_characters
- reconstruct_received_characters
- reconstruct_received_characters_info
- character_creation_successful
- character_creation_failed
- received_characters_info
- parse_account_server_info
- reconstruct_account_server_info
- account_server_info
- connection_refused
- map_loaded
- map_load_error
- stat_info
- stats_added
- stats_info
- stat_info2
- actor_display_compatibility
- actor_display
- actor_died_or_disappeared
- actor_action
- actor_info
- unit_levelup
- parse_minimap_indicator
- account_payment_info
- reconstruct_minimap_indicator
- homunculus_property
- enforce_homun_state
- homunculus_state_handler
- homunculus_info
- minimap_indicator
- parse_npc_image
- reconstruct_npc_image
- npc_image
- local_broadcast
- parse_sage_autospell
- reconstruct_sage_autospell
- sage_autospell
- show_eq
- misc_config
- misc_config_reply
- show_eq_msg_self
- show_script
- skill_post_delay
- skill_post_delaylist
- gospel_buff_aligned
- system_chat
- warp_portal_list
- char_delete2_result
- char_delete2_accept_result
- char_delete2_cancel_result
- arrow_equipped
- inventory_item_added
- inventory_item_removed
- rental_expired
- cart_off
- shop_skill
- shop_sold
- shop_sold_long
- vending_start
- vender_items_list
- revolving_entity
- monster_typechange
- monster_hp_info
- monster_hp_info_tiny
- account_id
- marriage_partner_name
- login_pin_code_request
- login_pin_new_code_result
- actor_status_active
- map_property3
- area_spell
- area_spell_multiple2
- area_spell_multiple3
- sync_request_ex
- cash_shop_list
- cash_shop_open_result
- cash_shop_buy_result
- sprite_change
- progress_bar
- progress_bar_stop
- quest_all_list
- quest_all_mission
- quest_add
- quest_update_mission_hunt
- quest_delete
- quest_active
- parse_npc_chat
- npc_chat
- makable_item_list
- storage_opened
- storage_closed
- storage_items_stackable
- storage_items_nonstackable
- storage_item_added
- storage_item_removed
- cart_items_stackable
- cart_items_nonstackable
- cart_item_added
- cart_item_removed
- cart_info
- cart_add_failed
- inventory_items_stackable
- item_list_start
- item_list_stackable
- item_list_nonstackable
- item_list_end
- login_error
- login_error_game_login_server
- character_deletion_successful
- character_deletion_failed
- character_moves
- character_name
- character_status
- whisper_list
- chat_created
- chat_info
- chat_users
- chat_join_result
- chat_modified
- chat_newowner
- chat_user_join
- chat_user_leave
- chat_removed
- deal_add_other
- deal_begin
- deal_cancelled
- deal_complete
- deal_finalize
- deal_request
- devotion
- egg_list
- emoticon
- errors
- friend_list
- friend_logon
- friend_request
- friend_removed
- friend_response
- homunculus_food
- slave_calcproperty_handler
- EAC_key
- gameguard_grant
- gameguard_request
- guild_allies_enemy_list
- guild_ally_request
- guild_broken
- guild_create_result
- guild_info
- guild_members_list
- guild_invite_result
- guild_location
- guild_leave
- guild_expulsion
- guild_member_online_status
- guild_update_member_position
- guild_members_title_list
- guild_name
- guild_request
- guild_master_member
- guild_alliance
- guild_member_setting_list
- guild_skills_list
- guild_expulsion_list
- guild_member_map_change
- guild_member_add
- guild_notice
- misc_effect
- sound_effect
- identify_list
- identify
- ignore_all_result
- ignore_player_result
- item_used
- married
- item_appeared
- item_exists
- item_disappeared
- item_upgrade
- high_jump
- hp_sp_changed
- map_change
- map_changed
- parse_hat_effect
- hat_effect
- npc_talk
- npc_talk_close
- npc_talk_continue
- npc_talk_number
- npc_talk_responses
- npc_talk_text
- npc_store_begin
- npc_store_info
- npc_sell_list
- npc_clear_dialog
- buy_result
- npc_market_info
- npc_market_purchase_result
- deal_add_you
- skill_exchange_item
- refineui_opened
- refineui_info
- refine_status
- character_ban_list
- flag
- offline_clone_found
- offline_clone_lost
- remain_time_info
- received_login_token
- hotkeys
- received_character_ID_and_Map
- received_sync
- actor_look_at
- actor_movement_interrupted
- actor_trapped
- party_join
- party_allow_invite
- party_chat
- party_exp
- party_leader
- party_hp_info
- party_invite
- party_invite_result
- party_leave
- party_location
- party_organize_result
- party_show_picker
- party_users_info
- party_dead
- rodex_mail_list
- rodex_read_mail
- unread_rodex
- rodex_remove_item
- rodex_add_item
- rodex_open_write
- rodex_check_player
- rodex_write_result
- rodex_get_zeny
- rodex_get_item
- rodex_delete
- booking_register_request
- booking_search_request
- booking_delete_request
- booking_insert
- booking_update
- booking_delete
- clan_user
- clan_info
- clan_chat
- clan_leave
- change_title
- pet_capture_process
- pet_capture_result
- pet_emotion
- pet_evolution_result
- pet_food
- pet_info
- pet_info2
- elemental_info
- upgrade_list
- cooking_list
- refine_result
- upgrade_message
- open_buying_store_fail
- search_store_open
- search_store_fail
- search_store_result
- search_store_pos
- skill_msg
- message_string
- skills_list
- skill_update
- overweight_percent
- partylv_info
- achievement_reward_ack
- achievement_update
- achievement_list
- quit_response
- private_airship_type
- sell_result
- GM_req_acc_name
- boss_map_info
- adopt_reply
- GM_silence
- guild_storage_log
- skill_delete
- captcha_session_ID
- captcha_image
- captcha_answer
- open_buying_store
- buyer_items
- open_buying_store_item_list
- buying_store_found
- buying_store_lost
- buying_store_items_list
- buying_store_item_delete
- buying_store_fail
- buying_store_update
- buyer_found
- buyer_lost
- buying_buy_fail
- special_item_obtain
- inventory_item_favorite
- private_message_sent
- vender_buy_fail
- cash_dealer
- merge_item_open
- parse_merge_item_open
- merge_item_result
- parse_merge_item_result
- memo_success
- change_to_constate25
- adopt_request
- rank_points
- blacksmith_points
- alchemist_point
- area_spell_disappears
- arrow_none
- arrowcraft_list
- attack_range
- auction_my_sell_stop
- auction_windows
- auction_add_item
- premium_rates_info
- rates_info2
- auction_result
- battleground_message
- battleground_emblem
- guild_emblem
- guild_emblem_update
- char_emblem_update
- guild_position_changed
- guild_position
- guild_unally
- guild_opposition_result
- guild_alliance_added
- map_change_cell
- blade_stop
- divorced
- hack_shield_alarm
- talkie_box
- manner_message
- instance_window_start
- instance_window_queue
- instance_window_join
- instance_window_leave
- card_merge_list
- card_merge_status
- combo_delay
- book_read
- rental_time
- cash_buy_fail
- equip_item
- equip_item_switch
- equip_switch_run_res
- equip_switch_log
- font
- initialize_message_id_encryption
- mail_delete
- mail_window
- mail_return
- mail_read
- mail_refreshinbox
- mail_getattachment
- mail_setattachment
- mail_send
- mail_new
- top10
- top10_alchemist_rank
- top10_blacksmith_rank
- top10_pk_rank
- top10_taekwon_rank
- taekwon_packets
- taekwon_rank
- storage_password_request
- storage_password_result
- mercenary_init
- mercenary_off
- monster_ranged_attack
- mvp_item
- mvp_other
- mvp_you
- no_teleport
- private_message
- progress_bar_unit
- pvp_rank
- repair_list
- repair_result
- resurrection
- secure_login_key
- self_chat
- sync_request
- sense_result
- skill_cast
- cast_cancelled
- switch_character
- unequip_item
- unequip_item_switch
- use_item
- users_online
- vender_found
- vender_lost
- skill_add
- isvr_disconnect
- skill_use_failed
- open_store_status
- stylist_res
- open_ui
- action_ui
- attendance_ui
- move_interrupt
- banking_check
- banking_deposit
- banking_withdraw
- navigate_to
- roulette_window
- roulette_info
- roulette_recv_item
- roulette_window_update
- load_confirm
- inventory_expansion_result
- item_preview
- ping
- starplace
- captcha_upload_request
- captcha_upload_request_status
- macro_reporter_status
- macro_detector
- macro_detector_image
- macro_detector_show
- macro_detector_status
- captcha_preview
- captcha_preview_image
- macro_reporter_select
- repute_info
- gold_pc_cafe_point
- dynamicnpc_create_result