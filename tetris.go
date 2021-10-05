package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// 定数定義
const dispW = 10       // 幅が10ブロック
const dispH = 20       // 高さが20ブロック
const waitBlockFPS = 5 // テトリス停止時間(FPS)
const showX = 2        // 表示位置x
const showY = 2        // 表示位置y

// ブロック定義
var block = [][][]int{
	{
		// 棒
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 1, 0},
	},
	// 四角
	{
		{0, 0, 0, 0},
		{0, 1, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// T
	{
		{0, 0, 0, 0},
		{0, 1, 1, 1},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	},
	// N
	{
		{0, 0, 1, 0},
		{0, 1, 1, 0},
		{0, 1, 0, 0},
		{0, 0, 0, 0},
	},
	// N
	{
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 0},
	},
	// L
	{
		{0, 1, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
	// L
	{
		{0, 0, 1, 0},
		{0, 0, 1, 0},
		{0, 1, 1, 0},
		{0, 0, 0, 0},
	},
}

// ブロック領域
var field [dispH][dispW]int

// 降っているテトリスブロック
var nowBlock [4][4]int

// 降っているテトリスの位置
var nowBlockX = 0
var nowBlockY = 0

// 速度
var initSpeed = 5
var speed = initSpeed
var speedCnt = speed

// ゲームの状態
// 0:なし
// 1:降るテトリス選定
// 2:降っている
// 3:ゲームオーバー
var gameStatus = 1

// 降っているテトリスのウエイトカウンタ
var waitBlockCnt = waitBlockFPS

// 画面クリア
func clearDisplay() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 画面描画
func drawDisp() {
	// 画面クリア
	clearDisplay()
	// 表示文字列
	dispStr := ""

	// Todo

	// 描画位置調整
	for y := 0; y < showY; y++ {
		dispStr += "\n"
	}

	for y := 0; y < dispH+1; y++ {

		// 描画位置調整
		for x := 0; x < showX; x++ {
			dispStr += "  "
		}

		for x := -1; x < dispW+1; x++ {
			if y == dispH || x == -1 || x == dispW {
				// 境界線
				dispStr += "■ "
				continue
			}

			disp := 0
			if field[y][x] != 0 {
				// フィールド
				disp = 1
			}

			innerX := x - nowBlockX
			innerY := y - nowBlockY
			if 0 <= innerX && innerX < 4 &&
				0 <= innerY && innerY < 4 &&
				nowBlock[innerY][innerX] != 0 {
				disp = 1
			}

			if disp != 0 {
				dispStr += "□"
			} else {
				dispStr += "  "
			}
		}

		dispStr += "\n"

	}

	// 操作方法説明
	dispStr += "\n"
	dispStr += "w:すぐに落とす,"
	dispStr += "a:左,"
	dispStr += "d:右,"
	dispStr += "s:下,"
	dispStr += "スペース:回転\n,"
	dispStr += "q:終了\n,"

	// まとめて表示
	fmt.Printf("%s", dispStr)
}

func mainLoop() {
	// ループミリ秒
	var loopms time.Duration = 33
	// タイマ作成
	timer := time.NewTimer(loopms * time.Millisecond)

	// メインループ
	for {
		select {
		case <-timer.C:
			timer = time.NewTimer(loopms * time.Millisecond)
		}
		// 描画
		drawDisp()
		break
	}
}

func main() {
	mainLoop()
}
