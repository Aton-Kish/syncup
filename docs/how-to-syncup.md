# How to syncup

Prerequisites:

- AWS account
- AWS IAM credentials
- AWS AppSync GraphQL API

## Dumping AWS AppSync GraphQL API

This command retrieves the AppSync Schema, Resolvers, and Functions to your local.

```shell
syncup pull --api-id aaaaaa123123123example123
```

output example:

```text
v saved schema
v saved function MyFunction
v saved all functions
v saved resolver Query.listTodos
v saved resolver Query.getTodo
v saved resolver Mutation.createTodo
v saved resolver Mutation.updateTodo
v saved resolver Mutation.deleteTodo
v saved all resolvers
```

file tree:

```text
.
├── functions
│   └── MyFunction
│       ├── code.js
│       └── metadata.json
├── resolvers
│   ├── Mutation
│   │   ├── createTodo
│   │   │   ├── metadata.json
│   │   │   ├── request.vtl
│   │   │   └── response.vtl
│   │   ├── deleteTodo
│   │   │   ├── metadata.json
│   │   │   ├── request.vtl
│   │   │   └── response.vtl
│   │   └── updateTodo
│   │       ├── metadata.json
│   │       ├── request.vtl
│   │       └── response.vtl
│   └── Query
│       ├── getTodo
│       │   ├── code.js
│       │   └── metadata.json
│       └── listTodos
│           ├── code.js
│           └── metadata.json
└── schema.graphqls
```

## Restoring AWS AppSync GraphQL API

You can restore the AppSync Schema, Resolvers, and Functions from your local.

```shell
syncup push --api-id aaaaaa123123123example123
```

output example:

```text
v pushed schema
v pushed function MyFunction
v pushed all functions
v pushed resolver Query.listTodos
v pushed resolver Query.getTodo
v pushed resolver Mutation.createTodo
v pushed resolver Mutation.updateTodo
v pushed resolver Mutation.deleteTodo
v pushed all resolvers
```

## Migrating to another AWS AppSync GraphQL API

> [!IMPORTANT]
> The source and target AppSync must have data sources with the same name.

Below is an example of migrating from API ID `aaaaaa123123123example123` to `bbbbbb456456456example456`:

```shell
syncup push --api-id bbbbbb456456456example456
```

output example:

```text
v pushed schema
v pushed function MyFunction
v pushed all functions
v pushed resolver Query.listTodos
v pushed resolver Query.getTodo
v pushed resolver Mutation.createTodo
v pushed resolver Mutation.updateTodo
v pushed resolver Mutation.deleteTodo
v pushed all resolvers
```

## See also

- [Command reference](./reference/README.md)
