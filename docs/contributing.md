### Development setup

_NOTE: This repository uses Go Modules, so it must be cloned **outside of your `$GOPATH`**._

    $ cd <some/dir/outside/GOPATH>
    $ git clone git@github.com:camirmas/go-stop-go.git

Install a version of Go `1.11` or higher -- `1.11` is needed for Go Module support.

    $ go version

Install and start PostgreSQL

    $ psql -V

Set up your PostgreSQL database:

Run the `./script/bootstrap` script, or manually create a `postgres` user:

    $ psql postgres -c 'CREATE ROLE postgres SUPERUSER LOGIN;'

### Running the test suite

    $ ./script/test
