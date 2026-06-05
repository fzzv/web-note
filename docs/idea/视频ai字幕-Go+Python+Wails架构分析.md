# Go + Python + Wails 完全离线桌面应用可行性分析

---

## 一、直接结论

**完全可行，且是这类项目的最优架构之一。** 原项目 pyVideoTrans 本身已经在用这种思路——Python AI 推理全部跑在**子进程**中（`ProcessPoolExecutor`），主进程只做调度。Go 天然就是做这个「调度层」的最佳语言。

---

## 二、为什么说可行：原项目已经验证了这种分离

通过源码分析，pyVideoTrans 的 Python AI 推理全部通过 `GlobalProcessManager`（基于 `multiprocessing.ProcessPoolExecutor`）在**独立子进程**中执行：

```
当前 pyVideoTrans 的实际架构：

PySide6 主进程 (UI + 调度)
    │
    ├── subprocess → openai_whisper()      # Whisper ASR
    ├── subprocess → faster_whisper()      # Faster-Whisper ASR
    ├── subprocess → funasr_mlt()          # FunASR
    ├── subprocess → qwen3tts_fun()        # Qwen3 TTS
    ├── subprocess → vocal_bgm()           # 人声分离
    ├── subprocess → remove_noise()        # 降噪
    ├── subprocess → pyannote_speakers()   # 说话人分离
    ├── subprocess → cam_speakers()        # 阿里 CAM++ 说话人分离
    ├── subprocess → built_speakers()      # sherpa-onnx 说话人分离
    └── subprocess → ffmpeg (外部二进制)    # 音视频处理
```

**核心发现：Python 主进程和 AI 推理进程之间的通信方式极其简单**——仅通过**文件**（JSON 参数文件 + 日志文件 + 音频/字幕结果文件）和**进程返回值**（`(result, error)` 元组）交互。

这意味着：把「Python 主进程」替换为「Go 主进程」，通信方式完全不需要改变。

---

## 三、架构设计：Go 主控 + Python Sidecar

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Wails 桌面应用                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  前端 (Vue3/React + TypeScript)                        │  │
│  │  视频上传 │ 字幕编辑器 │ 视频预览 │ 设置面板 │ 进度条  │  │
│  └───────────────────────┬───────────────────────────────┘  │
│                          │ Wails Bindings (直接调用 Go 函数) │
│  ┌───────────────────────▼───────────────────────────────┐  │
│  │                 Go 后端 (主控层)                        │  │
│  │                                                        │  │
│  │  ┌──────────────┐ ┌──────────────┐ ┌───────────────┐  │  │
│  │  │ 任务编排引擎  │ │ FFmpeg 封装   │ │ 字幕处理      │  │  │
│  │  │ (goroutine)  │ │ (os/exec)    │ │ SRT/ASS/VTT  │  │  │
│  │  └──────┬───────┘ └──────────────┘ └───────────────┘  │  │
│  │         │                                              │  │
│  │  ┌──────▼────────────────────────────────────────────┐│  │
│  │  │         Python Sidecar 管理器                      ││  │
│  │  │  - 启动/停止 Python 工作进程                        ││  │
│  │  │  - 任务分发 (JSON-RPC / stdio / HTTP)              ││  │
│  │  │  - 超时控制 + 健康检查                              ││  │
│  │  │  - GPU 资源调度                                    ││  │
│  │  └──────┬────────────────────────────────────────────┘│  │
│  └─────────┼────────────────────────────────────────────┘  │
└────────────┼────────────────────────────────────────────────┘
             │
             ▼
┌────────────────────────────────────────────────────────────┐
│              Python Worker 进程 (Sidecar)                    │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Flask/FastAPI HTTP 微服务  或  stdio JSON-RPC        │  │
│  │                                                       │  │
│  │  路由/命令:                                            │  │
│  │  POST /asr          → faster_whisper / openai_whisper │  │
│  │  POST /tts          → qwen3_tts / cosyvoice / f5tts   │  │
│  │  POST /translate    → m2m100 本地翻译                  │  │
│  │  POST /separate     → sherpa_onnx 人声分离             │  │
│  │  POST /denoise      → modelscope 降噪                 │  │
│  │  POST /diarize      → pyannote 说话人分离              │  │
│  │  POST /punctuate    → modelscope 标点恢复              │  │
│  │  GET  /health       → 健康检查                         │  │
│  │  GET  /gpu-info     → GPU 状态                         │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  本地模型 (models/ 目录)                               │  │
│  │  - faster-whisper (large-v3)     ~3GB                 │  │
│  │  - Qwen3-TTS                     ~3-7GB               │  │
│  │  - pyannote/speaker-diarization  ~200MB               │  │
│  │  - UVR-MDX-NET (人声分离)         ~100MB               │  │
│  │  - M2M100 (离线翻译)             ~2GB                  │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
             │
             ▼
