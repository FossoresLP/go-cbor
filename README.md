CBOR encoder and decoder for Golang
===================================

**Important information:**
- This package still needs additional tests
- Decoding into a user supplied type is not yet possible. The decoder will return an interface you can convert yourself. This feature is WIP and will be added shortly.
- This is a go module, not a normal package. The master branch is the main development branch and may therefore receive breaking changes. Since all supported go versions now support modules this shoud not be an issue.

This go module implements a general purpose CBOR encoder and decoder as defined in [RFC7049](https://tools.ietf.org/html/rfc7049). There are a few notable departures from the recommendations made there:

- The decoder will ignore unknown tags and won't forward them to the caller
- The decoder will fail when receiving unknown simple values

How to use a go module
----------------------

You can use modules just the way packages work, in fact they only wrap around packages and provide versioning support. To use versions higher than v1 simply add the major version to the end of the import path `import github.com/FossoresLP/go-cbor/v2`. Please note that this repository does not yet provide such versions.

Encoder
-------

When encountering custom types the encoder will first try to call `EncodeCBOR() []byte`, then `MarshalCBOR() ([]byte, error)` and panic when an error is returned. In case no specific functions for CBOR are present, the encoder will try `MarshalBinary() ([]byte, error)` with the same behavior on errors and encode the result as a byte string. If that's not supported either, the normal type-based encoding will attempt to determine the type of the input and encode it accordingly. As a last resort, the encoder will try `MarshalText() ([]byte, error)` and encode the output as a string. In case that function is not available, the encoder will panic. This is considered acceptable in the encoder as you should be in control of the input types.

The types supported by the type-based encoding are:

- Booleans
- All standard signed and unsigned integers (BigNums are WIP)
- All standard floating point values (BigFloats are WIP)
- Strings, arrays and slices of supported types (will always be encoded with fixed length)
- Maps of supported types (will always be encoded with fixed length)
- Structs (only exported fields, cbor tag can declare custom name, will be encoded as indefinite length maps) (omitempty and other properties from the JSON package may be added later)

Decoder
-------

The decoder supports all types defined in the CBOR standard including 16-bit floats which will be decoded as 64-bit floats. Undefined and null are both decoded as nil.

All individual decoding functions are exported by this package and do return values with the most solid typing possible. Please note that these functions do not check the major type of the input.

License
-------

This module is released under the Boost Software License 1.0 which can be found in `LICENSE.md`. This allows you to use the code freely according to the license terms while not having to give attribution in binary distributions although this is always highly appreciated.
