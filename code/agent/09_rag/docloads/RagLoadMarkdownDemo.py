# pip install unstructured[md]
from langchain_community.document_loaders import UnstructuredMarkdownLoader
from pathlib import Path

file_path = Path(__file__).parent / "assets/sample.md"

docs = UnstructuredMarkdownLoader(
    # 文件路径
    file_path=file_path,
    # 加载模式:
    #   single 返回单个Document对象
    #   elements 按标题等元素切分文档
    mode="elements",
).load()

print(docs)
