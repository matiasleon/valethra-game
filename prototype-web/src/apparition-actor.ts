import * as THREE from "three";
import type { EnemyCombatState } from "./enemy-combat";
import { sampleEnemyPose } from "./enemy-pose";
import { turnTowards } from "./locomotion";

/** Visual rig consumes combat state; it cannot deal damage or change gameplay timers. */
export class ApparitionActor {
  readonly root = new THREE.Group();
  private readonly body = new THREE.Group();
  private readonly weaponHand = new THREE.Group();
  private readonly arms: THREE.Group[] = [];
  private readonly shell: THREE.Mesh;
  private readonly core: THREE.Mesh<THREE.BufferGeometry, THREE.MeshStandardMaterial>;
  private readonly warning: THREE.Mesh<THREE.BufferGeometry, THREE.MeshBasicMaterial>;
  private readonly warningBeacon: THREE.Mesh;
  private readonly light: THREE.PointLight;

  constructor(private readonly index: number) {
      const group = this.body;
      const cloakMaterial = new THREE.MeshStandardMaterial({
        color: index % 2 === 0 ? 0x172329 : 0x211b27,
        emissive: 0x090d12,
        emissiveIntensity: 0.35,
        roughness: 0.82,
        metalness: 0.1,
        flatShading: true,
      });
      const shell = this.mesh(
        new THREE.ConeGeometry(0.78, 1.95, 7),
        cloakMaterial,
      );
      shell.position.y = 0.95;
      const hood = this.mesh(
        new THREE.ConeGeometry(0.52, 0.9, 7),
        cloakMaterial,
      );
      hood.position.y = 1.93;
      hood.rotation.x = -0.12;
      const face = this.mesh(
        new THREE.SphereGeometry(0.3, 10, 7),
        new THREE.MeshStandardMaterial({ color: 0x030506, roughness: 1 }),
      );
      face.position.set(0, 1.86, 0.2);
      face.scale.set(0.82, 1, 0.48);
      const eyeMaterial = new THREE.MeshBasicMaterial({ color: index % 2 === 0 ? 0xf5b068 : 0xd893ef });
      for (const side of [-1, 1]) {
        const eye = this.mesh(new THREE.SphereGeometry(0.035, 6, 4), eyeMaterial, false);
        eye.position.set(side * 0.105, 1.91, 0.355);
        group.add(eye);
      }
      const core = this.mesh(
        new THREE.OctahedronGeometry(0.18, 0),
        new THREE.MeshStandardMaterial({ color: 0xd6a16f, emissive: index % 2 === 0 ? 0xe36f38 : 0xa74ad5, emissiveIntensity: 2.8, roughness: 0.32 }),
        false,
      );
      core.position.set(0, 1.25, 0.66);
      const light = new THREE.PointLight(index % 2 === 0 ? 0xe36f38 : 0xa74ad5, 5, 4.5, 2);
      light.position.set(0, 1.3, 0.5);
      const armMaterial = new THREE.MeshStandardMaterial({ color: 0x192022, roughness: 0.85, flatShading: true });
      for (const side of [-1, 1]) {
        const shoulder = this.mesh(new THREE.DodecahedronGeometry(0.21, 0), armMaterial);
        shoulder.position.set(side * 0.49, 1.63, 0.04);
        shoulder.scale.set(1.35, 0.78, 1.05);
        const arm = this.mesh(new THREE.CylinderGeometry(0.11, 0.16, 1.2, 5), armMaterial);
        arm.position.set(side * 0.62, 1.28, 0.04);
        arm.rotation.z = side * 0.32;
        const pivot = new THREE.Group();
        pivot.position.set(side * 0.49, 1.63, 0.04);
        shoulder.position.set(0, 0, 0);
        arm.position.set(0, -0.45, 0);
        arm.rotation.z = 0;
        pivot.add(shoulder, arm);
        this.arms.push(pivot);
        group.add(pivot);
      }
      const enemyBladeShape = new THREE.Shape();
      enemyBladeShape.moveTo(-0.045, -0.42);
      enemyBladeShape.lineTo(0.045, -0.42);
      enemyBladeShape.lineTo(0.065, 0.36);
      enemyBladeShape.lineTo(0, 0.53);
      enemyBladeShape.lineTo(-0.065, 0.36);
      const weapon = this.mesh(
        new THREE.ExtrudeGeometry(enemyBladeShape, { depth: 0.035, bevelEnabled: true, bevelSize: 0.012, bevelThickness: 0.01, bevelSegments: 1 }),
        new THREE.MeshStandardMaterial({ color: 0x758080, emissive: 0x202a2d, emissiveIntensity: 0.5, roughness: 0.4, metalness: 0.75 }),
      );
      const weaponGuard = this.mesh(
        new THREE.BoxGeometry(0.34, 0.055, 0.09),
        new THREE.MeshStandardMaterial({ color: 0x7d684c, roughness: 0.48, metalness: 0.66 }),
      );
      weaponGuard.position.y = -0.45;
      const weaponGrip = this.mesh(
        new THREE.CylinderGeometry(0.045, 0.052, 0.36, 6),
        new THREE.MeshStandardMaterial({ color: 0x35231b, roughness: 0.9 }),
      );
      weaponGrip.position.y = -0.65;
      weapon.add(weaponGuard, weaponGrip);
      weapon.position.set(0, 0.65, 0);
      weapon.rotation.z = -0.35;
      const warning = this.mesh(
        new THREE.RingGeometry(0.72, 0.96, 28),
        new THREE.MeshBasicMaterial({ color: 0xff7548, transparent: true, opacity: 0, side: THREE.DoubleSide, depthWrite: false }),
        false,
      );
      warning.rotation.x = -Math.PI / 2;
      warning.position.y = 0.04;
      warning.visible = false;
      const warningBeacon = this.mesh(
        new THREE.ConeGeometry(0.17, 0.36, 3),
        new THREE.MeshBasicMaterial({ color: 0xffa45c, transparent: true, opacity: 0.9, depthTest: false }),
        false,
      );
      warningBeacon.position.set(0, 2.78, 0.08);
      warningBeacon.rotation.z = Math.PI;
      warningBeacon.visible = false;
      warningBeacon.renderOrder = 5;
      this.weaponHand.position.set(0, -1, 0);
      this.weaponHand.add(weapon);
      this.arms[1]!.add(this.weaponHand);
      group.add(shell, hood, face, core, light);
      this.root.add(group, warning, warningBeacon);
      this.shell = shell;
      this.core = core as THREE.Mesh<THREE.BufferGeometry, THREE.MeshStandardMaterial>;
      this.warning = warning as THREE.Mesh<THREE.BufferGeometry, THREE.MeshBasicMaterial>;
      this.warningBeacon = warningBeacon;
      this.light = light;

  }

