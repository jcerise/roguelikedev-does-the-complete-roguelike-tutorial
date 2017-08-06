package ecs

type Component interface {
	IsAIComponent() bool
}

// Player Component
type PlayerComponent struct {
}

func (pl PlayerComponent) IsAIComponent() bool {
	return false
}

// Position Component
type PositionComponent struct {
	X int
	Y int
}

func (pc PositionComponent) IsAIComponent() bool {
	return false
}

// Appearance Component
type AppearanceComponent struct {
	Color     string
	Character string
	Layer     int
	Name      string
}

func (a AppearanceComponent) IsAIComponent() bool {
	return false
}

// Movement Component
type MovementComponent struct {
}

func (m MovementComponent) IsAIComponent() bool {
	return false
}

// Controllable Component
type ControllableComponent struct {
}

func (c ControllableComponent) IsAIComponent() bool {
	return false
}

type HitPointComponent struct {
	Hp    int
	MaxHP int
}

func (h HitPointComponent) IsAIComponent() bool {
	return false
}

// Attacker Component
type AttackerComponent struct {
	Attack  int
	Defense int
}

func (a AttackerComponent) IsAIComponent() bool {
	return false
}

// Blocking Component
type BlockingComponent struct {
}

func (b BlockingComponent) IsAIComponent() bool {
	return false
}

// Random Movement Component - wanders aimlessly around the map
type RandomMovementComponent struct {
}

func (r RandomMovementComponent) IsAIComponent() bool {
	return true
}

// Basic Melee Attack AI Component
type BasicMeleeAIComponent struct {
	target *GameEntity
}

func (b BasicMeleeAIComponent) IsAIComponent() bool {
	return true
}

// Reproduces Component
type ReproducesComponent struct {
	MaxTimes       int
	TimesRemaining int
	PercentChance  int
}

func (r ReproducesComponent) IsAIComponent() bool {
	return false
}

// Killable Component
type KillableComponent struct {
	Character string
	Color     string
	Name      string
}

func (k KillableComponent) IsAIComponent() bool {
	return false
}