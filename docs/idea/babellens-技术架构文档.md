# BabelLens — 功能分析与技术架构文档

> 技术栈：Go + Python + Wails (React) | 双模式：GUI + CLI | 完全离线 + 在线 API

---

## 一、产品定位

BabelLens 是一款基于 AI 的视频字幕翻译与配音工具，提供从语音识别到字幕翻译、AI 配音、音视频合成的完整工作流。采用 Go 主控 + Python AI 推理的混合架构，同时支持桌面 GUI 和命令行 CLI 两种使用模式。

---

## 二、核心功能清单

| # | 功能 | 描述 | 缩写 |
|---|------|------|------|
| 1 | 视频翻译 | A 语言视频 → B 语言配音 + B 语言字幕视频 | VTV |
| 2 | 语音转录 / 字幕生成 | 音频/视频 → 带时间戳的字幕文件 (SRT/ASS/VTT) | STT |
| 3 | 字幕翻译 | SRT 字幕文件 A 语言 → B 语言 | STS |
| 4 | AI 配音 | 字幕 → 语音合成，支持多角色配音和声音克隆 | TTS |
| 5 | 音视频对齐 | 配音时长与原始视频时间轴同步 | ALIGN |
| 6 | 人声/背景音分离 | 视频/音频 → 人声 + 背景音乐 | SEP |
| 7 | 降噪处理 | 语音预处理降噪，提升 ASR 准确率 | DNR |
| 8 | 辅助工具集 | 音频提取、字幕转换/合并/覆盖、音视频合并、水印等 | TOOLS |

---

## 三、功能详细设计

### 3.1 视频翻译 (VTV)

完整流水线，串联所有子功能：

```
输入: 视频文件 (MP4/MOV/AVI/MKV...)
    │
    ▼
┌─ prepare ─────────────────────────────────┐
│  FFmpeg 提取音频 → 16kHz mono WAV          │
│  FFmpeg 分离无声视频                        │
│  获取视频元信息 (分辨率/帧率/编码/时长)      │
│  [可选] 人声/背景音分离 (SEP)               │
│  [可选] 降噪处理 (DNR)                      │
└───────────────────────────────────────────┘
    │
    ▼
┌─ recogn (STT) ────────────────────────────┐
│  ASR 引擎语音识别                           │
│  输出: 源语言 SRT 字幕 (含时间戳)            │
│  [可选] 说话人分离标注                       │
│  ── 暂停点: 用户可手动校对字幕 ──            │
└───────────────────────────────────────────┘
    │
    ▼
┌─ translate (STS) ─────────────────────────┐
│  翻译引擎: 源语言字幕 → 目标语言字幕         │
│  ── 暂停点: 用户可手动校对翻译 ──            │
└───────────────────────────────────────────┘
    │
    ▼
┌─ dubbing (TTS) ───────────────────────────┐
│  TTS 引擎: 目标语言字幕 → 语音片段           │
│  支持多角色: 不同说话人 → 不同音色            │
│  [可选] 声音克隆: 参考音频 → 克隆音色         │
│  ── 暂停点: 用户可试听并调整 ──              │
└───────────────────────────────────────────┘
    │
    ▼
┌─ align (ALIGN) ───────────────────────────┐
│  音频加速: rubberband 变速不变调             │
│  视频慢速: FFmpeg setpts 滤镜               │
│  静音填充 + 时间轴对齐                       │
└───────────────────────────────────────────┘
    │
    ▼
┌─ assemble ────────────────────────────────┐
│  FFmpeg 合并: 无声视频 + 配音 + 背景音乐     │
│  字幕嵌入: 硬字幕/软字幕/双语字幕             │
│  编码输出: H.264/H.265 + 硬件加速检测        │
└───────────────────────────────────────────┘
    │
    ▼
输出: 翻译后的视频 + 字幕文件
```

**CLI 示例**:
```bash
babellens vtv \
  --input video.mp4 \
  --source-lang zh \
  --target-lang en \
  --asr-engine faster-whisper \
  --asr-model large-v3 \
  --translate-engine ollama \
  --tts-engine edge-tts \
  --tts-voice en-US-AriaNeural \
  --subtitle-type hard \
  --output output.mp4
```

### 3.2 语音转录 / 字幕生成 (STT)

```
输入: 音频/视频文件
    │
    ▼
FFmpeg 转换 → 16kHz mono WAV
    │
    ▼
[可选] 降噪 (DNR)
    │
    ▼
VAD 语音活动检测 → 切分语音片段
    │
    ▼
ASR 引擎识别 → 文本 + word-level 时间戳
    │
    ▼
[可选] 说话人分离 → 标注 spk0/spk1/...
    │
    ▼
[可选] LLM 后处理 → 智能断句 + 标点修正
    │
    ▼
输出: SRT / ASS / VTT / TXT
```

**ASR 引擎架构（本地 + 在线 + 自定义）**:

```
ASREngine 接口
├── 本地引擎
│   ├── whisper-cpp         (Go 进程调用 whisper.cpp 二进制)
│   ├── faster-whisper       (Python Sidecar, ctranslate2)
│   ├── sherpa-onnx          (Go 原生绑定)
│   ├── funasr               (Python Sidecar, 阿里达摩院)
│   └── openai-whisper       (Python Sidecar, PyTorch)
├── 在线引擎
│   ├── openai-api           (Go HTTP, OpenAI Whisper API)
│   ├── google-speech        (Go HTTP, Google Cloud STT)
│   ├── azure-speech         (Go HTTP, Azure Cognitive)
│   ├── deepgram             (Go HTTP, Deepgram API)
│   └── volcengine           (Go HTTP, 字节火山引擎)
└── 自定义引擎
    └── custom-api            (Go HTTP, 用户配置 URL + 请求格式)
```

