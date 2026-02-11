# Análisis: Modelo de Quest extendido con ActionEvent

## 📊 Estado actual vs Propuesta

### Estado ACTUAL (lo que rompí en el refactoring):

```go
type Event struct {
    Description string
    Enemies     []Character  // ❌ Solo combates, acoplado
}
```

**Uso:**
```go
event := quest.Events[0]
enemy := event.Enemies[0]  // acceso directo
combat := NewCombatScene(quest, hero, enemy, ...)
```

### Propuesta ORIGINAL (la que vos tenías antes):

```go
type Event struct {
    Description string
    ActionEvent []ActionEvent
}

type ActionEvent struct {
    Description string
    Action      func() ActionResult
}
```

**Ventajas:**
✅ **Polimórfico**: Event no sabe si es combate, diálogo, puzzle  
✅ **Extensible**: podés agregar nuevos tipos de acciones sin tocar Event  
✅ **Flexible**: un evento puede tener múltiples acciones secuenciales  

---

## 🎯 Comparación: Diseño actual vs Extensible

### Limitaciones del modelo actual (`Enemies []Character`):

| Limitación | Impacto |
|------------|---------|
| **Solo combates** | No podés tener diálogos, puzzles, cutscenes |
| **Acoplado** | `Event` conoce detalles de combate (`Character`) |
| **No secuencial** | ¿Qué pasa si querés: diálogo → puzzle → combate? |
| **No condicional** | ¿Qué pasa si el jugador puede elegir pelear o negociar? |

**Ejemplo de lo que NO podés hacer hoy:**

```
Quest: "El Pacto de Grukh"
  Event 1: Diálogo con Grukh (opción: pelear o pagar peaje)
    → Si pagás: skip combate, pierdes oro
    → Si atacás: combate
  Event 2: Puzzle (abrir puerta secreta con runas)
  Event 3: Combate final (guardián del tesoro)
```

### Ventajas del modelo `ActionEvent`:

| Ventaja | Ejemplo |
|---------|---------|
| **Polimorfismo** | `CombatAction`, `DialogueAction`, `PuzzleAction` |
| **Composición** | Event = [Dialogue, Choice, Combat] |
| **Reutilización** | `CombatAction` se usa en N eventos diferentes |
| **Testing** | Podés testear `ActionEvent` sin Ebiten |

---

## 🏗️ Diseño propuesto: ActionEvent como interfaz

### Opción 1: Interface Go (recomendado)

```go
// ActionEvent es cualquier cosa que el jugador puede hacer en un evento.
type ActionEvent interface {
    Execute(hero *Character) ActionResult
    Description() string
}

// CombatAction es una acción de combate.
type CombatAction struct {
    Enemies []Character
}

func (c *CombatAction) Execute(hero *Character) ActionResult {
    // lógica de combate
    return ActionResult{Success: true, Message: "Victoria"}
}

func (c *CombatAction) Description() string {
    return fmt.Sprintf("Enfrentar a %d enemigos", len(c.Enemies))
}

// DialogueAction es una acción de diálogo.
type DialogueAction struct {
    NPC     string
    Options []DialogueOption
}

func (d *DialogueAction) Execute(hero *Character) ActionResult {
    // lógica de diálogo (UI diferente)
    return ActionResult{Success: true, Message: "Diálogo completado"}
}

func (d *DialogueAction) Description() string {
    return fmt.Sprintf("Hablar con %s", d.NPC)
}
```

**Uso:**

```go
quest := Quest{
    Title: "El Pacto de Grukh",
    Events: []Event{
        {
            Description: "Llegas al puente...",
            Actions: []ActionEvent{
                &DialogueAction{
                    NPC: "Grukh",
                    Options: []DialogueOption{
                        {Text: "Pagar peaje", Effect: ...},
                        {Text: "Atacar", Effect: ...},
                    },
                },
                &CombatAction{
                    Enemies: []Character{grukh},
                },
            },
        },
    },
}
```

---

## 🔄 Impacto en el código actual

### Archivos que hay que modificar:

#### 1. `core/entities/quest.go`

```go
type Event struct {
    Description string
    Actions     []ActionEvent  // cambio: Enemies → Actions
}

type ActionEvent interface {
    Execute(ctx *ActionContext) ActionResult
    Type() ActionType  // "combat", "dialogue", "puzzle"
}

type ActionType string
const (
    ActionCombat   ActionType = "combat"
    ActionDialogue ActionType = "dialogue"
    ActionPuzzle   ActionType = "puzzle"
)
```

#### 2. `core/entities/actions.go` (nuevo archivo)

```go
// CombatAction representa un combate.
type CombatAction struct {
    Enemies []Character
}

func (c *CombatAction) Execute(ctx *ActionContext) ActionResult {
    // Delega a CombatEngine
    return ctx.Engine.ExecuteCombat(ctx.Hero, c.Enemies)
}

func (c *CombatAction) Type() ActionType {
    return ActionCombat
}
```

#### 3. `game/combat_scene.go`

**Antes:**
```go
func NewCombatScene(quest Quest, hero *Character, ...) {
    enemy := quest.Events[0].Enemies[0]  // ❌ Asumir que es combate
}
```

