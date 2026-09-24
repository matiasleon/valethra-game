import * as THREE from "three";
import { damping, turnTowards } from "./locomotion";

/** A continuous, articulated silhouette; no gameplay rules live in this rig. */
export class AldricActor {
  readonly root = new THREE.Group();
  readonly swordHand = new THREE.Group();
  readonly shieldHand = new THREE.Group();
  private readonly torso = new THREE.Group();
  private readonly legs = [new THREE.Group(), new THREE.Group()];
  private readonly arms = [new THREE.Group(), new THREE.Group()];
  private readonly cape: THREE.Mesh;
  private stride = 0;
  private blend = 0;

  constructor(sword: THREE.Group, shield: THREE.Group) {
    const cloth = new THREE.MeshStandardMaterial({ color: 0x293d3b, roughness: 0.94 });
    const leather = new THREE.MeshStandardMaterial({ color: 0x49372e, roughness: 0.87 });
    const steel = new THREE.MeshStandardMaterial({ color: 0x73807e, metalness: 0.72, roughness: 0.4 });
    const bronze = new THREE.MeshStandardMaterial({ color: 0xa68c5c, metalness: 0.68, roughness: 0.48 });
    const skin = new THREE.MeshStandardMaterial({ color: 0xb99579, roughness: 0.88 });
    const hair = new THREE.MeshStandardMaterial({ color: 0x2c2421, roughness: 1 });
    const part = (parent: THREE.Object3D, geometry: THREE.BufferGeometry, material: THREE.Material,
      x: number, y: number, z: number) => {
      const mesh = new THREE.Mesh(geometry, material);
      mesh.position.set(x, y, z); mesh.castShadow = true; mesh.receiveShadow = true;
      parent.add(mesh); return mesh;
    };
    this.root.add(this.torso);
    this.torso.position.y = 1.08;
    part(this.torso, new THREE.CylinderGeometry(0.27, 0.22, 0.58, 12), cloth, 0, 0.17, 0).scale.z = 0.68;
    part(this.torso, new THREE.CylinderGeometry(0.25, 0.3, 0.32, 12), leather, 0, -0.22, 0).scale.z = 0.7;
    part(this.torso, new THREE.CylinderGeometry(0.255, 0.255, 0.09, 12), leather, 0, -0.05, 0).scale.z = 0.75;
    part(this.torso, new THREE.BoxGeometry(0.09, 0.08, 0.03), bronze, 0, -0.05, -0.2);
    part(this.torso, new THREE.CapsuleGeometry(0.15, 0.15, 6, 12), skin, 0, 0.62, -0.01);
    const hairCap = part(this.torso, new THREE.SphereGeometry(0.168, 12, 8), hair, 0, 0.72, 0.015);
    hairCap.scale.set(1, 0.68, 1);
    part(this.torso, new THREE.SphereGeometry(0.12, 10, 8), hair, 0, 0.57, 0.09).scale.set(1.1, 1, 0.6);
    for (let i = 0; i < 2; i++) {
      const side = i === 0 ? -1 : 1;
      const leg = this.legs[i]!; leg.position.set(side * 0.14, 0.83, 0); this.root.add(leg);
      part(leg, new THREE.CapsuleGeometry(0.105, 0.25, 4, 8), leather, 0, -0.18, 0);
      part(leg, new THREE.CapsuleGeometry(0.09, 0.22, 4, 8), cloth, 0, -0.52, 0);
      part(leg, new THREE.BoxGeometry(0.2, 0.18, 0.33), leather, 0, -0.73, -0.055);
      part(leg, new THREE.SphereGeometry(0.115, 8, 6), steel, 0, -0.39, -0.06).scale.z = 0.6;
      const arm = this.arms[i]!; arm.position.set(side * 0.33, 0.39, 0); this.torso.add(arm);
      part(arm, new THREE.SphereGeometry(0.18, 10, 8), steel, 0, -0.025, 0).scale.y = 0.64;
      part(arm, new THREE.CapsuleGeometry(0.075, 0.35, 4, 8), cloth, 0, -0.28, 0);
      part(arm, new THREE.CylinderGeometry(0.085, 0.07, 0.2, 8), leather, 0, -0.43, 0);
      part(arm, new THREE.SphereGeometry(0.075, 8, 6), leather, 0, -0.56, 0);
    }
    this.arms[1]!.add(this.swordHand); this.swordHand.position.y = -0.55;
    this.arms[0]!.add(this.shieldHand); this.shieldHand.position.y = -0.4;
    this.swordHand.add(sword); sword.position.set(0, 0, -0.08); sword.rotation.set(-1.25, 0, 0); sword.scale.setScalar(0.9);
    this.shieldHand.add(shield); shield.position.set(-0.08, 0, -0.08); shield.rotation.set(0, -0.55, 0); shield.scale.setScalar(0.82); shield.visible = true;
    // Tailored cape with tapered shoulders, open silhouette and a weathered hem.
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute("position", new THREE.Float32BufferAttribute([
      -0.27,0.43,0.14, 0.27,0.43,0.14, -0.39,-0.76,0.3,
      0.27,0.43,0.14, 0.39,-0.72,0.3, -0.39,-0.76,0.3,
    ],3)); geometry.computeVertexNormals();
    this.cape = part(this.torso, geometry, new THREE.MeshStandardMaterial({ color: 0x233334, roughness: 1, side: THREE.DoubleSide }),0,0,0);
    part(this.torso, new THREE.SphereGeometry(0.06, 8, 6), bronze, -0.21, 0.42, -0.14);
  }

  update(position: THREE.Vector3, facing: number, speed: number, elapsed: number, delta: number,
    attackTime: number, blocking: boolean, dodgeTime: number): void {
    this.root.position.copy(position);
    this.root.rotation.y = turnTowards(this.root.rotation.y, facing, delta);
    this.blend = THREE.MathUtils.lerp(this.blend, Math.min(speed / 4.1, 1), damping(12, delta));
    this.stride += speed * delta * 2.7;
    const step = Math.sin(this.stride) * this.blend;
    const swing = attackTime > 0 ? Math.sin((1 - attackTime / 0.34) * Math.PI) : 0;
    const dodge = dodgeTime > 0 ? Math.sin((1 - dodgeTime / 0.42) * Math.PI) : 0;
    this.legs[0]!.rotation.x = step * 0.65; this.legs[1]!.rotation.x = -step * 0.65;
    this.torso.position.y = 1.08 + Math.abs(step) * 0.035 + Math.sin(elapsed * 2.1) * 0.008 - dodge * 0.32;
    this.torso.rotation.x = -Math.min(speed / 6.7, 1) * 0.08 - dodge * 0.5;
    this.torso.rotation.y = -swing * 0.6;
    this.arms[1]!.rotation.set(-step * 0.28 - swing * 1.8, -swing * 0.8, swing * 0.6);
    this.arms[0]!.rotation.x = THREE.MathUtils.lerp(this.arms[0]!.rotation.x, blocking ? -1.25 : step * 0.28, damping(18, delta));
    this.cape.rotation.x = -0.06 - speed * 0.025 + Math.sin(elapsed * 4) * 0.025;
  }
}
