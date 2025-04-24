# Mercenary Related Handlers

**Method Implementations:**
- mercenary_off - Mercenary removal handler (lines 11131-11137)
  - Processes mercenary removal notifications
  - Removes mercenary from slave manager
  - Removes mercenary from slaves list
  - Deletes mercenary from character
  - Contains commented out line about deleting from slaves hash
  - Contains +message_string and -message_string comments
  - Simple implementation focused on cleanup
- mercenary_init - Mercenary initialization handler (lines 11093-11128)
  - Processes mercenary initialization notifications
  - Gets or creates mercenary actor
  - Sets mercenary map to current field
  - Copies all parameters from packet to mercenary object
  - Converts name from bytes to string
  - Calls slave_calcproperty_handler for additional property calculations
  - Adds mercenary to slave manager if not already present
  - Handles class conversion if needed (AI::Slave::Mercenary to Actor::Slave::Mercenary)
  - Handles attack distance configuration:
    * If mercenary_attackDistanceAuto is enabled:
      - Adjusts mercenary_attackDistance if greater than attack_range
      - Sets mercenary_attackMaxDistance to match attack_range
      - Displays success messages for autodetected values
  - Contains comment about ST0's counterpart for ST kRO
  - Contains TODO comment about consolidating attack range handling
  - Contains detailed packet format comment (029B)