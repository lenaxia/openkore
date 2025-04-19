**Booking System Handlers:**

**Booking Registration:**
- booking_register_request() - Handles booking creation (lines 8834-8845)
  - Processes registration results:
    * Success (0)
    * Already active (2)
    * Unknown errors
  - Provides appropriate feedback

**Booking Search:**
- booking_search_request() - Processes search results (lines 8847-8867)
  - Handles empty results case
  - Displays formatted search results:
    * Character name
    * Index
    * Creation time
    * Level
    * Map ID
    * Job information

**Booking Deletion:**
- booking_delete_request() - Handles deletion (lines 8869-8881)
  - Processes deletion results:
    * Success (0)
    * No active booking (3)
    * Unknown errors
  - Provides appropriate feedback

**Booking Notifications:**
- booking_insert() - Notifies of new bookings (lines 8883-8888)
  - Logs creation of new bookings
  - Includes character name and index

- booking_update() - Notifies of booking changes (lines 8890-8895)
  - Logs settings changes
  - Includes booking index