package score_extractor

import (
	"io/ioutil"
	"testing"
)

func Reader(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	return data
}

func TestScoreExtractor_Pacman(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     Pacman
	}{
		{
			scenario: "Pacman",
			filename: "./checkpoints_test_data/1615018620.ram",
			want: Pacman{
				Lives: 3,
				Level: 1,
				Score: 5270,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			got := ExtractPacman(Reader(t, tc.filename))

			if tc.want.Lives != got.Lives {
				t.Errorf("want lives %v, got %v", tc.want.Lives, got.Lives)
			}
			if tc.want.Level != got.Level {
				t.Errorf("want level %v, got %v", tc.want.Level, got.Level)
			}
			if tc.want.Score != got.Score {
				t.Errorf("want score %v, got %v", tc.want.Score, got.Score)
			}
		})
	}
}

func TestScoreExtractor_Tetris(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     Tetris
	}{
		{
			scenario: "Tetris",
			filename: "./checkpoints_test_data/1615022713.ram",
			want: Tetris{
				Level: 1,
				Score: 1475,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			got := ExtractTetris(Reader(t, tc.filename))

			if tc.want.Level != got.Level {
				t.Errorf("want level %v, got %v", tc.want.Level, got.Level)
			}
			if tc.want.Score != got.Score {
				t.Errorf("want score %v, got %v", tc.want.Score, got.Score)
			}
		})
	}
}
