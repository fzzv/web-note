# 方式1：外部加载Prompt,将 prompt 保存为 JSON
from langchain_core.prompts import load_prompt
from pathlib import Path

prompt_file = Path(__file__).parent / "prompt.json"

template = load_prompt(prompt_file, encoding="utf-8")
print(template.format(name="张三", what="搞笑的"))
# 请张三讲一个搞笑的的故事

# load_prompt在新版LangChain中将被弃用，改成序列化和反序列化使用 dumpd() / dumps() 和 load() / loads() 来实现
