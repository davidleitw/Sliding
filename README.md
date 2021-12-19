# <img src="https://img.icons8.com/color/96/000000/chessboard.png"/> Sliding

## <img src="https://img.icons8.com/color/48/000000/queen.png"/> 簡介

`Sliding` 是一個基於 `sliding window algorithm` 的輕量級 `library`。開發它的靈感來自於一些微服務限流專案(Ex: [guava](https://github.com/google/guava),[Sentinal](https://github.com/alibaba/Sentinel)) 中關於統計 `QPS` 的作法，透過將一些簡單的操作封裝，可以有效的省去開發的成本。

對於一些需要統計單位時間內發生事件的場景，都可以利用 `Sliding` 去簡單的統計資料，`Sliding` 可以自訂 `Upload function`，用來處理 `window` 滑動後，原本位於 `window` 的資料要怎麼處理，寫入資料庫，或者根據設定的 `threshold` 去觸發某些行為等等。

## <img src="https://img.icons8.com/color/48/000000/king.png"/> 概念

#### Installation

```go
go get github.com/davidleitw/Sliding/pkg/slidingwindow
```

#### include

```go
import "github.com/davidleitw/Sliding/pkg/slidingwindow"
```

在建立 `Sliding window` 的時候可以指定 `windowSize` 與 `windowLength`。

- `windowSize`: 單一窗口時間，用 `ms` 為單位
- `windowLength`: 總共有幾個 `window`

創立一個 `windowSize=100ms`, `windowLength=10` 的 `Sliding Window` 可以參考以下範例，這樣一輪總共有 `1000ms = 1s`。

```go
slw := slidingwindow.NewSlidingWindows(100, 10, nil)
```

## <img src="https://img.icons8.com/color/48/000000/knight.png"/> 範例

動工中

## <img src="https://img.icons8.com/color/48/000000/pawn.png"/> 參考文章

動工中