┌──────────────────────────────────────────────────────────────┐
│  外部二进制                                                    │
│  ├── ffmpeg / ffprobe    (Go 直接调用)                        │
│  └── rubberband          (Go 直接调用，可选)                   │
└──────────────────────────────────────────────────────────────┘
```

### 3.2 Go ↔ Python 通信方案对比

| 方案 | 原理 | 延迟 | 复杂度 | 推荐场景 |
|------|------|------|--------|----------|
| **HTTP 微服务** | Python 启动 Flask/FastAPI，Go 用 `net/http` 调用 | ~1-5ms | ⭐ 低 | **推荐首选** |
| **stdio JSON-RPC** | Go 启动 Python 子进程，通过 stdin/stdout 传 JSON | <1ms | ⭐⭐ 中 | 不想开端口时 |
| **gRPC** | Protocol Buffers 序列化，双向流 | <1ms | ⭐⭐⭐ 高 | 需要流式传输进度时 |
| **文件交换** | 参数写 JSON 文件，结果写文件，Go 轮询 | ~10-50ms | ⭐ 低 | 原项目就是这样做的 |
| **Unix Socket** | 本地 socket 通信 | <1ms | ⭐⭐ 中 | 高性能场景 |

**推荐方案：HTTP 微服务（FastAPI）**

理由：
1. 调试最方便——可以直接 curl 测试
2. Python 进程独立，崩溃不影响 Go 主进程
3. 天然支持异步——长时间的 ASR 任务不阻塞
4. 进度推送可以用 SSE (Server-Sent Events)
5. Go 的 `net/http` 客户端极其成熟

```
Go 端:                              Python 端:

app.POST("/start-asr", handler)     @app.post("/asr")
  │                                 async def asr(req: ASRRequest):
  │  POST http://127.0.0.1:9876/asr     model = load_whisper(req.model)
  ├──────────────────────────────►      result = model.transcribe(req.audio)
  │                                     return {"subtitles": result}
  │  {"subtitles": [...]}
  ◄──────────────────────────────┤
  │
  ▼ 继续流水线：翻译 → 配音 → 合成
