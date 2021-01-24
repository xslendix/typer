# Typer

Typer is a cli tool/game for testing your WPM, CPM and \*accuracy.

[![asciicast](https://asciinema.org/a/aw2whfEonVSOfQ3ga0msBX1Sz.svg)](https://asciinema.org/a/aw2whfEonVSOfQ3ga0msBX1Sz?t=5)

## Requirements

 - go

## Installing

```
go get -u github.com/xslendix/typer/...
typer
```

## Builing from source

```
cd go/src
git clone --depth 1 https://github.com/xSlendiX/typer
cd typer
go mod tidy
mkdir -p ~/.local/share
cp textdata ~/.local/share
go install
typer
```

