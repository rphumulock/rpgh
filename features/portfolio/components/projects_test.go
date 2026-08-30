package components

import "testing"

// TestTechNamesResolve keeps Project.Tech and Stack in sync: a name that does
// not match a StackItem still renders a filter button, but it can never be
// reached from the stack tab, which is a silent dead end.
func TestTechNamesResolve(t *testing.T) {
	known := map[string]bool{}
	for _, cat := range Stack {
		for _, item := range cat.Items {
			known[item.Name] = true
		}
	}

	for _, p := range Projects {
		for _, tech := range p.Tech {
			if !known[tech] {
				t.Errorf("project %q lists tech %q, which no Stack category declares", p.Title, tech)
			}
		}
	}
}

// TestEveryProjectIsReachableFromStack checks that each project can be found
// by clicking something on the stack tab -- the only way to set a filter now.
func TestEveryProjectIsReachableFromStack(t *testing.T) {
	clickable := map[string]bool{}
	for _, cat := range Stack {
		for _, item := range cat.Items {
			if TechUseCount(item.Name) > 0 {
				clickable[item.Name] = true
			}
		}
	}

	for _, p := range Projects {
		if len(p.Tech) == 0 {
			t.Errorf("project %q lists no tech, so the stack tab can never surface it", p.Title)
			continue
		}
		found := false
		for _, tech := range p.Tech {
			if clickable[tech] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("project %q has no clickable tech on the stack tab", p.Title)
		}
	}
}
