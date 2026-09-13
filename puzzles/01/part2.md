## Part 2 {#2}

The coffee-topper isn't done. He sweeps the rest of the tiles out of the
bag and across the table: *"forget pairs - the whole pile has to balance."*
For each tile in the pile you can choose to

* *skip* it,
* *add* its value to a running total, or
* *subtract* its value from the running total.

How many distinct sign-and-selection combinations produce a running total of
exactly `7.00`?

Two combinations are considered distinct if any tile is treated differently
(skip / add / subtract).

For the example from part 1, the valid combinations are

* `+1.00 +2.00 +4.00` (skip both `3.50`s),
* `+1.00 +2.00 +4.00 +3.50 -3.50`,
* `+1.00 +2.00 +4.00 -3.50 +3.50`, and
* (skip `1.00`, `2.00`, `4.00`) `+3.50 +3.50`,

so the answer is `4`.
