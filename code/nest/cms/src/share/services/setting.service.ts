import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from 'mongoose';
import { SettingDocument } from '../schemas/setting.schema';
import { CreateSettingDto, UpdateSettingDto } from '../dtos/setting.dto';
import { MongoDBBaseService } from './mongodb-base.service';

@Injectable()
export class SettingService extends MongoDBBaseService<SettingDocument, CreateSettingDto, UpdateSettingDto> {
  constructor(@InjectModel('Setting') settingModel: Model<SettingDocument>) {
    super(settingModel);
  }
  async findFirst(): Promise<SettingDocument | null> {
    return await this.model.findOne().exec();
  }
}
