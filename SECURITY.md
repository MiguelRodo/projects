# Security policy

## Supported versions

The project is still pre-v1. Only the latest published CLI release and the
current repository skill are supported while the interface settles.

## Reporting a vulnerability

Do not open a public issue for an undisclosed vulnerability.

Use the repository's private vulnerability reporting facility where available, or contact the maintainers privately. Include the affected revision, reproduction steps, impact and any known mitigation.

## Security and privacy boundaries

- Shared contracts, documentation and fixtures must contain only collaborator-safe information.
- Credentials, tokens, passwords, private keys and equivalent authentication material must never be committed.
- Private source locations, private destinations, local workspace mappings and per-user provider configuration belong outside shared repositories.
- Tests must not mutate live GitHub resources or private systems.
- Provider access must use least privilege.
- Mutating workflows must plan first, refuse stale writes and verify completed writes by readback.
