OlegDB: Enterprise Edition
==========================

One of the major concerns blocking enterprise adoption of OlegDB is the lack of
integration with existing enterprise workloads and toolsets. This project aims
to fix this, allowing anyone desiring a MAYO redundant parallel uptime galaxy-scale
supercluster to have OlegDB: Enterprise Edition as a Marginally Available option.

License is MIT. For consulting services, please contact me@christine.website.

```console
$ redis-cli -p 6660
127.0.0.1:6660> jar foo bar
OK
127.0.0.1:6660> unjar foo
bar
```

Usage
-----

### `JAR`

```
JAR <key> <value>
```

This is OlegDB's canonical 'set' function. Put a value into the mayo (the database). It's easy to piss in a bucket, it's not easy to piss in 19 jars.

### `UNJAR`

```
UNJAR <key>
```

This function retrieves a value from the database.

### `SCOOP`

```
SCOOP <key>
```

Removes an object from the database. Get that crap out of the mayo jar.

### `MEBBE`

```
MEBBE <prefix>
```

Return keys that match a given prefix.

### `DUMP`

```
DUMP
```

Like `MEBBE`, except that it takes no prefix and just dumps the entire tree.

### `SPOIL`

```
SPOIL <key> <duration>
```

Sets the expiration value of a key. Will fail if no value under the chosen key exists.

### `SNIFF`

```
SNIFF <key>
```

Retrieves the expiration time (RFC 3339) for a given key from the database.

### `SQUISH`

```
SQUISH
```

Compacts both the aol file (if enabled) and the values file. This is a blocking operation.

### `CANHAS`

```
CANHAS <key>
```

Returns whether the given key exists on the database.

### `UPTIME`

```
UPTIME
```

Gets the time, in seconds, that a database has been up.
