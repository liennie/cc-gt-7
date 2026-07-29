# AGENTS.md — repo guide

Repository: **`cc-gt-7`** — five "Code && Chill" puzzles for the
`garage::trip::7.0.0` event (unlock dates Sep 13–17 2026). Everything below
is the current state of the project; keep it in sync when you make changes.

## What "garage trip" is

**Garage Trip** is a recurring, non-profit annual retreat run by
attendees of the **GDG Garage Prague One** meetup group. The original
motivation: the weekly GDG Garage sessions never had enough time or
focus for the longer algorithmic tasks people wanted to tackle
together, so the group started renting a house for a week each year to
do just that — plus everything else the meetup didn't fit.

It is *not* a hackathon. The stated goals are "fun and relax",
supported by a rotating mix of:

- algorithmic tasks (solo or in teams) — the puzzle track this repo
  belongs to
- outdoor games (paintball, archery, roundnet, night games, …)
- LAN parties (Bulánci, Factorio, Warcraft 3, AoE II, the group's own
  RTS *Unnatural Worlds*, …)
- board games
- ciphering / logic games — puzzle hunts and escape games
- traditions the group keeps coming back to: hookah, BBQ, pizza night,
  invite ciphers, a barcamp of short talks

Editions to date (from `garage-trip.cz/history/`):

| version | year | dates              | attendees      | venue                                          |
| ------- | ---- | ------------------ | -------------- | ---------------------------------------------- |
| 1.0     | 2020 | 10.–13. 9.         | 17             | chatkyupotoka.cz                               |
| 2.0     | 2021 | 24.–29. 9.         | 15             | baba-jaga.cz                                   |
| 3.0     | 2022 | 2.– 10. 9.         | 22             | treninkcentrum.cz                              |
| 4.2     | 2023 | 6.– 14. 10.        | 15             | chalupasimia.cz                                |
| 5.0     | 2024 | 10.–17. 8.         | 24 (+2 kids)   | chalupasimia.cz (first English edition)        |
| 6.9     | 2025 | 20.–27. 9.         | 25 (+3 kids)   | Novy Svet (bunker visit, most hookahs ever)    |
| **7.0** | 2026 | **12.–19. 9.**     | tbd            | **Chalupa Nový Svět** (same as gt::6.9)        |

The versions are semver-style but only loosely (`4.2`, `6.9`) —
part of the joke. gt::7's invitation leans into it: *"Major version
update: from six to seven … 6 🤷 7. Legacy Support: All your favorite
friends and puzzles remain fully compatible. Multi-generational Mode:
the kid-friendly plugin is still enabled. Performance Tweak: optimized
for maximum nerdy joy and minimum sleep."*

### gt::7 specifics

- **Venue**: Chalupa Nový Svět, Prostřední Lipka 74, Králíky, Orlické
  hory. 32 beds / 14 rooms, ~14 000 m² garden with pond and outdoor
  playground, grill, 4 taps, sauna, two-story kids' corner. Same
  cottage as gt::6.9.
- **Trip window**: **Sat 12 Sep → Sat 19 Sep 2026** (8 days).
- **Preliminary schedule** (from the site — subject to change):
  | day       | main activity          |
  | --------- | ---------------------- |
  | Sat 12.9  | Arrival + Welcome BBQ  |
  | Sun 13.9  | Puzzle hunt            |
  | Mon 14.9  | Board games            |
  | Tue 15.9  | AI coding              |
  | Wed 16.9  | AI coding              |
  | Thu 17.9  | Barcamp                |
  | Fri 18.9  | Free day               |
  | Sat 19.9  | Leaving                |

The five puzzles in this repo are the **daily code-puzzle track** for
gt::7. `event.yaml` unlocks them at 10:00 local time from Sun 13.9
through Thu 17.9 — matching the five "core" days of the trip
(arrival Sat and free/leaving Sat get no puzzle). Puzzle 05 lands on
Thursday alongside the barcamp itself; the courtyard feast on Friday
is lore, not a puzzle beat.

