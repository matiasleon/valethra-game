# Valethra V2 prototype

Primer corte jugable en Godot. Implementa un único encuentro determinista entre Aldric y Grukh para validar que una decisión pueda modificar el contexto de una confrontación.

## Alcance implementado

- Presentación narrativa sin revelar estadísticas de Grukh.
- Cuatro decisiones neutrales y sin pistas mecánicas previas.
- Resolución determinista basada en Ataque, Stamina y los atributos de Aldric.
- Pantalla posterior que explica el modificador y la consecuencia.
- Reinicio inmediato para probar todas las rutas.
- Escenario pintado y retratos transparentes integrados en una UI cinematográfica.

No incluye mapa, persistencia, simulación del mundo, inventario ni consecuencias entre encuentros.

## Ejecutar

Requiere Godot 4.

```bash
make v2-run
```

Controles:

- `ENTER` o `SPACE`: continuar desde la introducción.
- `1`–`4`: elegir una respuesta.
- `R`: reiniciar después del resultado.
- `ESC`: salir.

## Validar

```bash
make v2-check
```

`v2-test` comprueba las cuatro resoluciones, construye cada estado visual y verifica que la pantalla de decisiones no revele pistas mecánicas. `v2-smoke` carga la escena principal durante dos frames en modo headless para detectar errores de importación, parseo o inicialización.
