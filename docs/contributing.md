### Development setup

_NOTE: This repository uses Go Modules, so it must be cloned **outside of your `$GOPATH`**._

    $ cd <some/dir/outside/GOPATH>
    $ git clone git@github.com:tengen-io/server.git

Install a version of Go `1.11` or higher -- `1.11` is needed for Go Module support.

    $ go version

Install and start PostgreSQL

    $ psql -V

#### Automatically

Run the `./script/bootstrap` script, or

#### Manually

    $ psql postgres -c 'CREATE ROLE postgres SUPERUSER LOGIN;'
    $ createdb tengen_test
    $ createdb tengen
    $ ./script/db_migrate tengen
    $ make gen

### Running the test suite

    $ ./script/test

### Run the development server

    $ make
    $ ./tengen

Visit http://localhost:8180
