"""
把文本转换成向量有什么用呢？
最核心的作用是可以通过向量之间的计算，来分析文本与文本之间的相似性。
计算的方法有很多种，其中用得最多的是向量余弦相似度。
Python语言中提供了一个库sklearn，可以很方便的计算向量之间的余弦相似度
"""

import dashscope
import os
from http import HTTPStatus
import numpy as np
from dotenv import load_dotenv

load_dotenv()


# 准备输入文本数据
texts = ["我喜欢吃苹果", "苹果是我最喜欢吃的水果", "我喜欢用苹果手机"]

# 获取每个文本的embedding向量
embeddings = []
dashscope.api_key = os.getenv("QWEN_API_KEY")
# 假如要处理图片，请参考https://bailian.console.aliyun.com/cn-beijing/?productCode=p_efm&tab=doc#/doc/?type=model&url=2842587
for text in texts:
    input_data = [{"text": text}]
    resp = dashscope.MultiModalEmbedding.call(
        model="multimodal-embedding-v1",
        base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
        input=input_data,
    )

    if resp.status_code == HTTPStatus.OK:
        embedding = resp.output["embeddings"][0]["embedding"]
        embeddings.append(embedding)


# 计算余弦相似度
def cosine_similarity(vec1, vec2):
    # 计算两个向量的余弦相似度
    dot_product = np.dot(vec1, vec2)
    norm_vec1 = np.linalg.norm(vec1)
    norm_vec2 = np.linalg.norm(vec2)
    return dot_product / (norm_vec1 * norm_vec2)


# 比较所有文本之间的相似度
print("文本相似度比较结果:")
print("=" * 40)

for i in range(len(texts)):
    for j in range(i + 1, len(texts)):
        similarity = cosine_similarity(embeddings[i], embeddings[j])
        print(f"文本{i + 1} vs 文本{j + 1}:")
        print(f"  文本{i + 1}: {texts[i]}")
        print(f"  文本{j + 1}: {texts[j]}")
        print(f"  余弦相似度: {similarity:.4f}")
        print("-" * 40)

# 文本相似度比较结果:
# ========================================
# 文本1 vs 文本2:
#   文本1: 我喜欢吃苹果
#   文本2: 苹果是我最喜欢吃的水果
#   余弦相似度: 0.9064
# ----------------------------------------
# 文本1 vs 文本3:
#   文本1: 我喜欢吃苹果
#   文本3: 我喜欢用苹果手机
#   余弦相似度: 0.7656
# ----------------------------------------
# 文本2 vs 文本3:
#   文本2: 苹果是我最喜欢吃的水果
#   文本3: 我喜欢用苹果手机
#   余弦相似度: 0.7421
# ----------------------------------------