**CLI 示例**:
```bash
babellens stt \
  --input audio.wav \
  --engine faster-whisper \
  --model large-v3 \
  --lang auto \
  --diarize \
  --output subtitle.srt
```

### 3.3 字幕翻译 (STS)

```
输入: SRT 字幕文件
    │
    ▼
解析字幕 → []SubtitleEntry{start, end, text}
    │
    ▼
翻译引擎 (批量翻译, 保留上下文)
    │
    ▼
输出: 目标语言 SRT (单语/双语)
```

**翻译引擎架构**:

```
TranslateEngine 接口
├── LLM 翻译 (上下文理解好)
│   ├── ollama-local         (Go HTTP → 本地 Ollama, 完全离线)
│   ├── openai-chatgpt       (Go HTTP → OpenAI API)
│   ├── deepseek             (Go HTTP → DeepSeek API)
│   ├── gemini               (Go HTTP → Google Gemini)
│   └── local-llm            (Go HTTP → 任意 OpenAI 兼容 API)
├── 传统机翻 (速度快)
│   ├── google               (Go HTTP → Google Translate)
│   ├── microsoft            (Go HTTP → Bing Translate)
│   ├── deepl                (Go HTTP → DeepL API)
│   └── baidu                (Go HTTP → 百度翻译)
└── 自定义
    └── custom-api            (Go HTTP → 用户配置)
```

### 3.4 AI 配音 (TTS)

```
输入: 目标语言字幕 + 角色分配
    │
    ▼
构建配音队列 [{text, role, start, end, ref_audio?}, ...]
    │
    ▼
TTS 引擎: 逐条合成 → WAV 音频片段
    │
    ▼
输出: 每条字幕对应的 WAV 文件
```

**TTS 引擎架构**:

```
TTSEngine 接口
├── 本地引擎 (无需网络)
│   ├── sherpa-piper          (Go 原生, Piper TTS, 质量良好)
│   ├── qwen3-tts             (Python Sidecar, 高质量 + 克隆)
│   ├── f5-tts                (Python Sidecar, 声音克隆)
│   ├── cosyvoice             (Python Sidecar, 声音克隆)
│   └── chattts               (Python Sidecar, 中英文)
├── 在线引擎
│   ├── edge-tts              (Go HTTP, 微软免费, 推荐默认)
│   ├── openai-tts            (Go HTTP, OpenAI TTS)
│   ├── azure-tts             (Go HTTP, Azure 语音)
│   └── elevenlabs            (Go HTTP, ElevenLabs)
└── 自定义
    └── custom-api             (Go HTTP, 用户配置)
```

**声音克隆流程**:
```
用户上传参考音频 (3-10s)
    │
    ▼
Python Sidecar: F5-TTS/CosyVoice/Qwen3-TTS 提取声纹
    │
    ▼
合成时传入声纹 → 生成克隆语音
```

### 3.5 音视频对齐 (ALIGN)

**纯 Go 实现**，无需 Python。

**对齐策略**:

```
对于每条字幕的配音:

配音时长 <= 字幕时长
    → 无需处理, 末尾填充静音

配音时长 > 字幕时长, 差距 <= 20%
    → 音频加速 (rubberband 变速不变调)

配音时长 > 字幕时长, 差距 > 20%
    → 音频加速 + 视频慢速 各承担一半
    → 音频: rubberband
    → 视频: FFmpeg setpts=X*PTS
```

**拼接流程**:
```
1. 第一条字幕前如有空隙 → 填充静音
2. 逐条拼接配音片段 + 间隔静音
3. 最终音频长度 = 视频长度 (不足则补静音, 超出则视频定格延长)
4. 字幕时间轴同步调整
```

### 3.6 人声/背景音分离 (SEP)

```
输入: 音频/视频文件
    │
    ▼
[Go] FFmpeg 提取音频 → 44.1kHz WAV
    │
    ▼
[Go] sherpa-onnx UVR-MDX-NET 模型推理 (Go 原生绑定, 仅 CPU)
    │
    ▼
输出:
├── vocal.wav        (人声)
└── instrument.wav   (背景音乐)
```

纯 Go 实现，sherpa-onnx 有官方 Go 绑定，无需 Python。

### 3.7 降噪处理 (DNR)

```
方案 A (Go 原生, 基础):
  FFmpeg arnndn 滤镜 → 基础降噪, 质量一般

方案 B (Python Sidecar, 高质量):
  ModelScope FRCRN 模型 → 深度学习降噪, 质量优秀
```

两种方案并存，用户可选择。无 Python 时降级到 FFmpeg 方案。

### 3.8 辅助工具集 (TOOLS)

全部用 **Go + FFmpeg** 实现，无需 Python。