```

---

## 四、每个离线功能的具体实现方案

### 4.1 语音识别（ASR）—— 完全离线

| 方案 | 实现方式 | 质量 | 是否需要 Python |
|------|----------|------|----------------|
| **A. faster-whisper (推荐)** | Python sidecar 调用 faster-whisper | ⭐⭐⭐⭐⭐ 最佳 | 是 |
| **B. whisper.cpp** | Go 通过 CGo 绑定或 os/exec 调用 | ⭐⭐⭐⭐ 优秀 | **否** |
| **C. sherpa-onnx** | Go 原生绑定，ONNX 推理 | ⭐⭐⭐⭐ 优秀 | **否** |

**分析**：whisper.cpp 和 sherpa-onnx 都可以让 ASR **完全不依赖 Python**。质量上 faster-whisper 略好（支持的参数更多），但差距不大。如果极致追求去 Python，ASR 环节可以用纯 Go 方案。

### 4.2 字幕翻译 —— 完全离线

| 方案 | 实现方式 | 质量 | 是否需要 Python |
|------|----------|------|----------------|
| **A. Ollama 本地 LLM** | Go HTTP 调用本地 Ollama API | ⭐⭐⭐⭐⭐ 最佳 | **否** |
| **B. llama.cpp** | Go CGo 绑定 | ⭐⭐⭐⭐ 优秀 | **否** |
| **C. M2M100** | Python sidecar 调用 ctranslate2 | ⭐⭐⭐ 良好 | 是 |

**分析**：翻译环节 **完全不需要 Python**。Ollama 是独立的 Go 程序，提供 HTTP API，可以运行 Qwen2.5、DeepSeek 等模型做高质量翻译。这是最优方案。

### 4.3 语音合成（TTS）—— 离线

| 方案 | 实现方式 | 质量 | 克隆能力 | 是否需要 Python |
|------|----------|------|----------|----------------|
| **A. sherpa-onnx TTS** | Go 原生绑定 (Piper/VITS) | ⭐⭐⭐ 良好 | 无 | **否** |
| **B. Qwen3-TTS** | Python sidecar | ⭐⭐⭐⭐⭐ 最佳 | ✅ 有 | 是 |
| **C. F5-TTS / CosyVoice** | Python sidecar (Gradio API) | ⭐⭐⭐⭐⭐ 最佳 | ✅ 有 | 是 |
| **D. Kokoro / ChatTTS** | Python sidecar | ⭐⭐⭐⭐ 优秀 | 部分 | 是 |

**分析**：TTS 是**最需要 Python 的环节**。纯 Go 方案 (sherpa-onnx Piper) 质量能用但远不如 Qwen3-TTS/F5-TTS。**如果需要声音克隆，Python 不可替代。**

### 4.4 人声分离 —— 完全离线

| 方案 | 实现方式 | 是否需要 Python |
|------|----------|----------------|
| **sherpa-onnx** | Go 原生绑定，UVR-MDX-NET ONNX | **否** ✅ |

原项目已经用 sherpa-onnx 做人声分离（仅 CPU），sherpa-onnx 有官方 Go 绑定，**完全不需要 Python**。

### 4.5 说话人分离 —— 完全离线

| 方案 | 实现方式 | 质量 | 是否需要 Python |
|------|----------|------|----------------|
| **A. sherpa-onnx** | Go 原生绑定 | ⭐⭐⭐⭐ 优秀 | **否** ✅ |
| **B. pyannote-audio** | Python sidecar | ⭐⭐⭐⭐⭐ 最佳 | 是 |

原项目中 `built_speakers()` 已经用 sherpa-onnx 做说话人分离，与 pyannote 质量接近，**可以不依赖 Python**。

### 4.6 降噪 / 标点恢复

| 功能 | 纯 Go 方案 | Python 方案 |
|------|-----------|-------------|
| 降噪 | FFmpeg `afftdn`/`arnndn` 滤镜（效果一般）| ModelScope FRCRN（效果好）|
| 标点恢复 | LLM 后处理（Ollama）| ModelScope CT-Transformer |

降噪可以用 FFmpeg 内置滤镜做基础版，标点恢复可以让本地 LLM 做后处理。都可以绕开 Python，但质量会有损失。

### 4.7 汇总：哪些必须 Python，哪些不必

```
完全不需要 Python 即可实现（纯 Go）:
├── ✅ FFmpeg 音视频处理        → os/exec
├── ✅ ASR 语音识别             → whisper.cpp 或 sherpa-onnx Go 绑定
├── ✅ 字幕翻译                 → Ollama 本地 LLM (HTTP)
├── ✅ 人声分离                 → sherpa-onnx Go 绑定
├── ✅ 说话人分离               → sherpa-onnx Go 绑定
├── ✅ VAD 语音检测             → sherpa-onnx Go 绑定
├── ✅ 字幕解析 SRT/ASS/VTT    → Go 字符串处理
├── ✅ 音视频对齐               → rubberband 二进制 + Go 逻辑
└── ✅ TTS (基础)               → sherpa-onnx Piper TTS

