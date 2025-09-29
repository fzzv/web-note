import { Controller, Get, Render, Post, Redirect, Body, UseFilters, Param, ParseIntPipe, Put, Delete, Headers, Res, Query, NotFoundException } from '@nestjs/common';
import { CreateArticleDto, UpdateArticleDto } from 'src/share/dtos/article.dto';
import { ArticleService } from 'src/share/services/article.service';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import { CategoryService } from 'src/share/services/category.service';
import { TagService } from 'src/share/services/tag.service';
import type { Response } from 'express';
import { ArticleStateEnum } from 'src/share/enums/article.enum';
import { EventEmitter2 } from '@nestjs/event-emitter';

@UseFilters(AdminExceptionFilter)
@Controller('admin/articles')
export class ArticleController {
  constructor(
    private readonly articleService: ArticleService,
    private readonly categoryService: CategoryService,
    private readonly tagService: TagService,
    private readonly eventEmitter: EventEmitter2,
  ) { }

  @Get()
  @Render('article/article-list')
  async findAll(@Query('keyword') keyword: string = '',
    @Query('page', new ParseOptionalIntPipe(1)) page: number,
    @Query('limit', new ParseOptionalIntPipe(10)) limit: number) {
    const { articles, total } = await this.articleService.findAllWithPagination(page, limit, keyword);
    const pageCount = Math.ceil(total / limit);
    return { articles, keyword, page, limit, pageCount };
  }

  @Get('create')
  @Render('article/article-form')
  async createForm() {
    const categoryTree = await this.categoryService.findAll();
    const tags = await this.tagService.findAll();
    return { article: { categories: [], tags: [] }, categoryTree, tags };
  }

  @Post()
  @Redirect('/admin/articles')
  async create(@Body() createArticleDto: CreateArticleDto) {
    await this.articleService.create(createArticleDto);
    return { success: true }
  }

  @Get(':id/edit')
  @Render('article/article-form')
  async editForm(@Param('id', ParseIntPipe) id: number) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not Found');
    const categoryTree = await this.categoryService.findAll();
    const tags = await this.tagService.findAll();
    return { article, categoryTree, tags };
  }

  @Put(':id')
  async update(@Param('id', ParseIntPipe) id: number, @Body() updateArticleDto: UpdateArticleDto, @Res({ passthrough: true }) res: Response, @Headers('accept') accept: string) {
    await this.articleService.update(id, updateArticleDto);
    if (accept === 'application/json') {
      return { success: true };
    } else {
      return res.redirect(`/admin/articles`);
    }
  }

  @Delete(":id")
  async delete(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.delete(id);
    return { success: true }
  }

  @Get(':id')
  @Render('article/article-detail')
  async findOne(@Param('id', ParseIntPipe) id: number) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not Found');
    return { article };
  }
  
  @Put(':id/submit')
  async submitForReview(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.PENDING } as UpdateArticleDto);
    this.eventEmitter.emit('article.submitted', { articleId: id });
    return { success: true };
  }

  @Put(':id/approve')
  async approveArticle(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.PUBLISHED, rejectionReason: undefined } as UpdateArticleDto);
    return { success: true };
  }

  @Put(':id/reject')
  async rejectArticle(
    @Param('id', ParseIntPipe) id: number,
    @Body('rejectionReason') rejectionReason: string
  ) {
    await this.articleService.update(id, { state: ArticleStateEnum.REJECTED, rejectionReason } as UpdateArticleDto);
    return { success: true };
  }

  @Put(':id/withdraw')
  async withdrawArticle(@Param('id', ParseIntPipe) id: number) {
    await this.articleService.update(id, { state: ArticleStateEnum.WITHDRAWN } as UpdateArticleDto);
    return { success: true };
  }
}
