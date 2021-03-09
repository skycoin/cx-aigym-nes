package score_extractor

import "fmt"

var (
	offset = 8 // address offset based on save file
)

var RomHashMap = map[string]string{
	"Tlfwh1Si/37HiCRWKftw+Z1OAD9m+GdCVmvKmcgQokQ=": "bomberman",
	"ETZFOAnlaBMYuBFa0jt/dpp62h+GSAkPwAbDnP6A2mE=": "donkeykong",
	"B0U4hfMsdlEl7cPyNQpECdzRWTbjbXMwNdxJWd38HiM=": "mario",
	"1enE6B/3MW1NixXyGYRxYG/bcouM2dqSsVO9xFEktWU=": "pacman",
	"rKsQ8tutTU711WsxzUzqIIu24uLF26XrM51xHsBONTU=": "tetris",
}

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
	ScorePlayerOne int64 `json:"score_player_one"`
	ScorePlayerTwo int64 `json:"score_player_two"`
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

func ExtractGameDetails(romHash string, ram []byte) {
	offset = 0 // Offset is zero when directly from the game, not a save file
	switch RomHashMap[romHash] {
	case "bomberman":
		fmt.Printf("%+v\n", ExtractBomberman(ram))
	case "donkeykong":
		fmt.Printf("%+v\n", ExtractDonkeyKong(ram))
	case "mario":
		fmt.Printf("%+v\n", ExtractSuperMarioBros(ram))
	case "pacman":
		fmt.Printf("%+v\n", ExtractPacman(ram))
	case "tetris":
		fmt.Printf("%+v\n", ExtractTetris(ram))
	default:
		fmt.Printf("rom hash cannot be found")
	}
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
	score := get7BCDFrom6Bytes([]byte{ram[0x0075+offset], ram[(0x0074)+offset], ram[0x0073+offset], ram[0x0072+offset], ram[0x0071+offset], ram[0x0070+offset]})
	return Pacman{
		Lives: int64(ram[0x0067+offset]),
		Level: int64(ram[0x0068+offset]),
		Score: score,
	}
}

// -----------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Bomberman:RAM_map
// -----------------------------------------------------------------------------
// Address | Information
// 0x0068  | Lives - Default value is 02
// 0x0058  | Level - Default value is 01 (Values: 01-32 hexadecimal)
// -----------------------------------------------------------------------------
func ExtractBomberman(ram []byte) Bomberman {
	return Bomberman{
		Lives: int64(ram[0x0068+offset]),
		Level: int64(ram[0x0058+offset]),
	}
}

// -------------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Donkey_Kong:RAM_map
// -------------------------------------------------------------------------------
// Address        | Information
// 0x0404         | Player 1 Marios remaining
// 0x0405         | Player 2 Marios remaining
// 0x0402         | Level number for player 1
// 0x0403         | Level number for player 2
// 0X0025-0X0027  | 6 digit 1P Score using BCD	1 nybble(4 bits) per digit
// 0X0029-0X002B  | 6 digit 2P Score using BCD	1 nybble(4 bits) per digit
// -------------------------------------------------------------------------------
func ExtractDonkeyKong(ram []byte) DonkeyKong {
	P1Score_byte_one := ram[0X0025+offset]
	P1Score_byte_two := ram[0X0026+offset]
	P1Score_byte_three := ram[0X0027+offset]

	P1Score := get6BCDFrom3Bytes([]byte{P1Score_byte_one, P1Score_byte_two, P1Score_byte_three})

	P2Score_byte_one := ram[0X0029+offset]
	P2Score_byte_two := ram[0X002A+offset]
	P2Score_byte_three := ram[0X002B+offset]

	P2Score := get6BCDFrom3Bytes([]byte{P2Score_byte_one, P2Score_byte_two, P2Score_byte_three})

	return DonkeyKong{
		LivesPlayerOne: int64(ram[0x0404+offset]),
		LivesPlayerTwo: int64(ram[0x0405+offset]),
		LevelPlayerOne: int64(ram[0x0402+offset]),
		LevelPlayerTwo: int64(ram[0x0403+offset]),
		ScorePlayerOne: P1Score,
		ScorePlayerTwo: P2Score,
	}
}

// -------------------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Super_Mario_Bros.:RAM_map
// -------------------------------------------------------------------------------------
// Address        | Information
// 0x075A         | Lives
// 0x075F		  | World
// 0x0760         | Level
// 0x07DD-0x07E2  | Mario score (1000000 100000 10000 1000 100 10) in BCD Format.
// 0x07D3/8       | Luigi score (1000000 100000 10000 1000 100 10) in BCD Format.
// -------------------------------------------------------------------------------------
func ExtractSuperMarioBros(ram []byte) SuperMarioBros {
	lvl := int64(ram[0x0760+offset])

	// If lvl is 0, it means its on the first level as observed in manual checking from the .ram file and scanner
	if lvl == 0 {
		lvl = 1
	}

	MScore := get7BCDFrom6Bytes([]byte{ram[0x07DD+offset], ram[(0x07DE)+offset], ram[0x07DF+offset], ram[0x07E0+offset], ram[0x07E1+offset], ram[0x07E2+offset]})
	LScore := get7BCDFrom6Bytes([]byte{ram[0x07D3+offset], ram[(0x07D4)+offset], ram[0x07D5+offset], ram[0x07D6+offset], ram[0x07D7+offset], ram[0x07D8+offset]})

	return SuperMarioBros{
		World:      int64(ram[0x075F+offset] + 1), // Plus 1 for World as observed in manual checking from the .ram file and scanner
		Level:      lvl,
		Lives:      int64(ram[0x075A+offset]),
		MarioScore: MScore,
		LuigiScore: LScore,
	}
}

// --------------------------------------------------------------------------------
// Link Ram Map Table: https://datacrystal.romhacking.net/wiki/Tetris_(NES):RAM_map
// --------------------------------------------------------------------------------
// Address        | Information
// 0x0053-0x0055  | Score (little endian bcd)
// 0x0044	      | current speed level
// --------------------------------------------------------------------------------
func ExtractTetris(ram []byte) Tetris {
	var score int64
	score_one := ram[0x0053+offset]
	score_two := ram[0x0054+offset]
	score_three := ram[0x0055+offset]

	score = get6BCDFrom3Bytes([]byte{score_three, score_two, score_one})

	return Tetris{
		Score: score,
		Level: int64(ram[0x0044+offset]),
	}
}

// get6BCDFrom3Bytes returns the 6 digit BCD value from three bytes where 1 nibble is equal to 1 digit
func get6BCDFrom3Bytes(bytes []byte) int64 {
	nibble_one := (bytes[0] & 0xF0) >> 4
	nibble_two := bytes[0] & 0xF
	nibble_three := (bytes[1] & 0xF0) >> 4
	nibble_four := bytes[1] & 0xF
	nibble_five := (bytes[2] & 0xF0) >> 4
	nibble_six := bytes[2] & 0xF

	return (int64(nibble_one) * 100000) + (int64(nibble_two) * 10000) + (int64(nibble_three) * 1000) + (int64(nibble_four) * 100) + (int64(nibble_five) * 10) + (int64(nibble_six))
}

// get7BCDFrom6Bytes returns the 7 digit BCD value from 6 bytes
func get7BCDFrom6Bytes(bytes []byte) int64 {
	return (int64(bytes[0]) * 1000000) + (int64(bytes[1]) * 100000) + (int64(bytes[2]) * 10000) + (int64(bytes[3]) * 1000) + (int64(bytes[4]) * 100) + (int64(bytes[5]) * 10)
}
