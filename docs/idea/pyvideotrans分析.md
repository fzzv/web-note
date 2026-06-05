# pyVideoTrans 项目功能分析与技术架构

> 项目版本：v3.98 | 开源协议：GPL-V3 | Python 3.10+
> 仓库：https://github.com/jianchang512/pyvideotrans

---

## 一、项目概述

pyVideoTrans 是一款**开源视频翻译 / 语音转录 / AI 配音 / 字幕翻译**桌面工具，提供完整的工作流：**语音识别 (ASR) → 字幕翻译 → 语音合成 (TTS) → 音视频合成**。支持本地离线部署和多种在线 API，覆盖 30+ 种语言。

项目提供两种运行模式：
- **GUI 模式**（`sp.py`）：基于 PySide6 (Qt6) 的桌面图形界面
- **CLI 模式**（`cli.py`）：命令行界面，适合服务器和批量处理

---

## 二、项目目录结构

```
pyvideotrans/
├── sp.py                    # GUI 入口（PySide6 启动窗口 + 主窗口）
├── cli.py                   # CLI 入口（argparse 命令行解析）
├── pyproject.toml           # 项目依赖配置（uv 管理）
├── models/                  # 本地模型文件存放
├── ffmpeg/                  # 内置 FFmpeg 二进制
├── videotrans/              # 核心业务包
│   ├── configure/           # 全局配置、常量、异常处理
│   ├── recognition/         # 语音识别（ASR）引擎适配层（22 种渠道）
│   ├── translator/          # 翻译引擎适配层（24 种渠道）
│   ├── tts/                 # 语音合成（TTS）引擎适配层（32 种渠道）
│   ├── task/                # 任务编排层（核心业务流程）
│   │   ├── _base.py         # 任务基类 BaseTask
│   │   ├── trans_create.py  # 视频翻译任务 (VTV)
│   │   ├── _speech2text.py  # 语音转录任务 (STT)
│   │   ├── _translate_srt.py# 字幕翻译任务 (STS)
│   │   ├── _dubbing.py      # 配音任务 (TTS)
│   │   ├── _rate.py         # 音视频对齐（加速/慢速）
│   │   ├── vad.py           # VAD 语音活动检测
│   │   └── separate_worker.py # 人声/背景音分离
│   ├── process/             # 子进程任务（降噪、ASR、TTS 工作函数）
│   ├── mainwin/             # 主窗口逻辑、信号绑定、事件处理
│   ├── component/           # UI 组件（进度条、字幕编辑、配音设置等）
│   ├── ui/                  # UI 界面定义（各个设置窗口）
│   ├── winform/             # 各设置窗口的业务逻辑
│   ├── util/                # 工具函数（FFmpeg 封装、SRT 处理、GPU 检测等）
│   ├── language/            # 多语言翻译文件
│   ├── styles/              # QSS 样式表、图标资源
│   ├── prompts/             # AI 翻译 prompt 模板
│   └── codes/               # 模型相关配置
```

---

## 三、核心功能详解

### 3.1 视频翻译（VTV — Video to Video）

**功能**：将 A 语言发音的视频，自动翻译为 B 语言配音 + B 语言字幕的视频。

**处理流程**（`trans_create.py` → `TransCreate` 类）：

