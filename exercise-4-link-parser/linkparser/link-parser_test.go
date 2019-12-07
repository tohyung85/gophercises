package linkparser

import (
	"os"
	"testing"
)

func TestProcessNodeExample1(t *testing.T) {
	links, err := parseHTMLFile("../data/ex1.html")
	if err != nil {
		t.Fatalf("Parsing failed with errors %s", err)
	}

	expectedOutput := []Link{
		Link{"/other-page", "A link to another page"},
	}

	runChecks(expectedOutput, links, t)
}

func TestProcessNodeExample2(t *testing.T) {
	links, err := parseHTMLFile("../data/ex2.html")
	if err != nil {
		t.Fatalf("Parsing failed with errors %s", err)
	}
	expectedOutput := []Link{
		Link{"https://www.twitter.com/joncalhoun", "Check me out on twitter"},
		Link{"https://github.com/gophercises", "Gophercises is on Github!"},
	}

	runChecks(expectedOutput, links, t)
}

func TestProcessNodeExample3(t *testing.T) {
	links, err := parseHTMLFile("../data/ex3.html")
	if err != nil {
		t.Fatalf("Parsing failed with errors %s", err)
	}

	expectedOutput := []Link{
		Link{"#", "Login"},
		Link{"/lost", "Lost? Need help?"},
		Link{"https://twitter.com/marcusolsson", "@marcusolsson"},
	}

	runChecks(expectedOutput, links, t)
}

func TestProcessNodeExample4(t *testing.T) {
	links, err := parseHTMLFile("../data/ex4.html")
	if err != nil {
		t.Fatalf("Parsing failed with errors %s", err)
	}
	expectedOutput := []Link{
		Link{"/dog-cat", "dog cat"},
	}

	runChecks(expectedOutput, links, t)
}

func parseHTMLFile(path string) ([]Link, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	linkParser := NewParser()

	links, err := linkParser.ParseHTML(file)

	return links, err
}

// Shared Test Functions

func runChecks(expect []Link, actual []Link, t *testing.T) {
	passed := checkLength(len(expect), len(actual), t)
	if passed {
		checkContents(expect, actual, t)
	}
}

func checkLength(expected int, actual int, t *testing.T) bool {
	if expected == actual {
		t.Logf("Sucess: %d link(s) in HTML", actual)
		return true
	} else {
		t.Errorf("Failed: Expected %d vs result of %d link(s)", expected, actual)
	}

	return false
}

func checkContents(expected []Link, actual []Link, t *testing.T) {
	for idx, eLink := range expected {
		actLink := actual[idx]
		if actLink.Href != eLink.Href {
			t.Errorf("Failed: Expected Href: %s vs result of %s", eLink.Href, actLink.Href)
		} else {
			t.Logf("Success: Href of %s", actLink.Href)
		}
		if actLink.Text != eLink.Text {
			t.Errorf("Failed: Expected Text: %s vs result of %s", eLink.Text, actLink.Text)
		} else {
			t.Logf("Success: Text of %s", actLink.Text)
		}
	}
}
