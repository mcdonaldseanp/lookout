package operation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mcdonaldseanp/lookout/sanitize"
)

// Definitions of the language paradigms.
type Operation interface {
	// Operation conflicts are checked by building
	// a hash and checking if the hashed value of the
	// Operation already exists.
	//
	// Some operations have multiple hash keys, because
	// they can conflict in multiple ways
	HashKeys() []string
	Empty() error
}

// OAR definitions
// Observations
// ---------------------------------------------------------------
type Observation struct {
	Entity   string `yaml:"entity" json:"entity"`
	Query    string `yaml:"query" json:"query"`
	Instance string `yaml:"instance" json:"instance"`
	Expect   string `yaml:"expect,omitempty" json:"expect,omitempty"`
}

type ObservationResult struct {
	Succeeded   bool        `yaml:"succeeded" json:"succeeded"`
	Result      string      `yaml:"result" json:"result"`
	Expected    bool        `yaml:"expected" json:"expected"`
	Logs        string      `yaml:"logs" json:"logs"`
	Observation Observation `yaml:"observation" json:"observation"`
}

type ObservationResults struct {
	Observations            map[string]ObservationResult `yaml:"observations" json:"observations"`
	Total_Observations      int                          `yaml:"total_observations" json:"total_observations"`
	Failed_Observations     int                          `yaml:"failed_observations" json:"failed_observations"`
	Unexpected_Observations int                          `yaml:"unexpected_observations" json:"unexpected_observations"`
}

// Observations can only conflict if
// 1. they expect something
// 2. they share all fields with another observation
//    but expect something different
//
// This makes checking for conflicts a pain because
// in all other cases for the other operations we want
// to check things via hash collision but for observations
// it's okay if two observations collide as long as they
// don't expect two different things.
//
// What we end up having to do is special case the
// collision check for observations so that if we find
// a collision we go the extra step to see if the two
// observations expect different things.
func (obsv Observation) HashKeys() []string {
	result := []string{}
	// Don't even return a hash key if there is no expect field
	if obsv.Expect != "" {
		hash := "OBS" + "EN" + sanitize.ReplaceAllSpaces(obsv.Entity) +
			"QU" + sanitize.ReplaceAllSpaces(obsv.Query) +
			"IN" + sanitize.ReplaceAllSpaces(obsv.Instance)
		result = append(result, hash)
	}
	return result
}

func (obsv Observation) Empty() error {
	if obsv.Entity == "" {
		return fmt.Errorf("missing entity")
	} else if obsv.Query == "" {
		return fmt.Errorf("missing query")
	} else if obsv.Instance == "" {
		return fmt.Errorf("missing instance")
	}
	return nil
}

// ---------------------------------------------------------------

// Actions
// ---------------------------------------------------------------
type Action struct {
	Path   string   `yaml:"path" json:"path"`
	Script string   `yaml:"script" json:"script"`
	Exe    string   `yaml:"exe,omitempty" json:"exe,omitempty"`
	Args   []string `yaml:"args,omitempty" json:"args,omitempty"`
}

type ActionResult struct {
	Succeeded bool   `yaml:"succeeded" json:"succeeded"`
	Output    string `yaml:"output" json:"output"`
	Logs      string `yaml:"logs" json:"logs"`
	Action    Action `yaml:"action" json:"action"`
}

type ActionResults struct {
	Actions map[string]ActionResult `json:"actions"`
}

func (actn Action) HashKeys() []string {
	// Actions can't conflict unless it's the name
	return []string{}
}

func (actn Action) Empty() error {
	if actn.Exe == "" {
		return fmt.Errorf("missing exe")
	}
	return nil
}

// ---------------------------------------------------------------

// Reactions
// ---------------------------------------------------------------
type Condition struct {
	Check string      `yaml:"check" json:"check"`
	Value interface{} `yaml:"value" json:"value"`
}

type Reaction struct {
	Observation string    `yaml:"observation" json:"observation"`
	Action      string    `yaml:"action" json:"action"`
	Condition   Condition `yaml:"condition" json:"condition"`
}

type ReactionResult struct {
	Succeeded bool     `yaml:"succeeded" json:"succeeded"`
	Skipped   bool     `yaml:"skipped" json:"skipped"`
	Output    string   `yaml:"output" json:"output"`
	Logs      string   `yaml:"logs" json:"logs"`
	Message   string   `yaml:"message" json:"message"`
	Reaction  Reaction `yaml:"reaction" json:"reaction"`
}

type ReactionResults struct {
	Reactions               map[string]ReactionResult    `yaml:"reactions" json:"reactions"`
	Observations            map[string]ObservationResult `yaml:"observations" json:"observations"`
	Total_Observations      int                          `yaml:"total_observations" json:"total_observations"`
	Failed_Observations     int                          `yaml:"failed_observations" json:"failed_observations"`
	Unexpected_Observations int                          `yaml:"unexpected_observations" json:"unexpected_observations"`
	Total_Reactions         int                          `yaml:"total_reactions" json:"total_reactions"`
	Failed_Reactions        int                          `yaml:"failed_reactions" json:"failed_reactions"`
	Skipped_Reactions       int                          `yaml:"skipped_reactions" json:"skipped_reactions"`
}

func (rctn Reaction) HashKeys() []string {
	// Reactions can't conflict unless it's the name
	return []string{}
}

