import * as THREE from "three";
import { AldricEquipment } from "./aldric-equipment";
import { AldricActor } from "./aldric-actor";
import { ThirdPersonCamera } from "./third-person-camera";
import { TrollActor } from "./troll-actor";
import { ENEMY_SPAWNS, ENEMY_DAMAGE, ENEMY_DEATH_DURATION, advanceEnemy, createEnemyState, type EnemyCombatState } from "./enemy-combat";
import { damping } from "./locomotion";

import {
  ATTACK_STAMINA_COST,
  DODGE_STAMINA_COST,
  MAX_HEALTH,
  MAX_STAMINA,
  canOpenWardGate,
  canSpendStamina,
  objectiveForWards,
  resolveIncomingHit,
  spendStamina,
} from "./combat-rules";

import type { LevelCallbacks, LevelEnding } from "./level-contracts";
import { SanctuaryWorld, type Ward } from "./sanctuary-world";
import { disposeScene } from "./scene-resources";

interface Enemy extends EnemyCombatState {
  actor: TrollActor;
  index: number;
  defenseDemonstrated: boolean;
}

const PLAYER_RADIUS = 0.38;
const WALK_SPEED = 4.1;
const SPRINT_SPEED = 6.7;
const DODGE_SPEED = 10.8;
const DODGE_DURATION = 0.42;

export class ThirdPersonLevel {
  private readonly renderer: THREE.WebGLRenderer;
  private readonly scene = new THREE.Scene();
  private readonly playerPosition = new THREE.Vector3(0, 0, 22);
  private readonly velocity = new THREE.Vector3();
  private readonly cameraRig = new ThirdPersonCamera();
  private readonly actor: AldricActor;
  private facing = 0;
  private suspended = false;
  private readonly camera = new THREE.PerspectiveCamera(67, 1, 0.08, 130);
  private readonly world: SanctuaryWorld;
  private readonly events = new AbortController();
  private frameId = 0;
  private previousFrame = performance.now();
  private disposed = false;
  private playerSpeed = 0;
  private readonly callbacks: LevelCallbacks;
  private readonly keys = new Set<string>();
  private readonly enemies: Enemy[] = [];
  private readonly storyTriggers = new Map<string, boolean>();
  private readonly moveVector = new THREE.Vector3();
  private readonly dodgeVector = new THREE.Vector3();
  private readonly forward = new THREE.Vector3();
  private readonly right = new THREE.Vector3();
  private readonly equipment = new AldricEquipment();
  private readonly resizeObserver: ResizeObserver;
  private yaw = 0;
  private pitch = 0.25;
  private health = MAX_HEALTH;
  private stamina = MAX_STAMINA;
  private blocking = false;
  private attackTime = 0;
  private attackCooldown = 0;
  private damageFlash = 0;
  private blockFlash = 0;
  private dodgeTime = 0;
  private activatedWards = 0;
  private gateOpen = false;
  private decisionOpen = false;
  private completed = false;
  private running = false;
  private autoplay = false;
  private autoplayIndex = 0;
  private elapsed = 0;

  private constructor(
    canvas: HTMLCanvasElement,
    callbacks: LevelCallbacks,
    textures: { stone: THREE.Texture; earth: THREE.Texture; wood: THREE.Texture },
  ) {
    this.callbacks = callbacks;
    this.renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    this.renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    this.renderer.outputColorSpace = THREE.SRGBColorSpace;
    this.renderer.toneMapping = THREE.ACESFilmicToneMapping;
    this.renderer.toneMappingExposure = 1.24;
    this.renderer.shadowMap.enabled = true;
    this.renderer.shadowMap.type = THREE.PCFShadowMap;
    this.scene.background = new THREE.Color(0x0b1718);
    this.scene.fog = new THREE.FogExp2(0x0b1718, 0.019);

    this.configureTexture(textures.stone, 5, 5);
    this.configureTexture(textures.earth, 18, 24);
    this.configureTexture(textures.wood, 2, 2);
    this.world = new SanctuaryWorld(this.scene, textures);
    this.createEnemies();
    this.actor = new AldricActor(this.equipment.weapon, this.equipment.shield);
    this.scene.add(this.actor.root, this.camera);

    this.resizeObserver = new ResizeObserver(() => this.resize());
    this.resizeObserver.observe(canvas);
    window.addEventListener("keydown", this.onKeyDown, { signal: this.events.signal });
    window.addEventListener("keyup", this.onKeyUp, { signal: this.events.signal });
    window.addEventListener("mousemove", this.onMouseMove, { signal: this.events.signal });
    window.addEventListener("mousedown", this.onMouseDown, { signal: this.events.signal });
    window.addEventListener("mouseup", this.onMouseUp, { signal: this.events.signal });
    canvas.addEventListener("contextmenu", this.onContextMenu, { signal: this.events.signal });
    window.addEventListener("blur", this.suspend, { signal: this.events.signal });
    document.addEventListener("visibilitychange", () => { if (document.hidden) this.suspend(); }, { signal: this.events.signal });
    canvas.addEventListener("click", this.capturePointer, { signal: this.events.signal });
    this.resize();
    this.reset();
    this.animate();
  }

