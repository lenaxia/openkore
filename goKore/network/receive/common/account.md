**Account Server Info Handlers:**

- **parse_account_server_info** - Server info packet parser (lines 1048-1068)
  - Handles different server info packet formats for various server types:
    - tRO/twRO (packet 0B60):
      - Length: 164 bytes
      - Type specifiers: 'a4 v Z20 v3 a128 V'
      - Keys: [ip, port, name, state, users, property, ip_port, unknown]
    - kRO Zero/ST/vRO (packets 0AC4/0B07):
      - Length: 160 bytes  
      - Type specifiers: 'a4 v Z20 v3 a128'
      - Keys: [ip, port, name, users, state, property, ip_port]
    - cRO (packet 0AC9):
      - Length: 154 bytes
      - (Additional details to be filled when more lines are read)
  - Parameters:
    - args->{switch}: Packet ID determining format
  - Returns structured server info hashref
  - Related methods:
    - reconstruct_account_server_info()
    - account_server_info()
  
- **reconstruct_account_server_info** - Server info packet reconstructor (lines 1109-1149)
  - Reconstructs server info packets for different server types:
    - tRO 2020 (packet 0B60):
      - Length: 164 bytes
      - Type specifiers: 'a4 v Z20 v3 a128 V'
      - Keys: [ip, port, name, state, users, property, ip_port, unknown]
    - kRO Zero/ST/vRO (packets 0AC4/0B07):
      - Length: 160 bytes
      - Type specifiers: 'a4 v Z20 v3 a128'
      - Keys: [ip, port, name, users, state, property, ip_port]
    - cRO (packet 0AC9):
      - Length: 154 bytes
      - Type specifiers: 'a20 V a2 a126'
      - Keys: [name, users, unknown, ip_port]
    - tRO (packet 0276):
      - Length: 36 bytes
      - Type specifiers: 'a4 v Z20 v5'
      - Keys: [ip, port, name, state, users, property, sid, unknown]
    - Default format:
      - Length: 32 bytes
      - Type specifiers: 'a4 v Z20 v2 x2'
      - Keys: [ip, port, name, users, display]
  - Parameters:
    - args->{switch}: Packet ID determining format
    - args->{lastLoginIP}: IP address to convert to binary format
  - Converts IP addresses to binary format
  - Related methods:
    - parse_account_server_info()
    - account_server_info()
  
- **account_server_info** - Account server info handler (lines 1156-1206)
  - Processes and displays account information from login server:
    - Sets network state to 2 (connected)
    - Stores session/account IDs:
      - sessionID, accountID, sessionID2
    - Processes account sex (mod 2 to handle inRO's female=2)
    - Displays formatted account info:
      - Account ID (numeric and hex)
      - Sex (converted to text)
      - Session IDs (numeric and hex)
    - Processes server list:
      - Displays server name, users, IP, port, SID, state
      - States: Idle, Normal, Busy, Full
    - Handles server selection when needed
  - Parameters:
    - args->{sessionID}: Primary session ID
    - args->{accountID}: Account ID
    - args->{sessionID2}: Secondary session ID
    - args->{accountSex}: Account sex (0=female, 1=male)
    - args->{servers}: Array of server info hashes
  - Related methods:
    - parse_account_server_info()
    - reconstruct_account_server_info()


**Method Implementations:**

- login_pin_new_code_result() - PIN code validation handler (lines 4198-4218)
  - Processes PIN code validation results
  - Handles invalid PIN codes (flag == 2)
  - Prevents use of sequences or repeated numbers
  - Enforces PIN code format (exactly 4 numbers)
  - Contains special handling for bRO bug with non-numeric PINs
  - Warns about potential ban risk from invalid PIN formats

- login_pin_code_request() - PIN code request handler (lines 4118-4196)
  - Processes PIN code request packets from server
  - Handles multiple flag states:
    - 0: Correct PIN
    - 1: PIN requested (already defined)
    - 2/4: PIN requested (not defined)
    - 3: PIN expired
    - 5: PIN invalid (sequence/repeated numbers)
    - 7: PIN disabled
    - 8: PIN incorrect
  - Manages PIN code validation and error handling
  - Handles special cases for XKore modes
  - Integrates with character selection timing
  - Uses queryAndSaveLoginPinCode() for user input

- queryLoginPinCode() - PIN code handler (lines 638-655)
  - Requests login PIN code from user via interface
  - Validates input:
    - Must contain only digits
    - Must be between 4-8 characters long
  - Handles cancellation by quitting
  - Shows error dialogs for invalid input
  - Returns:

- queryAndSaveLoginPinCode() - Persistent PIN handler (lines 662-672)
  - Wrapper around queryLoginPinCode()
  - Saves valid PIN codes to config file
  - Handles PIN code reuse:
    - Skips prompt if saved PIN exists
    - Validates saved PIN matches requirements
  - Returns:
    - Saved PIN if valid and exists
    - New PIN if successfully entered
    - undef if cancelled
    - Valid PIN code string on success
    - undef if cancelled

**Account Handlers:**

- account_id() - Account ID debug handler (lines 4100-4106)
  - Processes account ID packets
  - Logs account ID in both decimal and hexadecimal formats
  - Note: Preserves original accountID format to avoid corruption

- account_payment_info() - Displays account payment information (lines 2867-2885)
  - Processes D_minute (pay-per-day) and H_minute (pay-per-hour) inputs
  - Converts minutes to days/hours/minutes format:
    * Days = minutes / 1440
    * Hours = (minutes % 1440) / 60
    * Minutes = (minutes % 1440) % 60
  - Formats a bordered message showing both payment schemes
  - Outputs via message() with "info" type# Account Rates Related Handlers

**Method Implementations:**
- rates_info2 - Detailed rates information handler (lines 10285-10323)
  - Processes detailed rates information notifications
  - Gets raw message and size
  - Sets up header and detail pack formats
  - Calculates header and detail lengths
  - Initializes rates hash with:
    * exp: Experience rate (divided by 1000)
    * death: Death penalty rate (divided by 1000)
    * drop: Drop rate (divided by 1000)
  - Processes each detail entry in the message:
    * Unpacks type, exp, death, drop values
    * Adds values to rates hash by type
  - Contains comments about detail types:
    * 0: Base server rate
    * 1: Premium account additional rate
    * 2: Server additional rate
    * 3: Extra event rate (possibly)
  - Displays formatted header separator
  - Displays detailed rate information:
    * EXP rates with breakdown
    * Drop rates with breakdown
    * Death penalty with breakdown
  - Displays footer separator
  - Uses "info" message category for all messages
  - Contains comments about packet formats:
    * 08CB: Original format
    * 097B: Updated format
    * 0981: Chinese server format
- premium_rates_info - Premium account rates handler (lines 10274-10277)
  - Processes premium account rates notifications
  - Displays message with:
    * Experience rate bonus (percentage)
    * Death penalty reduction (percentage)
    * Drop rate bonus (percentage)
  - Uses "info" message category
  - Simple implementation focused on notification
