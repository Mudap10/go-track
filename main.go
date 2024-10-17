package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	colorReset         = "\033[0m"
	colorSkyBlue       = "\033[38;5;117m"  // Light sky blue
	colorDeepSkyBlue   = "\033[38;5;39m"   // Deep sky blue
	colorSunYellow     = "\033[38;5;220m"  // Bright sun yellow
	colorCloudWhite    = "\033[38;5;255m"  // Bright white for clouds
	colorGrass         = "\033[38;5;40m"   // Bright green for grass
	colorSunsetOrange  = "\033[38;5;208m"  // Sunset orange
	colorNightPurple   = "\033[38;5;93m"   // Night sky purple
	colorUserPrompt    = "\033[38;5;51m"   // Cyan for user prompts
	colorUsernamePrompt = "\033[38;5;75m"  // Medium blue for username prompt
	colorRed           = "\033[38;5;196m" // Bright red for errors
)

type Task struct {
	Description string `json:"description"`
	Priority    bool   `json:"priority"`
}

func main() {
	fmt.Print(colorUsernamePrompt + "Username: " + colorReset)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name := scanner.Text()

	fileName := name + ".json"
	tasks := loadTasks(fileName)

	fmt.Println(colorCloudWhite + "\nWelcome to your task manager!" + colorReset)
	fmt.Println(colorCloudWhite + "Available commands:" + colorReset)
	fmt.Println(colorDeepSkyBlue + "a" + colorReset + " - Add a new task")
	fmt.Println(colorDeepSkyBlue + "r" + colorReset + " - Remove a task")
	fmt.Println(colorDeepSkyBlue + "e" + colorReset + " - Exit the program")

	fmt.Println(colorCloudWhite + "\nCurrent tasks:" + colorReset)
	printTasks(tasks)

	for {
		fmt.Print(colorUserPrompt + "\nWhat would you like to do? (a/r/e): " + colorReset)
		scanner.Scan()
		action := strings.ToLower(scanner.Text())

		switch action {
		case "e":
			saveTasks(fileName, tasks)
			fmt.Println(colorDeepSkyBlue + "Goodbye!" + colorReset)
			return
		case "a":
			addTask(scanner, &tasks)
		case "r":
			removeTask(scanner, &tasks)
		default:
			fmt.Println(colorRed + "Invalid action. Please try again." + colorReset)
		}

		fmt.Println(colorSunYellow + "\nUpdated tasks:" + colorReset)
		printTasks(tasks)
	}
}

func loadTasks(fileName string) []Task {
	file, err := os.Open(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}
		}
		fmt.Println(colorSunsetOrange + "Error opening file: " + err.Error() + colorReset)
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ". ", 2)
		if len(parts) != 2 {
			continue
		}
		description := strings.TrimSpace(parts[1])
		priority := false
		if strings.HasPrefix(description, "!") {
			priority = true
			description = strings.TrimPrefix(description, "!")
		} else {
			description = strings.TrimPrefix(description, "  ")
		}
		tasks = append(tasks, Task{Description: description, Priority: priority})
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(colorSunsetOrange + "Error reading file: " + err.Error() + colorReset)
	}

	return tasks
}

func saveTasks(fileName string, tasks []Task) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(colorSunsetOrange + "Error creating file: " + err.Error() + colorReset)
		return
	}
	defer file.Close()

	for i, task := range tasks {
		prefix := "  "
		if task.Priority {
			prefix = "!"
		}
		line := fmt.Sprintf("%d. %s%s\n", i+1, prefix, task.Description)
		file.WriteString(line)
	}
}

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println(colorSunYellow + "No tasks." + colorReset)
		return
	}
	for i, task := range tasks {
		if task.Priority {
			fmt.Printf(colorRed+"%d. %s\n"+colorReset, i+1, task.Description)
		} else {
			fmt.Printf(colorSkyBlue+"%d. %s\n"+colorReset, i+1, task.Description)
		}
	}
}

func addTask(scanner *bufio.Scanner, tasks *[]Task) {
	fmt.Print(colorSunsetOrange + "Enter new task description (! for prio): " + colorReset)
	scanner.Scan()
	input := scanner.Text()

	priority := false
	description := input

	if strings.HasPrefix(input, "!") {
		priority = true
		description = strings.TrimPrefix(input, "!")
	}

	*tasks = append(*tasks, Task{Description: description, Priority: priority})
	fmt.Println(colorNightPurple + "Task added successfully." + colorReset)
}

func removeTask(scanner *bufio.Scanner, tasks *[]Task) {
	if len(*tasks) == 0 {
		fmt.Println(colorSunYellow + "No tasks to remove." + colorReset)
		return
	}

	fmt.Print(colorUserPrompt + "Enter the number of the task to remove: " + colorReset)
	scanner.Scan()
	var index int
	fmt.Sscanf(scanner.Text(), "%d", &index)

	if index < 1 || index > len(*tasks) {
		fmt.Println(colorSunsetOrange + "Invalid task number." + colorReset)
		return
	}

	*tasks = append((*tasks)[:index-1], (*tasks)[index:]...)
	fmt.Println(colorGrass + "Task removed successfully." + colorReset)
}
