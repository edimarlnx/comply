# {{.Name}} Compliance Program

High-level TODO created by [comply](https://github.com/strongdm/comply)

## Initialization Phase (hours)
- [ ] Add project to source control
- [ ] Verify `comply build` generates valid output
- [ ] Create ticketing credentials, configure via `comply.yml`
- [ ] Verify `comply sync` executes without errors

## Authoring Phase (weeks)
- [x] Validate standards/, pruning as necessary
    - [x] Updated to TSC-2022 with revised points of focus
- [ ] Customize narratives/
- [x] Customize policies/
    - [x] Enhanced vendor management policy for CC9.2 compliance
    - [x] Created asset inventory policy for CC2.1 requirements  
    - [x] Enhanced access policy for CC6 guidance
    - [x] Created patch management policy for CC8.1 compliance
    - [ ] Distribute remaining controls among policies
    - [ ] Ensure all policies address latest TSC-2022 controls
- [ ] Customize procedures/
    - [ ] Distribute controls among procedures
    - [ ] Create valid ticket templates
    - [ ] Assign schedules
- [ ] Verify `comply todo` indicates all controls satisfied

## Deployment Phase (weeks)
- [ ] Deploy `comply scheduler` (see README.md for example script)
- [ ] Deploy `comply build` output to shared location
- [ ] Distribute policies to team
- [ ] Train team on use of ticketing system to designate compliance-relevant activity

## Operating Phase (eternal)
- [ ] Monitor timely ticket workflow
- [ ] Adjust and re-publish narratives, policies and procedures as necessary

## Audit Phase (weeks, annually)
- [ ] Import request list (tickets will be generated)
- [ ] Fulfill all request tickets
    - [ ] Attach policies, procedures, and narratives
    - [ ] Attach evidence collected by previously-executed procedure tickets
