import * as THREE from "three";

export interface Ward {
  id: string;
  title: string;
  message: string;
  position: THREE.Vector3;
  mesh: THREE.Group;
  light: THREE.PointLight;
  active: boolean;
  core: THREE.Mesh<THREE.OctahedronGeometry, THREE.MeshStandardMaterial>;
}

interface Obstacle {
  gate?: boolean;
  minX: number;
  maxX: number;
  minZ: number;
  maxZ: number;
}

/** Owns scenery and its collision/interaction handles, never player or combat state. */
export class SanctuaryWorld {
  readonly obstacles: Obstacle[] = [];
  readonly cameraSolids: THREE.Object3D[] = [];
  readonly wards: Ward[] = [];
  readonly gate = new THREE.Group();
  readonly family = new THREE.Group();
  doors!: THREE.Group;
  lantern!: THREE.PointLight;

  constructor(private readonly scene: THREE.Scene, textures: { stone: THREE.Texture; earth: THREE.Texture; wood: THREE.Texture }) {
    this.build(textures);
  }

  collides(x: number, z: number, radius: number, gateOpen: boolean): boolean {
    if (x < -26 || x > 26 || z > 27 || z < -47) return true;
    return this.obstacles.some(obstacle => !(obstacle.gate && gateOpen)
      && x + radius > obstacle.minX && x - radius < obstacle.maxX
      && z + radius > obstacle.minZ && z - radius < obstacle.maxZ);
  }

  reset(): void {
    this.doors.position.y = 0;
    this.family.visible = false;
    for (const ward of this.wards) {
      ward.active = false;
      ward.light.intensity = 1.8;
      ward.core.scale.setScalar(0.58);
      ward.core.material.emissiveIntensity = 0.18;
    }
  }

  private build(textures: { stone: THREE.Texture; earth: THREE.Texture; wood: THREE.Texture }): void {
    const hemi = new THREE.HemisphereLight(0x7899a0, 0x30302b, 2.35);
    const moon = new THREE.DirectionalLight(0xb8d5dc, 3.1);
    moon.position.set(-18, 28, 16);
    moon.castShadow = true;
    moon.shadow.mapSize.set(2048, 2048);
    moon.shadow.camera.left = -38;
    moon.shadow.camera.right = 38;
    moon.shadow.camera.top = 38;
    moon.shadow.camera.bottom = -38;
    this.scene.add(hemi, moon);
    for (const [x, z] of [[-14, 18], [14, 1], [-13, -15], [11, -29]] as const) {
      const ember = new THREE.PointLight(0xe89a54, 9, 11, 2);
      ember.position.set(x, 2.1, z);
      const brazier = this.mesh(
        new THREE.CylinderGeometry(0.24, 0.36, 0.55, 7),
        new THREE.MeshStandardMaterial({ color: 0x392b23, roughness: 0.72, metalness: 0.45 }),
      );
      brazier.position.set(x, 0.3, z);
      const flame = this.mesh(
        new THREE.OctahedronGeometry(0.2, 0),
        new THREE.MeshStandardMaterial({ color: 0xffb56d, emissive: 0xe15c27, emissiveIntensity: 5 }),
        false,
      );
      flame.position.set(x, 0.88, z);
      this.scene.add(ember, brazier, flame);
    }

    const earthMaterial = new THREE.MeshStandardMaterial({ map: textures.earth, roughness: 0.98, metalness: 0 });
    const stoneMaterial = new THREE.MeshStandardMaterial({ map: textures.stone, roughness: 0.91, metalness: 0.02 });
    const woodMaterial = new THREE.MeshStandardMaterial({ map: textures.wood, roughness: 0.9, metalness: 0.01 });
    const ground = this.mesh(new THREE.PlaneGeometry(70, 105), earthMaterial);
    ground.rotation.x = -Math.PI / 2;
    ground.position.z = -12;
    this.scene.add(ground);
    this.addAtmosphere(stoneMaterial);

    this.addWall(-18, 1.7, 5, 2.4, 3.4, 34, stoneMaterial);
    this.addWall(18, 1.7, 5, 2.4, 3.4, 34, stoneMaterial);
    this.addWall(-17, 1.7, -25, 11, 3.4, 2.2, stoneMaterial);
    this.addWall(17, 1.7, -25, 11, 3.4, 2.2, stoneMaterial);
    this.addWall(-10.5, 1.45, -8, 5.5, 2.9, 1.5, stoneMaterial);
    this.addWall(10.5, 1.45, 3, 5.5, 2.9, 1.5, stoneMaterial);
    this.addWall(0, 1.15, -15, 5, 2.3, 1.3, stoneMaterial);

    for (let index = 0; index < 38; index += 1) {
      const side = index % 2 === 0 ? -1 : 1;
      const x = side * (21 + ((index * 7) % 8) * 0.8);
      const z = 28 - index * 2.25;
      this.scene.add(this.createTree(x, z, 4.5 + (index % 5) * 0.65));
    }

    const cart = this.createBrokenCart(-7.2, 0.2, 12, woodMaterial);
    this.scene.add(cart);
    this.cameraSolids.push(cart);
    this.addNarrativeProps(stoneMaterial, woodMaterial);
    this.createWards(stoneMaterial);
    this.createGate(stoneMaterial, woodMaterial);
    this.createFamily();

    const lantern = new THREE.PointLight(0xe5aa68, 6, 8, 2);
    lantern.position.set(0, 1.25, 0);
    this.scene.add(lantern);
    this.lantern = lantern;
  }

