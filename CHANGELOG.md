# Changelog

All notable changes to the Athlete Unknown API will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [auth0workflow]

### Added

- Username management functionality for Auth0 integration
- `PUT /v1/user/username` endpoint for authenticated users to update their display name
- Username validation utility with profanity filtering
  - Validates length (3-20 characters)
  - Validates format (alphanumeric + spaces only)
  - Checks for inappropriate content
- Default username generation utility (`GenerateDefaultUsername`)
  - Format: "Guest{1-999}{FirstLetterOfEmail}"
  - Example: "Guest42J" for john@example.com
- Comprehensive unit tests for username utilities
- New permission: `update:athlete-unknown:profile` for username updates

### Technical Details

- Username validation includes:
  - Length constraints (3-20 chars)
  - Character restrictions (letters, numbers, spaces only)
  - No consecutive spaces
  - Basic profanity filter (can be enhanced with external library)
- Creates user stats entry if it doesn't exist when setting username
- Uses existing Auth0 JWT middleware for authentication
- Supports both new user creation and username updates

## [v1.0.2] - 2026-01-11

## [PR-21]

### Changed

- Add S3 bucket to upload lambda to

## [PR-20]

### Changed

- Fix lambda build commands. Settle on x86_64 architecture for simplicity

## [PR-18]

### Added

- New deploy-backend.yml file. Deployments will only be made on release/ branches

### Changed

- Move physical attributes: height and weight to "Bio" tile, not "Player Information"
- Changed model. Variable name "tilesFlipped" -> "filppedTiles"
- Modify response of /rounds and /upcoming-rounds to be just round summaries
- Use DynamoDB query with GSI instead of scan for /rounds and /upcoming-rounds
- Handle currentDailyStreak to increment with daily interaction vs sequential play
- Fix user stats bugs

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
  - **deploy-backend.yml**: Automatically deploys to dev environment when release branches are created
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
