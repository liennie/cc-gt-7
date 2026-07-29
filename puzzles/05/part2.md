## Part 2 {#2}

The organiser turns back to the sheet and starts building the actual
lineup. He'll pick a subset of talks to run, subject to one rule: no
two picked talks may share a direct conflict - one such pair is
enough to derail the whole group into an argument. Among all valid
subsets, he wants the one that *maximises total hype*.

Formally, find the maximum total hype of a subset `S` of talks such
that no two members of `S` are connected by a conflict edge.

For the example in part 1, treat each group independently:

* From `factorio_deep_dive - factorio_101 - bulanci_101`: taking
  `factorio_101` alone gives `25`; taking `factorio_deep_dive` plus
  `bulanci_101` gives `60 + 15`. Best = `75`.
* From `bulanci_hot_takes - sauna_field_notes`: taking
  `bulanci_hot_takes` alone gives `40`; taking `sauna_field_notes`
  alone gives `30`. Best = `40`.

Total = `75 + 40 = 115`.
