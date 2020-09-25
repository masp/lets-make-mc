# Let's Make Minecraft
Welcome! This blog tries to be a fun guide/tour on building a Minecraft server. It's more of a blog than a book, because it covers my thought processes while I work on building certain features. It doesn't have a central theme other than building something cool.

My goal:
- Fun!
- Learn how to build large software projects
- Explore interesting topics

## Why Minecraft?
My goal is to explore and learn interesting topics in software using a big project, but also have fun while doing it.
Minecraft is a game with both demanding technical requirements as well as a game that happens to be very fun and interesting to develop for!

It has a lot of interesting technical requirements. In particular, it has a large game state that has to be asynchronously accessed
and modified by many clients simultaneously. It needs to handle disk and network resources efficiently because of how much of it
is needed.

Minecraft also has an extremely rich modding community, that makes tasks like understanding the network protocol and file formats a much easier endeavor. There's also a variety of other Minecraft servers that were able to borrow ideas and techniques from to let us focus on the more interesting tasks at hand.

## What is planned
Anything that is interesting technically. Components like chunk persistence, automated testing, and networking are all pieces that are interesting beyond Minecraft to other domains.

## What is not planned
Anything in Minecraft that is more of a feature of the game is not really important. If you want to play Minecraft with your friends, there's a huge number of great alternatives out there that can give you all the features you'd ever want. 

## How to Read This Blog
This blog can be thought of a guide of milestones or key features that are necessary for a Minecraft server. I will skip over a lot of features and pieces that a proper Minecraft would have, but also spend a lot more time on pieces that are really not that important if you're just wanting a Minecraft server.