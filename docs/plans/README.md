# Planes y backlog

Los planes de esta carpeta son neutrales respecto de la herramienta: pueden ejecutarlos una persona, Codex u otro agente. Un plan no describe funcionalidad disponible.

## Prioridad actual

1. **Asegurar el vertical slice**: ampliar tests de `GameMaster`, progresión entre eventos y acciones no-combate.
2. **Completar atribución de assets**: registrar autor, origen y licencia antes de distribuir builds.
3. **Separar contenido de composición**: extraer héroe, quests y textos de `main.go` cuando haya una segunda campaña o fuente de datos.
4. **Validar una escena de diálogo**: implementar un corte mínimo que demuestre el modelo `ActionEvent`.
5. **Probar una porción v2**: elegir un único slice de mundo abierto después de estabilizar la demo.

## Plantilla de tarea para agentes

```markdown
# Resultado esperado
Qué cambia para el jugador, el mundo o el mantenedor.

## Contexto y límites
Archivos relevantes, invariantes y cosas fuera de alcance.

## Criterios de aceptación
- Resultado observable.
- Casos límite.
- Documentación necesaria.

## Validación
Comandos concretos; normalmente `make check` y una prueba manual si cambia UI.
```

## Decisiones migradas de planes históricos

- El modelo `Event -> []ActionEvent` ya fue adoptado.
- La configuración centralizada de layout ya fue adoptada mediante `LayoutConfig` y `LayoutManager`.
- La propuesta de mundo abierto se conserva en `docs/architecture-v2-open-world.md` y sigue pendiente.
