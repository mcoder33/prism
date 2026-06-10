import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

// Package root = two levels up from this file (src/utils/ or dist/utils/).
export const packageRoot = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  '..',
  '..',
);

export const templatesDir = path.join(packageRoot, 'templates');

let cachedVersion: string | undefined;

export function packageVersion(): string {
  if (!cachedVersion) {
    const pkg = JSON.parse(
      readFileSync(path.join(packageRoot, 'package.json'), 'utf8'),
    ) as { version: string };
    cachedVersion = pkg.version;
  }
  return cachedVersion;
}

export function loadTemplate(relativePath: string): string {
  return readFileSync(path.join(templatesDir, relativePath), 'utf8');
}
