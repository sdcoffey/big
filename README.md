# BIG

Big is a simple, immuatable wrapper around golang's built-in `*big.Float` type desinged to offer a more user-friendly API and immutability guarantees at the cost of some runtime performance. 

### Example

Usage is dead simple:
```go
dec := big.NewDecimal(1.24)
addend := big.NewDecimal(3.14)

dec.Add(addend).String() // prints "4.38"
```
