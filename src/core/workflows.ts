import { loadTemplate } from '../utils/package.js';

export interface Workflow {
  /** Stable command id; also the file/slash-command suffix. */
  id: string;
  /** Human title used in frontmatter, e.g. "PRISM: Propose". */
  title: string;
  /** One-line description shown by agents in command pickers. */
  description: string;
}

export const WORKFLOWS: Workflow[] = [
  {
    id: 'use',
    title: 'PRISM: Use',
    description:
      'Select the active change (like git checkout) via an interactive picker — switch, stop, or "+ New change" (→ propose). All prism commands then default to it.',
  },
  {
    id: 'propose',
    title: 'PRISM: Propose',
    description:
      'Grill on requirements, survey best practices, pick a strategy + data-flow, then write the seed (proposal + concept) for a new decomposition change.',
  },
  {
    id: 'decompose',
    title: 'PRISM: Decompose',
    description:
      'Split the proposal (or a node) into a few small digestible node.md parts. Recursive.',
  },
  {
    id: 'drill',
    title: 'PRISM: Drill',
    description:
      'Drill ONE part to atomic and generate its artifact set (spec, detail, concept.drawio, signatures, tasks).',
  },
  {
    id: 'integrate',
    title: 'PRISM: Integrate',
    description:
      'Produce the cross-part artifacts — integration.drawio + combined signatures.md + overall tasks.md.',
  },
  {
    id: 'apply',
    title: 'PRISM: Apply',
    description:
      'Implement the change in code per the tasks, in dependency order, marking tasks done and running checks.',
  },
  {
    id: 'verify',
    title: 'PRISM: Verify',
    description:
      'Thorough post-implementation verification on a running dev environment — full test suite, blocking static checks, diff-driven functional and browser smoke, targeted concurrency/parallelism checks, load corner-cases, ability to fix findings and re-verify, final report with recommendations. Project-agnostic: commands and entry points are detected from repository configs.',
  },
  {
    id: 'archive',
    title: 'PRISM: Archive',
    description:
      'Archive a completed change — move .prism/<change>/ to .prism/archive/<change>/.',
  },
];

export function workflowBody(id: string): string {
  return loadTemplate(`commands/${id}.md`).trimEnd();
}

export function conventionsBody(): string {
  return loadTemplate('conventions.md').trimEnd();
}
