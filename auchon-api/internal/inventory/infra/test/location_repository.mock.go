// Package inventorytest holds the inventory context's test doubles.
package inventorytest

import (
	"context"
	"testing"

	"github.com/auchon/auchon-api/internal/inventory/domain"
	"github.com/auchon/auchon-api/lib/test"
)

// LocationRepositoryMock is an in-memory LocationRepository that records saves.
type LocationRepositoryMock struct {
	locations map[string]domain.LocationDTO
	saved     []domain.LocationDTO
	Fail      error
}

func NewLocationRepositoryMock() *LocationRepositoryMock {
	return &LocationRepositoryMock{locations: map[string]domain.LocationDTO{}}
}

// Reset empties the repository, forgets the saves and clears Fail, for the next scenario.
func (m *LocationRepositoryMock) Reset() {
	m.locations = map[string]domain.LocationDTO{}
	m.saved = nil
	m.Fail = nil
}

// WithLocation seeds an existing location (not recorded as saved).
func (m *LocationRepositoryMock) WithLocation(locations ...domain.LocationDTO) *LocationRepositoryMock {
	for _, l := range locations {
		m.locations[l.ID] = l
	}
	return m
}

func (m *LocationRepositoryMock) Save(_ context.Context, location domain.LocationDTO) error {
	if m.Fail != nil {
		return m.Fail
	}
	m.saved = append(m.saved, location)
	m.locations[location.ID] = location
	return nil
}

func (m *LocationRepositoryMock) FindByID(_ context.Context, id string) (*domain.LocationDTO, error) {
	if m.Fail != nil {
		return nil, m.Fail
	}
	l, ok := m.locations[id]
	if !ok {
		return nil, nil
	}
	return &l, nil
}

func (m *LocationRepositoryMock) AncestorIDs(_ context.Context, id string) ([]string, error) {
	if m.Fail != nil {
		return nil, m.Fail
	}
	var ids []string
	for current, ok := m.locations[id]; ok && current.ParentID != ""; current, ok = m.locations[current.ParentID] {
		ids = append(ids, current.ParentID)
		if len(ids) > len(m.locations) {
			break // seeded data with a cycle
		}
	}
	return ids, nil
}

func (m *LocationRepositoryMock) AssertSaved(t testing.TB, want domain.LocationDTO) {
	t.Helper()
	if len(m.saved) == 0 {
		t.Fatalf("nothing saved, want %+v", want)
	}
	test.AssertEqual(t, m.saved[len(m.saved)-1], want)
}

func (m *LocationRepositoryMock) AssertNothingSaved(t testing.TB) {
	t.Helper()
	if len(m.saved) != 0 {
		t.Errorf("saved %+v, want nothing", m.saved)
	}
}
