import path from 'node:path';
import pc from 'picocolors';

import {
  configuredTools,
  detectTools,
  installShared,
  installTool,
  installedVersion,
} from '../installer.js';
import { packageVersion } from '../../utils/package.js';

export interface UpdateOptions {
  force?: boolean;
}

export async function runUpdate(targetPath: string | undefined, options: UpdateOptions): Promise<void> {
  const projectRoot = path.resolve(targetPath ?? process.cwd());
  const version = packageVersion();

  const configured = configuredTools(projectRoot);
  if (configured.length === 0) {
    console.log(pc.yellow('No prism commands found in this project. Run `prism init` first.'));
    return;
  }

  const stale = options.force
    ? configured
    : configured.filter((t) => installedVersion(projectRoot, t) !== version);

  if (stale.length === 0) {
    console.log(pc.green(`All tools are up to date (v${version}).`) + pc.dim(' Use --force to regenerate anyway.'));
  } else {
    installShared(projectRoot);
    for (const tool of stale) {
      const from = installedVersion(projectRoot, tool) ?? 'unknown';
      installTool(projectRoot, tool);
      console.log(`  ${pc.green('✔')} ${tool.name}: ${from} → v${version}`);
    }
  }

  const newTools = detectTools(projectRoot).filter(
    (t) => !configured.some((c) => c.id === t.id),
  );
  if (newTools.length > 0) {
    console.log(
      pc.dim(`Detected but not configured: ${newTools.map((t) => t.name).join(', ')} — run \`prism init\` to add them.`),
    );
  }
}
