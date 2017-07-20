package main

import (
	blt "bearlibterminal"
	"camera"
	"entity"
	"fov"
	"gamemap"
	"strconv"
)

const (
	WindowSizeX = 100
	WindowSizeY = 35
	ViewAreaX = 100
	ViewAreaY = 30
	MapWidth = 150
	MapHeight = 150
	Title = "BearRogue"
	Font = "fonts/UbuntuMono.ttf"
	FontSize = 24
)

var (
	player *entity.GameEntity
	entities []*entity.GameEntity
	gameMap *gamemap.Map
	gameCamera *camera.GameCamera
	fieldOfView *fov.FieldOfVision
	messageLog []string
)

func init() {
	blt.Open()

	// BearLibTerminal uses configuration strings to set itself up, so we need to build these strings here
	// First set up the string for window properties (size and title)
	size := "size=" + strconv.Itoa(WindowSizeX) + "x" + strconv.Itoa(WindowSizeY)
	title := "title='" + Title + "'"
	window := "window: " + size + "," + title

	// Next, setup the font config string
	fontSize := "size=" + strconv.Itoa(FontSize)
	font := "font: " + Font + ", " + fontSize

	// Now, put it all together
	blt.Set(window + "; " + font)
	blt.Clear()

	// Create a player Entity and an NPC entity, and add them to our slice of Entities
	player = &entity.GameEntity{X: 1, Y: 1, Layer: 1, Char: "@", Color: "white"}
	npc := &entity.GameEntity{X: 10, Y: 10, Layer: 0, Char: "N", Color: "red"}
	entities = append(entities, player, npc)

	// Create a GameMap, and initialize it (and set the player position within it, for now)
	gameMap = &gamemap.Map{Width: MapWidth, Height: MapHeight}
	gameMap.InitializeMap()

	playerX, playerY := gameMap.GenerateCavern()
	player.X = playerX
	player.Y = playerY

	// Initialize a camera object
	gameCamera = &camera.GameCamera{X: 1, Y:1, Width: ViewAreaX, Height: ViewAreaY}

	// Initialize a FoV object
	fieldOfView = &fov.FieldOfVision{}
	fieldOfView.Initialize()
	fieldOfView.SetTorchRadius(500)

	messageLog = make([]string, 0)
	sendMessage("You find yourself in the caverns of eternal sadness...you start to feel a little more sad.")
}
	
func main() {
	// Main game loop

	renderAll()
	printMessages()

	for {
		blt.Refresh()

		key := blt.Read()

		// Clear each Entity off the screen
		for _, e := range entities {
			mapX, mapY := gameCamera.ToCameraCoordinates(e.X, e.Y)
			e.Clear(mapX, mapY)
		}

		if key != blt.TK_CLOSE {
			handleInput(key, player)
		} else {
			break
		}

		renderAll()
		printMessages()
	}

	blt.Close()
}

func handleInput(key int, player *entity.GameEntity) {
	// Handle basic character movement in the four main directions

	var (
		dx, dy int
	)

	switch key {
	case blt.TK_RIGHT:
		dx, dy = 1, 0
	case blt.TK_LEFT:
		dx, dy = -1, 0
	case blt.TK_UP:
		dx, dy = 0, -1
	case blt.TK_DOWN:
		dx, dy = 0, 1
	}

	// Check to ensure that the tile the player is trying to move in to is a valid move (not blocked)
	if !gameMap.IsBlocked(player.X + dx, player.Y + dy) {
		player.Move(dx, dy)
	}
}

func renderEntities() {
	// Draw every Entity present in the game. This gets called on each iteration of the game loop.
	for _, e := range entities {
		if e != player {
			cameraX, cameraY := gameCamera.ToCameraCoordinates(e.X, e.Y)
			if gameMap.Tiles[e.X][e.Y].Visible {
				e.Draw(cameraX, cameraY)
			}
		}
	}

	cameraX, cameraY := gameCamera.ToCameraCoordinates(player.X, player.Y)
	player.Draw(cameraX, cameraY)
}

func renderMap() {
	// Render the game map. If a tile is blocked and blocks sight, draw a '#', if it is not blocked, and does not block
	// sight, draw a '.'

	// First, set the entire map to not visible. We'll decide what is visible based on the torch radius.
	// In the process, clear every Tile on the map as well
	for x := 0; x < gameMap.Width; x++ {
		for y := 0; y < gameMap.Height; y++ {
			gameMap.Tiles[x][y].Visible = false
			blt.Print(x, y, " ")
		}
	}

	// Next figure out what is visible to the player, and what is not.
	fieldOfView.RayCast(player.X, player.Y, gameMap)

	// Now draw each tile that should appear on the screen, if its visible, or explored
	for x := 0; x < gameCamera.Width; x++ {
		for y := 0; y < gameCamera.Height; y++ {
			mapX, mapY := gameCamera.X + x, gameCamera.Y + y

			if gameMap.Tiles[mapX][mapY].Visible {
				if gameMap.Tiles[mapX][mapY].IsWall() {
					blt.Color(blt.ColorFromName("white"))
					blt.Print(x, y, "#")
				} else {
					blt.Color(blt.ColorFromName("white"))
					blt.Print(x, y, ".")
				}
			} else if gameMap.Tiles[mapX][mapY].Explored {
				if gameMap.Tiles[mapX][mapY].IsWall() {
					blt.Color(blt.ColorFromName("gray"))
					blt.Print(x, y, "#")
				} else {
					blt.Color(blt.ColorFromName("gray"))
					blt.Print(x, y, ".")
				}
			}
		}
	}
}

func renderAll() {
	// Convenience function to render all entities, followed by rendering the game map

	// Before anything is rendered, update the camera position, so it is centered (if possible) on the player
	// Only things within the cameras viewport will be drawn to the screen
	gameCamera.MoveCamera(player.X, player.Y, MapWidth, MapHeight)

	renderMap()
	renderEntities()
}

func sendMessage(message string) {
	// Prepend the message onto the messageLog slice
	if len(messageLog) >= 99 {
		// Throw away any messages that exceed our total queue size
		messageLog = messageLog[:len(messageLog)-1]
	}
	messageLog = append([]string{message}, messageLog...)
}

func clearMessages() {
	// Clear the message area, so our messages do not overlap
	blt.ClearArea(0, ViewAreaY, WindowSizeX, WindowSizeY - ViewAreaY)
}

func printMessages() {
	// Print the latest five messages from the messageLog. These will be printed in reverse order (newest at the top),
	// to make it appear they are scrolling down the screen
	clearMessages()

	toShow := 0

	if len(messageLog) <= 5 {
		// Just loop through the messageLog, printing them in reverse order
		toShow = len(messageLog)
	} else {
		// If we have more than 5 messages stored, just show the five most recent
		toShow = 5
	}

	blt.Color(blt.ColorFromName("white"))
	blt.Layer(1)
	for i := toShow; i > 0; i-- {
		blt.Print(1, (ViewAreaY - 1) + i, messageLog[i - 1])
	}
}

