# Release Checklist

## Overview

This document describes the checklist to publish a release for Ratify CLI via GitHub workflow.

## Release Process

- Check if there are any security vulnerabilities fixed and security advisories published before a release. Security advisories should be linked on the release notes.
- Determine a [SemVer2](https://semver.org/)-valid version prefixed with the letter `v` for release. For example, `version="v2.0.0-alpha.1"`.
- If there is new release in [ratify-go](https://github.com/ratify-project/ratify-go) library that are required to be upgraded in Ratify CLI, submit a PR to update the dependency versions in the `go.mod` and `go.sum` files of Ratify CLI
- Create another PR to update the Ratify CLI version with a single commit. The commit message MUST follow the [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) and could be `bump: tag and release $version`. Record the digest of that commit as `<commit_digest>`. This PR is also used for voting purpose of the new release. Add the link of change logs and repo-level maintainer list in the PR's description. The PR title could be `bump: tag and release $version`. Make sure to reach a majority of approvals from the [repo-level maintainers](MAINTAINERS) before releasing it. This PR should be merged using [Create a merge commit](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/about-merge-methods-on-github) method in GitHub. 
- After the voting PR is merged, execute `git clone git@github.com:ratify-project/ratify-cli.git` to clone the repository to your local file system.
- Enter the cloned repository and execute `git checkout <commit_digest>` to switch to the specified branch based on the voting result.
- Create a tag by running `git tag -am $version $version -s`.
- Run `git tag` and ensure the desired tag name in the list looks correct, then push the new tag directly to the repository by running `git push origin $version`.
- Wait for the completion of the GitHub action [release-github](https://github.com/ratify-project/ratify-cli/actions/workflows/release-github.yml).
- Check the new draft release, revise the release description, and publish the release.
- Announce the new release in the Ratify Project community.
