import { describe, expect, it } from "vitest";

import {
  ATTACK_STAMINA_COST,
  canOpenWardGate,
  canSpendStamina,
  objectiveForWards,
  resolveIncomingHit,
  spendStamina,
} from "./combat-rules";

describe("sanctuary combat rules", () => {
  it("only spends stamina when an action is affordable", () => {
    expect(canSpendStamina(18, ATTACK_STAMINA_COST)).toBe(true);
    expect(spendStamina(18, ATTACK_STAMINA_COST)).toBe(0);
    expect(spendStamina(17, ATTACK_STAMINA_COST)).toBe(17);
  });

  it("lets a healthy guard trade stamina for reduced damage", () => {
    expect(resolveIncomingHit(100, 100, 20, true)).toEqual({
      health: 96,
      stamina: 82,
      damageTaken: 4,
    });
  });

  it("breaks an exhausted guard without underflowing resources", () => {
    const result = resolveIncomingHit(12, 4, 30, true);
    expect(result.stamina).toBe(0);
    expect(result.health).toBeGreaterThanOrEqual(0);
    expect(result.damageTaken).toBeGreaterThan(9);
  });

  it("opens the ward gate after all three seals", () => {
    expect(canOpenWardGate(2)).toBe(false);
    expect(canOpenWardGate(3)).toBe(true);
  });

  it("describes progress without exposing implementation details", () => {
    expect(objectiveForWards(0)).toBe("Restaura 3 sellos");
    expect(objectiveForWards(2)).toBe("Restaura el último sello");
    expect(objectiveForWards(3)).toBe("Regresa al portón del santuario");
  });
});
