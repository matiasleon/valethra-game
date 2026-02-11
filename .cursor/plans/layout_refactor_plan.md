# Refactoring Plan: Combat Scene Layout

## Problema actual

El código de `game/combat_scene.go` tiene layout hardcodeado y mezclado con rendering:

- **Magic numbers**: `heroX=120`, `enemyX=480`, `spriteScale=2.5` sin contexto
- **Responsabilidades mezcladas**: `Draw()` calcula posiciones mientras dibuja
- **Difícil de modificar**: cambiar resolución, agregar personajes, o ajustar espaciado requiere editar múltiples lugares
- **No extensible**: no hay forma de pasar layout custom sin editar el código fuente

## Solución propuesta

Separar **configuración de layout** de **cálculo de posiciones** de **rendering**.

### Arquitectura objetivo

```
LayoutConfig (struct)        → valores configurables (offsets, scales, margins)
    ↓
LayoutManager (struct)       → calcula posiciones a partir de config + sprites
    ↓
CombatLayout (struct)        → posiciones precalculadas para 1 frame
    ↓
Draw() + drawXXX() methods   → solo renderiza usando posiciones precalculadas
```

---

## Fase 1: Extraer LayoutConfig

**Objetivo**: Centralizar todas las constantes de layout en un struct configurable.

### Archivos a crear/modificar

- **Nuevo archivo**: `game/layout_config.go`
- **Modificar**: `game/combat_scene.go` (reemplazar constantes por `LayoutConfig`)

### Cambios

1. Crear `LayoutConfig` struct con todos los valores de layout
2. Crear `DefaultLayoutConfig()` constructor
3. Agregar `layoutCfg LayoutConfig` a `CombatScene`
4. Reemplazar todas las referencias a constantes por `c.layoutCfg.XXX`

**Beneficio**: Un solo lugar para ver/modificar todos los valores de layout.

---

## Fase 2: Introducir LayoutManager

**Objetivo**: Separar cálculo de posiciones del rendering.

### Archivos a crear/modificar

- **Nuevo archivo**: `game/layout_manager.go`
- **Modificar**: `game/combat_scene.go` (usar `LayoutManager` en `Draw()`)

### Cambios

1. Crear `CharacterLayout` struct (posiciones calculadas para 1 personaje)
2. Crear `CombatLayout` struct (todas las posiciones precalculadas)
3. Crear `LayoutManager` con método `Compute(heroSprites, enemySprites) CombatLayout`
4. Agregar `layoutManager *LayoutManager` a `CombatScene`
5. Modificar `Draw()` para llamar `layout := c.layoutManager.Compute(...)` al inicio

**Beneficio**: Layout se calcula 1 vez por frame, no durante cada draw call.

---

## Fase 3: Simplificar Draw() y draw* methods

**Objetivo**: Métodos de rendering solo reciben posiciones, no calculan nada.

### Cambios

1. `Draw()` solo orquesta: computa layout → llama drawXXX con posiciones
2. `drawCharacterWithSprite()` recibe `CharacterLayout` en vez de calcular offsets
3. `drawCombatLog()` recibe `logX, logY, maxLines` sin calcular nada
4. `drawStatus()` recibe `statusBarY, statusTextY`

**Beneficio**: Código más claro, fácil de testear layout por separado.

---

## Fase 4: Extensibilidad (opcional, futuro)

**Objetivo**: Permitir layout custom por escena.

### Cambios

1. Agregar `CombatSceneOption` pattern (functional options)
2. `NewCombatScene(..., opts ...CombatSceneOption)`
3. Opciones: `WithSpriteScale()`, `WithScreenSize()`, `WithCharacterPositions()`

**Beneficio**: Diferentes escenas pueden tener layout diferente sin editar código.

---

## Orden de implementación recomendado

```
Paso 1: layout_config.go
  └─ LayoutConfig struct
  └─ DefaultLayoutConfig() constructor

Paso 2: combat_scene.go
  └─ Agregar layoutCfg a CombatScene
  └─ Reemplazar constantes por layoutCfg.XXX

Paso 3: layout_manager.go
  └─ CharacterLayout struct
  └─ CombatLayout struct
  └─ LayoutManager.Compute()

Paso 4: combat_scene.go
  └─ Usar LayoutManager en Draw()
  └─ Simplificar draw* methods

Paso 5: Limpiar (opcional)
  └─ Remover constantes viejas
  └─ Tests de LayoutManager
```

---

## Ejemplo de uso final

### Antes (hardcodeado):

```go
const spriteScale = 2.5
const (
    titleY = 12
    eventDescY = 35
    // ... 20+ más
)

func (c *CombatScene) Draw(screen *ebiten.Image) {
    heroX := 120
    enemyX := 480
    // cálculos mezclados con rendering...
}
```

### Después (configurable):

```go
cfg := DefaultLayoutConfig(800, 600)
cfg.SpriteScale = 3.0  // fácil de cambiar
cfg.HeroX = 150        // fácil de ajustar

lm := NewLayoutManager(cfg)
layout := lm.Compute(heroSprites, enemySprites)

// Draw() solo renderiza
ebitenutil.DebugPrintAt(screen, title, layout.TitleX, layout.TitleY)
c.drawCharacterWithSprite(screen, hero, layout.Hero, sprites)
```

---

## Ventajas finales

✅ **Claridad**: layout en un solo lugar, fácil de entender  
✅ **Modificabilidad**: cambiar valores sin tocar lógica de rendering  
✅ **Testeable**: `LayoutManager` se puede testear sin Ebiten  
✅ **Extensible**: diferentes escenas pueden tener layout diferente  
✅ **Menos bugs**: layout separado de rendering = menos side effects  

---

## Esfuerzo estimado

- **Fase 1**: ~30 minutos (crear config, reemplazar constantes)
- **Fase 2**: ~45 minutos (crear manager, structs)
- **Fase 3**: ~30 minutos (refactorizar draw methods)
- **Total**: ~2 horas de refactoring incremental

Podemos implementar fase por fase, verificando que compile y funcione después de cada paso.
