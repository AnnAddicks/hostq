# hostq
Project to determine who is hosting our weekly dinner group. 


=======
## Process
1. Cron task to send group emails every Sunday at 4pm who is next in queue.
2. Someone replies from any of the emails in the next hosting group.  (Multiple emails are use for any family member to reply or a business and personal email can be used for the hosting party.)
3. Look for the following responses `Yes`, `No`, `Skip` from **only** the hosting party for now:  
  *  Yes - update the queue and who is hosting next
  *  No - send an email to the next in line and update group
  *  Skip - send an email response with the queue order for the following week.


## Future Goals
 * Prevent the next in queue to repond early.  Example:  It is family A's turn to host this week, they respond with yes.  Family B knows thier turn is next week and respond with yes.  This will cause the next weeks email to be sent as Family C's turn.
 * Group defined cron schedule
 * Group Customizations
  * Can/Should anyone in the group be able to reply for skip
  * Can the person who wants to host reply with "me"
  * Give the responder a timeout, if they do not reply then after x number of hours go to the next person. 
 * Front end
 * Keep track of who hosted when to have a year's reflection.
  * You hosted x number of times
  * November is when you are least likely to host
  * Congratulations on 50 hosted events

## Code Review Suggestions (my thoughts in parens)
- [ ] more error checking (OMG!  the error checking in GO is crazy!)
- [ ] Check out Go's timers (Can I use timers on app engine?)
- [ ] Decouple Hosts & Group
- [ ] Api Links to accept hosting.  (If we do this, then probably do not send the the group email, but to an individual)
- [ ] Make the hosting rotation customizable.  (Our group rotates a queue so this will revisited once our needs are met)
 *  Exclude hosts which have chosen to be excluded (by web interface, perhaps, or a vacation email).
 * Order the remaining hosts by the number of times they have hosted.
 * If there exists a tie, select the host who last hosted longest ago.
 * If there exists a tie, choose based on stored slice order or random.
 * Send a notification and remove the host from the queue.
 * If the notification is rejected or times out, select the next remaining host and repeat 5.
 * If the hosts are exhausted, announce that the dinner has been canceled.

## Requests from dining in group
- [x] Mine:  Ack back when someone responds with 'yes'. 
- [x] Make a single responder that gives the current order to whoever emails with "status" or similar.  (After talking with Josh, this could also be a static api link.)
- [ ] Make it remember the current menu. This particular idea isn't fully baked, admittedly.
- [ ] Possibly include queue order on every email.
- [ ] One more thing you might add is actual dates - it doesn't have to be anything fancy, just know that the order is family a, b, c, and the dates of the next three saturdays, etc.
- [ ] Add a link or response to say who will be there so the host could dish out assignments.  (Pun intended!)

