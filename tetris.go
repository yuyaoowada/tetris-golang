package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/mattn/go-tty"
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
				dispStr += "□ "
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
	dispStr += "q:終了\n"

	// まとめて表示
	fmt.Printf("%s", dispStr)
}

// すでに積んであるブロックと現在のブロックの当たり判定
func intersectBlock(directionX, directionY int) int {
	for innerY := 0; innerY < 4; innerY++ {
		for innerX := 0; innerX < 4; innerX++ {
			x := nowBlockX + innerX + directionX
			y := nowBlockY + innerY + directionY
			if y == dispH || x <= -1 || x >= dispW {
				// 境界との当たり判定
				if nowBlock[innerY][innerX] != 0 {
					return 1
				}
			}
			if (0 <= x && x < dispW) && (0 <= y && y < dispH) {
				// フィールドとの当たり判定
				if nowBlock[innerY][innerX] != 0 && field[y][x] != 0 {
					return 1
				}
			}
		}
	}
	return 0
}

// ブロックがはみ出しているか
func intersectBlockUpper() int {
	for innerY := 0; innerY < 4; innerY++ {
		for innerX := 0; innerX < 4; innerX++ {
			y := nowBlockY + innerY
			if y < 0 {
				// 境界との当たり判定
				if nowBlock[innerY][innerX] != 0 {
					return 1
				}
			}
		}
	}
	return 0
}

// 降っているブロックからフィールドへ
func nowBlockToField() {
	for innerY := 0; innerY < 4; innerY++ {
		for innerX := 0; innerX < 4; innerX++ {
			x := nowBlockX + innerX
			y := nowBlockY + innerY
			if (0 <= x && x < dispW) && (0 <= y && y < dispH) {
				if nowBlock[innerY][innerX] != 0 {
					field[y][x] = nowBlock[innerY][innerX]
				}
			}
		}
	}
}

// ブロックが消える処理
func eraseBlock() {
	// 空白がある行をテンポラリにコピーする
	var tmpField [dispH][dispW]int

	tmpY := dispH - 1
	for y := dispH - 1; y >= 0; y-- {
		blank := 0 // 空白ブロックがあるか

		// 一行処理
		for x := 0; x < dispW; x++ {
			if field[y][x] == 0 {
				// 空白を見つけた
				blank = 1
				break
			}
		}

		// テンポラリにコピー
		if blank != 0 {
			for x := 0; x < dispW; x++ {
				tmpField[tmpY][x] = field[y][x]
			}
			tmpY--
		}
	}

	// テンポラリからフィールドに移す
	for y := 0; y < dispH; y++ {
		for x := 0; x < dispW; x++ {
			field[y][x] = tmpField[y][x]
		}
	}
}

// ブロックの回転
func rotateBlockL() {
	// テンポラリ
	var tmpBlock [4][4]int
	var srcX, srcY, dstX, dstY int

	// 今降っているブロックテンポラリ
	srcY = 0
	dstX = 0
	for srcY < 4 {
		srcX = 0
		dstY = 3
		for srcX < 4 {
			tmpBlock[dstY][dstX] = nowBlock[srcY][srcX]
			srcX++
			dstY--
		}
		srcY++
		dstX++
	}

	// テンポラリ→今降っているブロック
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			nowBlock[y][x] = tmpBlock[y][x]
		}
	}
}

// 初期化
func initGame() {
	// 乱数初期化
	rand.Seed(time.Now().UnixNano())

	// フィールド初期化
	for y := 0; y < dispH; y++ {
		for x := 0; x < dispW; x++ {
			field[y][x] = 0
		}
	}

	// 降っているテトリス初期化
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			nowBlock[y][x] = 0
		}
	}
	nowBlockX = 0
	nowBlockY = 0
	speed = initSpeed
	speedCnt = speed
	waitBlockCnt = waitBlockFPS

	// ゲーム状態初期化
	gameStatus = 1
}

// ゲーム処理
func execGame(str string) {
	if str == "q" {
		// 終了
		os.Exit(0)
	}

	switch gameStatus {
	case 1: // 降るテトリス選定
		{
			// ブロック決定
			bn := rand.Intn(len(block))
			// ブロックコピー
			for y := 0; y < 4; y++ {
				for x := 0; x < 4; x++ {
					nowBlock[y][x] = block[bn][y][x]
				}
			}
			// 位置の初期化
			nowBlockX = dispW/2 - 2
			nowBlockY = -4
			// 降っているテトリスが止まるカウント
			waitBlockCnt = waitBlockFPS
		}
		gameStatus++
	case 2: // 降っている状態
		{
			if str == " " { // 回転
				rotateBlockL()
				waitBlockCnt = waitBlockFPS
			}
			if str == "a" { // 左
				if intersectBlock(-1, 0) == 0 {
					nowBlockX--
				}
			}
			if str == "d" { // 右
				if intersectBlock(1, 0) == 0 {
					nowBlockX++
				}
			}
			if str == "s" { // 下
				if intersectBlock(0, 1) == 0 {
					nowBlockY++
				}
			}
			if str == "w" { // 上
				for {
					if intersectBlock(0, 1) != 0 {
						break
					}
					nowBlockY++
					waitBlockCnt = 1
				}
			}

			if speedCnt <= 0 {
				if intersectBlock(0, 1) == 0 {
					nowBlockY++
				}
				speedCnt = speed
			}
			speedCnt--

			// 回転して重なるなら戻す
			for i := 0; i < 4; i++ {
				if intersectBlock(0, 0) != 0 {
					if intersectBlock(-1, 0) == 0 {
						nowBlockX--
						break
					}
					if intersectBlock(1, 0) == 0 {
						nowBlockX++
						break
					}
					nowBlockY--
				}
			}
			// ブロックが着地しているか
			if intersectBlock(0, 1) != 0 {
				waitBlockCnt--
				if waitBlockCnt <= 0 {
					// 位置指定
					nowBlockToField()
					eraseBlock()
					// ゲームオーバー判定
					if intersectBlockUpper() != 0 {
						// ゲームオーバー
						gameStatus = 3
					} else {
						// 継続
						gameStatus = 1
					}
				}
			}
		}
	case 3: //ゲームオーバー
		fmt.Printf("GameOver.\n")
		os.Exit(0)
	}
}

// メインループ
func mainLoop() {
	// ループミリ秒
	var loopms time.Duration = 33
	// タイマ作成
	timer := time.NewTimer(loopms * time.Millisecond)
	// 初期化
	initGame()

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	// メインループ
	for {
		select {
		case <-timer.C:
			timer = time.NewTimer(loopms * time.Millisecond)
			// 描画
			drawDisp()
			// ゲーム処理
			for {
				r, err := tty.ReadRune()
				if err != nil {
					log.Fatal(err)
				}
				if string(r) != "" {
					// 入力によって操作
					execGame(string(r))
				}
				break
			}
		}
	}
}

func main() {
	mainLoop()
}
