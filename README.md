Data [![GoDoc](https://godoc.org/github.com/elos/data?status.svg)](https://godoc.org/github.com/elos/data)[![Build Status](https://travis-ci.org/elos/data.svg?branch=master)](https://travis-ci.org/elos/data)
----

Data is a high-level package that defines the interfaces needed to effectively implement data stores for go applications.

Please see the godoc reference for more information.

-----------------------------------

Schema is a semi-high-level package that defines the framework for defining an object relation layer for go system applications. Note: this is likely not your traditional concept of an "ORM" persay. It is a much more light-weight.

Essentially a schema is a relationship map with validity checking which exposes the Link and Unlink functions to link two "models" of data.

As such, schema also introduces the concept of a Linkable model. It extends the notion of a [data](https://github.com/elos/data) Record.
