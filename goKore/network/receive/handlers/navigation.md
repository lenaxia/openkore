**Navigation System Handlers:**

- navigate_to - Navigation request handler (lines 12037-12054)
  - Processes navigation requests from server (08E2)
  - Handles two navigation types:
    * Monster hunting: When mob_id is provided
      - Displays message about navigating to map to find monster
    * Location navigation: When x,y coordinates are provided
      - Displays message about navigating to specific coordinates
  - Calls 'navigate_to' plugin hook with all arguments
  - Simple implementation focused on notification
  - Contains TODO comment about documenting type and flag parameters