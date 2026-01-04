# BMad Workflow & Agent Guide

This document defines the standard operating procedures, workflows, and "Rules of Engagement" for AI Agents (and humans) working within this project.

## 🤖 Agent Roles & Responsibilities

| Agent | Command | Role | Primary Output |
| :--- | :--- | :--- | :--- |
| **Product Owner** | `/po` | Requirements, Backlog, Priorities | Story Files, Epics |
| **Architect** | `/architect` | System Design, Tech Specs | `architecture.md`, Design Docs |
| **Developer** | `/dev` | Implementation, Testing | Code, Tests, Status Updates |
| **QA** | `/qa` | Gating, Risk Assessment, Testing | Gate Files, QA Results |

---

## 🔄 The Standard Development Loop

The lifecycle of a unit of work (Story) follows this strict chain of custody:

```mermaid
graph LR
    PO((PO)) -->|1. Draft & Validate| Arch{Architect?}
    Arch -->|2. Design (Complex)| Dev((Dev))
    Arch -->|2. Direct (Simple)| Dev
    Dev -->|3. Implement & Test| QA((QA))
    QA -->|4. Gate: FAIL| Dev
    QA -->|4. Gate: PASS| PO
    PO -->|5. Acceptance| Done((Done))
```

### Phase 1: Definition (`/po`)
1.  **Draft**: User or PO creates a story.
2.  **Validate**: PO runs `*validate-story-draft`.
3.  **Route**:
    *   If design is unclear or complex: **Handoff to `/architect`**.
    *   If clear and simple: **Handoff to `/dev`**.

### Phase 2: Design (`/architect`)
*   *Optional step for complex features.*
1.  Architect reviews requirements.
2.  Creates/Updates `architecture.md` or specific design specs.
3.  **Handoff to `/dev`** with link to design docs.

### Phase 3: Implementation (`/dev`)
1.  **Check Design**: run `*check-design` to align with Architect's intent.
2.  **Develop**: run `*develop-story` (Test-Driven Development).
3.  **Golden Rules**:
    *   ❌ **NEVER** commit untested code.
    *   ❌ **NEVER** hardcode credentials.
    *   ❌ **NEVER** leave debug code.
4.  **Completion**: run `*completion`.
5.  **Handoff to `/qa`**.

### Phase 4: Quality Gate (`/qa`)
1.  **Review**: QA runs `*review {story}`.
2.  **Gate Decision**:
    *   **FAIL / CONCERNS**: QA provides Gate File. **Handoff back to `/dev`** (Dev runs `*review-qa`).
    *   **PASS**: **Handoff to `/po`**.

### Phase 5: Acceptance (`/po`)
1.  PO reviews the "Done" story and QA signals.
2.  Mark as **Closed / Released**.

---

## 📝 Handoff Protocols

When finishing a task, Agents **MUST** provide a "Handoff Receipt" to ensure context isn't lost.

**Format:**
```markdown
# Handoff Checklist
- Modified/Created: [List files]
- Status: [Current Status]
- Key Decisions: [One-line summary]
- Next Recommended Action: [Command for next agent, e.g., "Run /qa to gate"]
```

---

## 🔑 Command Cheat Sheet

| I want to... | Agent | Command |
| :--- | :--- | :--- |
| Create a new feature | `sarah` / `po` | `*create-story` |
| Design a system | `winston` / `architect` | `*create-backend-architecture` |
| Write code | `james` / `dev` | `*develop-story` |
| Fix a bug | `james` / `dev` | `*hotfix` (manual workflow) |
| Check quality | `quinn` / `qa` | `*review {story}` |
| Fix QA issues | `james` / `dev` | `*review-qa` |

---

## 🚫 Critical Constraints (All Agents)

1.  **Context Hygiene**: Do not hallucinate file content. Read it (`view_file`) before acting.
2.  **File Boundaries**:
    *   **Dev** writes Code + Story Status/Logs.
    *   **QA** writes Gate Files + Story QA Results.
    *   **PO** writes Story Requirements + Priorities.
    *   **Architect** writes Design Docs.
3.  **Sequential Execution**: Do not skip the pipeline. Devs should not merge without QA. QA should not gate without Dev completion.
