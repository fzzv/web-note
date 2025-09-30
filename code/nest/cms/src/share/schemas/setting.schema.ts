import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { HydratedDocument } from 'mongoose';
export type SettingDocument = HydratedDocument<Setting>;
@Schema()
export class Setting {
  id: string;
  @Prop({ required: true })
  siteName: string;
  @Prop()
  siteDescription: string;
  @Prop()
  contactEmail: string;
}
export const SettingSchema = SchemaFactory.createForClass(Setting);
SettingSchema.virtual('id').get(function () {
  return this._id.toHexString();
});
SettingSchema.set('toJSON', { virtuals: true });
SettingSchema.set('toObject', { virtuals: true });
