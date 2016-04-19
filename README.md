# mysql-rename

Utility to rename MySQL databases; **NOTE**: it moves only tables, not the views and other structures (stored procedures/functions etc).

# Usage

```
Usage:
  mysql-rename [OPTIONS]

Application Options:
      --mysql-dsn= MySQL Data Source Name URI (default: root@tcp(localhost:3306)/) [$MYSQL_DSN]
      --from=      Database to be renamed [$FROM]
      --to=        New database name [$TO]

Help Options:
  -h, --help       Show this help message
```

The source database will be dropped once all its tables have been moved and it's found to be empty at end of such moving operation.

# Building

Once the source code is in a proper `GOPATH` directory structure, run:
```
glide up --no-recursive
make
```

The generated binary will be in `dist/bin/mysql-rename`.

# License

[GNU/GPLv2](./LICENSE)
