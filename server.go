package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	if len(s.rooms) < 1 {
		r := &room{
			name:    "lobby",
			members: make(map[net.Addr]*client),
		}
		s.rooms["lobby"] = r
	}

	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_MEMBERS:
			s.roomMember(cmd.client)
		case CMD_DM:
			s.dm(cmd.client, cmd.args)
		case CMD_HELP:
			s.help(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}

	s.join(c, []string{"/join", "lobby"})

	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms are: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must join the room first"))
		return
	}

	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:len(args)], " "))
}

func (s *server) quit(c *client) {
	log.Printf("client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go :(")

	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) roomMember(c *client) {
	if c.room != nil {
		var members []string

		for member := range c.room.members {
			members = append(members, c.room.members[member].nick)
		}

		c.msg(fmt.Sprintf("Current members of room %s: %s", c.room.name, strings.Join(members, ", ")))
	}
}

func (s *server) dm(c *client, args []string) {
	var receiver *client

	clientRoomMembers := c.room.members

	for member := range clientRoomMembers {
		if clientRoomMembers[member].nick == args[1] {
			receiver = clientRoomMembers[member]
		}
	}

	if receiver == nil {
		c.err(errors.New(fmt.Sprintf("user %s not found", args[1])))
		return
	}

	receiver.dmsg(strings.Join(args[2:], " "), c)

}

func (s *server) help(c *client) {
	c.msg("Available Commands:")
	c.msg(" /nick <name> \t - get a name, otherwise you will stay anonymous")
	c.msg(" /join <name> \t - join a room, if room does not exist, it will be created. User can only be in one room at the same time")
	c.msg(" /rooms \t \t - shows list of available rooms to join")
	c.msg(" /msg <msg> \t \t - broadcasts message to everyone in the room")
	c.msg(" /quit \t \t - disconnect from the chat server")
	c.msg(" /quitRoom \t \t - quit a room and go back to the lobby chat")
	c.msg(" /roomMembers \t - shows list of all members in the current room")
	c.msg(" /dm <name> <msg> \t - send a direct message to a member of the current room")
}
