# 05: Barcamp Lineup {#1}

Thursday. The barcamp kicks off after coffee, and the shared
spreadsheet everyone's been dumping talk pitches into all week is
finally closed for edits. Of the two TVs on the living-room wall, one
has been claimed for the barcamp - one screen, one HDMI cable, and one
afternoon to get through the lineup, so the organisers need to decide
which talks actually make it onto the screen.

Someone has already walked the sheet and noted, next to each row, a
*hype* score - a running tally of reactions the pitch has collected;
nobody bothered to cap how many each attendee could throw at their
favourites.

Someone else - probably the same person - has drafted a list of pairs
that would step on each other: shared speaker, shared laptop, or two
talks leaning on the same running joke. Chains of clashes end up
pulling pitches together into little *topic groups*, and the hosts
want to know which pitch sits in which.

The hosts have a stack of last-minute questions on top: *"would these
two land in the same group?"*

## Input

Your puzzle input has three sections separated by blank lines:

1. *hype* - one line per talk as `<slug> <int>`. Slugs are unique
   lowercase identifiers built from letters, digits and underscores.
   The integer is that talk's hype score.
2. *conflicts* - one line per direct conflict as `<slug> <slug>`.
   Conflicts are symmetric.
3. *questions* - one line per host question as `<slug> <slug>`.

A *topic group* is a group of talks tied together by conflicts. Two
talks share a group if you can walk from one to the other by hopping
along direct conflicts.

## Part 1

For each question, check whether the two talks belong to *different*
topic groups. Every such question contributes the sum of *both* its
talks' hype scores. Return the *total* across all questions.

### Example

```
factorio_deep_dive 60
factorio_101 25
bulanci_101 15
bulanci_hot_takes 40
sauna_field_notes 30

factorio_deep_dive factorio_101
factorio_101 bulanci_101
bulanci_hot_takes sauna_field_notes

factorio_deep_dive bulanci_101
factorio_deep_dive sauna_field_notes
```

Here `factorio_deep_dive`, `factorio_101` and `bulanci_101` form one
group, and `bulanci_hot_takes`, `sauna_field_notes` form another.

The pair `factorio_deep_dive / bulanci_101` share a
group (via `factorio_101`) and contribute nothing, while
`factorio_deep_dive / sauna_field_notes` span two groups and
contribute `60 + 30 = 90` - so the answer is `90`.
