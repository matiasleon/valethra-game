/** Frame-rate independent damping, with seconds as the time unit. */
export function damping(rate: number, delta: number): number {
  return 1 - Math.exp(-rate * Math.max(0, delta));
}

export function turnTowards(current: number, target: number, delta: number): number {
  const shortest = Math.atan2(Math.sin(target - current), Math.cos(target - current));
  return current + shortest * damping(16, delta);
}
