# tragedyofthecoWmmons
An exploration of parallelism in Go, using the example of monte-carlo simulated cows feeding on a a 2-d board of "grass."

The first cow, Clarabelle, has no preference for direction when she moves.  She is equally likely to move in any direction at all times.  Let p be the probability that she moves in a direction.
The second cow, Bernadette, prefers to move towards squares with grass.  She is twice as likely to move towards a space that has grass as she is towards a space that she has already eaten.

By defualt, we simulate the two cows' 144 movements (48 hours of "eating") 10,000 times, parallelized over the number of processors locally available.
