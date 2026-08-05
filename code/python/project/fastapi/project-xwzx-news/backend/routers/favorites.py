from fastapi import APIRouter, Depends, Query, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from starlette import status

from config.db_conf import get_db
from crud import favorites
from models.user import User
from schemas.favorites import FavoriteListResponse, FavoriteCheckResponse, FavoriteAddRequest
from utils.auth import get_current_user
from utils.response import success_response

router = APIRouter(prefix="/api/favorite", tags=["favorites"])

@router.get("/list")
async def get_favorite_list(
        page: int = 1,
        page_size: int = Query(10, alias="pageSize", le=100),
        user: User = Depends(get_current_user),
        db: AsyncSession = Depends(get_db)
):
    offset = (page - 1) * page_size
    rows, total = await favorites.get_favorites_list(db, user, offset, page_size)

    favorite_list = [{
        **news.__dict__,
        "favorite_time": favorite_time,
        "favorite_id": favorite_id
    } for news, favorite_time, favorite_id in rows]
    has_more = total > page * page_size

    data = FavoriteListResponse(list=favorite_list, total=total, hasMore=has_more)
    return success_response(message="获取列表成功", data=data)

@router.get("/check")
async def get_favorite_check(news_id: int = Query(..., alias="newsId"), user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    is_favorite = await favorites.check_status(db, news_id, user)
    return success_response(message="查询成功", data=FavoriteCheckResponse(isFavorite=is_favorite))

@router.post("/add")
async def add_favorite(data: FavoriteAddRequest, user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    result = await favorites.add_news_to_favorite(db, data.news_id, user)
    return success_response(message="添加收藏成功", data=result)

@router.delete("/remove")
async def remove_favorite(news_id: int = Query(..., alias="newsId"), user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    result = await favorites.cancel_favorite(db, news_id, user)
    if not result:
        raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail="收藏记录不存在")
    return success_response(message="取消收藏成功")

@router.delete("/clear")
async def remove_all_favorite(user: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    count = await favorites.clear_favorite(db, user.id)
    return success_response(message=f"清空了{count}条收藏")