package main

import (
	"container/list"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

// var logger = log.New(ioutil.Discard, "", log.LstdFlags)

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
	Title string
	Time  int
}

func (t Talk) String() string {
	return fmt.Sprint(t.Time) + " "
}

type Track struct {
	morning   *Session
	afternoon *Session
	nextDay   *Track
}

func newTrack() *Track {
	return &Track{
		morning: &Session{
			MaxLen: 180,
			MinLen: 180,
			talk:   list.New(),
			Done:   false,
			Type:   "morning",
		},
		afternoon: &Session{
			MaxLen: 240,
			MinLen: 180,
			talk:   list.New(),
			Done:   false,
			Type:   "afternoon",
		},
	}
}

func (t *Track) String() (str string) {
	str = fmt.Sprintf("morning : %+v, afternoon : %+v, [%+v]", t.morning, t.afternoon, t.nextDay)
	return
}

func (t *Track) set(talk Talk) (err error) {
	if t.morning.Done == false {
		err = t.morning.set(talk)
	} else if t.afternoon.Done == false {
		err = t.afternoon.set(talk)
	} else {
		if t.nextDay == nil {
			t.nextDay = newTrack()
		}
		err = t.nextDay.set(talk)
	}
	return err
}

func init() {
	initTalkList(testInput)
}

func initTalkList(input []string) {
	for _, v := range input {
		n := strings.LastIndex(v, " ")
		title := v[:n]
		timeStr := v[n+1:]

		if timeStr == "lightning" {
			putToTalkList(talkList, Talk{
				Title: title,
				Time:  5,
			})
		} else if strings.HasSuffix(timeStr, "min") {
			t := timeStr[:len(timeStr)-3]
			time, err := strconv.Atoi(t)
			if err != nil {
				logger.Fatal(err)
			}
			putToTalkList(talkList, Talk{
				Title: title,
				Time:  time,
			})
		}
	}
	logger.Printf("%+v\n", talkList)
	logger.Println("====================================================================")
	return
}

func main() {
	track, err := plan(talkList)
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Printf("result : %+v", track)

	Print(track, 1)
}

func Print(track *Track, counter int) {
	fmt.Printf("Track %d\n", counter)
	beginTime := "9:00AM"
	for e := track.morning.talk.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s %s %dmin\n", beginTime, e.Value.(Talk).Title, e.Value.(Talk).Time)
		beginTime = getClock(beginTime, e.Value.(Talk).Time)
	}

	fmt.Println("12:00PM Lunch")

	beginTime = "1:00PM"
	for e := track.afternoon.talk.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s %s %dmin\n", beginTime, e.Value.(Talk).Title, e.Value.(Talk).Time)
		beginTime = getClock(beginTime, e.Value.(Talk).Time)
	}

	if track.afternoon.getTotalTime() < 180 {
		beginTime = "4:00PM"
	}
	fmt.Printf("%s Networking Event\n", beginTime)

	if track.nextDay != nil {
		Print(track.nextDay, counter+1)
	}
}

type Session struct {
	talk   *list.List
	MaxLen int
	MinLen int
	Done   bool
	Type   string
}

func (s Session) String() (str string) {
	for e := s.talk.Front(); e != nil; e = e.Next() {
		str += e.Value.(Talk).String()
	}
	return str
}

func (s *Session) set(t Talk) (err error) {
	//如果插入的talk总时间大于session的最大时间，则先从session中移除比当前talk时间大的talk, 再插入
	if s.getTotalTime()+t.Time > s.MaxLen {
		for e := s.talk.Back(); e != nil; e = e.Prev() {
			if (e.Value.(Talk).Time - t.Time) > (s.MaxLen - s.getTotalTime()) {
				s.talk.Remove(e)
				putToTalkList(talkList, e.Value.(Talk))
				logger.Printf("remove %d from %s session %+v, talks list is : %+v\n", e.Value.(Talk).Time, s.Type, s, talkList)
			}
		}
	}
	if s.getTotalTime()+t.Time > s.MaxLen {
		err = errors.New("error")
		logger.Printf("errro :: session is %v, talks list is : %+v\n", s, talkList)
		s.Done = true
	} else {
		s.talk.PushBack(t)
		logger.Printf("put %+v to session : %v, talks list is : %+v\n", t, s, talkList)
		if s.getTotalTime() == s.MaxLen {
			s.Done = true
		}
	}
	return err
}

func (s Session) getTotalTime() (total int) {
	if s.talk.Len() == 0 {
		return 0
	}
	for e := s.talk.Front(); e != nil; e = e.Next() {
		total += e.Value.(Talk).Time
	}
	logger.Printf("total Time is %d \n", total)
	return
}

func putToTalkList(l TalkList, t Talk) {
	if l.data.Len() == 0 {
		l.data.PushFront(t)
		return
	}

	for e := l.data.Front(); e != nil; e = e.Next() {
		if t.Time >= e.Value.(Talk).Time {
			l.data.InsertBefore(t, e)
			return
		}
	}
	l.data.PushBack(t)
	logger.Printf("push %v to talk list \n", t)
}

func plan(talkList TalkList) (track *Track, err error) {
	track = newTrack()
	for talkList.data.Len() > 0 {
		var next *list.Element
		for e := talkList.data.Front(); e != nil; e = next {
			next = e.Next()
			err = track.set(e.Value.(Talk))
			if err != nil {
				break
			}
			talkList.data.Remove(e)
			logger.Printf("remove %+v from talk list : %+v \n", e.Value.(Talk), talkList)
		}
		if err != nil {
			logger.Println(err)
			logger.Printf("talks list is : %+v\n", talkList)
		}
	}
	return
}

func getClock(begin string, duration int) string {
	t, _ := time.Parse("3:04PM", begin)
	t = t.Add(time.Minute * time.Duration(duration))
	return t.Format("3:04PM")
}

var testInput = []string{
	"Writing Fast Tests Against Enterprise Rails 60min",
	"Overdoing it in Python 45min",
	"Lua for the Masses 30min",
	"Ruby Errors from Mismatched Gem Versions 45min",
	"Common Ruby Errors 45min",
	"Rails for Python Developers lightning",
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

var testInput1 = []string{
	"spend 75min",
	"spend 75min",
	"spend 69min",
	"spend 69min",
	"spend 68min",
	"spend 60min",
	"spend 59min",
	"spend 57min",
	"spend 48min",
	"spend 48min",
	"spend 43min",
	"spend 43min",
	"spend 40min",
	"spend 40min",
	"spend 39min",
	"spend 38min",
	"spend 37min",
	"spend 37min",
	"spend 36min",
	"spend 36min",
	"spend 34min",
	"spend 33min",
	"spend 33min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 31min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 30min",
	"spend 27min",
	"spend 27min",
	"spend 27min",
	"spend 27min",
	"spend 27min",
	"spend 27min",
	"spend 27min",
	"spend 55min",
	"spend 55min",
	"spend 55min",
	"spend 55min",
}
