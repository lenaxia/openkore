# CAPTCHA Related Handlers

**Additional Macro Detector Handlers:**
- macro_detector_status - Macro detector status handler (lines 12316-12331)
  - Processes macro detector status notifications (0A5D)
  - Translates status codes to readable messages:
    * MCD_TIMEOUT: "Timeout"
    * MCD_INCORRECT: "Incorrect"
    * MCD_GOOD: "Correct"
    * Other: "Unknown"
  - Displays status message
  - Uses "captcha" message category
  - Packet: PACKET_ZC_MACRO_DETECTOR_STATUS

- macro_reporter_select - Macro reporter player list handler (lines 12374-12385)
  - Processes macro reporter player list notifications (0A6D)
  - Displays header message for account list
  - Parses account IDs from binary data in 4-byte chunks
  - Looks up player names from playersList using account IDs
  - Displays list of player names
  - Packet: PACKET_ZC_MACRO_REPORTER_SELECT

**Captcha Preview System:**
- captcha_preview - Captcha preview request handler (lines 12333-12349)
  - Processes captcha preview request notifications (0A6A)
  - Stores image size and captcha key in global variables
  - Handles multiple status codes:
    * 0: Ready to download - "Now you can download the image"
    * 1: Request failed - "Failed to Request Captcha (ID is out of range)"
    * Other: Unknown status with code
  - Outputs debug message with image size and captcha key
  - Uses "captcha" debug category
  - Packet: PACKET_ZC_CAPTCHA_PREVIEW_REQUEST

- captcha_preview_image - Captcha preview image handler (lines 12351-12372)
  - Processes captcha preview image chunks (0A6B)
  - Accumulates image data until complete
  - When complete:
    * Uncompresses image data
    * Converts to hex format
    * Saves to file in logs folder with character name and key
    * Displays message with file location
    * Resets captcha state variables
  - Uses "captcha" message category
  - Packet: PACKET_ZC_CAPTCHA_PREVIEW_REQUEST_DOWNLOAD

**Macro Detector System:**
- macro_detector_show - Macro detector UI handler (lines 12306-12312)
  - Processes macro detector UI notifications (0A5B)
  - Displays macro detector header message
  - Shows remaining chances and time information
  - Uses "captcha" message category
  - Skips further processing for non-direct connections
  - Packet: PACKET_ZC_MACRO_DETECTOR_SHOW

- macro_detector_image - Macro detector captcha image handler (lines 12262-12304)
  - Processes macro detector captcha image chunks (0A59)
  - Accumulates image data until complete
  - When complete:
    * Uncompresses image data
    * Processes image (fixes specific byte patterns)
    * Saves to file in logs folder
    * Triggers captcha_image and captcha_file hooks
    * Returns early if hook sets return flag
    * Displays instructions for solving captcha
    * Resets captcha state variables
    * Sends download confirmation to server
  - Uses "captcha" message category
  - Packet: PACKET_ZC_MACRO_DETECTOR_REQUEST_DOWNLOAD

- macro_detector - Macro detector image info handler (lines 12253-12260)
  - Processes macro detector image info notifications (0A58)
  - Outputs debug message with image size and captcha key
  - Stores image size and captcha key in global variables
  - Uses "captcha" debug category
  - Packet: PACKET_ZC_MACRO_DETECTOR_REQUEST

- macro_reporter_status - Macro reporter status handler (lines 12236-12251)
  - Processes macro reporter status notifications (0A57)
  - Translates status codes to readable messages:
    * MCR_MONITORING: "Monitoring"
    * MCR_NO_DATA: "No Data"
    * MCR_INPROGRESS: "In Progress"
    * Other: "Unknown"
  - Displays status message
  - Uses "captcha" message category
  - Packet: PACKET_ZC_MACRO_REPORTER_STATUS

**Captcha Upload System:**
- captcha_upload_request - Captcha upload request handler (lines 12215-12228)
  - Processes captcha upload request notifications (0A53)
  - Handles multiple status codes:
    * 0: Ready to upload - "Now you can upload the image"
    * 1: Upload failed - "Failed to upload the image"
    * Other: Unknown status with code
  - Skips further processing for non-direct connections
  - Packet: PACKET_ZC_CAPTCHA_UPLOAD_REQUEST

- captcha_upload_request_status - Captcha upload result handler (lines 12230-12234)
  - Processes captcha upload result notifications (0A55)
  - Displays success message
  - Simple implementation with minimal functionality
  - Packet: PACKET_ZC_CAPTCHA_UPLOAD_REQUEST_STATUS

**Method Implementations:**
- captcha_answer - CAPTCHA answer result handler (lines 9658-9665)
  - Processes CAPTCHA answer result notifications
  - Outputs debug message with packet information
  - Outputs debug message with answer result ("good" or "bad")
  - Updates captcha_state global variable with flag value
  - Triggers captcha_answer hook with flag value
  - Contains TODO comment: "debug + remove debug message"
  - Packet: 0x07e9,5
- captcha_image - CAPTCHA image handler (lines 9636-9654)
  - Processes CAPTCHA image data
  - Outputs debug message with packet information
  - Triggers captcha_image hook with image data
  - Returns early if hook sets return flag
  - Saves image to captcha.bmp in logs folder
  - Triggers captcha_file hook with file path
  - Returns early if hook sets return flag
  - Displays warning message with instructions for solving captcha
  - Uses file I/O operations to save image data
  - Uses "warning" message category
- captcha_session_ID - CAPTCHA session ID handler (lines 9629-9632)
  - Processes CAPTCHA session ID notifications
  - Outputs debug message with packet information
  - Contains comment indicating it's from kRO::RagexeRE_2009_09_22a
  - Packet: 07E6 (tentative)
  - Simple implementation focused on debugging