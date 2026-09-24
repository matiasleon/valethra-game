import * as THREE from "three";
import type { EnemyCombatState } from "./enemy-combat";
import { sampleEnemyPose } from "./enemy-pose";
import { sampleTrollDeath } from "./troll-death";
import { damping, turnTowards } from "./locomotion";

/** Presentation-only troll rig. The root is the stable gameplay origin. */
export class TrollActor {
  readonly root = new THREE.Group();
  private readonly body = new THREE.Group();
  private readonly torso = new THREE.Group();
  private readonly head = new THREE.Group();
  private readonly arms = [new THREE.Group(), new THREE.Group()];
  private readonly elbows = [new THREE.Group(), new THREE.Group()];
  private readonly thighs = [new THREE.Group(), new THREE.Group()];
  private readonly knees = [new THREE.Group(), new THREE.Group()];
  private readonly hand = new THREE.Group();
  private readonly club = new THREE.Group();
  private readonly warning: THREE.Mesh<THREE.RingGeometry, THREE.MeshBasicMaterial>;
  private readonly beacon: THREE.Mesh;
  private readonly bounds = new THREE.Box3();
  private readonly dropStart = new THREE.Vector3();
  private readonly dropRotation = new THREE.Quaternion();
  private readonly groundRotation = new THREE.Quaternion().setFromEuler(new THREE.Euler(0.08, 0.2, Math.PI / 2));
  private readonly parts: THREE.Object3D[];
  private deathRotations: THREE.Euler[] = [];
  private deathBodyPosition = new THREE.Vector3();
  private dying = false;
  private dropped = false;
  private stride = 0;
  private gait = 0;
  private readonly skin: THREE.MeshStandardMaterial;

