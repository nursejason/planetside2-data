# planetside2-data

### Player Logins
Judging by the PS2 data, there does not appear to be data that simply lists the last X logins and Y duration that player spent online. There is however `last_login` and `last_save`. This isn't ideal, as it would require a script to constantly be running and polling the API. Any service interruption in the script would cause the data to become invalid for the downed duration.        
Additionally, I'll need to determine the interval in between each `last_save` as a sleep time between each run of the script.      
With this information, I will be able to query Planetside 2 API:      
1. All players every X duration for `last_save >  NOW - X - 1`
2. Store Player: Name, Faction, Server, and Continent in addition to `last_played_timestamp == NOW`
