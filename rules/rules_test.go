package rules

import (
	"github.com/camirmas/go_stop/models"
	"testing"
)

func TestRules(t *testing.T) {
	t.Run("Running the rules", func(t *testing.T) {
		t.Run("with basic liberties", testBasicLiberties)
	})
}

func testBasicLiberties(t *testing.T) {
	s1 := models.Stone{X: 1, Y: 2, Color: "white"}
	s2 := models.Stone{X: 2, Y: 1, Color: "white"}
	s3 := models.Stone{X: 0, Y: 1, Color: "white"}
	s4 := models.Stone{X: 1, Y: 0, Color: "white"}
	// No liberties
	s5 := models.Stone{X: 1, Y: 1, Color: "black"}

	// Example: corner stone surrounded by enemy stones on two sides
	// No liberties
	s6 := models.Stone{X: 12, Y: 0, Color: "white"}
	s7 := models.Stone{X: 12, Y: 1, Color: "black"}
	s8 := models.Stone{X: 11, Y: 0, Color: "black"}

	b := &models.Board{
		Size:   models.SmallBoardSize,
		Stones: []models.Stone{s1, s2, s3, s4, s5, s6, s7, s8},
	}

	strings, _ := Run(b, s4)

	if len(strings) != 2 {
		t.Errorf("Expected 2 Strings, found %d", len(strings))
	}

	expected := []String{String{s5}, String{s6}}

	for i, str := range strings {
		e := expected[i]
		for j, stone := range str {
			if e[j] != stone {
				t.Errorf("Expected removed String to be %v, found %v", e[j], stone)
			}
		}
	}
}
