package evaluate

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/cdk-team/CDK/pkg/util"
)

const (
	ProfileBasic      = "basic"
	ProfileExtended   = "extended"
	ProfileAdditional = "additional"
)

// Context carries shared dependencies for evaluation checks.
type Context struct {
	Logger *log.Logger
}

// NewContext constructs a Context instance with a default logger when none is provided.
func NewContext(logger *log.Logger) *Context {
	if logger == nil {
		logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return &Context{Logger: logger}
}

// CheckFunc represents the executable unit for a security check.
type CheckFunc func(*Context) error

// Check describes an actionable evaluation task.
type Check struct {
	ID          string
	Title       string
	Description string
	Run         CheckFunc
}

func (c Check) execute(ctx *Context) error {
	if c.Run == nil {
		return nil
	}
	return c.Run(ctx)
}

// Category groups related checks under a shared heading.
type Category struct {
	ID     string
	Title  string
	Checks []Check
}

func (c Category) run(ctx *Context) {
	util.PrintH2(c.Title)
	logger := loggerFromContext(ctx)
	for _, check := range c.Checks {
		if err := check.execute(ctx); err != nil {
			logger.Printf("check %s failed: %v", readableCheckLabel(check), err)
		}
	}
}

// Profile combines categories into a runnable unit.
type Profile struct {
	ID         string
	Title      string
	Categories []Category
}

func (p Profile) run(ctx *Context) {
	for _, category := range p.Categories {
		category.run(ctx)
	}
}

// Evaluator coordinates profile registration and execution.
type Evaluator struct {
	profiles map[string]Profile
}

// NewEvaluator returns an Evaluator with the default profiles registered.
func NewEvaluator() *Evaluator {
	e := &Evaluator{profiles: make(map[string]Profile)}
	for _, profile := range defaultProfiles() {
		e.RegisterProfile(profile)
	}
	return e
}

// RegisterProfile adds or replaces a profile definition.
func (e *Evaluator) RegisterProfile(profile Profile) {
	if e.profiles == nil {
		e.profiles = make(map[string]Profile)
	}
	e.profiles[profile.ID] = profile
}

// Profile returns a copy of the profile and a boolean indicating whether it exists.
func (e *Evaluator) Profile(id string) (Profile, bool) {
	profile, ok := e.profiles[id]
	return profile, ok
}

// Profiles returns the registered profiles sorted by their identifier.
func (e *Evaluator) Profiles() []Profile {
	out := make([]Profile, 0, len(e.profiles))
	for _, profile := range e.profiles {
		out = append(out, profile)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

// RunProfile executes every category within the selected profile.
func (e *Evaluator) RunProfile(id string, ctx *Context) error {
	profile, ok := e.profiles[id]
	if !ok {
		return fmt.Errorf("unknown profile %q", id)
	}
	if ctx == nil {
		ctx = NewContext(nil)
	}
	profile.run(ctx)
	return nil
}

func loggerFromContext(ctx *Context) *log.Logger {
	if ctx != nil && ctx.Logger != nil {
		return ctx.Logger
	}
	return log.Default()
}

func readableCheckLabel(check Check) string {
	if check.ID != "" {
		return fmt.Sprintf("%s (%s)", check.Title, check.ID)
	}
	return check.Title
}
