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
