package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var talkList = TalkList{
	data: list.New(),
}

type TalkList struct {
	data *list.List
}

func (t TalkList) String() (str string) {
	for e := t.data.Front(); e != nil; e = e.Next() {
		str += e.Value.(Talk).String()
	}
	return str
}

type Talk struct {
	title string
	time  int
}

func (t Talk) String() string {
	return fmt.Sprint(t.time) + " "
}

type Session struct {
	talk   *list.List
	MaxLen int
	MinLen int
	Done   bool
}

func (s Session) String() (str string) {
	for e := s.talk.Front(); e != nil; e = e.Next() {
		str += e.Value.(Talk).String()
	}
	return str
}

type Track struct {
	morning   Session
	afternoon Session
	nextDay   *Track
}

func init() {
	initTalkList(testInput)
}

func main() {
	track, err := plan(talkList)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("%+v", track)
}
func (s Session) set(t Talk) (err error) {
	//如果插入的talk总时间大于session的最大时间，则先从session中移除比当前talk时间大的talk, 再插入
	if s.getTotalTime()+t.time > s.MaxLen {
		for e := s.talk.Back(); e != nil; e = e.Prev() {
			if e.Value.(Talk).time > t.time {
				s.talk.Remove(e)
				putTalkToTalkList(talkList, e.Value.(Talk))
				fmt.Printf("remove %d from %+v, talks list is : %+v\n", e.Value.(Talk).time, s, talkList)
			}
		}
	}
	if s.getTotalTime()+t.time > s.MaxLen {
		err = errors.New("error")
		fmt.Printf("session is %v, talks list is : %+v\n", s, talkList)
		return err
	}
	s.talk.PushBack(t)
	//如果刚好填满session，设置session为已完成
	if s.getTotalTime()+t.time == s.MaxLen {
		s.Done = true
	}
	return nil
}

func (s Session) getTotalTime() (total int) {
	if s.talk.Len() == 0 {
		return 0
	}
	for e := s.talk.Front(); e != nil; e = e.Next() {
		total += e.Value.(Talk).time
	}
	return
}

func newTrack() *Track {
	return &Track{
		morning: Session{
			MaxLen: 180,
			MinLen: 180,
			talk:   list.New(),
			Done:   false,
		},
		afternoon: Session{
			MaxLen: 240,
			MinLen: 180,
			talk:   list.New(),
			Done:   false,
		},
	}
}

func (t Track) set(talk Talk) (err error) {
	if t.morning.Done == false {
		err = t.morning.set(talk)
	} else if t.afternoon.Done == false {
		err = t.afternoon.set(talk)
	} else {
		t.nextDay = newTrack()
		err = t.nextDay.set(talk)
	}
	return err
}

func initTalkList(input []string) {
	for _, v := range input {
		n := strings.LastIndex(v, " ")
		title := v[:n]
		timeStr := v[n+1:]

		if timeStr == "lightning" {
			putTalkToTalkList(talkList, Talk{
				title: title,
				time:  5,
			})
		} else if strings.HasSuffix(timeStr, "min") {
			t := timeStr[:len(timeStr)-3]
			time, err := strconv.Atoi(t)
			if err != nil {
				log.Fatal(err)
			}
			putTalkToTalkList(talkList, Talk{
				title: title,
				time:  time,
			})
		}
	}
	fmt.Printf("%+v\n", talkList)
	fmt.Println("====================================================================")
	return
}

func putTalkToTalkList(l TalkList, t Talk) {
	if l.data.Len() == 0 {
		l.data.PushFront(t)
		return
	}

	for e := l.data.Front(); e != nil; e = e.Next() {
		if t.time >= e.Value.(Talk).time {
			l.data.InsertBefore(t, e)
			return
		}
	}
	l.data.PushBack(t)
	fmt.Printf("push %v to talk list", t)
}

func plan(talks TalkList) (track *Track, err error) {
	track = newTrack()
	for e := talks.data.Front(); e != nil; e = e.Next() {
		err = track.set(e.Value.(Talk))
		if err != nil {
			break
		}
	}
	if err != nil {
		// return plan(talks)
		return nil, err
	}
	return nil, nil
}

var testInput = []string{
	"Writing Fast Tests Against Enterprise Rails 60min",
	"Overdoing it in Python 45min",
	"Lua for the Masses 30min",
	"Ruby Errors from Mismatched Gem Versions 45min",
	"Common Ruby Errors 45min",
	"Rails for Python Developers lightning",
	"Communicating Over Distance 60min",
	"Communicating Over Distance 60min",
	"Accounting-Driven Development 45min",
	"Woah 30min",
	"Sit Down and Write 30min",
	"Pair Programming vs Noise 45min",
	"Rails Magic 60min",
	"Ruby on Rails: Why We Should Move On 60min",
	"Clojure Ate Scala (on my project) 45min",
	"Programming in the Boondocks of Seattle 30min",
	"Ruby vs. Clojure for Back-End Development 30min",
	"Ruby on Rails Legacy App Maintenance 60min",
	"A World Without HackerNews 30min",
	"User Interface CSS in Rails Apps 30min",
}
