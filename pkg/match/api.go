package match

import (
	"sync"

	"github.com/google/uuid"
)

type MatchSystem struct {
	mu         *sync.Mutex
	PersonList map[uuid.UUID]*SinglePerson
	matchTree  *node
}

func NewMatchSystem() MatchSystem {
	return MatchSystem{
		mu:         &sync.Mutex{},
		PersonList: map[uuid.UUID]*SinglePerson{},
		matchTree:  nil,
	}
}

func (m *MatchSystem) AddSinglePerson(name string, height uint, gender gender, wantDates int) SinglePerson {
	p := SinglePerson{
		ID:             uuid.New(),
		Name:           name,
		Height:         height,
		Gender:         gender,
		WantedDate:     wantDates,
		MatchPerson:    map[uuid.UUID]*SinglePerson{},
		NewMatchPerson: map[uuid.UUID]*SinglePerson{},
	}
	m.matchTree = insertNode(m.matchTree, &node{content: &p})
	m.PersonList[p.ID] = &p

	return p
}

func (m *MatchSystem) RemovePerson(id uuid.UUID) {
	p, ok := m.PersonList[id]
	if ok {
		m.matchTree = deleteNode(m.matchTree, &node{content: p})
		delete(m.PersonList, id)
	}
}

func (m *MatchSystem) removeFromMatchSystem(id uuid.UUID) {
	p := m.PersonList[id]
	m.matchTree = deleteNode(m.matchTree, &node{content: p})
}

func (m MatchSystem) MatchPerson(id uuid.UUID, want int) ([]*SinglePerson, int) {
	res := []*SinglePerson{}

	p := m.PersonList[id]
	limit := p.WantedDate - len(p.MatchPerson) - len(p.NewMatchPerson)
	if limit <= 0 {
		return res, 0
	}
	if want < limit {
		limit = want
	}
	nodes, remaining := findMatchesWithLimitAndCondition(m.matchTree, &node{content: p}, limit)
	for _, n := range nodes {
		res = append(res, n.content)
	}

	return res, remaining
}

func (m MatchSystem) RecordNewMatch(p *SinglePerson, NewMatch []*SinglePerson) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, np := range NewMatch {
		if np.WantedDate > (len(np.MatchPerson) + len(np.NewMatchPerson)) {
			p.NewMatchPerson[np.ID] = np
			np.NewMatchPerson[p.ID] = p
			// achieve WantedDate then removeFromMatchSystem
			if np.WantedDate == (len(np.MatchPerson) + len(np.NewMatchPerson)) {
				m.removeFromMatchSystem(np.ID)
			}
			if p.WantedDate == (len(p.MatchPerson) + len(p.NewMatchPerson)) {
				m.removeFromMatchSystem(p.ID)
				break
			}
		}
	}
}

func Display(p *SinglePerson, N int) []*SinglePerson {
	res := []*SinglePerson{}
	count := 0
	for _, v := range p.NewMatchPerson {
		count++
		res = append(res, v)

		// move from MatchPerson to NewMatchPerson
		p.MatchPerson[v.ID] = v
		delete(p.NewMatchPerson, v.ID)
		if count >= N {
			break
		}
	}
	return res
}