The story flavour throughout the puzzles pulls directly from Garage
Trip lore: the Sunday-morning porch tile game (01), the shisha lounge
(02, echoing the gt::6.9 "shisha master and bartender"), the woodland
bunker (03, echoing the gt::6.9 WWII bunker visit), the AI-coding
workshop with its palm-sized GT-Nano robots on the lawn (04), and the
Thursday barcamp lineup on the cottage TV (05). The Sun-Thu markdown
is one coherent trip, not five unrelated prompts.

The author writes both the puzzles and the reference solutions before
the event; solutions live in `solutions/sNN/` and `TestEvent` locks the
committed `puzzle.yaml` answers to the generated inputs. During the
event, participants receive only the story markdown (`part1.md`,
`part2.md`) and their own generated `inputs/NN.txt`.

## Environment

- Go module: `module puzzles`, `go 1.26.0`.
- Sibling repo dependency: `replace github.com/liennie/code-and-chill => ../code-and-chill`
  (used by `event_test.go` via `github.com/liennie/code-and-chill/pkg/eventtest`).
- Transitive: `github.com/liennie/AdventOfCode` (used by s03 for
  `pkg/path.Shortest`).
- Randomness: every generator uses `math/rand/v2` seeded
  `NewPCG(0x0N_5eed_5eed, uint64(idx)+1)` where N is the puzzle number
  (1–5) and `idx` is the 0-based input index. Determinism is critical —
  regenerating must reproduce the current inputs byte-for-byte.

## Repository layout

```
event.yaml                    puzzle order + unlock dates
event_test.go                 wires all 5 solutions to eventtest
go.mod
README.md                     one line: "Code && Chill puzzles for gt::7"
cmd/generate/main.go          runs g01..g05, writes puzzles/NN/inputs/MM.txt
cmd/computeanswers/main.go    runs sNN.Solution on each generated input
generators/gNN/generator.go   Generate(idx int) []byte
puzzles/NN/
  puzzle.yaml                 name, per-input answers (part1, part2)
  part1.md, part2.md          puzzle story text
  inputs/01..05.txt           generated inputs (committed)
solutions/sNN/
  solution.go                 Solution(input []byte) []string returns
                              [part1, part2]
  solution_test.go            doc-example and edge-case unit tests
```

Standard commands:

- `go test ./...` — runs all solution tests + `TestEvent` (validates every
  generated input against `puzzle.yaml`).
- `go run ./cmd/generate` — regenerates every `puzzles/NN/inputs/*.txt`.
- `go run ./cmd/computeanswers` — prints each solver's outputs against the
  current inputs, useful after generator changes.

## Puzzle roster (story + mechanics)

All stories share a "Code && Chill" garage-trip framing. Each puzzle is
self-contained; solvers are pure `input → [part1, part2]`.

### 01 — "Lucky Sevens" (Sunday, cottage porch)

Fixed-point tile values, multiples of `0.25` in `[0.25, 6.75]`, one per
line. **Part 1**: count unordered pairs summing to `7.00`. **Part 2**:
count sign-and-selection combinations (skip / add / subtract each item)
that sum to `7.00`; empty selection excluded.

- Generator: 40 numbers, each `0.25*k` for `k∈[1..27]`. Max is `6.75` (never
  `7.00`) so `+v` cannot cancel `-v` to hit target.
- Solver: freq-map for Part 1; memoised DP `f(i, t)` over signed subsets
  for Part 2, with `suffAbs` pruning.

### 02 — "The Master's Blend" (Monday, shisha lounge) — most involved

A blend is a sequence of 2-char compounds (`Ap Vn Sm Ax` …). Rules
`LHS N> RHS` fire when LHS appears in the blend: it smolders in place for
`N` ticks, then atomically becomes RHS. A smoldering region can't
overlap another rule. **Part 1**: first tick any `Ax` appears × its
1-indexed position. **Part 2**: first tick where at least half of the
blend is `Ax`.

Design invariant (never surfaced in the puzzle text): the 5 starting
tokens seed 5 **independent groups**. Each group's rules form a closed
chain with **prime cycle length** from `{211,223,227,229,233,239,241,251,
257,263,269,271,277,281,283}` — Part 2 answer is `LCM(cycles) − 1`
(≈10¹¹–10¹²).

Generator (`generators/g02/generator.go`, ~830 lines):

