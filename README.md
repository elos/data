Data [![GoDoc](https://godoc.org/github.com/elos/data?status.svg)](https://godoc.org/github.com/elos/data) [![Build Status](https://travis-ci.org/elos/data.svg?branch=master)](https://travis-ci.org/elos/data)
----

Data is a high-level package that defines the interfaces needed to effectively implement data stores for go applications.

A store is a very high-level interface which defines the framework for a an object relation layer. This is likely not your traditional conception of an ORM, as it is much more light-weight and we avoid reflection and type inferencing scenarios. This leads to specific, often large interfaces.  A store has a schema, essentially a relationship map with validity checking which exposes the Link and Unlink functions to manage the relationships between two "models" of data. The schema introduces the concept of a Linkable model.

Please see the godoc reference for more information.

Todo
----
make recorder db concurrent-safe

transactions
