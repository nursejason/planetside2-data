# planetside2-data

### Player Logins
##### The ultimate goal is to create graphs based on server for total hourly population per empire based on time. Then be able to overlay day over day or week over week to determine rising or falling population as time goes on, and if things like world events (alerts for example) cause people to log off after a continent is locked. 

*Wrote a script to determine `last_save` interval. Turns out it was fairly erratic and that there's a websocket stream for login / logout events*            
      
MetagameEvent - Will be good to notate when an alert starts and stops       
PlayerLogin + PlayerLogout - Will give overall information of total population on each server, but it doesn't include faction information.        
      
   
   
```JSON
"MetagameEvent":{
		"event_name":"MetagameEvent",
		"experience_bonus":"",
		"faction_nc":"",
		"faction_tr":"",
		"faction_vs":"",
		"metagame_event_id":"",
		"metagame_event_state":"",
		"timestamp":"",
		"world_id":"",
		"zone_id":""
	},	
"PlayerLogin":{
		"character_id":"",
		"event_name":"PlayerLogin",
		"timestamp":"",
		"world_id":""
	},
"PlayerLogout":{
		"character_id":"",
		"event_name":"PlayerLogout",
		"timestamp":"",
		"world_id":""
	}
```
