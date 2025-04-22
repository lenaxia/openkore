# MVP Related Handlers

**Method Implementations:**
- mvp_you - Player MVP notification handler (lines 11179-11184)
  - Processes player MVP notifications
  - Creates message with congratulations and experience amount
  - Displays message
  - Logs to chat log with "k" category
  - Simple implementation focused on notification
- mvp_other - Other player MVP notification handler (lines 11172-11177)
  - Processes other player MVP notifications
  - Gets actor name using Actor::get function
  - Displays message with actor name
  - Logs to chat log with "k" category
  - Simple implementation focused on notification
- mvp_item - MVP item reward handler (lines 11165-11170)
  - Processes MVP item reward notifications
  - Gets item name using itemNameSimple function
  - Displays message with item name
  - Logs to chat log with "k" category
  - Simple implementation focused on notification