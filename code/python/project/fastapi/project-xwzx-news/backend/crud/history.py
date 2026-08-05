from datetime import datetime

from sqlalchemy import select, func, delete
from sqlalchemy.ext.asyncio import AsyncSession

from models.history import History
from models.news import News
from models.user import User


async def get_history_list(db: AsyncSession, user: User, offset: int, page_size: int):
    count_query = select(func.count()).where(History.user_id == user.id)
    count_result = await db.execute(count_query)
    total = count_result.scalar()

    query = (select(News, History.id.label("history_id"), History.view_time.label("view_time"))
             .join(History, History.news_id == News.id)
             .where(History.user_id == user.id)
             .order_by(History.view_time.desc())
             .offset(offset).limit(page_size)
             )
    result = await db.execute(query)
    rows = result.all()
    return rows, total

async def add_news_to_history(db: AsyncSession, news_id: int, user_id: int):
    """
    添加历史记录
    """
    query = select(History).where(History.user_id == user_id, History.news_id == news_id)
    result = await db.execute(query)
    existing_history = result.scalar_one_or_none()
    if existing_history:
        existing_history.view_time = datetime.now()
        await db.commit()
        await db.refresh(existing_history)
        return existing_history
    else:
        history = History(user_id=user_id, news_id=news_id)
        db.add(history)
        await db.commit()
        await db.refresh(history)
        return history

async def delete_history(db: AsyncSession, user_id: int, history_id: int):
    """
    删除单条浏览历史
    :param db:
    :param user_id:
    :param history_id:
    :return:
    """
    query = delete(History).where(History.user_id == user_id, History.id == history_id)
    result = await db.execute(query)
    await db.commit()
    return result.rowcount > 0

async def delete_all_history(db: AsyncSession, user_id: int):
    """
    清空所有浏览历史
    :param db:
    :param user_id:
    :return:
    """
    query = delete(History).where(History.user_id == user_id)
    result = await db.execute(query)
    await db.commit()
    return result.rowcount or 0