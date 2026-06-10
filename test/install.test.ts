import { describe, expect, it } from 'vitest';
import { mkdtempSync, mkdirSync, existsSync, readFileSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import path from 'node:path';

import { ADAPTERS, claudeAdapter, geminiAdapter, cursorAdapter } from '../src/core/adapters/index.js';
import { resolveCommandRefs, parseGeneratedVersion } from '../src/core/adapters/types.js';
import { WORKFLOWS, workflowBody } from '../src/core/workflows.js';
import {
  detectTools,
  configuredTools,
  installShared,
  installTool,
  installedVersion,
  addToGitExclude,
} from '../src/core/installer.js';
import { packageVersion } from '../src/utils/package.js';

function tmpProject(): string {
  return mkdtempSync(path.join(tmpdir(), 'prism-test-'));
}

describe('workflows', () => {
  it('every workflow has a template body', () => {
    for (const w of WORKFLOWS) {
      expect(workflowBody(w.id).length, w.id).toBeGreaterThan(100);
    }
  });

  it('cross-references resolve for every adapter (no leftover placeholders)', () => {
    for (const a of ADAPTERS) {
      for (const w of WORKFLOWS) {
        const body = resolveCommandRefs(workflowBody(w.id), a);
        expect(body, `${a.id}/${w.id}`).not.toMatch(/\{\{cmd:/);
      }
    }
  });
});

describe('adapters', () => {
  it('claude uses namespaced commands, cursor uses flat ones', () => {
    expect(claudeAdapter.commandFilePath('drill')).toBe('.claude/commands/prism/drill.md');
    expect(claudeAdapter.commandRef('drill')).toBe('/prism:drill');
    expect(cursorAdapter.commandFilePath('drill')).toBe('.cursor/commands/prism-drill.md');
    expect(cursorAdapter.commandRef('drill')).toBe('/prism-drill');
  });

  it('gemini renders TOML with {{args}} instead of $ARGUMENTS', () => {
    const w = WORKFLOWS.find((x) => x.id === 'verify')!;
    const out = geminiAdapter.formatCommand(w, workflowBody('verify'), '0.0.0');
    expect(out).toContain('{{args}}');
    expect(out).not.toContain('$ARGUMENTS');
    expect(out).toMatch(/^description = /m);
  });

  it('generated files carry a parseable version stamp', () => {
    const w = WORKFLOWS[0];
    for (const a of ADAPTERS) {
      const out = a.formatCommand(w, 'body', '1.2.3');
      expect(parseGeneratedVersion(out), a.id).toBe('1.2.3');
    }
  });
});

describe('install', () => {
  it('installs commands for selected tools and detects them back', () => {
    const root = tmpProject();
    mkdirSync(path.join(root, '.claude'));
    mkdirSync(path.join(root, '.cursor'));

    expect(detectTools(root).map((t) => t.id)).toEqual(['claude', 'cursor']);
    expect(configuredTools(root)).toHaveLength(0);

    installShared(root);
    installTool(root, claudeAdapter);

    expect(existsSync(path.join(root, '.prism/conventions.md'))).toBe(true);
    for (const w of WORKFLOWS) {
      expect(existsSync(path.join(root, claudeAdapter.commandFilePath(w.id)))).toBe(true);
    }
    expect(configuredTools(root).map((t) => t.id)).toEqual(['claude']);
    expect(installedVersion(root, claudeAdapter)).toBe(packageVersion());

    const drill = readFileSync(path.join(root, '.claude/commands/prism/drill.md'), 'utf8');
    expect(drill).toContain('`/prism:decompose`');
    expect(drill).toContain('.prism/conventions.md');
  });

  it('adds .prism/ to .git/info/exclude once', () => {
    const root = tmpProject();
    mkdirSync(path.join(root, '.git/info'), { recursive: true });
    writeFileSync(path.join(root, '.git/info/exclude'), '# comment\n');

    addToGitExclude(root);
    addToGitExclude(root);

    const exclude = readFileSync(path.join(root, '.git/info/exclude'), 'utf8');
    expect(exclude.match(/\.prism\//g)).toHaveLength(1);
  });
});
