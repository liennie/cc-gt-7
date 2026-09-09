# 02: Bunker Sprint {#1}

Monday. Over breakfast someone floats a plan for the afternoon: a
*bunker sprint* through the old concrete WWII bunker in the woods a couple
of kilometres north. The idea is exactly what it sounds like - line up at
the front door, tear through the interior corridors on foot, and clock
out the moment you punch the emergency exit on the far side. Fastest
time wins.

Nobody wants to run it blind, so one of the more thorough souls surveyed
the place earlier in the week and produced a tidy grid map. It is a
`101 × 101` square, one cell per bootstep, with the front door marked
`S` and the emergency exit on the far side marked `E`. The outer wall is
solid rock - you can only travel through the interior cells.

## Input

Your puzzle input is a grid of `101` lines, each `101` characters wide,
drawn with the following cell tiles:

```
#   solid rock or rubble
.   open corridor
S   entrance - start here
E   emergency exit - finish here
```

Movement is one step at a time, orthogonally (up / down / left / right),
only onto `.`, `S`, or `E` cells. You may not step onto walls or leave
the grid.

## Part 1

Find the length of the *shortest walk* from `S` to `E`, counting each
step as one.

### Example

Consider a `15 × 15` grid with an outer wall of solid rock, `S` on the west
side and `E` tucked into the south-east corner:

````
###############
*S*...#.........#
#.###.#######.#
#.....#.....#.#
#.#####.###.#.#
#.#.....#...#.#
###.#.###.#.#.#
#...#...#.#...#
#.#.###.#.###.#
#.#.....#...#.#
#.###.###.#.#.#
#.#...#...#...#
#.#.###.#######
#............*E*#
###############
````

Narrow corridors, dead ends and a couple of small loops fill the
interior. One of the shortest walks from `S` to `E` weaves through the
maze like this, with each step marked by an arrow pointing where the
racer moves next:

````
###############
*Sv*..#*>>>>>>>>v*#
#*v*###*^*#######*v*#
#*>>>>^*#.....#*v*#
#.#####.###.#*v*#
#.#.....#...#*v*#
###.#.###.#.#*v*#
#...#...#.#..*v*#
#.#.###.#.###*v*#
#.#.....#*v<<*#*v*#
#.###.###*v*#*^*#*v*#
#.#...#*v<<*#*^<<*#
#.#.###*v*#######
#......*>>>>>>E*#
###############
````

The shortest walk here is `45` steps.
