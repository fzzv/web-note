import os

from langchain_core.messages import HumanMessage, AIMessage
from langchain_openai import ChatOpenAI
from langgraph.prebuilt import create_react_agent
from langgraph_supervisor import create_supervisor
from dotenv import load_dotenv

load_dotenv()


"""
langgraph.prebuilt.create_react_agent 在 LangGraph v1.0+ 中已被弃用，并将在 v2.0 中移除。

使用 langchain.agents.create_agent 来替代它
"""


def print_chinese_messages(chunk: dict):
    """
    只打印各个角色（supervisor / flight_assistant / hotel_assistant）
    输出的中文 content，过滤 tool / 空消息 / 英文控制信息
    """
    if not isinstance(chunk, dict):
        return

    for role, payload in chunk.items():
        if not isinstance(payload, dict):
            continue

        messages = payload.get("messages", [])
        for msg in messages:
            if isinstance(msg, (HumanMessage, AIMessage)):
                content = (msg.content or "").strip()
                if not content:
                    continue
                # 过滤明显的系统英文控制语
                if content.startswith("Transferring"):
                    continue
                if "Successfully transferred" in content:
                    continue

                print(f"{role}：{content}\n")


def init_llm_model() -> ChatOpenAI:
    """初始化大语言模型（ChatOpenAI）"""
    try:
        model = ChatOpenAI(
            model="qwen-plus",
            api_key=os.getenv("QWEN_API_KEY"),
            base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
            temperature=0.1,  # 低随机性，保证任务执行稳定性
            max_tokens=1024,
        )
        print("✅ 语言模型初始化成功")
        return model
    except Exception as e:
        print(f"❌ 语言模型初始化失败：{str(e)}")
        raise SystemExit(1)


def book_hotel(hotel_name: str):
    """预订酒店"""
    print(f"✅ 成功预订了 {hotel_name} 的住宿")
    return f"成功预订了 {hotel_name} 的住宿。"


def book_flight(from_airport: str, to_airport: str):
    """预订航班"""
    print(f"✅ 成功预订了从 {from_airport} 到 {to_airport} 的航班")
    return f"成功预订了从 {from_airport} 到 {to_airport} 的航班。"


flight_assistant = create_react_agent(
    model=init_llm_model(),
    tools=[book_flight],
    prompt=(
        "你是专业的航班预订助手，专注于帮助用户预订机票。\n"
        "工作流程：\n"
        "1. 从用户需求中提取出发地和目的地信息\n"
        "2. 调用book_flight工具完成预订\n"
        "3. 收到预订成功的确认后，向主管汇报结果并结束\n"
        "注意：每次只处理一个预订请求，完成后立即结束，不要重复调用工具。"
    ),
    name="flight_assistant",
)
hotel_assistant = create_react_agent(
    model=init_llm_model(),
    tools=[book_hotel],
    prompt=(
        "你是专业的酒店预订助手，专注于帮助用户预订酒店。\n"
        "工作流程：\n"
        "1. 从用户需求中提取酒店信息（如果未指定，选择经济型酒店）\n"
        "2. 调用book_hotel工具完成预订\n"
        "3. 收到预订成功的确认后，向主管汇报结果并结束\n"
        "注意：每次只处理一个预订请求，完成后立即结束，不要重复调用工具。"
    ),
    name="hotel_assistant",
)
supervisor = create_supervisor(
    agents=[flight_assistant, hotel_assistant],
    model=init_llm_model(),
    prompt=(
        "你是一个智能任务调度主管，负责协调航班预订助手(flight_assistant)和酒店预订助手(hotel_assistant)。\n\n"
        "工作流程：\n"
        "1. 分析用户需求，确定需要哪些服务（航班、酒店或两者）\n"
        "2. 如果需要预订航班，调用flight_assistant一次\n"
        "3. 如果需要预订酒店，调用hotel_assistant一次\n"
        "4. 收到助手的预订确认后，记录结果\n"
        "5. 当所有任务都完成后，向用户汇总所有预订结果，然后立即结束\n\n"
        "关键规则：\n"
        "- 每个助手只能调用一次，不要重复调用\n"
        "- 看到'成功预订'的消息后，该任务就已完成\n"
        "- 所有任务完成后，必须直接结束，不要再调用任何助手\n"
        "- 如果已经看到航班和酒店的预订确认，立即汇总并结束"
    ),
).compile()

for chunk in supervisor.stream(
    {
        "messages": [
            {
                "role": "user",
                "content": "帮我预定一个北京到深圳的机票，并且预定一个酒店",
            }
        ]
    }
):
    print(chunk)
    print("\n")
