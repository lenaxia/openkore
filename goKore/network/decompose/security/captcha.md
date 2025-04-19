**Captcha/Macro Detection Handlers:**

- macro_reporter_select() - Macro reporter account list handler (lines 12374-12385)
  - Displays list of accounts being monitored
  * Shows account IDs and player names
  * Processes PACKET_ZC_MACRO_REPORTER_SELECT

- captcha_upload_request() - Captcha upload request handler (lines 12215-12228)
  - Handles captcha image upload requests
  * Status 0: Can upload image
  * Status 1: Upload failed
  * Only active in DirectConnection mode

- captcha_upload_request_status() - Upload status handler (lines 12230-12234)
  - Confirms successful image upload

- macro_reporter_status() - Macro reporter status handler (lines 12236-12251)
  - Tracks macro detection status:
    * MCR_MONITORING (Monitoring)
    * MCR_NO_DATA (No Data)
    * MCR_INPROGRESS (In Progress)

- macro_detector() - Macro detector info handler (lines 12253-12260)
  - Receives captcha image metadata
  * Stores image size and captcha key

- macro_detector_image() - Macro detector image handler (lines 12262-12304)
  - Receives captcha image in chunks
  * Processes and saves image to logs folder
  * Calls hooks for external processing
  * Prompts user to solve captcha

- macro_detector_show() - Macro detector UI handler (lines 12306-12314)
  - Displays remaining chances and time
  * Only active in DirectConnection mode

- macro_detector_status() - Macro detector status handler (lines 12316-12331)
  - Reports detection results:
    * MCD_TIMEOUT (Timeout)
    * MCD_INCORRECT (Incorrect)
    * MCD_GOOD (Correct)

- captcha_preview() - Captcha preview request handler (lines 12333-12349)
  - Handles preview requests
  * Status 0: Can download image
  * Status 1: Invalid ID
  * Stores image metadata

- captcha_preview_image() - Captcha preview image handler (lines 12351-12372)
  - Receives preview image in chunks
  * Saves image to logs folder
  * Includes character name in filename