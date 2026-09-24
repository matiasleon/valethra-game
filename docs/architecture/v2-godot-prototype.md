# V2 Godot prototype

Estado verificado: agosto de 2026.

## Alcance implementado

La V2 experimental vive bajo `prototype-v2/` y utiliza Godot con GDScript tipado. Es un vertical slice separado de la V1 en Go: presenta el encuentro con Grukh, permite elegir una de cuatro respuestas y resuelve el resultado de forma determinista.

Antes de elegir, el jugador no ve estadísticas, probabilidades, requisitos ni recomendaciones. Después de elegir, la pantalla de resultado muestra el modificador aplicado y compara el poder final cuando hubo combate.

## Regla implementada

```text
Poder = Ataque + Stamina
```

- Atacar directamente: Aldric tiene 55 y Grukh 65; Aldric pierde.
- Observar: Inteligencia 10 descubre una herida y suma 15 de Ataque; Aldric gana 70 a 65.
- Provocar: Agilidad 15 reduce en 15 la Stamina de Grukh; Aldric gana 55 a 50.
- Negociar: Voluntad 12 produce un acuerdo sin combate.

## Arquitectura

```text
main.tscn
  └─ main.gd (presentación, input y fases)
       └─ encounter_rules.gd (reglas deterministas puras)
```

La escena no tiene un gestor global, bus de eventos, ECS ni sistema genérico de quests. Las reglas están separadas de la UI para poder verificarlas en modo headless.

## Relación con otras propuestas

Este prototipo no implementa `architecture-v2-open-world.md`. La propuesta histórica de mundo abierto continúa documentada como no implementada y no describe esta V2 experimental.

## Fuera de alcance

- Mapa o viajes.
- Sucesos independientes del jugador.
- Consecuencias sobre situaciones posteriores.
- Guardado de partida.
- Combate animado o por rondas.
- Contenido adicional.
