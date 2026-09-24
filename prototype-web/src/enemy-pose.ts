import { ENEMY_RECOVERY, ENEMY_WINDUP, type EnemyCombatState } from "./enemy-combat";

const smooth = (value: number) => { const t = Math.max(0, Math.min(1, value)); return t * t * (3 - 2 * t); };

/** Absolute poses avoid accumulated rotations, and place the blade at impact on the damage tick. */
export function sampleEnemyPose(state: Readonly<EnemyCombatState>) {
  const progress = state.attackWindup > 0 ? 1 - state.attackWindup / ENEMY_WINDUP : 0;
  const preparation = smooth(progress / 0.72);
  const swing = smooth((progress - 0.72) / 0.28);
  const recovery = state.attackCooldown > 0 ? 1 - smooth((ENEMY_RECOVERY - state.attackCooldown) / 0.5) : 0;
  const attacking = state.attackWindup > 0;
  return {
    arm: attacking ? -2.2 * preparation + 1.25 * swing : -0.95 * recovery,
    wrist: 0.2 + 1.6 * (attacking ? swing : recovery),
    twist: attacking ? -0.38 * preparation + 0.66 * swing : 0.28 * recovery,
    lean: attacking ? -0.12 * preparation + 0.4 * swing : 0.28 * recovery,
    warning: attacking ? preparation : 0,
    recoil: smooth(state.hitFlash / 0.22),
    dissolve: state.alive ? 1 : Math.max(0, state.deathTime / 0.62),
  };
}