```
输入视频 (MP4/MOV/AVI...)
    │
    ▼
① prepare() — 预处理
   ├── FFmpeg 提取音频 → 16kHz 单声道 WAV
   ├── FFmpeg 分离无声视频
   ├── [可选] 人声/背景音乐分离 (sherpa-onnx UVR 模型)
   └── 获取视频元信息（分辨率、帧率、编码、时长）
    │
    ▼
② recogn() — 语音识别
   ├── 调用 recognition.run()，按 recogn_type 分发到对应 ASR 引擎
   ├── 输出：源语言 SRT 字幕（含时间戳）
   └── [可选] 二次识别 (recogn2pass)
    │
    ▼
③ trans() — 字幕翻译
   ├── 调用 translator.run()，按 translate_type 分发到对应翻译引擎
   └── 输出：目标语言 SRT 字幕
    │
    ▼
④ dubbing() — AI 配音
   ├── 解析目标字幕，为每条字幕生成配音任务 (queue_tts)
   ├── 调用 tts.run()，按 tts_type 分发到对应 TTS 引擎
   └── 输出：每条字幕对应的 WAV 音频片段
    │
    ▼
⑤ align() — 音视频对齐 (_rate.py → TtsSpeedRate)
   ├── 音频加速：配音时长 > 字幕时长时，用 rubberband 加速音频
   ├── 视频慢速：用 FFmpeg setpts 滤镜延长视频片段
   ├── 拼接所有音频片段 + 静音填充
   └── 确保音频总长 = 视频总长
    │
    ▼
⑥ assembling() — 合成输出
   ├── FFmpeg 合并：无声视频 + 配音音频 + 背景音乐
   ├── 字幕嵌入：硬字幕(烧录) / 软字幕(内封) / 双语字幕
   └── 编码输出：H.264 / H.265，支持硬件加速 (NVENC/QSV/AMF/VideoToolbox)
    │
    ▼
⑦ task_done() — 清理临时文件，输出最终视频
```

**关键实现细节**：
- 每个阶段可暂停，允许用户手动校对字幕后继续
- 支持说话人识别（`pyannote-audio`），为不同说话人分配不同配音角色
- 支持批量视频处理（`_mult_video.py`）

---

### 3.2 语音转录 / 字幕生成（STT — Speech to Text）

**功能**：从音频/视频中提取语音，生成带时间戳的字幕文件。

**实现**（`_speech2text.py` → `SpeechToText` 类）：

```
prepare() → 转换为 16kHz WAV
    ↓
[可选] 降噪 → ModelScope FRCRN 降噪模型 (iic/speech_frcrn_ans_cirm_16k)
    ↓
recogn() → 调用 ASR 引擎识别
    ↓
[可选] 说话人分离 (diariz) → pyannote-audio
    ↓
输出 SRT / TXT / ASS / VTT 字幕文件
```

**支持的 22 种 ASR 引擎**：

| 类型 | 引擎 | 实现文件 | 说明 |
|------|------|----------|------|
| 本地 | Faster-Whisper | `_overall.py` | 默认引擎，ctranslate2 加速，速度快 |
| 本地 | OpenAI Whisper | `_overall.py` | 官方 openai-whisper 包 |
| 本地 | FunASR | `_funasr.py` | 阿里达摩院，中文优化 |
| 本地 | Qwen-ASR | `_qwenasrlocal.py` | 通义千问本地 ASR |
| 本地 | Huggingface ASR | `_huggingface.py` | 支持多个 HF 模型（Parakeet 等）|
| 本地 | WhisperX | `_whisperx.py` | 支持时间戳对齐和说话人分离 |
| 本地 | Whisper.cpp | 通过外部进程 | C++ 实现，低资源占用 |
| 本地 | Whisper.NET | `_whispernet.py` | .NET 实现 |
| 在线 | OpenAI API | `_openairecognapi.py` | OpenAI 官方 Whisper API |
| 在线 | Gemini AI | `_gemini.py` | Google Gemini 多模态 |
| 在线 | 阿里 Qwen3-ASR | `_qwen3asr.py` | 阿里百炼 API |
| 在线 | 字节火山引擎 STT | `_zijiemodel.py` | 字节跳动 ASR |
| 在线 | 字节火山字幕 API | `_doubao.py` | 豆包字幕 API |
| 在线 | 智谱 GLM-ASR | `_glmasr.py` | 智谱 AI |
| 在线 | Deepgram | `_deepgram.py` | Deepgram.com |
| 在线 | Google Speech | `_google.py` | Google Cloud STT |
| 在线 | ElevenLabs | `_elevenlabs.py` | ElevenLabs.io |
| 在线 | 302.AI | `_ai302.py` | 302.AI 聚合 API |
| 外部 | Parakeet-tdt | `_parakeet.py` | NVIDIA Parakeet |
| 外部 | Faster-Whisper-XXL | 外部进程 | XXL 增强版 |
| 自定义 | STT API | `_stt.py` | 用户自定义 API |
| 自定义 | Custom API | `_recognapi.py` | 自定义识别 API |