| 工具 | CLI 命令 | 实现方式 |
|------|----------|----------|
| 从视频提取音频 | `babellens tools extract-audio` | FFmpeg `-vn -acodec` |
| 字幕格式转换 | `babellens tools convert-sub` | Go 字幕解析 SRT↔ASS↔VTT 互转 |
| 字幕合并 | `babellens tools merge-sub` | Go 读取多个 SRT，按时间轴合并 |
| 字幕配音 | `babellens tts` | 复用 TTS Pipeline |
| 多角色配音 | `babellens tts --roles` | TTS Pipeline + 说话人角色映射 |
| 语音识别 | `babellens stt` | 复用 STT Pipeline |
| 字幕覆盖/嵌入 | `babellens tools burn-sub` | FFmpeg subtitles 滤镜 (硬) / `-c:s` (软) |
| 文字匹配 | `babellens tools match-text` | Go 文本 ↔ 时间轴对齐算法 |
| 视频+音频合并 | `babellens tools mux` | FFmpeg `-i video -i audio -c copy` |
| 视频+字幕合并 | `babellens tools burn-sub` | 同字幕覆盖 |
| 视频水印 | `babellens tools watermark` | FFmpeg overlay 滤镜 (图片/文字) |

**CLI 示例**:

```bash
# 提取音频
babellens tools extract-audio --input video.mp4 --output audio.wav

# 字幕格式转换
babellens tools convert-sub --input sub.srt --format ass --output sub.ass

# 合并字幕 (双语)
babellens tools merge-sub --input zh.srt --input en.srt --output dual.srt

# 字幕烧录进视频
babellens tools burn-sub --input video.mp4 --sub sub.srt --type hard --output out.mp4

# 视频+音频合并
babellens tools mux --video video.mp4 --audio dubbing.wav --output out.mp4

# 水印
babellens tools watermark --input video.mp4 --text "BabelLens" --position bottom-right --output out.mp4
```

**对应 Go 模块**:

```
internal/tools/
├── extract_audio.go       # 提取音频
├── convert_subtitle.go    # 字幕格式转换
├── merge_subtitle.go      # 字幕合并
├── burn_subtitle.go       # 字幕烧录
├── mux.go                 # 音视频混流
├── watermark.go           # 水印
└── match_text.go          # 文字匹配
```

---

## 四、系统架构

### 4.1 整体架构图

```
┌──────────────────────────────────────────────────────────────────┐
│                           用户界面层                              │
│                                                                   │
│  ┌─────────────────────────────┐  ┌────────────────────────────┐ │
│  │     GUI (Wails + React)     │  │     CLI (cobra)            │ │
│  │                             │  │                            │ │
│  │  ┌───────┐ ┌─────────────┐ │  │  babellens vtv --input ... │ │
│  │  │视频上传│ │字幕编辑器    │ │  │  babellens stt --input ... │ │
│  │  │(拖拽)  │ │(时间轴+文本)│ │  │  babellens sts --input ... │ │
│  │  ├───────┤ ├─────────────┤ │  │  babellens tts --input ... │ │
│  │  │视频预览│ │配音设置      │ │  │  babellens sep --input ... │ │
│  │  │+字幕   │ │(角色/音色)  │ │  │                            │ │
│  │  ├───────┤ ├─────────────┤ │  │  支持 JSON/表格 输出        │ │
│  │  │进度面板│ │引擎设置      │ │  │  支持 --config 配置文件     │ │
│  │  └───────┘ └─────────────┘ │  │                            │ │
│  └─────────────┬───────────────┘  └──────────────┬─────────────┘ │
│                │ Wails Bindings                   │ 直接调用       │
└────────────────┼──────────────────────────────────┼───────────────┘
                 │                                  │
                 ▼                                  ▼
┌──────────────────────────────────────────────────────────────────┐
│                        Go 核心层                                  │
│                                                                   │
│  ┌────────────────────────────────────────────────────────────┐  │
│  │                    任务编排引擎 (pipeline)                   │  │
│  │                                                             │  │
│  │  Pipeline{                                                  │  │
│  │    stages: [prepare, recogn, translate, dub, align, assemble]│  │
│  │    onProgress: func(stage, percent)                         │  │
│  │    pausePoints: [after_recogn, after_translate, after_dub]  │  │
│  │  }                                                          │  │
│  │                                                             │  │
│  │  - goroutine + context.Context 控制生命周期                   │  │
│  │  - channel 传递阶段间数据                                     │  │
│  │  - 支持暂停/恢复/取消                                         │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                   │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────┐  │
│  │ ASR 适配层    │ │ 翻译适配层    │ │ TTS 适配层                │  │
│  │ (asr.Engine)  │ │ (trans.Engin)│ │ (tts.Engine)             │  │
│  │              │ │              │ │                          │  │
│  │ 本地/在线/   │ │ LLM/机翻/   │ │ 本地/在线/自定义          │  │
│  │ 自定义       │ │ 自定义       │ │ + 声音克隆               │  │
│  └──────┬───────┘ └──────┬───────┘ └──────────┬───────────────┘  │
│         │                │                     │                  │
│  ┌──────▼────────────────▼─────────────────────▼───────────────┐ │
│  │                    引擎路由层                                 │ │
│  │  在线引擎 → Go net/http 直接调用                               │ │
│  │  Go 本地引擎 → whisper.cpp / sherpa-onnx (进程调用/Go绑定)    │ │
│  │  Python 引擎 → Python Sidecar HTTP 调用                      │ │
│  └──────────────────────────┬──────────────────────────────────┘ │
│                             │                                    │
│  ┌──────────────────────────▼──────────────────────────────────┐ │
│  │                  基础能力层                                   │ │
│  │                                                              │ │
│  │  ┌──────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐ │ │
│  │  │ FFmpeg   │ │ 字幕处理  │ │ 音频对齐    │ │ Python       │ │ │
│  │  │ 封装     │ │ SRT/ASS/ │ │ rubberband │ │ Sidecar      │ │ │
│  │  │ (os/exec)│ │ VTT 解析 │ │ (os/exec)  │ │ 管理器       │ │ │
│  │  └──────────┘ └──────────┘ └────────────┘ └──────────────┘ │ │
│  │  ┌──────────┐ ┌──────────┐ ┌────────────┐ ┌──────────────┐ │ │
│  │  │ 配置管理  │ │ 模型管理  │ │ GPU 检测   │ │ 日志系统     │ │ │
│  │  │ (TOML)   │ │ (下载/   │ │ (CUDA/     │ │ (zerolog)   │ │ │
│  │  │          │ │  缓存)   │ │  Metal)    │ │              │ │ │
│  │  └──────────┘ └──────────┘ └────────────┘ └──────────────┘ │ │
│  └──────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Python Sidecar (可选)                           │
│                                                                   │
│  FastAPI HTTP 微服务, 监听 127.0.0.1:{port}                      │
│                                                                   │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │  POST /asr/faster-whisper    → faster-whisper 推理           │ │
│  │  POST /asr/openai-whisper    → openai-whisper 推理           │ │
│  │  POST /asr/funasr            → FunASR 推理                   │ │
│  │  POST /tts/qwen3             → Qwen3-TTS 合成                │ │
│  │  POST /tts/f5                → F5-TTS 合成 (声音克隆)        │ │
│  │  POST /tts/cosyvoice         → CosyVoice 合成 (声音克隆)     │ │
│  │  POST /denoise               → FRCRN 降噪                    │ │
│  │  POST /diarize/pyannote      → pyannote 说话人分离            │ │
│  │  POST /punctuate             → 标点恢复                       │ │
│  │  GET  /health                → 健康检查                       │ │
│  │  GET  /gpu                   → GPU 信息                       │ │
│  │  SSE  /progress/{task_id}    → 实时进度推送                   │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                   │
│  依赖: torch, faster-whisper, pyannote-audio, modelscope,        │
│        funasr, qwen-tts, soundfile, librosa 等                    │
│  打包: PyInstaller → 独立可执行文件 (worker.exe / worker)         │
└──────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌──────────────────────────────────────────────────────────────────┐
│                      外部二进制                                    │
│  ├── ffmpeg / ffprobe       音视频处理                             │
│  ├── rubberband             音频变速不变调 (可选)                   │
│  ├── whisper.cpp 二进制      本地 ASR (可选, 替代 Python Whisper)  │
│  └── ollama                 本地 LLM 翻译 (独立服务)               │
└──────────────────────────────────────────────────────────────────┘
```

