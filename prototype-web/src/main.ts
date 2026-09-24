import "./styles.css";

import { AmbientScore } from "./ambient-audio";
import { ThirdPersonLevel } from "./third-person-level";
import type { HudState, LevelEnding } from "./level-contracts";
import { CanvasRecorder } from "./qa-recorder";

declare global {
  interface Window {
    __valethraQA: {
      getRecordingDataUrl: () => Promise<string | undefined>;
      isReady: () => boolean;
      status: () => string;
    };
  }
}

const game = requireElement<HTMLElement>("game");
const canvas = requireElement<HTMLCanvasElement>("world");
const hud = requireElement<HTMLElement>("hud");
const overlay = requireElement<HTMLElement>("overlay");
const objective = requireElement<HTMLElement>("objective");
const wardCount = requireElement<HTMLElement>("ward-count");
const healthBar = requireElement<HTMLElement>("health");
const healthValue = requireElement<HTMLElement>("health-value");
const staminaBar = requireElement<HTMLElement>("stamina");
const staminaValue = requireElement<HTMLElement>("stamina-value");
const combatState = requireElement<HTMLElement>("combat-state");
const combatStateTitle = combatState.querySelector<HTMLElement>("strong")!;
const combatStateDetail = combatState.querySelector<HTMLElement>("small")!;
const prompt = requireElement<HTMLElement>("prompt");
const message = requireElement<HTMLElement>("message");
const messageTitle = requireElement<HTMLElement>("message-title");
const messageBody = requireElement<HTMLElement>("message-body");
const wards = [...document.querySelectorAll<HTMLElement>(".wards i")];
const query = new URLSearchParams(location.search);
const autoplay = query.get("autoplay") === "1";
const shouldRecord = query.get("record") === "1";
const score = new AmbientScore();
const recorder = new CanvasRecorder();
let level: ThirdPersonLevel;
let messageTimer = 0;
let endingTimer = 0;
let recordingTimer = 0;
let runId = 0;
let starting = false;
let disposed = false;
const events = new AbortController();

function requireElement<T extends HTMLElement>(id: string): T {
  const element = document.getElementById(id);
  if (!element) throw new Error(`Falta el elemento #${id}`);
  return element as T;
}

function updateHud(state: HudState): void {
  const health = Math.round(state.health);
  const stamina = Math.round(state.stamina);
  objective.textContent = state.objective;
  wardCount.textContent = `${state.wards} / 3`;
  healthBar.style.width = `${health}%`;
  healthValue.textContent = String(health);
  staminaBar.style.width = `${stamina}%`;
  staminaValue.textContent = String(stamina);
  game.dataset.blocking = String(state.blocking);
  game.dataset.dodging = String(state.dodging);
  game.dataset.threatened = String(state.threatened);
  game.dataset.guardImpact = String(state.guardImpact);
  if (state.dodging) {
    combatStateTitle.textContent = "ESQUIVA";
    combatStateDetail.textContent = "Invulnerable durante el impulso";
  } else if (state.guardImpact) {
    combatStateTitle.textContent = "IMPACTO BLOQUEADO";
    combatStateDetail.textContent = "Daño reducido · resistencia consumida";
  } else if (state.blocking) {
    combatStateTitle.textContent = "GUARDIA ACTIVA";
    combatStateDetail.textContent = "Soltá F para atacar";
  } else if (state.threatened) {
    combatStateTitle.textContent = "ATAQUE ENTRANTE";
    combatStateDetail.textContent = "Bloqueá con F o esquivá con Q";
  } else {
    combatStateTitle.textContent = "";
    combatStateDetail.textContent = "";
  }
  combatState.classList.toggle("is-visible", Boolean(combatStateTitle.textContent));
  wards.forEach((ward, index) => ward.classList.toggle("is-active", index < state.wards));
}

function showMessage(title: string, body: string): void {
  window.clearTimeout(messageTimer);
  messageTitle.textContent = title;
  messageBody.textContent = body;
  message.classList.add("is-visible");
  messageTimer = window.setTimeout(() => message.classList.remove("is-visible"), 5600);
}

