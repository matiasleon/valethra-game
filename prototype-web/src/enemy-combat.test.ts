import { describe, expect, it } from "vitest";
import { advanceEnemy, createEnemyState, ENEMY_RECOVERY, ENEMY_WINDUP } from "./enemy-combat";
import { sampleEnemyPose } from "./enemy-pose";

describe("enemy combat and animation contract", () => {
  it("telegraphs before dealing one strike and enters recovery", () => {
    const state = createEnemyState();
    expect(advanceEnemy(state, 0, 1)).toBe("windup");
    expect(state.attackWindup).toBe(ENEMY_WINDUP);
    expect(advanceEnemy(state, 0.89, 1)).toBe("none");
    expect(advanceEnemy(state, 0.02, 1)).toBe("strike");
    expect(state.attackCooldown).toBe(ENEMY_RECOVERY);
    expect(advanceEnemy(state, 0.01, 1)).toBe("none");
  });
  it("does not start a new attack while staggered, distant or dead", () => {
    const state = createEnemyState();
    expect(advanceEnemy(state, 0.1, 5)).toBe("none");
    state.stagger = 0.3;
    expect(advanceEnemy(state, 0.1, 1)).toBe("none");
    state.alive = false; state.deathTime = 0.62;
    expect(advanceEnemy(state, 2, 1)).toBe("none");
    expect(state.deathTime).toBe(0);
  });
  it.each([30, 60, 120])("emits the strike once at %i FPS", (fps) => {
    const state = createEnemyState();
    advanceEnemy(state, 0, 1);
    let strikes = 0;
    for (let i = 0; i < fps; i++) if (advanceEnemy(state, 1 / fps, 1) === "strike") strikes++;
    expect(strikes).toBe(1);
  });
  it("keeps pose continuous across the damage tick and recovers to idle", () => {
    const state = createEnemyState();
    state.attackWindup = 0.000001;
    const before = sampleEnemyPose(state);
    advanceEnemy(state, 0.001, 1);
    const impact = sampleEnemyPose(state);
    expect(before.arm).toBeCloseTo(impact.arm, 5);
    expect(before.wrist).toBeCloseTo(impact.wrist, 5);
    expect(before.twist).toBeCloseTo(impact.twist, 5);
    expect(before.lean).toBeCloseTo(impact.lean, 5);
    advanceEnemy(state, 0.6, 3);
    expect(sampleEnemyPose(state).arm).toBeCloseTo(0);
  });
  it("sampling is deterministic, finite, and cannot mutate combat state", () => {
    const state = Object.freeze({ ...createEnemyState(), attackWindup: 0.3, hitFlash: 0.1 });
    const first = sampleEnemyPose(state);
    expect(sampleEnemyPose(state)).toEqual(first);
    expect(Object.values(first).every(Number.isFinite)).toBe(true);
    expect(state.attackWindup).toBe(0.3);
  });
});