- Each group starts `Fs > [c1..cn F]` (n∈1..3). Every catalyst is drawn
  globally unique; every intermediate rule's LHS contains a token freshly
  produced by the immediately preceding rule → exactly one rule fires per
  tick per group.
- 10–15 middle steps chosen from three excursion tiers: 30% `simple`
  (`C>C'`, `CF>C'F`, `CF>C'F'`), 25% short excursion (expand + walk +
  consolidate on size 2 or 3), 45% `longExcursion` — 4-way dispatch:
  `walk2`, `walk3`, `decay3`, `bubble`.
- Bubbling (`subgroupBubble` / `bubbleChain`): rewrites a strict subgroup
  (window size 1..n-1) of an expanded state to introduce a fresh
  catalyst, threading an anchor so consecutive bubbles must include the
  last rewritten position. `bubbleExcursion` dispatches 4-way over
  `bubbleMidShift`, `bubbleRightFence`, `bubbleWaitDecay`, `bubbleSameShape`.
- Ash producer/consumer: one rule with `Ax` in RHS, one rule with `Ax` in
  LHS producing `[Fs]` back. Ash consumer keeps smolder 1 (Ax lives one
  tick per cycle per group).
- Post-generation, `distributeSmolder` pads the total smolder so the
  observed cycle length equals the prime target `L_i`. **Rule-count =
  cycle-length invariant** must survive any generator change.
- Fence tokens: 15 pool; 5 distinct **starting** fences, intermediate
  rules draw from the remaining 10 so no intermediate can accidentally
  match another group's starting rule.

Solver (`solutions/s02/solution.go`, ~296 lines):

- `sim{blend, rules, commits}`. `commit{pos, length, rhs, completeAt}`.
- `matchRules(tick)`: scans uncommitted positions, collects **every**
  matching rule, sorts by `(pos asc, length desc)`, and **panics** on any
  overlap `matches[i].pos < matches[i-1].pos + matches[i-1].length`. This
  enforces the generator's one-rule-per-tick-per-group invariant.
  Non-overlapping matches are all committed at the same tick (needed for
  5 parallel groups).
- `applyCompletions(tick)`: sorts by pos, splices R-to-L, delta-adjusts
  positions of remaining commits to the right of each splice.
- `simulatePart2` uses per-token `cycleLength` + LCM (returns `lcm − 1`).
  Never runs the full LCM ticks.

Rule wire format: `strings.Fields(line)`; if `<3` fields, RHS is `nil`
(empty). RHS is written without trailing space when empty.

### 03 — "Bunker Sprint" (Tuesday)

101×101 grid with `#`, `.`, `S`, `E`. **Part 1**: shortest 4-directional
walk from `S` to `E` (returns `-1` if unreachable — never happens on
committed inputs). **Part 2**: same goal, but the racer plants up to 3
bombs anywhere on the grid before the run; each bomb turns the walls
in a `3 × 3` area around its center into walkable corridor, then the
walker paths the modified grid the same way as Part 1.

- Generator: recursive-backtracker maze on 50×50 cells seeded from centre;
  ~5% of interior walls knocked out for loops; S at `(1,1)`, E at
  `(99,99)`.
- Solver: BFS for Part 1; Part 2 is the shipped dive/warp graph (a bomb
  "dives" from `(r,c)` to any `(r+dr, c+dc)` with `dr,dc∈[-3..3]`,
  minus `(0,0)` and the four `(±3, ±3)` corners, at cost `|dr|+|dc|`).
  It is *not* the wall-destruction mechanic the puzzle text describes,
  but on every hand-crafted test up to 15×15 the two produce the same
  Part 2, pinned by `TestWallDestructionEquivalence` in `verify_test.go`
  (brute-force enumeration of every triple of in-bounds bomb centers,
  BFS on the modified grid, min over placements). See the "Puzzle 03"
  bullet under "Conventions and gotchas" for the caveat.

### 04 — "Boot Choreography" (Wednesday)

Tiny VM (in-story: the *GT-Nano*, a house-built SBC driving palm-sized
robots with LED strip + servos; the OUT stream is the robot's
choreography bus) with 6 signed-int64 registers (A–F), infinite 1-D
signed-int64 tape indexed by signed cursor `x`, and an output stream.
Ops: `LEFT/RIGHT`, `LOAD R/SAVE V`, `SET/ADD/SUB R V`, `JMP O/JNZ R O`,
`OUT V`. Jump offset `O` is signed nonzero and relative to the jump
instruction's own pc (`pc <- pc + O`). Non-jump instructions advance to
the next line. Halts when pc walks off the end.

