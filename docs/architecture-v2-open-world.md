# Architecture: v1.0 Linear → v2.0 Open World

Reference document for migrating Valethra from a linear scene-based RPG (v1.0) to an open-world 2D game with event loop (v2.0).

**Do not implement anything from v2.0 until v1.0 is stable and tagged.**

---

## v1.0 — Linear Scene-Based (Current)

### Architecture

```
main.go → Game (state machine) → Scene stack (linear)
                                    ├─ IntroScene
                                    ├─ QuestIntroScene
                                    ├─ CombatScene
                                    └─ ResultScene (victory/defeat)
```

### Flow

```
Intro → QuestIntro → Combat → [Victory?]
                                 ├─ Yes + more quests → QuestIntro (next)
                                 ├─ Yes + no more    → VictoryScene → Exit
                                 └─ No               → DefeatScene → Exit
```

### Key Files

| File | Responsibility |
|------|----------------|
| `game/scene.go` | Scene interface + SceneResult |
| `game/game.go` | Game loop, phase state machine, scene transitions |
| `game/combat_scene.go` | Combat rendering, input, round execution, multi-event |
| `game/transition_scene.go` | Text screens (intro, quest intro, victory, defeat) |
| `game/sprite.go` | Animator with 4 states (idle, attack, hurt, dead) |
| `game/assets.go` | Embedded asset loading, sprite sheet slicing |
| `core/game_master.go` | Combat engine, quest progression, hero healing |
| `core/entities/` | Character, Attributes, Quest, Event, ActionEvent |
| `core/ports/combat.go` | CombatEngine interface, RoundResult |

### What stays in v2.0

- `core/entities/` — domain model unchanged
- `core/ports/combat.go` — CombatEngine interface unchanged
- `core/game_master.go` — ExecuteRound logic unchanged
- `game/sprite.go` — Animator extended (new states, same core)
- `game/assets.go` — asset loading unchanged

---

## v2.0 — Open World 2D

### Target Architecture

```
main.go → Game → SceneManager (stack)
                    │
                    ├─ WorldScene (always at bottom, never "ends")
                    │   ├─ TileMap
                    │   ├─ Camera
                    │   ├─ Player (real-time movement)
                    │   ├─ WorldEntities (NPCs, enemies, objects)
                    │   ├─ TriggerSystem
                    │   └─ EventBus
                    │
                    ├─ CombatScene (pushed on encounter)
                    ├─ DialogueScene (pushed on NPC interaction)
                    └─ TransitionScene (pushed for cutscenes)
```

### Core Design Patterns

#### 1. Entity-Component System (ECS Hybrid)

Every "thing" in the world is an entity with composable components.

```go
// Entities are IDs. Components are data. Systems operate on components.

type EntityID int

type TransformComponent struct {
    X, Y     float64
    Rotation float64
}

type RenderComponent struct {
    Sprites *CharacterSprites
    Layer   int // draw order
}

type PhysicsComponent struct {
    Collider AABB
    Solid    bool
}

type AIComponent struct {
    BehaviorTree *BehaviorTree
    Blackboard   map[string]any
}

type CombatComponent struct {
    Character *entities.Character // reuse existing domain entity
}

type InteractableComponent struct {
    Type    string // "talk", "open", "examine"
    Handler func(world *WorldScene)
}
```

Why: avoids deep inheritance hierarchies. A chest is Transform+Render+Interactable. An NPC is Transform+Render+AI+Combat+Interactable. A tree is Transform+Render+Physics.

#### 2. State Machine — Combat

Hierarchical state machine for player combat actions.

```
PlayerCombatState
  ├─ Idle
  ├─ Attacking
  │   ├─ FastAttack (sub-state)
  │   ├─ HeavyAttack (sub-state)
  │   └─ Combo chain (fast → fast → heavy)
  ├─ Dodging
  ├─ Casting (Signs/Spells)
  └─ Staggered (hit reaction)
```

Each state defines:
- Valid transitions (idle → attack, attack → dodge, NOT attack → attack)
- Animation to play
- Hitbox activation frames
- Input buffer window

#### 3. Behavior Tree + Blackboard — AI

Each NPC/monster has a behavior tree for decision-making.

