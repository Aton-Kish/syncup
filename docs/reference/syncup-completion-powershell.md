## `syncup completion powershell`

<sub><sup>Last updated on 2023-10-19</sup></sub>

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	syncup completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```shell
syncup completion powershell [flags]
```

### Options

```shell
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### See also

- [syncup completion](syncup-completion.md) - Generate the autocompletion script for the specified shell
