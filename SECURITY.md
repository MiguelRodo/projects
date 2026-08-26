# Security Policy

## Supported Versions

We provide security updates and patches for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| `v0.1.x`| :white_check_mark: |
| `< 0.1` | :x:                |

---

## Reporting a Vulnerability

If you discover a potential security vulnerability in `projects` or `projectctl`, please report it responsibly:

1. **Do not create a public GitHub issue** for undisclosed security vulnerabilities.
2. Please open a [GitHub Private Vulnerability Report](https://github.com/MiguelRodo/projects/security/advisories/new) or send an email to the project maintainers.
3. Include a detailed description of the issue, reproduction steps, affected versions, and any suggested fixes or mitigations.
4. We will acknowledge receipt of your report within 48 hours and provide a timeline for triage and remediation.

---

## Security & Privacy Boundary Guarantees

`projects` is engineered with explicit privacy boundary guarantees:

1. **No Secret Leaks**: Configuration manifests and workspace files must never require or store plaintext tokens, passwords, or credentials.
2. **Neutrality & Privacy Separation**:
   - Public repository manifests, examples, and fixtures must not contain private internal repository names, personal identifiers, or private cloud storage identifiers (e.g. Google Drive/Sheets IDs).
   - Sensitive companion data or private control-plane records must remain exclusively on user-controlled private storage and reference only non-sensitive stable linkage identifiers.
3. **Least-Privilege GitHub Access**: Tools and workflows must operate with the minimum required GitHub token scopes.
