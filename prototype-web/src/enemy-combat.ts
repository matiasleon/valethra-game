export const ENEMY_WINDUP = 0.9;
export const ENEMY_RECOVERY = 1.18;
export const ENEMY_DEATH_DURATION = 1.8;
export const ENEMY_DAMAGE = 15;
export const ENEMY_SPAWNS = [[-10, 3], [9, -9], [-5, -22], [5, -27]] as const;

export interface EnemyCombatState {
  health: number;
  attackCooldown: number;
  attackWindup: number;
  hitFlash: number;
  stagger: number;
  deathTime: number;
  alive: boolean;
}

export function createEnemyState(): EnemyCombatState {
  return { health: 90, attackCooldown: 0, attackWindup: 0, hitFlash: 0, stagger: 0, deathTime: 0, alive: true };
}

/** One transition per tick. A strike is emitted once when anticipation expires. */
export function advanceEnemy(state: EnemyCombatState, delta: number, distance: number): "none" | "windup" | "strike" {
  const step = Math.max(0, delta);
  if (!state.alive) { state.deathTime = Math.max(0, state.deathTime - step); return "none"; }
  state.attackCooldown = Math.max(0, state.attackCooldown - step);
  state.hitFlash = Math.max(0, state.hitFlash - step);
  state.stagger = Math.max(0, state.stagger - step);
  if (state.attackWindup > 0) {
    state.attackWindup = Math.max(0, state.attackWindup - step);
    if (state.attackWindup === 0) { state.attackCooldown = ENEMY_RECOVERY; return "strike"; }
  } else if (distance < 2.05 && state.attackCooldown === 0 && state.stagger === 0) {
    state.attackWindup = ENEMY_WINDUP;
    return "windup";
  }
  return "none";
}
