# Arquitectura actual

Estado verificado: agosto de 2026.

## Alcance implementado

Valethra es hoy un vertical slice de escritorio. `main.go` compone al héroe, dos quests y los textos narrativos. `game.Game` administra una secuencia lineal de escenas. Los encuentros usan acciones polimórficas, aunque solo el combate tiene integración visual completa.

```text
main.go
  └─ game.Game (ciclo de Ebitengine y escenas)
       ├─ core.GameMaster (orquestación de acciones)
       ├─ core/ports (contratos consumidos por core)
       └─ core/entities (personajes, quests y acciones)
```

Las dependencias de proyecto fluyen desde `main` y `game` hacia `core`. `core` no importa Ebitengine.

## Responsabilidades

| Área | Responsabilidad |
| --- | --- |
| `core/entities` | Estado y reglas del dominio: atributos, ataque, quest y acciones |
| `core/ports` | Contratos mínimos que usa la orquestación |
| `core/game_master.go` | Resolución de acciones sin conocimiento de la UI |
| `game` | Escenas, input, render, animación, layout y assets embebidos |
| `main.go` | Composición de la demo y arranque de Ebitengine |

## Decisiones vigentes

- El combate es determinista y reduce armadura antes que salud.
- Una quest contiene eventos; cada evento contiene acciones mediante `ActionEvent`.
- La UI conserva una máquina de estados lineal orientada a la demo.
- Los PNG bajo `game/assets` se embeben en el ejecutable.
- No hay persistencia, red, generación procedural ni mundo abierto implementados.

## Límites conocidos

- El contenido narrativo está compuesto en `main.go`; conviene extraerlo cuando exista una segunda fuente de contenido real.
- `DialogueAction` existe en dominio, pero no tiene una escena dedicada.
- La propuesta v2 es deliberadamente aspiracional y debe implementarse en cortes verticales pequeños.