  private addWall(x: number, y: number, z: number, width: number, height: number, depth: number, material: THREE.Material): void {
    const wall = this.mesh(new THREE.BoxGeometry(width, height, depth), material);
    wall.position.set(x, y, z);
    this.scene.add(wall);
    this.cameraSolids.push(wall);
    this.obstacles.push({ minX: x - width / 2, maxX: x + width / 2, minZ: z - depth / 2, maxZ: z + depth / 2 });
  }

  private createTree(x: number, z: number, height: number): THREE.Group {
    const tree = new THREE.Group();
    const trunk = this.mesh(
      new THREE.CylinderGeometry(0.28, 0.48, height * 0.58, 7),
      new THREE.MeshStandardMaterial({ color: 0x2f241d, roughness: 1, flatShading: true }),
    );
    trunk.position.y = height * 0.29;
    const crown = this.mesh(
      new THREE.ConeGeometry(1.35, height * 0.8, 7),
      new THREE.MeshStandardMaterial({ color: 0x1b3028, roughness: 1, flatShading: true }),
    );
    crown.position.y = height * 0.78;
    tree.add(trunk, crown);
    tree.position.set(x, 0, z);
    return tree;
  }

  private createBrokenCart(x: number, y: number, z: number, material: THREE.Material): THREE.Group {
    const cart = new THREE.Group();
    const bed = this.mesh(new THREE.BoxGeometry(3.2, 0.28, 1.7), material);
    bed.position.y = 0.8;
    bed.rotation.z = 0.13;
    const axle = this.mesh(new THREE.CylinderGeometry(0.12, 0.12, 3.8, 8), material);
    axle.rotation.z = Math.PI / 2;
    axle.position.y = 0.45;
    const wheelMaterial = new THREE.MeshStandardMaterial({ color: 0x241c17, roughness: 1 });
    for (const side of [-1, 1]) {
      const wheel = this.mesh(new THREE.TorusGeometry(0.65, 0.12, 7, 12), wheelMaterial);
      wheel.position.set(side * 1.55, 0.45, 0);
      wheel.rotation.y = Math.PI / 2;
      cart.add(wheel);
    }
    cart.add(bed, axle);
    cart.position.set(x, y, z);
    cart.rotation.y = -0.32;
    this.obstacles.push({ minX: x - 2, maxX: x + 2, minZ: z - 1.4, maxZ: z + 1.4 });
    return cart;
  }

