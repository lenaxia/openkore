# Security Handlers

## PIN Code Management

### Server PIN Requests
- login_pin_code_request (lines 4118-4196)
  - Handles all server PIN code states:
    - 0: Correct PIN
    - 1: PIN requested (already defined)
    - 2: PIN requested (not defined)
    - 3: PIN expired
    - 4: PIN requested (private servers)
    - 5: Invalid PIN
    - 7: PIN disabled
    - 8: Incorrect PIN
  - Processes ZC_PC_PIN_CODE_REQUEST
  - Validates PIN format (4 digits)
  - Shows appropriate messages for each state
  - Calls queryAndSaveLoginPinCode when needed

- login_pin_new_code_result (lines 4198-4218)
  - Handles new PIN code validation
  - Flag 2: Invalid PIN (sequences/repeats)
  - Enforces 4-digit numeric requirement
  - Shows error messages for invalid PINs
  - Prevents problematic PINs that could lock accounts

### Client PIN Handling
### queryLoginPinCode() - lines 632-655
```perl
sub queryLoginPinCode {
    my $message = $_[0] || T("You've never set a login PIN code before.\nPlease enter a new login PIN code:");
    do {
        my $input = $interface->query($message, isPassword => 1,);
        if (!defined($input)) {
            quit();
            return;
        } else {
            if ($input !~ /^\d+$/) {
                $interface->errorDialog(T("The PIN code may only contain digits."));
            } elsif ((length($input) <= 3) || (length($input) >= 9)) {
                $interface->errorDialog(T("The PIN code must be between 4 and 9 characters."));
            } else {
                return $input;
            }
        }
    } while (1);
}
```

### queryAndSaveLoginPinCode() - lines 663-671
```perl
sub queryAndSaveLoginPinCode {
    my ($self, $message) = @_;
    my $pin = queryLoginPinCode($message);
    if (defined $pin) {
        configModify('loginPinCode', $pin, silent => 1);
        return 1;
    } else {
        return 0;
    }
}
```

### Key Features:
- Validates PIN code format (digits only)
- Enforces length requirements (4-8 characters)
- Handles cancellation gracefully
- Provides user feedback via error dialogs
- Configures silent saving of valid PIN codes