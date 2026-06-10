import type { Workflow } from '../workflows.js';
import { generatedStamp, type ToolAdapter } from './types.js';

function yamlQuote(value: string): string {
  return `"${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`;
}

/** Markdown command file with YAML frontmatter (Claude Code style). */
function markdownWithFrontmatter(
  workflow: Workflow,
  body: string,
  version: string,
  extraFrontmatter: string[] = [],
): string {
  const lines = [
    '---',
    `name: ${yamlQuote(workflow.title)}`,
    `description: ${yamlQuote(workflow.description)}`,
    ...extraFrontmatter,
    '---',
    generatedStamp(version),
    '',
    body,
    '',
  ];
  return lines.join('\n');
}

/** Plain markdown command file (tools without frontmatter support). */
function plainMarkdown(workflow: Workflow, body: string, version: string): string {
  return [
    `# ${workflow.title}`,
    '',
    `> ${workflow.description}`,
    '',
    generatedStamp(version),
    '',
    body,
    '',
  ].join('\n');
}

export const claudeAdapter: ToolAdapter = {
  id: 'claude',
  name: 'Claude Code',
  detectPaths: ['.claude'],
  commandFilePath: (id) => `.claude/commands/prism/${id}.md`,
  commandRef: (id) => `/prism:${id}`,
  formatCommand: (w, body, version) =>
    markdownWithFrontmatter(w, body, version, [
      'category: Workflow',
      'tags: [workflow, design, prism]',
    ]),
};

export const cursorAdapter: ToolAdapter = {
  id: 'cursor',
  name: 'Cursor',
  detectPaths: ['.cursor'],
  commandFilePath: (id) => `.cursor/commands/prism-${id}.md`,
  commandRef: (id) => `/prism-${id}`,
  formatCommand: (w, body, version) => plainMarkdown(w, body, version),
};

export const codexAdapter: ToolAdapter = {
  id: 'codex',
  name: 'Codex CLI',
  detectPaths: ['.codex'],
  commandFilePath: (id) => `.codex/prompts/prism-${id}.md`,
  commandRef: (id) => `/prism-${id}`,
  formatCommand: (w, body, version) => plainMarkdown(w, body, version),
};

export const geminiAdapter: ToolAdapter = {
  id: 'gemini',
  name: 'Gemini CLI',
  detectPaths: ['.gemini'],
  commandFilePath: (id) => `.gemini/commands/prism/${id}.toml`,
  commandRef: (id) => `/prism:${id}`,
  formatCommand: (w, body, version) => {
    // TOML literal multi-line strings process no escapes; guard the delimiter.
    // Gemini CLI injects invocation args via {{args}}, not $ARGUMENTS.
    const safeBody = body.replaceAll("'''", "''​'").replaceAll('$ARGUMENTS', '{{args}}');
    return [
      `# prism:generated v${version} — managed by the prism CLI, do not edit (run \`prism update\` to regenerate)`,
      `description = ${JSON.stringify(w.description)}`,
      "prompt = '''",
      `# ${w.title}`,
      '',
      safeBody,
      "'''",
      '',
    ].join('\n');
  },
};

export const copilotAdapter: ToolAdapter = {
  id: 'copilot',
  name: 'GitHub Copilot',
  detectPaths: ['.github/copilot-instructions.md', '.github/prompts'],
  commandFilePath: (id) => `.github/prompts/prism-${id}.prompt.md`,
  commandRef: (id) => `/prism-${id}`,
  formatCommand: (w, body, version) => markdownWithFrontmatter(w, body, version),
};

export const windsurfAdapter: ToolAdapter = {
  id: 'windsurf',
  name: 'Windsurf',
  detectPaths: ['.windsurf'],
  commandFilePath: (id) => `.windsurf/workflows/prism-${id}.md`,
  commandRef: (id) => `/prism-${id}`,
  formatCommand: (w, body, version) => markdownWithFrontmatter(w, body, version),
};

export const opencodeAdapter: ToolAdapter = {
  id: 'opencode',
  name: 'OpenCode',
  detectPaths: ['.opencode', 'opencode.json'],
  commandFilePath: (id) => `.opencode/command/prism-${id}.md`,
  commandRef: (id) => `/prism-${id}`,
  formatCommand: (w, body, version) => markdownWithFrontmatter(w, body, version),
};

export const ADAPTERS: ToolAdapter[] = [
  claudeAdapter,
  cursorAdapter,
  codexAdapter,
  geminiAdapter,
  copilotAdapter,
  windsurfAdapter,
  opencodeAdapter,
];

export function adapterById(id: string): ToolAdapter | undefined {
  return ADAPTERS.find((a) => a.id === id);
}
