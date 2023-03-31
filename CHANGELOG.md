# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.7.0] - 2023-03-31

- Support for CAPO `v0.7.x` added (supports CAPO api version `v1alpha6`)

## [0.6.0] - 2022-08-15

## Changed

- Increase memory limit from `80Mi` to `200Mi` since we observed `OOMKilled` in production.

## [0.5.0] - 2022-06-27

## Removed

- Remove `OpenStackMachineTemplate` controller.

## [0.4.0] - 2022-06-09

## Changed

- Support for CAPO `v0.6.x` added (supports now CAPO api version `v1beta5`)

## [0.3.0] - 2022-04-26

## Changed

- Requeue OpenStackClusters when deletion is in-progress.
- Improve logging.

## [0.2.0] - 2022-03-22

### Added

- Add `OpenStackMachineTemplate` controller to remove finalizers from unused templates.

## [0.1.0] - 2022-02-15

### Added

- Project initilization.

[Unreleased]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.7.0...HEAD
[0.7.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/giantswarm/cluster-api-cleaner-openstack/releases/tag/v0.1.0
