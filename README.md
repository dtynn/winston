#### Winston

This project comes from the paper [Gorilla: A Fast, Scalable, In-Memory Time Series Database](http://www.vldb.org/pvldb/vol8/p1816-teller.pdf) .

The name comes from Overwatch.



##### benchmark

```
BenchmarkChunkPush-8             5000000           279 ns/op          14 B/op          0 allocs/op
BenchmarkChunkIterRead-8        10000000           166 ns/op           0 B/op          0 allocs/op
BenchmarkChunkIterRead1K-8         10000        172933 ns/op           0 B/op          0 allocs/op
BenchmarkChunkIterRead10M-8            1    1496245569 ns/op           0 B/op          0 allocs/op
```



##### license

MIT
