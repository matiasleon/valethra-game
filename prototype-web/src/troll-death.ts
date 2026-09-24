import { ENEMY_DEATH_DURATION } from "./enemy-combat";

const ease = (t: number) => { const x = Math.max(0, Math.min(1, t)); return x * x * (3 - 2 * x); };

/** Authored stages in seconds: recoil, knees yielding, fall, impact, settling. */
export function sampleTrollDeath(remaining: number) {
  const elapsed = Math.max(0, Math.min(ENEMY_DEATH_DURATION, ENEMY_DEATH_DURATION - remaining));
  const knees = ease((elapsed - 0.12) / 0.42);
  const fall = ease((elapsed - 0.42) / 0.72);
  const settle = ease((elapsed - 1.14) / 0.66);
  return {
    blend: ease(elapsed / 0.24),
    knees,
    fall,
    settle,
    roll: -Math.PI / 2 * fall,
    pitch: -0.16 * ease(elapsed / 0.18) + 0.34 * knees - 0.08 * settle,
    height: 1.15 - 0.42 * knees - 0.22 * fall,
    arm: 0.3 * knees + 0.55 * fall,
    drop: ease((elapsed - 0.22) / 0.64),
    impact: elapsed >= 1.14 ? Math.sin(Math.min(1, (elapsed - 1.14) / 0.25) * Math.PI) * (1 - settle) : 0,
  };
}
