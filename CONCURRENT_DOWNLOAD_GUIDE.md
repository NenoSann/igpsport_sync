# 并行下载使用指南

## 概述

`DownloadAllActivitiesWithConcurrency` 方法支持使用多个 goroutine 并发下载文件，适合在需要高效传输的场景中使用。

## 方法签名

```go
func (s *IgpsportSync) DownloadAllActivitiesWithConcurrency(
    ext Extension, 
    maxConcurrency int, 
    callback DownloadCallback,
) error
```

### 参数说明

- **ext**: 文件格式（FIT、GPX、TCX）
- **maxConcurrency**: 最大并发数（必须 > 0）
  - 建议值：5-10（避免对服务器造成过大压力）
  - 过小会影响下载速度
  - 过大可能被服务器限流
- **callback**: 处理每个下载结果的回调函数

## 工作原理

```
主线程 (获取分页)
    ↓
工作队列 (ActivityRow channel)
    ↓
Worker 1 (下载) ──┐
Worker 2 (下载) ──┤─→ 回调函数 → 上传到 OSS / 保存文件
Worker 3 (下载) ──┘
...
Worker N (下载)
```

主线程负责获取分页，Worker 线程负责实际的下载操作。所有 Worker 并发运行，提高下载效率。

## 使用示例

### 示例1：基础用法 - 并行下载并上传到 OSS

```go
const maxConcurrency = 5

err := igpsport.DownloadAllActivitiesWithConcurrency(
    igpsportsync.FIT,
    maxConcurrency,
    func(activity *igpsportsync.DownloadedActivity) bool {
        if activity.Error != nil {
            fmt.Printf("✗ Download failed for activity %d: %v\n", 
                activity.RideID, activity.Error)
            // 继续处理下一个
            return true
        }
        
        // 上传到 OSS
        ossKey := fmt.Sprintf("activities/%d.fit", activity.RideID)
        if err := uploadToOSS(ossKey, activity.Data); err != nil {
            fmt.Printf("✗ Upload failed for %s: %v\n", ossKey, err)
        } else {
            fmt.Printf("✓ Activity %d uploaded successfully\n", activity.RideID)
        }
        
        return true // 继续处理下一个
    },
)

if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

### 示例2：带进度统计

```go
type DownloadStats struct {
    Total       int
    Successful  int
    Failed      int
    TotalBytes  int64
    mu          sync.Mutex
}

stats := &DownloadStats{}

err := igpsport.DownloadAllActivitiesWithConcurrency(
    igpsportsync.GPX,
    8,
    func(activity *igpsportsync.DownloadedActivity) bool {
        stats.mu.Lock()
        defer stats.mu.Unlock()
        
        stats.Total++
        
        if activity.Error != nil {
            stats.Failed++
            fmt.Printf("[Failed] Activity %d: %v\n", 
                activity.RideID, activity.Error)
        } else {
            stats.Successful++
            stats.TotalBytes += int64(len(activity.Data))
            fmt.Printf("[%d/%d] Activity %d ✓ (%d KB)\n",
                stats.Successful, stats.Total,
                activity.RideID,
                len(activity.Data)/1024)
        }
        
        return true
    },
)

if err == nil {
    fmt.Printf("\n=== Download Complete ===\n")
    fmt.Printf("Total: %d\n", stats.Total)
    fmt.Printf("Successful: %d\n", stats.Successful)
    fmt.Printf("Failed: %d\n", stats.Failed)
    fmt.Printf("Total size: %.2f MB\n", float64(stats.TotalBytes)/1024/1024)
}
```

### 示例3：速度对比 - 串行 vs 并行

```go
import "time"

// 串行下载
startSerial := time.Now()
err1 := igpsport.DownloadAllActivities(igpsportsync.FIT, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error == nil {
        saveFile(activity)
    }
    return true
})
serialTime := time.Since(startSerial)

// 并行下载（5个并发）
startConcurrent := time.Now()
err2 := igpsport.DownloadAllActivitiesWithConcurrency(igpsportsync.FIT, 5, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error == nil {
        saveFile(activity)
    }
    return true
})
concurrentTime := time.Since(startConcurrent)

fmt.Printf("Serial download time: %v\n", serialTime)
fmt.Printf("Concurrent download time: %v\n", concurrentTime)
fmt.Printf("Speed improvement: %.1fx\n", float64(serialTime)/float64(concurrentTime))
```

### 示例4：条件停止 - 只下载前 50 个

```go
var count int
var countMu sync.Mutex
const maxCount = 50

