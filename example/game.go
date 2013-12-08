package main

import (
	"fmt"
	"github.com/hajimehoshi/go-ebiten/graphics"
	"github.com/hajimehoshi/go-ebiten/graphics/matrix"
	"github.com/hajimehoshi/go-ebiten/ui"
	"image/color"
	"math"
)

var TexturePaths = map[string]string{
	"ebiten": "images/ebiten.png",
	"text":   "images/text.png",
}

type Size struct {
	Width  int
	Height int
}

var RenderTargetSizes = map[string]Size{
	"whole": Size{256, 254},
}

type drawInfo struct {
	textures      map[string]graphics.TextureId
	renderTargets map[string]graphics.RenderTargetId
	inputStr      string
	textureX      int
	textureY      int
	textureAngle  float64
	textureGeo    matrix.Geometry
}

type Game struct {
	inputX     int
	inputY     int
	inputPrevX int
	inputPrevY int
	counter    int
	drawInfo
}

func NewGame() *Game {
	return &Game{
		inputX:     -1,
		inputY:     -1,
		inputPrevX: -1,
		inputPrevY: -1,
		counter:    0,
		drawInfo: drawInfo{
			textures:      map[string]graphics.TextureId{},
			renderTargets: map[string]graphics.RenderTargetId{},
			textureX:      0,
			textureY:      0,
			textureAngle:  0,
			textureGeo:    matrix.IdentityGeometry(),
		},
	}
}

func (game *Game) OnTextureCreated(e graphics.TextureCreatedEvent) {
	if e.Error != nil {
		panic(e.Error)
	}
	game.textures[e.Tag.(string)] = e.Id
}

func (game *Game) OnRenderTargetCreated(e graphics.RenderTargetCreatedEvent) {
	if e.Error != nil {
		panic(e.Error)
	}
	game.renderTargets[e.Tag.(string)] = e.Id
}

func (game *Game) OnInputStateUpdated(e ui.InputStateUpdatedEvent) {
	game.inputX, game.inputY = e.X, e.Y
}

func (game *Game) isInitialized() bool {
	if len(game.drawInfo.textures) < len(TexturePaths) {
		return false
	}
	if len(game.drawInfo.renderTargets) < len(RenderTargetSizes) {
		return false
	}
	return true
}

func (game *Game) Update() {
	if !game.isInitialized() {
		return
	}

	const textureWidth = 57
	const textureHeight = 26

	game.counter++
	game.drawInfo.inputStr = fmt.Sprintf(`Input State:
  X: %d
  Y: %d`, game.inputX, game.inputY)

	if game.inputPrevX != -1 && game.inputPrevY != -1 &&
		game.inputX != -1 && game.inputY != -1 {
		dx, dy := game.inputX-game.inputPrevX, game.inputY-game.inputPrevY
		game.textureX += dx
		game.textureY += dy

	}
	game.drawInfo.textureAngle = 2 * math.Pi * float64(game.counter) / 600

	geo := matrix.IdentityGeometry()
	geo.Translate(-textureWidth/2, -textureHeight/2)
	geo.Rotate(game.drawInfo.textureAngle)
	geo.Translate(textureWidth/2, textureHeight/2)
	geo.Translate(float64(game.textureX), float64(game.textureY))

	game.drawInfo.textureGeo = geo

	// Update for the next frame.
	game.inputPrevX, game.inputPrevY = game.inputX, game.inputY
}

func (game *Game) Draw(g graphics.Canvas) {
	if !game.isInitialized() {
		return
	}

	whole := game.drawInfo.renderTargets["whole"]
	g.SetOffscreen(whole)
	g.Fill(0x40, 0x60, 0xb0)
	game.drawTexture(g, game.drawInfo.textureGeo, matrix.IdentityColor())
	game.drawText(g, game.drawInfo.inputStr, 6, 6, &color.RGBA{0x0, 0x0, 0x0, 0x80})
	game.drawText(g, game.drawInfo.inputStr, 5, 5, color.White)

	g.ResetOffscreen()
	g.DrawRenderTarget(whole, matrix.IdentityGeometry(), matrix.IdentityColor())
	wholeGeo := matrix.IdentityGeometry()
	wholeGeo.Scale(0.5, 0.5)
	wholeGeo.Translate(256/2, 240/2)
	wholeColor := matrix.IdentityColor()
	wholeColor.Scale(&color.RGBA{0x80, 0x80, 0x80, 0x80})
	g.DrawRenderTarget(whole, wholeGeo, wholeColor)
}

func (game *Game) drawText(g graphics.Canvas, text string, x, y int, clr color.Color) {
	const letterWidth = 6
	const letterHeight = 16

	parts := []graphics.TexturePart{}
	textX := 0
	textY := 0
	for _, c := range text {
		if c == '\n' {
			textX = 0
			textY += letterHeight
			continue
		}
		code := int(c)
		x := (code % 32) * letterWidth
		y := (code / 32) * letterHeight
		source := graphics.Rect{x, y, letterWidth, letterHeight}
		parts = append(parts, graphics.TexturePart{
			LocationX: textX,
			LocationY: textY,
			Source:    source,
		})
		textX += letterWidth
	}

	geometryMatrix := matrix.IdentityGeometry()
	geometryMatrix.Translate(float64(x), float64(y))
	colorMatrix := matrix.IdentityColor()
	colorMatrix.Scale(clr)
	g.DrawTextureParts(game.drawInfo.textures["text"], parts,
		geometryMatrix, colorMatrix)
}

func (game *Game) drawTexture(g graphics.Canvas, geo matrix.Geometry, color matrix.Color) {
	g.DrawTexture(game.drawInfo.textures["ebiten"], geo, color)
}
