## 24/01/2013

On gocardless/gocardless codebase :

```bash
$ time grep -r -i Controller . | wc -l
$ 13562
$ 17.48 real        17.01 user         0.45 sys

$ time ack -i Controller . | wc -l
$ 6973
$ 2.83 real         2.23 user         0.60 sys

$ time ag -u -i Controller . | wc -l
$ 497
$ 0.03 real         0.04 user         0.08 sys
```

Observations : Ag seems to be ignoring a lot of files by default. Here are the results if we search across all files (text and binaries).

```bash
$ time ag -u -i Controller . | wc -l
$ 13556
$ 1.02 real         1.54 user         1.41 sys
```

### nu [very naïve approach]

Using the `go run` tool, on the same codebas :

```bash
$ time go run nu.go Controller ~/src/gc/gocardless | wc -l
$ 9100
$ 2.17 real         1.69 user         0.46 sys
```

Using the compiled version, on the same codebase :

```bash
$ time ./nu Controller ~/src/gc/gocardless | wc -l
$ 9100
$ 1.96 real         1.52 user         0.43 sys
```

Observations : I have yet to explain the difference of line count between each tool.