### 4.2 Go ↔ Python 通信协议

采用 **HTTP JSON API**，Python Sidecar 作为本地微服务。

**请求格式**:
```json
POST /asr/faster-whisper
{
  "audio_path": "/tmp/babellens/xxx/audio.wav",
  "model": "large-v3",
  "language": "auto",
  "device": "cuda",
  "options": {
    "beam_size": 5,
    "vad_filter": true,
    "initial_prompt": ""
  }
}
```

**响应格式**:
```json
{
  "success": true,
  "data": {
    "subtitles": [
      {"index": 1, "start": 0.0, "end": 2.5, "text": "Hello world"},
      {"index": 2, "start": 3.0, "end": 5.8, "text": "How are you"}
    ],
    "language": "en",
    "duration": 10.5
  },
  "error": null
}
```

**进度推送 (SSE)**:
```
GET /progress/{task_id}

data: {"stage": "asr", "percent": 35, "message": "Transcribing segment 7/20"}
data: {"stage": "asr", "percent": 70, "message": "Transcribing segment 14/20"}
data: {"stage": "asr", "percent": 100, "message": "Done"}
```

**Go 端 Sidecar 管理器生命周期**:
```
应用启动
  │
  ├─ 检测 Python Worker 是否存在
  │   ├── 存在 → 记录可用
  │   └── 不存在 → 标记为不可用, 仅启用 Go 原生引擎
  │
  ├─ 用户触发需要 Python 的功能时
  │   ├── 查找空闲端口
  │   ├── 启动 worker 进程 (worker.exe --port {port})
  │   ├── 轮询 /health 等待就绪 (超时 30s)
  │   └── 就绪后开始接受任务
  │
  ├─ 空闲超时 (5 分钟无任务)
  │   └── 优雅关闭 Python Worker, 释放 GPU 内存
  │
  └─ 应用退出
      └── 发送 shutdown 信号, 等待 Worker 退出, 强杀兜底
```

---

## 五、项目目录结构

