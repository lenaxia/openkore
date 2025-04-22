**Method Implementations:**

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
