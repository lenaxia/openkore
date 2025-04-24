# Elemental Related Handlers

**Method Implementations:**
- elemental_info - Elemental information handler (lines 9100-9111)
  - Processes elemental information notifications
  - Gets or creates elemental actor:
    * Retrieves existing actor if ID changed
    * Creates new Actor::Elemental if not defined
  - Updates elemental data with all fields from packet
  - Uses foreach loop to process multiple fields
  - Simple implementation focused on data updates