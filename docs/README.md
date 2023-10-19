<!-- markdownlint-disable MD041 -->
<div align="center">
  <img src="https://github.com/Aton-Kish/syncup/raw/main/logo.svg" alt="syncup" title="syncup" width="256" />
</div>

# syncup

A development workflow tool for syncing up with AWS AppSync.

## Why syncup?

AWS AppSync makes it easy to build GraphQL APIs through the management console.
However, API version-control is not straightforward in this scenario.

The syncup provides API snapshot capture and restoration features, helping with version management.

## Installation

### Manually

Download the pre-compiled binaries from the [releases page](https://github.com/Aton-Kish/syncup/releases).

### `go install`

```shell
go install github.com/Aton-Kish/syncup/cmd/syncup@latest
```

## Guides

- [How to syncup](./how-to-syncup.md)

## Reference

- [Command Reference](./reference/README.md)

## Changelog

Refer to the [CHANGELOG](https://github.com/Aton-Kish/syncup/blob/main/CHANGELOG.md).
