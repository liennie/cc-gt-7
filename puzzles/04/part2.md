## Part 2 {#2}

The master is patient - one stray `Ax` is not the end of the pipe. What
ruins the blend is when the ash takes over.

Find the *first tick at which at least half of the compounds in the blend
are `Ax`*.

Continuing the example trace from part 1, counting ash as a fraction of the blend:

| tick | blend         | ash count | fraction |
| ---: | :------------ | :-------- | :------- |
| 6    | ``Hz*Ax*Cf``  | 1         | 1/3      |
| ...  |               |           |          |
| 9    | ``HzAx*Ax*``  | 2         | 2/3      |
| ...  |               |           |          |
| 12   | ``*Ax*AxAx``  | 3         | 3/3      |

Ash first reaches half at tick *9*, when ``Cf 7> *Ax*`` completes, so the answer is `9`.
