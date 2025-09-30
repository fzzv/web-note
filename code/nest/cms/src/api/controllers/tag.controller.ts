import { Controller, Get, Query } from '@nestjs/common';
import { TagService } from '../../share/services/tag.service';

@Controller('api/tags')
export class TagController {
  constructor(private readonly tagService: TagService) { }

  @Get()
  async getTags(@Query('selectedTag') selectedTag: string = '') {
    const tags = await this.tagService.findAll();
    return { tags, selectedTag };
  }
}
