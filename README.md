<!-- markdownlint-disable MD041 -->
<div align="center">
  <img src="./logo.svg" alt="syncup" title="syncup" width="256" />
</div>

# syncup

[![MIT License](https://img.shields.io/github/license/Aton-Kish/syncup)](./LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Aton-Kish/syncup)](https://goreportcard.com/report/github.com/Aton-Kish/syncup)

A development workflow tool for syncing up with AWS AppSync.

By combining it with a version control system, you can version control the AppSync Environment Variables, Schema, Resolvers, and Functions, allowing you to quickly restore versions.

## Installation

Download the pre-compiled binaries from the [releases page](https://github.com/Aton-Kish/syncup/releases).

## Usage

```shell
: dumping
syncup pull --api-id aaaaaa123123123example123

: restoring
syncup push --api-id aaaaaa123123123example123
```

If you want to learn more, check out the [user guide](docs/README.md).

## Troubleshooting

If you think you've found a bug, or something isn't behaving the way you think it should, please raise an [issue](https://github.com/Aton-Kish/syncup/issues/new/choose) on GitHub.

## Changelog

Refer to the [CHANGELOG](./CHANGELOG.md).

## License

The syncup is licensed under the MIT License, see [LICENSE](./LICENSE).