```
babellens/
├── main.go                          # Wails 桌面入口
├── cli/
│   └── main.go                      # CLI 独立入口
├── app.go                           # Wails App 绑定 (暴露给前端)
├── wails.json                       # Wails 配置
├── build/                           # Wails 构建产物
│
├── internal/                        # Go 核心业务 (不对外暴露)
│   ├── pipeline/                    # 任务编排
│   │   ├── pipeline.go              # Pipeline 结构体 + 流水线调度
│   │   ├── vtv.go                   # 视频翻译任务
│   │   ├── stt.go                   # 语音转录任务
│   │   ├── sts.go                   # 字幕翻译任务
│   │   ├── tts_task.go              # 配音任务
│   │   └── task.go                  # Task 接口 + 基础结构
│   │
│   ├── asr/                         # ASR 引擎适配层
│   │   ├── engine.go                # Engine 接口定义
│   │   ├── registry.go              # 引擎注册表 (工厂模式)
│   │   ├── whisper_cpp.go           # whisper.cpp 进程调用
│   │   ├── sherpa.go                # sherpa-onnx Go 绑定
│   │   ├── openai_api.go            # OpenAI Whisper API
│   │   ├── google_stt.go            # Google Cloud STT
│   │   ├── deepgram.go              # Deepgram API
│   │   ├── python_sidecar.go        # 转发到 Python (faster-whisper 等)
│   │   └── custom_api.go            # 自定义 API
│   │
│   ├── translator/                  # 翻译引擎适配层
│   │   ├── engine.go                # Engine 接口定义
│   │   ├── registry.go              # 引擎注册表
│   │   ├── ollama.go                # Ollama 本地 LLM
│   │   ├── openai.go                # OpenAI ChatGPT
│   │   ├── deepseek.go              # DeepSeek
│   │   ├── google.go                # Google Translate
│   │   ├── deepl.go                 # DeepL
│   │   └── custom_api.go            # 自定义 API
│   │
│   ├── tts/                         # TTS 引擎适配层
│   │   ├── engine.go                # Engine 接口定义
│   │   ├── registry.go              # 引擎注册表
│   │   ├── edge_tts.go              # Edge-TTS (免费)
│   │   ├── sherpa_piper.go          # sherpa-onnx Piper (本地)
│   │   ├── openai_tts.go            # OpenAI TTS
│   │   ├── python_sidecar.go        # 转发到 Python (Qwen3/F5/CosyVoice)
│   │   └── custom_api.go            # 自定义 API
│   │
│   ├── ffmpeg/                      # FFmpeg 封装层
│   │   ├── runner.go                # 执行器 (subprocess 管理)
│   │   ├── probe.go                 # ffprobe 元信息查询
│   │   ├── audio.go                 # 音频提取/转换/混音
│   │   ├── video.go                 # 视频转码/裁切/合并
│   │   ├── subtitle.go              # 字幕嵌入 (硬/软/双语)
│   │   └── hwaccel.go               # 硬件加速检测 (NVENC/QSV/AMF)
│   │
│   ├── subtitle/                    # 字幕处理
│   │   ├── srt.go                   # SRT 解析/生成
│   │   ├── ass.go                   # ASS 解析/生成
│   │   ├── vtt.go                   # VTT 解析/生成
│   │   └── types.go                 # SubtitleEntry 结构体
│   │
│   ├── align/                       # 音视频对齐
│   │   ├── speed.go                 # 音频加速 (rubberband)
│   │   ├── video_slow.go            # 视频慢速 (setpts)
│   │   └── mixer.go                 # 音频拼接 + 静音填充
│   │
│   ├── separate/                    # 人声分离
│   │   └── uvr.go                   # sherpa-onnx UVR (Go 绑定)
│   │
│   ├── denoise/                     # 降噪
│   │   ├── ffmpeg.go                # FFmpeg 滤镜降噪 (基础)
│   │   └── python_sidecar.go        # FRCRN 降噪 (高质量)
│   │
│   ├── tools/                       # 辅助工具集
│   │   ├── extract_audio.go         # 从视频提取音频
│   │   ├── convert_subtitle.go      # 字幕格式转换 SRT↔ASS↔VTT
│   │   ├── merge_subtitle.go        # 字幕合并 (双语)
│   │   ├── burn_subtitle.go         # 字幕烧录/嵌入视频
│   │   ├── mux.go                   # 音视频混流
│   │   ├── watermark.go             # 视频水印 (文字/图片)
│   │   └── match_text.go            # 文字与时间轴匹配
│   │
│   ├── diarize/                     # 说话人分离
│   │   ├── sherpa.go                # sherpa-onnx (Go 绑定)
│   │   └── python_sidecar.go        # pyannote (Python)
│   │
│   ├── sidecar/                     # Python Sidecar 管理
│   │   ├── manager.go               # 进程启动/停止/健康检查
│   │   ├── client.go                # HTTP 客户端封装
│   │   └── types.go                 # 请求/响应结构体
│   │
│   ├── config/                      # 配置管理
│   │   ├── config.go                # 全局配置结构体
│   │   ├── loader.go                # TOML 配置加载
│   │   └── defaults.go              # 默认值
│   │
│   ├── model/                       # 模型管理
│   │   ├── download.go              # 模型下载 (HuggingFace/自定义源)
│   │   ├── cache.go                 # 模型缓存目录管理
│   │   └── registry.go              # 已下载模型列表
│   │
│   └── util/                        # 工具函数
│       ├── gpu.go                   # GPU 检测 (CUDA/Metal)
│       ├── lang.go                  # 语言代码映射
│       ├── progress.go              # 进度回调封装
│       └── tempdir.go               # 临时目录管理
│
├── frontend/                        # React 前端 (Wails 标准结构)
│   ├── src/
│   │   ├── App.tsx
│   │   ├── pages/
│   │   │   ├── VideoTranslate.tsx   # 视频翻译页
│   │   │   ├── Transcribe.tsx       # 语音转录页
│   │   │   ├── SubtitleTranslate.tsx# 字幕翻译页
│   │   │   ├── Dubbing.tsx          # 配音页
│   │   │   ├── Tools.tsx            # 工具集 (分离/降噪/提取/合并/水印等)
│   │   │   └── Settings.tsx         # 设置页
│   │   ├── components/
│   │   │   ├── VideoPlayer.tsx      # 视频播放器 + 字幕同步
│   │   │   ├── SubtitleEditor.tsx   # 字幕编辑器 (时间轴)
│   │   │   ├── ProgressPanel.tsx    # 任务进度面板
│   │   │   ├── EngineSelector.tsx   # 引擎选择器
│   │   │   └── FileDropZone.tsx     # 文件拖拽区域
│   │   ├── hooks/
│   │   │   ├── useTask.ts           # 任务状态管理
│   │   │   └── useWails.ts          # Wails binding 封装
│   │   └── wailsjs/                 # Wails 自动生成的 JS 绑定
│   ├── index.html
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── python/                          # Python Sidecar 项目
│   ├── worker.py                    # FastAPI 主入口
│   ├── routers/
│   │   ├── asr.py                   # ASR 路由
│   │   ├── tts.py                   # TTS 路由
│   │   ├── denoise.py               # 降噪路由
│   │   └── diarize.py               # 说话人分离路由
│   ├── engines/
│   │   ├── faster_whisper_engine.py
│   │   ├── openai_whisper_engine.py
│   │   ├── funasr_engine.py
│   │   ├── qwen3_tts_engine.py
│   │   ├── f5_tts_engine.py
│   │   └── cosyvoice_engine.py
│   ├── pyproject.toml               # uv 依赖管理
│   └── build.py                     # PyInstaller 打包脚本
│
├── configs/
│   └── default.toml                 # 默认配置文件
├── scripts/
│   ├── build.sh                     # 构建脚本
│   └── build-python.sh              # Python Worker 打包脚本
└── README.md
```