  static async create(canvas: HTMLCanvasElement, callbacks: LevelCallbacks): Promise<ThirdPersonLevel> {
    const loader = new THREE.TextureLoader();
    const loaded = await Promise.allSettled([
      loader.loadAsync("/assets/textures/valethra-stone-albedo.png"),
      loader.loadAsync("/assets/textures/valethra-earth-albedo.png"),
      loader.loadAsync("/assets/textures/valethra-wood-albedo.png"),
    ]);
    const failed = loaded.find(result => result.status === "rejected");
    if (failed) {
      for (const result of loaded) if (result.status === "fulfilled") result.value.dispose();
      throw failed.reason;
    }
    const [stone, earth, wood] = loaded.map(result => (result as PromiseFulfilledResult<THREE.Texture>).value);
    try { return new ThirdPersonLevel(canvas, callbacks, { stone: stone!, earth: earth!, wood: wood! }); }
    catch (error) { for (const texture of [stone, earth, wood]) texture?.dispose(); throw error; }
  }

  start(options: { autoplay?: boolean } = {}): void {
    this.reset();
    this.autoplay = options.autoplay ?? false;
    this.running = true;
    this.previousFrame = performance.now();
    this.callbacks.onMessage(
      "El Santuario del Umbral",
      "Ocho meses después de partir, Aldric encuentra cerrado el último camino hacia su pueblo.",
    );
  }

  pause(): void {
    this.running = false;
    this.keys.clear();
    this.blocking = false;
  }

  chooseEnding(ending: LevelEnding): void {
    if (!this.decisionOpen) return;
    this.decisionOpen = false;
    this.completed = true;
    this.gateOpen = ending === "open";
    this.callbacks.onComplete(ending);
  }

  dispose(): void {
    if (this.disposed) return;
    this.disposed = true;
    this.pause();
    cancelAnimationFrame(this.frameId);
    this.events.abort();
    this.resizeObserver.disconnect();
    if (document.pointerLockElement === this.renderer.domElement) void document.exitPointerLock();
    disposeScene(this.scene);
    this.renderer.dispose();
  }

  private reset(): void {
    this.health = MAX_HEALTH;
    this.stamina = MAX_STAMINA;
    this.activatedWards = 0;
    this.gateOpen = false;
    this.decisionOpen = false;
    this.completed = false;
    this.autoplayIndex = 0;
    this.yaw = 0;
    this.pitch = 0.25;
    this.elapsed = 0;
    this.attackTime = 0;
    this.attackCooldown = 0;
    this.damageFlash = 0;
    this.blockFlash = 0;
    this.dodgeTime = 0;
    this.dodgeVector.set(0,0,0);
    this.blocking = false;
    this.playerPosition.set(0, 0, 22);
    this.velocity.set(0, 0, 0);
    this.keys.clear();
    this.facing = 0;
    this.actor.root.rotation.set(0, 0, 0);
    this.cameraRig.reset();
    this.suspended = false;
    this.camera.rotation.set(0, 0, 0);
    this.world.reset();
    this.playerSpeed = 0;
    for (const enemy of this.enemies) {
      Object.assign(enemy, createEnemyState());
      enemy.defenseDemonstrated = false;
      enemy.actor.reset();
      const [x, z] = ENEMY_SPAWNS[enemy.index]!;
      enemy.actor.root.position.set(x, 0, z);
    }
    for (const key of this.storyTriggers.keys()) this.storyTriggers.set(key, false);
    this.emitHud();
  }