func (rctn Reaction) Empty() error {
	if rctn.Observation == "" {
		return fmt.Errorf("missing observation")
	} else if rctn.Action == "" {
		return fmt.Errorf("missing action")
	} else if rctn.Condition.Check == "" {
		return fmt.Errorf("missing condition check")
	} else if rctn.Condition.Value == "" {
		return fmt.Errorf("missing condition value")
	}
	return nil
}

// ---------------------------------------------------------------

// Implements
// ---------------------------------------------------------------
type Correction struct {
	Entity      string   `yaml:"entity" json:"entity"`
	Query       string   `yaml:"query" json:"query"`
	Starts_From []string `yaml:"starts_from" json:"starts_from"`
	Results_In  string   `yaml:"results_in" json:"results_in"`
}

type ReactionImplement struct {
	Corrects Correction `yaml:"corrects,omitempty" json:"corrects,omitempty"`
	Args     []string   `yaml:"args" json:"args"`
}

type ObservationImplement struct {
	Entity string   `yaml:"entity" json:"entity"`
	Query  string   `yaml:"query" json:"query"`
	Args   []string `yaml:"args" json:"args"`
}

type Implement struct {
	Path        string               `yaml:"path,omitempty" json:"path,omitempty"`
	Script      string               `yaml:"script,omitempty" json:"script,omitempty"`
	Exe         string               `yaml:"exe,omitempty" json:"exe,omitempty"`
	Source_File string               `yaml:"source_file,omitempty" json:"source_file,omitempty"`
	Source_Url  string               `yaml:"source_url,omitempty" json:"source_url,omitempty"`
	Reacts      ReactionImplement    `yaml:"reacts,omitempty" json:"reacts,omitempty"`
	Observes    ObservationImplement `yaml:"observes,omitempty" json:"observes,omitempty"`
}

func emptyObserves(impl Implement) bool {
	if impl.Observes.Entity == "" ||
		impl.Observes.Query == "" ||
		impl.Observes.Args == nil {
		return true
	} else {
		return false
	}
}

func emptyReacts(impl Implement) bool {
	if impl.Reacts.Args == nil {
		return true
	} else {
		return false
	}
}

func emptyCorrects(impl Implement) bool {
	if impl.Reacts.Corrects.Entity == "" ||
		impl.Reacts.Corrects.Query == "" ||
		impl.Reacts.Corrects.Starts_From == nil ||
		impl.Reacts.Corrects.Results_In == "" {
		return true
	}
	return false
}

// Implements have very specific collision behavior:
//
// If the implement can react _and_ correct then it
// can conflict with other implements if they both
// attempt to correct the same thing. Implements
// that have the same value for all correction fields
// conflict.
//
// If the implement can observe then it conflicts with
// any other implement that observes the same entity/query
func (impl Implement) HashKeys() []string {
	result := []string{}
	if emptyReacts(impl) == false && emptyCorrects(impl) == false {
		// Sort the starts_from hash to ensure that the hash
		// we create can be checked against another hash
		// with the same values but in a different order
		var hash_starts string
		raw_starts := impl.Reacts.Corrects.Starts_From
		if raw_starts != nil {
			sort.Strings(raw_starts)
			hash_starts = strings.Join(raw_starts, "-")
		}
		react_hash := "IMPLRCT" + "EN" + sanitize.ReplaceAllSpaces(impl.Reacts.Corrects.Entity) +
			"QU" + sanitize.ReplaceAllSpaces(impl.Reacts.Corrects.Query) +
			"SF" + sanitize.ReplaceAllSpaces(hash_starts) +
			"RI" + sanitize.ReplaceAllSpaces(impl.Reacts.Corrects.Results_In)
		result = append(result, react_hash)
	}
	if emptyObserves(impl) == false {
		observe_hash := "IMPLOBS" + "EN" + sanitize.ReplaceAllSpaces(impl.Observes.Entity) +
			"QU" + sanitize.ReplaceAllSpaces(impl.Observes.Query)
		result = append(result, observe_hash)
	}
	return result
}

func (impl Implement) Empty() error {
	// impls must have one of:
	// * A complete source, including file and url
	// * A path
	// * A script
	if (len(impl.Source_File) < 1 && len(impl.Source_Url) < 1) && len(impl.Path) < 1 && len(impl.Script) < 1 {
		return fmt.Errorf("missing at least one of: path, script, source_file/url")
	}
	// Exe can _only_ be empty if Source_File is provided, since
	// the source file will substitute for the missing exe
	if len(impl.Exe) < 1 && len(impl.Source_File) < 1 {
		return fmt.Errorf("exe can only be empty with valid source_file/url")
	}
	// I'm not positive this is true, but can't think
	// of whether or not an implement can omit
	// both reacting and observing (I'm pretty
	// sure they are useless without that)
	if emptyReacts(impl) && emptyObserves(impl) {
		return fmt.Errorf("missing at least one of reacts, observes")
	}
	return nil
}

// ---------------------------------------------------------------

// Everything together
// ---------------------------------------------------------------
type Operations struct {
	Reactions    map[string]Reaction    `yaml:"reactions,omitempty" json:"reactions,omitempty"`
	Observations map[string]Observation `yaml:"observations,omitempty" json:"observations,omitempty"`
	Implements   map[string]Implement   `yaml:"implements,omitempty" json:"implements,omitempty"`
	Actions      map[string]Action      `yaml:"actions,omitempty" json:"actions,omitempty"`
}
