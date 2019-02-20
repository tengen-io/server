### Development setup

_NOTE: This repository uses Go Modules, so it must be cloned **outside of your `$GOPATH`**._

    $ cd <some/dir/outside/GOPATH>
    $ git clone git@github.com:camirmas/go-stop-go.git

Install a version of Go `1.11` or higher -- `1.11` is needed for Go Module support.

    $ go version

Install and start PostgreSQL

    $ psql -V

Set up your PostgreSQL database:

* Run `./script/bootstrap`, or:

Create a `postgres` user

    $ psql postgres -c 'CREATE ROLE postgres superuser;'

Create the `go_stop_test` database:

    $ createdb go_stop_test

### Running the test suite

    $ ./script/test