  private configureTexture(texture: THREE.Texture, repeatX: number, repeatY: number): void {
    texture.colorSpace = THREE.SRGBColorSpace;
    texture.wrapS = THREE.RepeatWrapping;
    texture.wrapT = THREE.RepeatWrapping;
    texture.repeat.set(repeatX, repeatY);
    texture.anisotropy = this.renderer.capabilities.getMaxAnisotropy();
  }

  private createEnemies(): void {
    ENEMY_SPAWNS.forEach((_, index) => {
      const actor = new TrollActor(index);
      this.scene.add(actor.root);
      this.enemies.push({ ...createEnemyState(), actor, index, defenseDemonstrated: false });
    });
  }

  private readonly capturePointer = (): void => {
    if (this.running && !this.autoplay && !window.matchMedia("(pointer: coarse)").matches && document.pointerLockElement !== this.renderer.domElement) {
      void this.renderer.domElement.requestPointerLock()?.catch(() => {
        this.callbacks.onPrompt("Flechas izquierda/derecha · Girar cámara");
      });
    }
  };

  private readonly suspend = (): void => {
    if (!this.running || this.autoplay) return;
    this.suspended = true;
    this.pause();
    this.velocity.set(0, 0, 0);
    if (document.pointerLockElement) void document.exitPointerLock();
    this.callbacks.onPause();
  };

  input(code: string, pressed: boolean): void {
    if (!pressed) { this.keys.delete(code); if (code === "KeyF") this.blocking = false; return; }
    if (!this.running) return;
    this.keys.add(code);
    if (code === "Space") this.attack();
    if (code === "KeyQ") this.dodge();
    if (code === "KeyE") this.interact();
    if (code === "KeyF") this.blocking = true;
    if (code === "Escape") this.suspend();
  }

  resume(): void {
    if (!this.suspended) return;
    this.suspended = false;
    this.running = true;
    this.previousFrame = performance.now();
  }

  private readonly onKeyDown = (event: KeyboardEvent): void => {
    if (event.code === "Escape") { this.suspend(); return; }
    if (!this.running) return;
    if (["Space", "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown"].includes(event.code)) event.preventDefault();
    if (!event.repeat) this.input(event.code, true);
  };

  private readonly onKeyUp = (event: KeyboardEvent): void => {
    this.input(event.code, false);
  };

  private readonly onMouseMove = (event: MouseEvent): void => {
    if (!this.running || document.pointerLockElement !== this.renderer.domElement) return;
    this.yaw -= event.movementX * 0.0022;
    this.pitch = THREE.MathUtils.clamp(this.pitch + event.movementY * 0.0018, -0.05, 0.95);
  };

  private readonly onMouseDown = (event: MouseEvent): void => {
    if (!this.running || event.target !== this.renderer.domElement) return;
    if (event.button === 0) this.attack();
    if (event.button === 2) this.blocking = true;
  };

  private readonly onMouseUp = (event: MouseEvent): void => {
    if (event.button === 2) this.blocking = false;
  };

  private readonly onContextMenu = (event: MouseEvent): void => event.preventDefault();

