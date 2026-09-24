# V2 Three.js prototype

Estado verificado: septiembre de 2026.

## Alcance implementado

`prototype-web/` es un vertical slice 3D en tercera persona construido con TypeScript, Three.js y HTML/CSS. El nivel **El Santuario del Umbral** prueba exploración libre, combate en tiempo real, recursos del personaje, narrativa ambiental y una decisión final con dos consecuencias explícitas.

La secuencia narrativa es local al prototipo: Aldric vuelve a su pueblo después de ocho meses fuera, encuentra bloqueado el último camino, restaura tres sellos y descubre una familia detrás del portón. El carro orientado hacia el pueblo, la campana cortada desde dentro, los escudos de la guardia y una cuna vacía permiten inferir una retirada sin fijar todavía el suceso en el canon global.

## Bucle jugable

```text
Explorar el santuario
  → encontrar un sello y una pista ambiental
  → combatir o evitar trolls
  → administrar vida y resistencia
  → restaurar el sello
  → repetir hasta abrir la decisión del portón
```

Atacar cuesta 18 de resistencia y esquivar 28. Bloquear convierte un golpe de 15 puntos en daño reducido a cambio de resistencia. Restaurar un sello recupera vida y resistencia. La muerte reinicia el nivel. El portón sólo acepta la decisión después de restaurar los tres sellos.

## Arquitectura

```text
index.html + styles.css     HUD, introducción y decisiones
main.ts                     coordinación de fases
third-person-level.ts        coordinación de input, movimiento y combate
combat-rules.ts             reglas puras y testeables
ambient-audio.ts            música y viento procedurales temporales
qa-recorder.ts              captura automatizada del canvas a 60 FPS
```

Three.js mantiene un único mundo 3D continuo. El HUD y las decisiones permanecen en el DOM por accesibilidad. El modo QA usa la misma lógica que el jugador, con una ruta automática determinista para repetir grabaciones comparables.

## Decisiones técnicas

- Se conservó el stack pequeño: Three.js, TypeScript, Vite y Vitest; no se agregaron React, ECS, motor de físicas ni backend.
- Las texturas generadas son mapas de color base. Los materiales usan el modelo PBR de Three.js, pero no se afirma que existan mapas separados de normales, rugosidad o desplazamiento.
- La música procedural es un placeholder documentado porque Lyria no está disponible en el entorno actual.
- La V1 en Go/Ebitengine sigue siendo el juego estable; esta V2 continúa como experimento.

## Fuera de alcance

- Mapa regional y viajes entre niveles.
- Sucesos autónomos del mundo y consecuencias diferidas.
- Guardado, inventario, árboles de habilidades o físicas rígidas.
- Incorporación del santuario al canon global.

## Tercera persona (septiembre de 2026)

La posición física de Aldric es independiente de la cámara. `third-person-camera.ts` sigue un pivote amortiguado y usa tres rayos para retraerse ante muros, carro y portón; se recupera gradualmente al despejarse el encuadre. `aldric-actor.ts` contiene el modelo articulado procedural y animaciones de marcha, respiración, ataque, guardia y esquiva. `locomotion.ts` concentra amortiguación independiente del framerate y giro por el arco más corto.

El movimiento acelera y frena de forma gradual y conserva velocidades, costos y daño. La cámara se controla con mouse capturado o flechas; la dirección del ataque sigue la cámara y no su posición. Los golpes y la activación de sellos requieren un camino libre de obstáculos. Escape, pérdida de foco y cambio de pestaña pausan la partida. Los controles táctiles permiten recorrer el mismo nivel sin teclado. No se agregaron escenas, historia ni dependencias.

El rig es estilizado y procedural: no es un modelo con animaciones de captura de movimiento. La música procedural y el alcance narrativo local siguen siendo los del prototipo.


## Revisión de responsabilidades y animación

- `sanctuary-world.ts` posee escenario, sellos y obstáculos identificados, incluido el portón; reinicia su estado visual explícitamente.
- `enemy-combat.ts` contiene estado y transiciones temporales sin importar Three.js ni DOM. Emite anticipación o impacto una sola vez.
- `enemy-pose.ts` transforma un estado de combate de sólo lectura en una pose absoluta. Anticipación, ataque y recuperación comparten los tiempos del daño.
- `troll-actor.ts` construye el rig y aplica poses. Brazos articulados desde el hombro y arma sujeta a la mano; el desplazamiento visual nunca altera el origen de combate.
- `aldric-equipment.ts` posee espada, escudo y feedback de bloqueo.
- `level-contracts.ts` permite a audio y HUD consumir contratos sin importar el coordinador concreto.
- `scene-resources.ts` libera geometrías, materiales y texturas compartidas una sola vez. El nivel cancela su bucle, observadores y eventos al desmontarse.

Se retiró del web el grafo no referenciado del puente (`bridge-scene`, `exploration`, `encounter-rules`, sus pruebas y retratos duplicados). Godot conserva el prototipo histórico de Grukh. Por pedido de dirección artística, los cuatro enemigos del santuario ahora se representan como trolls. Este cambio visual local no establece rasgos morales de su especie ni modifica el canon global.

Los enemigos ahora respetan los obstáculos al perseguir y al recibir empuje. No se añadió navegación de rutas: pueden quedar detrás de un muro hasta que el jugador se acerque por un paso libre. La grabación QA está limitada a tres minutos y libera sus pistas al terminar.

El sistema usa composición y contratos pequeños; no introduce ECS, jerarquías de motores, plugins de entidades ni dependencias para necesidades hipotéticas. El coordinador aún contiene flujo del nivel y locomoción: una futura segunda escena sería el punto para extraer un controlador compartido. No se afirma haber validado rendimiento con cientos de enemigos.

`troll-death.ts` define fases temporales puras para la muerte. `TrollActor` interpola desde la pose del golpe, separa el arma y corrige el contacto con el suelo sin cambiar la escala. El cadáver conserva el origen y orientación de combate hasta el reinicio. El daño, la cantidad de encuentros y la historia permanecen iguales.
