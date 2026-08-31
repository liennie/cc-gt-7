## Part 2 {#2}

Someone squats next to a running unit and watches its opening pose:
elbow up, tracks crooked, LED spine flickering in a sad rainbow.
*"That's not the pose it should be striking. The neutral 'stand at
attention' is meant to be six zeros on the bus. Whoever set up the
assembly line got creative and now the boot ROM is producing
nonsense."* They tap the debug port on the GT-Nano's chassis.
*"But look - if we seed the first six memory cells before the ROM
runs, the LOADs will pick those up instead of zero. So, what do we
put in there?"*

Determine the six signed integers `X0 X1 X2 X3 X4 X5` that, when
written into cells `0 1 2 3 4 5` before execution begins, cause every
output value to be `0`. Everything else - the registers, all other
tape cells, the cursor - still starts at `0`.

Return the positional base-10 digest of the six values:

    1×X0 + 10×X1 + 100×X2 + 1000×X3 + 10000×X4 + 100000×X5

For the same `17`-line program from part 1, writing `X0 = -2` into
cell `0` and `X1 = 3` into cell `1` before booting zeros both outputs,
so the scaled-down digest is *1×(-2) + 10×3* = `28`.