function showIntro(): void {
  game.dataset.phase = "intro";
  overlay.hidden = false;
  overlay.innerHTML = `
    <article class="intro-card">
      <p class="eyebrow">Valethra · Una jornada</p>
      <h1>El Santuario<br>del Umbral</h1>
      <p class="lead">Aldric regresa a su pueblo después de ocho meses en una misión peligrosa.</p>
      <p>El sendero termina ante un santuario abandonado. El portón está sellado, la guardia desapareció y algo se mueve entre los árboles.</p>
      <button type="button" data-action="start">Entrar al santuario <kbd>Click</kbd></button>
      <small>Un nivel jugable en tercera persona · Usá auriculares</small>
    </article>`;
  overlay.querySelector<HTMLButtonElement>("[data-action='start']")?.addEventListener("click", () => void startGame(), { once: true });
}

async function startGame(): Promise<void> {
  if (starting || disposed) return;
  starting = true;
  runId++;
  window.clearTimeout(endingTimer);
  window.clearTimeout(recordingTimer);
  if (shouldRecord) await recorder.stop();
  game.dataset.phase = "playing";
  overlay.hidden = true;
  hud.hidden = false;
  if (!autoplay) {
    try { await score.start(); }
    catch { showMessage("Sonido no disponible", "Podés continuar jugando sin sonido."); }
  }
  if (disposed) { score.stop(); starting = false; return; }
  if (shouldRecord) {
    recorder.start(canvas, 60);
    game.dataset.recorder = recorder.state;
  }
  level.start({ autoplay });
  starting = false;
}

function showPause(): void {
  game.dataset.phase = "paused";
  overlay.hidden = false;
  overlay.innerHTML = `<article class="ending-card"><p class="eyebrow">Valethra</p><h2>Una pausa en el camino</h2><p>WASD para moverte · Mouse o flechas para mirar.<br>Espacio: atacar · F: bloquear · Q: esquivar · E: interactuar.</p><button type="button" data-action="resume">Continuar</button></article>`;
  overlay.querySelector("button")?.addEventListener("click", () => {
    overlay.hidden = true; game.dataset.phase = "playing"; level.resume();
  }, { once: true });
}

function showDecision(): void {
  game.dataset.phase = "decision";
  overlay.hidden = false;
  overlay.innerHTML = `
    <article class="decision-card">
      <p class="eyebrow">Los tres sellos responden</p>
      <h2>El portón recuerda su propósito</h2>
      <p>Del otro lado espera una familia. Abrir el camino puede llevarla a salvo, pero también liberará aquello que la guardia encerró.</p>
      <div class="decision-actions">
        <button type="button" data-ending="open"><span>Abrir el portón</span><small>Dejar pasar a la familia</small></button>
        <button type="button" data-ending="seal"><span>Sellar el camino</span><small>Contener la amenaza</small></button>
      </div>
    </article>`;
  for (const button of overlay.querySelectorAll<HTMLButtonElement>("[data-ending]")) {
    button.addEventListener("click", () => level.chooseEnding(button.dataset.ending as LevelEnding), { once: true });
  }
  if (autoplay) endingTimer = window.setTimeout(() => level.chooseEnding(query.get("ending") === "seal" ? "seal" : "open"), 1300);
}

