import { existsSync, mkdirSync, readFileSync, writeFileSync, appendFileSync } from 'node:fs';
import path from 'node:path';

import { ADAPTERS } from './adapters/index.js';
import { resolveCommandRefs, generatedStamp, parseGeneratedVersion, type ToolAdapter } from './adapters/types.js';
import { WORKFLOWS, workflowBody, conventionsBody } from './workflows.js';
import { packageVersion } from '../utils/package.js';

export const PRISM_DIR = '.prism';
export const CONVENTIONS_PATH = `${PRISM_DIR}/conventions.md`;

export interface InstallResult {
  tool: ToolAdapter;
  files: string[];
}

/** Tools whose dot-dirs/config files are present in the project. */
export function detectTools(projectRoot: string): ToolAdapter[] {
  return ADAPTERS.filter((a) =>
    a.detectPaths.some((p) => existsSync(path.join(projectRoot, p))),
  );
}

/** Tools that already have prism command files installed. */
export function configuredTools(projectRoot: string): ToolAdapter[] {
  return ADAPTERS.filter((a) =>
    WORKFLOWS.some((w) => existsSync(path.join(projectRoot, a.commandFilePath(w.id)))),
  );
}

/** Installed prism version per configured tool (from the generated stamp). */
export function installedVersion(projectRoot: string, tool: ToolAdapter): string | null {
  for (const w of WORKFLOWS) {
    const file = path.join(projectRoot, tool.commandFilePath(w.id));
    if (existsSync(file)) {
      return parseGeneratedVersion(readFileSync(file, 'utf8'));
    }
  }
  return null;
}

function writeFileEnsured(filePath: string, content: string): void {
  mkdirSync(path.dirname(filePath), { recursive: true });
  writeFileSync(filePath, content, 'utf8');
}

/** Write .prism/conventions.md (shared by all tools) and exclude .prism/ from git. */
export function installShared(projectRoot: string): string[] {
  const version = packageVersion();
  const conventionsFile = path.join(projectRoot, CONVENTIONS_PATH);
  writeFileEnsured(
    conventionsFile,
    `${generatedStamp(version)}\n\n${conventionsBody()}\n`,
  );
  addToGitExclude(projectRoot);
  return [CONVENTIONS_PATH];
}

/** Write all command files for one tool. Files are tool-owned: always overwritten. */
export function installTool(projectRoot: string, tool: ToolAdapter): InstallResult {
  const version = packageVersion();
  const files: string[] = [];
  for (const w of WORKFLOWS) {
    const body = resolveCommandRefs(workflowBody(w.id), tool);
    const rendered = tool.formatCommand(w, body, version);
    const rel = tool.commandFilePath(w.id);
    writeFileEnsured(path.join(projectRoot, rel), rendered);
    files.push(rel);
  }
  return { tool, files };
}

/** Add .prism/ to .git/info/exclude so artifacts are never committed. */
export function addToGitExclude(projectRoot: string): void {
  const gitDir = path.join(projectRoot, '.git');
  if (!existsSync(gitDir)) return;
  const excludeFile = path.join(gitDir, 'info', 'exclude');
  const entry = '.prism/';
  let current = '';
  if (existsSync(excludeFile)) {
    current = readFileSync(excludeFile, 'utf8');
    if (current.split('\n').some((l) => l.trim() === entry || l.trim() === '.prism')) {
      return;
    }
  } else {
    mkdirSync(path.dirname(excludeFile), { recursive: true });
  }
  const prefix = current.length > 0 && !current.endsWith('\n') ? '\n' : '';
  appendFileSync(excludeFile, `${prefix}${entry}\n`);
}
