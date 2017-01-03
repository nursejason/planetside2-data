# planetside2-data

### Player Logins
##### The ultimate goal is to create graphs based on server for total population per empire based on time. Then be able to overlay day over day or week over week to determine rising or falling population as time goes on, and if things like world events (alerts for example) cause people to log off after a continent is locked. 

*Wrote a script to determine `last_save` interval. Turns out it was fairly erratic and that there's a websocket stream for login / logout events*            
      
MetagameEvent - Will be good to notate when an alert starts and stops       
PlayerLogin + PlayerLogout - Will give overall information of total population on each server, but it doesn't include faction information. Since it does provide CharacterID however, that should be able to be retrieved.       
      
###### Aggregating the data    
ElasticSearch's Aggregation feature should make this reasonably easy. With the example ElasticDoc provided below, the primary bucket will use a terms aggregation to separate based on `world_id`, and a sub bucket being a date range between `login_ts` and `logout_ts`. 
       
```JSON
ElasticDoc
ID = characterID
{
	"character_name": <name_str>,
	"character_id": <id_str>,
	"faction": <faction_str>,
	"faction_id": <faction_id>,
	"world_id": <world_id>,
	"world_name": <world_name>,
	"login_logout_events: [
		{
			"login_ts": <login_ts>,
			"logout_ts": <logout_ts>
		}, ...
	]
}
```

###### Getting character info
In an effort to reduce external calls, an attempt will first be made to update an existing document. If that document does not exist, then the character information will be retrieved from the external API. 

###### Resiliency
Step 2 will be to improve resiliency. This will be achieved by isolating the main worker into two pieces.       
- First is the websocket subscriber. This will simply subscribe to the websocket and then produce the information into a Kafka topic. This way, if the Planetside 2 API goes down, and somehow the websocket is still alive, then login / logout data is less likely to be lost. Or, more likely, if the calling script is encounters an error / Elasticsearch is down, then the entire process doesn't come to a halt.        
- Secondly the piece that checks against Elastic, calls Planetside API, and saves to Elasticsearch will be a second running script.      
     
###### UI  
:shrug:
   
###### Stream Payloads
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
