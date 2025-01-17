# Changelog

All notable changes to this project will be documented in this file.
See updating [Changelog example here](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]

### Added
- dbaas: property validators
- dbaas: PostgreSQL properties `default_toast_compression` and `max_slot_wal_keep_size`

### Fixed
- dbaas: fractional values in PostgreSQL properties `autovacuum_analyze_scale_factor`, `autovacuum_vacuum_scale_factor` and `bgwriter_lru_multiplier`

### Changed
- dbaas: updated property descriptions
- structured logging with `tflog`

## [2.5.0] - 2022-06-20

### Added
- lbaas: frontend and backend properties
- lbaas: `set_forwarded_headers` frontend rule action
- firewall: allow specifying default rules

### Changed
- New upcloud-go-api version 4.7.0 with context support

## [2.4.2] - 2022-05-10

### Changed
- Update GoReleaser to v1.8.3

## [2.4.1] - 2022-05-05

### Fixed

- server: Remove all tags when tags change into an empty value
- server: Delete unused tags created with server resource on server delete
- server: Improve tags validation: check for case-insensitive duplicates, supress diff when only order of tags changes, print warning when trying to create existing tag with different letter casing
- dbaas: require that both `maintenance_window_time` and `maintenance_window_dow` are set when defining maintenance window
- dbaas: `maintenance_window_time` format

### Changed
- New upcloud-go-api version v4.5.1
- Update terraform-plugin-sdk to v2.15.0
- Update Go version to 1.17

## [2.4.0] - 2022-04-12

### Added
- Support for UpCloud Managed Load Balancers (beta)

### Fixed
- dbaas: upgrading database version

## [2.3.0] - 2022-03-14

### Added
- object storage: allow passing access and secret key as environment variables
- object storage: enable import feature
- storage: add support for autoresizing partition and filesystem

### Fixed
- dbaas: fix PostgreSQL properties: pg_stat_statements_track, pg_partman_bgw_role, pg_partman_bgw_interval

## [2.2.0] - 2022-02-14

### Added

- storage: upcloud_storage data source to retrieve specific storage details

### Fixed

- docs: set provider username and password as required arguments
- provider: return underlying error from initial login check instead of custom error
- provider: fix dangling resource references by removing a binding to an remote object if it no longer exists
- provider: fix runtime error when importing managed database


## [2.1.5] - 2022-01-27

### Fixed

- storage: fix missing backup_rule when importing resource
- provider: fix user-agent for release builds
- server: fix missing template id if resource creation fails on tag errors

### Changed

- Update documentation


## [2.1.4] - 2022-01-18

### Added

- server: validate plan and zone field values before executing API commands
- Support for UpCloud Managed Databases
- Support for debuggers like Delve

### Fixed

- firewall: fix missing server_id when importing firewall resource
- firewall: change port types from int to string to avoid having zero values in state when importing rules with undefined port number(s).
- firewall: remove proto field's default value "tcp" as this prevents settings optional fields value to null and update validator to accept empty string which corresponds to any protocol
- object storage: fix issue where order of storage buckets in an object storage resource would incorrectly trigger changes
- server: return more descriptive error message if subaccount tries to edit server tags

### Changed

- Upgraded terraform-plugin-sdk from v2.7.1 to v2.10.0

## [2.1.3] - 2021-11-18

### Added

- Added title field to the server resource

### Fixed

- server: fix custom plan updates (cpu/mem)

### Changed

- server: new hostname validator

## [2.1.2] - 2021-11-01

### Added

