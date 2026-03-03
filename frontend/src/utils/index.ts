import type { Fyo } from 'fyo';
import { Money } from 'pesa';

/**
 * And so should not contain and platforma specific imports.
 */
export function getValueMapFromList<T, K extends keyof T, V extends keyof T>(
  list: T[],
  key: K,
  valueKey: V,
  filterUndefined = true
): Record<string, T[V]> {
  if (filterUndefined) {
    list = list.filter(
      (f) =>
        (f[valueKey] as unknown) !== undefined &&
        (f[key] as unknown) !== undefined
    );
  }

  return list.reduce((acc, f) => {
    const keyValue = String(f[key]);
    const value = f[valueKey];
    acc[keyValue] = value;
    return acc;
  }, {} as Record<string, T[V]>);
}

export function getRandomString(): string {
  const randomNumber = Math.random().toString(36).slice(2, 8);
  const currentTime = Date.now().toString(36);
  return `${randomNumber}-${currentTime}`;
}

export async function sleep(durationMilliseconds = 1000) {
  return new Promise((r) => setTimeout(() => r(null), durationMilliseconds));
}

export function getMapFromList<T, K extends keyof T>(
  list: T[],
  name: K
): Record<string, T> {
  /**
   * Do not convert function to use copies of T
   * instead of references.
   */
  const acc: Record<string, T> = {};
  for (const t of list) {
    const key = t[name];
    if (key === undefined) {
      continue;
    }

    acc[String(key)] = t;
  }
  return acc;
}

export function getDefaultMapFromList<T, K extends keyof T, D>(
  list: T[] | string[],
  defaultValue: D,
  name?: K
): Record<string, D> {
  const acc: Record<string, D> = {};
  if (typeof list[0] === 'string') {
    for (const l of list as string[]) {
      acc[l] = defaultValue;
    }

    return acc;
  }

  if (!name) {
    return {};
  }

  for (const l of list as T[]) {
    const key = String(l[name]);
    acc[key] = defaultValue;
  }

  return acc;
}

export function getListFromMap<T>(map: Record<string, T>): T[] {
  return Object.keys(map).map((n) => map[n]);
}

export function getIsNullOrUndef(value: unknown): value is null | undefined {
  return value === null || value === undefined;
}

export function titleCase(phrase: string): string {
  return phrase
    .split(' ')
    .map((word) => {
      const wordLower = word.toLowerCase();
      if (['and', 'an', 'a', 'from', 'by', 'on'].includes(wordLower)) {
        return wordLower;
      }
      return wordLower[0].toUpperCase() + wordLower.slice(1);
    })
    .join(' ');
}

export function invertMap(map: Record<string, string>): Record<string, string> {
  const keys = Object.keys(map);
  const inverted: Record<string, string> = {};
  for (const key of keys) {
    const val = map[key];
    inverted[val] = key;
  }

  return inverted;
}

export function time<K, T>(func: (...args: K[]) => T, ...args: K[]): T {
  /* eslint-disable no-console */
  const name = func.name;
  console.time(name);
  const stuff = func(...args);
  console.timeEnd(name);
  return stuff;
}

export async function timeAsync<K, T>(
  func: (...args: K[]) => Promise<T>,
  ...args: K[]
): Promise<T> {
  /* eslint-disable no-console */
  const name = func.name;
  console.time(name);
  const stuff = await func(...args);
  console.timeEnd(name);
  return stuff;
}

export function changeKeys<T>(
  source: Record<string, T>,
  keyMap: Record<string, string | undefined>
) {
  const dest: Record<string, T> = {};
  for (const key of Object.keys(source)) {
    const newKey = keyMap[key] ?? key;
    dest[newKey] = source[key];
  }

  return dest;
}

export function deleteKeys<T>(
  source: Record<string, T>,
  keysToDelete: string[]
) {
  const dest: Record<string, T> = {};
  for (const key of Object.keys(source)) {
    if (keysToDelete.includes(key)) {
      continue;
    }
    dest[key] = source[key];
  }

  return dest;
}

function safeParseNumber(value: unknown, parser: (v: string) => number) {
  let parsed: number;
  switch (typeof value) {
    case 'string':
      parsed = parser(value);
      break;
    case 'number':
      parsed = value;
      break;
    default:
      parsed = Number(value);
      break;
  }

  if (Number.isNaN(parsed)) {
    return 0;
  }

  return parsed;
}

export function safeParseFloat(value: unknown): number {
  return safeParseNumber(value, Number);
}

export function safeParseInt(value: unknown): number {
  return safeParseNumber(value, (v: string) => Math.trunc(Number(v)));
}

export function safeParsePesa(value: unknown, fyo: Fyo): Money {
  if (value instanceof Money) {
    return value;
  }

  if (typeof value === 'number') {
    return fyo.pesa(value);
  }

  if (typeof value === 'bigint') {
    return fyo.pesa(value);
  }

  if (typeof value !== 'string') {
    return fyo.pesa(0);
  }

  try {
    return fyo.pesa(value);
  } catch {
    return fyo.pesa(0);
  }
}

