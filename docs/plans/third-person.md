# Valethra: tercera persona

Alcance: transformar el santuario existente de `prototype-web/`. No agregar niveles, personajes, historia ni canon. La V1 de Go y el experimento anterior de Godot quedan como referencias.

## Plan previo a implementación

1. Separar la posición física de Aldric de la cámara. Todas las distancias de combate, IA, sellos y narrativa se calculan desde el personaje.
2. Extraer un controlador orbital con amortiguación exponencial, límites verticales y obstrucción por geometría del escenario.
3. Extraer un actor articulado de Aldric: silueta de viajero, capa oscura, cuero, acero y bronce; marcha, respiración, giro, ataque, guardia y esquiva.
4. Suavizar velocidad y orientación sin cambiar los costos, daño, enemigos, sellos ni finales. Agregar pausa segura y controles de cámara alternativos.
5. Validar reglas, movimiento y cámara; compilar web y ejecutar `make check`. Inspeccionar el juego en navegador y publicar en el Site existente conservando su acceso público.

## Arquitectura prevista

- `third-person-level.ts`: coordinación del nivel y reglas existentes.
- `third-person-camera.ts`: cámara, seguimiento y colisiones visuales.
- `aldric-actor.ts`: modelo y animación de presentación.
- `locomotion.ts`: funciones puras de suavizado y orientación.
- `fps-rules.ts`: reglas existentes, sin cambios de balance.

No se requiere cambiar motor ni agregar dependencias. Three.js ya proporciona geometría, jerarquías, luces y raycasting.

## Implementado y validado

Se completó la separación de posición física, cámara, locomoción y actor visual. Se añadieron controles táctiles y pausa, sin expandir el contenido del nivel. El alcance sigue siendo el prototipo web; no se presenta la V1 de Go como convertida a 3D.

Las pruebas cubren amortiguación entre 30 y 120 FPS, giros al cruzar ±π y retracción/recuperación de cámara sin desplazar el cuerpo. La revisión de navegador verificó el personaje, pausa, restauración de los tres sellos y ambos finales (abrir y sellar), sin errores de consola. `make check`, `make v2-check` (22 tests) y `make docs-check` pasaron. Las grabaciones antiguas se conservan como evidencia histórica, sin atribuirlas a esta versión.
