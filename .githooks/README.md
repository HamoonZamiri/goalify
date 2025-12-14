# Git Hooks Setup

## First-time Setup

Run this command once to enable the hooks:

```bash
git config core.hooksPath .githooks
```

## What Runs on Pre-commit

The pre-commit hook runs comprehensive checks before every commit (~30-45s total):

**Frontend** (~10-15s):
- Format (auto-fix)
- Lint (auto-fix)
- Type checking
- Unit tests

**Backend** (~20-30s):
- Format (auto-fix)
- Lint
- All tests (unit + store + integration)
- **Note**: Requires Docker running for integration tests

## Bypassing Hooks

In emergencies only:
```bash
git commit --no-verify
```
