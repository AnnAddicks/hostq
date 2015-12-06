# hostq
Project to determine who is hosting our weekly dinner group. 


## Process
1. Cron task to send group emails every Sunday at 4pm who is next in queue.
2. Someone replies from any of the emails in the next hosting group.  (Multiple emails are use for any family member to reply or a business/personal email can be used for the hosting party.)
3. Look for the following responses `Yes`, `No`, `Skip` from **only** the hosting party for now:  
  *  Yes - update the queue and who is hosting next
  *  No - send an email to the next in line and update group
  *  Skip - send an email response with the queue order for the following week.
