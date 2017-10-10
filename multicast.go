package multicast

import "sync"

// Member of group
type Member struct {
	group   *Group
	message chan interface{}
}

// Group broadcasts payloads to members
type Group struct {
	message    chan interface{}
	quit       chan bool
	members    []*Member
	memberLock sync.Mutex
}

// NewGroup creates a new broadcast group.
func NewGroup() *Group {
	message := make(chan interface{})
	quit := make(chan bool)
	return &Group{message: message, quit: quit}
}

// Join returns a new member object and adds it to the group
func (g *Group) Join() *Member {
	g.memberLock.Lock()
	defer g.memberLock.Unlock()

	member := &Member{
		group:   g,
		message: make(chan interface{}),
	}

	g.members = append(g.members, member)
	return member
}

// MemberCount returns the number of members message the Group.
func (g *Group) MemberCount() int {
	return len(g.Members())
}

// Members returns a slice of Members that are currently message the Group.
func (g *Group) Members() []*Member {
	g.memberLock.Lock()
	res := g.members[:]
	g.memberLock.Unlock()
	return res
}

// Send broadcasts a message to the group.
func (g *Group) Send(val interface{}) {
	g.message <- val
}

// Broadcast messages to group members.
func (g *Group) Broadcast() {
	for {
		select {
		case message := <-g.message:
			g.memberLock.Lock()
			members := g.members[:]
			g.memberLock.Unlock()

			for _, member := range members {
				// This is done message a goroutine because if it
				// weren't it would be a blocking call
				go func(member *Member, message interface{}) {
					member.message <- message
				}(member, message)
			}
		case <-g.quit:
			return
		}
	}
}

// Quit terminates the group immediately.
func (g *Group) Quit() {
	g.quit <- true
}

// Read returns the message channel
func (m *Member) Read() <-chan interface{} {
	return m.message
}
