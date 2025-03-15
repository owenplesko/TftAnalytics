import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function sentenceCase(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
}

export function formatStageNumber(rounds: number) {
  const stage = Math.floor((rounds - 4) / 7) + 2;
  const round = (rounds - 4) % 7;

  return `${stage}-${round}`;
}

// formats timestamp str to most significant unit of time ago
export function formatTimeSince(date: Date) {
  const now = new Date();

  const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  const years = Math.floor(seconds / (60 * 60 * 24 * 365));
  if (years >= 1) return `${years} year${years > 1 ? "s" : ""} ago`;

  const months = Math.floor(seconds / (60 * 60 * 24 * 30));
  if (months >= 1) return `${months} month${months > 1 ? "s" : ""} ago`;

  const days = Math.floor(seconds / (60 * 60 * 24));
  if (days >= 1) return `${days} day${days > 1 ? "s" : ""} ago`;

  const hours = Math.floor(seconds / (60 * 60));
  if (hours >= 1) return `${hours} hour${hours > 1 ? "s" : ""} ago`;

  const minutes = Math.floor(seconds / 60);
  if (minutes >= 1) return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;

  return `${seconds} second${seconds > 1 ? "s" : ""} ago`;
}

// formats float value seconds to string format m:s
export function formatSeconds(seconds: number) {
  const s = Math.floor(seconds % 60);
  const m = Math.floor((seconds % 3600) / 60);

  const str = `${m}:${s}`;

  return str;
}
