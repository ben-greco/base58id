# base58id

**Fast generation of surprisingly short unique IDs for standalone or distributed systems with low memory use 
and no disk use**

## Quick Start (Single instance only)

```go
s, _ := base58id.New(1)

id := s.Get()

fmt.Println(id)
// KRuRjLq
```

## Quick Start (Multi-instance)

```go
// Make sure to use unique integer IDs for multiple instances
a, _ := base58id.New(1, 1)
b, _ := base58id.New(1, 2)

// These two ids will never be the same
idOne := a.Get()
idTwo := b.Get()
```

## Installation

`go get -u github.com/ben-greco/base58id`

## What base58id is

This package...
- Makes string IDs that are between 59% and 81% shorter than a type 4 UUID, depending on 
the number of IDs generated per second and number of concurrent generators.
- Makes IDs that are guaranteed to be unique for **at least** a millennium when instance IDs for multi-instance configurations are used correctly.
- Generates IDs faster than 180,000 per second per instance.
- Can pre-generate and store up to `math.MaxInt32` unique IDs in memory to be retrieved during 
burst usage, depending on configuration. See parameters section. 
- Never uses disk storage for unique ID generation.
- Uses a small amount of memory for unique ID generation per instance.

## What base58id is NOT

This package does **NOT**...
- Make it difficult to guess what the next id will be.
- Serve as a drop in replacement for type 4 UUID strings in multi-instance systems without some configuration first.
- Make unique IDs across many instances and/or distributed systems **UNLESS** you properly use a unique 
integer ID for each instance/system creating IDs.

## The parameters

- capacity (int): This is how many pre-generated IDs will be stored in memory. Defaults to 1 if 
less than 1 is given. Set to a larger number to store up to `math.MaxInt32` unique IDs in memory 
depending on your tradeoffs between storing in memory and your expectation that the application 
will be hampered by the on-demand burst rate of 180k unique IDs per second. Increase instances 
to increase on-demand throughput. See instanceID.

- instanceID (...int): Must not contain a zero (see How it works). Must either not be submitted, 
or be a single integer. Use this to make multiple instances of base58id generators and increase 
output. **NEVER** create a multi-instance system where two instances have the same instance ID 
(that includes having no instance ID, 
they are guaranteed to make a duplicate ID under any kind of production volume. Overall, shorter 
instance IDs create shorter base58id IDs and omitting it creates the shortest IDs.

## Methods

- `New(capacity int, instanceID ...int)` : returns a base58id generator according to the above parameters

- `Get()` : returns a new unique ID as a string. 

- `GetMany(n int)` : returns n new unique IDs as a slice of strings.

## How it works

A base58id ID is created by combining the following things together in an integer and encoding 
that integer in base58:

- an "as short we can" integer that is unique in the last second that contains no zeros
- a zero, as a separator
- an instance ID that contains no zeros, but this may be omitted
- a zero, as a separator, if needed
- the current unix timestamp in seconds

### Examples

Here is an example for a theoretical instance with instance ID `123` that was created during a
relatively high volume of ID generation (hence the relatively long `574385` as the unique integer)

```go
"574385" + "0" + "123" + "0" + "1582317552"
```
Which becomes the following 128 bit integer
```
574385012301582317552
```
Which is base58 encoded into
```
PzHZVxm1wDxT
```

However for less strenuous use on a single instance, quite short unique IDs can be generated
 as you can see in this example that follows the same process:
```go
"3" + "0" + "1582317552"
```
```
301582317552
```
```
8vUoD3D
```

### End of life for base58id

Because of the non-zero restrictions and the zeroes that are used as separators, no IDs can be 
duplicated for different instance IDs until an ID is created after a Unix time that is twelve 
digits long with a zero for the second digit, the first possible time is 100000000000 which
 corresponds to 16 November 5138.


## Changing your configuration

You can switch back and forth between omitted instance IDs, multi-instance and single instance at 
any time as long as there are no conflicts in unique instance IDs for a multi-instance system or there 
is at least one full second where no IDs are generated. No two base58id IDs can ever be duplicates 
unless they are generated within one second of eachother or end of life has passed. 