**VAD 语音活动检测**（`vad.py`）：
- 使用 `silero-VAD`（通过 faster-whisper 内置）和 `TenVad` 两种方案
- 先检测语音片段再送入 ASR，提高效率和准确率
- 支持配置：阈值、最短语音时长、最长语音时长、最短静默时长

---

### 3.3 字幕翻译（STS — Subtitle to Subtitle）

**功能**：将 SRT 字幕文件从一种语言翻译为另一种语言。

**实现**（`_translate_srt.py` → `TranslateSrt` 类）：

```
读取源 SRT → 解析字幕条目列表
    ↓
调用 translator.run() 批量翻译
    ↓
输出目标语言 SRT（支持单语/双语格式）
```

**支持的 24 种翻译引擎**：

| 分类 | 引擎 | 说明 |
|------|------|------|
| **AI/LLM** | ChatGPT, DeepSeek, Gemini, 智谱AI, Azure GPT, Local LLM (Ollama), OpenRouter, 硅基流动, 302.AI, 火山引擎 LLM, 阿里百炼, MiniMax | 上下文理解好，翻译自然 |
| **传统机翻** | Google, Microsoft, 百度, 腾讯, DeepL, DeepLx, 阿里机翻 | 速度快，免费/低成本 |
| **离线** | M2M100, LibreTranslate | 完全本地，无需网络 |
| **自定义** | OTT, MyMemory, 自定义 API | 对接自有翻译服务 |

**语言代码映射系统**（`translator/__init__.py`）：
- 每种语言维护一个 11 元素数组，分别对应不同翻译渠道所需的语言代码格式
- 索引：0=Google, 1=字幕嵌入, 2=百度, 3=DeepL, 4=腾讯, 5=OTT, 6=微软, 7=AI渠道, 8=阿里, 9=Qwen, 10=M2M100
- 支持 30 种语言 + 自定义语言扩展（`newlang.txt`）

---

### 3.4 AI 配音（TTS — Text to Speech）

**功能**：将字幕文本合成为语音，支持多角色配音和声音克隆。

**实现**（`_dubbing.py` → `DubbingSrt` 类）：

```
读取 SRT 字幕 → 按条目构建 queue_tts
    ↓
调用 tts.run() → 为每条字幕生成 WAV 音频
    ↓
align() → 对齐音频时长与字幕时长
    ↓
输出配音 WAV 文件
```

**支持的 32 种 TTS 引擎**：

| 分类 | 引擎 | 特点 |
|------|------|------|
| **免费** | Edge-TTS, gTTS, Azure(free) | 微软/Google 免费接口 |
| **本地 (声音克隆)** | F5-TTS, CosyVoice, GPT-SoVITS, ChatterBox, Index TTS, Spark TTS, VoxCPM, Dia TTS, Qwen3-TTS | 零样本声音克隆，自部署 |
| **本地 (普通)** | ChatTTS, Kokoro, Fish TTS, VITS, Piper, clone-voice, Supertonic | 本地部署 TTS |
| **在线** | OpenAI TTS, Azure TTS, Minimaxi, Qwen3-TTS(百炼), 火山引擎 TTS, 豆包2, 智谱 GLM-TTS, ElevenLabs, Gemini TTS, 302.AI, X.AI TTS | 商业 API |
| **自定义** | TTS API, Google Cloud TTS | 自定义/Google Cloud |

**声音克隆**：支持 10 种克隆渠道（`SUPPORT_CLONE` 列表），用户上传参考音频即可克隆说话风格。

---

### 3.5 音视频对齐（Audio-Video Alignment）

**功能**：确保配音时长与原始视频时间轴同步。

**实现**（`_rate.py` → `TtsSpeedRate` 类）：

**对齐策略**：

