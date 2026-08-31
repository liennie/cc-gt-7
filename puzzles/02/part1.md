# 02: The Master's Blend {#1}

Monday evening. The house shisha master has spent the afternoon in the
kitchen preparing a lineup of tobacco *blends* for the night, each one a
careful mix of his favourite flavours packed into its own bowl. You grab
the first bowl, set an HMS on top, drop the coals into the HMS, and settle
in. Now you just have to hope the master got the chemistry right - one
wisp of *ash* too early and the whole pipe is spoiled.

## Compounds and blends

Every tobacco mix the master prepares is really a *blend* of flavour
*compounds* packed one after another into the bowl.

A *compound* is a two-letter token: an uppercase letter followed by a
lowercase letter (e.g. `Ap`, `Vn`, `Sm`, `Ax`). A *blend* is a sequence
of compounds written back-to-back with no separators, e.g. `ApVnSmAx` is
four compounds: `Ap`, `Vn`, `Sm`, `Ax`.

The *ash* compound is always `Ax`. Every other compound in your input is
one of many flavours the master keeps in his cabinet.

## Rules

Each *rule* describes how a run of adjacent compounds transforms under
heat. Rules are written on one line each in the form

````
*LHS* *N*> *RHS*
````

where

* `LHS` is one or more compounds concatenated (the pattern to match).
* `RHS` is zero or more compounds concatenated (the pattern produced).
* `N` is a positive integer - the *smolder time* in ticks.

At every tick, every occurrence of a rule's `LHS` in the blend causes that
rule to fire: the matched compounds *smolder in place, unchanged, for
exactly `N` ticks* and then, in a single atomic step, are replaced by the
rule's `RHS`.

A compound that is smoldering as part of one rule cannot simultaneously
participate in another rule. The rules in your input are designed so that
at every tick no two rule matches overlap.

Positions in the blend are numbered starting at `1` from the left.

## Input

The first line of your puzzle input is the initial blend at tick `0`.
Then a blank line. Then one rule per line, in no particular order.

## Part 1

Simulate the blend under coal. Find the *first tick at which any `Ax`
compound is present in the blend*. Return the tick number multiplied by
the `1`-indexed position of that `Ax` in the blend at that tick.

If more than one `Ax` appears in the blend on that tick, use the position
of the leftmost one.

### Example

```
ApVn

Ap 3> HzSm
Vn 2> LcCf
SmLc 3> Ax
Hz 9> Ax
Cf 7> Ax
```

Tracing tick by tick:

| tick | blend          | notes                                         |
| ---: | :-----------   | :-------------------------------------------- |
| 0    | `ApVn`         | `Vn 2>` and `Ap 3>` both fire                 |
| 1    | `ApVn`         | both still smoldering                         |
| 2    | ``Ap*LcCf*``   | ``Vn > *LcCf*``; `Cf 7>` fires                |
| 3    | ``*HzSm*LcCf`` | ``Ap > *HzSm*``; `Hz 9>` and `SmLc 3>` fire   |
| ...  |                | everything is smoldering                      |
| 6    | ``Hz*Ax*Cf``   | ``SmLc > *Ax*`` (first ash)                   |

`Ax` first appears at tick `6`, at position `2` in the blend
``Hz*Ax*Cf``, so the answer is *6 × 2* = `12`.
