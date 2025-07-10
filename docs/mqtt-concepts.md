Broker: The central server (our Mosquitto container) that receives all messages and forwards them to interested clients.
Client: Any program that connects to the broker (ping and pong Go apps). Each client needs a unique ID.
Publish: The act of sending a message to a specific "address" on the broker.
Subscribe: The act of telling the broker "I am interested in messages sent to this address."
Topic: The "address" of a message (locus/test). It's a string that lets the broker know where to route messages.