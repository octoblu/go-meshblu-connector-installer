# Meshblu Connector Assembler

[![Build Status](https://travis-ci.org/octoblu/go-meshblu-connector-assembler.svg?branch=master)](https://travis-ci.org/octoblu/go-meshblu-connector-assembler)
[![GoDoc](https://godoc.org/github.com/octoblu/go-meshblu-connector-assembler?status.svg)](https://godoc.org/github.com/octoblu/go-meshblu-connector-assembler)
[![Test Coverage](https://codecov.io/gh/octoblu/go-meshblu-connector-assembler/branch/master/graph/badge.svg)](https://codecov.io/gh/octoblu/go-meshblu-connector-assembler)
[![Slack Status](http://community-slack.octoblu.com/badge.svg)](http://community-slack.octoblu.com)

# Table of Contents

* [Introduction](#introduction)
* [Getting Started](#getting-started)
  * [Install](#install)
* [Usage](#usage)
  * [Help](#help)
  * [Arguments](#arguments)
  * [Example](#example)

# Introduction

The Meshblu Connector Assembler is a utility that downloads, configures and sets up the service files. This utility is used during the install connector step in the [Meshblu Connector Installer](https://github.com/octoblu/electron-meshblu-connector-installer).

# Getting Started

## Install

```bash
go install github.com/octoblu/go-meshblu-connector-assembler
```

# Usage

## Help

```bash
go-meshblu-connector-assembler --help
```

## Debug

```bash
env DEBUG='meshblu-connector-assembler*' go-meshblu-connector-assembler
```

## Arguments

* `--connector`, `-c` *String* The connector name without `meshblu-` or `meshblu-connector-`.
* `--uuid`, `-u` *String* The Meshblu UUID of the connector you wish to install.
* `--token`, `-t` *String* The Meshblu Token of of connector you wish to install.
* `--github-slug`, `-g` *String* The github owner and repo, separated by a slash. Example: `octoblu/meshblu-connector-say-hello`
* `--tag`, `-T` *String* The release tag of the connector to download. Example `v1.0.0`.
* `--ignition`, `-i` *String* The tag for the github release of [go-meshblu-connector-ignition](https://github.com/octoblu/go-meshblu-connector-ignition). Example: `v1.0.0`.
* `--debug`, `-d` *Bool* **(optional)** If specified this will wait for user input when setting up the service files on windows.

## Example

Assemble a connector

```bash
go-meshblu-connector-assembler \
  --connector say-hello \
  --uuid 6f87a9bc-02e6-41a8-b265-a1a58b8e569a \
  --token 8b265a1a55696f87a9bc041a82e6b8ea \
  --github-slug octoblu/meshblu-connector-say-hello \
  --tag v6.0.0 \
  --ignition v6.0.0
```

Run with debug

```bash
env DEBUG='meshblu-connector-assembler*' \
  go-meshblu-connector-assembler \
  --connector say-hello \
  --uuid 6f87a9bc-02e6-41a8-b265-a1a58b8e569a \
  --token 8b265a1a55696f87a9bc041a82e6b8ea \
  --github-slug octoblu/meshblu-connector-say-hello \
  --tag v6.0.0 \
  --ignition v6.0.0 \
  --debug
```
