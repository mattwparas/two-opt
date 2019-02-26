# two-opt
2-Opt TSP solver implementation in Go

### How it works

See here for an explanation of the 2-opt algorithm: [2-Opt](https://en.wikipedia.org/wiki/2-opt)

It is by no means the most complex or sophisticated TSP solver, however it is relatively fast to implement. It does not scale well since it's time complexity is O(n^2), but we can abuse Go's concurrency tools in order to solve many instances simultaneously. The two-opt algorithm is deterministic from a given starting position, however you can re-solve it many times with different randomly shuffled starting tours and select the best result. 

In this implementation (currently just using randomized distance matrices and an arbitrary number of nodes), you can select how many times you would like to randomly restart, and the solver runs them all concurrently. On my machine, a TSP with 200 nodes can be solved 200 different times with random restarts in about 2.5 seconds. From there, the best solution can be selected.
