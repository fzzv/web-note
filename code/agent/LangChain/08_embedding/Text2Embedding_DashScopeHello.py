# https://bailian.console.aliyun.com/cn-beijing/?productCode=p_efm&tab=doc#/doc/?type=model&url=2842587

import dashscope
import os
from http import HTTPStatus
from dotenv import load_dotenv

load_dotenv()

input_text = "衣服的质量杠杠的"

dashscope.api_key = os.getenv("QWEN_API_KEY")

resp = dashscope.TextEmbedding.call(
    model="text-embedding-v4",
    input=input_text,
)

if resp.status_code == HTTPStatus.OK:
    print(resp)
