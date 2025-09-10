name: Patch Management Policy
acronym: PMP
satisfies:
  TSC:
    - CC8.1
    - CC7.1
majorRevisions:
  - date: Sep 8 2025
    comment: Initial document for 2022 TSC compliance
---

# Purpose and Scope

a. The purpose of this policy is to establish requirements for identifying, evaluating, testing, and deploying software patches and security updates to maintain system security and resilience.

a. This policy applies to all information technology systems, including operating systems, applications, firmware, and security tools within the organization's infrastructure.

a. This policy applies to all employees and contractors responsible for system administration, security, and change management.

# Background

a. Effective patch management is critical to maintaining system security and preventing exploitation of known vulnerabilities. The 2022 TSC revised points of focus emphasize the importance of software patch identification and system resiliency as part of change control processes.

a. Unpatched systems represent significant security risks and can lead to data breaches, service disruptions, and compliance violations.

# References

a. Change Control Policy

a. Risk Management Policy

a. Security Incident Response Policy

a. Asset Inventory and Management Policy

# Policy

a. *Patch Identification and Monitoring:*

    i. The organization shall implement automated tools to identify available patches for all systems in the asset inventory.

    i. Security patches and critical updates must be identified and evaluated within 24 hours of release.

    i. All patches must be cataloged with information including:
        - Affected systems and software versions
        - Severity level and security impact
        - Vendor recommendations and installation instructions
        - Dependencies and prerequisites

a. *Patch Evaluation and Risk Assessment:*

    i. Each patch must be evaluated for applicability, security impact, and potential business disruption.

    i. Critical security patches addressing active threats or high-risk vulnerabilities must be prioritized for immediate deployment.

    i. The risk of not applying a patch must be weighed against the risk of applying it, considering system stability and business operations.

a. *Patch Testing Requirements:*

    i. All patches must be tested in a non-production environment before deployment to production systems.

    i. Testing must verify that patches do not negatively impact system functionality, performance, or security controls.

    i. Emergency patches addressing critical security vulnerabilities may bypass extensive testing with appropriate management approval and enhanced monitoring.

    i. Test results must be documented and approved before production deployment.

a. *Patch Deployment and Scheduling:*

    i. Critical security patches must be deployed within 72 hours of successful testing.

    i. Non-critical patches must be deployed within 30 days of release or during the next scheduled maintenance window.

    i. Patch deployment must follow established change control procedures including:
        - Change request approval
        - Implementation planning
        - Rollback procedures
        - Communication to stakeholders

    i. Emergency patches may be deployed outside of normal change windows with appropriate approvals and notifications.

a. *System Resiliency and Recovery:*

    i. Before applying patches, system backups must be verified and recovery procedures tested.

    i. Rollback procedures must be prepared and tested for all patch deployments.

    i. Systems must be monitored during and after patch deployment to ensure proper operation.

    i. Failed patches must be rolled back immediately, and issues must be investigated and resolved.

a. *Patch Management for Different System Types:*

    i. Critical infrastructure and production systems require additional approval and monitoring during patch deployment.

    i. Development and testing systems may follow accelerated patch schedules to validate patches before production deployment.

    i. Legacy systems requiring specialized handling must have documented patch procedures and compensating controls.

a. *Documentation and Reporting:*

    i. All patch management activities must be logged and documented including:
        - Patch inventory and status
        - Risk assessments and approval decisions  
        - Testing results and deployment records
        - Any issues encountered and resolution steps

    i. Monthly reports must be generated showing patch compliance status across all systems.

    i. Metrics must be maintained on patch deployment timelines and system availability.

a. *Compliance and Audit:*

    i. Patch status must be regularly audited to ensure compliance with this policy.

    i. Systems that cannot be patched due to technical or business constraints must have documented compensating controls.

    i. Patch management processes must be reviewed annually and updated as needed.