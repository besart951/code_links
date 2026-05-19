import { isLocale } from '$lib/site';

export function match(value: string): boolean {
  return isLocale(value);
}
