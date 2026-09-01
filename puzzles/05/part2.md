## Part 2 {#2}

The organiser turns back to the sheet and starts building the actual
lineup. He'll pick a subset of talks to run, subject to one rule: no
two picked talks may share a direct conflict - one such pair is
enough to derail the whole group into an argument. Among all valid
subsets, he wants the one that *maximises total hype*.

Formally, find the maximum total hype of a subset `S` of talks such
that no two members of `S` are connected by a conflict edge.

For the example in part 1, the best subset is:

* `factorio_deep_dive` with hype score of *60*,
* `bulanci_101` with hype score of *15*,
* `bulanci_hot_takes` with hype score of *40*.

Total = *60 + 15 + 40* = `115`.
