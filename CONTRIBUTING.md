# Contributing to RisingWave Console

We're excited that you're interested in contributing to RisingWave Console! This guide will help you get started with making contributions to the project.

## Contributing to Documentation

RisingWave Console's documentation is partially generated from templates located in the [docs/templates](docs/templates) directory. To make documentation changes:

1. Find the relevant template file in the templates directory
2. Make your changes to the template
3. Run `make doc` to regenerate the documentation files

This ensures documentation stays consistent and up-to-date across the project.

## Development Setup

1. Install Anchor CLI and toolchains (only for first time setup)

```shell
go install github.com/cloudcarver/anchor@latest
anchor install --config dev/anchor.yaml .
```

2. Generate Code

```shell
make gen
```

3. Build the web page

```shell
make build-web
```

## Code Style Guidelines

[TODO]
