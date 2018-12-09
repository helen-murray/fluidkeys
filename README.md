# Fluidkeys command line

[![Build Status](https://travis-ci.org/fluidkeys/fluidkeys.svg?branch=master)](https://travis-ci.org/fluidkeys/fluidkeys)

Fluidkeys helps teams protect themselves with strong encryption. It builds on the OpenPGP standard and is compatible with other OpenPGP software.

0.2 helps you create a key or import one from `gpg`, then automatically maintain it.

Once installed, run it with `fk`.

## Install

Head over to [download.fluidkeys.com](https://download.fluidkeys.com) <!--It might be easier if the links opened in a new tab-->

## Install from source

You'll need the [Go compiler](https://golang.org/dl/)

Clone the repo:

```
REPODIR=$(go env GOPATH)/src/github.com/fluidkeys/fluidkeys

git clone https://github.com/fluidkeys/fluidkeys.git $REPODIR
cd $REPODIR
```

Build and install to `/usr/local/bin/fk`:

```
make && sudo make install <!--I got stuck here with an error 'make: *** No rule to make target `install'.  Stop.' I made it work in the end (I think) by going to finder and searching for the makefile then going to that directory and running the command there. My makefile was in go/src/github.com/fluidkeys/fluidkeys/. I got error 23 when it completed -->
```

If you prefer to run without `sudo` (root), install into `$HOME/bin/fk`:

```
PREFIX=$HOME make install
```

## Develop

If you want to hack on Fluidkeys locally you'll need [Go 1.10+](https://golang.org/dl/) and [dep](https://github.com/golang/dep#installation).

Run:

```
make run <!-- I got error 2 when this completed but it still seemed to work I think -->
```
