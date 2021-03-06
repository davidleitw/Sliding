# <img src="https://img.icons8.com/color/96/000000/chessboard.png"/> Sliding

## <img src="https://img.icons8.com/color/48/000000/queen.png"/> 簡介

`Sliding` 是一個基於 `sliding window algorithm` 的輕量級 `library`。開發它的靈感來自於一些微服務限流專案(Ex: [guava](https://github.com/google/guava),[Sentinal](https://github.com/alibaba/Sentinel)) 中關於統計 `QPS` 的作法，透過將一些簡單的操作封裝，可以有效的省去開發的成本。

對於一些需要統計單位時間內發生事件的場景，都可以利用 `Sliding` 去簡單的統計資料，`Sliding` 可以自訂 `Upload function`，用來處理 `window` 滑動後，原本位於 `window` 的資料要怎麼處理，寫入資料庫，或者根據設定的 `threshold` 去觸發某些行為等等。

#### Installation

```go
go get github.com/davidleitw/Sliding/pkg/slidingwindow
```

#### include

```go
import "github.com/davidleitw/Sliding/pkg/slidingwindow"
```
## <img src="https://img.icons8.com/color/48/000000/king.png"/> 概念


在建立 `Sliding window` 的時候可以指定 `windowSize` 與 `windowLength`。

- `windowSize`: 單一窗口時間，用 `ms` 為單位
- `windowLength`: 總共有幾個 `window`

創立一個 `windowSize=125ms`, `windowLength=8` 的 `Sliding Window` 可以參考以下範例，這樣一輪總共有 `1000ms = 1s`。

```go
slw := slidingwindow.NewSlidingWindows(125, 8, nil)
```

![](https://i.imgur.com/7GQW2su.png)

可以簡單用上圖的例子來說明 `Sliding window` 的運作原理

假設今天要紀錄的資料發生在`487ms`，我們可以透過 `(487 / windowSize) % wnidowLength` 來求出 `index=3`，我們也可以透過 `487 - (487 % windowLength)` 來算出這輪的 `start = 375`

那 `Sliding window` 創立之後怎麼滑動呢? 可以用下面的圖片說明，一開始的步驟都像上面，但是不一樣的地方在於算出來的 `start=48763` 會跟 `windows[index].start` 先比大小，如果發現新出來的 `start` 大於原本 `window` 的 `start`，就會更新相關的資料，在 `Sliding` 內會分成兩個部份執行

- 更新 `window` 的 `start`
- 執行使用者設定的 `upload function`，處理 `window` 更新前存在內部的統計資料，會在後續的範例中看到詳細的用法

![](https://i.imgur.com/EQSFTXK.png)

## <img src="https://img.icons8.com/color/48/000000/knight.png"/> 範例

動工中

## <img src="https://img.icons8.com/color/48/000000/pawn.png"/> 參考文章

動工中

## TODO

- 思考兩個使用範例
    - 整合 InfluxDB, Grafana 呈現
- 思考 window 的更新問題，尤其是超過兩輪沒有更新時如何補上空缺的資料