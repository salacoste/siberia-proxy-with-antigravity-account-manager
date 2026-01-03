# Git Branching & Epic Workflow Protocol

This document defines the regulations for managing branches, epics, and merges within the Siberia Proxy project.

## 1. Branching Strategy

We use a feature-branch workflow centered around **Epics**.

### 1.1 Branch Naming Conventions
| Type | Pattern | Created By | Description |
|------|---------|------------|-------------|
| **Epic** | `feat/epic-XX-<slug>` | **Dev** (First Story) | Long-lived branch for a collection of stories. Base for all story branches in that epic. |
| **Story** | `antigravity/epic-XX/story-YY-<slug>` | **Dev** | Short-lived branch for a single story. Merges into the **Epic Branch**. |
| **Hotfix** | `fix/<issue-id>-<slug>` | **Dev** | Critical fixes for production. Merges to `main` & `develop`. |
| **Release** | `release/vX.Y.Z` | **Orchestrator** | Stabilization branch before main merge. |

### 1.2 Responsibilities

#### Product Owner (PO)
*   **Does NOT** create Git branches.
*   **Defines** the Epic ID (e.g., `Epic-09`) and Slug (e.g., `advanced-inspection`) in documentation.
*   *Why?* Keeps PO focused on "What" (Documentation), not "How" (Git).

#### Developer (Dev)
*   **Creates** the Epic Branch (`feat/epic-XX...`) when implementing the *first* story of that Epic.
*   **Checks** for existing Epic branch before starting subsequent stories.
*   **Merges** Story branches into the Epic Branch upon Story verification.

#### QA / Architect
*   **Validates** the entire Epic Branch when all stories are complete.

---

## 2. Epic Lifecycle & Merge Regulations

### 2.1 The Merge Gate
An Epic Branch (`feat/epic-XX`) can ONLY be merged into `develop` (or `main`) when:
1.  **All Stories Completed**: all associated `story-YY` files in `docs/stories` are marked `[x] Status: Completed`.
2.  **QA Validation**: A "Phase Validation Report" (e.g., `docs/qa/phase-X-validation.md`) confirms the feature set is stable.
3.  **No Conflicts**: The branch is up-to-date with `develop`.

*Action*: The **Tech Lead** or **Orchestrator** performs the merge of the Epic Branch.

### 2.2 Re-opening & Extending Epics
**Rule: Epics are Finite.**

Once an Epic is merged and "Closed", it **CANNOT** be reopened.

*   **Scenario A: Bug in Merged Epic**
    *   Action: Treat as a standard **Bug**. Create `fix/...` branch. Do NOT reopen the Epic branch.
    
*   **Scenario B: New Feature for Merged Epic** (e.g., "Add more filters to Inspection")
    *   **Do NOT** re-use `Epic-09`.
    *   **Action**: Create a **New Epic** (e.g., `Epic-15: Enhanced Inspection`).
    *   *Reason*: Preserves history, allows independent tracking, prevents "Zombie Epics" that never die.

---

## 3. Automation & Enforcement
*   **Agents (Dev)**: Must query `git branch -r` to check for specific `feat/epic-XX` before creating.
*   **CI/CD**: Workflows should trigger on pushes to `feat/epic-*` to run full regression suites.
