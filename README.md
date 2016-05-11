# Morton
  A Go library designed to implment N-dimensional coordinate encoding/decoding to/from Morton (Z-Order Curve) Order in Go.

## How Do I Use Morton?
  Please see the examples/ directory in this repository for a succinct and simple example of how to use this library.

## What Does This Library Do?
  This library encodes and decodes N-dimensional coordinates in the context of Z-Order Curve ordering.

## Morton Order/Z-Order Curve
  This type of encoding is used for ordering data, and is accomplished by interleaving the bits of each coordinate component. In 3 dimensions, {1,1,1} becomes 0b000000111, and {2,3,4} becomes 0b100011010, and so on. Here's another way to describe the 3 dimensional Morton encoded result: 0bzyxzyxzyxzyx, hence bit interleaving.

### Encoding
  The coordinate, {n1, n2, ...}, can be encoded into a single, memory saving, unsigned integer.

  Encoding an N-dimensional coordinate can have a compressive effect, with the exception of 1-dimensional coordinates.

  This library generates lookup tables for use in encoding. These lookup tables are concurrently generated.

### Decoding
  An encoded Z-value 0b000000001 can be decoded to reveal its N-dimensional coordinate components, which you may have guessed, is {1, n2, ..., n9}.

  This library generates magic bits for decoding. This process uses several bitwise operations on the encoded number, using these generated magic bits. Like the lookup tables, the magic bits are concurrently generated

## Whatever For?
  While there are many possible uses for Morton encoding, I originally wrote this as an adjunct to a voxelization project of mine (which I'll publish eventually).

  Without getting into specifics that are out of scope here, this library is being used in that particular project to encode each nodes position. This saves a bunch of memory, as it reduces the original 3 unsigned integer coordinate to a single unsigned integer.
