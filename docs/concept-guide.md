# Concept guide

## Why syncup?

AWS AppSync makes it easy to build GraphQL APIs through the management console.
However, API version-control is not straightforward in this scenario.

The syncup provides API snapshot capture and restoration features, helping with version management.

## Compatibility with AWS AppSync

The syncup is compatible with AWS AppSync in most cases.
However, the only exception is with the identifier for AppSync Function:
in contrast to AWS AppSync, which uses the Function ID as an identifier, syncup uses the Function Name as an identifier.
This makes it easier to migrate to another AppSync GraphQL API.
