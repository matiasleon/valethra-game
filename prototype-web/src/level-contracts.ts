export type LevelEnding = "open" | "seal";
export type CombatCue = "attack" | "hit" | "block" | "dodge" | "guard-break" | "enemy-windup";

export interface HudState {
  health: number;
  stamina: number;
  wards: number;
  objective: string;
  blocking: boolean;
  dodging: boolean;
  threatened: boolean;
  guardImpact: boolean;
}

export interface LevelCallbacks {
  onHud: (state: HudState) => void;
  onMessage: (title: string, body: string) => void;
  onPrompt: (text: string) => void;
  onCombatCue: (cue: CombatCue) => void;
  onDecision: () => void;
  onComplete: (ending: LevelEnding) => void;
  onDeath: () => void;
  onPause: () => void;
}
