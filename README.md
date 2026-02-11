# Valethra

Un RPG de combate táctico por turnos ambientado en un mundo de fantasía adulta donde el heroísmo es caro y la supervivencia es lo único que importa.

Juega como **Aldric**, un mercenario sin tierra que acepta los contratos que nadie más quiere. Cruza puentes infestados de trolls, adentrarte en bosques corrompidos, y enfrenta enemigos cada vez más peligrosos en un mundo que no recompensa la bondad — pero que la necesita.

---

## ¿De qué va esto?

Este proyecto nació de las ganas de hacer un juego que:

1. **Tenga combate con peso**  
   No es spam de clicks. Cada ataque cuenta. La armadura se rompe. Si tu salud llega a cero, perdés. No hay segundas oportunidades.

2. **Se sienta táctico, no aleatorio**  
   No hay RNG de daño. Si tenés 25 de ataque y el enemigo tiene 15 de armadura, le hacés 10 de daño. Simple, predecible, táctico.

3. **Tenga narrativa sin forzarla**  
   No hay cinemáticas largas. La historia está en los textos breves entre combates. Si querés skipearla, adelante. Si querés leerla, está bien escrita.

4. **Sea divertido de jugar**  
   Sprites animados, sonido de combate (próximamente), controles simples. Apretás una tecla y atacás. No hace falta menú de 10 niveles.

---

## ¿Qué tecnologías usa?

- **Go** — lenguaje principal
- **Ebiten** — engine 2D (parecido a Love2D, pero en Go)
- **Clean Architecture** — core de negocio separado de la UI, todo testeado

No hay motores pesados como Unity o Unreal. Es puro código y pixel art.

---

## Instalación

### Requisitos

- **Go 1.21+** — [Descargar aquí](https://go.dev/dl/)
- **Compilador C** (para Ebiten):
  - **macOS**: Ya viene con Xcode Command Line Tools
  - **Linux**: `sudo apt install gcc` o equivalente
  - **Windows**: Instalar [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)

### Clonar el repo

```bash
git clone https://github.com/tu-usuario/saturday-chill.git
cd saturday-chill
```

### Instalar dependencias

```bash
go mod download
```

---

## Cómo jugar

### Ejecutar el juego

```bash
go run main.go
```

Si todo anda bien, se abre una ventana de 800x600 con el juego corriendo.

### Controles

| Tecla             | Acción                          |
|-------------------|---------------------------------|
| `SPACE` / `ENTER` | Atacar / Continuar              |
| `A`               | Atacar (alternativo)            |
| `Q` / `ESC`       | Salir                           |

---

## Estructura del proyecto

```
saturday-chill/
├── main.go              # Entry point (crea el juego, corre Ebiten)
├── core/                # Lógica de negocio (independiente de UI)
│   ├── entities/        # Character, Quest, Attributes
│   ├── ports/           # Interfaces (CombatEngine)
│   └── game_master.go   # Orquestador de quests y combate
├── game/                # UI de Ebiten (escenas, sprites, animaciones)
│   ├── scene.go         # Interfaz Scene
│   ├── game.go          # SceneManager (maneja flujo de escenas)
│   ├── combat_scene.go  # Escena de combate (multi-evento)
│   ├── transition_scene.go # Escena de transiciones (intro/victoria/derrota)
│   ├── sprite.go        # Animador de sprites
│   └── assets.go        # Carga de imágenes
├── assets/              # Sprites (soldier/, orc/)
└── docs/                # Documentación del mundo y reglas
```

---

## Tests

El combate tiene cobertura de tests unitarios. Podés correrlos con:

```bash
go test ./core/entities
```

Esperá ver algo como:

```
PASS
ok      saturday-chill/core/entities    0.123s
```

---

## Roadmap (cosas que quiero agregar)

- [ ] Más enemigos (esqueletos, magos, dragones)
- [ ] Sistema de loot (armas, armaduras)
- [ ] Habilidades especiales (ataques críticos, hechizos)
- [ ] Sonido de combate
- [ ] Más quests (el mundo de Valethra tiene mucho más para contar)
- [ ] Combate en tiempo real con timing (experimental)

---

## Contribuir

Si te gusta el proyecto y querés aportar:

1. **Fork** el repo
2. Creá una **branch** para tu feature (`git checkout -b feature/nueva-cosa`)
3. **Commiteá** tus cambios (`git commit -m "Add nueva-cosa"`)
4. **Pusheá** la branch (`git push origin feature/nueva-cosa`)
5. Abrí un **Pull Request**

No hace falta que sea perfecto. Si tenés una idea, mandá el PR y charlamos.

---

## Licencia

Este proyecto está bajo licencia **MIT**. Hacé lo que quieras con el código, pero si hacés algo copado, avisame.

---

## Créditos

- **Sprites**: [itch.io](https://itch.io/game-assets/free) (revisar carpeta `assets/` para atribuciones específicas)
- **Engine**: [Ebiten](https://ebiten.org/)
- **Inspiración**: DoD, The Witcher (la novela, no el juego), Warhammer Fantasy

---

## Contacto

Si querés charlar sobre el proyecto, tirar ideas o reportar bugs:

- **GitHub Issues**: [Crear issue](https://github.com/tu-usuario/saturday-chill/issues)
- **Email**: tu-email@ejemplo.com (opcional)

---

**Hecho con ☕ y ganas de que exista un juego así.**
