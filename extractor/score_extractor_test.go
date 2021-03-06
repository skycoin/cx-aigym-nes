package score_extractor

import (
	"io/ioutil"
	"testing"
)

func TestScoreExtractor_Pacman(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     Pacman
	}{
		{
			scenario: "Pacman",
			filename: "../checkpoints/1615018620.ram",
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
				t.Errorf("want levels %v, got %v", tc.want.Level, got.Level)
			}
			if tc.want.Score != got.Score {
				t.Errorf("want score %v, got %v", tc.want.Score, got.Score)
			}
		})
	}
}

func Reader(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading file: %v", err)
	}
	return data
}
