# Changelog

## v1.1.0
- Fixed potential crash in `WriteColorFile` if there is an error.
- Replaced `bytes.Buffer` with `strings.Builder` for slightly faster string concatenation.
- Use `string` instead of `[]byte` for encoding internally. Avoids allocating copies.  
- More minor code cleanups

## v1.0.0

Initial forked version. Contains mainly code cleanups and modernizations. 