# pip install jq / uv add jq
from langchain_community.document_loaders import JSONLoader
from pathlib import Path

file_path = Path(__file__).parent / "assets/sample.json"

# 提取所有字段
docs = JSONLoader(
    file_path=file_path,  # 文件路径
    jq_schema=".",  # 提取所有字段
    text_content=False,  # 提取内容是否为字符串格式
).load()

print(docs)
