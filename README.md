sharedmap
=========

An embarrassingly simple concurrent map implementation in Go, based on synchronization through channel I/O instead of mutexes.

A benchmark is included in the test file. On my laptop (darwin/amd64 on 2.4GHz Core i7 with 1333 MHz DDR3 RAM) I get the following results:

    $ go test -bench . -benchtime 10s -cpu 1
    PASS
    BenchmarkSharedMapSingle    20000000          1199 ns/op
    BenchmarkSharedMap100     200000         99812 ns/op
    BenchmarkSharedMap100000         100     164587644 ns/op
    BenchmarkMutexMapSingle 100000000          239 ns/op
    BenchmarkMutexMap100     1000000         23957 ns/op
    BenchmarkMutexMap100000     1000      24149722 ns/op
    ok      github.com/codemartial/sharedmap    151.693s

With 1 concurrent accessor, SharedMap is about 5x slower than a Mutex protected Map. This slowdown increases to nearly 7x as we hit 100k concurrents. The above test run used only 1 CPU core. However, if we force the Go runtime to use all 4 available cores, the results are quite different:

    $ go test -bench . -benchtime 10s -cpu 4
    PASS
    BenchmarkSharedMapSingle-4  10000000          2456 ns/op
    BenchmarkSharedMap100-4   200000         93664 ns/op
    BenchmarkSharedMap100000-4       200     116053928 ns/op
    BenchmarkMutexMapSingle-4   100000000          241 ns/op
    BenchmarkMutexMap100-4    500000         60585 ns/op
    BenchmarkMutexMap100000-4        500      61700457 ns/op
    ok      github.com/codemartial/sharedmap    169.553s

In this scenario, the single accessor test becomes considerably slow but for higher concurrents, SharedMap performs better than it did with just 1 CPU core, while mutex protected map performs slower than it did earlier. At 100k concurrents, the performance gap between SharedMap and Mutex protected map reduces to 1.8x (SharedMap is still slower).

The benchmark assumes a read:write ratio of 4:1 and add:delete ratio of 100:1.

Conclusion
----------

Mutexes perform better than channel synchronised access to shared state. However, the performance gap diminishes where access concurrency is high under a multi-core usage scenario. This is a performance vs. verifiability trade-off. It's much easier to wrap shared state in a goroutine with channel based access and to argue about its correctness than to ensure that all mutex based access is done correctly. In case of performance issues, though, the latter would be preferable.

