import { Controller, Get, Render, Post, Redirect, Body, UseFilters, Param, ParseIntPipe, Put, Delete, Headers, Res, Query, NotFoundException, StreamableFile, Header } from '@nestjs/common';
import { CreateArticleDto, UpdateArticleDto } from 'src/share/dtos/article.dto';
import { ArticleService } from 'src/share/services/article.service';
import { AdminExceptionFilter } from '../filters/admin-exception.filter';
import { ParseOptionalIntPipe } from 'src/share/pipes/parse-optional-int.pipe';
import { CategoryService } from 'src/share/services/category.service';
import { TagService } from 'src/share/services/tag.service';
import type { Response } from 'express';
import { ArticleStateEnum } from 'src/share/enums/article.enum';
import { EventEmitter2 } from '@nestjs/event-emitter';
import { WordExportService } from 'src/share/services/word-export.service';
import { PptExportService } from 'src/share/services/ppt-export.service';
import { ExcelExportService } from 'src/share/services/excel-export.service';

@UseFilters(AdminExceptionFilter)
@Controller('admin/articles')
export class ArticleController {
  constructor(
    private readonly articleService: ArticleService,
    private readonly categoryService: CategoryService,
    private readonly tagService: TagService,
    private readonly eventEmitter: EventEmitter2,
    private readonly wordExportService: WordExportService,
    private readonly pptExportService: PptExportService,
    private readonly excelExportService: ExcelExportService
  ) { }

  @Get('export-excel')
  @Header('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
  async exportExcel(@Query('search') search: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number, @Res({ passthrough: true }) res: Response) {
    const { articles } = await this.articleService.findAllWithPagination(page, limit, search);
    const data = articles.map(article => ({
      title: article.title,
      categories: article.categories.map(c => c.name).join(', '),
      tags: article.tags.map(t => t.name).join(', '),
      state: article.state,
      createdAt: article.createdAt,
    }));
    const columns = [
      { header: '标题', key: 'title', width: 30 },
      { header: '分类', key: 'categories', width: 20 },
      { header: '标签', key: 'tags', width: 20 },
      { header: '状态', key: 'state', width: 15 },
      { header: '创建时间', key: 'createdAt', width: 20 },
    ];
    const buffer = await this.excelExportService.exportAsExcel(data, columns);
    res.setHeader('Content-Disposition', `attachment; filename="articles.xlsx"`);
    return new StreamableFile(new Uint8Array(buffer));
  }
  
  @Get('export-ppt')
  @Header('Content-Type', 'application/vnd.openxmlformats-officedocument.presentationml.presentation')
  async exportPpt(@Query('keyword') keyword: string = '', @Query('page', new ParseOptionalIntPipe(1)) page: number, @Query('limit', new ParseOptionalIntPipe(10)) limit: number, @Res({ passthrough: true }) res: Response) {
    const { articles } = await this.articleService.findAllWithPagination(page, limit, keyword);
    const buffer = await this.pptExportService.exportToPpt(articles);
    res.setHeader('Content-Disposition', 'attachment; filename=articles.pptx');
    return new StreamableFile(buffer);
  }

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

  @Get(':id/export-word')
  @Header('Content-Type', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document')
  async exportWord(@Param('id', ParseIntPipe) id: number, @Res({ passthrough: true }) res: Response) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not found');

    const htmlContent = `
           <h1>${article.title}</h1>
           <p><strong>状态:</strong> ${article.state}</p>
           <p><strong>分类:</strong> ${article.categories.map(c => c.name).join(', ')}</p>
           <p><strong>标签:</strong> ${article.tags.map(t => t.name).join(', ')}</p>
           <hr/>
           ${article.content}
       `;

    const buffer = await this.wordExportService.exportToWord(htmlContent);
    res.setHeader('Content-Disposition', `attachment; filename="${encodeURIComponent(article.title)}.docx"`);
    return new StreamableFile(buffer);
  }
}
