package actor

import (
	"testing"
	"time"
)

func TestBaseActor(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	actor := NewBaseActor(id, "TestActor")

	// Test ID
	if string(actor.ID()) != string(id) {
		t.Errorf("Expected ID %v, got %v", id, actor.ID())
	}

	// Test NameID
	expectedNameID := uint32(0x04030201)
	if actor.NameID() != expectedNameID {
		t.Errorf("Expected NameID %d, got %d", expectedNameID, actor.NameID())
	}

	// Test Name
	actor.SetName("TestName")
	if actor.Name() != "TestName" {
		t.Errorf("Expected Name TestName, got %s", actor.Name())
	}

	// Test Position
	pos := &Position{X: 10, Y: 20}
	actor.SetPosition(pos)
	if actor.Position().X != 10 || actor.Position().Y != 20 {
		t.Errorf("Expected Position {10, 20}, got {%d, %d}", actor.Position().X, actor.Position().Y)
	}

	// Test PositionTo
	posTo := &Position{X: 30, Y: 40}
	actor.SetPositionTo(posTo)
	if actor.PositionTo().X != 30 || actor.PositionTo().Y != 40 {
		t.Errorf("Expected PositionTo {30, 40}, got {%d, %d}", actor.PositionTo().X, actor.PositionTo().Y)
	}

	// Test ActorType
	if actor.ActorType() != "TestActor" {
		t.Errorf("Expected ActorType TestActor, got %s", actor.ActorType())
	}

	// Test Avoid
	actor.SetAvoid(true)
	if !actor.IsAvoid() {
		t.Errorf("Expected IsAvoid true, got false")
	}

	// Test WalkSpeed
	actor.SetWalkSpeed(0.2)
	if actor.WalkSpeed() != 0.2 {
		t.Errorf("Expected WalkSpeed 0.2, got %f", actor.WalkSpeed())
	}

	// Test TimeMove
	now := time.Now()
	actor.SetTimeMove(now)
	if actor.TimeMove() != now {
		t.Errorf("Expected TimeMove %v, got %v", now, actor.TimeMove())
	}

	// Test TimeMoveCalc
	actor.SetTimeMoveCalc(1.5)
	if actor.TimeMoveCalc() != 1.5 {
		t.Errorf("Expected TimeMoveCalc 1.5, got %f", actor.TimeMoveCalc())
	}

	// Test DeepCopy
	copy := actor.DeepCopy().(*BaseActor)
	if string(copy.ID()) != string(id) {
		t.Errorf("DeepCopy: Expected ID %v, got %v", id, copy.ID())
	}
	if copy.Name() != "TestName" {
		t.Errorf("DeepCopy: Expected Name TestName, got %s", copy.Name())
	}
	if copy.Position().X != 10 || copy.Position().Y != 20 {
		t.Errorf("DeepCopy: Expected Position {10, 20}, got {%d, %d}", copy.Position().X, copy.Position().Y)
	}
}

