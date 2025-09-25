import { Injectable } from "@nestjs/common";
import { InjectRepository } from "@nestjs/typeorm";
import { Category } from "../entities/category.entity";
import { MysqlBaseService } from "./mysql-base.service";
import { CreateCategoryDto, UpdateCategoryDto } from '../dtos/category.dto';
import { TreeRepository, UpdateResult } from 'typeorm';

@Injectable()
export class CategoryService extends MysqlBaseService<Category> {
  constructor(
    @InjectRepository(Category) protected repository: TreeRepository<Category>
  ) {
    super(repository);
  }

  async findAllList() {
    const categories = await this.repository.find({ relations: ['children', 'parent'] });
    return categories;
  }

  async findAll() {
    const categoryTree = await this.repository.findTrees({ relations: ['children', 'parent'] });
    return categoryTree.filter(category => !category.parent);
  }

  async create(createCategoryDto: CreateCategoryDto): Promise<Category> {
    const { parentId, ...dto } = createCategoryDto;
    const category = this.repository.create(dto);
    if (parentId) {
      const parent = await this.repository.findOneBy({ id: parentId });
      if (!parent) throw new Error('Parent category not found');
      category.parent = parent;
    }
    await this.repository.save(category);
    return this.findOne({ where: { id: category.id } }) as Promise<Category>;
  }

  async update(id: number, updateCategoryDto: UpdateCategoryDto) {
    const { id: _ignoreId, parentId, ...dto } = updateCategoryDto;
    const category = await this.repository.findOneBy({ id });
    console.log(category, 'categoty');
    if (!category) throw new Error('Category not found');
    Object.assign(category, dto);
    if (parentId) {
      const parent = await this.repository.findOneBy({ id: parentId });
      if (!parent) throw new Error('Parent category not found');
      category.parent = parent;
    }
    await this.repository.save(category);
    return UpdateResult.from({ raw: [], affected: 1, records: [] });
  }
}
