import { Controller, Get, Query } from '@nestjs/common';
import { CategoryService } from '../../share/services/category.service';

@Controller('api/categories')
export class CategoryController {
  constructor(private readonly categoryService: CategoryService) { }

  @Get()
  async getCategories(@Query('selectedCategory') selectedCategory: string = '') {
    const categories = await this.categoryService.findList();
    return { categories, selectedCategory };
  }
}
