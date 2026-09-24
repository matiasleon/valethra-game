import { describe, expect, it, vi } from "vitest";
import * as THREE from "three";
import { disposeScene } from "./scene-resources";
import { SanctuaryWorld } from "./sanctuary-world";
import { TrollActor } from "./troll-actor";
import { createEnemyState } from "./enemy-combat";

describe("scene ownership", () => {
  it("releases shared geometry, materials and textures exactly once", () => {
    const scene = new THREE.Scene();
    const texture = new THREE.Texture();
    const material = new THREE.MeshStandardMaterial({ map: texture });
    const geometry = new THREE.BoxGeometry();
    const disposeTexture = vi.spyOn(texture, "dispose");
    const disposeMaterial = vi.spyOn(material, "dispose");
    const disposeGeometry = vi.spyOn(geometry, "dispose");
    scene.add(new THREE.Mesh(geometry, material), new THREE.Mesh(geometry, material));
    disposeScene(scene); disposeScene(scene);
    for (const dispose of [disposeTexture, disposeMaterial, disposeGeometry]) expect(dispose).toHaveBeenCalledTimes(1);
    expect(scene.children).toHaveLength(0);
  });
  it("restarting closes an open gate and resets all wards immediately", () => {
    const scene = new THREE.Scene();
    const texture = new THREE.Texture();
    const world = new SanctuaryWorld(scene, { stone: texture, earth: texture, wood: texture });
    world.doors.position.y = 5.6;
    world.family.visible = true;
    world.wards.forEach(ward => { ward.active = true; ward.core.material.emissiveIntensity = 3.2; });
    expect(world.collides(0, -33.7, 0.38, false)).toBe(true);
    expect(world.collides(0, -33.7, 0.38, true)).toBe(false);
    expect(world.collides(-18, 5, 0.38, true)).toBe(true);
    world.reset();
    expect(world.doors.position.y).toBe(0);
    expect(world.family.visible).toBe(false);
    expect(world.wards.every(ward => !ward.active && ward.core.material.emissiveIntensity === 0.18)).toBe(true);
    disposeScene(scene);
  });
  it("animating/falling an enemy never moves its combat origin and resets for replay", () => {
    const actor = new TrollActor(0);
    actor.root.position.set(1, 0, 2);
    const state = createEnemyState(); state.attackWindup = 0.2;
    actor.update(state, 4, 1 / 60, 0.2, true);
    expect(actor.root.position.toArray()).toEqual([1, 0, 2]);
    state.alive = false; state.deathTime = 0;
    actor.update(state, 5, 1 / 60, 0, false);
    expect(actor.root.visible).toBe(true);
    actor.reset(); actor.update(createEnemyState(), 0, 0, 0, false);
    expect(actor.root.visible).toBe(true);
    const scene = new THREE.Scene(); scene.add(actor.root); disposeScene(scene);
  });
});
