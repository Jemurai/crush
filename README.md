# CRUSH - Code Review Helps Us

A simple tool to assist with code review.

![Crush](./crush-logo.gif)

## Installation

Download the appropriate executable from the [releases](https://github.com/Jemurai/crush/releases) page.

Alternatively `docker pull jemurai/crush`.

Or use the GitHub Action.

## Basic Usage

The most basic use will run the tool on a directory of code.

`crush examine --directory /your/code/here`

Generally you might want to specify extensions, tags or thresholds.  These help you to run the checks you really want or care about.  As you can imagine, just searching for certain strings can get noisy and some tuning can go a long way.

So you can specify: 

1. Threshold: `--threshold 2` - this tells Crush only to check the things that have a higher threshold than specified.  The default is 5.  You can see the values on the checks in the JSON files.

1. Tag:  `--tag badwords` - this tells crush to run the checks that have this tag.  Some checks are tagged with language or issue types.  The `badwords` tag is taken directly from [this blog post](https://btlr.dev/blog/how-to-find-vulnerabilities-in-code-bad-words) by Will Butler.

1. Extension: `--ext .java` - this tells crush just to run the checks that apply to .java files.

## Docker Usage

To run the docker image against a local directory, just do this: 

`docker run -v <local-directory>:/tmp/toanalyze jemurai/crush:v examine --directory /tmp/toanalyze`

Of course, you can also run this with the above tags and thresholds:

`docker run -v <localdir>:/tmp/target jemurai/crush:v examine --directory /tmp/target --tag badwords --threshold 1 --debug true`

This will generate a lot of output.

## Setting Expectations

We do a fair amount of code review.  As we do that, some things present
that are worthy of review pretty much whenever we see them.  They are
_candidate_ issues.  The certainty for any given item may be low, but we
put them into this tool because we want to review them.

## Check Anatomy

You can find the checks in `/checks/*.json`.  They look like this: 

```json
[{
        "name": "Raw handling of something",
        "description": "Raw handling something where security might be applied at a higher abstraction",
        "magic": "(?i)raw",
        "threshold": 1.0,
        "exts" : [
        ],
        "tags": [
            "badwords"
        ]
    }
]
```

You can see here how the tags, extensions and threshold are set for each check, which is essentially a Golang Regex in the "magic" field.

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

`go run crush.go examine --directory /your/code/here`

## Building Cross Platforms

If you want to build for cross platform use, you can use the build.sh script packaged with crush.

```sh
git clone github.com/jemurai/crush
cd $GOPATH/src/github.com/jemurai/crush
build.sh github.com/jemurai/crush
```

## Building The Docker Image

This should be as simple as: 

`docker build .` or `docker build -t jemurai/crush:v . -f Dockerfile`

To build and push the docker image to dockerhub:
`docker build -t jemurai/crush:v .`
`docker push jemurai/crush:v`

## Advanced Usage

Crush provides some advanced options (tags, extensions and 
threshold) as configurable knobs you can turn to try to ensure 
that you get the results you want.

Additional documentation will added here.

### Compare

```sh
crush examine --compare <file of old findings> --directory /path/to/code
```

This will produce JSON for added new findings in the current source (what is found in the directory).

## Issues and Roadmap 

Are tracked in [GitHub Issues](https://github.com/jemurai/crush/issues/).

## License

Crush is open source and licenced under an Apache license.