  private addNarrativeProps(stone: THREE.Material, wood: THREE.Material): void {
    const bellFrame = new THREE.Group();
    const leftPost = this.mesh(new THREE.BoxGeometry(0.32, 3.2, 0.32), wood);
    const rightPost = leftPost.clone();
    leftPost.position.set(-1.25, 1.6, 0);
    rightPost.position.set(1.25, 1.6, 0);
    const beam = this.mesh(new THREE.BoxGeometry(2.9, 0.32, 0.32), wood);
    beam.position.y = 3.1;
    const bell = this.mesh(
      new THREE.ConeGeometry(0.48, 0.8, 9),
      new THREE.MeshStandardMaterial({ color: 0x67543b, roughness: 0.42, metalness: 0.72 }),
    );
    bell.position.y = 2.45;
    bell.rotation.x = Math.PI;
    bellFrame.add(leftPost, rightPost, beam, bell);
    bellFrame.position.set(-12, 0, 8);
    this.scene.add(bellFrame);

    for (let index = 0; index < 5; index += 1) {
      const shield = this.mesh(new THREE.CylinderGeometry(0.62, 0.62, 0.13, 8), stone);
      shield.rotation.set(Math.PI / 2, 0, index % 2 === 0 ? 0.2 : -0.14);
      shield.position.set(6.5 + index * 0.72, 0.34, 8.2 + (index % 2) * 0.35);
      this.scene.add(shield);
    }

    const crib = this.mesh(new THREE.BoxGeometry(1.25, 0.55, 0.72), wood);
    crib.position.set(-2.4, 0.32, -27.8);
    crib.rotation.y = 0.25;
    this.scene.add(crib);
  }

  private addAtmosphere(stone: THREE.Material): void {
    const pathMaterial = (stone as THREE.MeshStandardMaterial).clone();
    pathMaterial.color.set(0x69716a);
    pathMaterial.roughness = 0.96;
    for (let index = 0; index < 25; index += 1) {
      const slab = this.mesh(new THREE.BoxGeometry(3.1 + (index % 3) * 0.32, 0.08, 2.05), pathMaterial);
      slab.position.set(Math.sin(index * 1.8) * 1.25, 0.015, 24 - index * 2.35);
      slab.rotation.y = Math.sin(index * 2.7) * 0.09;
      this.scene.add(slab);
    }

    const rubbleMaterial = new THREE.MeshStandardMaterial({ color: 0x3a403c, roughness: 1, flatShading: true });
    for (let index = 0; index < 46; index += 1) {
      const side = index % 2 === 0 ? -1 : 1;
      const size = 0.13 + (index % 5) * 0.055;
      const rubble = this.mesh(new THREE.DodecahedronGeometry(size, 0), rubbleMaterial);
      rubble.position.set(side * (5.5 + ((index * 7) % 13)), size * 0.65, 27 - index * 1.55);
      rubble.rotation.set(index * 0.3, index * 0.7, 0);
      this.scene.add(rubble);
    }

    const moon = this.mesh(
      new THREE.SphereGeometry(4.2, 20, 12),
      new THREE.MeshBasicMaterial({ color: 0xb8d2cf, fog: false }),
      false,
    );
    moon.position.set(28, 34, -78);
    this.scene.add(moon);

    const motePositions = new Float32Array(270 * 3);
    for (let index = 0; index < 270; index += 1) {
      motePositions[index * 3] = ((index * 47) % 520) / 10 - 26;
      motePositions[index * 3 + 1] = ((index * 29) % 70) / 10 + 0.2;
      motePositions[index * 3 + 2] = 28 - ((index * 73) % 750) / 10;
    }
    const moteGeometry = new THREE.BufferGeometry();
    moteGeometry.setAttribute("position", new THREE.BufferAttribute(motePositions, 3));
    const motes = new THREE.Points(
      moteGeometry,
      new THREE.PointsMaterial({ color: 0xb6c9bc, size: 0.035, transparent: true, opacity: 0.42, depthWrite: false }),
    );
    this.scene.add(motes);
  }

  private createWards(material: THREE.Material): void {
    const wardData = [
      { id: "bell", title: "El sello de la campana", message: "La cuerda fue cortada desde el interior. Nadie quiso que el bosque oyera la alarma.", x: -12, z: 7 },
      { id: "watch", title: "El sello de la guardia", message: "Los escudos caídos miran hacia el santuario. La guardia no defendía el camino: huía de él.", x: 11, z: -7 },
      { id: "hearth", title: "El sello del hogar", message: "Junto al brasero hay una cuna vacía. Alguien intentó cruzar antes de que cerraran el portón.", x: -6, z: -24 },
    ];
    for (const data of wardData) {
      const group = new THREE.Group();
      const base = this.mesh(new THREE.CylinderGeometry(0.85, 1.1, 1.25, 7), material);
      base.position.y = 0.62;
      const coreMaterial = new THREE.MeshStandardMaterial({
        color: 0x8aa79c,
        emissive: 0x6bd8b0,
        emissiveIntensity: 0.18,
        roughness: 0.48,
        metalness: 0.16,
        flatShading: true,
      });
      const core = this.mesh(new THREE.OctahedronGeometry(0.42, 0), coreMaterial, false);
      core.position.y = 1.65;
      core.scale.setScalar(0.58);
      const light = new THREE.PointLight(0x67d9ae, 0, 9, 2);
      light.position.y = 1.75;
      group.add(base, core, light);
      group.position.set(data.x, 0, data.z);
      this.scene.add(group);
      this.wards.push({
        id: data.id,
        title: data.title,
        message: data.message,
        position: new THREE.Vector3(data.x, 0, data.z),
        mesh: group,
        light,
        active: false,
        core: core as THREE.Mesh<THREE.OctahedronGeometry, THREE.MeshStandardMaterial>,
      });
    }
  }