  constructor(private readonly index: number) {
    this.skin = new THREE.MeshStandardMaterial({ color: [0x74806a, 0x7b7760, 0x68796f, 0x85816e][index % 4], roughness: 0.9 });
    const darkSkin = new THREE.MeshStandardMaterial({ color: 0x4c5947, roughness: 0.96 });
    const leather = new THREE.MeshStandardMaterial({ color: 0x403027, roughness: 0.94 });
    const fur = new THREE.MeshStandardMaterial({ color: 0x302e28, roughness: 1, flatShading: true });
    const metal = new THREE.MeshStandardMaterial({ color: 0x77756c, metalness: 0.55, roughness: 0.64 });
    const ivory = new THREE.MeshStandardMaterial({ color: 0xcbb998, roughness: 0.72 });
    const eye = new THREE.MeshStandardMaterial({ color: 0xc9ab69, roughness: 0.5 });
    const pupil = new THREE.MeshStandardMaterial({ color: 0x171915 });
    const wood = new THREE.MeshStandardMaterial({ color: 0x5b402a, roughness: 1 });
    this.root.add(this.body);
    this.body.position.y = 1.15;
    this.body.add(this.torso);
    this.torso.position.y = 0.32;
    this.ellipsoid(this.body, [0.52, 0.32, 0.36], [0, 0, 0], leather);
    this.ellipsoid(this.torso, [0.65, 0.65, 0.38], [0, 0.3, 0], this.skin);
    this.ellipsoid(this.torso, [0.76, 0.36, 0.42], [0, 0.62, -0.08], darkSkin);
    this.ellipsoid(this.torso, [0.5, 0.37, 0.36], [0, 0.02, 0.09], this.skin);
    const belt = this.mesh(this.body, new THREE.CylinderGeometry(0.53, 0.53, 0.16, 12), leather, [0, 0.07, 0]); belt.scale.z = 0.7;
    this.mesh(this.body, new THREE.BoxGeometry(0.16, 0.14, 0.045), metal, [0, 0.07, 0.39]);
    for (let i = 0; i < 7; i++) {
      const angle = i / 7 * Math.PI * 2;
      const flap = this.mesh(this.body, new THREE.ConeGeometry(0.16, 0.42, 4), fur, [Math.sin(angle) * 0.45, -0.22, Math.cos(angle) * 0.31]);
      flap.rotation.z = Math.PI;
    }
    this.torso.add(this.head); this.head.position.set(0, 0.93, 0.15);
    this.ellipsoid(this.head, [0.32, 0.35, 0.28], [0, 0.1, 0], this.skin);
    this.ellipsoid(this.head, [0.32, 0.18, 0.26], [0, -0.09, 0.12], darkSkin);
    this.ellipsoid(this.head, [0.14, 0.14, 0.2], [0, 0.09, 0.26], this.skin);
    for (const side of [-1, 1]) {
      this.ellipsoid(this.head, [0.13, 0.055, 0.1], [side * 0.16, 0.22, 0.23], darkSkin);
      this.ellipsoid(this.head, [0.058, 0.036, 0.035], [side * 0.16, 0.155, 0.29], eye);
      this.ellipsoid(this.head, [0.021, 0.027, 0.015], [side * 0.16, 0.155, 0.32], pupil);
      const ear = this.mesh(this.head, new THREE.ConeGeometry(0.13, 0.32, 5), this.skin, [side * 0.34, 0.14, -0.02]); ear.rotation.z = side * -1.05;
      const tusk = this.mesh(this.head, new THREE.ConeGeometry(0.055, 0.19, 8), ivory, [side * 0.22, -0.015, 0.32]); tusk.rotation.x = 0.18;
    }
    for (let i = 0; i < 2; i++) {
      const side = i === 0 ? -1 : 1;
      const arm = this.arms[i]!; this.torso.add(arm); arm.position.set(side * 0.7, 0.65, 0);
      this.ellipsoid(arm, [0.31, 0.28, 0.3], [0, -0.03, 0], this.skin);
      this.mesh(arm, new THREE.CapsuleGeometry(0.2, 0.37, 6, 10), this.skin, [0, -0.34, 0]);
      const elbow = this.elbows[i]!; arm.add(elbow); elbow.position.y = -0.63;
      this.mesh(elbow, new THREE.CapsuleGeometry(0.17, 0.3, 6, 10), this.skin, [0, -0.2, 0]);
      this.mesh(elbow, new THREE.CylinderGeometry(0.19, 0.17, 0.22, 10), leather, [0, -0.31, 0]);
      this.ellipsoid(elbow, [0.19, 0.2, 0.15], [0, -0.52, 0.025], this.skin);
      const thigh = this.thighs[i]!; this.body.add(thigh); thigh.position.set(side * 0.28, -0.12, 0);
      this.mesh(thigh, new THREE.CapsuleGeometry(0.2, 0.25, 6, 10), this.skin, [0, -0.2, 0]);
      const knee = this.knees[i]!; thigh.add(knee); knee.position.y = -0.47;
      this.mesh(knee, new THREE.CapsuleGeometry(0.15, 0.22, 6, 10), darkSkin, [0, -0.17, 0]);
      this.mesh(knee, new THREE.CylinderGeometry(0.17, 0.16, 0.2, 10), leather, [0, -0.26, 0]);
      this.ellipsoid(knee, [0.2, 0.13, 0.3], [0, -0.43, 0.12], darkSkin);
      for (const toe of [-1, 0, 1]) this.ellipsoid(knee, [0.047, 0.045, 0.08], [toe * 0.095, -0.43, 0.36], ivory);
    }
    // Asymmetric worn shoulder guard, without assigning a new faction or backstory.
    this.ellipsoid(this.arms[0]!, [0.34, 0.15, 0.34], [0, 0.08, 0], metal);
    this.elbows[1]!.add(this.hand); this.hand.position.set(0, -0.52, 0.02);
    this.hand.add(this.club);
    this.mesh(this.club, new THREE.CylinderGeometry(0.07, 0.055, 1.05, 9), wood, [0, 0.36, 0]);
    this.ellipsoid(this.club, [0.19, 0.34, 0.18], [0, 0.88, 0], wood);
    for (const height of [0.7, 1.02]) this.mesh(this.club, new THREE.CylinderGeometry(0.19, 0.19, 0.07, 10), metal, [0, height, 0]);
    this.warning = new THREE.Mesh(new THREE.RingGeometry(0.75, 1, 32), new THREE.MeshBasicMaterial({ color: 0xe8a25d, transparent: true, opacity: 0, side: THREE.DoubleSide, depthWrite: false }));
    this.warning.rotation.x = -Math.PI / 2; this.warning.position.y = 0.03;
    this.beacon = new THREE.Mesh(new THREE.ConeGeometry(0.13, 0.26, 3), new THREE.MeshBasicMaterial({ color: 0xe8b578, depthTest: false }));
    this.beacon.rotation.z = Math.PI; this.beacon.position.y = 3;
    this.root.add(this.warning, this.beacon);
    this.parts = [this.body, this.torso, this.head, ...this.arms, ...this.elbows, ...this.thighs, ...this.knees, this.hand];
    this.reset();
  }

  reset(): void {
    this.dying = this.dropped = false;
    this.stride = this.gait = 0;
    this.root.visible = true;
    this.root.rotation.set(0, 0, 0);
    this.body.position.set(0, 1.15, 0);
    for (const part of this.parts) part.rotation.set(0, 0, 0);
    this.hand.add(this.club);
    this.club.position.set(0, 0, 0); this.club.rotation.set(0, 0, 0); this.club.scale.setScalar(1);
    this.skin.emissiveIntensity = 0;
    this.warning.visible = this.beacon.visible = false;
  }