需要 Python 才能达到最佳质量:
├── 🐍 TTS 声音克隆            → F5-TTS / CosyVoice / Qwen3-TTS (PyTorch)
├── 🐍 高质量 TTS              → Qwen3-TTS / GPT-SoVITS (PyTorch)
├── 🐍 降噪 (高质量)           → ModelScope FRCRN (PyTorch)
└── 🐍 标点恢复 (高质量)       → ModelScope CT-Transformer (PyTorch)
```

---

## 五、Python Sidecar 的打包分发方案

这是整个方案中**最大的工程挑战**。

### 5.1 方案对比

| 方案 | 原理 | 打包体积 | 用户体验 | 推荐度 |
|------|------|----------|----------|--------|
| **A. 内嵌 Python** | 用 PyInstaller 把 Python worker 打包为独立 exe | +200-500MB | ⭐⭐⭐⭐ 好 | **推荐** |
| **B. 要求用户装 Python** | 用户自行安装 Python + pip install | +0MB | ⭐ 差 | 不推荐 |
| **C. 内嵌 uv + venv** | 应用内携带 uv，首次启动自动创建 venv | +50MB | ⭐⭐⭐ 可 | 备选 |
| **D. conda-pack** | 预打包 conda 环境 | +1-3GB | ⭐⭐⭐ 可 | 体积太大 |
| **E. Docker** | 内嵌 Docker 容器运行 Python | +500MB+ | ⭐⭐ 差 | 不推荐桌面 |

### 5.2 推荐方案：PyInstaller 打包 Python Worker

```
最终分发结构:

MyApp/
├── myapp.exe                    # Wails 主程序 (Go, ~15MB)
├── python-worker/
│   ├── worker.exe               # PyInstaller 打包的 Python worker (~300-500MB)
│   └── models/                  # 模型文件 (按需下载)
│       ├── faster-whisper-large-v3/   (~3GB)
│       ├── Qwen3-TTS/                 (~3-7GB)
│       └── ...
├── ffmpeg/
│   ├── ffmpeg.exe               # (~80MB)
│   └── ffprobe.exe
└── rubberband.exe               # (~5MB, 可选)
```

**Go 启动 Python Worker 的代码示例**：

```go
type PythonSidecar struct {
    cmd     *exec.Cmd
    baseURL string
    port    int
}

func (s *PythonSidecar) Start() error {
    s.port = findFreePort()
    s.baseURL = fmt.Sprintf("http://127.0.0.1:%d", s.port)

    workerPath := filepath.Join(appDir, "python-worker", "worker.exe")
    s.cmd = exec.Command(workerPath, "--port", strconv.Itoa(s.port))
    s.cmd.Stdout = os.Stdout
    s.cmd.Stderr = os.Stderr

    if err := s.cmd.Start(); err != nil {
        return fmt.Errorf("failed to start python worker: %w", err)
    }

    // 等待 worker 就绪
    return s.waitForReady(30 * time.Second)
}

func (s *PythonSidecar) ASR(audioPath, model string) ([]Subtitle, error) {
    resp, err := http.Post(s.baseURL+"/asr", "application/json",
        toJSON(map[string]string{"audio": audioPath, "model": model}))
    // ...
}
```

**Python Worker 端 (FastAPI)**：

```python
# worker.py — 用 PyInstaller 打包为 worker.exe
from fastapi import FastAPI
import uvicorn

app = FastAPI()

@app.post("/asr")
async def asr(req: ASRRequest):
    from faster_whisper import WhisperModel
    model = WhisperModel(req.model, device="cuda" if torch.cuda.is_available() else "cpu")
    segments, _ = model.transcribe(req.audio)
    return {"subtitles": [{"start": s.start, "end": s.end, "text": s.text} for s in segments]}

@app.post("/tts")
async def tts(req: TTSRequest):
    # Qwen3-TTS / F5-TTS / CosyVoice...
    ...

@app.get("/health")
async def health():
    return {"status": "ok", "gpu": torch.cuda.is_available()}

if __name__ == "__main__":
    import argparse
    parser = argparse.ArgumentParser()
    parser.add_argument("--port", type=int, default=9876)
    args = parser.parse_args()
    uvicorn.run(app, host="127.0.0.1", port=args.port)
```

---

## 六、两种实施路线

### 路线 A：极简纯 Go（不要 Python）

**适用**：接受 TTS 质量降级，不需要声音克隆。

```
Go + Wails + whisper.cpp + Ollama + sherpa-onnx + FFmpeg

