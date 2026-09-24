import type { CombatCue } from "./level-contracts";

export class AmbientScore {
  private context?: AudioContext;
  private master?: GainNode;
  private sources: AudioScheduledSourceNode[] = [];
  private bellTimer?: number;

  async start(): Promise<void> {
    if (this.context) {
      await this.context.resume();
      return;
    }

    const context = new AudioContext();
    const master = context.createGain();
    master.gain.value = 0.16;
    master.connect(context.destination);
    this.context = context;
    this.master = master;

    this.addDrone(73.42, "sine", 0.07, -5);
    this.addDrone(110, "triangle", 0.035, 4);
    this.addWind();
    this.scheduleBell();
    await context.resume();
  }

  stop(): void {
    if (this.bellTimer !== undefined) window.clearTimeout(this.bellTimer);
    for (const source of this.sources) source.stop();
    this.sources = [];
    void this.context?.close();
    this.context = undefined;
    this.master = undefined;
  }

  playCombatCue(cue: CombatCue): void {
    const context = this.context;
    const master = this.master;
    if (!context || !master) return;
    const profiles: Record<CombatCue, { from: number; to: number; gain: number; duration: number; type: OscillatorType }> = {
      attack: { from: 180, to: 75, gain: 0.07, duration: 0.13, type: "sawtooth" },
      hit: { from: 95, to: 42, gain: 0.12, duration: 0.19, type: "square" },
      block: { from: 520, to: 210, gain: 0.1, duration: 0.22, type: "triangle" },
      dodge: { from: 240, to: 740, gain: 0.055, duration: 0.16, type: "sine" },
      "guard-break": { from: 170, to: 38, gain: 0.13, duration: 0.38, type: "sawtooth" },
      "enemy-windup": { from: 80, to: 230, gain: 0.045, duration: 0.55, type: "sine" },
    };
    const profile = profiles[cue];
    const now = context.currentTime;
    const oscillator = context.createOscillator();
    const gain = context.createGain();
    oscillator.type = profile.type;
    oscillator.frequency.setValueAtTime(profile.from, now);
    oscillator.frequency.exponentialRampToValueAtTime(profile.to, now + profile.duration);
    gain.gain.setValueAtTime(profile.gain, now);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + profile.duration);
    oscillator.connect(gain).connect(master);
    oscillator.start(now);
    oscillator.stop(now + profile.duration);
  }

  private addDrone(frequency: number, type: OscillatorType, volume: number, detune: number): void {
    const context = this.context;
    const master = this.master;
    if (!context || !master) return;

    const oscillator = context.createOscillator();
    const gain = context.createGain();
    const filter = context.createBiquadFilter();
    oscillator.type = type;
    oscillator.frequency.value = frequency;
    oscillator.detune.value = detune;
    gain.gain.value = volume;
    filter.type = "lowpass";
    filter.frequency.value = 430;
    oscillator.connect(filter).connect(gain).connect(master);
    oscillator.start();
    this.sources.push(oscillator);
  }

  private addWind(): void {
    const context = this.context;
    const master = this.master;
    if (!context || !master) return;

    const buffer = context.createBuffer(1, context.sampleRate * 4, context.sampleRate);
    const data = buffer.getChannelData(0);
    for (let index = 0; index < data.length; index += 1) data[index] = Math.random() * 2 - 1;

    const wind = context.createBufferSource();
    const filter = context.createBiquadFilter();
    const gain = context.createGain();
    const lfo = context.createOscillator();
    const lfoGain = context.createGain();
    wind.buffer = buffer;
    wind.loop = true;
    filter.type = "bandpass";
    filter.frequency.value = 310;
    filter.Q.value = 0.7;
    gain.gain.value = 0.018;
    lfo.frequency.value = 0.09;
    lfoGain.gain.value = 0.012;
    lfo.connect(lfoGain).connect(gain.gain);
    wind.connect(filter).connect(gain).connect(master);
    wind.start();
    lfo.start();
    this.sources.push(wind, lfo);
  }

  private scheduleBell(): void {
    const context = this.context;
    const master = this.master;
    if (!context || !master) return;

    const now = context.currentTime;
    const oscillator = context.createOscillator();
    const gain = context.createGain();
    oscillator.type = "sine";
    oscillator.frequency.value = Math.random() > 0.45 ? 293.66 : 220;
    gain.gain.setValueAtTime(0.0001, now);
    gain.gain.exponentialRampToValueAtTime(0.045, now + 0.03);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + 2.8);
    oscillator.connect(gain).connect(master);
    oscillator.start(now);
    oscillator.stop(now + 3);
    this.bellTimer = window.setTimeout(() => this.scheduleBell(), 7200 + Math.random() * 5200);
  }
}
