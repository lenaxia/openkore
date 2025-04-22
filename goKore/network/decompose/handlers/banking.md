**Banking System Handlers:**

- banking_withdraw - Bank withdrawal result handler (lines 12015-12035)
  - Processes bank withdrawal result notifications (09AA)
  - Handles multiple result codes:
    * BWA_SUCCESS (0x0): Successful withdrawal
      - Updates character zeny balance
      - Displays success message
      - Calls 'banking_withdraw_success' plugin hook
    * BWA_NO_MONEY (0x1): Insufficient funds
      - Displays "No Money for Withdraw" error message
    * BWA_UNKNOWN_ERROR (0x2): Bank overflow
      - Displays "Money in the bank overflow" error message
  - Calls 'banking_withdraw_failed' plugin hook on failure
    * Passes reason code to hook
  - Contains TODO comment about stat_info packet handling

- banking_check - Bank account balance check handler (lines 11968-11988)
  - Processes bank account balance check notifications (09A6)
  - Sets bankingopened flag to 1
  - Stores bank zeny amount in banking.zeny
  - Displays formatted banking information:
    * Header with separator lines
    * Bank balance ("In Bank")
    * Character's current zeny ("On Hand")
  - Uses "info" message category
  - Calls 'banking_opened' plugin hook
  - Handles reason code 1 (mark opening and closing)

- banking_deposit - Bank deposit result handler (lines 11990-12012)
  - Processes bank deposit result notifications (09A8)
  - Handles multiple result codes:
    * BDA_SUCCESS (0x0): Successful deposit
      - Updates character zeny balance
      - Displays success message
      - Calls 'banking_deposit_success' plugin hook
    * BDA_ERROR (0x1): Generic deposit error
      - Displays "Try it again" error message
    * BDA_NO_MONEY (0x2): Insufficient funds
      - Displays "No Money For Deposit" error message
    * BDA_OVERFLOW (0x3): Bank overflow
      - Displays "Money in the bank overflow" error message
  - Calls 'banking_deposit_failed' plugin hook on failure
    * Passes reason code to hook
  - Contains TODO comment about stat_info packet handling