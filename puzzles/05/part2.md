## Part 2 {#2}

The organiser turns back to the sheet and starts building the actual
lineup. He'll pick a subset of talks to run, subject to one rule: no
two picked talks may share a *direct* conflict - one such pair is
enough to derail the whole group into an argument. A *direct* conflict
is a pair *explicitly listed* in the conflicts section of the input;
sitting in the same topic group is not enough on its own, since group
membership follows *transitively* along chains of conflicts. Among all
valid subsets, he wants the one that *maximises total hype*.

Formally, find the maximum total hype of a subset `S` of talks such
that no two members of `S` are connected by a conflict edge.

For the example in part 1, the direct conflicts are exactly the three
lines from the conflicts section:

* `factorio_deep_dive` and `factorio_101`,
* `factorio_101` and `bulanci_101`,
* `bulanci_hot_takes` and `sauna_field_notes`.

Note that `factorio_deep_dive` and `bulanci_101` land in the *same*
topic group (linked through `factorio_101`), but they are *not* in
direct conflict, so both can still be picked together.

The best subset is:

* `factorio_deep_dive` with hype score of *60*,
* `bulanci_101` with hype score of *15*,
* `bulanci_hot_takes` with hype score of *40*.

None of the three pairs within this subset appears in the conflicts
list, so it is valid. Total = *60 + 15 + 40* = `115`.
