# cx-aigym-nes

Testing Gym for Cartesian Genetic Algorithms.

What is the size in bytes of the smallest CX assembly program that can beat mario world, without looking at the screen or taking any input from the game?

#Usage

```
NAME:
   cx-aigym-nes - A new cli application

USAGE:
   cx [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --disable-audio               disable audio (default: false)
   --disable-video               disable video (default: false)
   --savedirectory value         Path to store the state of games
   --loadrom value, --lr value   load .rom file/s
   --loadjson value, --lj value  load .json file/s
   --help, -h                    show help (default: false)
   --version, -v                 print the version (default: false)
```

# Examples
```
./cx-aigym-nes --disable-audio --disable-video 
--loadrom ../roms/Super_mario_brothers.nes --savedirectory ../checkpoints
```
1. --disable-audio: Disable audio, default is false if missing
2. --disable-video: Disable video, default is false if missing
3. --savedirectory: Where to store the state of games, default is ../checkpoints if missing
4. --loadrom: Path of rom file