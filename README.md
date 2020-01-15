# CRUSH - Code Review Helps Us

A simple tool to assist with code review.

![Crush](./crush-logo.gif)

## Installation

Download the appropriate executable from our "releases" page.

## Usage

The most basic use will run the tool on a directory of code.

`crush examine --directory /your/code/here`

## Docker Usage



## Setting Expectations

We do a fair amount of code review.  As we do that, some things present
that are worthy of review pretty much whenever we see them.  They are
_candidate_ issues.  The certainty for any given item may be low, but we
put them into this tool because we want to review them.

## Local Installation

`go get github.com/jemurai/crush`

## Building Cross Platforms

If you want to build for cross platform use, you can use the build.sh script packaged with crush.

```sh
git clone github.com/jemurai/crush
cd $GOPATH/src/github.com/jemurai/crush
build.sh github.com/jemurai/crush
```

## Building The Docker Image

## Running

```
docker run -v /Users/mk/area51/triage:/tmp/triage jemurai/crush:0.1 examine --debug true --directory /tmp/triage
```

