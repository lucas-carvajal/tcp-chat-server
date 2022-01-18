## tcp-chat-server

### What is it and what does it do?
TCP chat server built in Go.
`WIP`

### How to use it?
1. Make sure you have Go 1.16 installed
2. Run the main function in the 'main.go' file
3. Use telnet to connect to the server
   1. `telnet localhost 8888`
   2. `/nick lucas` to give you a nick
   3. `/join football` to join the football chat
   4. `/msg hi!` to send "hi!" to the chat


### What did I learn?
* Go
* Goroutines

### Disclaimers
This project was created following Alex Pliutau's video on YouTube.  
You can find it here: https://www.youtube.com/watch?v=Sphme0BqJiY&t=179s

Functionality I added includes:
* Users can omit the `/msg` command and the input will be sent as a message automatically
* After connecting, all users will automatically be added to the 'lobby' chat
* Users can exit a room and will be placed back in the 'lobby' chat
* Users can see all members in a room
* Users can send direct messages to members of their current room
* Users can get help by seeing a list of all commands



---

#### Commands

- `/nick <name>` - get a name, otherwise will stay anonymous
- `/join <name>` - join a room, if room does not exist, it will be created. User can only be in one room at the same time
- `/rooms` - shows list of available rooms to join
- `/msg <msg>` - broadcasts message to everyone in the room
- `/quit` - disconnect from the chat server
- `/quitRoom` - quit a room and go back to the lobby chat
- `/roomMembers` - shows list of all members in the current room 
- `/dm <name> <msg>` - send a direct message to a member of the current room
- `/help` - shows list of all available commands
