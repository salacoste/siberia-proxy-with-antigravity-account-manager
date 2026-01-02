# Coding Standards & Workflow Protocol

**CRITICAL**: All agents MUST adhere to these rules. Deviations will cause conflicts with parallel workstreams.

## 1. Parallel Workflow Protocol (Conductor + Antigravity)

We operate in a parallel environment with **Conductor** agents.

### Branching Strategy
- **NEVER** work directly on `main`, `master`, or `develop`. These are protected synchronization points.
- **ALWAYS** create a feature branch before writing code:
    - Naming convention: `antigravity/<type>/<description>` (e.g., `antigravity/feat/user-login`).
    - Distinguish from Conductor branches: Conductor uses `conductor/...`.

### Worktree Safety
- **DO NOT** attempt to checkout a branch that is currently locked/active in a Conductor worktree.
- **Sync**: Run `git fetch origin` frequently to update your view of the remote state.

### Merging
- **DO NOT** merge locally. Push your branch and open a Pull Request.

## 2. Frontend & UI Standards

### UI Framework
- **Mandatory Library**: [Creative Tim UI](https://github.com/creativetimofficial/ui) (shadcn/ui compatible).
- **Aesthetics**: "Saint Vandal" style - vibrant, premium, glassmorphism, dynamic animations.

### Component Management
- **Tool**: You **MUST** use the **Shadcn MCP Server** (`shadcn-server`) to install components.
- **Registry**: `https://creative-tim.com/ui/r/all.json`.
- **Prohibition**: Do not manually copy-paste shadcn code if a Creative Tim component exists. Use the tool.
