# grubprint.io [![Build Status](https://drone.dasa.cc/api/badges/dskinner/grubprint/status.svg)](https://drone.dasa.cc/dskinner/grubprint)

Check-out repository to `$GOPATH/src/grubprint.io`

### Server

Generate datastore for first-time setups:

```bash
$ go generate grubprint.io/usda
```

Generate a new key pair for first-time setups:

```bash
$ cd $GOPATH/src/grubprint.io
$ go run main.go -keygen
$ mv id_rsa* assets/
```

There exists a `reflex.conf` in the project root that can assist with recompiling the server and
local assets on file change events with the `reflex` cli.

```bash
$ go get github.com/cespare/reflex
$ reflex -c reflex.conf
```

## TODO

### grubprint.io/keystore

CoreOS devs recently released this:

* https://github.com/coreos/dex
* https://github.com/coreos/go-oidc

`go-oidc` has packages that are very similar to the keystore package. May be worth
considering replacing keystore if desired.

Counter-point, keystore is small and easy to maintain. It is strict in what it supports
and will not be open to unknown developments. This makes it easier to determine security
threats.

For example, a recent vulnerability in many jwt packages was discovered whereby
a client could declare `alg: none` in the token header, by-passing verification server-side
since many of these libraries support multiple algorithms and do not confirm if an algorithm
is ok to use by default. Package keystore would never suffer from such an oversight since
it only, and always, accepts a single algorithm and even rejects storage of public keys that
do not meet the package's requirements.
