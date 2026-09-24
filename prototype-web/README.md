# Valethra web prototype

Vertical slice experimental de Valethra construido con TypeScript, Three.js y HTML/CSS. Presenta un único nivel 3D en tercera persona: **El Santuario del Umbral**.

Aldric regresa a su pueblo después de ocho meses en una misión peligrosa. El último camino atraviesa un santuario abandonado. Para abrir su portón debe explorar las ruinas, restaurar tres sellos y sobrevivir a cuatro apariciones. El carro roto, los escudos abandonados, una campana silenciada y una cuna vacía cuentan qué ocurrió sin convertirlo en una explicación explícita.

## Controles

- `WASD`: movimiento.
- Mouse: cámara.
- `Shift`: correr.
- Click izquierdo o `Espacio`: atacar.
- Click derecho o `F`: bloquear.
- `Q`: esquivar.
- `E`: interactuar.

Atacar, correr, bloquear y esquivar comparten la resistencia. Restaurar un sello recupera recursos y funciona como punto de progreso. Con los tres sellos activos, el jugador decide entre abrir el portón para la familia o mantener la amenaza contenida.

## Ejecutar y validar

```bash
make v2-setup
make v2-run
make v2-check
```

El servidor abre el juego en `http://127.0.0.1:4173`.

La ruta `/?autoplay=1&record=1` ejecuta una partida automática y graba el canvas a 60 FPS. Las grabaciones se conservan localmente y no se versionan. Guardá una captura como `qa/valethra-autoplay-60fps.webm`; `qa/review.html?t=8` permite inspeccionar un instante concreto.

## Assets y límites

- Las texturas de piedra, tierra y madera fueron generadas como mapas de color base y se aplican a materiales PBR de Three.js con rugosidad y metalicidad configuradas en runtime.
- La ambientación sonora actual es procedural mediante Web Audio. Es un reemplazo temporal y no una composición de Lyria.
- No hay backend, guardado, físicas rígidas, inventario ni simulación regional.
- El contenido narrativo del santuario es local al prototipo y no modifica el canon global de Valethra.

## Presentación en tercera persona

Aldric permanece visible y la cámara orbital se acerca ante obstáculos. El movimiento acelera, frena y gira gradualmente; espada, escudo y cuerpo se animan juntos. Las flechas permiten orientar la cámara sin capturar el mouse. `Esc` pausa; perder foco también pausa. En pantallas pequeñas y dispositivos táctiles aparecen botones de movimiento, cámara y acciones. El mouse se captura al hacer click dentro del mundo.

La cámara, el actor visual y la locomoción se separan de las reglas del santuario. No se agregaron escenas ni historia. Las grabaciones anteriores en `qa/` corresponden a la versión previa en primera persona; no son evidencia de esta actualización.


## Código y animaciones

El santuario separa escenario, contratos, reglas de combate, poses y actores visuales. Las apariciones anticipan el ataque, golpean con el arma vinculada a la mano, recuperan la postura y reaccionan al daño; la pose de impacto coincide con el temporizador del combate. No se agregaron encuentros ni historias. Los archivos web del puente anterior y sus retratos duplicados fueron retirados porque no tenían consumidores; el experimento Godot permanece en `../prototype-v2/`.

`make v2-check` cubre reglas, temporización de enemigos, continuidad de poses, cámara, reinicio del mundo y liberación de recursos compartidos. `/?autoplay=1&record=1` verifica el final abierto y permite descargar una grabación; `/?autoplay=1&ending=seal` recorre el final sellado. Las grabaciones se detienen como máximo a los tres minutos.
