# Problem description 

Write a TCP service in a language of your choice. If you are comfortable with Go, use that, but other languages are
acceptable as long as you can explain the pros/cons and scaling characteristics.The service should have the following end points:

1. Start listening for tcp/udp connections.
2. Be able to accept connections.
3. Read json payload ({"user_id": 1, "friends": [2, 3, 4]})
3. After establishing successful connection - "store" it in memory the way you like.
4. When another connection established with the user_id from the list of any other user's "friends" section,
they should be notified about it with message {"online": true}
5. When the user goes offline, his "friends" (if it has any and any of them online) should receive a message {"online": false}

# How to run

To run server use: 
```
make server
```

To run one of 4 clients use: 
```
make user<n>
```
where N is from 1 to 4.

User `Ctrl+C` to stop server or client.
