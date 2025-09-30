import { Controller, Get, NotFoundException, Param, Query } from '@nestjs/common';
import { ArticleService } from '../../share/services/article.service';

@Controller('api/articles')
export class ArticleController {
  constructor(
    private readonly articleService: ArticleService,
  ) { }

  @Get()
  async getArticles(
    @Query('categoryId') categoryId: string = '',
    @Query('tagId') tagId: string = '',
    @Query('keyword') keyword: string = '',
  ) {
    const articles = await this.articleService.findList(
      keyword,
      categoryId,
      tagId
    );
    return {
      keyword,
      categoryId,
      tagId,
      articles
    };
  }

  @Get(':id')
  async getArticleById(@Param('id') id: number) {
    const article = await this.articleService.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) {
      throw new NotFoundException('Article not found');
    }
    return article;
  }
}