---

## 六、核心接口设计

### 6.1 ASR 引擎接口

```go
// internal/asr/engine.go

type Engine interface {
    // 语音识别
    Recognize(ctx context.Context, req RecognizeRequest) (*RecognizeResult, error)
    // 引擎名称
    Name() string
    // 引擎类型: local, online, custom
    Type() EngineType
    // 是否可用 (模型已下载 / API 已配置)
    Available() bool
}

type RecognizeRequest struct {
    AudioPath  string            // 音频文件路径
    Language   string            // 语言代码 ("auto", "zh", "en", ...)
    Model      string            // 模型名称
    Device     string            // "cpu" / "cuda" / "cuda:0"
    Options    map[string]any    // 引擎特定选项
    OnProgress func(percent int) // 进度回调
}

type RecognizeResult struct {
    Subtitles []SubtitleEntry    // 字幕条目
    Language  string             // 检测到的语言
    Duration  float64            // 音频总时长 (秒)
}

type SubtitleEntry struct {
    Index   int     `json:"index"`
    Start   float64 `json:"start"`    // 开始时间 (秒)
    End     float64 `json:"end"`      // 结束时间 (秒)
    Text    string  `json:"text"`
    Speaker string  `json:"speaker"`  // 说话人标识 (可选)
}
```

### 6.2 翻译引擎接口

```go
// internal/translator/engine.go

type Engine interface {
    Translate(ctx context.Context, req TranslateRequest) (*TranslateResult, error)
    Name() string
    Type() EngineType
    Available() bool
    // 该引擎支持的语言列表
    SupportedLanguages() []string
}

type TranslateRequest struct {
    Entries    []SubtitleEntry  // 源语言字幕
    SourceLang string           // 源语言
    TargetLang string           // 目标语言
    OnProgress func(percent int)
}

type TranslateResult struct {
    Entries []SubtitleEntry     // 翻译后的字幕
}
```

### 6.3 TTS 引擎接口

```go
// internal/tts/engine.go

type Engine interface {
    Synthesize(ctx context.Context, req SynthesizeRequest) error
    Name() string
    Type() EngineType
    Available() bool
    // 获取可用音色列表
    ListVoices(lang string) ([]Voice, error)
    // 是否支持声音克隆
    SupportsClone() bool
}

type SynthesizeRequest struct {
    Entries    []TTSEntry       // 配音队列
    Language   string           // 目标语言
    OutputDir  string           // WAV 输出目录
    OnProgress func(percent int)
}

type TTSEntry struct {
    Index    int    `json:"index"`
    Text     string `json:"text"`
    Voice    string `json:"voice"`     // 音色名
    RefAudio string `json:"ref_audio"` // 克隆参考音频 (可选)
    Output   string `json:"output"`    // 输出 WAV 路径
}

type Voice struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Language string `json:"language"`
    Gender   string `json:"gender"`
    Preview  string `json:"preview"`   // 试听 URL
}
```

### 6.4 任务流水线接口

