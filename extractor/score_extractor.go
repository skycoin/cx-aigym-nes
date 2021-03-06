package score_extractor

import (
	"encoding/hex"
	"errors"
	"strconv"
)

const (
	offset = 8 // address offset based on save file
)

type Pacman struct {
	Lives int64 `json:"lives"`
	Level int64 `json:"level"`
	Score int64 `json:"score"`
}

type Bomberman struct {
	Lives int64 `json:"lives"`
	Level int64 `json:"level"`
}

type DonkeyKong struct {
	LivesPlayerOne int64 `json:"lives_player_one"`
	LivesPlayerTwo int64 `json:"lives_player_two"`
	LevelPlayerOne int64 `json:"level_player_one"`
	LevelPlayerTwo int64 `json:"level_player_two"`
}

type SuperMarioBros struct {
	Lives      int64 `json:"lives"`
	World      int64 `json:"world"`
	Level      int64 `json:"level"`
	MarioScore int64 `json:"mario_score"`
	LuigiScore int64 `json:"luigi_score"`
}

type Tetris struct {
	Score int64 `json:"score"`
	Level int64 `json:"level"`
}

// ----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Pac-Man:RAM_map
// ----------------------------------------------------------------------------
// Address | Information
// 0x0067  | Lives
// 0x0068  | Level
// 0x0070  | Current Score 00000*0, can be used to set the current score.
// 0x0071  | Current Score 0000*00, can be used to set the current score.
// 0x0072  | Current Score 000*000, can be used to set the current score.
// 0x0073  | Current Score 00*0000, can be used to set the current score.
// 0x0074  | Current Score 0*00000, can be used to set the current score.
// 0x0075  | Current Score *000000, can be used to set the current score.
// ----------------------------------------------------------------------------
func ExtractPacman(ram []byte) Pacman {
	score := (int64(ram[0x0075+offset]) * 1000000) + (int64(ram[(0x0074)+offset]) * 100000) + (int64(ram[0x0073+offset]) * 10000) + (int64(ram[0x0072+offset]) * 1000) + (int64(ram[0x0071+offset]) * 100) + (int64(ram[0x0070+offset]) * 10)
	return Pacman{
		Lives: int64(ram[0x0067+offset]),
		Level: int64(ram[0x0068+offset]),
		Score: score,
	}
}

// ----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Bomberman:RAM_map
// ----------------------------------------------------------------------------
// Address | Information
// 0x0068  | Lives - Default value is 02
// 0x0058  | Level - Default value is 01 (Values: 01-32 hexadecimal)
// ----------------------------------------------------------------------------
func ExtractBomberman(ram []byte) Bomberman {
	return Bomberman{
		Lives: int64(ram[0x0068+offset]),
		Level: int64(ram[0x0058+offset]),
	}
}

// ----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Donkey_Kong:RAM_map
// ----------------------------------------------------------------------------
// Address | Information
// 0x0404  | Player 1 Marios remaining
// 0x0405  | Player 2 Marios remaining
// 0x0402  | Level number for player 1
// 0x0403  | Level number for player 2
// ----------------------------------------------------------------------------
func ExtractDonkeyKong(ram []byte) DonkeyKong {
	return DonkeyKong{
		LivesPlayerOne: int64(ram[0x0404+offset]),
		LivesPlayerTwo: int64(ram[0x0405+offset]),
		LevelPlayerOne: int64(ram[0x0402+offset]),
		LevelPlayerTwo: int64(ram[0x0403+offset]),
	}
}

// ----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Super_Mario_Bros.:RAM_map
// ----------------------------------------------------------------------------
// Address        | Information
// 0x075A         | Lives
// 0x075F		  | World
// 0x0760         | Level
// 0x07DD-0x07E2  | Mario score (1000000 100000 10000 1000 100 10) in BCD Format.
// 0x07D3/8       | Luigi score (1000000 100000 10000 1000 100 10) in BCD Format.
// ----------------------------------------------------------------------------
func ExtractSuperMarioBros(ram []byte) SuperMarioBros {
	lvl := int64(ram[0x0760+offset])
	if lvl == 0 {
		lvl = 1
	}

	MScore := int64(ram[0x07DD+offset])<<40 | int64(ram[(0x07DE)+offset])<<32 | int64(ram[0x07DF+offset])<<24 | int64(ram[0x07E0+offset])<<16 | int64(ram[0x07E1+offset])<<8 | int64(ram[0x07E2+offset])
	LScore := int64(ram[0x07D3+offset])<<40 | int64(ram[(0x07D4)+offset])<<32 | int64(ram[0x07D5+offset])<<24 | int64(ram[0x07D6+offset])<<16 | int64(ram[0x07D7+offset])<<8 | int64(ram[0x07D8+offset])

	return SuperMarioBros{
		World:      int64(ram[0x075F+offset] + 1),
		Level:      lvl,
		Lives:      int64(ram[0x075A+offset]),
		MarioScore: MScore,
		LuigiScore: LScore,
	}
}

// ----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Tetris_(NES):RAM_map
// ----------------------------------------------------------------------------
// Address        | Information
// 0x0053-0x0055  | Score (little endian bcd)
// 0x0044	      | current speed level
// ----------------------------------------------------------------------------
func ExtractTetris(ram []byte) Tetris {
	var score int64
	score_one := ram[0x0053+offset]
	score_two := ram[0x0054+offset]
	score_three := ram[0x0055+offset]

	score, err := decodeBcd([]byte{score_three, score_two, score_one})
	if err != nil {
		panic(err)
	}

	return Tetris{
		Score: score,
		Level: int64(ram[0x0044+offset]),
	}
}

func decodeBcd(bcd []byte) (int64, error) {
	s := hex.EncodeToString(bcd)
	if s[len(s)-1] == 'f' {
		s = s[:len(s)-1]
	}
	result, err := strconv.ParseUint(s, 10, 64)
	if err == nil {
		return int64(result), nil
	}
	switch err.(*strconv.NumError).Err {
	case strconv.ErrRange:
		return 0, errors.New("Overflow occurred in BCD decoding")
	case strconv.ErrSyntax:
		return 0, errors.New("Bad digit in BCD decoding")
	default:
		panic("unexpected error from strconv.ParseUint")
	}
}
