package rules

import (
	_ "fmt"
	"github.com/camirmas/go_stop/models"
	"testing"
)

func TestRules(t *testing.T) {
	t.Run("Running the rules", func(t *testing.T) {
		t.Run("with basic liberties", testBasicLiberties)
		t.Run("with strings", testStrings)
		t.Run("self-capture violation", testSelfCaptureViolation)
		t.Run("valid self-capture", testValidSelfCapture)
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

func testStrings(t *testing.T) {
	// Example: String of white stones surrounded by enemy stones
	b1 := models.Stone{X: 1, Y: 0, Color: "black"}
	b2 := models.Stone{X: 2, Y: 0, Color: "black"}
	b3 := models.Stone{X: 3, Y: 0, Color: "black"}
	b4 := models.Stone{X: 4, Y: 1, Color: "black"}
	b5 := models.Stone{X: 4, Y: 2, Color: "black"}
	b6 := models.Stone{X: 3, Y: 3, Color: "black"}
	b7 := models.Stone{X: 2, Y: 2, Color: "black"}
	b8 := models.Stone{X: 0, Y: 1, Color: "black"}
	b9 := models.Stone{X: 0, Y: 2, Color: "black"}
	b10 := models.Stone{X: 0, Y: 3, Color: "black"}
	b11 := models.Stone{X: 1, Y: 4, Color: "black"}
	b12 := models.Stone{X: 2, Y: 4, Color: "black"}

	w1 := models.Stone{X: 1, Y: 1, Color: "white"}
	w2 := models.Stone{X: 2, Y: 1, Color: "white"}
	w3 := models.Stone{X: 3, Y: 1, Color: "white"}
	w4 := models.Stone{X: 3, Y: 2, Color: "white"}
	w5 := models.Stone{X: 2, Y: 3, Color: "white"}
	w6 := models.Stone{X: 1, Y: 3, Color: "white"}
	w7 := models.Stone{X: 1, Y: 2, Color: "white"}

	b := &models.Board{
		Size: models.SmallBoardSize,
		Stones: []models.Stone{
			b1,
			b2,
			b3,
			b4,
			b5,
			b6,
			b7,
			b8,
			b9,
			b10,
			b11,
			b12,
			w1,
			w2,
			w3,
			w4,
			w5,
			w6,
			w7,
		},
	}

	strings, _ := Run(b, b7)

	if len(strings) != 1 {
		t.Errorf("Expected 1 String, found %d", len(strings))
	}

	if len(strings[0]) != 7 {
		t.Errorf("Expected 7 Stones in String, found %d", len(strings[0]))
	}
}

func testSelfCaptureViolation(t *testing.T) {
	b1 := models.Stone{X: 0, Y: 0, Color: "black"}
	b2 := models.Stone{X: 1, Y: 0, Color: "black"}
	b3 := models.Stone{X: 2, Y: 0, Color: "black"}
	b4 := models.Stone{X: 0, Y: 1, Color: "black"}
	b5 := models.Stone{X: 0, Y: 2, Color: "black"}
	b6 := models.Stone{X: 0, Y: 3, Color: "black"}
	b7 := models.Stone{X: 1, Y: 2, Color: "black"}
	b8 := models.Stone{X: 2, Y: 1, Color: "black"}
	b9 := models.Stone{X: 2, Y: 2, Color: "black"}

	w1 := models.Stone{X: 1, Y: 1, Color: "white"}

	b := &models.Board{
		Size: models.SmallBoardSize,
		Stones: []models.Stone{
			b1,
			b2,
			b3,
			b4,
			b5,
			b6,
			b7,
			b8,
			b9,
			w1,
		},
	}

	_, err := Run(b, w1)

	expectedErr := selfCaptureError{}
	if err.Error() != expectedErr.Error() {
		t.Error("Expected selfCaptureError")
	}
}

func testValidSelfCapture(t *testing.T) {
	b1 := models.Stone{X: 1, Y: 0, Color: "black"}
	b2 := models.Stone{X: 1, Y: 1, Color: "black"}
	b3 := models.Stone{X: 1, Y: 2, Color: "black"}
	b4 := models.Stone{X: 2, Y: 2, Color: "black"}
	b5 := models.Stone{X: 3, Y: 2, Color: "black"}
	b6 := models.Stone{X: 3, Y: 1, Color: "black"}
	b7 := models.Stone{X: 3, Y: 0, Color: "black"}

	w1 := models.Stone{X: 0, Y: 0, Color: "white"}
	w2 := models.Stone{X: 0, Y: 1, Color: "white"}
	w3 := models.Stone{X: 0, Y: 2, Color: "white"}
	w4 := models.Stone{X: 1, Y: 3, Color: "white"}
	w5 := models.Stone{X: 2, Y: 3, Color: "white"}
	w6 := models.Stone{X: 3, Y: 3, Color: "white"}
	w7 := models.Stone{X: 4, Y: 2, Color: "white"}
	w8 := models.Stone{X: 4, Y: 1, Color: "white"}
	w9 := models.Stone{X: 4, Y: 0, Color: "white"}
	w10 := models.Stone{X: 2, Y: 1, Color: "white"}
	w11 := models.Stone{X: 2, Y: 0, Color: "white"}

	b := &models.Board{
		Size: models.SmallBoardSize,
		Stones: []models.Stone{
			b1,
			b2,
			b3,
			b4,
			b5,
			b6,
			b7,
			w1,
			w2,
			w3,
			w4,
			w5,
			w6,
			w7,
			w8,
			w9,
			w10,
			w11,
		},
	}

	strings, err := Run(b, w11)

	if err != nil {
		t.Errorf("Expected 1 String, got error: %s", err.Error())
	}

	if len(strings) != 1 {
		t.Errorf("Expected 1 String, got %d", len(strings))
	}

	if len(strings[0]) != 7 {
		t.Errorf("Expected String with 7 stones, got %d", len(strings[0]))
	}
}