func TestPlayer(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	player := NewPlayer(id)

	// Test basic Actor methods
	player.SetName("TestPlayer")
	if player.Name() != "TestPlayer" {
		t.Errorf("Expected Name TestPlayer, got %s", player.Name())
	}

	// Test Player-specific methods
	player.SetJob(10)
	if player.Job() != "Job 10" {
		t.Errorf("Expected Job 'Job 10', got %s", player.Job())
	}

	player.SetLevel(50)
	if player.Level() != 50 {
		t.Errorf("Expected Level 50, got %d", player.Level())
	}

	player.SetSex(1)
	if player.Sex() != 1 {
		t.Errorf("Expected Sex 1, got %d", player.Sex())
	}

	player.SetHP(800)
	player.SetMaxHP(1000)
	if player.HP() != 800 || player.MaxHP() != 1000 {
		t.Errorf("Expected HP 800/1000, got %d/%d", player.HP(), player.MaxHP())
	}
	if player.HPPercent() != 80 {
		t.Errorf("Expected HPPercent 80, got %d", player.HPPercent())
	}

	player.SetSP(300)
	player.SetMaxSP(500)
	if player.SP() != 300 || player.MaxSP() != 500 {
		t.Errorf("Expected SP 300/500, got %d/%d", player.SP(), player.MaxSP())
	}
	if player.SPPercent() != 60 {
		t.Errorf("Expected SPPercent 60, got %d", player.SPPercent())
	}

	player.SetSitting(true)
	if !player.IsSitting() {
		t.Errorf("Expected IsSitting true, got false")
	}

	player.SetDead(true)
	if !player.IsDead() {
		t.Errorf("Expected IsDead true, got false")
	}

	player.SetGuildID(100)
	player.SetGuildName("TestGuild")
	player.SetGuildTitle("TestTitle")
	player.SetEmblemID(200)
	player.SetPartyName("TestParty")

	if player.GuildID() != 100 || player.GuildName() != "TestGuild" || player.GuildTitle() != "TestTitle" ||
		player.EmblemID() != 200 || player.PartyName() != "TestParty" {
		t.Errorf("Guild/Party info not set correctly")
	}

	player.SetAppearance(5, 6, 7)
	player.SetHeadgear(8, 9, 10)
	player.SetEquipment(11, 12, 13)

	player.SetTeleported(true)
	if !player.IsTeleported() {
		t.Errorf("Expected IsTeleported true, got false")
	}

	player.SetDisappeared(true)
	if !player.IsDisappeared() {
		t.Errorf("Expected IsDisappeared true, got false")
	}

	now := time.Now()
	player.SetGoneTime(now)
	if player.GoneTime() != now {
		t.Errorf("Expected GoneTime %v, got %v", now, player.GoneTime())
	}

	player.AddDamageTaken("monster1", 100)
	player.AddDamageTaken("monster1", 50)
	player.AddDamageTaken("monster2", 200)

	if player.GetDamageTaken("monster1") != 150 {
		t.Errorf("Expected DamageTaken from monster1 150, got %d", player.GetDamageTaken("monster1"))
	}
	if player.GetDamageTaken("monster2") != 200 {
		t.Errorf("Expected DamageTaken from monster2 200, got %d", player.GetDamageTaken("monster2"))
	}

	player.AddDamageDone("monster1", 300)
	player.AddDamageDone("monster2", 400)

	if player.GetDamageDone("monster1") != 300 {
		t.Errorf("Expected DamageDone to monster1 300, got %d", player.GetDamageDone("monster1"))
	}
	if player.GetDamageDone("monster2") != 400 {
		t.Errorf("Expected DamageDone to monster2 400, got %d", player.GetDamageDone("monster2"))
	}

	// Test DeepCopy
	copy := player.DeepCopy().(*Player)
	if string(copy.ID()) != string(id) {
		t.Errorf("DeepCopy: Expected ID %v, got %v", id, copy.ID())
	}
	if copy.Name() != "TestPlayer" {
		t.Errorf("DeepCopy: Expected Name TestPlayer, got %s", copy.Name())
	}
	if copy.Level() != 50 {
		t.Errorf("DeepCopy: Expected Level 50, got %d", copy.Level())
	}
	if copy.GetDamageTaken("monster1") != 150 {
		t.Errorf("DeepCopy: Expected DamageTaken from monster1 150, got %d", copy.GetDamageTaken("monster1"))
	}
}

func TestMonster(t *testing.T) {
	id := []byte{1, 2, 3, 4}
	monster := NewMonster(id)

	// Test basic Actor methods
	monster.SetName("TestMonster")
	if monster.Name() != "TestMonster" {
		t.Errorf("Expected Name TestMonster, got %s", monster.Name())
	}

	// Test Monster-specific methods
	monster.SetBinType(100)
	if monster.BinType() != 100 {
		t.Errorf("Expected BinType 100, got %d", monster.BinType())
	}

	monster.SetNameGiven("GivenName")
	if monster.NameGiven() != "GivenName" {
		t.Errorf("Expected NameGiven GivenName, got %s", monster.NameGiven())
	}

	monster.SetLevel(30)
	if monster.Level() != 30 {
		t.Errorf("Expected Level 30, got %d", monster.Level())
	}

	monster.SetHP(500)
	monster.SetMaxHP(1000)
	if monster.HP() != 500 || monster.MaxHP() != 1000 {
		t.Errorf("Expected HP 500/1000, got %d/%d", monster.HP(), monster.MaxHP())
	}
	if monster.HPPercent() != 50 {
		t.Errorf("Expected HPPercent 50, got %d", monster.HPPercent())
	}

	monster.SetDead(true)
	if !monster.IsDead() {
		t.Errorf("Expected IsDead true, got false")
	}

	monster.SetTeleported(true)
	if !monster.IsTeleported() {
		t.Errorf("Expected IsTeleported true, got false")
	}

	monster.SetDisappeared(true)
	if !monster.IsDisappeared() {
		t.Errorf("Expected IsDisappeared true, got false")
	}

	now := time.Now()
	monster.SetGoneTime(now)
	if monster.GoneTime() != now {
		t.Errorf("Expected GoneTime %v, got %v", now, monster.GoneTime())
	}

	monster.AddDamageFromYou(100)
	monster.AddDamageFromYou(50)
	if monster.DamageFromYou() != 150 {
		t.Errorf("Expected DamageFromYou 150, got %d", monster.DamageFromYou())
	}

	monster.AddDamageFromParty(200)
	monster.AddDamageFromParty(100)
	if monster.DamageFromParty() != 300 {
		t.Errorf("Expected DamageFromParty 300, got %d", monster.DamageFromParty())
	}

	monster.AddDamageFrom("player1", 100)
	monster.AddDamageFrom("player1", 50)
	monster.AddDamageFrom("player2", 200)

	if monster.GetDamageFrom("player1") != 150 {
		t.Errorf("Expected DamageFrom player1 150, got %d", monster.GetDamageFrom("player1"))
	}
	if monster.GetDamageFrom("player2") != 200 {
		t.Errorf("Expected DamageFrom player2 200, got %d", monster.GetDamageFrom("player2"))
	}

	// Test DeepCopy
	copy := monster.DeepCopy().(*Monster)
	if string(copy.ID()) != string(id) {
		t.Errorf("DeepCopy: Expected ID %v, got %v", id, copy.ID())
	}
	if copy.Name() != "TestMonster" {
		t.Errorf("DeepCopy: Expected Name TestMonster, got %s", copy.Name())
	}
	if copy.Level() != 30 {
		t.Errorf("DeepCopy: Expected Level 30, got %d", copy.Level())
	}
	if copy.DamageFromYou() != 150 {
		t.Errorf("DeepCopy: Expected DamageFromYou 150, got %d", copy.DamageFromYou())
	}
}

