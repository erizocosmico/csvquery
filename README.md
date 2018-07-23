# csvquery [![Build Status](https://travis-ci.org/erizocosmico/csvquery.svg?branch=master)](https://travis-ci.org/erizocosmico/csvquery) [![codecov](https://codecov.io/gh/erizocosmico/csvquery/branch/master/graph/badge.svg)](https://codecov.io/gh/erizocosmico/csvquery) [![Go Report Card](https://goreportcard.com/badge/github.com/erizocosmico/csvquery)](https://goreportcard.com/report/github.com/erizocosmico/csvquery) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Query CSV files using SQL.

## Features

- Interactive REPL to quickly query CSV files with on demand file loading.
- MySQL-compatible server to query CSV files, powered by [go-mysql-server](https://github.com/src-d/go-mysql-server).
- Insertions of new rows (using `INSERT INTO`).

For more info about what subset of SQL is supported, check out the documentation of [go-mysql-server](https://github.com/src-d/go-mysql-server).

## Install

You can install csvquery by downloading one of the prebuilt binaries available in the [releases page](https://github.com/erizocosmico/csvquery/releases).

Or install it manually:

```
go get github.com/erizocosmico/csvquery/...
```

## Usage

There are two way you can use csvquery: REPL or server.

### REPL

You can start the REPL with the following command:

```
csvquery repl
```

If you want to add some CSV files as tables directly when you're starting the command, you can do so with the `-f` flag.

```
csvquery repl -f some/file.csv -f another/file.csv:custom_table_name
```

As you can see in the previous example, you can choose a custom name for the table with the syntax `FILE_PATH:TABLE_NAME`. If no custom name is provided, the table name will be the file name without the extension and with all characters that are not alphanumeric or underscore removed.

If you don't want to load the files on startup time, don't worry, you can load them later using `load`.

Once the REPL has started, you can use the following commands:

- `load <FILE_PATH>[ <TABLE_NAME>]` to load a new file into the database with the given table name (or the default one if it's not provided).
- `help` to show the usage.
- `quit` or `exit` to exit the REPL.
- `<SQL QUERY>` to execute a SQL query.

### Server

The server, unlike the REPL, needs all tables loaded at the start of the command, so you must provide all the CSV files that you want to query as flags using the `-f` flag. You can see an example of how to do that in the [REPL](#REPL) section.

These are all available options for the server:

```
-d, --dbname=   Database name. (default: csv)
-f, --file=     Add file as a table. You can use the flag in the format '/path/to/file:NAME' to give the file a
                specific table name. Otherwise, the file name without extension will be the table name with only
                alphanumeric characters and underscores.
-u, --user=     User name to access the server. (default: root)
-p, --password= Password to access the server.
-P, --port=     Port in which the server will listen. (default: 3306)
-h, --host=     Host name of the server. (default: 127.0.0.1)
```

For example, let's start the server with the test data contained in the `_testdata` directory of this repository.

```
csvquery server -f _testdata/ratings.csv -f _testdata/superheroes.csv
```

Now, let's connect to the server using the `mysql` client tool and execute a query:

```
$ mysql --host=127.0.0.1 --port=3306 -u root
mysql> SELECT COUNT(*) as num_ratings, id
    -> FROM superheroes s
    -> INNER JOIN ratings r ON s.id = r.superhero_id
    -> GROUP BY id;
+-------------+------+
| num_ratings | id   |
+-------------+------+
|           1 | sm   |
|           2 | bm   |
|           2 | ww   |
|           2 | ma   |
+-------------+------+
4 rows in set (0.00 sec)
```

## Roadmap

- [x] Read-only queries with REPL and server.
- [x] Support insertions.
- [ ] Support updates and deletes.
- [ ] Support creating indexes.

## LICENSE

MIT License, see [LICENSE](/LICENSE).