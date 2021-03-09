package score_extractor

import (
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
			got := ExtractPacman(Reader(tc.filename))

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
			got := ExtractTetris(Reader(tc.filename))

			if tc.want.Level != got.Level {
				t.Errorf("want level %v, got %v", tc.want.Level, got.Level)
			}
			if tc.want.Score != got.Score {
				t.Errorf("want score %v, got %v", tc.want.Score, got.Score)
			}
		})
	}
}

func TestScoreExtractor_SuperMarioBros(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     SuperMarioBros
	}{
		{
			scenario: "Super Mario Bros 2-2",
			filename: "./checkpoints_test_data/1615043111.ram",
			want: SuperMarioBros{
				Lives:      1,
				World:      2,
				Level:      2,
				MarioScore: 43750,
				LuigiScore: 0,
			},
		},
		{
			scenario: "Super Mario Bros 1-2",
			filename: "./checkpoints_test_data/1615041930.ram",
			want: SuperMarioBros{
				Lives:      1,
				World:      1,
				Level:      2,
				MarioScore: 20250,
				LuigiScore: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			got := ExtractSuperMarioBros(Reader(tc.filename))

			if tc.want.Lives != got.Lives {
				t.Errorf("want lives %v, got %v", tc.want.Lives, got.Lives)
			}
			if tc.want.World != got.World {
				t.Errorf("want world %v, got %v", tc.want.World, got.World)
			}
			if tc.want.Level != got.Level {
				t.Errorf("want level %v, got %v", tc.want.Level, got.Level)
			}
			if tc.want.MarioScore != got.MarioScore {
				t.Errorf("want mario score %v, got %v", tc.want.MarioScore, got.MarioScore)
			}
			if tc.want.LuigiScore != got.LuigiScore {
				t.Errorf("want luigi score %v, got %v", tc.want.LuigiScore, got.LuigiScore)
			}
		})
	}
}

func TestScoreExtractor_Bomberman(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     Bomberman
	}{
		{
			scenario: "Bomberman",
			filename: "./checkpoints_test_data/1615047878.ram",
			want: Bomberman{
				Lives: 2,
				Level: 1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			got := ExtractBomberman(Reader(tc.filename))

			if tc.want.Lives != got.Lives {
				t.Errorf("want lives %v, got %v", tc.want.Lives, got.Lives)
			}
			if tc.want.Level != got.Level {
				t.Errorf("want level %v, got %v", tc.want.Level, got.Level)
			}
		})
	}
}

func TestScoreExtractor_DonkeyKong(t *testing.T) {
	tests := []struct {
		scenario string
		filename string
		want     DonkeyKong
	}{
		{
			scenario: "DonkeyKong",
			filename: "./checkpoints_test_data/1615048176.ram",
			want: DonkeyKong{
				LivesPlayerOne: 2,
				LevelPlayerOne: 0,
				LivesPlayerTwo: 3,
				LevelPlayerTwo: 0,
				ScorePlayerOne: 5400,
				ScorePlayerTwo: 0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.scenario, func(t *testing.T) {
			got := ExtractDonkeyKong(Reader(tc.filename))

			if tc.want.LivesPlayerOne != got.LivesPlayerOne {
				t.Errorf("want lives player one %v, got %v", tc.want.LivesPlayerOne, got.LivesPlayerOne)
			}
			if tc.want.LevelPlayerOne != got.LevelPlayerOne {
				t.Errorf("want level player one %v, got %v", tc.want.LevelPlayerOne, got.LevelPlayerOne)
			}
			if tc.want.LivesPlayerTwo != got.LivesPlayerTwo {
				t.Errorf("want lives player two %v, got %v", tc.want.LivesPlayerTwo, got.LivesPlayerTwo)
			}
			if tc.want.LevelPlayerTwo != got.LevelPlayerTwo {
				t.Errorf("want level player two %v, got %v", tc.want.LevelPlayerTwo, got.LevelPlayerTwo)
			}
			if tc.want.ScorePlayerOne != got.ScorePlayerOne {
				t.Errorf("want score player one %v, got %v", tc.want.ScorePlayerOne, got.ScorePlayerOne)
			}
			if tc.want.ScorePlayerTwo != got.ScorePlayerTwo {
				t.Errorf("want score player two %v, got %v", tc.want.ScorePlayerTwo, got.ScorePlayerTwo)
			}
		})
	}
}
