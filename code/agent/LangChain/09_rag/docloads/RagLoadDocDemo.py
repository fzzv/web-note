# pip install langchain_community unstructured[docx]
# pip install -U unstructured
# pip install python-docx
# pip install regex==2026.1.14
from langchain_community.document_loaders import UnstructuredWordDocumentLoader
from pathlib import Path

file_path = Path(__file__).parent / "assets/alibaba-more.docx"
# file_path = Path(__file__).parent / "assets/Doc21.docx"
# file_path = Path(__file__).parent / "assets/sample.doc"

docs = UnstructuredWordDocumentLoader(
    # 文件路径
    file_path=file_path,
    # 加载模式:
    #   single 返回单个Document对象
    #   elements 按标题等元素切分文档
    mode="single",
).load()

# print(type(docs))
print(docs)
