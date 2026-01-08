# Siberia Proxy with Antigravity Account Manager

This project is the Siberia Proxy with Antigravity Account Manager, built using the [BMad Method](https://github.com/bmadcode/bmad-method) within the Antigravity AI agent environment.

It comes pre-configured with agent definitions, tasks, and transposed workflows that allow Antigravity to fully leverage the BMad agile development process through native slash commands.

## Features

*   **Integrated BMad Core**: Contains the full `.bmad-core` definition (Agents, Tasks, Templates).
*   **Integrated BMad Core**: Contains the full `.bmad-core` definition (Agents, Tasks, Templates).
*   **Enterprise Proxy**: High-performance (Go/Rust) engine with Buffer Pooling and Privacy Masking.
*   **Traffic Inspector**: Virtualized UI capable of handling 5000+ events with deep body inspection.
*   **Universal IDE Support**: Native injection for VS Code, Cursor, and Windsurf.
*   **Internal AI**: Built-in Z.ai MCP Server with Web Search & Vision.


## Project Standards

### Frontend Stack
*   **Library**: [Creative Tim UI](https://github.com/creativetimofficial/ui) (shadcn/ui compatible).
*   **Integration**: Uses **Shadcn MCP Server** for component management.

### MCP Configuration
To enable the UI tools, add this to your `mcp_config.json`:

```json
{
  "mcpServers": {
    "shadcn-server": {
      "command": "npx",
      "args": [
        "-y",
        "@heilgar/shadcn-ui-mcp-server"
      ]
    }
  }
}
```

### Conductor & Parallel Development
This project is designed to support parallel development with **Conductor**.

*   **Worktrees**: Conductor runs agents in isolated Git Worktrees.
*   **Protocol**:
    *   **Antigravity** runs in the primary workspace.
    *   **DO NOT** check out the same branch in both places simultaneously.
    *   Always use feature branches (`antigravity/feat-xyz` vs `conductor/feat-abc`).
    *   Merge via Pull Requests to `main`.

## Available Command Agents

Use these slash commands in Antigravity to activate specific BMad roles:

| Command | Role | Description |
| :--- | :--- | :--- |
| `/pm` | **Product Manager** | Requirements gathering, PRD creation, user stories. |
| `/architect` | **System Architect** | Technical design, architecture documentation. |
| `/po` | **Product Owner** | Backlog management, story validation, process alignment. |
| `/dev` | **Developer** | Implementation of stories, writing code. |
| `/qa` | **Test Architect** | Testing strategy, risk assessment, quality gates. |
| `/ux` | **UX Expert** | UI/UX design, frontend specifications. |
| `/analyst` | **Analyst** | Deep research, market analysis. |
| `/sm` | **Scrum Master** | Process checking, alignment verification. |
| `/master` | **BMad Master** | General-purpose helper for the framework. |
| `/orchestrator` | **Orchestrator** | High-level agent coordination. |

## Quick Start

1.  **Start Planning**:
    *   Run `/pm` to begin defining your project requirements.
    *   The agent will guide you through creating a `docs/prd.md`.

2.  **Follow the Flow**:
    *   Agents will read your config and artifacts automatically.
    *   Upon finishing a task, they will suggest the next command (e.g., `Now run /architect`).


## Maintenance

The repository includes a helper script `.gemini/transpose_bmad.py`. If you update the core BMad definitions in `.bmad-core`, run this script to regenerate the Antigravity workflows:

```bash
python3 .gemini/transpose_bmad.py
```
