import dashscope
import json
import os
from http import HTTPStatus
from dotenv import load_dotenv

load_dotenv()

# Embedding 文本向量化
text = "通用多模态表征模型示例"
input = [{"text": text}]
# 调用多模态embedding模型接口进行向量编码
# https://bailian.console.aliyun.com/?productCode=p_efm&tab=model#/model-market/all?capabilities=ME
dashscope.api_key = os.getenv("QWEN_API_KEY")
resp = dashscope.MultiModalEmbedding.call(
    model="tongyi-embedding-vision-plus",  # 支持 v1 或 v2
    input=input,
)

result = ""
# 处理模型返回结果，提取关键信息并格式化输出
if resp.status_code == HTTPStatus.OK:
    result = {
        "status_code": resp.status_code,
        "request_id": getattr(resp, "request_id", ""),
        "code": getattr(resp, "code", ""),
        "message": getattr(resp, "message", ""),
        "output": resp.output,
        "usage": resp.usage,
    }
    print(json.dumps(result, ensure_ascii=False, indent=4))

print("=================================")
print()

# result 就是你已经拿到的完整 dict
# embedding_values = result["output"]["embeddings"][0]["embedding"]
# print(embedding_values)
# print("=================================")
# print("=================================")
# # 只打印 embedding 数组
# print(json.dumps(embedding_values, ensure_ascii=False))
