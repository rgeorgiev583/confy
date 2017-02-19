# Confy

***Confy*** is a web API service for managing the configuration of a remote system
using a convenient tree-based representation.

*Confy* uses internally the API of the *Augeas* project in order to parse the
configuration files present in the system and to lookup and modify the
configuration tree.


## Augeas

Please consult the
[Augeas project documentation](http://augeas.net/docs/index.html) if you want to
understand how any of the following things work:
* the Augeas public API;
* the Augeas tree representation;
* Augeas path expressions;
* Augeas configuration file format descriptions.


## RESTful API

*Confy* provides a RESTful JSON-based API which supports the following operations
with the configuration tree (where `$PATH` is the path to the respective tree node
which is being manipulated):

* `GET /api/list/$PATH`:  retrieves the list of paths to the direct children of the
Augeas node with the given path as a JSON array.
* `GET /api/match/$PATTERN`:  retrieves the list of paths to the Augeas nodes
matching the given pattern string as a JSON array.
The pattern string's special syntax is as follows: the `*` character represents any
character in a label, and `[i]` represents the i-th element of a node array.
* `GET /api/get/$PATH`:  retrieves the value of the Augeas node with the given path
as a JSON object with one property: `value`.
* `GET /api/all/$PATH`:  retrieves the values of all Augeas nodes with the given
path (i.e. all nodes from the array with that path) as a JSON array.
* `GET /api/label/$PATH`:  retrieves the label (the last component of the path) of
the Augeas node with the given path.
* `PUT /api/set/$PATH`:  assigns the value passed as the JSON `value` property in
the request string to the Augeas node with the given path.
Its parent nodes are created as necessary if they do not exist.
* `PUT /api/multiset/$PATH`:  assigns the value passed as the JSON `value` property
in the request string to all nodes in the Augeas node array with the given path
whose relative paths to the array match the pattern passed as the JSON `pattern`
property.
If the pattern is an empty string, all nodes in the array are matched.
* `PUT /api/clear/$PATH`:  clears (i.e. sets to NULL) the value of the Augeas node
with the given path.
Its parent nodes are created as necessary if they do not exist.
* `POST /api/insert-before/$PATH`:  inserts a new Augeas node with a label passed
as the JSON `label` property in the request string before the node with the given
path.
* `POST /api/insert-after/$PATH`:  inserts a new Augeas node with a label passed
as the JSON `label` property in the request string after the node with the given
path.
* `DELETE /api/remove/$PATH`:  removes the whole Augeas subtree with the given
path, including all of its descendants.
* `PATCH /api/move/$PATH`:  moves the whole Augeas subtree with the path passed as
the JSON `source` property in the request string to the path passed as the JSON
`destination` property.
* `PATCH /api/reload/$PATH`:  reloads the Augeas configuration tree from the mapped
configuration files.
* `PATCH /api/save/$PATH`:  persists the Augeas configuration tree into the mapped
configuration files.