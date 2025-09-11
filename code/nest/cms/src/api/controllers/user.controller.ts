import { Controller, Get, Post, Put, Delete, Body, Param, ParseIntPipe, UseInterceptors, SerializeOptions, HttpException, HttpStatus } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { CreateUserDto, UpdateUserDto } from '../../share/dtos/user.dto';
import { ApiOperation, ApiResponse, ApiParam, ApiBody, ApiTags, ApiBearerAuth } from '@nestjs/swagger';
import { ClassSerializerInterceptor } from '@nestjs/common';
import { User } from '../../share/entities/user.entity';
import { Result } from '../../share/vo/result';

@ApiTags('api/users')
@SerializeOptions({ strategy: 'exposeAll' })
@UseInterceptors(ClassSerializerInterceptor)
@Controller('api/users')
export class UserController {
  constructor(private readonly userService: UserService) { }

  @Get()
  @ApiOperation({ summary: '获取所有用户列表' })
  @ApiResponse({ status: 200, description: '成功返回用户列表', type: [User] })
  async findAll() {
    return this.userService.findAll();
  }

  @Get(':id')
  @ApiOperation({ summary: '根据ID获取用户信息' })
  @ApiParam({ name: 'id', description: '用户ID', type: Number })
  @ApiResponse({ status: 200, description: '成功返回用户信息', type: User })
  @ApiResponse({ status: 404, description: '用户未找到' })
  async findOne(@Param('id', ParseIntPipe) id: number) {
    return this.userService.findOne({ where: { id } });
  }

  @Post()
  @ApiBearerAuth()
  @ApiOperation({ summary: '创建新用户' })
  @ApiBody({ type: CreateUserDto })
  @ApiResponse({ status: 201, description: '用户成功创建', type: User })
  @ApiResponse({ status: 400, description: '请求参数错误' })
  async create(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }

  @Put(':id')
  @ApiOperation({ summary: '更新用户信息' })
  @ApiParam({ name: 'id', description: '用户ID', type: Number })
  @ApiBody({ type: UpdateUserDto })
  @ApiResponse({ status: 200, description: '用户信息更新成功', type: Result })
  @ApiResponse({ status: 400, description: '请求参数错误' })
  @ApiResponse({ status: 404, description: '用户未找到' })
  async update(@Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto) {
    const updateResult = await this.userService.update(id, updateUserDto);
    if (!updateResult.affected) {
      throw new HttpException('用户未找到', HttpStatus.NOT_FOUND);
    }
    return Result.success('用户信息更新成功');
  }

  @Delete(':id')
  @ApiOperation({ summary: '删除用户' })
  @ApiParam({ name: 'id', description: '用户ID', type: Number })
  @ApiResponse({ status: 200, description: '用户删除成功', type: Result })
  @ApiResponse({ status: 404, description: '用户未找到' })
  async delete(@Param('id', ParseIntPipe) id: number) {
    const deleteResult = await this.userService.delete(id);
    if (!deleteResult.affected) {
      throw new HttpException('用户未找到', HttpStatus.NOT_FOUND);
    }
    return Result.success('用户删除成功');
  }
}
