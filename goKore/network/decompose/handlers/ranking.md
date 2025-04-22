# Ranking Related Handlers

**Method Implementations:**
- top10_taekwon_rank - Taekwon top 10 ranking display (lines 10932-10940)
  - Processes Taekwon top 10 ranking notifications
  - Gets formatted text list using top10Listing helper
  - Converts bytes to string
  - Displays formatted message with:
    * "TAEKWON RANK" header
    * Column headers: #, Name, Points
    * Listing of top 10 Taekwon fighters with points
    * Footer separator line
  - Uses "list" message category
  - Contains comment about packet format (0226)
- top10_pk_rank - PK top 10 ranking display (lines 10920-10928)
  - Processes PK top 10 ranking notifications
  - Gets formatted text list using top10Listing helper
  - Converts bytes to string
  - Displays formatted message with:
    * "PVP RANK" header
    * Column headers: #, Name, Points
    * Listing of top 10 PK players with points
    * Footer separator line
  - Uses "list" message category
  - Contains comment about packet format (0238)
- top10_blacksmith_rank - Blacksmith top 10 ranking display (lines 10908-10916)
  - Processes blacksmith top 10 ranking notifications
  - Gets formatted text list using top10Listing helper
  - Converts bytes to string
  - Displays formatted message with:
    * "BLACKSMITH RANK" header
    * Column headers: #, Name, Points
    * Listing of top 10 blacksmiths with points
    * Footer separator line
  - Uses "list" message category
  - Contains comment about packet format (0219)
- top10_alchemist_rank - Alchemist top 10 ranking display (lines 10896-10904)
  - Processes alchemist top 10 ranking notifications
  - Gets formatted text list using top10Listing helper
  - Converts bytes to string
  - Displays formatted message with:
    * "ALCHEMIST RANK" header
    * Column headers: #, Name, Points
    * Listing of top 10 alchemists with points
    * Footer separator line
  - Uses "list" message category
  - Contains comment about packet format (021A)
- top10 - Top 10 ranking dispatcher (lines 10878-10892)
  - Processes top 10 ranking notifications
  - Dispatches to specific handlers based on type:
    * 0: Blacksmith ranking (top10_blacksmith_rank)
    * 1: Alchemist ranking (top10_alchemist_rank)
    * 2: Taekwon ranking (top10_taekwon_rank)
    * 3: PK ranking (top10_pk_rank)
    * Other: Displays "Unknown top10 type" error
  - Passes raw message data (minus first 2 bytes) to specific handlers
  - Simple implementation focused on dispatching
  - Packet: 097D
- alchemist_point - Alchemist ranking points handler (lines 10146-10149)
  - Processes alchemist ranking points notifications
  - Displays message with points increase and total points
  - Uses "list" message category
  - Contains comment about packet format (021C)
  - Similar to blacksmith_points but for alchemist ranking
- blacksmith_points - Blacksmith ranking points handler (lines 10139-10142)
  - Processes blacksmith ranking points notifications
  - Displays message with points increase and total points
  - Uses "list" message category
  - Contains comment about packet format (021B)
  - Simple implementation focused on notification
- rank_points - Ranking points update handler (lines 10128-10135)
  - Processes ranking points update notifications
  - Handles multiple ranking types:
    * 0: Blacksmith - Calls blacksmith_points
    * 1: Alchemist - Calls alchemist_point
    * 2: Taekwon - Calls taekwon_rank with total as rank
    * Other: Displays "Unknown rank type" message
  - Contains comment about packet format (097E)
  - Simple implementation focused on routing to specific handlers