  reset(): void {
    this.root.visible = true;
    this.root.rotation.set(0, 0, 0);
    this.body.position.set(0, 0, 0);
    this.body.rotation.set(0, 0, 0);
    this.body.scale.setScalar(1);
    this.warning.visible = this.warningBeacon.visible = false;
  }

  update(state: Readonly<EnemyCombatState>, elapsed: number, delta: number, facing: number, moving: boolean): void {
    const pose = sampleEnemyPose(state);
    const drift = Math.sin(elapsed * 2.4 + this.index * 1.7);
    this.root.rotation.y = turnTowards(this.root.rotation.y, facing, delta);
    // The body floats inside the rig. Its gameplay origin remains on the ground.
    this.body.position.y = drift * 0.055 + (1 - pose.dissolve) * 0.2;
    this.body.position.z = -pose.recoil * 0.13;
    this.body.rotation.set(pose.lean - pose.recoil * 0.23, pose.twist, moving ? drift * 0.035 : drift * 0.015);
    this.body.scale.set(1 + (1 - pose.dissolve) * 0.25, Math.max(0.03, pose.dissolve), 1);
    this.shell.scale.y = 1 + Math.sin(elapsed * 2 + this.index) * 0.025;
    this.arms[1]!.rotation.set(pose.arm - pose.recoil * 0.3, pose.twist * 0.6, 0.22 + pose.warning * 0.25);
    this.weaponHand.rotation.x = pose.wrist;
    this.arms[0]!.rotation.set(-pose.arm * 0.3, 0, -0.22 - pose.warning * 0.4);
    this.core.material.emissiveIntensity = state.hitFlash > 0 ? 9 : 2.8 + pose.warning * 5;
    this.core.scale.setScalar(1 + pose.recoil * 0.55 + pose.warning * 0.2);
    this.warning.visible = this.warningBeacon.visible = state.alive && state.attackWindup > 0;
    this.warning.scale.setScalar(0.78 + pose.warning * 0.58);
    this.warning.material.opacity = 0.18 + pose.warning * 0.68;
    this.warningBeacon.position.y = 2.78 + Math.sin(elapsed * 6) * 0.025;
    this.warningBeacon.scale.setScalar(0.85 + pose.warning * 0.55);
    this.light.intensity = 5 * pose.dissolve;
    this.root.visible = state.alive || state.deathTime > 0;
  }

  private mesh(geometry: THREE.BufferGeometry, material: THREE.Material, shadows = true): THREE.Mesh {
    const mesh = new THREE.Mesh(geometry, material);
    mesh.castShadow = mesh.receiveShadow = shadows;
    return mesh;
  }
}
