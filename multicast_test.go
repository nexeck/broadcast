package multicast

import (
	"testing"
	"time"
)

func TestNewGroup(t *testing.T) {
	group := NewGroup()
	if group == nil {
		t.Error("Could not create group")
	}
}

func TestGroupQuit(t *testing.T) {
	group := NewGroup()

	quitChannel := make(chan bool)

	go func(quitChannel chan bool) {
		group.Broadcast()
		quitChannel <- true
	}(quitChannel)

	group.Quit()

	if <-quitChannel != true {
		t.Errorf("Quit failed")
	}
}

func TestGroupJoin(t *testing.T) {
	group := NewGroup()
	member := group.Join()

	if member == nil {
		t.Error("Could not create member and join group")
	}
}

func TestGroupMemberCount(t *testing.T) {
	group := NewGroup()

	group.Join()
	memberCount := group.MemberCount()
	if memberCount != 1 {
		t.Errorf("MemberCount should be 1, but is %d", memberCount)
	}

	group.Join()
	memberCount = group.MemberCount()
	if memberCount != 2 {
		t.Errorf("MemberCount should be 2, but is %d", memberCount)
	}
}

func TestGroupSend(t *testing.T) {
	group := NewGroup()

	go group.Broadcast()
	group.Send("test payload")
}

func TestGroupSendAndOneMemberRead(t *testing.T) {
	testMessage := "test payload"

	group := NewGroup()
	go group.Broadcast()

	member := group.Join()
	group.Send(testMessage)

	go func(m *Member, t *testing.T) {
		for {
			select {
			case val := <-m.Read():
				if val != testMessage {
					t.Errorf("Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Read: '%s' equal '%s'", m, testMessage)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Read timedout")
			}
		}
	}(member, t)
}

func TestGroupSendAndMultipleMembersRead(t *testing.T) {
	testMessage := "test payload"

	group := NewGroup()
	go group.Broadcast()

	member1 := group.Join()
	member2 := group.Join()
	group.Send(testMessage)

	go func(m *Member, t *testing.T) {
		for {
			select {
			case value := <-m.Read():
				if value != testMessage {
					t.Errorf("Member1 Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Member1 Read: '%s' equal '%s'", m, testMessage)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Member1 Read timedout")
			}
		}
	}(member1, t)

	go func(m *Member, t *testing.T) {
		for {
			select {
			case val := <-m.Read():
				if val != testMessage {
					t.Errorf("Member2 Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Member2 Read: '%s' equal '%s'", m, testMessage)
				}
			case <-time.After(5 * time.Second):
				t.Errorf("Member2 Read timedout")
			}
		}
	}(member2, t)
}
