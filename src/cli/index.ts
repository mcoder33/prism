import { Command } from 'commander';
import pc from 'picocolors';

import { packageVersion } from '../utils/package.js';
import { runInit } from '../core/commands/init.js';
import { runUpdate } from '../core/commands/update.js';
import { runList } from '../core/commands/list.js';
import { ADAPTERS } from '../core/adapters/index.js';

export function runCli(argv: string[] = process.argv): void {
  const program = new Command();

  program
    .name('prism')
    .description(
      'PRISM — recursive decomposition workflow for AI coding agents.\nInstalls /prism slash commands into a project for the agents you use.',
    )
    .version(packageVersion());

  program
    .command('init')
    .argument('[path]', 'project directory (default: current directory)')
    .option(
      '--tools <list>',
      `comma-separated tool ids (${ADAPTERS.map((a) => a.id).join(', ')}), or "all"/"none"; omit for interactive selection`,
    )
    .description('install prism slash commands into a project')
    .action(async (targetPath: string | undefined, opts: { tools?: string }) => {
      await runInit(targetPath, opts);
    });

  program
    .command('update')
    .argument('[path]', 'project directory (default: current directory)')
    .option('--force', 'regenerate even if versions match')
    .description('refresh previously installed prism command files')
    .action(async (targetPath: string | undefined, opts: { force?: boolean }) => {
      await runUpdate(targetPath, opts);
    });

  program
    .command('list')
    .argument('[path]', 'project directory (default: current directory)')
    .description('list active prism changes in a project')
    .action(async (targetPath: string | undefined) => {
      await runList(targetPath);
    });

  program.parseAsync(argv).catch((err: unknown) => {
    const message = err instanceof Error ? err.message : String(err);
    console.error(pc.red(`error: ${message}`));
    process.exitCode = 1;
  });
}
