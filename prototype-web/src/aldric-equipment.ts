import * as THREE from "three";

export class AldricEquipment {
  readonly weapon = new THREE.Group();
  readonly shield = new THREE.Group();
  private face!: THREE.Mesh<THREE.BufferGeometry, THREE.MeshStandardMaterial>;
  private boss!: THREE.Mesh;
  constructor() { this.buildWeapon(); this.buildShield(); }
  updateImpact(flash: number): void {
    this.face.material.emissiveIntensity = 0.7 + flash * 18;
    this.boss.scale.setScalar(1 + flash * 1.8);
  }
  private buildWeapon(): void {
    const grip = this.mesh(
      new THREE.CylinderGeometry(0.045, 0.055, 0.55, 8),
      new THREE.MeshStandardMaterial({ color: 0x3c2a20, roughness: 0.9 }),
      false,
    );
    const bladeShape = new THREE.Shape();
    bladeShape.moveTo(-0.055, -0.43);
    bladeShape.lineTo(0.055, -0.43);
    bladeShape.lineTo(0.08, 0.3);
    bladeShape.lineTo(0, 0.52);
    bladeShape.lineTo(-0.08, 0.3);
    const blade = this.mesh(
      new THREE.ExtrudeGeometry(bladeShape, { depth: 0.025, bevelEnabled: true, bevelSize: 0.012, bevelThickness: 0.012, bevelSegments: 1 }),
      new THREE.MeshStandardMaterial({ color: 0xa9afb1, roughness: 0.34, metalness: 0.86 }),
      false,
    );
    grip.position.y = -0.15;
    blade.position.y = 0.55;
    this.weapon.add(grip, blade);

  }

  private buildShield(): void {
    const faceMaterial = new THREE.MeshStandardMaterial({
      color: 0x35443f,
      emissive: 0x183c35,
      emissiveIntensity: 0.7,
      roughness: 0.58,
      metalness: 0.52,
    });
    const face = this.mesh(new THREE.CylinderGeometry(0.37, 0.37, 0.09, 8), faceMaterial, false);
    face.rotation.x = Math.PI / 2;
    const rim = this.mesh(
      new THREE.TorusGeometry(0.37, 0.038, 7, 8),
      new THREE.MeshStandardMaterial({ color: 0x8e7654, roughness: 0.38, metalness: 0.78 }),
      false,
    );
    const boss = this.mesh(
      new THREE.OctahedronGeometry(0.1, 0),
      new THREE.MeshStandardMaterial({ color: 0x91c6b5, emissive: 0x4bc49d, emissiveIntensity: 1.5, roughness: 0.3, metalness: 0.7 }),
      false,
    );
    rim.rotation.x = Math.PI / 2;
    boss.position.z = 0.09;
    this.shield.add(face, rim, boss);
    this.face = face as THREE.Mesh<THREE.BufferGeometry, THREE.MeshStandardMaterial>;
    this.boss = boss;

  }

  private mesh(geometry: THREE.BufferGeometry, material: THREE.Material, shadows = true): THREE.Mesh {
    const mesh = new THREE.Mesh(geometry, material);
    mesh.castShadow = mesh.receiveShadow = shadows;
    return mesh;
  }
}
