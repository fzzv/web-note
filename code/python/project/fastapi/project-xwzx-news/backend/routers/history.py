from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession

from config.db_conf import get_db
from crud import history
from models.user import User
from schemas.history import HistoryListResponse, HistoryAddRequest
from utils.auth import get_current_user
from utils.response import success_response

router = APIRouter(prefix="/api/history", tags=["history"])

@router.get("/list")
async def get_history(
        page: int = Query(1, ge=1),
        page_size: int = Query(10, ge=1, le=100, alias="pageSize"),
        db: AsyncSession = Depends(get_db),
        user: User = Depends(get_current_user)):
    offset = (page - 1) * page_size
    rows, total = await history.get_history_list(db, user, offset, page_size)

    history_list = [{
        **news.__dict__,
        "history_id": history_id,
        "view_time": view_time
    } for news, history_id, view_time in rows]

    has_more = total > page * page_size
    data = HistoryListResponse(list=history_list,total=total,hasMore=has_more)

    return success_response(message="获取浏览历史列表成功", data=data)

@router.post("/add")
async def add_history(data: HistoryAddRequest, user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    result = await history.add_news_to_history(db, data.news_id, user.id)
    return success_response(message="添加浏览历史成功", data=result)

@router.delete("/delete/{history_id}")
async def delete_history(history_id: int, user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    result = await history.delete_history(db, user.id, history_id)
    return success_response(message="删除浏览历史成功", data=result)

@router.delete("/clear")
async def clear_history(db: AsyncSession = Depends(get_db), user: User = Depends(get_current_user)):
    count = await history.delete_all_history(db, user.id)
    return success_response(message=f"清空了{count}条浏览记录")