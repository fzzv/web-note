from langchain_community.document_loaders.csv_loader import CSVLoader
from pathlib import Path

file_path = Path(__file__).parent / "assets/sample.csv"

# 加载所有列
docs = CSVLoader(
    file_path=file_path,  # 文件路径
).load()  # 返回List[Document]

print(docs)

# 加载部分列
docs = CSVLoader(
    file_path=file_path,  # 文件路径
    metadata_columns=["title", "author"],  # 将指定列作为元数据
    content_columns=["content"],  # 将指定列作为内容
).load()  # 返回List[Document]

print(docs)
