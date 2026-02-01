# Security Policy

StepKeys is a locally running daemon that cannot be accessed remotely, so it does not present an attack surface. While it poses virtually no risk from external threats, like any software it could contain bugs that may affect system stability or behavior. Because StepKeys interacts closely with the operating system and your browser, maintaining its security and reliability is important.

## Supported Versions

Older versions of StepKeys will generally become unsupported when a new (non-patch) release is published. Support is typically provided for the latest minor version. In the event of a breaking change, the previous minor version might also receive temporary support limited to security updates.

| Version |     Supported      | End of Life |
| :-----: | :----------------: | :---------: |
|  1.0.0  | :white_check_mark: |  in 1.1.0   |

## Reporting a Vulnerability

If you discover a vulnerability, please report it as a regular bug. For severe vulnerabilities, or ones that could potentially be exploited, **do not** open a public issue. Instead, contact me directly via [email](mailto:barni@sigsegv.hu) or draft a [security advisory](https://github.com/BrNi05/StepKeys/security/advisories).

## What not to report

- Vulnerabilities affecting **unsupported versions** of StepKeys

- Dependency vulnerabilities — Dependabot already checks dependencies to ensure secure builds.

## Response and solution

Security issue related tickets are top priority and I will do my best to find a solution as fast as possible. Expect a response within 24 hours and a resolution within a few days.
