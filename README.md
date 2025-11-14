# Rocky's Go Common Library

🚀 常用的 Go 工具庫集合，提供日誌、設定、SSH 客戶端等通用功能。

## 📦 安裝

```bash
go get github.com/rocky-chen/rocky-go-common
```

## 📚 模組說明

### Logger - 彩色日誌輸出

```go
import "github.com/rocky-chen/rocky-go-common/logger"

log := logger.New(true) // verbose mode
log.Info("應用程式啟動中...")
log.Success("連線成功！")
log.Warning("注意：即將進行危險操作")
log.Error("發生錯誤：%v", err)
log.Debug("除錯資訊（只在 verbose 模式顯示）")
```

**特色：**
- ✅ 彩色輸出（Info=藍、Success=綠、Warning=黃、Error=紅、Debug=紫）
- ✅ 時間戳記
- ✅ Verbose 模式控制除錯輸出

---

### Config - 設定檔管理

```go
import "github.com/rocky-chen/rocky-go-common/config"

type AppConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

var cfg AppConfig

// 自動搜尋設定檔（當前目錄、$HOME/.myapp/、/etc/myapp/）
err := config.LoadConfig("myapp", &cfg)

// 或指定設定檔路徑
err := config.LoadConfigFromFile("/path/to/config.yaml", &cfg)
```

**支援格式：**
- ✅ YAML
- ✅ 環境變數自動覆蓋
- ✅ 多路徑自動搜尋

**範例設定檔 (config.yaml)：**
```yaml
host: localhost
port: 8080
```

---

### SSH - SSH 客戶端工具

```go
import "github.com/rocky-chen/rocky-go-common/ssh"

// 建立 SSH 連線
client, err := ssh.NewClient(ssh.Config{
    Host:     "192.168.1.100",
    Port:     22,
    User:     "rocky",
    Password: "password", // 或使用 KeyFile: "/path/to/key"
    Timeout:  30 * time.Second,
})
if err != nil {
    log.Fatal("SSH 連線失敗: %v", err)
}
defer client.Close()

// 執行遠端命令
output, err := client.RunCommand("ls -la /tmp")

// 下載檔案
err = client.DownloadFile("/remote/file.txt", "./local/file.txt")

// 上傳檔案
err = client.UploadFile("./local/file.txt", "/remote/file.txt")

// 檢查檔案是否存在
exists, err := client.FileExists("/remote/file.txt")

// 列出目錄內容
files, err := client.ListFiles("/remote/path")
```

**功能：**
- ✅ 密碼認證
- ✅ 金鑰認證
- ✅ 遠端命令執行
- ✅ 檔案上傳/下載（SCP）
- ✅ 檔案操作（檢查存在、列出目錄）

---

### TimeUtil - 時間處理工具

```go
import "github.com/rocky-chen/rocky-go-common/utils/timeutil"

// 解析時間範圍 "14-15" 或 "14:30-15:45"
start, end, err := timeutil.ParseTimeRange("14-15")

// 格式化時間長度
duration := 125 * time.Second
fmt.Println(timeutil.FormatDuration(duration)) // "2.1m"

// 解析日期（支援多種格式）
date, err := timeutil.ParseDate("2024-11-12")

// 檢查是否為今天/昨天
if timeutil.IsToday(someTime) {
    // ...
}

// 取得當天開始/結束時間
start := timeutil.StartOfDay(time.Now())  // 00:00:00
end := timeutil.EndOfDay(time.Now())      // 23:59:59
```

**功能：**
- ✅ 時間範圍解析（支援小時和分鐘）
- ✅ 人性化時間長度顯示
- ✅ 多格式日期解析
- ✅ 日期判斷工具

---

## 🏗️ 專案結構

```
rocky-go-common/
├── logger/           # 日誌模組
│   └── logger.go
├── config/           # 設定管理
│   └── loader.go
├── ssh/              # SSH 客戶端
│   └── client.go
├── utils/
│   └── timeutil/     # 時間工具
│       └── timeutil.go
├── go.mod
└── README.md
```

---

## 🎯 設計原則

### 1. 最小依賴
- 優先使用 Go 標準庫
- 只在必要時引入外部依賴（如 viper、color）

### 2. 簡單易用
- API 設計直觀
- 提供合理的預設值
- 錯誤訊息清晰

### 3. 零配置啟動
- 每個模組可獨立使用
- 不需要複雜的初始化

---

## 📝 使用範例

### 完整的應用程式範例

```go
package main

import (
    "github.com/rocky-chen/rocky-go-common/logger"
    "github.com/rocky-chen/rocky-go-common/config"
    "github.com/rocky-chen/rocky-go-common/ssh"
)

type AppConfig struct {
    SSHHost string `mapstructure:"ssh_host"`
    SSHUser string `mapstructure:"ssh_user"`
    SSHPass string `mapstructure:"ssh_pass"`
}

func main() {
    // 初始化日誌
    log := logger.New(true)
    log.Info("應用程式啟動")
    
    // 載入設定
    var cfg AppConfig
    if err := config.LoadConfig("myapp", &cfg); err != nil {
        log.Fatal("載入設定失敗: %v", err)
    }
    
    // 連線 SSH
    client, err := ssh.NewClient(ssh.Config{
        Host:     cfg.SSHHost,
        Port:     22,
        User:     cfg.SSHUser,
        Password: cfg.SSHPass,
    })
    if err != nil {
        log.Fatal("SSH 連線失敗: %v", err)
    }
    defer client.Close()
    
    log.Success("連線成功！")
    
    // 執行命令
    output, err := client.RunCommand("hostname")
    if err != nil {
        log.Error("命令執行失敗: %v", err)
        return
    }
    
    log.Info("主機名稱: %s", output)
}
```

---

## 🔧 開發

### 執行測試

```bash
go test ./...
```

### 安裝依賴

```bash
go mod tidy
```

### 建置

```bash
go build ./...
```

---

## 📋 版本歷史

### v1.0.0 (2024-11-13)
- ✅ 初始版本發布
- ✅ Logger 模組
- ✅ Config 模組
- ✅ SSH 模組
- ✅ TimeUtil 模組

---

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

---

## 📄 授權

MIT License

---

## 👤 作者

Rocky Chen

---

## 🔗 相關資源

- [Go Modules 官方文檔](https://go.dev/doc/modules)
- [Effective Go](https://go.dev/doc/effective_go)

---

**最後更新**: 2024-11-13
