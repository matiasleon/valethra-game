# Juego

Cada jugador tiene un personaje (Character). El juego esta modelado por un sistema de rondas, 

# Reglas de Ataque

## Sistema de Defensa

Cuando un personaje ataca a otro, el ataque es defendido por el personaje objetivo.

La defensa consta únicamente de:
- **Armor (Armadura)**: Valor que se va consiguiendo durante el juego y se reduce con los ataques

## Mecánica de Ataque

1. **Si el ataque es mayor o igual a la armadura (y armor > 0):**
   - El personaje objetivo **muere** (Health = 0)
   - La armadura se reduce a 0

2. **Si el ataque es menor a la armadura:**
   - El ataque **impacta la armadura** (se resta del Armor)
   - No hay daño al Health

3. **Si el personaje no tiene armadura (Armor = 0):**
   - El ataque **va directamente al Health** (daño progresivo)
   - El Health se reduce por el valor completo del ataque

## Resumen

- **Armor**: Se reduce con cada ataque. Si el ataque iguala o supera la armadura, el personaje muere.
- **Health**: 
  - Se afecta si el ataque supera la armadura (muerte instantánea)
  - Se afecta progresivamente si Armor = 0 (el ataque va directamente a Health)