export function joinMapLists<A, B>(
  listA: A[],
  listB: B[],
  keyA: keyof A,
  keyB: keyof B
): (A & B)[] {
  const mapA = getMapFromList(listA, keyA);
  const mapB = getMapFromList(listB, keyB);

  const keyListA = listA
    .map((i) => i[keyA])
    .filter((k) => (k as unknown as string) in mapB);

  const keyListB = listB
    .map((i) => i[keyB])
    .filter((k) => (k as unknown as string) in mapA);

  const keys = new Set([keyListA, keyListB].flat().sort());

  const joint: (A & B)[] = [];
  for (const k of keys) {
    const a = mapA[k as unknown as string];
    const b = mapB[k as unknown as string];
    const c = { ...a, ...b };

    joint.push(c);
  }

  return joint;
}

export function removeAtIndex<T>(array: T[], index: number): T[] {
  if (index < 0 || index >= array.length) {
    return array;
  }

  return [...array.slice(0, index), ...array.slice(index + 1)];
}

/**
 * Asserts that `value` is of type T. Use with care.
 */
export const assertIsType = <T>(value: unknown): value is T => true;

import { t } from 'fyo';
import { Doc } from 'fyo/model/doc';
import { isPesa } from 'fyo/utils';
import {
  BaseError,
  DuplicateEntryError,
  LinkValidationError,
} from 'fyo/utils/errors';
import { Field, FieldType, FieldTypeEnum, NumberField } from 'schemas/types';
import { fyo } from 'src/initFyo';

export function stringifyCircular(
  obj: unknown,
  ignoreCircular = false,
  convertDocument = false
): string {
  const cacheKey: string[] = [];
  const cacheValue: unknown[] = [];

  return JSON.stringify(obj, (key: string, value: unknown) => {
    if (typeof value !== 'object' || value === null) {
      cacheKey.push(key);
      cacheValue.push(value);
      return value;
    }

    if (cacheValue.includes(value)) {
      const circularKey: string =
        cacheKey[cacheValue.indexOf(value)] || '{self}';
      return ignoreCircular ? undefined : `[Circular:${circularKey}]`;
    }

    cacheKey.push(key);
    cacheValue.push(value);

    if (convertDocument && value instanceof Doc) {
      return value.getValidDict();
    }

    return value;
  });
}

export function fuzzyMatch(input: string, target: string) {
  const keywordLetters = [...input];
  const candidateLetters = [...target];

  let keywordLetter = keywordLetters.shift();
  let candidateLetter = candidateLetters.shift();

  let isMatch = true;
  let distance = 0;

  while (keywordLetter && candidateLetter) {
    if (keywordLetter === candidateLetter) {
      keywordLetter = keywordLetters.shift();
    } else if (keywordLetter.toLowerCase() === candidateLetter.toLowerCase()) {
      keywordLetter = keywordLetters.shift();
      distance += 0.5;
    } else {
      distance += 1;
    }

    candidateLetter = candidateLetters.shift();
  }

  if (keywordLetter !== undefined) {
    distance = Number.MAX_SAFE_INTEGER;
    isMatch = false;
  } else {
    distance += candidateLetters.length;
  }

  return { isMatch, distance };
}

export function convertPesaValuesToFloat(obj: Record<string, unknown>) {
  Object.keys(obj).forEach((key) => {
    const value = obj[key];
    if (!isPesa(value)) {
      return;
    }

    obj[key] = value.float;
  });
}

export function getErrorMessage(e: Error, doc?: Doc): string {
  const errorMessage = e.message || t`An error occurred.`;

  let { schemaName, name } = doc ?? {};
  if (!doc) {
    schemaName = (e as BaseError).more?.schemaName as string | undefined;
    name = (e as BaseError).more?.value as string | undefined;
  }

  if (!schemaName || !name) {
    return errorMessage;
  }

  const label = fyo.db.schemaMap[schemaName]?.label ?? schemaName;
  if (e instanceof LinkValidationError) {
    return t`${label} ${name} is linked with existing records.`;
  } else if (e instanceof DuplicateEntryError) {
    return t`${label} ${name} already exists.`;
  }

  return errorMessage;
}

export function isNumeric(
  fieldtype: FieldType
): fieldtype is NumberField['fieldtype'];
export function isNumeric(fieldtype: Field): fieldtype is NumberField;
export function isNumeric(
  fieldtype: Field | FieldType
): fieldtype is NumberField | NumberField['fieldtype'] {
  if (typeof fieldtype !== 'string') {
    fieldtype = fieldtype?.fieldtype;
  }

  const numericTypes: FieldType[] = [
    FieldTypeEnum.Int,
    FieldTypeEnum.Float,
    FieldTypeEnum.Currency,
  ];

  return numericTypes.includes(fieldtype);
}

