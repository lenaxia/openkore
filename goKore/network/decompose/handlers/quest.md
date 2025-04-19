**Quest System Handlers:**
- quest_list (lines 6332-6345)
  - Manages active quest listings
  - Processes:
    - Quest IDs
    - Objectives
    - Progress tracking
    - Reward information
  - Maintains:
    - @activeQuests array
    - %questDetails hash

- quest_update (lines 6347-6360)
  - Handles quest progress updates

- quest_complete (lines 6362-6375)
  - Processes quest completion
  - Handles:
    - Reward distribution
    - Quest log updates
    - Completion status
  - Features:
    - Triggers 'quest_completed' hook
    - Updates character stats
    - Provides completion feedback
  - Features:
    - Processes objective completion
    - Updates quest status
    - Triggers reward notifications
    - Executes 'quest_update' hooks