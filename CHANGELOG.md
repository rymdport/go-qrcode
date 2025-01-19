# Changelog

## v1.2.0
- Added `SetContent` method to allow reusing a QRCode object.
- Some internal lookup tables are now arrays to help the compiler know they won't grow.
- Some minor code cleanups. 

## v1.1.0
- Fixed potential crash in `WriteColorFile` if there is an error.
- Replaced `bytes.Buffer` with `strings.Builder` for slightly faster string concatenation.
- Use `string` instead of `[]byte` for encoding internally. Avoids allocating copies.  
- More minor code cleanups

## v1.0.0

Initial forked version. Contains mainly code cleanups and modernizations. 