```go
// internal/pipeline/pipeline.go

type Pipeline struct {
    ID         string
    Config     *PipelineConfig
    Status     PipelineStatus       // pending / running / paused / done / error
    Stage      string               // 当前阶段
    Progress   int                  // 当前阶段进度 0-100
    OnProgress func(PipelineEvent)  // 进度事件回调 (推送到前端)
}

type PipelineConfig struct {
    // 输入
    InputFile   string
    SourceLang  string
    TargetLang  string

    // 引擎选择
    ASREngine       string  // "faster-whisper", "whisper-cpp", "openai-api", ...
    ASRModel        string
    TranslateEngine string  // "ollama", "openai", "google", ...
    TTSEngine       string  // "edge-tts", "qwen3-tts", "openai-tts", ...
    TTSVoice        string

    // 选项
    EnableDiarize    bool
    EnableDenoise    bool
    EnableSeparate   bool
    SubtitleType     string  // "none", "hard", "soft", "hard-dual", "soft-dual"
    VoiceCloneRef    string  // 克隆参考音频路径
    AudioAutoRate    bool    // 音频自动加速
    VideoAutoRate    bool    // 视频自动慢速

    // 输出
    OutputDir string
}

type PipelineEvent struct {
    PipelineID string `json:"pipeline_id"`
    Stage      string `json:"stage"`       // prepare/asr/translate/tts/align/assemble
    Progress   int    `json:"progress"`    // 0-100
    Message    string `json:"message"`
    Status     string `json:"status"`      // running/paused/done/error
}

// Wails 绑定方法 (前端可调用)
func (a *App) StartVTV(config PipelineConfig) (string, error)      // 返回 pipeline ID
func (a *App) StartSTT(config STTConfig) (string, error)
func (a *App) PausePipeline(id string) error
func (a *App) ResumePipeline(id string) error
func (a *App) CancelPipeline(id string) error
func (a *App) GetPipelineStatus(id string) (*PipelineEvent, error)

// 前端监听事件
// runtime.EventsOn("pipeline:progress", callback)
// runtime.EventsOn("pipeline:paused", callback)
// runtime.EventsOn("pipeline:done", callback)
```

### 6.5 CLI 命令结构

```go
// cli/main.go (cobra)

babellens
├── vtv          # 视频翻译
│   --input, --output, --source-lang, --target-lang
│   --asr-engine, --asr-model, --translate-engine
│   --tts-engine, --tts-voice, --subtitle-type
│   --diarize, --denoise, --separate
│   --config (TOML 配置文件)
│
├── stt          # 语音转录
│   --input, --output, --engine, --model, --lang
│   --format (srt/ass/vtt/txt)
│   --diarize, --denoise
│
├── sts          # 字幕翻译
│   --input, --output, --engine
│   --source-lang, --target-lang
│   --format (single/dual)
│
├── tts          # AI 配音
│   --input (SRT), --output (WAV)
│   --engine, --voice, --rate, --volume
│   --clone-ref (克隆参考音频)
│
├── separate     # 人声分离
│   --input, --output-dir
│
├── denoise      # 降噪
│   --input, --output, --engine (ffmpeg/frcrn)
│
├── engines      # 引擎管理
│   ├── list     # 列出所有可用引擎
│   └── test     # 测试引擎连通性
│
├── models       # 模型管理
│   ├── list     # 列出已下载模型
│   ├── download # 下载模型
│   └── remove   # 删除模型
│
├── tools        # 辅助工具
│   ├── extract-audio   # 从视频提取音频
│   ├── convert-sub     # 字幕格式转换
│   ├── merge-sub       # 字幕合并
│   ├── burn-sub        # 字幕烧录进视频
│   ├── mux             # 音视频混流
│   ├── watermark       # 视频水印
│   └── match-text      # 文字匹配
│
└── config       # 配置
    ├── show     # 显示当前配置
    └── edit     # 编辑配置
```

---

## 七、配置文件设计

```toml
# configs/default.toml

[app]
language = "zh"              # 界面语言
theme = "dark"               # 主题
temp_dir = ""                # 临时目录 (空=系统默认)
output_dir = "~/BabelLens"   # 默认输出目录

[python]
enabled = true               # 是否启用 Python Sidecar
worker_path = ""             # worker 路径 (空=自动检测)
port = 0                     # 端口 (0=自动分配)
idle_timeout = 300           # 空闲超时关闭 (秒)
startup_timeout = 30         # 启动超时 (秒)

[ffmpeg]
path = ""                    # ffmpeg 路径 (空=PATH 查找)
hw_accel = "auto"            # 硬件加速: auto/nvenc/qsv/amf/none
video_codec = "h264"         # 输出编码: h264/h265
crf = 23                     # 质量
preset = "medium"            # 速度

# ── ASR 默认配置 ──
[asr]
default_engine = "faster-whisper"
default_model = "large-v3"
default_device = "auto"      # auto/cpu/cuda

[asr.faster_whisper]
beam_size = 5
vad_filter = true

[asr.openai_api]
api_key = ""
base_url = "https://api.openai.com/v1"

[asr.custom]
url = ""
method = "POST"
headers = {}
body_template = ""           # Go template
response_path = ""           # JSONPath 到字幕数组

# ── 翻译默认配置 ──
[translator]
default_engine = "ollama"

[translator.ollama]
base_url = "http://127.0.0.1:11434"
model = "qwen2.5:7b"

[translator.openai]
api_key = ""
base_url = "https://api.openai.com/v1"
model = "gpt-4o"

[translator.deepl]
api_key = ""

[translator.custom]
url = ""
method = "POST"
headers = {}
body_template = ""
response_path = ""

# ── TTS 默认配置 ──
[tts]
default_engine = "edge-tts"
default_voice = ""           # 空=根据目标语言自动选择

[tts.edge_tts]
# 无需配置

[tts.openai]
api_key = ""
model = "tts-1"
voice = "alloy"

[tts.custom]
url = ""
method = "POST"
headers = {}
body_template = ""

# ── 处理选项默认值 ──
[processing]
enable_denoise = false
enable_diarize = false
enable_separate = false
audio_auto_rate = true
video_auto_rate = false
subtitle_type = "hard"
```

