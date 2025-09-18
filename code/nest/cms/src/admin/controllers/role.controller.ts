import { Controller, Get, Post, Put, Delete, Body, Param, Query, Res, Render, Redirect, UseFilters, Headers, ParseIntPipe } from '@nestjs/common';
import { RoleService } from 'src/share/services/role.service';
import { CreateRoleDto, UpdateRoleAccessesDto, UpdateRoleDto } from 'src/share/dtos/role.dto';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import type { Response } from 'express';
import { AccessService } from 'src/share/services/access.service';

@UseFilters(AdminExceptionFilter)
@Controller('admin/roles')
export class RoleController {
  constructor(
    private readonly roleService: RoleService,
    private readonly accessService: AccessService
  ) { }

  @Get()
  @Render('role/role-list')
  async findAll(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number) {
    const { roles, total } = await this.roleService.findAllWithPagination(page, limit, search);
    const pageCount = Math.ceil(total / limit);
    const accessTree = await this.accessService.findAll();
    return { search, page, limit, roles, pageCount, accessTree };
  }

  @Get('create')
  @Render('role/role-form')
  createForm() {
    return { role: {} };
  }

  @Post()
  @Redirect('/admin/roles')
  async create(@Body() createRoleDto: CreateRoleDto) {
    await this.roleService.create(createRoleDto);
    return { success: true };
  }

  @Get(':id/edit')
  @Render('role/role-form')
  async editForm(@Param('id', ParseIntPipe) id: number) {
    const role = await this.roleService.findOne({ where: { id }, relations: ['accesses'] });
    if (!role) throw new Error('Role not found');
    return { role };
  }

  @Put(':id')
  async update(@Param('id', ParseIntPipe) id: number, @Body() updateRoleDto: UpdateRoleDto, @Res() res: Response, @Headers('accept') accept: string) {
    await this.roleService.update(id, updateRoleDto);
    if (accept === 'application/json') {
      return res.json({ success: true });
    } else {
      return res.redirect('/admin/roles');
    }
  }

  @Get(':id')
  @Render('role/role-detail')
  async findOne(@Param('id', ParseIntPipe) id: number, @Res() res: Response, @Headers('accept') accept: string) {
    const role = await this.roleService.findOne({ where: { id }, relations: ['accesses'] });
    if (!role) throw new Error('Role not found');
    if (accept === 'application/json') {
      return res.json(role);
    } else {
      return res.render('role/role-detail', { role });
    }
  }


  @Put(':id/accesses')
  async updateAccesses(@Param('id', ParseIntPipe) id: number, @Body() updateRoleAccessesDto: UpdateRoleAccessesDto) {
    await this.roleService.updateAccesses(id, updateRoleAccessesDto);
    return { success: true };
  }

  @Delete(':id')
  async delete(@Param('id', ParseIntPipe) id: number) {
    await this.roleService.delete(id);
    return { success: true };
  }
}
