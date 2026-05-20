from langchain_community.document_loaders import TextLoader
from pathlib import Path

# 返回List[Document]
file_path = Path(__file__).parent / "assets/sample.txt"  # 文件路径
encoding = "utf-8"  # 文件编码方式

docs = TextLoader(file_path, encoding).load()

print(docs)
# [Document(metadata={'source': 'asset/sample.txt'}, page_content='...')]
