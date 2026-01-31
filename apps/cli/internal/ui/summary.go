package ui

import "fmt"

func PrintSummary(path string) {
	fmt.Printf("\n✅ Project created at %s\n", path)
	fmt.Println("Next steps:")
	fmt.Println("  1) cd into the project")
	fmt.Println("  2) review generated .env files")
	fmt.Println("  3) follow the project README")
}
