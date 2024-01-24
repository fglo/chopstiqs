# chopstiqs

[![Go Reference](https://pkg.go.dev/badge/github.com/fglo/chopstiqs.svg)](https://pkg.go.dev/github.com/fglo/chopstiqs)
[![build](https://github.com/fglo/chopstiqs/actions/workflows/go-build.yaml/badge.svg?branch=main)](https://github.com/fglo/chopstiqs/actions?query=workflow%3Ago-build)
[![release](https://github.com/fglo/chopstiqs/actions/workflows/deploy-webasm.yml/badge.svg?branch=main)](https://github.com/fglo/chopstiqs/actions?query=workflow%3Adeploy-webasm)

![chopstiqs logo generated by DALL·E 2](img/chopstiqs-logo-1-4x4.png)

Chopstiqs aims to be a minimalistic GUI package for the [ebiten](https://ebitengine.org/) engine. Rather than using separate image files, it draws interface elements using built-in drawing functions. This allows for quick prototyping and use in projects that do not need polished graphics.

## Examples

Running example: <https://fglo.github.io/chopstiqs/>

## Components

- buttons
- checkboxes
- labels
- sliders
- containers
- container layouts
  - horizontal list
  - vertical list
  - grid

## Roadmap

What I want to achieve:

- more components:
  - text inputs
  - tooltips
  - radiogroups
  - modals and/or dialogs
  - range sliders
  - ...
- touchscreen support
- container layouts
  - vertical scroll
  - horizontal scroll
  - flexbox
- layers (mulitple containers on top of each other)
- tests for components
