# 使用 Golang 重写 pyVideoTrans 的可行性分析：Wails vs Electron

---

## 一、核心问题：Go 能否替代 Python 做这件事？

### 1.1 原项目对 Python 的依赖分析

根据源码分析，pyVideoTrans 对 Python 生态的依赖分为三层：

| 层级 | 依赖项 | 是否可替代 |
|------|--------|-----------|
| **必须 Python** | PyTorch、Whisper、FunASR、pyannote-audio、ModelScope、transformers、ctranslate2 | 本地 AI 推理强依赖 Python |
| **可用 HTTP 替代** | OpenAI API、Google API、Azure API、DeepL、Deepgram 等 26+ 在线翻译/ASR/TTS 渠道 | Go 完全可以做 HTTP 调用 |
| **已是外部二进制** | FFmpeg、ffprobe、rubberband | Go 同样可以调用 subprocess |

### 1.2 结论：Go 重写的定位

**Go 无法完全替代 Python**，原因是本地 AI 推理（Whisper、pyannote 等）的生态完全在 Python/PyTorch 上。但可以采用**混合架构**：

| 方案 | 描述 | 可行性 |
|------|------|--------|
| **方案 A：纯 Go + 仅在线 API** | 放弃所有本地模型，全部使用在线 API | ✅ 完全可行，但丧失离线能力 |
| **方案 B：Go 主控 + Python sidecar** | Go 做 UI/编排/FFmpeg，AI 推理调用 Python 子进程或微服务 | ✅ 可行，架构更复杂 |
| **方案 C：Go + whisper.cpp/ggml** | 用 C/C++ 实现的模型替代 PyTorch 模型 | ⚠️ 部分可行，但 TTS/翻译/说话人分离无替代 |
| **方案 D：完全用 Python** | 不重写，保持现状 | ✅ 最省事 |

**推荐方案 A 或 B**，具体取决于是否需要离线能力。

---

## 二、Go 技术栈能覆盖的部分

### 2.1 Go 擅长的部分（原项目的 70%+ 工作量）

| 功能模块 | Go 实现方式 | 难度 |
|----------|------------|------|
| FFmpeg 音视频处理 | `os/exec` 调用 ffmpeg/ffprobe | ⭐ 简单 |
| 在线 ASR API 调用 | `net/http` + SDK | ⭐ 简单 |
| 在线翻译 API 调用 | `net/http` + SDK（OpenAI Go SDK 等）| ⭐ 简单 |
| 在线 TTS API 调用 | `net/http` + SDK | ⭐ 简单 |
| SRT/ASS/VTT 字幕解析 | Go 字符串处理，已有 `astisub` 等库 | ⭐ 简单 |
| 任务编排（流水线调度）| goroutine + channel，Go 天然优势 | ⭐ 简单 |
| CLI 命令行 | `cobra` / `urfave/cli` | ⭐ 简单 |
| 文件管理/临时目录 | 标准库 `os`/`filepath` | ⭐ 简单 |
| 并发任务控制 | goroutine + sync + context | ⭐ Go 最强项 |
| 进度推送 (WebSocket) | `gorilla/websocket` 或标准库 | ⭐ 简单 |

### 2.2 Go 做不了或很难的部分

| 功能模块 | 困难原因 | 替代方案 |
|----------|----------|----------|
| 本地 Whisper 推理 | PyTorch 模型，Go 无原生支持 | whisper.cpp (CGo 绑定) 或 API |
| 本地 TTS（F5-TTS/CosyVoice 等）| 全部依赖 PyTorch | 只用在线 TTS 或 Python sidecar |
| 说话人分离 (pyannote) | PyTorch 模型 | sherpa-onnx (有 Go 绑定) 或 API |
| 人声分离 (UVR) | sherpa-onnx ONNX 模型 | sherpa-onnx Go 绑定 ✅ 可行 |
| 降噪 (ModelScope) | PyTorch 模型 | 只能 API 或 Python sidecar |
| M2M100 本地翻译 | ctranslate2 + PyTorch | 只能 API 或 Python sidecar |
| 音频变速 (rubberband) | C 库 | 调用 rubberband 二进制 ✅ |

