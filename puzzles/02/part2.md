## Part 2 {#2}

The first heat was a warm-up. Someone digs the actual demolition
charges out of the tool shed - the surveyor left them there after the
mapping run - and hands each racer *three bombs* to plant on the map
before the sprint starts.

A bomb may be planted on any cell inside the grid, wall or corridor.
All of them detonate before the run; each blast turns the walls in a
`3 × 3` area centered on the bomb's cell into walkable corridor. Once
the dust settles the sprint plays by the same rules as part 1.

You don't have to use all three bombs; leaving one or two in your
pocket is fine if you don't need them.

Same goal: find the *shortest path* from `S` to `E`.

In the example from part 1, two bombs planted at *row `2`, column `5`*
and *row `13`, column `14`* punch a hole through the wall blocking row `1` right
next to `S` and the wall band shielding `E`, so the racer can march
straight along row `1` and then straight down column `13`. Below, the
`@` cells are the two `3 × 3` craters and the arrows trace the
racer's route:

````
###*@@@*#########
*S>>>>>>>>>>>>v*#
#.#*@@@*#######*v*#
#.....#.....#*v*#
#.#####.###.#*v*#
#.#.....#...#*v*#
###.#.###.#.#*v*#
#...#...#.#..*v*#
#.#.###.#.###*v*#
#.#.....#...#*v*#
#.###.###.#.#*v*#
#.#...#...#.*@v@*
#.#.###.####*@v@*
#...........*@E@*
###############
````

The shortest walk in the modified grid is `25` steps. The third bomb
stays in the racer's pocket.
