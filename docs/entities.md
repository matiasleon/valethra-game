# Entidades del Juego y Relaciones

Documentación de las entidades del juego **Valethra Game** y cómo se relacionan entre sí.

---

## Diagrama de Relaciones (resumen)

```
GameMaster ──┬── tiene muchos ──► Quest
              │
              └── inicia ──► Quest

Quest ── tiene muchos ──► Event

Event ── tiene muchos ──► Enemy  (Enemy referenciado; puede ser Character o tipo propio)

Game ── tiene muchos ──► Round

Character ── tiene uno ──► Attributes
```

---

## Entidades Principales

### 1. Character (Personaje)

**Ubicación:** `models/character.go`

Representa a un personaje del juego (jugador, NPC o enemigo).

| Campo        | Tipo        | Descripción                          |
|-------------|-------------|--------------------------------------|
| ID          | string      | Identificador único                  |
| Name        | string      | Nombre del personaje                 |
| Description | string      | Descripción o trasfondo              |
| Attributes  | Attributes  | Atributos numéricos (ver más abajo)   |

**Comportamiento:**
- `Attack(target *Character) string` — Ataca a otro personaje (reglas en `game-rules.md` / `docs/fight-system.md`).
- `Spell(target *Character) string` — Lanza un hechizo sobre un objetivo.

**Relaciones:**
- **Tiene uno** → `Attributes` (composición).

---

### 2. Attributes (Atributos)

**Ubicación:** `models/character.go`

Atributos numéricos de un personaje. No es una entidad independiente; siempre pertenece a un `Character`.

| Campo        | Tipo | Descripción                    |
|-------------|------|--------------------------------|
| Health      | int  | Vida                           |
| AttackPower | int  | Poder de ataque                |
| Armor       | int  | Armadura (se reduce al recibir ataques) |
| Agility     | int  | Agilidad                       |
| Intelligence| int  | Inteligencia                   |
| Willpower   | int  | Fuerza de voluntad             |
| Speed       | int  | Velocidad                      |
| Level       | int  | Nivel                          |

**Relaciones:**
- **Pertenece a** → `Character` (1:1).

---

### 3. Game (Partida)

**Ubicación:** `models/game.go`

Representa una partida o sesión de juego.

| Campo  | Tipo    | Descripción        |
|--------|---------|--------------------|
| Rounds | []Round | Lista de rondas    |

**Relaciones:**
- **Tiene muchos** → `Round` (1:N).

---

### 4. Round (Ronda)

**Ubicación:** `models/game.go`

Representa una ronda dentro de una partida. Estructura actualmente vacía (placeholder para turnos, combates, etc.).

**Relaciones:**
- **Pertenece a** → `Game` (N:1).

---

### 5. Quest (Misión)

**Ubicación:** `models/quest.go`

Una misión o línea de historia que el jugador puede realizar.

| Campo        | Tipo    | Descripción                    |
|-------------|---------|--------------------------------|
| Title       | string  | Título de la misión            |
| Introduction| string  | Texto de introducción         |
| Events      | []Event | Secuencia de eventos          |

**Relaciones:**
- **Tiene muchos** → `Event` (1:N).
- **Usada por** → `GameMaster` (quien ofrece/inicia misiones).

---

### 6. Event (Evento)

**Ubicación:** `models/quest.go`

Un paso o escena dentro de una misión (narración, encuentro, combate, etc.).

| Campo       | Tipo    | Descripción                         |
|------------|---------|-------------------------------------|
| Description| string  | Descripción del evento              |
| Enemies    | []Enemy | Enemigos presentes en el evento     |

**Relaciones:**
- **Pertenece a** → `Quest` (N:1).
- **Tiene muchos** → `Enemy` (1:N).  
  *(Nota: el tipo `Enemy` está referenciado en código pero no definido aún en `models/`; podría ser un `Character` o un struct específico.)*

---

### 7. Enemy (Enemigo)

**Estado:** Referenciado en `Event` (`Enemies []Enemy`), pero el tipo **no está definido** en el proyecto.

Opciones de diseño:
- **A)** Definir `Enemy` como un struct con datos mínimos (nombre, atributos, loot).
- **B)** Usar `Character` como enemigo y tipar `Enemies` como `[]*Character`.

**Relaciones (previstas):**
- **Pertenece a** → `Event` (N:1).

---

### 8. GameMaster (Máster de Juego)

**Ubicación:** `core/game_master.go`

Entidad que orquesta misiones y flujo de juego (estilo director de partida).

| Campo  | Tipo         | Descripción     |
|--------|--------------|-----------------|
| Quests | []models.Quest | Misiones disponibles |

**Comportamiento:**
- `StartQuest(quest models.Quest)` — Inicia una misión (actualmente llama `quest.Start()`, que debe existir en el modelo).

**Relaciones:**
- **Tiene muchos** → `Quest` (1:N).
- **Inicia** → `Quest` (acción sobre la entidad).

---

## Resumen de Relaciones

| Entidad     | Relación   | Con entidad  | Cardinalidad |
|------------|------------|-------------|---------------|
| Character  | tiene      | Attributes  | 1:1           |
| Game       | tiene      | Round       | 1:N           |
| Quest      | tiene      | Event       | 1:N           |
| Event      | tiene      | Enemy       | 1:N           |
| GameMaster | tiene      | Quest       | 1:N           |
| GameMaster | inicia     | Quest       | acción        |

---

## Documentación Relacionada

- **Mundo y tono:** `docs/game-context.md` (Valethra).
- **Combate:** `game-rules.md` (raíz) y/o `docs/fight-system.md`.
- **Game Master:** `docs/game-master.md`.

---

*Última actualización según estado del código en `models/` y `core/`.*
