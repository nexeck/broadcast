package multicast

import "sync"

// Member of group
type Member struct {
	group *Group
	in    chan interface{}
}

// Group broadcasts payloads to members
type Group struct {
	in         chan interface{}
	quit       chan bool
	members    []*Member
	memberLock sync.Mutex
}

// NewGroup creates a new broadcast group.
func NewGroup() *Group {
	in := make(chan interface{})
	quit := make(chan bool)
	return &Group{in: in, quit: quit}
}

// Join returns a new member object and adds it to the group
func (g *Group) Join() *Member {
	g.memberLock.Lock()
	defer g.memberLock.Unlock()

	member := &Member{
		group: g,
		in:    make(chan interface{}),
	}

	g.members = append(g.members, member)
	return member
}

// MemberCount returns the number of members in the Group.
func (g *Group) MemberCount() int {
	return len(g.Members())
}

// Members returns a slice of Members that are currently in the Group.
func (g *Group) Members() []*Member {
	g.memberLock.Lock()
	res := g.members[:]
	g.memberLock.Unlock()
	return res
}

// Send broadcasts a message to the group.
func (g *Group) Send(val interface{}) {
	g.in <- val
}

// Broadcast messages to group members.
func (g *Group) Broadcast() {
	for {
		select {
		case received := <-g.in:
			g.memberLock.Lock()
			members := g.members[:]
			g.memberLock.Unlock()

			for _, member := range members {
				// This is done in a goroutine because if it
				// weren't it would be a blocking call
				go func(member *Member, received interface{}) {
					member.in <- received
				}(member, received)
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

// Read reads one value from the member's Read channel
func (m *Member) Read() interface{} {
	return <-m.in
}
