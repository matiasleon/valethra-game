# Migración a ActionEvent - Resumen

## ✅ Fases Completadas

### Fase 1: Agregar ActionEvent sin romper nada (Adapter Pattern)

**Archivos creados:**
- **`core/entities/actions.go`** (107 líneas)
  - Interface `ActionEvent` con métodos `Execute()`, `Type()`, `Description()`
  - `ActionType` enum: `combat`, `dialogue`, `puzzle`, `cutscene`
  - `CombatAction` struct (implementa `ActionEvent`)
  - `DialogueAction` struct (placeholder para futuro)
  - `ActionContext` y `ActionResult` para ejecutar acciones

**Archivos modificados:**
- **`core/entities/quest.go`**
  - Agregado `Actions []ActionEvent` a `Event`
  - Mantenido `Enemies []Character` como deprecado (backward compatibility)
  - Helper method `GetCombatAction()` para migración gradual
  - Helper methods `HasCombat()` y `GetActionsByType()`

- **`game/combat_scene.go`**
  - `NewCombatScene()` usa `GetCombatAction()` en lugar de acceso directo a `Enemies`
  - `advanceToNextEvent()` también migrado a `GetCombatAction()`

### Fase 2: Migrar main.go a usar Actions

**Archivos modificados:**
- **`main.go`**
  - `createTrollQuest()` migrado: `Enemies` → `Actions: [&CombatAction{...}]`
  - `createForestQuest()` migrado (2 eventos)

---

## 📊 Estado del código

### ✅ Lo que funciona ahora:

```go
// Sintaxis nueva (recomendada)
Events: []entities.Event{
    {
        Description: "...",
        Actions: []entities.ActionEvent{
            &entities.CombatAction{
                Enemies: []entities.Character{...},
                Desc: "Enfrentar a Grukh el Troll",
            },
        },
    },
}

// Sintaxis vieja (todavía funciona por backward compatibility)
Events: []entities.Event{
    {
        Description: "...",
        Enemies: []entities.Character{...},  // ⚠️ deprecado pero funcional
    },
}
```

### 🎯 Cómo acceder a enemigos:

```go
// Antes (acceso directo):
enemy := event.Enemies[0]  // ❌ Ya no recomendado

// Ahora (a través de GetCombatAction):
combatAction := event.GetCombatAction()
if combatAction != nil {
    enemy := combatAction.Enemies[0]  // ✅ Funciona con ambas sintaxis
}
```

---

## 🚀 Beneficios logrados

### 1. **Polimorfismo de acciones**

Ahora un `Event` puede tener **cualquier tipo de acción**:

```go
Events: []entities.Event{
    {
        Description: "Encuentro en la taberna...",
        Actions: []entities.ActionEvent{
            &entities.DialogueAction{
                NPC: "Tabernero",
                Options: []DialogueOption{
                    {Text: "Comprar bebida", ...},
                    {Text: "Preguntar por Grukh", ...},
                },
            },
            &entities.CombatAction{
                Enemies: []entities.Character{...}, // si eliges pelear
            },
        },
    },
}
```

### 2. **Backward compatibility total**

Todo el código viejo que usa `Enemies` **sigue funcionando** gracias al Adapter Pattern:

```go
// GetCombatAction() busca en Actions O en Enemies legacy
func (e *Event) GetCombatAction() *CombatAction {
    // 1. Intenta encontrar CombatAction en Actions (nuevo)
    for _, action := range e.Actions {
        if ca, ok := action.(*CombatAction); ok {
            return ca
        }
    }
    
    // 2. Fallback: crea CombatAction desde Enemies (legacy)
    if len(e.Enemies) > 0 {
        return &CombatAction{Enemies: e.Enemies}
    }
    
    return nil
}
```

### 3. **Fácil agregar nuevos tipos de acciones**

Para agregar un nuevo tipo (ej: `PuzzleAction`):

1. Crear struct que implemente `ActionEvent`:

