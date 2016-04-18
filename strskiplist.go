package main

import (
	"math/rand"
	"strings"
    "fmt"
)

const MAX_LEVELS = 32

func RandLevel() (lvl int) {
	for i := 1; i <= MAX_LEVELS; i++ {
		if rand.Intn(3) != 0 {
			return i
		}
	}
	return MAX_LEVELS
}

type StrSkipListLevel struct {
	Prev *StrSkipListElement
	Next *StrSkipListElement
}

type StrSkipListElement struct {
	Value  interface{}
	key    string
	levels []StrSkipListLevel
}

type StrSkipListRange struct {
    Begin *StrSkipListElement
    End *StrSkipListElement
}

func (self StrSkipListRange) String() string {
    ret := make([]string, 0)
    for beg := self.Begin; beg != self.End; beg = beg.Next(){
        ret = append(ret, beg.String())
    }
    return strings.Join(ret, ",")
}


func (self *StrSkipListElement) String() string {
    return fmt.Sprint(self.key, ":", self.Value)
}

func (self *StrSkipListElement) Key() string {
	return self.key
}

func (self *StrSkipListElement) nextLevel(lvl int) *StrSkipListElement {
	return self.levels[lvl-1].Next
}

func (self *StrSkipListElement) prevLevel(lvl int) *StrSkipListElement {
	return self.levels[lvl-1].Prev
}

func (self *StrSkipListElement) Next() *StrSkipListElement {
	return self.nextLevel(1)
}

func (self *StrSkipListElement) Prev() *StrSkipListElement {
	return self.prevLevel(1)
}

func (self *StrSkipListElement) setLevelNext(lvl int, p *StrSkipListElement) {
	self.levels[lvl-1].Next = p
}

func (self *StrSkipListElement) setLevelPrev(lvl int, p *StrSkipListElement) {
	self.levels[lvl-1].Prev = p
}

func (self *StrSkipListElement) insertBehind(lvl int, p *StrSkipListElement) {
	p.setLevelPrev(lvl, self)
	p.setLevelNext(lvl, self.nextLevel(lvl))
	self.setLevelNext(lvl, p)
    if p.nextLevel(lvl) != nil{
        p.nextLevel(lvl).setLevelPrev(lvl, p)
    }
}

type StrSkipList struct {
	head *StrSkipListElement
    LevelHeight int
	Size uint64
}

func NewStrSkipList() *StrSkipList {
	return new(StrSkipList).Init()
}

func (self *StrSkipList) Init() *StrSkipList {
	self.head = &StrSkipListElement{levels: make([]StrSkipListLevel, MAX_LEVELS)}
    self.LevelHeight = 0
	self.Size = 0
	return self
}

func (self *StrSkipList) Begin() *StrSkipListElement {
	return self.head.Next()
}

func (self *StrSkipList) Add(k string, v interface{}) *StrSkipListElement {
	rlvl := RandLevel()
	pele := &StrSkipListElement{Value: v, key: k}
	pinsertpos := make([]*StrSkipListElement, rlvl)
    var MaxLvl int
    if self.LevelHeight < rlvl{
        MaxLvl = rlvl
    } else {
        MaxLvl = self.LevelHeight
    }
	ppos := self.head
	for ilvl := MaxLvl; ilvl > 0; ilvl-- {
		for ; ppos.nextLevel(ilvl) != nil; ppos = ppos.nextLevel(ilvl) {
			comp := strings.Compare(ppos.nextLevel(ilvl).key, k)
			if comp < 0 {
				continue
			} else if comp > 0 {
				break
			} else {
				pele.levels = make([]StrSkipListLevel, 1)
				ppos.nextLevel(ilvl).insertBehind(1, pele)
				return pele
			}
		}
		if rlvl >= ilvl {
			pinsertpos[ilvl-1] = ppos
		}
	}
	pele.levels = make([]StrSkipListLevel, rlvl)
	for ilvl := rlvl; ilvl > 0; ilvl-- {
		pinsertpos[ilvl-1].insertBehind(ilvl, pele)
	}
    if self.LevelHeight < MaxLvl{
        self.LevelHeight = MaxLvl
    }
	self.Size++
	return pele
}

func (self *StrSkipList) findPos(k string) (bool, *StrSkipListElement) {
	ppos := self.head
	for ilvl := self.LevelHeight; ilvl > 0; ilvl-- {
        for {
            next := ppos.nextLevel(ilvl)
            if next == nil{
                break
            }
            comp := strings.Compare(next.key, k)
            if comp == 0 {
                return true, next
            }
            if comp > 0 {
                break
            }
            ppos = next
        }
	}
	return false, ppos.Next()
}

func (self *StrSkipList) FindOne(k string) *StrSkipListElement {
	equal, ret := self.findPos(k)
	if equal {
		return ret
	}
	return nil
}

func (self *StrSkipList) FindAll(k string) StrSkipListRange {
	ret := self.FindOne(k)
    if ret == nil{
        return StrSkipListRange{nil, nil}
    }
	pele := ret
	for ; pele != nil && pele.key == k; pele = pele.Next() {

	}
	return StrSkipListRange{ret, pele}
}

func (self *StrSkipList) FindRange(kbeg, kend string) StrSkipListRange {
	comp := strings.Compare(kbeg, kend)
	if comp > 0 {
		return StrSkipListRange{nil, nil}
	} else if comp == 0 {
		return self.FindAll(kbeg)
	} else {
		_, pbeg := self.findPos(kbeg)
        _, pend := self.findPos(kend)
		return StrSkipListRange{pbeg, pend}
	}
}

func (self *StrSkipList) getPrefixEnd(pre string) *StrSkipListElement{
    slicepre := []byte(pre)
    for i := len(slicepre) - 1; i >= 0; i--{
        slicepre[i]++
        if slicepre[i] != 0{
            _,ret := self.findPos(string(slicepre))
            return ret
        }
    }
    return nil
}

func (self *StrSkipList) FindPrefix(pre string) StrSkipListRange {
    _, ppre1 := self.findPos(pre)
    if ppre1 == nil{
        return StrSkipListRange{nil, nil}
    }
    return StrSkipListRange{ppre1,self.getPrefixEnd(pre)}
}

func (self *StrSkipList) Remove(ele *StrSkipListElement) (ret *StrSkipListElement) {
    ret = ele.Next()
	for ilvl := len(ele.levels); ilvl > 0; ilvl-- {
		if ele.nextLevel(ilvl) != nil {
			ele.nextLevel(ilvl).setLevelPrev(ilvl, ele.prevLevel(ilvl))
		}
		ele.prevLevel(ilvl).setLevelNext(ilvl, ele.nextLevel(ilvl))
		ele.setLevelPrev(ilvl, nil)
		ele.setLevelNext(ilvl, nil)
        if self.head.nextLevel(ilvl) == nil && self.LevelHeight == ilvl{
            self.LevelHeight--
        }
	}
	self.Size--
	return ret
}

func (self *StrSkipList) RemoveRange(rg StrSkipListRange) {
    for beg := rg.Begin; beg != rg.End; beg = self.Remove(beg){

    }
}

func (self *StrSkipList) String() string {
	ret := make([]string, 0)
    for ilvl := self.LevelHeight; ilvl > 0; ilvl--{
        for p := self.head; p.nextLevel(ilvl) != nil; p = p.nextLevel(ilvl) {
            ret = append(ret, p.nextLevel(ilvl).String())
        }
        ret = append(ret, "\n")
    }
	return strings.Join(ret, ",")
}
