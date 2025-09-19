import { Injectable, NotFoundException } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Article } from "../entities/article.entity";
import { Repository, Like, In, UpdateResult } from 'typeorm';
import { MysqlBaseService } from "./mysql-base.service";
import { CreateArticleDto, UpdateArticleDto } from "../dtos/article.dto";
import { Category } from "../entities/category.entity";
import { Tag } from "../entities/tag.entity";

@Injectable()
export class ArticleService extends MysqlBaseService<Article> {
  constructor(
    @InjectRepository(Article) protected repository: Repository<Article>,
    @InjectRepository(Category) private readonly categoryRepository: Repository<Category>,
    @InjectRepository(Tag) private readonly tagRepository: Repository<Tag>
  ) {
    super(repository);
  }

  async findAll(keyword?: string) {
    const where = keyword ? [
      { title: Like(`%${keyword}%`) }
    ] : {};
    return this.repository.find({ where, relations: ['categories', 'tags'] });
  }

  async findAllWithPagination(page: number, limit: number, keyword?: string) {
    const where = keyword ? [
      { title: Like(`%${keyword}%`) }
    ] : {};
    const [articles, total] = await this.repository.findAndCount({
      where,
      relations: ['categories', 'tags'],
      skip: (page - 1) * limit,
      take: limit
    });
    return { articles, total };
  }

  async create(createArticleDto: CreateArticleDto) {
    const { categoryIds = [], tagIds = [], ...articleData } = createArticleDto;
    const article = this.repository.create(articleData);
    article.categories = await this.categoryRepository.findBy({ id: In(categoryIds) });
    article.tags = await this.tagRepository.findBy({ id: In(tagIds) });
    return await this.repository.save(article);
  }

  async update(id: number, updateArticleDto: UpdateArticleDto) {
    const { categoryIds, tagIds, ...articleData } = updateArticleDto;
    const article = await this.repository.findOne({ where: { id }, relations: ['categories', 'tags'] });
    if (!article) throw new NotFoundException('Article not found');
    Object.assign(article, articleData);
    if (categoryIds) {
      article.categories = await this.categoryRepository.findBy({ id: In(categoryIds) });
    }
    if (tagIds) {
      article.tags = await this.tagRepository.findBy({ id: In(tagIds) });
    }
    await this.repository.update(id, article);
    return UpdateResult.from({ raw: [], affected: 1, records: [] });
  }
}