func TestActorList(t *testing.T) {
	list := NewActorList()

	// Test empty list
	if list.Count() != 0 {
		t.Errorf("Expected Count 0, got %d", list.Count())
	}

	// Add actors
	actor1 := NewBaseActor([]byte{1, 2, 3, 4}, "TestActor")
	actor1.SetName("Actor1")
	list.Add(actor1)

	actor2 := NewBaseActor([]byte{5, 6, 7, 8}, "TestActor")
	actor2.SetName("Actor2")
	list.Add(actor2)

	// Test Count
	if list.Count() != 2 {
		t.Errorf("Expected Count 2, got %d", list.Count())
	}

	// Test GetByID
	if list.GetByID([]byte{1, 2, 3, 4}).Name() != "Actor1" {
		t.Errorf("Expected GetByID to return Actor1")
	}
	if list.GetByID([]byte{5, 6, 7, 8}).Name() != "Actor2" {
		t.Errorf("Expected GetByID to return Actor2")
	}

	// Test GetAll
	actors := list.GetAll()
	if len(actors) != 2 {
		t.Errorf("Expected GetAll to return 2 actors, got %d", len(actors))
	}

	// Test Contains
	if !list.Contains([]byte{1, 2, 3, 4}) {
		t.Errorf("Expected Contains to return true for actor1")
	}
	if !list.Contains([]byte{5, 6, 7, 8}) {
		t.Errorf("Expected Contains to return true for actor2")
	}
	if list.Contains([]byte{9, 10, 11, 12}) {
		t.Errorf("Expected Contains to return false for non-existent actor")
	}

	// Test Remove
	list.Remove(actor1)
	if list.Count() != 1 {
		t.Errorf("Expected Count 1 after Remove, got %d", list.Count())
	}
	if list.Contains([]byte{1, 2, 3, 4}) {
		t.Errorf("Expected Contains to return false after Remove")
	}

	// Test RemoveByID
	list.RemoveByID([]byte{5, 6, 7, 8})
	if list.Count() != 0 {
		t.Errorf("Expected Count 0 after RemoveByID, got %d", list.Count())
	}

	// Test Clear
	list.Add(actor1)
	list.Add(actor2)
	list.Clear()
	if list.Count() != 0 {
		t.Errorf("Expected Count 0 after Clear, got %d", list.Count())
	}

	// Test specialized lists
	playersList := NewPlayersList()
	player1 := NewPlayer([]byte{1, 2, 3, 4})
	player1.SetName("Player1")
	playersList.Add(player1)

	if playersList.Count() != 1 {
		t.Errorf("Expected PlayersList Count 1, got %d", playersList.Count())
	}

	if playersList.GetByID([]byte{1, 2, 3, 4}).Name() != "Player1" {
		t.Errorf("Expected PlayersList GetByID to return Player1")
	}

	monstersList := NewMonstersList()
	monster1 := NewMonster([]byte{5, 6, 7, 8})
	monster1.SetName("Monster1")
	monstersList.Add(monster1)

	if monstersList.Count() != 1 {
		t.Errorf("Expected MonstersList Count 1, got %d", monstersList.Count())
	}

	if monstersList.GetByID([]byte{5, 6, 7, 8}).Name() != "Monster1" {
		t.Errorf("Expected MonstersList GetByID to return Monster1")
	}
}
