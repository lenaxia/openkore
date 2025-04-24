**Booking System Handlers:**

- booking_delete - Booking deletion notification handler (lines 8898-8902)
  - Processes booking deletion notifications
  - Displays message with booking ID
  - Simple implementation focused on notification
  - Packet: 0x80B

- booking_update - Booking update notification handler (lines 8891-8895)
  - Processes booking update notifications
  - Displays message with booking ID
  - Simple implementation focused on notification
  - Packet: 0x80A

- booking_insert - Booking creation notification handler (lines 8884-8888)
  - Processes booking creation notifications
  - Displays message with creator name and booking ID
  - Simple implementation focused on notification

- booking_delete_request - Booking deletion result handler (lines 8870-8879)
  - Processes booking deletion result notifications
  - Handles multiple result codes:
    * 0: Success - "Reserve deleted successfully!"
    * 3: Not active - "You're not with a group booking active!"
    * Other: Unknown error with code
  - Uses "booking" message category
  - Simple implementation focused on result notification

- booking_search_request - Booking search result handler (lines 8848-8867)
  - Processes booking search result notifications
  - Handles empty result case:
    * Displays "Without results!" error
    * Returns early
  - For successful search:
    * Displays header with separator
    * Processes each 48-byte entry in innerData
    * Unpacks entry data: index, charName, expireTime, level, mapID, job array
    * Formats and displays each entry with:
      - Character name and index
      - Creation time (formatted date) and level
      - Map ID
      - Job information (5 values)
    * Uses separator lines between entries
  - Uses "booking" message category
  - Uses swrite for formatted output

- booking_register_request - Booking registration result handler (lines 8834-8845)
  - Processes booking registration result notifications
  - Handles multiple result codes:
    * 0: Success - "Booking successfully created!"
    * 2: Already active - "You already got a reservation group active!"
    * Other: Unknown error with code
  - Uses "booking" message category
  - Simple implementation focused on result notification