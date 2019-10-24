# Contributing Guidelines

​First, thank you for wanting to contribute!

The following is a set of guidelines for contributing to Filecoin Project's HTTP REST API.  If you have suggestions for improving this document, please file an issue or PR.

## How can I contribute?

Here at `go-filecoin`, there’s always a lot of work to do. There are many ways you can support the project, from progamming, writing, organizing, and more. Consider these as starting points:

- **Submit bugs**: Perform a cursory [search](https://github.com/filecoin-project/go-http-api/issues) to see if the problem has already been reported. If it does exist, add a 👍 to the issue to indicate this is also an issue for you, and add a comment if there is extra information you can contribute. If it does not exist, [create a new issue](https://github.com/filecoin-project/go-http-api/issues/new/choose) (using the Bug report template).

- **Write code:** Once you've read this contributing guide, check out [Good First Issues](#good-first-issues) for well-prepared starter issues.

## What should I know before getting started?

This is a fairly simple API.  It obeys basic REST methodology. Please read the README and set up a working filecoin node (even if it is local only) before getting started. These two things will give you an idea of how the API is expected to be implemented and used, but also, it will help you understand how to add new features properly.

### Design Before Code
- Write down design intent before writing code, and subject it to constructive feedback.
- Major changes should have a [Design Doc](https://github.com/filecoin-project/designdocs).
- For minor changes, file an issue for your change if it doesn't yet have one, and outline your implementation plan on the issue.

### Pull Requests and Reviews

We review every change before merging to master, using GitHub pull requests.
**Try to keep PRs small**, no more than 400 lines or 8 files. Code reviews are easier, faster, and more effective when the diff is small.

Merging a PR to `master` requires maintainer approval. The following process aims to merge code quickly and efficiently while avoiding both accidental and malicious introduction of bugs, unintended consequences, or vulnerabilities.

1. Committers require approval from any single maintainer, before [landing](#landing-changes) their own PRs (maintainers require approval from any committer).
1. Non-committers require approval first from a committer familiar with the relevant area of code. Once they deem it ready for merge, the reviewer will assign a maintainer for approval.
    - Once approved, the reviewing committer is responsible for landing the commits.

If your PR hasn't been reviewed in 3 days, pinging reviewers via Github or community chat is welcome and encouraged.

We use the following conventions for code reviews:

- "Approve" means approval. If there are comments and approval, it is expected that you'll address the comments before merging. Ask for clarification if you're not sure.
  - *Example: reviewer points out an off by one error in a blocking comment, but Approves the PR. Reviewee must fix the error, but PR can progress to maintainer review. Committer confirms this before merge.*
- "Request Changes" means you don't have approval, and the reviewer wants another look.
  - *Example: the whole design of an abstraction is wrong and reviewer wants to see it reworked.*

- By default, code review comments are advisory: the reviewee should consider them but doesn't _have_ to respond or address them. Comments that start with "BLOCKING" must be addressed and responded to. If a reviewer makes a blocking comment but does not block merging (by marking the review "Add Comments" or "Approve") then the reviewee can merge if the issue is addressed.

In rare cases, a maintainer may request approval from all maintainers for a wide-reaching PR.

#### Reviewer Responsibilities:

**Avoid lengthy design discussions in PR reviews.** Major design questions should be addressed during the [Design Before Code](#design-before-code) step. If the conversation snowballs, prefer to merge and spin out an issue for discussion, and then follow up on the process failure that led to needing a big design discussion in a PR.

It is considered helpful to add blocking comments to PRs that introduce protocol changes that do not appear in the [spec](#the-spec).

### Landing Changes

We strongly discourage merge commits on the `master` branch. When merging to master:
- squash your PR's commits into one commit with an encompassing message, and
- rebase against master so that your commit lands as a "fast-forward", without a merge commit.

At the command line, update your branch with `git fetch origin +master:master; git rebase master`.
"Merge" your changes with `git checkout master; git merge --ff-only <mybranch>`, and push.

On GitHub, use the grey "updated branch" button, which adds a merge commit to your branch, but then land your changes with the "rebase and merge" or "squash and merge" options on the green merge button.
Both of these will rebase your branch, squashing out any trailing merge commits in the process.

We may enable GitHub branch protections requiring that a CI build has passed on the branch, that the branch is up to date with master, or both.
If either protection is not in force, committers should use their best judgement when landing an untested commit.

## Issues and tracking

We use GitHub issues to track all significant work, including design, implementation, documentation and community efforts.
We also use [ZenHub](https://app.zenhub.com/workspaces/filecoin-5ab0036a12f8e82ae4ed60f0/boards?repos=113219518&showPipelineDescriptions=false) to record issue priority and track team progress.
ZenHub adds some useful project management overlay data to GitHub issues.

To pick up an issue:

1. **Assign** it to yourself.
1. **Ask** for any clarifications via the issue, pinging in [community chat](https://github.com/filecoin-project/community#chat) as needed.
1. **Create a PR** with your changes, following the [Pull Request and Code Review guidelines]().

### Contributors

Anyone is welcome. If you have created an issue, commented on an issue or discussion forum thread, or submitted a PR, you are a contributor.

Contributors agree to:

* Participate in the project, following the [Code of Conduct](https://github.com/filecoin-project/community/blob/master/CODE_OF_CONDUCT.md).

Contributors may:
* Open issues and PRs
* Comment on issues and PRs

### Collaborators

A **Collaborator** is someone who demonstrates helpful contributions to the project.

Responsibilities:

* Make helpful contributions via issues, PRs, and other venues

Abilities:
* Write to the repo (branches only - cannot merge to master)
* Manage issues

### Committers

A **Committer** is someone with a broad understanding of the codebase and the judgment and humility to call on others' expertise when needed. They have a consistent track record of quality contributions, regular participation, and enabling others.

Committers agree to:

* Review PRs and guide work to ready-to-merge
* With maintainer approval, rebase and merge contributor PRs
* Issue triage and other project stewardship

Abilities:

* All the things Collaborators may do, plus
* Merge PRs to master

### Maintainers

A **Maintainer** is someone:

1. **who is invested in and broadly familiar with the project**, as demonstrated by a history of significant technical-, process-, and/or project-level contributions;
2. **who deeply understands the system**, especially knowing when and who to defer to as a reviewer, and with an eye towards unintended consequences;
3. **who is actively engaged in project progress and stewardship** by enabling others through project-wide planning, code reviews, design feedback, etc.; and
4. **who is a model of trustworthiness, technical judgement, civility, and helpfulness**.

Maintainers agree to:

* Review: Timely, friendly review of PRs and design docs to ensure high-quality code and grow knowledge of committers and frequent contributors
* Planning and Improvements: Participate meaningfully in technical and process-related improvements at the project level
* Make significant, direct technical contributions
* Backstop for hard problems and general project stewardship (TODO: improve wording)

Abilities:

* Do all the things Contributors and Collaborators may do

## Additional Developer Notes

#### Testing

- All new code should be accompanied by unit tests. Prefer focussed unit tests to integration tests for thorough validation of behaviour. Existing code is not necessarily a good model, here.
- Integration tests should test integration, not comprehensive functionality
- Tests should be placed in a separate package, and follow the naming pattern `$PACKAGE_test`. For example, a test of the chain package should live in a package named `chain_test`. In limited situations, exceptions may be made for some "white box" tests placed in the same package as the code it tests.

#### Conventions and Style

We use the following import ordering.
```
import (
        [stdlib packages, alpha-sorted]
        <single, empty line>
        [external packages]
        <single, empty line>
        [go-filecoin packages]
)
```

Where a package name does not match its directory name, an explicit alias is expected (`goimports` will add this for you).

Example:

```go
import (
	"context"
	"testing"

	cmds "github.com/ipfs/go-ipfs-cmds"
	cid "github.com/ipfs/go-cid"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-filecoin/testhelpers"
	"github.com/filecoin-project/go-filecoin/types"
)
```

Additionally, CI runs golangci-lint on the codebase.  If the linter fails, CI will fail. To prevent you from a bunch of push-round trips to please the linter on CI, we invite you to consider pre-commit hook.  Here is one which runs exactly what is run on CI and performs other checks:

```bash
#!/bin/bash
export GO_NAMES='\.(go)$'
export GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep -E $GO_NAMES)
function exit_err() { echo "❌ 💔" ; exit 1; }

if [[ ! -z $GO_FILES ]]
then
    echo "Examining $GO_FILES"
    echo $GO_FILES | xargs gofmt -s -w || exit_err
    /usr/local/bin/golangci-lint run || exit_err
    echo $GO_FILES | xargs git add
else
    echo "nothing to do 😑"
    exit 0
fi
echo "✅ 😇"
```
