# ReactGoNative

[![Build Status](https://travis-ci.org/steve-winter/reactgonative.svg?branch=master)](https://travis-ci.org/steve-winter/reactgonative)
[![Code Climate](https://codeclimate.com/github/steve-winter/reactgonative/badges/gpa.svg)](https://codeclimate.com/github/steve-winter/reactgonative)
[![Issue Count](https://codeclimate.com/github/steve-winter/reactgonative/badges/issue_count.svg)](https://codeclimate.com/github/steve-winter/reactgonative)
[![Go Report Card](https://goreportcard.com/badge/github.com/steve-winter/reactgonative)](https://goreportcard.com/report/github.com/steve-winter/reactgonative)
[![GoDoc](https://godoc.org/github.com/steve-winter/reactgonative?status.svg)](https://godoc.org/github.com/steve-winter/reactgonative)

## Current Status

Currently in pre-alpha. Code currently successfully generates Android bindings.

## Roadmap

1. Add tests to Android elements
2. Create iOS integration

## Context

This tool is born out of experimentation with creating a single development platform that allows one set of code to run across both Android and iOS platforms. (In principle this could be extended to Desktop but not considered at present)

The [Go](golang.com) language, through [GoMobile](https://github.com/golang/mobile) along with [react-native](https://facebook.github.io/react-native/) have created two methods to enable mobile development on a shared codebase. This tool aims to allow externalising any business logic to a common Go component, with the React components handling pure UI.

Communication between React and Go are via the use of [Promises](https://developer.mozilla.org/en/docs/Web/JavaScript/Reference/Global_Objects/Promise). The generated code is not packaged but instead placed within your code folders. They can be edited as required.

### Constraints
1. One return type (of simple type i.e. String, int) from Go method. Plans to introduce mapping to allow multiple returns, and object returns.
2. Tool does not check if generated code already exists, nor if the call is run from the wrong location.
3. At present relies on GOPATH being set, and you GO package being present in the GOPATH

### Usage
To install:

```sh
$ go get -u go get github.com/steve-winter/reactgonative
```

To use:

```sh
$ cd $MYANDROIDFOLDER
$ reactgonative $MYGOPACKAGE
```