```
Monster: Troll
  Selector
  ├─ Sequence: "Combat"
  │   ├─ Condition: player_in_range(20)?
  │   ├─ Selector: choose attack
  │   │   ├─ Seq: club_smash (if player close)
  │   │   ├─ Seq: throw_rock (if player far)
  │   │   └─ Seq: charge (if player medium range)
  │   └─ Action: execute_attack
  ├─ Sequence: "Patrol"
  │   ├─ Condition: NOT player_in_range?
  │   └─ Action: walk_to(next_patrol_point)
  └─ Sequence: "Flee"
      ├─ Condition: health < 15%?
      └─ Action: flee_to(safe_point)

Blackboard:
  target: nil | *Player
  last_seen: (x, y)
  health_pct: 0.45
  patrol_idx: 2
```

Node types:
- **Selector**: tries children left→right, returns first success
- **Sequence**: runs children left→right, fails on first failure
- **Condition**: checks a predicate
- **Action**: executes a behavior (returns running/success/fail)

#### 4. Quest System — DAG + World Facts

Quests are directed graphs with conditions based on global "world facts".

```
Quest: "El Peaje de Grukh"

  Node_1: "Llegar al puente"
    trigger: enter_area("varnock_bridge")
    on_complete → Node_2

  Node_2: "Enfrentar a Grukh"
    trigger: interact("grukh")
    choices:
      ├─ "Pelear" → Node_3a (combat)
      └─ "Negociar" → Node_3b (dialogue, skill check)

  Node_3a: combat result
    on_victory → Node_4a [set_fact("grukh_dead")]
    on_defeat  → Node_4b [set_fact("player_defeated")]

  Node_3b: negotiation
    on_success → Node_4c [set_fact("grukh_ally")]
    on_fail    → Node_3a (combat)

  Node_4a/4c: quest complete
    effects: reward_gold(200), update_reputation("varnock", +10)
```

World facts: global key-value store that persists across quests.

```go
type WorldFacts map[string]any

// Examples:
// facts["grukh_dead"] = true
// facts["baron_reputation"] = 15
// facts["ciri_location"] = "novigrad"
```

Every quest node can READ facts as preconditions and WRITE facts as effects.

#### 5. Event Bus — Observer (Global)

Decouples systems. A single event triggers reactions in many systems.

```go
type EventBus struct {
    listeners map[string][]EventHandler
}

type EventHandler func(data any)

func (eb *EventBus) Emit(event string, data any)
func (eb *EventBus) On(event string, handler EventHandler)
```

Example flow:

```
EventBus.Emit("enemy_killed", {enemy: grukh, position: (x,y)})
  │
  ├─ QuestSystem    → advance quest objective
  ├─ XPSystem       → grant experience
  ├─ LootSystem     → spawn drops at position
  ├─ AISystem       → nearby enemies react (flee/aggro)
  ├─ WorldFacts     → set_fact("grukh_dead", true)
  └─ AudioSystem    → play death sound
```

#### 6. Spatial Partitioning — World Streaming

World divided into chunks loaded/unloaded based on player position.

```go
type Chunk struct {
    X, Y     int        // chunk coordinates
    Tiles    [][]Tile
    Entities []EntityID
    Loaded   bool
}

type WorldStreamer struct {
    chunks    map[[2]int]*Chunk
    loadRadius int // chunks around player to keep loaded
}

func (ws *WorldStreamer) Update(playerChunkX, playerChunkY int)
// loads nearby chunks, unloads distant ones
```

#### 7. Dialogue System — Tree with Conditions

```go
type DialogueNode struct {
    Speaker    string
    Text       string
    Conditions []FactCondition      // required facts to show this node
    Effects    []FactEffect         // facts to set when this node plays
    Choices    []DialogueChoice
}

type DialogueChoice struct {
    Text       string
    SkillCheck *SkillCheck          // optional: requires stat >= threshold
    NextNode   *DialogueNode
}

type SkillCheck struct {
    Attribute string // "intelligence", "willpower"
    Threshold int
}
```

#### 8. Save/Load — Memento