- Added simple backups support (#188)

### Fixed

- Prevent empty tags from replanning a server (#178)
- Make sure either storage devices or template are required on the server resource

## [2.1.1] - 2021-06-22

### Fixed

- fix(client): fix user-agent value (#165)

## [2.1.0] - 2021-06-01

### Added

- Support for UpCloud ObjectStorage S3 compatible storage.
- Add host field to the server resource
- server: add tags attribute support (#150) 
- chore: Add more examples

### Fixed

- Server not started after updating storage device
- router: fix creation of attachedNetworks for routers #144
- chore: fix example in upcloud_tag #125
- server: prevent some attribute update from restarting (#146) 
- router: allow detaching router and deleting attached routers (#151) 
- storage: check size before cloning a device (#152)
- storage: fix address formating (#153)

### Changed

- Update documentation
- Update README

### Deprecated

- tag resource
- zone and zones datasources
- tag datasource

## [2.0.0] - 2021-01-27

### Added

- Missing documentation server resource [#89](https://github.com/UpCloudLtd/terraform-provider-upcloud/issues/89)
- Missing documentation for zone datasource [#120](https://github.com/UpCloudLtd/terraform-provider-upcloud/issues/120)
- New [examples](../blob/master/examples) of using the provider
- Updated workflow to run acceptance tests when opening pull request / pushing to master
- Add user-agent header to the requests
- Can now explicitly set IP address for network interfaces (requires special priviledes for your UpCloud account)
- Expose metadata field for server resource

### Changed

- **Breaking**: the template (os storage) is described with a separate block within the server resource, note that removing / recreating server resource also recreates the storage
- **Breaking**: other storages are now managed outside of the server resource and attached to server using `storage_devices` block

### Removed

- Moved multiple utility functions to `/internal`

### Fixed

- Better drift detection [#106](https://github.com/UpCloudLtd/terraform-provider-upcloud/issues/106)
- Fixed issue where a change in server storages would replace the server network interfaces and recreate the server
- Addressed issue where a change in server networking would replace the previous storages (the template will still be created anew)
- Inconsistent documentation

## [1.0.0] - 2020-10-19

Updated upcloud-go-api, added build/CI scripts, and repackaged 0.1.0 as 1.0.0.

## [0.1.0] - 2020-09-24

### Added

- Changelog to highlight key alterations across future releases
- Website directory for future provider documentation
- Vendor directory through go modules to cover CI builds
- datasource_upcloud_hosts to view hosts data
- datasource_upcloud_ip_addresses to retrieve account wide ip address data
- datasource_upcloud_networks to retrieve account wide networks data
- datasource_upcloud_tags to retrieve account wide tag data
- datasource_upcloud_zone to retrieve specific zone details
- datasource_upcloud_zones to retrieve account wide zone data
- resource_upcloud_firewall_rules add to allow rules to be applied to server
- resource_upcloud_floating_ip_address to allow the management of floating ip addresses
- resource_upcloud_network to allow the management of networks
- resource_upcloud_router to allow the management of routers

### Changed

- README and examples/README to cover local builds, setup and test execution
- Go version to 1.14 and against Go master branch in Travis CI
- Travis CI file to execute website-test covering provider documentation
- Provider uses Terraform Plugin SDK V2
- resource_upcloud_server expanded with new functionality from UpCloud API 1.3
- resource_upcloud_storage expaned with new functionality from UpCloud API 1.3
- resource_upcloud_tag expanded to implement read function

### Removed

- Removed storage HCL blocks that failed due to referencing older UpCloud template ID
- Removed the plan, price, price_zone and timezone UpCloud resources
- resource_upcloud_ip removed and replaced by resource_upcloud_floating_ip_address
- resource_upcloud_firewall_rule removed and replaced by resource_upcloud_firewall_rules
- resource_upcloud_zone removed and replaced by zone and zones datasources

[Unreleased]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.5.0...HEAD
[2.5.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.4.2...v2.5.0
[2.4.2]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.4.1...v2.4.2
[2.4.1]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.4.0...v2.4.1
[2.4.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.5...v2.2.0
[2.1.5]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.4...v2.1.5
[2.1.4]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.3...v2.1.4
[2.1.3]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.2...v2.1.3
[2.1.2]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.1...v2.1.2
[2.1.1]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/v2.1.0...v2.1.1
[2.1.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/2.0.0...v2.1.0
[2.0.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/1.0.0...2.0.0
[1.0.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/compare/0.1.0...1.0.0
[0.1.0]: https://github.com/UpCloudLtd/terraform-provider-upcloud/releases/tag/0.1.0
