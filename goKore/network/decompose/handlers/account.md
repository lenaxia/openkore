**Account Handlers:**

- account_payment_info() - Displays account payment information (lines 2867-2885)
  - Processes D_minute (pay-per-day) and H_minute (pay-per-hour) inputs
  - Converts minutes to days/hours/minutes format:
    * Days = minutes / 1440
    * Hours = (minutes % 1440) / 60
    * Minutes = (minutes % 1440) % 60
  - Formats a bordered message showing both payment schemes
  - Outputs via message() with "info" type