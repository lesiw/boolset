# lesiw.io/boolset

[![Go Reference](https://pkg.go.dev/badge/lesiw.io/boolset.svg)](https://pkg.go.dev/lesiw.io/boolset)
[![CI](https://github.com/lesiw/boolset/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/lesiw/boolset/actions/workflows/main.yml)
[![Release](https://img.shields.io/github/v/tag/lesiw/boolset?sort=semver&label=release)](https://github.com/lesiw/boolset/tags)
[![Go Version](https://img.shields.io/github/go-mod/go-version/lesiw/boolset)](../go.mod)
[![Discord](https://img.shields.io/discord/1145827224516300971?logo=discord&logoColor=white&color=5865F2&label=discord)](https://lesiw.dev/discord)
[![License](https://img.shields.io/github/license/lesiw/boolset)](../LICENSE)

An `analysis.Analyzer` that finds sets dressed as `map[T]bool`.

A `map[T]bool` that only ever stores `true` and is read for
membership is a set wearing the wrong type: `map[T]struct{}` says
so, and its values take no space.

## Checks

### True-only map[T]bool as a set

```go
m := make(map[string]bool) // m is a set and can be map[T]struct{}
m["k"] = true
if !m["k"] {
```

`delete`, `clear`, and `len` mean the same thing for a set, so they
do not disqualify a map. Each variable is reported once, at its
declaration. Only maps created inside a function — by `make`
or a composite literal — are examined, since only there is every
store in view.

Silent, deliberately: storing anything but literal `true` — `false`,
a variable, a computed value, whether by assignment or as a
composite-literal element — marks the value as meaningful state,
as does reading it as one (a consumed comma-ok value, a range over
values). A map that arrives from
elsewhere — a parameter, a struct field, a package variable, an
aliased value — or escapes into a call is not examined: its stores
are not all in view.

## Usage

```sh
go get -tool lesiw.io/boolset/cmd/boolset
go tool boolset ./...
```
