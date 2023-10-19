## `syncup completion bash`

<sub><sup>Last updated on 2023-10-19</sup></sub>

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(syncup completion bash)

To load completions for every new session, execute once:

#### Linux:

	syncup completion bash > /etc/bash_completion.d/syncup

#### macOS:

	syncup completion bash > $(brew --prefix)/etc/bash_completion.d/syncup

You will need to start a new shell for this setup to take effect.


```shell
syncup completion bash
```

### Options

```shell
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### See also

- [syncup completion](syncup-completion.md) - Generate the autocompletion script for the specified shell