### 2.3 Go 生态中可用的 AI 相关库

| 库 | 功能 | 成熟度 |
|----|------|--------|
| [whisper.cpp](https://github.com/ggerganov/whisper.cpp) + Go 绑定 | 本地 Whisper 推理 | ⭐⭐⭐ 成熟 |
| [sherpa-onnx Go](https://github.com/k2-fsa/sherpa-onnx/tree/master/go-api-examples) | VAD/ASR/TTS/说话人分离 | ⭐⭐⭐ 成熟 |
| [go-openai](https://github.com/sashabaranov/go-openai) | OpenAI API 客户端 | ⭐⭐⭐ 成熟 |
| [onnxruntime-go](https://github.com/yalue/onnxruntime_go) | ONNX 模型推理 | ⭐⭐ 可用 |
| [go-anthropic](https://github.com/anthropics/anthropic-sdk-go) | Claude API | ⭐⭐⭐ 官方 |

---

## 三、桌面框架对比：Wails vs Electron

### 3.1 对比总表

| 维度 | Wails | Electron |
|------|-------|----------|
| **后端语言** | Go（原生）| Node.js（可通过 child_process 调 Go）|
| **前端** | 任意 Web 框架（Vue/React/Svelte）| 任意 Web 框架 |
| **打包体积** | ~10-15MB（不含 FFmpeg/模型）| ~80-150MB（Chromium 运行时）|
| **内存占用** | 低（50-100MB 基线）| 高（200-400MB 基线，Chromium）|
| **启动速度** | 快（<1s）| 慢（2-5s）|
| **跨平台** | Windows/macOS/Linux | Windows/macOS/Linux |
| **系统 WebView** | 使用系统 WebView2/WebKit | 自带 Chromium |
| **进程模型** | 单进程（Go + WebView）| 多进程（Main + Renderer）|
| **前后端通信** | Go 函数直接绑定到 JS | IPC（preload + contextBridge）|
| **生态成熟度** | ⭐⭐ 较新（v2 2022+）| ⭐⭐⭐⭐ 非常成熟 |
| **调试体验** | 较好（DevTools 可用）| 极好（完整 Chrome DevTools）|
| **自动更新** | 需自行实现 | electron-updater 一行搞定 |
| **系统托盘/通知** | 支持 | 支持 |
| **文件拖拽** | 支持 | 支持 |
| **原生菜单** | 支持 | 支持 |
| **社区/文档** | GitHub 12k+ stars，文档较全 | GitHub 115k+ stars，极其丰富 |
| **Windows 兼容性** | Win10+ (WebView2) | Win7+ |

### 3.2 针对本项目的关键对比

#### 视频翻译工具的特殊需求

| 需求 | Wails 表现 | Electron 表现 |
|------|-----------|--------------|
| **调用 FFmpeg 子进程** | Go `os/exec` 原生高效 | Node.js `child_process` 同样可以 |
| **长时间后台任务** | Go goroutine 天然优势 | Node.js 单线程，需 Worker |
| **大文件处理** | Go 内存控制精准 | Node.js 内存管理不如 Go |
| **实时进度推送** | Go channel → 前端事件 | IPC 或 WebSocket |
| **并发 API 调用** | Go 最强项 | Node.js async 也不错 |
| **调用 Python sidecar** | Go `os/exec` 启动 Python 子进程 | Node.js `child_process` 同理 |
| **视频预览 + 字幕同步** | 前端 video.js 均可 | 前端 video.js 均可 |
| **字幕编辑器** | 前端实现，无差异 | 前端实现，无差异 |
| **打包分发** | 体积小，用户体验好 | 体积大但用户接受度高 |

### 3.3 推荐：Wails

**理由**：

1. **Go 后端一致性**：既然选择 Go 重写，用 Wails 可以让前后端在同一个技术栈中，Go 函数直接暴露给前端调用，无需额外的 IPC 层

2. **资源消耗**：视频处理本身就很吃资源（FFmpeg 占大量 CPU/内存），不应再让 Electron 的 Chromium 额外消耗 200MB+ 内存

3. **打包体积**：视频工具已经因为 FFmpeg + 模型文件很大了，Wails 省下的 100MB+ 很有意义

4. **Go 原生能力**：goroutine 做任务编排 + `os/exec` 调 FFmpeg/Python 都是 Go 最擅长的场景，Wails 完全不阻碍这些

5. **足够成熟**：Wails v2 已经稳定，对于这种工具类桌面应用完全够用

**Electron 适合的场景**（本项目不太需要）：
- 需要复杂的多窗口管理
- 需要 Win7 兼容
- 团队只熟悉 Node.js
- 需要深度 Web 能力（如内嵌网页）

---

## 四、推荐技术架构

### 4.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                  前端 (Wails WebView)                     │
│  Vue 3 / React + TypeScript                              │
│  ┌─────────┐ ┌──────────────┐ ┌───────────────────────┐ │
│  │视频上传  │ │字幕编辑器     │ │视频预览 + 字幕同步    │ │
│  │(拖拽)    │ │(时间轴+文本) │ │(video.js/plyr)       │ │
│  └────┬────┘ └──────┬───────┘ └───────────┬───────────┘ │
│       │             │                     │              │
│       └─────────────┴─────────────────────┘              │
│                      │ Wails Bindings                    │
└──────────────────────┼───────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                   Go 后端 (Wails Backend)                 │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │              任务编排引擎 (Task Pipeline)          │   │
│  │  prepare → recogn → translate → tts → align →    │   │
│  │  assemble → done                                 │   │
│  │  (goroutine + context + channel)                 │   │
│  └──────────┬───────────────┬───────────────────────┘   │
│             │               │                            │
│  ┌──────────▼──────┐  ┌────▼────────────────────────┐   │
│  │  FFmpeg 封装层   │  │    AI 引擎适配层             │   │
│  │  (os/exec)      │  │                              │   │
│  │  - 音频提取     │  │  ASR:                        │   │
│  │  - 视频转码     │  │  ├ whisper.cpp (CGo/进程)    │   │
│  │  - 字幕烧录     │  │  ├ sherpa-onnx (Go绑定)     │   │
│  │  - 音视频合并   │  │  ├ OpenAI Whisper API       │   │
│  │  - 硬件加速检测 │  │  ├ Google/Azure/Deepgram..  │   │
│  └────────────────┘  │  │                            │   │
│                       │  │  Translation:              │   │
│  ┌────────────────┐  │  ├ go-openai (ChatGPT)       │   │
│  │ 字幕处理       │  │  ├ Google/DeepL/Azure API    │   │
│  │ - SRT/ASS/VTT │  │  ├ 本地 LLM (Ollama HTTP)    │   │
│  │ - 解析/生成    │  │  │                            │   │
│  │ - 双语合并     │  │  │  TTS:                      │   │
│  └────────────────┘  │  ├ Edge-TTS (HTTP)            │   │
│                       │  ├ OpenAI/Azure/ElevenLabs   │   │
│  ┌────────────────┐  │  ├ sherpa-onnx (本地TTS)      │   │
│  │ 音频处理       │  │  └ [Python sidecar 可选]      │   │
│  │ - rubberband   │  │                              │   │
│  │ - 变速对齐     │  └──────────────────────────────┘   │
│  │ - 静音填充     │                                      │
│  └────────────────┘                                      │
│                                                          │
│  ┌──────────────────────────────────────────────────┐   │
│  │        [可选] Python Sidecar 通信层               │   │
│  │  本地模型推理时启动 Python 子进程/Flask 微服务     │   │
│  │  - 本地 Whisper (faster-whisper)                  │   │
│  │  - pyannote 说话人分离                            │   │
│  │  - 本地 TTS (F5-TTS/CosyVoice)                   │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### 4.2 Go 后端核心模块设计

```
cmd/
├── desktop/main.go        # Wails 桌面入口
├── cli/main.go            # CLI 入口
internal/
├── pipeline/              # 任务编排
│   ├── task.go            # 任务基类接口
│   ├── vtv.go             # 视频翻译流水线
│   ├── stt.go             # 语音转录
│   ├── sts.go             # 字幕翻译
│   └── tts_task.go        # 配音任务
├── asr/                   # ASR 引擎适配 (策略模式)
│   ├── engine.go          # 接口定义
│   ├── whisper_cpp.go     # whisper.cpp 绑定
│   ├── openai_api.go      # OpenAI Whisper API
│   ├── google_stt.go      # Google Cloud STT
│   └── ...
├── translator/            # 翻译引擎适配
│   ├── engine.go
│   ├── openai.go
│   ├── google.go
│   ├── deepl.go
│   └── ...
├── tts/                   # TTS 引擎适配
│   ├── engine.go
│   ├── edge_tts.go
│   ├── openai_tts.go
│   └── ...
├── ffmpeg/                # FFmpeg 封装
│   ├── runner.go
│   ├── probe.go
│   ├── encoder.go         # 硬件加速检测
│   └── subtitle.go        # 字幕嵌入
├── subtitle/              # 字幕处理
│   ├── srt.go
│   ├── ass.go
│   └── vtt.go
├── align/                 # 音视频对齐
│   └── speed_rate.go
└── sidecar/               # [可选] Python 子进程管理
    └── python.go
frontend/
├── src/                   # Vue 3 / React 前端
│   ├── views/
│   │   ├── VideoTranslate.vue
│   │   ├── SubtitleEditor.vue
│   │   └── Settings.vue
│   └── components/
│       ├── VideoPlayer.vue
│       ├── TimelineEditor.vue
│       └── ProgressBar.vue
```

### 4.3 核心接口设计

```go
// ASR 引擎接口
type ASREngine interface {
    Recognize(ctx context.Context, audioPath string, opts ASROptions) ([]SubtitleEntry, error)
    Name() string
}

// 翻译引擎接口
type TranslateEngine interface {
    Translate(ctx context.Context, entries []SubtitleEntry, opts TranslateOptions) ([]SubtitleEntry, error)
    Name() string
}

// TTS 引擎接口
type TTSEngine interface {
    Synthesize(ctx context.Context, entries []SubtitleEntry, opts TTSOptions) ([]AudioSegment, error)
    Name() string
}

// 任务流水线
type Pipeline struct {
    ASR        ASREngine
    Translator TranslateEngine
    TTS        TTSEngine
    FFmpeg     *ffmpeg.Runner
    OnProgress func(stage string, percent int)
}

func (p *Pipeline) Run(ctx context.Context, input string) error {
    // prepare → recogn → translate → tts → align → assemble
}
```

---

## 五、最终建议

### 5.1 推荐路线

| 阶段 | 目标 | 技术选择 |
|------|------|----------|
| **MVP** | 在线 API 模式跑通全流程 | Go + Wails + 仅在线 API |
| **V2** | 补充本地 ASR | 集成 whisper.cpp（CGo 或进程调用）|
| **V3** | 补充本地 TTS + 说话人分离 | 集成 sherpa-onnx Go 绑定 |
| **V4** | 完整离线能力 | Python sidecar（F5-TTS/CosyVoice 等）|

### 5.2 框架选择总结

| 场景 | 推荐 |
|------|------|
| 追求轻量、性能、Go 技术栈一致性 | **Wails** ✅ |
| 团队熟悉 Node.js、需要丰富插件生态 | Electron |
| 需要 Web 端 + 桌面端共享代码 | Electron 或 Tauri |

### 5.3 风险提醒

1. **最大风险**：Go 的 AI 生态远不如 Python，本地模型能力会大幅缩水。如果核心卖点是「离线+本地模型」，不建议用 Go 重写
2. **适合 Go 重写的场景**：定位为「在线 API 聚合工具」+ 轻量桌面端，不强调本地推理
3. **Wails 风险**：WebView2 在 Windows 10 部分旧版本需要用户手动安装；macOS WebKit 的兼容性偶尔有 CSS 差异
4. **工作量**：即使只做在线 API 模式，原项目 24 种翻译 + 22 种 ASR + 32 种 TTS 的适配量也很大，建议 MVP 先做 5-8 个核心渠道