  update(state: Readonly<EnemyCombatState>, elapsed: number, delta: number, facing: number, moving: boolean): void {
    if (!state.alive) { this.updateDeath(state.deathTime); return; }
    const pose = sampleEnemyPose(state);
    this.root.rotation.y = turnTowards(this.root.rotation.y, facing, delta);
    this.gait += ((moving ? 1 : 0) - this.gait) * damping(9, delta);
    this.stride += delta * 4.4 * (moving ? 1 : 0);
    const step = Math.sin(this.stride) * this.gait;
    this.body.position.set(0, 1.15 + Math.abs(step) * 0.025, -pose.recoil * 0.08);
    this.body.rotation.set(0, 0, step * 0.035);
    this.torso.rotation.set(0.12 + pose.lean * 0.6 - pose.recoil * 0.18, pose.twist, 0);
    this.head.rotation.x = -0.07 + Math.sin(elapsed * 1.8 + this.index) * 0.015;
    this.arms[1]!.rotation.set(pose.arm - pose.recoil * 0.2, pose.twist * 0.6, 0.15);
    this.arms[0]!.rotation.set(-pose.arm * 0.25 + step * 0.15, 0, -0.15);
    this.elbows[1]!.rotation.x = -0.18 - pose.warning * 0.3;
    this.elbows[0]!.rotation.x = -0.15;
    this.hand.rotation.x = pose.wrist;
    for (let i = 0; i < 2; i++) {
      const swing = i === 0 ? step : -step;
      this.thighs[i]!.rotation.x = swing * 0.28;
      this.knees[i]!.rotation.x = Math.max(0, -swing) * 0.35;
    }
    this.skin.emissive.set(0x8b3d24); this.skin.emissiveIntensity = pose.recoil * 0.3;
    this.warning.visible = this.beacon.visible = state.attackWindup > 0;
    this.warning.material.opacity = 0.18 + pose.warning * 0.5;
    this.warning.scale.setScalar(0.85 + pose.warning * 0.25);
  }

  private updateDeath(remaining: number): void {
    if (!this.dying) {
      this.dying = true;
      this.deathRotations = this.parts.map(part => part.rotation.clone());
      this.deathBodyPosition.copy(this.body.position);
    }
    const pose = sampleTrollDeath(remaining);
    this.warning.visible = this.beacon.visible = false;
    this.skin.emissiveIntensity = 0;
    // Blend from the actual struck pose rather than snapping to a canned standing pose.
    const rotations = [
      [pose.pitch, 0, pose.roll], [0.16, -0.1 * pose.fall, 0], [0.15 + pose.settle * 0.12, 0.1, 0],
      [pose.arm * (1 - pose.settle), 0.1, -0.25 + pose.fall * 1.6], [-0.3 + pose.fall * 0.7, 0.2, 0.3 - pose.fall * 0.55],
      [-0.5, 0, 0], [-0.35, 0, 0],
      [-0.72 * pose.knees, 0, -0.12 * pose.fall], [-0.58 * pose.knees, 0, 0.18 * pose.fall],
      [1.3 * pose.knees, 0, 0], [1.1 * pose.knees, 0, 0], [0.3, 0, 0],
    ];
    this.parts.forEach((part, i) => {
      const from = this.deathRotations[i]!; const to = rotations[i]!;
      part.rotation.set(THREE.MathUtils.lerp(from.x, to[0]!, pose.blend), THREE.MathUtils.lerp(from.y, to[1]!, pose.blend), THREE.MathUtils.lerp(from.z, to[2]!, pose.blend));
    });
    this.body.position.set(0.18 * pose.fall, THREE.MathUtils.lerp(this.deathBodyPosition.y, pose.height, pose.blend), -0.16 * pose.knees);
    // Ground contact is measured from the posed mesh, so limbs never sink or squash.
    this.root.updateMatrixWorld(true);
    this.bounds.setFromObject(this.body);
    this.body.position.y += Math.max(0, this.root.position.y + 0.025 - this.bounds.min.y);
    this.body.position.y += pose.impact * 0.025;
    if (pose.drop > 0 && !this.dropped) {
      this.root.updateMatrixWorld(true);
      this.root.attach(this.club);
      this.dropStart.copy(this.club.position); this.dropRotation.copy(this.club.quaternion);
      this.dropped = true;
    }
    if (this.dropped) {
      this.club.position.set(THREE.MathUtils.lerp(this.dropStart.x, 1.05, pose.drop), THREE.MathUtils.lerp(this.dropStart.y, 0.2, pose.drop) + Math.sin(pose.drop * Math.PI) * 0.18, THREE.MathUtils.lerp(this.dropStart.z, 0.35, pose.drop));
      this.club.quaternion.slerpQuaternions(this.dropRotation, this.groundRotation, pose.drop);
    }
  }

  private mesh(parent: THREE.Object3D, geometry: THREE.BufferGeometry, material: THREE.Material, position: [number, number, number]): THREE.Mesh {
    const mesh = new THREE.Mesh(geometry, material); mesh.position.set(...position);
    mesh.castShadow = mesh.receiveShadow = true; parent.add(mesh); return mesh;
  }

  private ellipsoid(parent: THREE.Object3D, scale: [number, number, number], position: [number, number, number], material: THREE.Material): THREE.Mesh {
    const mesh = this.mesh(parent, new THREE.SphereGeometry(1, 12, 9), material, position); mesh.scale.set(...scale); return mesh;
  }
}
