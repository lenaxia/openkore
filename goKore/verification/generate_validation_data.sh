#!/bin/bash
# Script to generate validation data for common Send.pm methods

# Make sure the send_wrapper.pl script is executable
chmod +x send_wrapper.pl

# Create output directory
mkdir -p validation_data

echo "Generating validation data for common Send.pm methods..."

# Login-related packets
echo "Generating login packets..."
./send_wrapper.pl call sendMasterLogin "username" "password" 1 1
mv packet_output_sendMasterLogin.json validation_data/

./send_wrapper.pl call sendGameLogin 123456 "sessionID" "sessionID2" 1
mv packet_output_sendGameLogin.json validation_data/

./send_wrapper.pl call sendCharLogin 0
mv packet_output_sendCharLogin.json validation_data/

./send_wrapper.pl call sendMapLogin 123456 789012 "sessionID" 1
mv packet_output_sendMapLogin.json validation_data/

# Movement and actions
echo "Generating movement and action packets..."
./send_wrapper.pl call sendMove 100 100
mv packet_output_sendMove.json validation_data/

./send_wrapper.pl call sendLook 0 0
mv packet_output_sendLook.json validation_data/

./send_wrapper.pl call sendAction 12345 0
mv packet_output_sendAction.json validation_data/

# Chat
echo "Generating chat packets..."
./send_wrapper.pl call sendChat "Hello, world!"
mv packet_output_sendChat.json validation_data/

./send_wrapper.pl call sendPrivateMsg "Player" "Hello there!"
mv packet_output_sendPrivateMsg.json validation_data/

./send_wrapper.pl call sendPartyChat "Party message"
mv packet_output_sendPartyChat.json validation_data/

# Items
echo "Generating item-related packets..."
./send_wrapper.pl call sendItemUse 1234 5678
mv packet_output_sendItemUse.json validation_data/

./send_wrapper.pl call sendDrop 1234 10
mv packet_output_sendDrop.json validation_data/

./send_wrapper.pl call sendTake 1234
mv packet_output_sendTake.json validation_data/

# Skills
echo "Generating skill-related packets..."
./send_wrapper.pl call sendSkillUse 1 10 12345
mv packet_output_sendSkillUse.json validation_data/

./send_wrapper.pl call sendSkillUseLoc 1 10 100 100
mv packet_output_sendSkillUseLoc.json validation_data/

# Other common packets
echo "Generating other common packets..."
./send_wrapper.pl call sendSync
mv packet_output_sendSync.json validation_data/

./send_wrapper.pl call sendRestart 0
mv packet_output_sendRestart.json validation_data/

echo "Validation data generation complete!"
echo "JSON files are available in the validation_data directory."