function showEnding(ending: LevelEnding): void {
  game.dataset.phase = "complete";
  hud.hidden = true;
  overlay.hidden = false;
  const opened = ending === "open";
  overlay.innerHTML = `
    <article class="ending-card">
      <p class="eyebrow">La jornada termina</p>
      <h2>${opened ? "El camino vuelve a respirar" : "El bosque queda en silencio"}</h2>
      <p>${opened
        ? "La familia cruza. Aldric ve barro del pueblo en sus botas y comprende que su regreso no será el final de esta historia."
        : "Aldric refuerza los sellos. La familia permanece fuera y el secreto de la guardia sigue encerrado con ella."}</p>
      <p class="consequence">${opened ? "Consecuencia: el sendero queda abierto para todos." : "Consecuencia: la amenaza queda contenida, por ahora."}</p>
      <button type="button" data-action="restart">Volver a jugar</button>
    </article>`;
  overlay.querySelector<HTMLButtonElement>("[data-action='restart']")?.addEventListener("click", () => void startGame(), { once: true });
  if (shouldRecord) recordingTimer = window.setTimeout(() => void finalizeRecording(), autoplay ? 2600 : 0);
}

function showDeath(): void {
  if (document.pointerLockElement) void document.exitPointerLock();
  game.dataset.phase = "death";
  hud.hidden = true;
  overlay.hidden = false;
  overlay.innerHTML = `
    <article class="ending-card ending-card--death">
      <p class="eyebrow">Aldric cae</p><h2>El santuario conserva su secreto</h2>
      <p>Administrá la resistencia: bloquear reduce el daño, esquivar evita el golpe y atacar sin aire te deja expuesto.</p>
      <button type="button" data-action="restart">Intentarlo otra vez</button>
    </article>`;
  overlay.querySelector<HTMLButtonElement>("[data-action='restart']")?.addEventListener("click", () => void startGame(), { once: true });
  if (shouldRecord) void finalizeRecording();
}

async function finalizeRecording(): Promise<void> {
  const recordingRun = runId;
  await recorder.stop();
  game.dataset.recorder = recorder.state;
  const dataUrl = await recorder.asDataUrl();
  if (!dataUrl || disposed || recordingRun !== runId) return;
  const link = document.createElement("a");
  link.id = "qa-recording";
  link.href = dataUrl;
  link.download = "valethra-autoplay-60fps.webm";
  link.className = "qa-download";
  link.textContent = "Descargar grabación QA";
  overlay.append(link);
}

async function initialize(): Promise<void> {
  level = await ThirdPersonLevel.create(canvas, {
    onHud: updateHud,
    onMessage: showMessage,
    onPrompt: (text) => { prompt.textContent = text; },
    onCombatCue: (cue) => { if (!autoplay) score.playCombatCue(cue); },
    onDecision: showDecision,
    onComplete: showEnding,
    onDeath: showDeath,
    onPause: showPause,
  });
  if (disposed) { level.dispose(); return; }
  for (const button of document.querySelectorAll<HTMLButtonElement>("[data-key]")) {
    button.addEventListener("pointerdown", (event) => {
      event.preventDefault(); button.setPointerCapture(event.pointerId);
      level.input(button.dataset.key!, true);
    }, { signal: events.signal });
    const release = () => level.input(button.dataset.key!, false);
    button.addEventListener("pointerup", release, { signal: events.signal });
    button.addEventListener("pointercancel", release, { signal: events.signal });
    button.addEventListener("lostpointercapture", release, { signal: events.signal });
  }
  window.addEventListener("pagehide", (event) => { if (!event.persisted) dispose(); }, { signal: events.signal });
  window.__valethraQA = {
    getRecordingDataUrl: () => recorder.asDataUrl(),
    isReady: () => recorder.ready,
    status: () => recorder.state,
  };
  if (autoplay) await startGame();
  else showIntro();
}

function dispose(): void {
  if (disposed) return;
  disposed = true;
  runId++;
  window.clearTimeout(endingTimer);
  window.clearTimeout(recordingTimer);
  window.clearTimeout(messageTimer);
  events.abort();
  level?.dispose();
  recorder.dispose();
  score.stop();
}

if (import.meta.hot) import.meta.hot.dispose(dispose);

void initialize().catch((error: unknown) => {
  console.error(error);
  overlay.innerHTML = `<article class="ending-card"><p class="eyebrow">No pudimos abrir Valethra</p><h2>Error al cargar el nivel</h2><p>Recargá la página para volver a intentarlo.</p></article>`;
});