**Después:**
```go
func NewCombatScene(action *CombatAction, hero *Character, ...) {
    enemy := action.Enemies[0]  // ✅ Explícito que es combate
}
```

#### 4. `game/game.go` (SceneManager)

**Antes:**
```go
phase := phaseCombat
scene := NewCombatScene(quest, hero, ...)
```

**Después:**
```go
action := event.Actions[actionIndex]
switch action.Type() {
case ActionCombat:
    scene = NewCombatScene(action.(*CombatAction), hero, ...)
case ActionDialogue:
    scene = NewDialogueScene(action.(*DialogueAction), hero, ...)
case ActionPuzzle:
    scene = NewPuzzleScene(action.(*PuzzleAction), hero, ...)
}
```

---

## 🎨 Mejoras del diseño propuesto

### 1. **Separación de Responsabilidades**

| Capa | Responsabilidad | No sabe de... |
|------|-----------------|---------------|
| `entities.Quest` | Estructura de quest | Cómo se ejecutan las acciones |
| `entities.ActionEvent` | Interfaz de acción | UI, Ebiten, sprites |
| `game.CombatScene` | UI de combate | Quests, eventos |
| `game.SceneManager` | Orquestación | Detalles de cada acción |

### 2. **Extensibilidad sin romper código existente**

```go
// Hoy: solo combates
Events: []Event{
    {Actions: []ActionEvent{&CombatAction{...}}},
}

// Mañana: agregar diálogos sin tocar CombatScene
Events: []Event{
    {Actions: []ActionEvent{
        &DialogueAction{...},
        &CombatAction{...},
    }},
}

// Pasado mañana: agregar puzzles
Events: []Event{
    {Actions: []ActionEvent{
        &PuzzleAction{Type: "lockpick"},
        &CombatAction{...},
    }},
}
```

### 3. **Testing más fácil**

```go
func TestCombatAction(t *testing.T) {
    action := &CombatAction{
        Enemies: []Character{weakEnemy},
    }
    result := action.Execute(ctx)
    assert.True(t, result.Success)
}

func TestDialogueAction(t *testing.T) {
    action := &DialogueAction{...}
    result := action.Execute(ctx)
    // No necesitás UI para testear lógica
}
```

---

## 🚀 Plan de migración recomendado

### Fase 1: Agregar interfaz sin romper nada (Adapter Pattern)

```go
// quest.go
type Event struct {
    Description string
    Actions     []ActionEvent  // nuevo
    Enemies     []Character    // deprecado, pero sigue funcionando
}

// Helper para migración gradual
func (e *Event) GetCombatAction() *CombatAction {
    if len(e.Actions) > 0 {
        for _, action := range e.Actions {
            if ca, ok := action.(*CombatAction); ok {
                return ca
            }
        }
    }
    // Fallback: usar Enemies legacy
    if len(e.Enemies) > 0 {
        return &CombatAction{Enemies: e.Enemies}
    }
    return nil
}
```

### Fase 2: Migrar `main.go` a usar Actions

```go
// Antes
Events: []Event{
    {Enemies: []Character{grukh}},
}

// Después
Events: []Event{
    {Actions: []ActionEvent{
        &CombatAction{Enemies: []Character{grukh}},
    }},
}
```

### Fase 3: Migrar `game/` a usar Actions

```go
func (g *Game) Update() error {
    action := currentEvent.GetCombatAction()
    if action != nil {
        g.currentScene = NewCombatScene(action, hero, ...)
    }
}
```

### Fase 4: Eliminar `Enemies` legacy

Una vez que todo use `Actions`, eliminar el campo `Enemies`.

---

## 💡 Resumen: ¿Vale la pena?

### Pros de migrar a ActionEvent:

✅ **Juego más rico**: diálogos, puzzles, cutscenes  
✅ **Código más limpio**: cada acción sabe cómo ejecutarse  
✅ **Testing más fácil**: actions testeables sin UI  
✅ **Extensible**: agregar tipos sin tocar código existente  

### Contras:

❌ **Refactoring**: ~3-4 archivos a modificar  
❌ **Complejidad inicial**: interfaz + type assertions  
❌ **Type safety**: necesitás casting `action.(*CombatAction)`  

### Recomendación:

**SÍ, migrá al modelo ActionEvent** si:
- Querés agregar diálogos, puzzles, o eventos no-combate
- Querés que el juego sea más rico narrativamente
- Te importa la separación de responsabilidades

**NO** si:
- Solo vas a tener combates (nunca otra cosa)
- Preferís simplicidad sobre extensibilidad

---

## 🎯 Mi opinión

El modelo `ActionEvent` es **arquitecturalmente superior**. Te da:
1. **Flexibilidad narrativa** (diálogos con opciones)
2. **Gameplay variado** (puzzles, stealth, trading)
3. **Código desacoplado** (UI no conoce estructura de quest)

El único trade-off es complejidad inicial, pero el payoff es enorme si querés un juego más completo que solo "atacar todo el tiempo".

¿Querés que implemente la migración al modelo ActionEvent?
