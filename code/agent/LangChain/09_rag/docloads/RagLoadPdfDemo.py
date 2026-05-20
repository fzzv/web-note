from langchain_community.document_loaders import PyPDFLoader
from pathlib import Path

file_path = Path(__file__).parent / "assets/sample.pdf"

docs = PyPDFLoader(
    # 文件路径，支持本地文件和在线文件链接，如"https://arxiv.org/pdf/alg-geom/9202012"
    file_path=file_path,
    # 提取模式:
    #   plain 提取文本
    #   layout 按布局提取
    extraction_mode="plain",
).load()

print(docs)
