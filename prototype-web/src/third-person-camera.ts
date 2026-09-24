import * as THREE from "three";
import { damping } from "./locomotion";

/** Presentation only: moving the camera never moves the player or combat origin. */
export class ThirdPersonCamera {
  private readonly pivot = new THREE.Vector3();
  private readonly desired = new THREE.Vector3();
  private readonly direction = new THREE.Vector3();
  private readonly ray = new THREE.Raycaster();
  private readonly probeOrigin = new THREE.Vector3();
  private readonly hits: THREE.Intersection[] = [];
  private distance = 5.2;
  private initialized = false;

  reset(): void { this.initialized = false; this.distance = 5.2; }

  update(camera: THREE.PerspectiveCamera, player: THREE.Vector3, yaw: number, pitch: number,
    solids: THREE.Object3D[], delta: number, sprinting: boolean): void {
    this.desired.copy(player);
    this.desired.y += 1.35;
    if (!this.initialized) this.pivot.copy(this.desired);
    else this.pivot.lerp(this.desired, damping(18, delta));
    this.direction.set(Math.sin(yaw) * Math.cos(pitch), Math.sin(pitch), Math.cos(yaw) * Math.cos(pitch));
    // Multiple parallel probes protect the near-plane edges, not only its center.
    let safeDistance = 5.2;
    for (const side of [-0.24, 0, 0.24]) {
      this.probeOrigin.set(Math.cos(yaw) * side, 0.15, -Math.sin(yaw) * side).add(this.pivot);
      this.ray.set(this.probeOrigin, this.direction);
      this.ray.far = 5.5;
      this.hits.length = 0;
      const hit = this.ray.intersectObjects(solids, true, this.hits)[0];
      if (hit) safeDistance = Math.min(safeDistance, Math.max(0.25, hit.distance - 0.35));
    }
    // Retract immediately at obstructions; ease outward when the view clears.
    this.distance = safeDistance < this.distance ? safeDistance : THREE.MathUtils.lerp(this.distance, safeDistance, damping(5, delta));
    camera.position.copy(this.pivot).addScaledVector(this.direction, this.distance);
    camera.position.y = Math.max(0.35, camera.position.y);
    camera.lookAt(this.pivot);
    camera.fov = THREE.MathUtils.lerp(camera.fov, sprinting ? 65 : 61, damping(4, delta));
    camera.updateProjectionMatrix();
    this.initialized = true;
  }
}
