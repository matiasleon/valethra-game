export const MAX_HEALTH = 100;
export const MAX_STAMINA = 100;
export const ATTACK_STAMINA_COST = 18;
export const DODGE_STAMINA_COST = 28;
export const REQUIRED_WARDS = 3;

export interface DefensiveResult {
  health: number;
  stamina: number;
  damageTaken: number;
}

export function clampResource(value: number, maximum = 100): number {
  return Math.min(Math.max(value, 0), maximum);
}

export function canSpendStamina(stamina: number, cost: number): boolean {
  return stamina >= cost;
}

export function spendStamina(stamina: number, cost: number): number {
  return canSpendStamina(stamina, cost) ? stamina - cost : stamina;
}

export function resolveIncomingHit(
  health: number,
  stamina: number,
  rawDamage: number,
  blocking: boolean,
): DefensiveResult {
  if (!blocking) {
    const damageTaken = Math.max(rawDamage, 0);
    return { health: clampResource(health - damageTaken), stamina, damageTaken };
  }

  const guardCost = Math.max(rawDamage * 0.9, 0);
  if (stamina >= guardCost) {
    const damageTaken = rawDamage * 0.2;
    return {
      health: clampResource(health - damageTaken),
      stamina: clampResource(stamina - guardCost),
      damageTaken,
    };
  }

  const prevented = stamina / 0.9;
  const damageTaken = Math.max(rawDamage - prevented * 0.8, rawDamage * 0.35);
  return { health: clampResource(health - damageTaken), stamina: 0, damageTaken };
}

export function canOpenWardGate(activatedWards: number): boolean {
  return activatedWards >= REQUIRED_WARDS;
}

export function objectiveForWards(activatedWards: number): string {
  const remaining = Math.max(REQUIRED_WARDS - activatedWards, 0);
  if (remaining === 0) return "Regresa al portón del santuario";
  return remaining === 1 ? "Restaura el último sello" : `Restaura ${remaining} sellos`;
}
