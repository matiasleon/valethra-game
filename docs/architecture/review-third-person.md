# Revisión SOLID del prototipo web

24 de septiembre de 2026. Alcance: grafo de módulos desde `prototype-web/src/main.ts`, código sin consumidores del puente, actores, cámara, recursos y contratos. Implementación autorizada junto con la revisión. La V1 Go y el prototipo Godot se conservan.

## Hallazgos corregidos

| Prioridad | Hallazgo previo | Cambio y motivo |
| --- | --- | --- |
| P1 | El coordinador de casi 1.200 líneas construía mundo, equipo y enemigos además de ejecutar combate. | Separación por responsabilidad mediante composición: `SanctuaryWorld`, `AldricEquipment`, `ApparitionActor`; reglas y poses de enemigos sin DOM ni renderer. |
| P1 | No existía cierre de RAF, eventos, observador ni recursos GPU. | `dispose()` idempotente, eventos cancelables y liberación única de materiales/geometrías/texturas compartidos; invocación desde desmontaje y cierre de página. |
| P1 | Reiniciar después del final abierto no restablecía la altura de las hojas del portón. | Reinicio explícito en el propietario del mundo y prueba de regresión. |
| P1 | Los enemigos y el empuje de ataques podían atravesar obstáculos. | Colisiones del mismo mundo aplicadas a desplazamiento y empuje con radio de enemigo. |
| P2 | Audio importaba tipos del coordinador concreto; referencias visuales viajaban como `userData` y conversiones sin control. | Contratos independientes y referencias tipadas en propietarios del mundo y los rigs. |
| P2 | El puente web completo no era importado por ningún punto de entrada activo. | Eliminación de renderer, reglas, utilidades, pruebas exclusivas y tres retratos duplicados. Godot preserva la referencia histórica. |
| P2 | API `getCanvas`, temporizador `cameraShake`, mira oculta y poses iniciales sobrescritas estaban sin uso. | Eliminación y renombrado de reglas FPS a reglas de combate. |
| P2 | Grabación sin límite ni cierre de pistas y callbacks diferidos que sobrevivían al reinicio. | Límite de tres minutos, cierre de pistas y cancelación/identificación de partida para callbacks diferidos. |
| P2 | Ordenamiento y asignaciones temporales innecesarias cada frame. | Búsqueda lineal del sello más cercano y buffers de raycasting reutilizados en cámara. |

## Principios y límites

SRP y DIP mejoran con propietarios de recursos y contratos separados. La animación recibe estado de sólo lectura y no decide daño. No se agregaron jerarquías de herencia, interfaces sin consumidores ni una plataforma de entidades genérica: LSP no requiere una nueva jerarquía y OCP no justifica extensiones hipotéticas.

Las cadenas visibles se asignan con `textContent` o plantillas de texto local; no se encontró entrada remota interpolada en HTML. El prototipo no tiene backend, autenticación propia ni almacenamiento de secretos. Esto no constituye una auditoría de dependencias o infraestructura de Sites.

## Validación y trabajo posterior

Pruebas de ataque único a 30/60/120 FPS, continuidad de brazo/muñeca/torso al impacto, recuperación, estado inmutable para poses, reinicio del portón y liberación única de recursos. Se mantienen pruebas de cámara y reglas de vida/resistencia.

Antes de aumentar sustancialmente la cantidad de enemigos: medir frame time y memoria, agregar partición espacial si la medición lo exige y resolver navegación alrededor de muros. El coordinador conserva el flujo de esta única escena; extraer flujo compartido cuando exista una segunda escena autorizada. El juego sigue siendo un prototipo procedural, no una base probada para multijugador o centenares de entidades.

Resultado de validación: `make check`, `make v2-check` (18 pruebas vigentes) y `make docs-check` aprobados. Se retiraron las 14 pruebas exclusivas del puente sin consumidores y se agregaron 10 pruebas del sistema vigente. Navegador: ambos finales, reinicio, anticipación de ataque y descarga de grabación disponibles; sin errores de consola. Permanecen avisos no bloqueantes del backend Metal de Ebitengine y del tamaño del bundle Three.js.
