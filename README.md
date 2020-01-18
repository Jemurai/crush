# CRUSH - Code Review Helps Us

A simple tool to assist with code review.

![Crush](./crush-logo.gif)

## Installation

Download the appropriate executable from the [releases](https://github.com/Jemurai/crush/releases) page.

Alternatively `docker pull jemurai/crush`.

## Basic Usage

The most basic use will run the tool on a directory of code.

`crush examine --directory /your/code/here`

## Docker Usage

To run the docker image against a local directory, just do this: 

`docker run -v <local-directory>:/tmp/toanalyze jemurai/crush:lastest examine --directory /tmp/toanalyze`

## Setting Expectations

We do a fair amount of code review.  As we do that, some things present
that are worthy of review pretty much whenever we see them.  They are
_candidate_ issues.  The certainty for any given item may be low, but we
put them into this tool because we want to review them.

## Pairing with FKIT

[fkit](https://github.com/jemuria/fkit) is a library we use for handling findings 
per the [OWASP OFF](https://github.com/owasp/off) format and integrating them into 
tool chains.  You can use fkit to create findings in proper JSON, push findings to
GitHub projects, or even push a finding to a particular PR comment.

We recommend tuning the `crush` command to produce the findings you want then test
integrating with `fkit` to get the results where you want them.

## Local Installation and Use

Get your source code by either:

`git clone https://github.com/jemurai/crush` or `go get github.com/jemurai/crush`.

Make changes.  Run locally without building.

`go run main.go examine --directory /your/code/here`

## Building Cross Platforms

If you want to build for cross platform use, you can use the build.sh script packaged with crush.

```sh
git clone github.com/jemurai/crush
cd $GOPATH/src/github.com/jemurai/crush
build.sh github.com/jemurai/crush
```

## Building The Docker Image

This should be as simple as: 

`docker build .` or `docker build -t jemurai/crush:0.1 . -f Dockerfile`

## Advanced Usage

Crush provides some advanced options (tags, extensions and threshold) as configurable knobs 
you can turn to try to ensure that you get the results you want.

Additional documentation will added here.

## Issues and Roadmap 

Are tracked in [GitHub Issues](https://github.com/jemurai/crush/issues/).

## License

Crush is open source and licenced under an Apache license.