| 场景 | 处理方式 |
|------|----------|
| 配音时长 ≤ 字幕时长 | 无需处理，直接使用 |
| 配音时长 > 字幕时长，差距 ≤ 20% | 音频加速（rubberband/pyrubberband）|
| 配音时长 > 字幕时长，差距 > 20% | 音频加速 + 视频慢速各承担一半 |
| 仅音频加速模式 | rubberband 变速不变调，有最大倍速限制 |
| 仅视频慢速模式 | FFmpeg `setpts=X*PTS` 滤镜，有最大 PTS 倍率限制 |

**技术细节**：
- 音频加速使用 `rubberband` 库（变速不变调）
- 视频慢速使用 FFmpeg 的 `setpts` 滤镜 + `-fps_mode vfr`
- 字幕间静音区间会被合并利用以减少加速/慢速幅度

---

### 3.6 人声/背景音分离

**功能**：将视频/音频中的人声和背景音乐分离。

**实现**（`separate_worker.py` + `process/prepare_audio.py`）：

- 使用 **sherpa-onnx** 的 UVR-MDX-NET 模型（`UVR-MDX-NET-Inst_HQ_4.onnx`）
- 仅 CPU 推理，无需 GPU
- 输出：`vocal-xxx.wav`（人声）+ `instrument-xxx.wav`（背景音乐）
- 用途：翻译时保留原始背景音乐，仅替换人声

---

### 3.7 降噪处理

**功能**：对语音进行预处理降噪，提升 ASR 准确率。

**实现**（`process/prepare_audio.py` → `remove_noise()`）：

- 使用 ModelScope 模型 `iic/speech_frcrn_ans_cirm_16k`
- 支持 CUDA 加速
- 首次使用自动下载模型

---

### 3.8 辅助工具集

项目在 `winform/` 中提供了多个独立工具：

| 工具 | 文件 | 功能 |
|------|------|------|
| 从视频提取音频 | `fn_audiofromvideo.py` | FFmpeg 提取音频轨道 |
| 字幕翻译 | `fn_fanyisrt.py` | 独立字幕翻译面板 |
| 字幕格式转换 | `fn_formatcover.py` | SRT/ASS/VTT 互转 |
| 字幕合并 | `fn_hebingsrt.py` | 合并多个 SRT 文件 |
| 音视频混流 | `fn_hunliu.py` | 将音频混入视频 |
| 字幕配音 | `fn_peiyin.py` | 独立配音面板 |
| 多角色配音 | `fn_peiyinrole.py` | 为不同说话人分配不同音色 |
| 语音识别 | `fn_recogn.py` | 独立语音识别面板 |
| 人声分离 | `fn_separate.py` | 人声/背景音分离面板 |
| 字幕覆盖 | `fn_subtitlescover.py` | 字幕嵌入视频 |
| 文字匹配 | `fn_vas.py` | 文本与字幕时间轴匹配 |
| 视频+音频合并 | `fn_videoandaudio.py` | 合并视频和音频轨道 |
| 视频+字幕合并 | `fn_videoandsrt.py` | 将字幕嵌入视频 |
| 水印 | `fn_watermark.py` | 为视频添加水印 |

---

## 四、技术架构

### 4.1 整体架构图

