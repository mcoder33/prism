import type { Workflow } from '../workflows.js';

export interface ToolAdapter {
  /** Stable id used in --tools, e.g. "claude". */
  id: string;
  /** Display name, e.g. "Claude Code". */
  name: string;
  /** Directories/files whose presence means the tool is used in the project. */
  detectPaths: string[];
  /** Project-relative path of the generated command file for a workflow. */
  commandFilePath(workflowId: string): string;
  /** How the user invokes a workflow in this tool, e.g. "/prism:drill". */
  commandRef(workflowId: string): string;
  /** Render the full command file (frontmatter + body) for this tool. */
  formatCommand(workflow: Workflow, body: string, version: string): string;
}

/** Replace {{cmd:<id>}} placeholders with the tool-specific slash command. */
export function resolveCommandRefs(body: string, adapter: ToolAdapter): string {
  return body.replace(/\{\{cmd:([a-z-]+)\}\}/g, (_, id: string) =>
    `\`${adapter.commandRef(id)}\``,
  );
}

export function generatedStamp(version: string): string {
  return `<!-- prism:generated v${version} — managed by the prism CLI, do not edit (run \`prism update\` to regenerate) -->`;
}

/** Extract the version from a previously generated file, if present. */
export function parseGeneratedVersion(content: string): string | null {
  const m = content.match(/prism:generated v([0-9][^\s]*)/);
  return m ? m[1] : null;
}