```go
type PuzzleAction struct {
    Type        string // "lockpick", "riddle", "cipher"
    Difficulty  int
    Reward      string
}

func (p *PuzzleAction) Execute(ctx *ActionContext) ActionResult {
    // lógica del puzzle
    return ActionResult{Success: true}
}

func (p *PuzzleAction) Type() ActionType {
    return ActionPuzzle
}

func (p *PuzzleAction) Description() string {
    return "Resolver " + p.Type
}
```

2. Usar en quests:

```go
Events: []entities.Event{
    {
        Description: "Una puerta con candado rúnico...",
        Actions: []entities.ActionEvent{
            &entities.PuzzleAction{
                Type: "lockpick",
                Difficulty: 5,
            },
        },
    },
}
```

3. Crear scene correspondiente:

```go
// game/puzzle_scene.go
type PuzzleScene struct {
    puzzle *entities.PuzzleAction
    // ...
}
```

---

## 🔄 Próximos pasos (opcional)

### Fase 3: Scene routing por ActionType

Modificar `game/game.go` para crear escenas según el tipo de acción:

```go
func (g *Game) Update() error {
    // ...
    event := g.gm.CurrentQuestData().Events[eventIndex]
    
    for _, action := range event.Actions {
        switch action.Type() {
        case entities.ActionCombat:
            combatAction := action.(*entities.CombatAction)
            g.currentScene = NewCombatScene(combatAction, hero, ...)
        
        case entities.ActionDialogue:
            dialogueAction := action.(*entities.DialogueAction)
            g.currentScene = NewDialogueScene(dialogueAction, hero, ...)
        
        case entities.ActionPuzzle:
            puzzleAction := action.(*entities.PuzzleAction)
            g.currentScene = NewPuzzleScene(puzzleAction, hero, ...)
        }
    }
    // ...
}
```

### Fase 4: Eliminar `Enemies` legacy

Cuando estés seguro de que todo usa `Actions`, eliminar el campo deprecado:

```go
type Event struct {
    Description string
    Actions     []ActionEvent
    // Enemies []Character  ← eliminar
}
```

---

## 📝 Ejemplo completo: Quest multi-acción

```go
func createRichQuest() entities.Quest {
    return entities.Quest{
        Title: "La Conspiración del Barón",
        Events: []entities.Event{
            // Evento 1: Diálogo en la taberna
            {
                Description: "Entras a la taberna \"El Dragón Ebrio\"...",
                Actions: []entities.ActionEvent{
                    &entities.DialogueAction{
                        NPC: "Informante Encapuchado",
                        Options: []DialogueOption{
                            {Text: "Pagar por información (10 oro)", ...},
                            {Text: "Intimidar", ...},
                            {Text: "Irse", ...},
                        },
                    },
                },
            },
            
            // Evento 2: Infiltración (puzzle)
            {
                Description: "Llegas al castillo del Barón. Hay guardias en la entrada...",
                Actions: []entities.ActionEvent{
                    &entities.PuzzleAction{
                        Type: "stealth",
                        Description: "Esquivar a los guardias",
                    },
                },
            },
            
            // Evento 3: Combate con guardias
            {
                Description: "Te descubren! Los guardias te atacan.",
                Actions: []entities.ActionEvent{
                    &entities.CombatAction{
                        Enemies: []entities.Character{guard1, guard2},
                    },
                },
            },
            
            // Evento 4: Diálogo final con el Barón
            {
                Description: "El Barón te espera en su trono...",
                Actions: []entities.ActionEvent{
                    &entities.DialogueAction{
                        NPC: "Barón Corvus",
                        Options: []DialogueOption{
                            {Text: "Negociar", ...},
                            {Text: "Atacar", ...},
                        },
                    },
                    &entities.CombatAction{
                        Enemies: []entities.Character{baron},
                        // Solo si eliges "Atacar"
                    },
                },
            },
        },
    }
}
```

---

## ✨ Conclusión

La migración a `ActionEvent` está **completa y funcionando**:

✅ Código 100% backward compatible (viejo sigue funcionando)  
✅ Nueva arquitectura polimórfica lista para extender  
✅ Fácil agregar diálogos, puzzles, cutscenes sin tocar código existente  
✅ Build exitoso, sin errores de compilación  

El juego puede seguir creciendo con nuevos tipos de gameplay sin refactorings grandes.
