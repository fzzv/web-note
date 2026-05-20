from langchain_redis import RedisConfig, RedisVectorStore
from langchain_community.embeddings import DashScopeEmbeddings
import os
from dotenv import load_dotenv

load_dotenv()

# 初始化 Embedding 模型
# 1. 初始化阿里千问 Embedding 模型
embeddingsModel = DashScopeEmbeddings(
    model="text-embedding-v3",  # 支持 v1 或 v2
    dashscope_api_key=os.getenv("QWEN_API_KEY"),  # 从环境变量读取
)

# 2. 创建Redis向量存储实例
vector_store = RedisVectorStore(
    embeddingsModel,
    config=RedisConfig(index_name="newsgroups", redis_url="redis://localhost:6389"),
)

# ========== 查询数据 ==========
# 定义查询文本
query = "我喜欢吃什么？"

# 3. 将查询语句向量化，并在Redis中做相似度检索
results = vector_store.similarity_search_with_score(query, k=3)

print("=== 查询结果 ===")
for i, (doc, score) in enumerate(results, 1):
    similarity = 1 - score  #  score 是距离，可以转成相似度
    print(f"结果 {i}:")
    print(f"内容: {doc.page_content}")
    print(f"元数据: {doc.metadata}")
    print(f"相似度: {similarity:.4f}")

# 结果 1:
# 内容: 我喜欢吃苹果
# 元数据: {'segment_id': '1'}
# 相似度: 0.7424
# 结果 2:
# 内容: 苹果是我最喜欢吃的水果
# 元数据: {'segment_id': '2'}
# 相似度: 0.6670
# 结果 3:
# 内容: 我喜欢用苹果手机
# 元数据: {'segment_id': '3'}
# 相似度: 0.5940
