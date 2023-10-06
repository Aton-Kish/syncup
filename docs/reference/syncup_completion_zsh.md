## `syncup completion zsh`

<sub><sup>Last updated on 2023-10-06</sup></sub>

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(syncup completion zsh)

To load completions for every new session, execute once:

#### Linux:

	syncup completion zsh > "${fpath[1]}/_syncup"

#### macOS:

	syncup completion zsh > $(brew --prefix)/share/zsh/site-functions/_syncup

You will need to start a new shell for this setup to take effect.


```shell
syncup completion zsh [flags]
```

### Options

```shell
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### See also

- [syncup completion](syncup_completion.md) - Generate the autocompletion script for the specified shell