  private attack(): void {
    if (!this.running || this.decisionOpen || this.completed || this.attackCooldown > 0 || this.blocking || this.blockFlash > 0 || this.dodgeTime > 0) return;
    if (!canSpendStamina(this.stamina, ATTACK_STAMINA_COST)) return;
    this.stamina = spendStamina(this.stamina, ATTACK_STAMINA_COST);
    this.attackTime = 0.34;
    this.attackCooldown = 0.48;
    this.callbacks.onCombatCue("attack");
    this.facing = this.yaw;
    this.forward.set(-Math.sin(this.facing), 0, -Math.cos(this.facing));
    let hit = false;
    for (const enemy of this.enemies) {
      if (!enemy.alive) continue;
      const toEnemy = enemy.actor.root.position.clone().sub(this.playerPosition);
      const distance = toEnemy.length();
      if (distance > 3.25 || toEnemy.normalize().dot(this.forward) < 0.55 || !this.clearReach(enemy.actor.root.position)) continue;
      enemy.health -= 38;
      const pushed = enemy.actor.root.position.clone().addScaledVector(this.forward, 0.55);
      if (!this.world.collides(pushed.x, pushed.z, 0.6, this.gateOpen)) enemy.actor.root.position.copy(pushed);
      enemy.hitFlash = 0.22;
      enemy.stagger = 0.26;
      this.callbacks.onCombatCue("hit");
      hit = true;
      if (enemy.health <= 0) {
        enemy.alive = false;
        enemy.deathTime = ENEMY_DEATH_DURATION;
        this.health = Math.min(this.health + 12, MAX_HEALTH);
        this.callbacks.onMessage("El troll cae", "El peso del cuerpo levanta el polvo del camino.");
      }
      break;
    }
    if (!hit) this.callbacks.onPrompt("El golpe corta solamente la niebla");
    this.emitHud();
  }

  private dodge(): void {
    if (!this.running || this.decisionOpen || this.completed || this.dodgeTime > 0 || !canSpendStamina(this.stamina, DODGE_STAMINA_COST)) return;
    this.stamina = spendStamina(this.stamina, DODGE_STAMINA_COST);
    this.forward.set(-Math.sin(this.yaw), 0, -Math.cos(this.yaw));
    this.right.set(Math.cos(this.yaw), 0, -Math.sin(this.yaw));
    this.dodgeVector.set(0, 0, 0);
    if (this.keys.has("KeyW")) this.dodgeVector.add(this.forward);
    if (this.keys.has("KeyS")) this.dodgeVector.sub(this.forward);
    if (this.keys.has("KeyD")) this.dodgeVector.add(this.right);
    if (this.keys.has("KeyA")) this.dodgeVector.sub(this.right);
    if (this.dodgeVector.lengthSq() === 0) this.dodgeVector.add(this.right);
    this.dodgeVector.normalize();
    this.dodgeTime = DODGE_DURATION;
    this.blocking = false;
    this.callbacks.onCombatCue("dodge");
    this.callbacks.onPrompt("Esquiva · invulnerable");
    this.emitHud();
  }

  private interact(): void {
    if (!this.running || this.decisionOpen || this.completed) return;
    const nearestWard = this.nearestWard();
    if (nearestWard && !nearestWard.active && this.horizontalDistance(nearestWard.position) < 2.6 && this.clearReach(nearestWard.position)) {
      nearestWard.active = true;
      this.activatedWards += 1;
      nearestWard.light.intensity = 12;
      nearestWard.core.scale.setScalar(1);
      nearestWard.core.material.emissiveIntensity = 3.2;
      this.health = Math.min(this.health + 25, MAX_HEALTH);
      this.stamina = MAX_STAMINA;
      this.callbacks.onMessage(nearestWard.title, nearestWard.message);
      this.emitHud();
      return;
    }

    if (this.distanceTo(0, -32.4) < 4.2) {
      if (!canOpenWardGate(this.activatedWards)) {
        this.callbacks.onPrompt("El portón no responde: faltan sellos");
        return;
      }
      this.decisionOpen = true;
      this.running = false;
      this.keys.clear();
      if (document.pointerLockElement) void document.exitPointerLock();
      this.callbacks.onDecision();
    }
  }

  private clearReach(target: THREE.Vector3): boolean {
    const origin = this.playerPosition.clone().add(new THREE.Vector3(0, 1.1, 0));
    const direction = target.clone().add(new THREE.Vector3(0, 1.1, 0)).sub(origin);
    this.scene.updateMatrixWorld(true);
    const ray = new THREE.Raycaster(origin, direction.clone().normalize(), 0, direction.length());
    return ray.intersectObjects(this.world.cameraSolids, true).length === 0;
  }

