package skills

import (
	"fmt"
	"sync"
	"testing"
)

func createMockSkill(id, name, cat string) *Skill {
	return &Skill{
		ID:          id,
		Name:        name,
		Category:    cat,
		Description: "Mock skill description for " + name,
		Triggers:    []string{id, name},
		Tools: []ToolDefinition{
			{
				Name:        id + "_tool",
				Description: "Tool for " + name,
				Parameters: ToolParameters{
					Type: "object",
					Properties: map[string]ParameterProperty{
						"input": {Type: "string"},
					},
					Required: []string{"input"},
				},
				Category: cat,
			},
		},
		Instructions: "Instructions for " + name,
		MerkleSeal:   "seal_" + id,
	}
}

func TestRegistryBasicOperations(t *testing.T) {
	reg := NewSkillRegistry()

	s1 := createMockSkill("skill_1", "Skill One", "DevOps")
	s2 := createMockSkill("skill_2", "Skill Two", "Finance")

	reg.Register(s1)
	reg.Register(s2)

	if reg.Count() != 2 {
		t.Errorf("expected 2 skills, got %d", reg.Count())
	}

	retrieved, found := reg.GetByID("skill_1")
	if !found || retrieved.Name != "Skill One" {
		t.Errorf("expected to find 'Skill One' by ID")
	}

	byName, found := reg.GetByName("skill two")
	if !found || byName.ID != "skill_2" {
		t.Errorf("expected to find 'skill_2' by name")
	}

	tool, parentSkill, toolFound := reg.GetToolByName("skill_1_tool")
	if !toolFound || tool.Name != "skill_1_tool" || parentSkill.ID != "skill_1" {
		t.Errorf("expected to find tool and parent skill")
	}

	devopsSkills := reg.GetByCategory("DevOps")
	if len(devopsSkills) != 1 || devopsSkills[0].ID != "skill_1" {
		t.Errorf("expected 1 DevOps skill")
	}

	allTools := reg.GetAllTools()
	if len(allTools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(allTools))
	}

	seal := reg.GlobalMerkleSeal()
	if seal == "" {
		t.Errorf("expected non-empty global seal")
	}

	// Test unregister
	if !reg.Unregister("skill_1") {
		t.Errorf("expected skill_1 to be unregistered")
	}
	if reg.Count() != 1 {
		t.Errorf("expected count to be 1 after unregister")
	}
	if _, found := reg.GetByID("skill_1"); found {
		t.Errorf("expected skill_1 to not exist after unregister")
	}
}

func TestRegistryConcurrency(t *testing.T) {
	reg := NewSkillRegistry()
	var wg sync.WaitGroup

	numWorkers := 20
	skillsPerWorker := 10

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < skillsPerWorker; i++ {
				id := fmt.Sprintf("worker_%d_skill_%d", workerID, i)
				name := fmt.Sprintf("Worker %d Skill %d", workerID, i)
				reg.Register(createMockSkill(id, name, "Parallel"))
			}
		}(w)
	}

	wg.Wait()

	expectedCount := numWorkers * skillsPerWorker
	if reg.Count() != expectedCount {
		t.Errorf("expected %d skills under concurrent registration, got %d", expectedCount, reg.Count())
	}
}

func TestRegistryExportMCP(t *testing.T) {
	reg := NewSkillRegistry()
	reg.Register(createMockSkill("sec_edgar", "SEC Edgar", "Finance"))

	mcpBytes, err := reg.ExportMCPTools()
	if err != nil {
		t.Fatalf("unexpected error exporting MCP tools: %v", err)
	}

	if len(mcpBytes) == 0 {
		t.Errorf("expected non-empty MCP export bytes")
	}
}
