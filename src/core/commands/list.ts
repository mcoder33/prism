import path from 'node:path';
import { existsSync, readdirSync, readFileSync } from 'node:fs';
import pc from 'picocolors';

import { PRISM_DIR } from '../installer.js';

export async function runList(targetPath: string | undefined): Promise<void> {
  const projectRoot = path.resolve(targetPath ?? process.cwd());
  const prismDir = path.join(projectRoot, PRISM_DIR);
  if (!existsSync(prismDir)) {
    console.log(pc.yellow('No .prism/ directory here. Start a change with the propose command in your agent.'));
    return;
  }

  const currentFile = path.join(prismDir, 'CURRENT');
  const current = existsSync(currentFile) ? readFileSync(currentFile, 'utf8').trim() : null;

  const changes = readdirSync(prismDir, { withFileTypes: true })
    .filter((e) => e.isDirectory() && e.name !== 'archive')
    .map((e) => e.name)
    .sort();

  if (changes.length === 0) {
    console.log(pc.dim('No active changes.'));
  } else {
    console.log(pc.bold('Active changes:'));
    for (const name of changes) {
      const marker = name === current ? pc.green(' (current)') : '';
      console.log(`  ${name}${marker}`);
    }
  }

  const archiveDir = path.join(prismDir, 'archive');
  if (existsSync(archiveDir)) {
    const archived = readdirSync(archiveDir, { withFileTypes: true }).filter((e) => e.isDirectory());
    if (archived.length > 0) {
      console.log(pc.dim(`Archived: ${archived.length} (.prism/archive/)`));
    }
  }
}
