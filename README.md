# Valethra

Prototipo de RPG 2D táctico escrito en Go y Ebitengine. Valethra combina combate determinista por turnos con una fantasía adulta donde las decisiones morales tienen un costo real.

El juego actual permite recorrer dos misiones como Aldric, combatir varios enemigos y ver transiciones narrativas de victoria o derrota. El diseño de mundo abierto documentado en el repositorio es una dirección futura, no funcionalidad ya implementada.

Además, `prototype-web/` contiene el vertical slice V2 experimental construido con TypeScript y Three.js. Es un nivel 3D en tercera persona: tras ocho meses fuera, Aldric debe atravesar un santuario fronterizo, combatir apariciones, restaurar tres sellos y decidir si abre el camino para una familia. Todavía no reemplaza la V1 en Go. El corte anterior en Godot permanece en `prototype-v2/` como referencia de comparación.

## Requisitos

- Go 1.25 o superior (línea soportada y compatible con el stack estable)
- Dependencias de escritorio de Ebitengine para tu sistema
- Node.js 22.12 o superior para ejecutar el prototipo web V2
- Godot 4 sólo para ejecutar el prototipo visual anterior

Consulta la [guía oficial de instalación de Ebitengine](https://ebitengine.org/en/documents/install.html) si la compilación gráfica requiere librerías del sistema.

## Inicio rápido

```bash
git clone https://github.com/matiasleon/valethra-game.git
cd valethra-game
make setup
make check
make run
```

| Comando | Propósito |
| --- | --- |
| `make setup` | Descarga y verifica dependencias |
| `make test` | Ejecuta tests con detección de carreras |
| `make lint` | Verifica formato, `go vet` y consistencia del módulo |
| `make build` | Compila el juego en `bin/valethra` |
| `make check` | Quality gate local completo |
| `make run` | Inicia el juego |
| `make v2-setup` | Instala las dependencias bloqueadas del prototipo web |
| `make v2-run` | Inicia el prototipo V2 con Three.js |
| `make v2-check` | Prueba, verifica y compila el prototipo Three.js |
| `make godot-v2-run` | Inicia el corte anterior en Godot |

## Controles

| Tecla | Acción |
| --- | --- |
| `SPACE` / `ENTER` / `A` | Atacar o continuar |
| `Q` / `ESC` | Salir |

## Estructura

```text
core/               Dominio y orquestación, sin dependencias gráficas
game/               Adaptador de presentación con Ebitengine
game/assets/        Recursos embebidos usados en runtime
assets/             Fuentes y recursos artísticos de trabajo
docs/architecture/  Arquitectura vigente y propuestas futuras
docs/world/         Canon y proceso de evolución del mundo
docs/plans/         Trabajo propuesto, listo para humanos o agentes
prototype-v2/       Vertical slice experimental en Godot
prototype-web/      Vertical slice V2 experimental con Three.js
```

La documentación comienza en [`docs/README.md`](docs/README.md). Las reglas para agentes están en [`AGENTS.md`](AGENTS.md), y las convenciones para contribuir en [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Estado

Valethra está en etapa de prototipo. El foco inmediato es consolidar el vertical slice actual antes de implementar la propuesta de mundo abierto. El backlog vivo se mantiene en [`docs/plans/README.md`](docs/plans/README.md).

## Licencia y assets

El código aún no incluye un archivo de licencia. No asumas una licencia hasta que se agregue una explícitamente. Los assets conservan sus condiciones de origen; antes de distribuir una build pública hay que completar su inventario y atribución.