  private nearestWard(): Ward | undefined {
    let nearest: Ward | undefined;
    let distance = Infinity;
    for (const ward of this.world.wards) {
      const candidate = this.horizontalDistance(ward.position);
      if (!ward.active && candidate < distance) { nearest = ward; distance = candidate; }
    }
    return nearest;
  }

  private animate = (timestamp = performance.now()): void => {
    if (this.disposed) return;
    this.frameId = requestAnimationFrame(this.animate);
    const delta = Math.max(0, Math.min((timestamp - this.previousFrame) / 1000, 0.05));
    this.previousFrame = timestamp;
    this.elapsed += delta;
    if (this.running) this.update(delta);
    this.updateVisuals(delta);
    this.renderer.render(this.scene, this.camera);
  };

  private update(delta: number): void {
    if (this.autoplay) this.updateAutoplay();
    if (!this.running) return;
    this.updateMovement(delta);
    this.updateEnemies(delta);
    this.updateStoryTriggers();
    this.updatePrompt();
    this.stamina = Math.min(this.stamina + delta * (this.blocking ? 5 : 15), MAX_STAMINA);
    this.attackCooldown = Math.max(this.attackCooldown - delta, 0);
    this.attackTime = Math.max(this.attackTime - delta, 0);
    this.dodgeTime = Math.max(this.dodgeTime - delta, 0);
    this.damageFlash = Math.max(this.damageFlash - delta, 0);
    this.blockFlash = Math.max(this.blockFlash - delta, 0);
    this.emitHud();
  }

  private updateMovement(delta: number): void {
    if (this.keys.has("ArrowLeft")) this.yaw += delta * 1.8;
    if (this.keys.has("ArrowRight")) this.yaw -= delta * 1.8;
    if (this.keys.has("ArrowUp")) this.pitch = Math.min(0.95, this.pitch + delta);
    if (this.keys.has("ArrowDown")) this.pitch = Math.max(-0.05, this.pitch - delta);
    this.forward.set(-Math.sin(this.yaw), 0, -Math.cos(this.yaw));
    this.right.set(Math.cos(this.yaw), 0, -Math.sin(this.yaw));
    this.moveVector.set(0, 0, 0);
    if (this.keys.has("KeyW")) this.moveVector.add(this.forward);
    if (this.keys.has("KeyS")) this.moveVector.sub(this.forward);
    if (this.keys.has("KeyD")) this.moveVector.add(this.right);
    if (this.keys.has("KeyA")) this.moveVector.sub(this.right);
    const moving = this.moveVector.lengthSq() > 0;
    const sprinting = moving && this.keys.has("ShiftLeft") && this.stamina > 0 && !this.blocking && this.dodgeTime <= 0;
    const speed = sprinting ? SPRINT_SPEED : WALK_SPEED;
    if (this.dodgeTime > 0) {
      this.velocity.copy(this.dodgeVector).multiplyScalar(DODGE_SPEED * (0.72 + this.dodgeTime / DODGE_DURATION * 0.28));
    } else {
      this.moveVector.normalize().multiplyScalar(speed);
      this.velocity.lerp(this.moveVector, damping(moving ? 14 : 20, delta));
    }
    if (sprinting) this.stamina = Math.max(this.stamina - delta * 19, 0);
    if (this.blocking || this.attackTime > 0) this.facing = this.yaw;
    else if (this.velocity.lengthSq() > 0.05) this.facing = Math.atan2(-this.velocity.x, -this.velocity.z);
    const previous = this.playerPosition.clone();
    const nextX = this.playerPosition.x + this.velocity.x * delta;
    const nextZ = this.playerPosition.z + this.velocity.z * delta;
    if (!this.collides(nextX, this.playerPosition.z)) this.playerPosition.x = nextX;
    else this.velocity.x = 0;
    if (!this.collides(this.playerPosition.x, nextZ)) this.playerPosition.z = nextZ;
    else this.velocity.z = 0;
    // Animation uses actual travel so the feet stop against a wall.
    this.playerSpeed = previous.distanceTo(this.playerPosition) / Math.max(delta, 0.001);
  }

  private collides(x: number, z: number): boolean {
    return this.world.collides(x, z, PLAYER_RADIUS, this.gateOpen);
  }

