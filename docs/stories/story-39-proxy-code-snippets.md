# Story-39: Proxy Code Snippets (DevTools)

**Epic:** [Epic-02: Siberia Proxy Core](./epic-02-proxy.md)
**Status:** Ready
**Reference:** `docs/ag-ref-docs/feature-frontend.md`

## Goal
Provide a "Code Generator" in the Proxy UI that creates ready-to-use Client snippets (Python/OpenAI, Curl, Node.js) configured to point to the local Siberia proxy.

## Context
The Reference App includes a helper that generates code like:
```python
client = OpenAI(
    base_url="http://localhost:3000/v1",
    api_key="dummy"
)
```
This reduces friction for developers trying to test the proxy with their scripts.

## Tasks
- [ ] **Task 1: Snippet Generator Utlity**
    -   Create a frontend utility `generateSnippet(language, port, model)`.
    -   Support: Curl, Python (openai lib), Node.js (openai lib), .env format.
- [ ] **Task 2: UI Integration**
    -   Add "Use this Proxy" button/modal to the Proxy Dashboard.
    -   Display the snippets with "Copy to Clipboard" functionality.
- [ ] **Task 3: Dynamic Port**
    -   Ensure the snippets reflect the *actual* running port of the proxy.

## Acceptance Criteria
- [ ] **Accuracy:** generated Python code runs without modification (assuming `openai` pip package is installed).
- [ ] **Usability:** User can copy the code with one click.

## Technical Notes
- Frontend-only story.
