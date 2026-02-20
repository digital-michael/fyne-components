package main

// TaskData represents a hierarchical task for the demo
type TaskData struct {
	ID       int
	Name     string
	Status   string
	Priority int
	Depth    int  // Indentation level (0 = root, 1 = child, 2 = grandchild)
	External bool // true = external task, false = internal task
}

// generateMockTasks creates sample hierarchical tasks for the demo
func generateMockTasks() []TaskData {
	return []TaskData{
		// Root level tasks (Depth = 0)
		{ID: 1, Name: "Project Alpha", Status: "In Progress", Priority: 1, Depth: 0, External: false},
		{ID: 2, Name: "Planning Phase", Status: "Complete", Priority: 2, Depth: 1, External: false},
		{ID: 3, Name: "Requirements Gathering", Status: "Complete", Priority: 3, Depth: 2, External: false},
		{ID: 4, Name: "Design Phase", Status: "In Progress", Priority: 2, Depth: 1, External: false},
		{ID: 5, Name: "UI Mockups", Status: "In Progress", Priority: 3, Depth: 2, External: true},
		{ID: 6, Name: "Database Schema", Status: "Complete", Priority: 3, Depth: 2, External: false},
		{ID: 7, Name: "Development Phase", Status: "Not Started", Priority: 2, Depth: 1, External: false},
		{ID: 8, Name: "Backend Development", Status: "Not Started", Priority: 3, Depth: 2, External: false},
		{ID: 9, Name: "Frontend Development", Status: "Not Started", Priority: 3, Depth: 2, External: true},

		// Second root task
		{ID: 10, Name: "Project Beta", Status: "Not Started", Priority: 1, Depth: 0, External: false},
		{ID: 11, Name: "Initial Research", Status: "Pending", Priority: 2, Depth: 1, External: true},
		{ID: 12, Name: "Feasibility Study", Status: "Not Started", Priority: 3, Depth: 2, External: true},

		// Third root task
		{ID: 13, Name: "Infrastructure Updates", Status: "In Progress", Priority: 1, Depth: 0, External: false},
		{ID: 14, Name: "Server Upgrades", Status: "In Progress", Priority: 2, Depth: 1, External: true},
		{ID: 15, Name: "Network Optimization", Status: "Complete", Priority: 2, Depth: 1, External: false},

		// Fourth root task
		{ID: 16, Name: "Documentation", Status: "Pending", Priority: 1, Depth: 0, External: false},
		{ID: 17, Name: "User Manual", Status: "Pending", Priority: 2, Depth: 1, External: true},
		{ID: 18, Name: "API Documentation", Status: "Not Started", Priority: 2, Depth: 1, External: false},
		{ID: 19, Name: "Examples and Tutorials", Status: "Not Started", Priority: 3, Depth: 2, External: false},
	}
}