Program is 6 equation-blocks, one per output line. Each block emits
`SET <out> b_i`, then for j = 0..5 (in a sweep direction that
alternates per row so no cursor-reset instructions are needed): `LOAD
<v>; ADD <v> off; SET <c> |a_ij|; ADD|SUB <out> <v>; SUB <c> 1;
JNZ <c> -2`, with `RIGHT` (even i) or `LEFT` (odd i) between columns.
`out_reg` is fixed for the whole program; `var_reg` and `ctr_reg` are
picked per (i, j) from the 5 non-out registers.

**Part 1**: run untouched (all zeros) and report the positional base-10
digest `Σ v_i · 10^i` of the 6 outputs.
**Part 2**: recover the six signed integers to pre-populate cells `0..5`
so every output is `0`; report `Σ x_j · 10^j`.

- Generator: `n = 6`. Reject-samples `x*_j ∈ [-100,100]\{0}`,
  `A[i][j] ∈ [-100,100]\{0}` and `offset[i][j] ∈ [-100,100]` until A is
  invertible over the rationals (big.Rat Gauss) and every
  `-Σ A[i][j]·x*_j` is non-zero. `x*` and `A` draw magnitudes as the
  max of two uniform `[1, 100]` picks (biased toward the extremes) with
  a fair sign; `offset` is uniform. Emitted immediate is
  `b_i = -Σ A[i][j] · (x*_j + offset[i][j])`, so with mem = 0 Part 1
  outputs are non-trivial and with mem = x* they collapse to zero.
- Solver: `runVM` with an optional `patches map[int64]int64` initial
  memory. One untouched run yields `b` (Part 1 outputs); 6 unit-vector
  probes `patches = {j: 1}` yield `A[:,j] = probe - b`. Then
  `solveLinear` runs big.Rat Gauss on `[A | -b]` to recover x*. A final
  patched run with `mem = x*` asserts every output is 0. Part 2 digest
  is `Σ x*_j · 10^j`.
- Text NEVER mentions systems, matrices, linearity, ranges, or the
  block/sweep structure. Those are pure generator secrets.

### 05 — "Barcamp Lineup" (Thursday, cottage TV)

Three blank-line sections: HYPE (`slug int`), CONFLICTS (`slug slug`),
QUESTIONS (`slug slug`). Slugs are unique `[a-z0-9_]+` identifiers.
Conflict graph is a **forest**. **Part 1**: for each question whose
two talks fall in *different* connected components, add the sum of
their hype scores; return the total. **Part 2**: max-weight
independent set across the forest.

- Generator: 80 talks with GT-lore slug titles composed from three
  package-level pools (`topics` 47, `suffixes` 21, `prefixes` 12);
  ~75% follow `topic_suffix` (e.g. `factorio_deep_dive`), ~25% follow
  `prefix_topic` (e.g. `weekend_shisha`). Every topic carries a
  category (`game` / `indev` / `outdoor` / `food` / `place` / `event`)
  and every fragment declares which categories it may combine with;
  rejection sampling filters out silly combos like `bbq_speedrun`,
  `unnatural_worlds_postmortem`, or `sauna_considered_harmful`. Titles
  are unique; the title-to-node mapping is shuffled so slug order
  doesn't leak tree grouping. Hype scores uniform in `[1, 1000]`.
  Conflict forest is *explicit*: exactly one big tree of `48..54`
  talks plus 3-5 small trees of `4..10` each (every tree size >= 4,
  guaranteed by absorbing any <4 tail into the last small tree). Each
  tree is a random recursive tree over its contiguous node block. 20
  queries mixed between same- and cross-tree pairs (`nSame` random in
  `[6, 14]` so Part 1 answers vary).
- Solver: union-find for Part 1; per-tree post-order DP
  `take[v] = w[v] + Σ skip[c]`, `skip[v] = Σ max(take[c], skip[c])` for
  Part 2. Solver is name-agnostic — the slug reformat didn't require
  any code change.

