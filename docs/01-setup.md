# Gocontentful Setup

## Installation

Prerequisite: you need Go 1.21. Upgrade if you still haven't, then run:

```shell
go get github.com/foomo/gocontentful
```

You can run gocontentful from the repository's main folder with

```
go run main.go LIST_OF_PARAMS
```

or build the binary yourself. If you trust us there are [precompiled versions on Github](https://github.com/foomo/gocontentful/releases)

On Mac OS:

```shell
brew install foomo/gocontentful/gocontentful
```

Test the installation (make sure $GOPATH/bin is in your $PATH):

```shell
$ gocontentful -version
v1.1.0
```

## Optional tools

Gocontentful requires a CMA API key to scan the Contentful space and generate the model.
This can be passed as a CLI parameter but that's tedious, and your management key will remain in your shell history.
A better approach is to log in to Contentful using the official _Contentful CLI_. Gocontentful will get the key automatically.
To install the CLI refer to the [official documentation at Contentful.com](https://www.contentful.com/developers/docs/tutorials/cli/installation/).

After installing the CLI, log in inside your terminal with:

```shell
$ contentful login
```

After a roundtrip through the Web authentication pages at Contentful, you'll be logged in:

```
A browser window will open where you will log in (or sign up if you don’t have an account), authorize this CLI tool and paste your CMA token here:

? Continue login on the browser? Yes
? Paste your token here: *******************************************

Great! You've successfully logged in!
╭─────────────────────────────────────────────────────────────────────────╮
│                                                                         │
│   Your management token: **************                                 │
│   Stored at: /Users/yourusername/.contentfulrc.json                     │
│                                                                         │
╰─────────────────────────────────────────────────────────────────────────╯
```

## Generate a client for your space

The gocontentful command accepts the following parameters:

```shell
$ gocontentful -help

ERROR: Please specify either a Contentful Space ID and CMA access token or an export file name

SYNOPSIS
  gocontentful -spaceid SpaceID -cmakey CMAKey [-contenttypes firsttype,secondtype...lasttype] path/to/target/package

  -cmakey string
    	[Optional] Contentful CMA key
  -contenttypes string
    	[Optional] Content type IDs to parse, comma separated
  -environment string
    	[Optional] Contentful space environment
  -exportfile string
    	Space export file to generate the API from
  -help
    	Print version and exit
  -spaceid string
    	Contentful space ID
  -version
    	Print version and exit

Notes:
- The last segment of the path/to/target/package will be used as package name
- The -cmakey parameter can be omitted if you logged in with the Contentful CLI
```

Notes:

- The last segment of the path/to/target/package will be used as package name
- You need to pass gocontentful either cmakey/spaceid (and optional environment) to generate the API from a live space or exportfile to generate it from a local space export file. The cmakey can be omitted if you are logged in through the Contentful CLI.

Assuming you are logged in with the Contentful CLI and your space id at Contentful is xyz123, you can now generate the Gocontentful client files for a package named "myclient" with:

```
gocontentful -spaceid xyz123 myclient
```

This will create all the necessary Go files inside the `myclient` directory. When that happens without errors, your setup is complete.
