# Code Quality, Static Code Analysis, and Security Analysis
This document details approach of code quality and security analysis used in `go-testapp` Application. It follow closely approach defined in https://kargox.atlassian.net/wiki/spaces/ENG/pages/2379710497/RFC+-+Code+Quality+and+Security+Analysis+on+Golang

## Code Quality and Static Code Analysis
Code quality and static code analysis is based on `golangci-lint` tools, which is checked every Pull Request to ensure Pull Request doesn't introduce Code quality also checked on before commit via `pre-commit`

To adapt code quality to new or existing codebase, below are relevant file/structure should be exported and adjusted:
```
# relevant files:
.github
|─ workflows
|── pr.yaml # Github action, with job `golangci-lint` to trigger PR check of golangci-lint
.golangci.yml # Set config of `golangci-lint`
.pre-commit-config.yaml # pre-commit configuration for triggering golangci-lint every commit
```

## Security Analysis
Security analysis is based on `horusec` tools, which is checked every Pull Request to ensure Pull Request doesn't introduce new issue to codebase. 
Security analysis is also checked on before commit via `pre-commit` to ensure leaks and credential issue not committed into git history.

To adapt Security analysis to new or existing codebase, below are relevant file/structure should be exported and adjusted:
```
# relevant files:
.github
|─ workflows
|── pr.yaml # Github action, with job `horusec` to trigger PR check of security analysis
horusec-config.json # Set config of `horusec`,particularly to ignore false positive
.pre-commit-config.yaml # pre-commit configuration for triggering horusec every commit
.trivyignore # Trivy ignore file, for ignoring certain CVE found by trivy, one of security tools used by horusec
```

## Issue tracking
Codequality and Security analysis issue for project would be tracked via `sonarqube` (https://sonarqube.helios.kargo.tech/). 
Issue tracked is code coverage, code quality (via `golangci-lint`), and security issue (via `horusec`)

To adapt Issue tracking to new or existing codebase, below are relevant file/structure should be exported and adjusted:
```
# relevant files:
.github
|─ workflows
|── push.yaml # Github action, with job `ci-sonarscanner` to trigger issue tracking every merge to master/main
Makefile # Task `ci-sonarqube-report`, which create report data for `sonarscanner`
sonar-project.properties # Config of `sonarscanner`, particularly need to ensure that `sonar.projectKey` is unique per project
```

Beside adapting file/structure, to adapt issue tracking, engineer also need to add/adjust github repository setting (via Github Repo Page -> Settings -> Secrets), especially: `SONAR_TOKEN` and `SONAR_HOST_URL`. 
Both of these secret (and `sonar.projectKey`) need to be made and registered on SonarQube dashboard. 
To create and register project, go to SonarQube dashboard (https://sonarqube.helios.kargo.tech/), 
and follow instruction to create project from github (Projects -> Create Project -> Github).

## Multi Project in Single Repository
In case your project contains multiple project in a single repository, the recommended strategy is to:
- Run `golangci` and test coverage for each project
- Run `horusec` analysis and `sonarqube` reporting for entire repository

Below example of such strategy applied.
```
.github
|─ workflows
|── project1.pr.yaml # Project specific PR workflow, for running `golangci-lint`
|── pr.yaml # Generic PR workflow for running `horusec`
|── push.yaml # Github action, with job `ci-sonarscanner` to run code quality capture
project1
|─ .golangci.yml # Setting for golangci per project basis
Makefile # Makefile adjusted for `ci-sonarqube-report` to collect all project `golangci-lint` and test coverage data
.pre-commit-config.yaml # pre-commit adjusted for `golangci-lint` to each project
sonar-project.properties # `sonar.go` external report path (`test`, `coverage`, `golangci-lint`) should be adjusted to include all project directory
horusec-config.json # `horusecCliWorkDir.go` should be adjusted to include all project directory
```
