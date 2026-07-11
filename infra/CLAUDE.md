# infra/ (Terraform + deploy)

terraform/modules are reusable; terraform/envs/{lean,enterprise} compose them. Lean is the default target.

- prevent_destroy on all stateful resources (RDS, SQS). Never run apply with --auto-approve.
- plan is posted to the PR; apply runs only after manual approval. Do not automate the approval away.
- deploy/ holds the nginx blue/green template, the systemd unit, and swap.sh. The swap is rewrite-then-reload.
- Change flow: edit module or env, run terraform plan, show me the plan, wait. Never plan and apply in one step.