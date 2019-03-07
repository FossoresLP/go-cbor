CBOR encoder and decoder for Golang
===================================

**Important information:**
- This package still needs additional tests
- Decoding into a user supplied type is not yet possible. The decoder will return an interface you will have to try to convert yourself. This feature is WIP and will be added shortly.
- This is a go module, not a normal package. The master branch is the main development branch and may therefore receive breaking changes. Since all supported go versions now support modules this shoud not be an issue.

This go module implements a general purpose CBOR encoder and decoder as defined in [RFC7049](https://tools.ietf.org/html/rfc7049). There are a few notable departures from the recommendations made there:

- The decoder will ignore unknown tags and won't forward them to the caller
- The decoder will fail when receiving unknown simple values

Encoder
-------

When encountering custom types the encoder will first try to call `EncodeCBOR() []byte`, then try `MarshalCBOR() ([]byte, error)` and panic when an error is returned. In case no specific functions for CBOR are present, the encoder will try `MarshalBinary() ([]byte, error)` with the same behavior on errors and encode it as a byte string. Then the normal type-based encoding will attempt to determine the type of the input and encode it accordingly. Should that not be possible for some reason, the encoder will try `MarshalText() ([]byte, error)` as a last resort and encode the output as a string. In case that function is not available either, the encoder will fall trough and panic.

The types supported by the type-based encoding are:

- All standard signed and unsigned integers (BigNums are WIP)
- All standard floating point value (BigFloats are WIP)
- Strings, arrays and slices of supported types (will always be encoded with fixed length)
- Maps of supported types (will always be encoded with fixed length)
- Structs (only exported fields, cbor tag will be checked for name, will be encoded as indefinite length maps)
- Booleans

Decoder
-------

The decoder supports all types defined in the CBOR standard including 16-bit floats which will be decoded as 64-bit floats.
All individual decoding functions are exported by this package and do return values with the most solid typing possible. Please not that these functions do not check the major type of the input.

License
-------

This module is released under the Boost Software License 1.0 which can be found in `LICENSE.md` in the repository. This means that you can freely use the code according to the license and do not have to give attribution for binary distributions although this is always highly appreciated.