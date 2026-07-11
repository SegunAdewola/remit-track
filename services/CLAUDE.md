# services/ (Python multi-agent suite)

Runs only inside the GitHub Actions sandbox. Agents: agent_healer, compliance_sentinel, onboarding_agent.

- Agents open PRs for human review. They never deploy or merge on their own.
- agent_healer/linter.go is the Go AST allowlist gate, shelled out from Python. Do not port it to Python ast.
- Shared LLM/S3/Chroma/PR helpers live in services/shared. Secrets come from env, never code.
- Tests: pytest, mock all network and LLM calls. Ask before adding deps; pin versions in requirements.txt.