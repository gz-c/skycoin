# Skycoin

[![GoDoc](https://godoc.org/github.com/skycoin/skycoin?status.svg)](https://godoc.org/github.com/skycoin/skycoin) [![Go Report Card](https://goreportcard.com/badge/github.com/skycoin/skycoin)](https://goreportcard.com/report/github.com/skycoin/skycoin)

Skycoin is a next-generation cryptocurrency.

Skycoin improves on Bitcoin in too many ways to be addressed here.

Skycoin is small part of OP Redecentralize and OP Darknet Plan.

## Links

* [skycoin.net](https://www.skycoin.net)
* [Skycoin Blog](https://blog.skycoin.net)
* [Skycoin Blockchain Explorer](https://explorer.skycoin.net)
* [Skycoin Distribution Event](https://event.skycoin.net)

## Installation

### Prerequisites

Install go1.8+.

*Note: In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source.*

### Fetch the source

```sh
go get https://github.com/skycoin/skycoin
```

This will download `github.com/skycoin/skycoin` to `$GOPATH/src/github.com/skycoin/skycoin`.

If you clone the repo with `git clone`, make sure to place it at `$GOPATH/src/github.com/skycoin/skycoin`.

### Run

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
./run.sh
```

### Options

```sh
cd $GOPATH/src/github.com/skycoin/skycoin
./run.sh --help
```

## API Documentation

### Wallet REST API

[Wallet REST API](src/gui/README.md).

### JSON-RPC 2.0 API

[JSON-RPC 2.0 README](src/api/webrpc/README.md).

### Skycoin command line interface

[CLI command API](cmd/cli/README.md).

## Development

We have two branches: `master` and `develop`.

`develop` is the default branch and will have the latest code.

`master` will always be equal to the current stable release on the website, and should correspond with the latest release tag.

### Modules

* `/src/cipher` - cryptography library
* `/src/coin` - the blockchain
* `/src/daemon` - networking and wire protocol
* `/src/visor` - the top level, client
* `/src/gui` - the web wallet and json client interface
* `/src/wallet` - the private key storage library
* `/src/api/webrpc` - JSON-RPC 2.0 API
* `/src/api/cli` - CLI library

### Formatting

All `.go` source files should be formatted with `gofmt` or `goimports`.

### Dependencies

Dependencies are managed with [dep](https://github.com/golang/dep).

To install `dep`:

```sh
go get -u github.com/golang/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

After adding a new dependency (with `dep ensure`), run `dep prune` to remove any unnecessary subpackages from the dependency.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
dep init
dep prune
```

Update all dependencies:

```sh
dep ensure -update -v
dep prune
```

Add a single dependency (latest version):

```sh
dep ensure github.com/foo/bar
dep prune
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
dep ensure github.com/foo/bar@tag
dep prune
```

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Compile and add the wallet source to git:

```sh
cd src/gui/static
npm install
gulp build
git add .
```

### Releases

*TODO: Full instructions on doing a release. Need instructions on updating version number in source code, and running scripts to build releases*

When ready to do a release, a pull request merging `develop` into `master` must be made.

After merging to `master`, tag the branch with the version number.

Once `master` branch is updated, `git checkout master` and create the release builds.

#### Creating release builds

Install Gox for cross-compilation:

```sh
go get github.com/mitchellh/gox
```

Make sure that the wallet dist is up to date:

```sh
cd /src/gui/static
npm install
gulp dist
```

*TODO: Finish release build instructions*
