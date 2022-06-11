# gozt

Collection of tools for taking Zettelkasten notes in terminal.

* print all note names with `zt-a`,
* print the last note name with `zt-l`,
* print the last but one note name with `zt-ll`,
* print a new note name with `zt-n`,
* print a Dot formatted graph with `zt-g`.

Notes are written in plain text files named with 4 digits. Connections inside
of a note are marked as `$0013`.

## How I build it

```sh
gmake  # or use `make`
```

## How I use it

To start a new note I do ``vim `zt-n` ``.

To come back to the last note I've just closed I do ``vim `zt-l` ``.

To count notes I do `zt-a | wc -l`.

To draw a graph of connections I do `zt-g | dot -Tpng > g.png`.
