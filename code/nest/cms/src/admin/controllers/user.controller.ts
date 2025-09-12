import { Body, Controller, Get, NotFoundException, Param, Post, Put, Redirect, Render } from '@nestjs/common';
import { UserService } from '../../share/services/user.service';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { UtilityService } from '../../share/services/utility.service';
import { CreateUserDto, UpdateUserDto } from 'src/share/dtos/user.dto';

@ApiTags('admin/user')
@Controller('admin/user')
export class UserController {

  constructor(
    private readonly userService: UserService,
    private readonly utilityService: UtilityService
  ) {}

  @Get()
  @ApiOperation({ summary: '获取所有用户列表(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户列表' })
  @Render('user/user-list')
  async findAll() {
    const users = await this.userService.findAll();
    return { users };
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
    console.log(createUserDto, 'createUserDto')
    const hashedPassword = await this.utilityService.hashPassword(createUserDto.password);
    await this.userService.create({ ...createUserDto, password: hashedPassword });
    return { url: '/admin/user', success: true, message: '用户添加成功' };
  }

  @Get('edit/:id')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  @Render('user/user-form')
  async edit(@Param('id') id: string) {
    const user = await this.userService.findOne({ where: { id: Number(id) } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }

  @Put(':id')
  @Redirect('/admin/user')
  @ApiOperation({ summary: '编辑用户(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回编辑用户页面' })
  async updateUser(@Param('id') id: string, @Body() updateUserDto: UpdateUserDto) {
    if (updateUserDto.password) {
      updateUserDto.password = await this.utilityService.hashPassword(updateUserDto.password);
    } else {
      delete updateUserDto.password;
    }
    await this.userService.update(Number(id), updateUserDto);
    return { url: '/admin/user', success: true, message: '用户更新成功' };
  }

  @Get(':id')
  @ApiOperation({ summary: '获取用户详情(管理后台)' })
  @ApiResponse({ status: 200, description: '成功返回用户详情' })
  @Render('user/user-detail')
  async findOne(@Param('id') id: string) {
    const user = await this.userService.findOne({ where: { id: Number(id) } });
    if (!user) {
      throw new NotFoundException('用户不存在');
    }
    return { user };
  }
}
