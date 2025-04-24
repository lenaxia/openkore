**Method Implementations:**

- exp() - Experience handler (lines 562-603)
  - Handles ZC_NOTIFY_EXP/ZC_NOTIFY_EXP2 packets
  - Calculates exp percentage based on max exp
  - Differentiates between:
    - Base EXP (VAR_EXP) vs Job EXP (VAR_JOBEXP)
    - Battle EXP (EXP_FROM_BATTLE) vs Quest EXP (EXP_FROM_QUEST)
  - Outputs formatted messages for different exp types
  - Packet formats:
    - ZC_NOTIFY_EXP: <account id>.L <amount>.L <var id>.W <exp type>.W
    - ZC_NOTIFY_EXP2: <account id>.L <amount>.Q <var id>.W <exp type>.W
