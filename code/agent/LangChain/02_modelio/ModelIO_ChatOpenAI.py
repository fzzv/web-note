from langchain_openai import ChatOpenAI
import os
from dotenv import load_dotenv
from pydantic import SecretStr

load_dotenv()

api_key = os.getenv("XIAOMI_API_KEY")

chatLLM = ChatOpenAI(
    api_key=SecretStr(api_key) if api_key else None,
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
    model="mimo-v2.5",
    # other params...
)

messages = [
    {"role": "system", "content": "You are a helpful assistant."},
    {"role": "user", "content": "你是谁？"},
]

response = chatLLM.invoke(messages)

print(response.content)