---

## 八、分阶段实施计划

### Phase 1: 骨架搭建

```
目标: Wails 项目初始化 + CLI 框架 + 核心接口定义

□ Wails 项目初始化 (React + TypeScript)
□ CLI 入口 (cobra) 搭建, 子命令注册
□ 定义 Go 核心接口 (asr.Engine, translator.Engine, tts.Engine)
□ 引擎注册表 (registry) 实现
□ FFmpeg 封装层 (runner, probe, audio 提取)
□ 字幕解析 (SRT 读写)
□ 配置管理 (TOML 加载)
□ 前端基础布局 (侧边栏导航 + 空白页面)
```

### Phase 2: STT 核心

```
目标: 语音转录功能完整可用

□ whisper.cpp 引擎集成 (Go 进程调用)
□ OpenAI Whisper API 引擎 (Go HTTP)
□ 自定义 API 引擎 (Go HTTP)
□ VAD 预切分 (sherpa-onnx Go 绑定)
□ Python Sidecar 管理器 (启动/停止/健康检查)
□ Python Worker 骨架 (FastAPI + faster-whisper 路由)
□ STT Pipeline 编排 (prepare → recogn → output)
□ CLI: babellens stt 完整实现
□ GUI: 语音转录页面
```

### Phase 3: 翻译 + 配音

```
目标: 字幕翻译 + AI 配音可用

□ Ollama 翻译引擎 (Go HTTP)
□ Google/DeepL 翻译引擎 (Go HTTP)
□ 自定义翻译 API
□ Edge-TTS 引擎 (Go HTTP)
□ sherpa-onnx Piper TTS (Go 绑定)
□ 自定义 TTS API
□ STS Pipeline (translate → output)
□ TTS Pipeline (synthesize → align → output)
□ CLI: babellens sts, babellens tts
□ GUI: 字幕翻译页, 配音页
```

### Phase 4: VTV 全流程

```
目标: 视频翻译完整流水线

□ 音视频对齐 (rubberband 音频加速 + FFmpeg 视频慢速)
□ FFmpeg 字幕嵌入 (硬/软/双语)
□ FFmpeg 音视频合并 + 硬件加速编码
□ VTV Pipeline 全流程串联
□ 暂停/恢复/取消机制
□ CLI: babellens vtv
□ GUI: 视频翻译页 (完整流程 + 进度)
□ GUI: 字幕编辑器组件
```

### Phase 5: 增强功能

```
目标: 人声分离、降噪、说话人识别、声音克隆

□ 人声分离 (sherpa-onnx Go 绑定)
□ 说话人分离 (sherpa-onnx Go 绑定)
□ 降噪 FFmpeg 方案
□ Python: FRCRN 降噪路由
□ Python: pyannote 说话人分离路由
□ Python: Qwen3-TTS / F5-TTS 路由 (声音克隆)
□ GUI: 工具页 (分离/降噪)
□ GUI: 配音页 (多角色 + 克隆)
□ CLI: babellens separate, babellens denoise
```

### Phase 6: 打包发布

```
目标: 可分发的桌面应用

□ Wails 构建 Windows/macOS/Linux
□ Python Worker PyInstaller 打包
□ 模型下载管理器 (首次使用时下载)
□ 自动检测 FFmpeg/rubberband/ollama
□ GPU 自动检测和配置
□ 设置页面 (引擎配置/路径配置)
□ 错误处理和用户友好提示
□ 构建脚本 + CI/CD
```

---

## 九、技术选型汇总

| 领域 | 选型 | 理由 |
|------|------|------|
| 主语言 | Go 1.24 | 并发、编译速度、单二进制 |
| AI 推理 | Python 3.12 (Sidecar) | PyTorch 生态不可替代 |
| 桌面框架 | Wails v2 | 轻量、Go 原生、体积小 |
| 前端 | React + TypeScript + Vite | 生态丰富、类型安全 |
| CLI | cobra | Go 标准 CLI 框架 |
| 配置 | TOML (BurntSushi/toml) | 人类可读、Go 原生支持好 |
| 日志 | zerolog | 结构化、零分配、高性能 |
| HTTP 客户端 | net/http (标准库) | 无需第三方依赖 |
| Python 框架 | FastAPI + uvicorn | 异步、类型提示、自动文档 |
| Python 包管理 | uv | 极快、兼容 pip |
| Python 打包 | PyInstaller | 独立可执行文件 |
| 音视频 | FFmpeg (外部二进制) | 工业标准 |
| 本地 ASR | whisper.cpp + sherpa-onnx | 无需 Python 依赖 |
| 本地翻译 | Ollama | Go 实现、HTTP API、模型管理完善 |
| 本地 TTS | sherpa-onnx Piper | Go 绑定、离线可用 |
| 人声分离 | sherpa-onnx UVR | Go 绑定、离线可用 |
| 说话人分离 | sherpa-onnx | Go 绑定、离线可用 |
| 音频变速 | rubberband (外部二进制) | 变速不变调、工业标准 |
