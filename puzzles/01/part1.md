# 01: Lucky Sevens {#1}

Sunday morning. The BBQ ash is still warm from last night, coffee's brewing
on the cottage porch, and someone has upended a cloth bag of numbered wooden
tiles across the table - a battered set dug out of the cottage's
game shelf, each tile stamped with a small fractional value. One of the
early risers waves a mug and calls the first side game of the day: *"Every
pair of tiles that adds up to exactly 7.00 - point them out and I'll top up
your coffee."*

## Input

Your puzzle input is a list of fixed-point numbers with exactly two decimal
places, one per line. Every value is a multiple of `0.25` and lies in
`[0.25, 6.75]`.

## Part 1

Find the number of *unordered pairs* of tiles whose values add up to
exactly `7.00`.

### Example

```
1.00
2.00
4.00
3.50
3.50
```

In the example, the only lucky pair is `3.50 + 3.50`, so the answer is
`1`.