```
┌────────────────────────────────────────────────────────────────────┐
│                         用户界面层                                  │
│  ┌──────────────────────┐     ┌──────────────────────┐            │
│  │  GUI (PySide6/Qt6)   │     │  CLI (argparse)      │            │
│  │  sp.py → MainWindow  │     │  cli.py              │            │
│  │  - ui/ (界面定义)     │     │  - 4种任务分发        │            │
│  │  - winform/ (逻辑)   │     │                      │            │
│  │  - component/ (组件)  │     │                      │            │
│  └──────────┬───────────┘     └──────────┬───────────┘            │
└─────────────┼────────────────────────────┼────────────────────────┘
              │                            │
              ▼                            ▼
┌────────────────────────────────────────────────────────────────────┐
│                         任务编排层 (task/)                          │
│                                                                    │
│  BaseTask  ← 所有任务基类                                           │
│    ├── TransCreate (VTV)   视频翻译：prepare→recogn→trans→          │
│    │                       dubbing→align→assembling                │
│    ├── SpeechToText (STT)  语音转录：prepare→recogn→diariz         │
│    ├── TranslateSrt (STS)  字幕翻译：trans                          │
│    └── DubbingSrt (TTS)    字幕配音：dubbing→align                  │
│                                                                    │
│  TaskCfg 数据类：TaskCfgSTT, TaskCfgTTS, TaskCfgSTS, TaskCfgVTT    │
└──────────┬──────────────────┬──────────────────┬──────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
│  recognition/    │ │  translator/     │ │  tts/            │
│  ASR 引擎适配层   │ │  翻译引擎适配层   │ │  TTS 引擎适配层   │
│                  │ │                  │ │                  │
│  22 种 ASR 渠道   │ │  24 种翻译渠道   │ │  32 种 TTS 渠道   │
│  统一接口 run()  │ │  统一接口 run()  │ │  统一接口 run()  │
│                  │ │                  │ │                  │
│  _base.py 基类   │ │  _base.py 基类   │ │  _base.py 基类   │
│  __init__.py分发 │ │  __init__.py分发 │ │  __init__.py分发 │
└──────────────────┘ └──────────────────┘ └──────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌────────────────────────────────────────────────────────────────────┐
│                         基础能力层                                  │
│                                                                    │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐       │
│  │ FFmpeg       │  │ VAD         │  │ 人声分离              │       │
│  │ (util/       │  │ (silero-vad │  │ (sherpa-onnx UVR)    │       │
│  │ help_ffmpeg) │  │  + TenVad)  │  │                      │       │
│  └─────────────┘  └─────────────┘  └──────────────────────┘       │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────┐       │
│  │ rubberband   │  │ 降噪处理     │  │ pyannote-audio       │       │
│  │ (音频变速)    │  │ (ModelScope) │  │ (说话人分离)          │       │
│  └─────────────┘  └─────────────┘  └──────────────────────┘       │
└────────────────────────────────────────────────────────────────────┘
```

### 4.2 设计模式

**策略模式 (Strategy Pattern)**：ASR / 翻译 / TTS 三大模块均采用统一接口 + 工厂分发：
- 每个渠道实现一个类，继承自 `_base.py` 基类
- `__init__.py` 中的 `run()` 函数根据 `type` 参数分发到具体实现
- 新增渠道只需：新建实现文件 → 在 `__init__.py` 注册

**模板方法模式 (Template Method)**：`BaseTask` 定义了标准处理流水线：
```python
prepare() → recogn() → diariz() → trans() → dubbing() → align() → assembling() → task_done()
```
子类按需重写各步骤。

**信号-槽机制**：GUI 模式使用 Qt 的信号机制（`_signal.py`）在主线程和工作线程间传递进度信息。

---

## 五、核心技术栈

### 5.1 AI / ML 相关

| 依赖 | 用途 |
|------|------|
| `faster-whisper` (ctranslate2) | 默认 ASR 引擎，Whisper 加速推理 |
| `openai-whisper` | OpenAI 官方 Whisper 实现 |
| `pyannote-audio` | 说话人分离 (Speaker Diarization) |
| `funasr` | 阿里 FunASR 语音识别 |
| `modelscope` | 模型下载（降噪模型等）|
| `sherpa-onnx` | 人声分离 (UVR-MDX-NET)、ONNX 推理 |
| `transformers` | HuggingFace 模型加载 |
| `torch` / `torchaudio` | PyTorch 深度学习框架 |
| `edge-tts` | 微软免费 TTS |
| `ten-vad` | 语音活动检测 |
| `sentencepiece` / `tiktoken` | 分词器 |

### 5.2 音视频处理

