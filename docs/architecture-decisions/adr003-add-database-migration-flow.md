---
id: adrs-adr003
title: 'ADR003: Database Migration Flow'
# prettier-ignore
description: Implement concrete database migration flow with default and postdeployment
---

## Decision

A decision was made to implement database migration flow with two kind of migration: `default` and `postdeployment`
following RFC: https://kargox.atlassian.net/wiki/spaces/ENG/pages/2434269209/RFC+-+DB+schema+migration+flow

This would add new Github Action for implementing manual database migration for `postdeployment`, 
and adjustment of script and starting of application to enable automatic database migration of `default`.

## Discussion

There is a need to have a guide for understanding which migration is part of `default` and which is part of `postdeployment`.

## Risks

Unauthorized people could have access manual dispatch Github Action for manual database migration. 
Currently people which has `write` access to repository could trigger Github Action manually.