  private updateEnemies(delta: number): void {
    for (const enemy of this.enemies) {
      const toPlayer = this.playerPosition.clone().sub(enemy.actor.root.position);
      toPlayer.y = 0;
      const distance = toPlayer.length();
      const wasPreparing = enemy.attackWindup > 0;
      const transition = advanceEnemy(enemy, delta, distance);
      const moving = enemy.alive && !wasPreparing && distance < 10 && distance > 1.75 && enemy.stagger <= 0;
      if (moving) {
        const position = enemy.actor.root.position;
        const travel = toPlayer.clone().normalize().multiplyScalar(delta * 1.22);
        if (!this.world.collides(position.x + travel.x, position.z, 0.6, this.gateOpen)) position.x += travel.x;
        if (!this.world.collides(position.x, position.z + travel.z, 0.6, this.gateOpen)) position.z += travel.z;
      }
      if (transition === "windup") this.callbacks.onCombatCue("enemy-windup");
      if (transition === "strike") {
        enemy.defenseDemonstrated = true;
        if (distance < 2.35 && this.dodgeTime <= 0 && this.clearReach(enemy.actor.root.position)) {
          const result = resolveIncomingHit(this.health, this.stamina, ENEMY_DAMAGE, this.blocking);
          this.health = result.health;
          this.stamina = result.stamina;
          if (this.blocking) {
            this.blockFlash = 0.2;
            const guardBroken = result.stamina === 0;
            this.callbacks.onCombatCue(guardBroken ? "guard-break" : "block");
            this.callbacks.onPrompt(guardBroken ? "¡Guardia rota!" : "Bloqueo · daño reducido");
          } else {
            this.damageFlash = 0.32;
            this.callbacks.onCombatCue("hit");
            this.callbacks.onPrompt("El troll te alcanza");
          }
          if (this.health <= 0) {
            this.pause();
            this.callbacks.onDeath();
            return;
          }
        }
      }
      enemy.actor.update(enemy, this.elapsed, delta, Math.atan2(toPlayer.x, toPlayer.z), moving);
    }
  }

  private updateAutoplay(): void {
    const targets: Array<{ position: THREE.Vector3; interact?: boolean }> = [
      { position: new THREE.Vector3(-13, 0, 14) },
      { position: this.world.wards[0]!.position, interact: true },
      { position: new THREE.Vector3(0, 0, 5) },
      { position: new THREE.Vector3(6, 0, 0) },
      { position: this.world.wards[1]!.position, interact: true },
      { position: new THREE.Vector3(7, 0, -16) },
      { position: new THREE.Vector3(-6, 0, -20) },
      { position: this.world.wards[2]!.position, interact: true },
      { position: new THREE.Vector3(0, 0, -30) },
      { position: new THREE.Vector3(0, 0, -31.2), interact: true },
    ];
    const target = targets[Math.min(this.autoplayIndex, targets.length - 1)]!;
    const direction = target.position.clone().sub(this.playerPosition);
    direction.y = 0;
    const distance = direction.length();
    this.yaw = Math.atan2(-direction.x, -direction.z);
    this.pitch = 0.25;
    this.keys.add("KeyW");
    if (this.stamina > 58) this.keys.add("ShiftLeft");
    else this.keys.delete("ShiftLeft");
    const nearbyEnemy = this.enemies.find((enemy) => enemy.alive && enemy.actor.root.position.distanceTo(this.playerPosition) < 3.1 && this.clearReach(enemy.actor.root.position));
    if (nearbyEnemy) {
      const enemyDirection = nearbyEnemy.actor.root.position.clone().sub(this.playerPosition);
      this.yaw = Math.atan2(-enemyDirection.x, -enemyDirection.z);
      this.keys.delete("ShiftLeft");
      this.keys.delete("KeyW");
      if (!nearbyEnemy.defenseDemonstrated && nearbyEnemy.attackWindup > 0) {
        if (nearbyEnemy.index % 2 === 0) {
          this.blocking = true;
        } else {
          this.blocking = false;
          if (nearbyEnemy.attackWindup < 0.42 && this.dodgeTime <= 0) this.dodge();
        }
      } else if (!nearbyEnemy.defenseDemonstrated) {
        this.blocking = false;
      } else {
        this.blocking = false;
        this.attack();
      }
      return;
    }
    this.blocking = false;
    if (distance < 2.2) {
      this.keys.delete("KeyW");
      this.keys.delete("ShiftLeft");
      if (target.interact) this.interact();
      if (this.autoplayIndex < targets.length - 1) this.autoplayIndex += 1;
    }
  }

