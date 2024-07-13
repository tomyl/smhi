# `smhi` ☀️

![CI](https://github.com/tomyl/smhi/actions/workflows/ci.yml/badge.svg?branch=main&event=push)
[![Go Reference](https://pkg.go.dev/badge/github.com/tomyl/smhi.svg)](https://pkg.go.dev/github.com/tomyl/smhi)

Go library for consuming the [SMHI meteorological forecasts API](https://opendata.smhi.se/apidocs/metfcst/index.html). This project is not affiliated with SMHI in any way.

**Pre-alpha software**. Expect API breakage, crashes, data loss, silent data corruption etc.

## Installation

```bash
$ go get github/tomyl/smhi
```

To install the example application:
```bash
$ go install github/tomyl/smhi/cmd/smhi@latest
$ smhi -lon 18.040468 -lat 59.340379        
Time       Weather                   Temperature  Max precipitation  Wind speed
Sat 17:00  ☁  Overcast               20.7°C       0.0 mm/h           4.8 m/s
...
```

## Usage

See the [example application](<./cmd/smhi/main.go>).

