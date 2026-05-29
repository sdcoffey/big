# BIG

Big is a simple, immutable wrapper around Go's built-in `*big.Float` type designed to offer a more user-friendly API and immutability guarantees at the cost of some runtime performance.

Because Big wraps `*big.Float`, it uses arbitrary-precision binary floating-point arithmetic. It is not a fixed-point decimal or money library. If you need decimal-exact financial arithmetic, use a decimal package with explicit scale and rounding semantics.

### Example

Usage is dead simple:
```go
dec := big.NewDecimal(1.24)
addend := big.NewDecimal(3.14)

dec.Add(addend).String() // prints "4.38"
```