err := igpsport.DownloadAllActivitiesWithConcurrency(
    igpsportsync.TCX,
    5,
    func(activity *igpsportsync.DownloadedActivity) bool {
        if activity.Error != nil {
            return true // 错误不计数，继续
        }
        
        countMu.Lock()
        count++
        currentCount := count
        countMu.Unlock()
        
        // 处理文件
        processFile(activity)
        
        fmt.Printf("Processed %d/%d\n", currentCount, maxCount)
        
        // 达到限制时停止
        if currentCount >= maxCount {
            fmt.Println("Reached max count limit, stopping...")
            return false // 返回 false 停止
        }
        
        return true
    },
)
```

### 示例5：动态调整并发数

```go
// 根据 CPU 核心数调整并发数
runtime.NumCPU() // 获取 CPU 核心数

concurrency := runtime.NumCPU() * 2 // 通常设为核心数的 2-4 倍
if concurrency > 16 {
    concurrency = 16 // 上限
}

fmt.Printf("Using %d concurrent workers\n", concurrency)

err := igpsport.DownloadAllActivitiesWithConcurrency(
    igpsportsync.FIT,
    concurrency,
    func(activity *igpsportsync.DownloadedActivity) bool {
        if activity.Error != nil {
            fmt.Printf("✗ Error: %v\n", activity.Error)
            return true
        }
        
        if err := uploadToOSS(activity); err != nil {
            fmt.Printf("✗ Upload failed: %v\n", err)
        }
        
        return true
    },
)
```

## 性能建议

### 并发数选择

| 场景 | 推荐并发数 | 说明 |
|------|----------|------|
| 网络不稳定 | 2-3 | 降低失败率 |
| 普通网络环境 | 5-8 | 平衡性能和稳定性 |
| 高速网络 | 10-16 | 充分利用带宽 |
| 云函数环境 | 3-5 | 考虑内存和超时限制 |

### 注意事项

1. **内存使用**: 每个并发任务在下载时会占用内存，大文件可能导致内存溢出
   ```go
   // 如果文件很大，考虑流式保存而不是全部加载到内存
   ```

2. **网络连接**: 并发过高可能被服务器限流，建议从较小的值开始逐步调整

3. **错误处理**: 
   - 单个下载失败不会影响其他下载
   - 返回 `false` 会立即停止所有 Worker

4. **超时设置**: HTTP 客户端超时为 30 秒，对于大文件可能需要调整
   ```go
   // 如需修改超时时间，可在初始化时设置
   ```

## 对比：串行 vs 并行

### 串行下载 (`DownloadAllActivities`)
```
Worker 1: [====] [====] [====] [====] ... (100 个文件，每个 2 秒)
总耗时: 200 秒
```

### 并行下载 (`DownloadAllActivitiesWithConcurrency`, 5 个并发)
```
Worker 1: [====]         [====]  ...
Worker 2:      [====]         [====] ...
Worker 3:           [====]         [====] ...
Worker 4:                [====]         ...
Worker 5:                     [====]     ...
总耗时: 40 秒 (理想情况下快 5 倍)
```

## 错误处理

```go
err := igpsport.DownloadAllActivitiesWithConcurrency(
    ext, 
    5, 
    callback,
)

if err != nil {
    // 这里的错误通常是获取分页列表失败
    fmt.Printf("Fatal error: %v\n", err)
}

// 单个文件的错误会在 callback 的 activity.Error 中返回
```

## 常见问题

**Q: 应该选择串行还是并行?**  
A: 如果需要上传到 OSS 等外部存储，建议使用并行以提高效率。对于本地保存，根据网络情况选择。

**Q: 并发数越多越好吗?**  
A: 不是。过高的并发可能导致：
- 内存溢出
- 被服务器限流（429 Too Many Requests）
- 反而降低效率

**Q: 如何动态调整并发数?**  
A: 可以根据网络响应时间或错误率动态调整，但需要额外的复杂性。建议从 5-8 开始测试。

**Q: 能否只并发某些步骤（如只并发上传）?**  
A: 可以在回调中实现。保持下载串行，在回调中异步上传（使用 goroutine 和 WaitGroup）。
