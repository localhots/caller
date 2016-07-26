# Caller

Package caller is used to dynamically call functions with data unmarshalled
into the functions' first argument. Its main purpose is to hide common
unmarshalling code from each function implementation thus reducing
boilerplate and making package interaction code sexier.

[Documentation](https://godoc.org/github.com/localhots/shezmu/caller)

Caller abstracts away the process of unmarshaling data before processing.

```go
type PriceUpdate struct {
    Product string  `json:"product"`
    Amount  float32 `json:"amount"`
}

func PriceUpdatePrinter(p PriceUpdate) {
    log.Printf("Price for %q is now $%.2f", p.Product, p.Amount)
}

// Error handling is skipped for clarity
func main() {
    printer, _ := caller.New(PriceUpdatePrinter)
    _ = printer.Call([]byte(`{"product": "Paperclip", "amount": 0.01}`))
}
```
