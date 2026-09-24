import { describe, expect, it } from "vitest";
import * as THREE from "three";
import { damping, turnTowards } from "./locomotion";
import { ThirdPersonCamera } from "./third-person-camera";

describe("third person presentation", () => {
  it("damps consistently at 30 and 120 frames per second", () => {
    const run = (fps: number) => { let value = 0; for (let i = 0; i < fps; i++) value += (1 - value) * damping(5, 1 / fps); return value; };
    expect(run(30)).toBeCloseTo(run(120), 10);
    expect(damping(10, 0)).toBe(0);
  });
  it("turns across the angle seam via the shortest direction", () => {
    const result = turnTowards(Math.PI - 0.1, -Math.PI + 0.1, 1 / 60);
    expect(result).toBeGreaterThan(Math.PI - 0.1);
    expect(result).toBeLessThan(Math.PI + 0.1);
  });
  it("retracts before a wall, recovers smoothly and never moves the body", () => {
    const body = new THREE.Vector3(0, 0, 0);
    const camera = new THREE.PerspectiveCamera();
    const rig = new ThirdPersonCamera();
    rig.update(camera, body, 0, 0.2, [], 1 / 60, false);
    const originalDistance = camera.position.z;
    const wall = new THREE.Mesh(new THREE.BoxGeometry(8, 6, 0.5));
    wall.position.set(0, 2, 2); wall.updateMatrixWorld(true);
    rig.update(camera, body, 0, 0.2, [wall], 1 / 60, false);
    expect(camera.position.z).toBeLessThan(1.75);
    const compressed = camera.position.z;
    rig.update(camera, body, 0, 0.2, [], 1 / 60, false);
    expect(camera.position.z).toBeGreaterThan(compressed);
    expect(camera.position.z).toBeLessThan(originalDistance);
    expect(body.toArray()).toEqual([0, 0, 0]);
  });
});