## Committed answers (`puzzles/NN/puzzle.yaml`)

| # | 01               | 02               | 03           | 04                                | 05        |
|---|------------------|------------------|--------------|-----------------------------------|-----------|
| 1 | 31 / 56842165683797805 | 1890 / 723491704828 | 865 / 265 | -557775349 / 2791931       | 8697 / 25746 |
| 2 | 31 / 50543652303855546 | 1890 / 631477939050 | 701 / 277 | 988403367 / 5926868        | 11430 / 26597 |
| 3 | 39 / 58569109103239985 | 452 / 944235325140 | 521 / 313 | 1033228639 / -7986246       | 7210 / 27155 |
| 4 | 32 / 55960769019965964 | 2856 / 1077260446036 | 789 / 257 | -1598083723 / -9450015    | 8192 / 27832 |
| 5 | 33 / 51855039244186640 | 1332 / 799098816148 | 697 / 285 | 408421131 / 7824336        | 8217 / 30313 |

## Conventions and gotchas

- **CRLF normalisation**: all solvers do
  `strings.ReplaceAll(string(input), "\r\n", "\n")` before parsing.
- **Answer stability**: any generator change must be validated with
  `go test ./...` (which runs `TestEvent`). If answers shift, update
  `puzzle.yaml` in the same commit.
- **Puzzle 02 invariant**: rules must never overlap in a tick. The solver
  panics on violation with
  `s02: overlapping rule matches at tick T: [pos p1 len l1] and [pos p2 len l2]`.
  Do not silence this panic — a trigger means the generator broke the
  chain-uniqueness property.
- **Puzzle 02 rule shapes** on input 01: ~42 `1>1`, 48 `2>2`, 5 `3>3`,
  12 `1>0` decays. Empty RHS emitted without trailing space.
- **Puzzle 03 dive/wall-destruction split**: the puzzle text tells the
  player that bombs pre-detonate `3 × 3` wall craters and then the
  sprint runs by Part 1's rules. The shipped solver actually runs a
  dive/warp graph (`|dr|, |dc| <= 3`, minus origin and the four
  `(±3, ±3)` corners, cost `|dr|+|dc|`, up to 3 uses); the two mechanics
  produce the same Part 2 numbers on every committed input, verified by
  `TestWallDestructionEquivalence` in `solutions/s03/verify_test.go`
  (exhaustive brute-force enumeration of every triple of in-bounds bomb
  centers on small hand-crafted grids up to `15 × 15`). The real
  `101 × 101` inputs are too big to brute-force; keep the equivalence
  test in mind before touching either the solver or `puzzle.yaml`.
- **Solver privacy**: the puzzle text never mentions groups, primes,
  LCM, cycle length, block independence, forest structure, etc. Those
  are pure generator secrets that make Part 2 tractable.
- **Puzzle markdown formatting**: `part1.md` / `part2.md` are rendered
  with a markdown extension that allows inline markdown formatting
  *inside* code spans and code blocks. A formatting-enabled code span is
  delimited by **double backticks** (``); a formatting-enabled
  code block is delimited by **four backticks** (````) on their
  own lines. Escaping inside these uses backslashes, so a literal
  backtick is written as `` \` ``. Use this to emphasise important
  tokens inside code samples (e.g. highlighting one rule in a listing).
  Prefer a single `*` for emphasis over `**` — keep the source lighter
  and reserve `**` for cases where single-`*` would collide with
  surrounding punctuation.
- **ASCII in puzzle text**: prefer plain ASCII in `part1.md` / `part2.md`.
  Use `^v<>` for directional arrows (not `↑↓←→`), `-` for dashes (not
  `–`/`—`), `'`/`"` for quotes (not `‘’“”`), `...` for ellipsis, `*` for
  bullets, etc. Non-ASCII is fine when it's the actual subject (e.g. a
  name that's genuinely spelled that way) or a math symbol that has no
  clean ASCII equivalent (`×`, `÷`, `≤`, `≥`, `≠`, Greek letters, …).
- **Do not add doc comments / helpers / abstractions** unless requested.
  Keep changes minimal and targeted.
- **No markdown docs** should be created without an explicit request.
  This file exists because the user asked for it.
