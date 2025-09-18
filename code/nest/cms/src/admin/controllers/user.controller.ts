import { Body, Controller, Delete, Get, NotFoundException, Query, Param, ParseIntPipe, Headers, Post, Put, Redirect, Render, Res, UseFilters, HttpException } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto, UpdateUserRolesDto } from 'src/share/dtos/user.dto';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import type { Response } from 'express';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import { RoleService } from 'src/share/services/role.service';

@ApiTags('admin/user')
@UseFilters(AdminExceptionFilter)
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService,
    private readonly roleService: RoleService
  ) { }

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number) {
    const { users, total } = await this.userService.findAllWithPagination(page, limit, search);
    const pageCount = Math.ceil(total / limit);
    const roles = await this.roleService.findAll();
    return { users, search, page, limit, pageCount, roles };
  }

  @Get('create')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  @Render('user/user-form')
  async create() {
    return { user: {} };
  }

  @Post()
  @Redirect('/admin/user')
  @ApiOperation({ summary: '添加用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回添加用户页面' })
  async createUser(@Body() createUserDto: CreateUserDto) {
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  @Render('user/user-form')
  async edit(@Param('id', ParseIntPipe) id: number) {
    const user = await this.userService.findOne({ where: { id } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Put(':id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  async updateUser(
    @Param('id', ParseIntPipe) id: number, @Body() updateUserDto: UpdateUserDto,
    @Res() res: Response, @Headers('accept') accept: string
  ) {
    if (updateUserDto.password) {
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password);
    } else {
      delete updateUserDto.password;
    }
    await this.userService.update(id, updateUserDto);
    if (accept.includes('application/json')) {
      return res.json({ success: true, message: '用户更新成功' });
    } else {
      return res.redirect('/admin/user');
    }
  }

  @Get(':id')
  @ApiOperation({ summary: '获取用户详情(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户详情' })
  async findOne(@Param('id', ParseIntPipe) id: number, @Res() res: Response, @Headers('accept') accept: string) {
    const user = await this.userService.findOne({ where: { id }, relations: ['roles'] });
    if (!user) throw new HttpException('User not Found', 404)
    if (accept === 'application/json') {
      return res.json(user);
    } else {
      res.render('user/user-detail', { user });
    }
  }

  @Delete(':id')
  @ApiOperation({ summary: '删除用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回删除用户页面' })
  async deleteUser(@Param('id', ParseIntPipe) id: number) {
    await this.userService.delete(id);
    return { success: true, message: '用户删除成功' };
  }

  @Put(':id/roles')
  @ApiOperation({ summary: '更新用户角色(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回更新用户角色页面' })
  async updateRoles(@Param('id', ParseIntPipe) id: number, @Body() updateUserRolesDto: UpdateUserRolesDto) {
    await this.userService.updateRoles(id, updateUserRolesDto);
    return { success: true };
  }
}
