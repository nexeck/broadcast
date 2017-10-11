package multicast

import (
	"sync"
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
	defer group.Quit()

	go group.Broadcast()
	group.Send("test payload")
}

func TestGroupSendAndOneMemberRead(t *testing.T) {
	testMessage := "test payload"

	group := NewGroup()
	defer group.Quit()

	go group.Broadcast()

	member := group.Join()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup, m *Member, t *testing.T) {
		defer wg.Done()

		for {
			select {
			case val := <-m.Read():
				if val != testMessage {
					t.Errorf("Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Read: '%s' equal '%s'", m, testMessage)
				}
				return
			case <-time.After(5 * time.Second):
				t.Error("Read timedout")
				return
			}
		}
	}(wg, member, t)

	group.Send(testMessage)

	wg.Wait()
}

func TestGroupSendAndMultipleMembersRead(t *testing.T) {
	testMessage := "test payload"

	group := NewGroup()
	defer group.Quit()

	go group.Broadcast()

	member1 := group.Join()
	member2 := group.Join()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(m *Member, t *testing.T) {
		defer wg.Done()

		for {
			select {
			case value := <-m.Read():
				if value != testMessage {
					t.Errorf("Member1 Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Member1 Read: '%s' equal '%s'", m, testMessage)
				}
				return
			case <-time.After(5 * time.Second):
				t.Error("Member1 Read timedout")
				return
			}
		}
	}(member1, t)

	wg.Add(1)
	go func(m *Member, t *testing.T) {
		defer wg.Done()

		for {
			select {
			case val := <-m.Read():
				if val != testMessage {
					t.Errorf("Member2 Read: '%s' not '%s'", m, testMessage)
				} else {
					t.Logf("Member2 Read: '%s' equal '%s'", m, testMessage)
				}
				return
			case <-time.After(5 * time.Second):
				t.Error("Member2 Read timedout")
				return
			}
		}
	}(member2, t)

	group.Send(testMessage)

	wg.Wait()
}