| 依赖 | 用途 |
|------|------|
| `ffmpeg` (外部二进制) | 音视频转码、合并、字幕嵌入、裁切 |
| `pydub` | 音频处理（时长检测、格式转换）|
| `librosa` | 音频分析 |
| `soundfile` / `pysoundfile` | 音频文件 I/O |
| `pyrubberband` | 音频时间拉伸（变速不变调）|
| `pytsmod` | 时间尺度修改 |
| `av` (PyAV) | FFmpeg Python 绑定 |
| `srt` | SRT 字幕解析 |

### 5.3 GUI 框架

| 依赖 | 用途 |
|------|------|
| `PySide6` | Qt6 桌面 GUI 框架 |
| `qdarkstyle` | 暗色主题样式 |

### 5.4 API 客户端

| 依赖 | 用途 |
|------|------|
| `openai` | OpenAI API 客户端 |
| `anthropic` | Claude API 客户端 |
| `google-genai` | Google Gemini API |
| `deepl` | DeepL 翻译 API |
| `deepgram-sdk` | Deepgram ASR API |
| `elevenlabs` | ElevenLabs TTS API |
| `dashscope` | 阿里灵积 API |
| `azure-cognitiveservices-speech` | Azure 语音服务 |
| `tencentcloud-sdk-python` | 腾讯云 API |
| `alibabacloud-*` | 阿里云 API 系列 |

---

## 六、关键流程的实现细节

### 6.1 语音识别流程（以 Faster-Whisper 为例）

```python
# recognition/_overall.py → FasterAll 类

1. 下载模型 → HuggingFace Hub 下载到 models/ 目录
2. [可选] VAD 预切分 → silero-VAD 检测语音片段
3. 启动子进程 → process/faster_whisper.py
   - 加载 CTranslate2 格式的 Whisper 模型
   - 配置参数：语言、CUDA、初始 prompt、温度等
   - 逐段识别，输出带 word-level 时间戳的结果
4. 后处理 → 合并相邻片段、修正时间戳、标点恢复
5. 返回 List[Dict]，每个 dict 含 start_time, end_time, text
```

### 6.2 视频编码输出

`util/help_ffmpeg.py` 实现了智能编码选择：

```
检测系统 GPU → 选择最优编码器
  ├── NVIDIA GPU → h264_nvenc / hevc_nvenc (NVENC)
  ├── Intel GPU → h264_qsv / hevc_qsv (QSV)
  ├── AMD GPU → h264_amf / hevc_amf (AMF)
  ├── macOS → h264_videotoolbox / hevc_videotoolbox
  └── 无 GPU → libx264 / libx265 (软编码)

CRF 值自动映射 → 不同编码器的质量参数格式
Preset 分类 → fast / medium / slow 三档
```

### 6.3 字幕嵌入方式

| 类型 | subtitle_type | 实现方式 |
|------|--------------|----------|
| 无字幕 | 0 | 不处理 |
| 硬字幕 | 1 | FFmpeg `subtitles` 滤镜烧录 |
| 软字幕 | 2 | FFmpeg `-c:s mov_text` 内封 |
| 硬双语 | 3 | 双语 ASS 渲染后烧录 |
| 软双语 | 4 | 双语 SRT 内封 |

---

## 七、项目特点总结

### 优势
1. **渠道覆盖极广**：ASR 22 种 + 翻译 24 种 + TTS 32 种，堪称市面最全
2. **支持完全离线**：Faster-Whisper + M2M100/Ollama + 本地 TTS 可实现零网络依赖
3. **声音克隆**：集成 F5-TTS、CosyVoice、GPT-SoVITS 等主流克隆方案
4. **音视频对齐**：自研的音频加速 + 视频慢速双重对齐策略，处理跨语言时长差异
5. **人工校对支持**：每个阶段可暂停编辑，保证翻译质量
6. **跨平台**：Windows / macOS / Linux 均支持
7. **CLI + GUI 双模式**：适应桌面使用和服务器批处理

### 架构特点
- 基于**策略模式**的插件化引擎适配，扩展性好
- 大量使用**子进程**执行 AI 推理，避免 GIL 瓶颈和内存泄漏
- FFmpeg 作为底层音视频处理核心，稳定可靠
- 多语言代码映射系统，一套代码兼容所有翻译渠道的语言标识差异
