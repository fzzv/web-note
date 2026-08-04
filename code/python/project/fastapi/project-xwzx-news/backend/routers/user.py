from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.ext.asyncio import AsyncSession
from starlette import status

from config.db_conf import get_db
from crud import user
from models.user import User
from schemas.user import UserRequest, UserAuthResponse, UserInfoResponse, UserUpdateRequest, UserChangePasswordRequest
from utils.auth import get_current_user
from utils.response import success_response

router = APIRouter(prefix="/api/user", tags=["user"])

@router.post("/register")
async def register_user(user_data: UserRequest, db: AsyncSession = Depends(get_db)):
    """
    用户注册
    :param user_data:
    :param db:
    :return:
    """
    # 检查用户名是否存在
    existing_user = await user.get_user_by_username(db, user_data.username)
    if existing_user:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="用户已存在")
    # 创建新用户
    users = await user.create_user(db, user_data)
    token = await user.create_token(db, users.id)
    response_data = UserAuthResponse(token=token, user_info=UserInfoResponse.model_validate(users))
    return success_response(message="注册成功", data=response_data)

@router.post("/login")
async def login_user(user_data: UserRequest, db: AsyncSession = Depends(get_db)):
    """
    用户登录
    :param user_data:
    :param db:
    :return:
    """
    existing_user = await user.get_user_by_username(db, user_data.username)
    if not existing_user:
        raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="用户不存在")
    users = await user.login(db, user_data)
    if not users:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="用户名或密码错误")
    token = await user.create_token(db, users.id)

    response_data = UserAuthResponse(token=token, user_info=UserInfoResponse.model_validate(users))
    return success_response(message="登录成功", data=response_data)

@router.get("/info")
async def user_info(users: User = Depends(get_current_user)):
    print(users)
    return success_response(message="获取用户信息成功", data=UserInfoResponse.model_validate(users))

@router.put("/update")
async def update_user_info(user_data: UserUpdateRequest, users: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    users = await user.update_user(db, users.username, user_data)
    return success_response(message="编辑用户信息成功", data=UserInfoResponse.model_validate(users))

@router.put("/password")
async def change_password(user_data: UserChangePasswordRequest, users: User = Depends(get_current_user), db: AsyncSession = Depends(get_db)):
    res_change_pwd = await user.change_password(db, user_data, users)
    if not res_change_pwd:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="修改密码失败，请稍后再试")
    return success_response(message="密码修改成功")