  private createGate(stone: THREE.Material, wood: THREE.Material): void {
    const leftTower = this.mesh(new THREE.BoxGeometry(5.5, 7.2, 4), stone);
    const rightTower = leftTower.clone();
    leftTower.position.set(-7.2, 3.6, -34);
    rightTower.position.set(7.2, 3.6, -34);
    const arch = this.mesh(new THREE.BoxGeometry(9.2, 2, 3), stone);
    arch.position.set(0, 6.2, -34);
    const doors = new THREE.Group();
    for (const side of [-1, 1]) {
      const door = this.mesh(new THREE.BoxGeometry(3.8, 5.2, 0.42), wood);
      door.position.set(side * 2, 2.6, 0);
      doors.add(door);
    }
    doors.position.set(0, 0, -33.7);
    this.gate.add(leftTower, rightTower, arch, doors);
    this.doors = doors;
    this.scene.add(this.gate);
    this.cameraSolids.push(this.gate);
    this.obstacles.push({ gate: true, minX: -4.2, maxX: 4.2, minZ: -34.2, maxZ: -33.2 });
  }

  private createFamily(): void {
    const colors = [0x5b493e, 0x3f5057, 0x705342];
    const scales = [1, 0.92, 0.67];
    const xs = [-0.05, 0.95, 0];
    for (let index = 0; index < 3; index += 1) {
      const person = new THREE.Group();
      const body = this.mesh(
        new THREE.BoxGeometry(0.72, 1.05, 0.46),
        new THREE.MeshStandardMaterial({ color: colors[index], roughness: 0.95, flatShading: true }),
      );
      body.position.y = 1.35;
      const head = this.mesh(
        new THREE.DodecahedronGeometry(0.3, 0),
        new THREE.MeshStandardMaterial({ color: 0x9a694c, roughness: 0.9, flatShading: true }),
      );
      head.position.y = 2.15;
      person.add(body, head);
      person.scale.setScalar(scales[index]!);
      person.position.set(xs[index]!, 0, index === 2 ? 0.45 : 0);
      this.family.add(person);
    }
    const refugeLight = new THREE.PointLight(0xefad6b, 22, 14, 2);
    refugeLight.position.set(0, 3.4, 1.5);
    this.family.add(refugeLight);
    for (let index = 0; index < 5; index += 1) {
      const windowLight = new THREE.PointLight(0xeaa467, 5, 7, 2);
      windowLight.position.set(-10 + index * 5, 2.2 + (index % 2), -6 - (index % 3) * 2);
      const home = this.mesh(
        new THREE.BoxGeometry(2.3, 2.4, 2.2),
        new THREE.MeshStandardMaterial({ color: 0x262927, roughness: 1 }),
      );
      home.position.set(windowLight.position.x, 1.2, windowLight.position.z);
      const roof = this.mesh(
        new THREE.ConeGeometry(2, 1.7, 4),
        new THREE.MeshStandardMaterial({ color: 0x1a1b1b, roughness: 1, flatShading: true }),
      );
      roof.position.set(windowLight.position.x, 3.15, windowLight.position.z);
      roof.rotation.y = Math.PI / 4;
      this.family.add(home, roof, windowLight);
    }
    this.family.position.set(0, 0, -42);
    this.family.visible = false;
    this.scene.add(this.family);
  }

  private mesh(geometry: THREE.BufferGeometry, material: THREE.Material, shadows = true): THREE.Mesh {
    const mesh = new THREE.Mesh(geometry, material);
    mesh.castShadow = shadows;
    mesh.receiveShadow = shadows;
    return mesh;
  }
}
