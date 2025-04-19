**System Handlers:**

- quit_response() - Server quit response (lines 9493-9498)
  - Processes server quit confirmation
  - Handles different response types:
    * Normal quit (0)
    * Force quit (1)
    * Error cases (2-4)
  - Triggers disconnect sequence