```go
type SaveState struct {
    PlayerPosition  [2]float64
    PlayerStats     entities.Attributes
    Inventory       []ItemID
    WorldFacts      WorldFacts
    QuestStates     map[string]QuestNodeID
    DeadEntities    []EntityID
    OpenedChests    []EntityID
    WorldTime       float64
    ChunkMods       map[[2]int][]ChunkModification
}
```

---

## New Files to Create (v2.0)

```
game/
  ├─ scene_manager.go      // Scene stack (push/pop)
  ├─ world_scene.go        // Main open-world scene
  ├─ tilemap.go            // Tile types, map loading, rendering
  ├─ camera.go             // Viewport, follow player, world↔screen coords
  ├─ player.go             // Real-time movement, input, collision
  ├─ collision.go          // AABB, tile collision, entity collision
  ├─ trigger.go            // Area triggers (combat, dialogue, quest)
  ├─ dialogue_scene.go     // Dialogue UI overlay
  ├─ npc.go                // NPC entity with behavior
  └─ world_entity.go       // Base for all world entities

core/
  ├─ event_bus.go          // Global event system
  ├─ world_facts.go        // Key-value world state
  ├─ quest_graph.go        // Quest as DAG with conditions
  ├─ behavior_tree.go      // BT nodes (selector, sequence, action)
  ├─ save.go               // Save/load serialization
  └─ entities/
      └─ item.go           // Item, EquipableItem, LootTable
```

## Files to Modify (v2.0)

| File | Change |
|------|--------|
| `game/scene.go` | Add `OnPause()` / `OnResume()` to Scene interface for stack |
| `game/game.go` | Replace phase enum with SceneManager stack |
| `game/sprite.go` | Add walk states: WalkUp, WalkDown, WalkLeft, WalkRight |
| `game/combat_scene.go` | Receive single enemy, return result to WorldScene |
| `game/assets.go` | Add tileset loading, map file loading |
| `main.go` | Load world map, pass to Game instead of quest list |

---

## Migration Steps

### Phase 1: Scene Stack (prerequisite for everything)

Replace `game.phase` state machine with a `SceneManager` that supports push/pop.
WorldScene stays at the bottom of the stack, CombatScene pushes on top.

### Phase 2: World Foundation

1. TileMap — define tile format, load from file, render visible tiles
2. Camera — follow player, world-to-screen coordinate conversion
3. Player — WASD movement, tile collision, walk animations
4. WorldScene — integrates tilemap + player + camera

Milestone: Aldric walks around a 2D map.

### Phase 3: Connect World to Combat

1. TriggerSystem — define zones/entities that start combat
2. Adapt CombatScene — receive single enemy, push/pop from world
3. Enemy entities visible on map before combat

Milestone: walk into enemy on map → combat → return to map.

### Phase 4: NPCs and Dialogue

1. NPC entities with basic AI (stand, patrol)
2. DialogueScene — tree-based dialogue UI
3. Interaction system — press key near NPC to talk

Milestone: talk to NPCs, get quest info through dialogue.

### Phase 5: Quest Graph

1. WorldFacts — global key-value state
2. QuestGraph — DAG with conditions and effects
3. EventBus — connect quest completions to world changes

Milestone: complete quest → world reacts (NPC moves, area changes).

### Phase 6: Full Systems

1. Behavior trees for monster AI
2. Inventory and items
3. Save/load
4. World streaming (if map gets large enough)

---

## Damage Calculation Reference (v2.0 real-time combat)

When combat transitions from turn-based to real-time, the damage formula expands:

```
finalDamage = baseDamage
    * weaponMultiplier       // weapon type vs enemy type
    * critMultiplier         // if critical hit (based on agility)
    * buffMultiplier         // active spell/potion effects
    - enemyArmor
    - enemyResistance[type]  // per damage type (fire, silver, etc.)
```

The existing `Character.Attack()` method becomes the base. Additional multipliers layer on top.

---

## Notes

- The core domain (`entities`, `ports`, `CombatEngine`) should remain untouched through v2.0
- 90% of v2.0 work is in the `game/` layer
- Each phase is independently shippable — the game is playable after each phase
- Data-drive everything: tile maps, quest graphs, dialogue trees, enemy stats should be loaded from files, not hardcoded in Go
