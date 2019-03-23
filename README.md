# Backoff

[![GoDoc][godoc image]][godoc]

This is a simple Go library implementing a ticker with exponential back off.
The API follows that of time.Ticker in the standard library, 
except with additional parameters for describing the backoff.

Example Usage:

```$xslt
start := time.Second
max := 120 * time.Second
factor := 2

ticker, err := backoff.NewTicker(start, max, factor)
if err != nil {
	// Handle error
}
defer ticker.Close()

for {
	select {
	case <-ticker.C:
		// Do something on the ticks
	}
}
```

[godoc]: https://godoc.org/github.com/robbieheywood/backoff
[godoc image]: https://godoc.org/github.com/cenkalti/backoff?status.png
