import path from 'node:path';
import { existsSync } from 'node:fs';
import pc from 'picocolors';

import { ADAPTERS, adapterById } from '../adapters/index.js';
import type { ToolAdapter } from '../adapters/types.js';
import {
  configuredTools,
  detectTools,
  installShared,
  installTool,
} from '../installer.js';
import { packageVersion } from '../../utils/package.js';

export interface InitOptions {
  tools?: string;
}

function parseToolsFlag(flag: string): ToolAdapter[] {
  if (flag === 'all') return ADAPTERS;
  if (flag === 'none') return [];
  return flag.split(',').map((raw) => {
    const id = raw.trim().toLowerCase();
    const adapter = adapterById(id);
    if (!adapter) {
      const known = ADAPTERS.map((a) => a.id).join(', ');
      throw new Error(`Unknown tool "${id}". Known tools: ${known}, or "all"/"none".`);
    }
    return adapter;
  });
}

async function selectToolsInteractive(projectRoot: string): Promise<ToolAdapter[]> {
  const { checkbox } = await import('@inquirer/prompts');
  const detected = new Set(detectTools(projectRoot).map((a) => a.id));
  const configured = new Set(configuredTools(projectRoot).map((a) => a.id));
  const choices = ADAPTERS.map((a) => ({
    name:
      a.name +
      (configured.has(a.id) ? pc.dim(' (installed)') : detected.has(a.id) ? pc.dim(' (detected)') : ''),
    value: a.id,
    checked: configured.has(a.id) || detected.has(a.id),
  }));
  const picked = await checkbox({
    message: 'Which AI tools should get the prism commands?',
    choices,
  });
  return picked.map((id) => adapterById(id)!);
}

export async function runInit(targetPath: string | undefined, options: InitOptions): Promise<void> {
  const projectRoot = path.resolve(targetPath ?? process.cwd());
  if (!existsSync(projectRoot)) {
    throw new Error(`Path does not exist: ${projectRoot}`);
  }

  let tools: ToolAdapter[];
  if (options.tools) {
    tools = parseToolsFlag(options.tools);
  } else if (process.stdout.isTTY && process.stdin.isTTY) {
    tools = await selectToolsInteractive(projectRoot);
  } else {
    tools = detectTools(projectRoot);
    if (tools.length === 0) {
      throw new Error(
        'No AI tools detected and not running interactively. Pass --tools <list|all>.',
      );
    }
  }

  if (tools.length === 0) {
    console.log(pc.yellow('No tools selected — nothing to do.'));
    return;
  }

  const sharedFiles = installShared(projectRoot);
  const results = tools.map((t) => installTool(projectRoot, t));

  console.log();
  console.log(pc.bold(`prism v${packageVersion()} installed into ${projectRoot}`));
  console.log(pc.dim(`  shared: ${sharedFiles.join(', ')} (.prism/ is git-excluded)`));
  for (const r of results) {
    console.log(`  ${pc.green('✔')} ${r.tool.name}: ${r.files.length} commands → ${path.dirname(r.files[0])}/`);
  }
  console.log();
  console.log(`Try ${pc.cyan(tools[0].commandRef('propose'))} in ${tools[0].name} to start a change.`);
  console.log(pc.dim('Restart your IDE/agent if slash commands do not show up.'));
}
