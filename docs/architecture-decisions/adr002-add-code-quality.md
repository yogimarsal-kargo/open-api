---
id: adrs-adr002
title: 'ADR002: Code Quality and Security Analysis log'
# prettier-ignore
description: Tools and workflow for Code Quality and Security ANalysis
---

## Decision

A decision was made to adapt `golangci-lint` as Code Quality tools and `horusec` for security analysis
following RFC: https://kargox.atlassian.net/wiki/spaces/ENG/pages/2379710497/RFC+-+Code+Quality+and+Security+Analysis+on+Golang

## Discussion

There is a need to standardize code quality and security analysis measure 
used in Kargo golang application to enable more maintainable software.

## Risks

People ignore issue given by Code Quality / Security Analysis
