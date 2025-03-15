# freebsd-archive-combiner

`freebsd-archive-combiner` is a tool for downloading and combining split files from the FreeBSD archives for older (EOL) versions. In older FreeBSD versions, distribution files are split into multiple smaller files (base.aa, base.ab, ...). This tool automatically downloads these files and combines them into a single tgz file.

## Features

- Automatically downloads split files from FreeBSD archives
- Combines downloaded split files
- Configuration using YAML files
- Support for multiple components (base, kernel, etc.)
- Reuses existing files (avoids redundant downloads)

## Installation

```bash
$ go build -o freebsd-archive-combiner cmd/main.go
```

## Usage

1. Create a YAML configuration file (see the sample in `examples/8.4-RELEASE.yaml`)
2. Run the following command:

```bash
$ ./freebsd-archive-combiner -c examples/8.4-RELEASE.yaml
```

## Output Directory Structure

The combined files are saved in the following directory structure:

```
output/
└── 8.4-RELEASE/
    └── i386/
        ├── fetch/
        │   ├── base/
        │   │   ├── base.aa
        │   │   ├── base.ab
        │   │   └── ...
        │   └── kernels/
        │       ├── generic.aa
        │       ├── generic.ab
        │       └── ...
        └── combine/
            ├── base.tgz
            └── generic.tgz
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
