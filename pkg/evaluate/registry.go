package evaluate

import (
	"fmt"
	"sort"
	"sync"
)

const defaultCategoryOrderStep = 100

// CategorySpec describes where a check is shown and its display metadata.
type CategorySpec struct {
	ID              string
	Title           string
	DefaultProfiles []string
	Order           int
}

type registeredCheck struct {
	check    Check
	profiles map[string]struct{}
	order    int
}

type categoryRegistration struct {
	spec           CategorySpec
	checks         []*registeredCheck
	nextCheckOrder int
}

type registry struct {
	mu         sync.RWMutex
	categories map[string]*categoryRegistration
}

var (
	globalRegistry = newRegistry()
	profileTitles  = map[string]string{
		ProfileBasic:      "Baseline information gathering",
		ProfileExtended:   "Extended reconnaissance",
		ProfileAdditional: "Additional evaluation checks",
	}
)

func newRegistry() *registry {
	return &registry{
		categories: make(map[string]*categoryRegistration),
	}
}

// RegisterCheck registers a check with the given category and optional profile list.
func RegisterCheck(category CategorySpec, check Check, profiles ...string) {
	if err := globalRegistry.register(category, check, profiles...); err != nil {
		panic(err)
	}
}

// RegisterSimpleCheck registers a check backed by a parameterless function.
func RegisterSimpleCheck(category CategorySpec, id, title string, fn func(), profiles ...string) {
	RegisterCheck(category, Check{
		ID:    id,
		Title: title,
		Run: func(*Context) error {
			fn()
			return nil
		},
	}, profiles...)
}

// RegisterContextCheck registers a check backed by a function that consumes Context.
func RegisterContextCheck(category CategorySpec, id, title string, fn CheckFunc, profiles ...string) {
	RegisterCheck(category, Check{
		ID:    id,
		Title: title,
		Run:   fn,
	}, profiles...)
}

func (r *registry) register(spec CategorySpec, check Check, profiles ...string) error {
	if spec.ID == "" {
		return fmt.Errorf("category ID cannot be empty")
	}
	if check.ID == "" {
		return fmt.Errorf("check ID cannot be empty")
	}
	if check.Title == "" {
		return fmt.Errorf("check title cannot be empty")
	}
	if check.Run == nil {
		return fmt.Errorf("check %s has no runnable function", check.ID)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.categories[spec.ID]
	if !ok {
		if spec.Order == 0 {
			spec.Order = r.nextCategoryOrder()
		}
		entry = &categoryRegistration{
			spec: spec,
		}
		r.categories[spec.ID] = entry
	} else {
		if entry.spec.Title != spec.Title {
			return fmt.Errorf("category %s already registered with title %q (new title: %q)", spec.ID, entry.spec.Title, spec.Title)
		}
		if entry.spec.Order == 0 && spec.Order != 0 {
			entry.spec.Order = spec.Order
		}
		if len(entry.spec.DefaultProfiles) == 0 && len(spec.DefaultProfiles) > 0 {
			entry.spec.DefaultProfiles = append([]string(nil), spec.DefaultProfiles...)
		}
	}

	activeProfiles := profiles
	if len(activeProfiles) == 0 {
		activeProfiles = entry.spec.DefaultProfiles
	}
	if len(activeProfiles) == 0 {
		return fmt.Errorf("no profiles specified for check %s in category %s", check.ID, spec.ID)
	}

	newCheck := &registeredCheck{
		check:    check,
		profiles: make(map[string]struct{}, len(activeProfiles)),
		order:    entry.nextCheckOrder,
	}
	entry.nextCheckOrder++

	for _, profileID := range activeProfiles {
		newCheck.profiles[profileID] = struct{}{}
	}

	entry.checks = append(entry.checks, newCheck)
	return nil
}

func (r *registry) nextCategoryOrder() int {
	maxOrder := 0
	for _, entry := range r.categories {
		if entry.spec.Order > maxOrder {
			maxOrder = entry.spec.Order
		}
	}
	return maxOrder + defaultCategoryOrderStep
}

func (r *registry) profiles() []Profile {
	r.mu.RLock()
	defer r.mu.RUnlock()

	builders := make(map[string]*profileBuilder)

	for _, entry := range r.categories {
		for _, regCheck := range entry.checks {
			for profileID := range regCheck.profiles {
				builder := builders[profileID]
				if builder == nil {
					builder = newProfileBuilder(profileID)
					builders[profileID] = builder
				}
				builder.addCheck(entry.spec, regCheck)
			}
		}
	}

	out := make([]Profile, 0, len(builders))
	for _, builder := range builders {
		out = append(out, builder.build())
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out
}

type profileBuilder struct {
	id         string
	title      string
	categories map[string]*categoryBuilder
}

func newProfileBuilder(id string) *profileBuilder {
	title, ok := profileTitles[id]
	if !ok {
		title = id
	}
	return &profileBuilder{
		id:         id,
		title:      title,
		categories: make(map[string]*categoryBuilder),
	}
}

func (p *profileBuilder) addCheck(spec CategorySpec, regCheck *registeredCheck) {
	builder := p.categories[spec.ID]
	if builder == nil {
		builder = &categoryBuilder{
			spec: spec,
		}
		p.categories[spec.ID] = builder
	}
	builder.checks = append(builder.checks, checkWithOrder{
		check: regCheck.check,
		order: regCheck.order,
	})
}

func (p *profileBuilder) build() Profile {
	builders := make([]*categoryBuilder, 0, len(p.categories))
	for _, cat := range p.categories {
		builders = append(builders, cat)
	}
	sort.Slice(builders, func(i, j int) bool {
		orderA := builders[i].spec.Order
		orderB := builders[j].spec.Order
		if orderA == orderB {
			return builders[i].spec.Title < builders[j].spec.Title
		}
		return orderA < orderB
	})

	categories := make([]Category, len(builders))
	for idx, cat := range builders {
		categories[idx] = cat.build()
	}

	return Profile{
		ID:         p.id,
		Title:      p.title,
		Categories: categories,
	}
}

type categoryBuilder struct {
	spec   CategorySpec
	checks []checkWithOrder
}

type checkWithOrder struct {
	check Check
	order int
}

func (c *categoryBuilder) build() Category {
	sort.SliceStable(c.checks, func(i, j int) bool {
		return c.checks[i].order < c.checks[j].order
	})

	checks := make([]Check, 0, len(c.checks))
	for _, ch := range c.checks {
		checks = append(checks, ch.check)
	}

	return Category{
		ID:     c.spec.ID,
		Title:  c.spec.Title,
		Checks: checks,
	}
}

func defaultProfiles() []Profile {
	return globalRegistry.profiles()
}