  private updateStoryTriggers(): void {
    this.triggerStory("cart", -7.2, 12, 5, "Un regreso interrumpido", "El carro quedó orientado hacia el pueblo. Sus dueños intentaban llegar, no escapar.");
    this.triggerStory("shields", 8, 8.4, 4.5, "La retirada", "Los escudos fueron arrojados mirando hacia el santuario. La guardia abandonó su puesto con prisa.");
    this.triggerStory("crib", -2.4, -27.8, 4, "Una espera demasiado larga", "La pequeña cuna conserva barro reciente. La familia está cerca.");
  }

  private triggerStory(id: string, x: number, z: number, radius: number, title: string, body: string): void {
    if (this.storyTriggers.get(id) || this.distanceTo(x, z) > radius) return;
    this.storyTriggers.set(id, true);
    this.callbacks.onMessage(title, body);
  }

  private updatePrompt(): void {
    const nearestWard = this.nearestWard();
    if (nearestWard && this.horizontalDistance(nearestWard.position) < 2.6 && this.clearReach(nearestWard.position)) {
      this.callbacks.onPrompt(`E · Restaurar ${nearestWard.title.toLowerCase()}`);
      return;
    }
    if (this.distanceTo(0, -32.4) < 4.2) {
      this.callbacks.onPrompt(canOpenWardGate(this.activatedWards) ? "E · Decidir el destino del portón" : "El portón exige tres sellos");
      return;
    }
    this.callbacks.onPrompt("");
  }

  private updateVisuals(delta: number): void {
    const speed = this.running ? this.playerSpeed : 0;
    this.cameraRig.update(this.camera, this.playerPosition, this.yaw, this.pitch, this.world.cameraSolids, delta, speed > 4.5);
    this.actor.update(this.playerPosition, this.facing, speed, this.elapsed, delta, this.attackTime, this.blocking, this.dodgeTime);
    this.world.lantern.position.copy(this.playerPosition).y = 1.3;
    this.equipment.updateImpact(this.blockFlash);
    const doors = this.world.doors;
    const targetY = this.gateOpen ? 5.6 : 0;
    doors.position.y += (targetY - doors.position.y) * Math.min(delta * 2.2, 1);
    if (this.gateOpen) this.world.family.visible = true;
    for (const ward of this.world.wards) {
      if (ward.active) {
        const core = ward.core;
        core.rotation.y += delta * 1.4;
        ward.light.intensity = 10 + Math.sin(this.elapsed * 4 + core.id) * 2;
      }
    }
    this.renderer.domElement.style.filter = this.damageFlash > 0
      ? "sepia(.3) saturate(1.8) brightness(.8)"
      : this.dodgeTime > 0 ? "saturate(1.35) brightness(1.08)" : "none";
  }

  private emitHud(): void {
    this.callbacks.onHud({
      health: this.health,
      stamina: this.stamina,
      wards: this.activatedWards,
      objective: objectiveForWards(this.activatedWards),
      blocking: this.blocking,
      dodging: this.dodgeTime > 0,
      threatened: this.enemies.some((enemy) => enemy.alive && enemy.attackWindup > 0),
      guardImpact: this.blockFlash > 0,
    });
  }

  private horizontalDistance(position: THREE.Vector3): number {
    return this.distanceTo(position.x, position.z);
  }

  private distanceTo(x: number, z: number): number {
    return Math.hypot(this.playerPosition.x - x, this.playerPosition.z - z);
  }

  private resize(): void {
    const canvas = this.renderer.domElement;
    const width = Math.max(canvas.clientWidth, 1);
    const height = Math.max(canvas.clientHeight, 1);
    this.camera.aspect = width / height;
    this.camera.updateProjectionMatrix();
    this.renderer.setSize(width, height, false);
  }
}
