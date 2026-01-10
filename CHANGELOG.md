# Changelog

All notable changes to the Athlete Unknown API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [PR-16]

### Added

- Added initials and nicknames to model and tile tracking model
- Changed model. Variable name "tilesFlipped" -> "filppedTiles"

### Changed

- Move physical attributes: height and weight to "Bio" tile, not "Player Information"

## [PR-15]

### Added

- Created new entry point for lambda builds
- Template.yaml file for AWS deployments

## [PR-14] (https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/14)

### Added

- POST route for migrating stats

## [PR-12]

### Added

- Scraped nicknames and added to model
- Saved initials as its own field

## [PR-11](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/11)

### Changed

- Upgraded security vulnerability packages

## [PR-10](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/10)

### Added

- GitHub Actions CI/CD workflows for automated testing, security scanning, and deployment
  - **ci.yml**: Runs on merge to main - builds, tests, lints, and performs security scans
  - **pr-checks.yml**: Validates pull requests with comprehensive checks and posts status comments
  - **deploy-dev.yml**: Automatically deploys to dev environment when release branches are created
  - **changelog-reminder.yml**: Reminds contributors to update changelog for PRs

## [PR-8](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/8)

### Changed

- Fixed lots of poor AI generated scraping logic
- DB dependency injection into handler
- Secured and validated web scraping URLs
- Refactored large handlers.go file to be more readable and organized

## [PR-5](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/5)

### Added

- User History field to user stats model
- Subsequent updates to input payload and userId extraction from token

## [PR-4](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/4)

### Added

- Submitting results now also updates user's stats
- Additional refactoring & safeguarding

## [PR-3](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/3)

### Added

- Authenticate API access via Bearer Token claims and Auth0 permissions

### Changed

- Changed net/http library to gin for convenience sake

## [PR-2](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/2)

### Added

- Added POST `/v1/round` endpoint that scrapes, formats, and creates a round
- Added unit tests

## [PR-1](https://github.com/Statsland-Fantasy/athlete-unknown-api/pull/1)

### Added

- Initial API implementation with (local) DynamoDB integration

## [0.0.1] - 2025-11-24

### Added

- Initial commit
