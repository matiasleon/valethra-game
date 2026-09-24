import { describe, expect, it } from "vitest";
import * as THREE from "three";
import { TrollActor } from "./troll-actor";
import { createEnemyState, ENEMY_DEATH_DURATION } from "./enemy-combat";
import { sampleTrollDeath } from "./troll-death";
import { disposeScene } from "./scene-resources";

describe("troll death", () => {
  it("starts at the struck pose, collapses and settles without a discontinuity", () => {
    expect(sampleTrollDeath(ENEMY_DEATH_DURATION).blend).toBe(0);
    expect(sampleTrollDeath(0).roll).toBeCloseTo(-Math.PI / 2);
    expect(sampleTrollDeath(0).settle).toBe(1);
    for (const time of [0.12, 0.22, 0.42, 0.54, 0.86, 1.14, 1.8]) {
      const a = sampleTrollDeath(ENEMY_DEATH_DURATION - time + 0.00001);
      const b = sampleTrollDeath(ENEMY_DEATH_DURATION - time - 0.00001);
      for (const key of Object.keys(a) as Array<keyof typeof a>) expect(Math.abs(a[key] - b[key])).toBeLessThan(0.001);
    }
  });
  it.each([30, 60, 120])("leaves a stable full-size corpse above the ground at %i FPS", (fps) => {
    const actor = new TrollActor(0);
    const state = createEnemyState();
    actor.root.position.set(3, 0, -4);
    actor.update(state, 1, 1 / fps, 0.2, true);
    state.alive = false;
    for (let i = 0; i <= Math.ceil(ENEMY_DEATH_DURATION * fps); i++) {
      state.deathTime = Math.max(0, ENEMY_DEATH_DURATION - i / fps);
      actor.update(state, 1 + i / fps, 1 / fps, 2, false);
    }
    expect(actor.root.visible).toBe(true);
    expect(actor.root.position.toArray()).toEqual([3, 0, -4]);
    actor.root.traverse(part => {
      if (part instanceof THREE.Group) for (const scale of part.scale.toArray()) expect(scale).toBeCloseTo(1, 10);
    });
    actor.root.updateMatrixWorld(true);
    const before = new THREE.Box3().setFromObject(actor.root);
    // Exclude the hidden warning ring only via a small tolerance; neither body nor club sinks.
    expect(before.min.y).toBeGreaterThanOrEqual(-0.03);
    const rotation = actor.root.rotation.y;
    actor.update(state, 10, 1, -2, false);
    expect(actor.root.rotation.y).toBe(rotation);
    expect(new THREE.Box3().setFromObject(actor.root).min.distanceTo(before.min)).toBeLessThan(0.0001);
    actor.reset(); actor.update(createEnemyState(), 0, 0, 0, false);
    expect(actor.root.visible).toBe(true);
    expect(actor.root.children.filter(child => child instanceof THREE.Group)).toHaveLength(1);
    const scene = new THREE.Scene(); scene.add(actor.root); disposeScene(scene);
  });
});
