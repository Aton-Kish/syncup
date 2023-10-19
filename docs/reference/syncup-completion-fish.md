## `syncup completion fish`

<sub><sup>Last updated on 2023-10-19</sup></sub>

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	syncup completion fish | source

To load completions for every new session, execute once:

	syncup completion fish > ~/.config/fish/completions/syncup.fish

You will need to start a new shell for this setup to take effect.


```shell
syncup completion fish [flags]
```

### Options

```shell
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### See also

- [syncup completion](syncup-completion.md) - Generate the autocompletion script for the specified shell
