# Account Rates Related Handlers

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