ASR: whisper.cpp (CGo 或进程调用)
翻译: Ollama (本地 HTTP，运行 Qwen2.5/DeepSeek)
TTS: sherpa-onnx Piper TTS (Go 绑定)
人声分离: sherpa-onnx (Go 绑定)
说话人分离: sherpa-onnx (Go 绑定)
音视频处理: FFmpeg (os/exec)
```

| 优点 | 缺点 |
|------|------|
| 单一技术栈，打包简单 | TTS 质量不如 PyTorch 方案 |
| 无 Python 依赖 | 无声音克隆能力 |
| 部署体积较小 | 降噪质量降级 |
| 启动速度快 | 模型选择受限 |

### 路线 B：Go + Python Sidecar（推荐）

**适用**：追求最佳质量，需要声音克隆。

```
Go + Wails (主控) + Python Worker (AI推理) + FFmpeg

Go 负责: UI、任务编排、FFmpeg调用、字幕处理、进度管理
Python 负责: ASR、TTS、翻译(可选)、降噪、说话人分离
```

| 优点 | 缺点 |
|------|------|
| AI 质量最佳 | 打包体积大 (+300-500MB Python Worker) |
| 完整声音克隆能力 | 双技术栈维护成本 |
| 可以使用所有 PyTorch 模型 | 首次启动较慢（加载模型）|
| 进程隔离，Python 崩溃不影响主程序 | 需要管理 Python 进程生命周期 |

### 路线 C：渐进混合（实际最推荐）

**先做纯 Go，再按需加 Python**。

```
Phase 1 (MVP): 纯 Go
  ASR → whisper.cpp
  翻译 → Ollama
  TTS → sherpa-onnx Piper
  → 验证核心流程，无 Python 依赖

Phase 2: 加入 Python Sidecar（可选安装）
  高质量 TTS → Qwen3-TTS / F5-TTS
  声音克隆 → CosyVoice
  高质量降噪 → ModelScope
  → 用户可选择是否安装 Python 增强包

Phase 3: 优化
  模型按需下载
  GPU 自动检测和调度
  多模型并行
```

**这种方式的优势**：

1. MVP 阶段零 Python 依赖，打包和分发极其简单
2. Python 作为可选增强，用户自主选择是否安装
3. 即使没有 Python，基础功能（ASR + 翻译 + 基础 TTS）完全可用
4. 降低首次使用门槛

---

## 七、关键技术风险和应对

| 风险 | 等级 | 应对 |
|------|------|------|
| PyInstaller 打包 PyTorch 体积过大 | 高 | 按功能拆分多个 worker.exe，用户按需下载 |
| Python Worker 启动慢（加载模型） | 中 | 预加载 + 模型常驻内存 + 健康检查 |
| Go ↔ Python 大文件传输开销 | 低 | 传文件路径而非文件内容，共享磁盘目录 |
| whisper.cpp CGo 交叉编译 | 中 | 改用 subprocess 调用 whisper.cpp 二进制 |
| Ollama 本地 LLM 内存占用 | 高 | 7B 模型需 ~8GB RAM，文档中说明硬件要求 |
| macOS/Linux 兼容性 | 中 | FFmpeg/whisper.cpp/sherpa-onnx 均跨平台 |
| GPU 版 PyTorch 打包体积 (~2GB) | 高 | GPU 版作为可选下载包，默认 CPU |

---

## 八、最终建议

```
推荐架构: Go (Wails) + Python Sidecar (可选)

核心原则:
1. Go 做所有能做的事（调度、FFmpeg、字幕、UI）
2. 纯 Go 方案覆盖基础功能（whisper.cpp + Ollama + sherpa-onnx）
3. Python Sidecar 作为「高级功能包」可选安装
4. 两者通过 HTTP (localhost) 通信，进程完全隔离

分发策略:
├── 基础版 (~200MB): Go 主程序 + FFmpeg + whisper.cpp + sherpa-onnx
├── 标准版 (+500MB): + Python Worker (CPU 版)
└── GPU 版 (+2GB): + CUDA 版 PyTorch
模型文件: 全部按需下载，不内置
```

这个架构相比原项目 pyVideoTrans 的优势：
- **启动更快**：Go 程序启动 <1s，不需要等 Python 解释器初始化
- **内存更省**：Go 主进程 ~50MB vs Python PySide6 ~300MB
- **更稳定**：Python AI 推理崩溃不影响主程序，自动重启
- **更好的并发**：goroutine 调度多个 FFmpeg/Python 任务
- **更小的基础包**：不需要 Python 也能用基础功能
