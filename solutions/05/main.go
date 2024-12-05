package main

import (
	"bufio"
	"errors"
	"flag"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	linkedlist "github.com/hmcalister/Go-DSA/list/LinkedList"
)

func main() {
	debugFlag := flag.Bool("debug", false, "Debug Flag")
	inputFilePath := flag.String("inputFile", "puzzleInput", "Path to input file.")
	selectedPart := flag.Int("part", 0, "Part to execute. Must be 1 or 2.")
	flag.Parse()
	logFileHandler := SetLogging(*debugFlag)
	defer logFileHandler.Close()

	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		slog.Error("error opening input file", "error", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	fileScanner := bufio.NewScanner(inputFile)
	if err != nil {
		slog.Error("error processing file input", "error", err)
		os.Exit(1)
	}

	var result int
	computationStartTime := time.Now()
	switch *selectedPart {
	case 1:
		result, err = Part01(fileScanner)
	case 2:
		result, err = Part02(fileScanner)
	default:
		slog.Error("invalid part selected, part must be one of 1 or 2", "part selected", *selectedPart)
		os.Exit(1)
	}
	computationEndTime := time.Now()
	if err != nil {
		slog.Error("error encountered during computation", "error", err, "part selected", *selectedPart)
		os.Exit(1)
	}

	slog.Info("computation completed", "result", result, "computation time elapsed (ns)", computationEndTime.Sub(computationStartTime).Nanoseconds())
}

// Defines a dependency graph --- key page depends on all pages in the value list and must occur *after* these\
func parsePageDependencyGraph(pageDependencies []string) map[int][]int {
	dependencyGraph := make(map[int][]int)
	for _, dependency := range pageDependencies {
		line := strings.TrimSpace(dependency)

		dependencyPages := strings.Split(line, "|")
		if len(dependencyPages) != 2 {
			slog.Error("found dependency line that is not of form 'a|b'", "offending line", dependency)
			continue
		}
		beforePage, err1 := strconv.Atoi(dependencyPages[0])
		afterPage, err2 := strconv.Atoi(dependencyPages[1])
		if errors.Join(err1, err2) != nil {
			slog.Error("found dependency line that is not of form '[int]|[int]'", "offending line", dependency)
		}

		// If afterPage is not already in the dependency graph, add a new list
		if _, ok := dependencyGraph[afterPage]; !ok {
			dependencyGraph[afterPage] = make([]int, 0)
		}
		dependencyGraph[afterPage] = append(dependencyGraph[afterPage], beforePage)
		slog.Debug("parsed dependency", "line", line, "updated dependency list", dependencyGraph[afterPage])
	}

	return dependencyGraph
}

func parseUpdateLine(updateLine []string) []int {
	updatePagesList := make([]int, len(updateLine))
	for i := 0; i < len(updateLine); i += 1 {
		updatePagesList[i], _ = strconv.Atoi(updateLine[i])
	}
	return updatePagesList
}

// Given a page's dependencies,
// the pages involved in this update (to filter out dependencies that do not matter),
// and the pages added to the output thus far,
// determine if the proposed page can be added to the list
func isPageValid(pageDependencies []int, updatePages []int, addedPages []int) bool {
	// For each page that is required to precede
	for _, precedingPage := range pageDependencies {
		// If that page is in the update and the page is not added
		if slices.Contains(updatePages, precedingPage) && !slices.Contains(addedPages, precedingPage) {
			return false
		}
	}
	return true
}

func Part01(fileScanner *bufio.Scanner) (int, error) {
	dependencies := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			break
		}
		dependencies = append(dependencies, line)
	}
	pageDependencyGraph := parsePageDependencyGraph(dependencies)
	slog.Debug("page dependency graph parsed", "dependency graph", pageDependencyGraph)

	middleNumbersSum := 0

updateValidationLoop:
	for fileScanner.Scan() {
		updateLineString := strings.Split(fileScanner.Text(), ",")

		// Pages involved in the update
		updatePages := parseUpdateLine(updateLineString)

		// Pages added to the update so far
		addedPages := make([]int, 0)
		slog.Debug("parsed update", "included pages", updatePages)

		// Walk over each page involved in the update in turn and ensure all dependencies are met
		// if not, the update is invalid, so move to the next line
		// if so, add page to the added pages list and continue to the next page
		for i, page := range updatePages {
			pageDependencies := pageDependencyGraph[page]
			slog.Debug("page validation loop", "page index", i, "page", page, "page dependencies", pageDependencies, "current added pages", addedPages)
			if !isPageValid(pageDependencies, updatePages, addedPages) {
				continue updateValidationLoop
			}
			addedPages = append(addedPages, page)
		}

		// All pages good, add middle number
		middleNumbersSum += addedPages[len(addedPages)/2]
		slog.Debug("good update found", "update list", updatePages, "new middle sum", middleNumbersSum)
	}

	return middleNumbersSum, nil
}

func Part02(fileScanner *bufio.Scanner) (int, error) {
	dependencies := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if len(line) == 0 {
			break
		}
		dependencies = append(dependencies, line)
	}
	pageDependencyGraph := parsePageDependencyGraph(dependencies)
	slog.Debug("page dependency graph parsed", "dependency graph", pageDependencyGraph)

	middleNumbersSum := 0

updateLoop:
	for fileScanner.Scan() {
		updateLineString := strings.Split(fileScanner.Text(), ",")

		// Pages involved in the update
		updatePages := parseUpdateLine(updateLineString)

		// The pages that are yet to be added, separate from updatePages
		// as this variable will be spliced out until it is eventually empty
		//
		// We use a linked list to easily remove items later
		// without having to reallocate potentially large arrays
		remainingPages := linkedlist.New[int]()
		for _, page := range updatePages {
			remainingPages.Add(page)
		}

		// Pages added to the update so far
		addedPages := make([]int, 0)

		// Try to add each number in turn, if that number can be added, do it and splice it out of the list
	addPagesLoop:
		for remainingPages.Length() > 0 {
			slog.Debug("add pages loop", "current added pages", addedPages, "remaining pages", remainingPages.Length())
			for i := 0; i < remainingPages.Length(); i += 1 {
				page, _ := remainingPages.ItemAtIndex(i)
				pageDependencies := pageDependencyGraph[page]
				if isPageValid(pageDependencies, updatePages, addedPages) {
					// We can add this page, do it
					addedPages = append(addedPages, page)
					remainingPages.RemoveAtIndex(i)
					continue addPagesLoop
				}
			}

			// If we have made it here, we have exhausted all options from the
			// pages to add and hence the update cannot be fixed
			slog.Debug("bad update cannot be fixed", "update pages", updatePages)
			continue updateLoop
		}

		if slices.Equal(updatePages, addedPages) {
			middleNumbersSum += addedPages[len(addedPages)/2]
			slog.Debug("bad update found (and fixed)", "fixed update list", addedPages, "updated middle sum", middleNumbersSum)
		} else {
			slog.Debug("update already good", "update list", addedPages)
		}

	}

	return middleNumbersSum, nil
}
