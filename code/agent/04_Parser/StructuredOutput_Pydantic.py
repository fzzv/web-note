"""
PydanticOutputParser 是 LangChain 输出解析器体系中最常用、最强大的结构化解析器之一。
它与 JsonOutputParser 类似，但功能更强 —— 能直接基于 Pydantic 模型 定义输出结构，
并利用其类型校验与自动文档能力。
对于结构更复杂、具有强类型约束的需求，PydanticOutputParser 则是最佳选择。
它结合了Pydantic模型的强大功能，提供了类型验证、数据转换等高级功能
"""

import os
from langchain.chat_models import init_chat_model
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from loguru import logger
from pydantic import BaseModel, Field, field_validator
from dotenv import load_dotenv

load_dotenv()


class Product(BaseModel):
    """
    产品信息模型类，用于定义产品的结构化数据格式

    属性:
        name (str): 产品名称
        category (str): 产品类别
        description (str): 产品简介，长度必须大于等于10个字符
    """

    name: str = Field(description="产品名称")
    category: str = Field(description="产品类别")
    description: str = Field(description="产品简介")

    @field_validator("description")
    def validate_description(cls, value):
        """
        验证产品简介字段的长度
        参数:
            value (str): 待验证的产品简介文本
        返回:
            str: 验证通过的产品简介文本
        异常:
            ValueError: 当产品简介长度小于10个字符时抛出
        """
        if len(value) < 10:
            raise ValueError("产品简介长度必须大于等于10")
        return value


# 创建Pydantic输出解析器实例，用于解析模型输出为Product对象
parser = PydanticOutputParser(pydantic_object=Product)

# 获取格式化指令，用于指导模型输出符合Product模型的JSON格式
format_instructions = parser.get_format_instructions()

# 创建聊天提示模板，包含系统消息和人类消息
prompt_template = ChatPromptTemplate.from_messages(
    [
        ("system", "你是一个AI助手，你只能输出结构化的json数据\n{format_instructions}"),
        ("human", "请你输出标题为：{topic}的新闻内容"),
    ]
)

# 格式化提示消息，填充主题和格式化指令
prompt = prompt_template.format_messages(
    topic="小米Yu7", format_instructions=format_instructions
)

# 记录格式化后的提示消息
logger.info(prompt)

# 创建模型
model = init_chat_model(
    model="mimo-v2.5-pro",
    model_provider="openai",
    api_key=os.getenv("XIAOMI_API_KEY"),
    base_url="https://token-plan-cn.xiaomimimo.com/v1",
)

# 调用模型获取结果
result = model.invoke(prompt)

# 记录模型返回的结果
logger.info(f"模型原始输出:\n{result.content}")

# 使用解析器将模型结果解析为Product对象
response = parser.invoke(result)

# 打印解析后的结构化结果
logger.info(f"解析后的结构化结果:\n{response}")

# 打印类型
logger.info(f"结果类型: {type(response)}")
