/*
Package dir_tree provides a function to generate and print a tree of directories and files in a given path to the console.
*/
package dir_tree

import (
	"fmt"
	"os"
	"path/filepath"
)

// Print prints a tree of directories and files in the current working directory to the console.
func Print() {
	fmt.Println()
	filepath.Walk(".", func(name string, info os.FileInfo, err error) error {
		if info.IsDir() {
			fmt.Println(name)
		} else {

			fmt.Println("   " + name)
		}
		return nil
	})
	fmt.Println()
}
