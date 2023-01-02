# GO-WHERE

A basic SQL Query builder where you type out your entire query before you hit the `WHERE` clause

The idea is that when a user providers query params via REST, those query params are used to create a `WHERE x AND y AND z` statement.

These params are usually provided in a struct by your router/framework of choice, so they can just be passed into the query constructor.

OFFSET and LIMIT and reserved keywords, so if they are passed into the query constructor, they will be applied at the end of the query and LIMIT and OFFSET, one can also not live without the other.