# Combat System Design

## Core Types
```go
type CombatStats struct {
    HP, MaxHP       int
    SP, MaxSP       int
    Attack, Defense int
    Hit, Flee       int
    Critical        int
}

type AttackResult struct {
    Damage      int
    IsCritical  bool
    IsMiss      bool
    SkillUsed   SkillID
    Timestamp   time.Time
}

type CombatHistory struct {
    AttacksReceived  []AttackResult
    AttacksDealt     []AttackResult
    LastAttacker     EntityID
    TotalDamageTaken int
    TotalDamageDealt int
}
```

## Combat System Components

```go
// CombatResolver handles all combat calculations
type CombatResolver struct {
    randSource rand.Source // For deterministic testing
}

func (cr *CombatResolver) CalculateDamage(
    attacker Entity, 
    defender Entity, 
    skill Skill,
) AttackResult {
    baseDmg := cr.calculateBaseDamage(attacker, defender, skill)
    
    return AttackResult{
        Damage:     baseDmg,
        IsCritical: cr.checkCritical(attacker, defender),
        IsMiss:     cr.checkMiss(attacker, defender),
        SkillUsed:  skill.ID,
        Timestamp:  time.Now(),
    }
}

func (cr *CombatResolver) calculateBaseDamage(
    attacker Entity, 
    defender Entity,
    skill Skill,
) int {
    attack := attacker.Stats().Attack
    defense := defender.Stats().Defense
    
    // Physical damage formula
    if skill.IsPhysical() {
        return attack - defense/2 + 
               skill.Power * (attack/100)
    }
    // Magical damage formula
    return attack + 
           skill.Power * (attacker.Stats().MATK/50) -
           defender.Stats().MDEF
}

func (cr *CombatResolver) checkCritical(
    attacker Entity,
    defender Entity,
) bool {
    critRate := attacker.Stats().Critical - 
                defender.Stats().Luk/3
    return cr.randSource.Intn(100) < critRate
}

// CombatResult includes detailed breakdown
type CombatResult struct {
    BaseDamage    int
    FinalDamage   int
    DamageType    DamageType
    IsCritical    bool
    IsMiss        bool
    Element       ElementType
    StatusApplied []StatusEffect
}
```

## Combat Events
| Event Type       | Payload Contents               |
|------------------|--------------------------------|
| AttackStarted    | Attacker, Target, Skill       |
| AttackLanded     | Damage, Critical, Miss        |
| DamageTaken      | Amount, DamageType, Attacker  |
| Death            | Killer, LastDamage            |
| StatusApplied    | StatusID, Duration, Caster    |
