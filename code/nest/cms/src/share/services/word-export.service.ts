import { Injectable } from '@nestjs/common';
import htmlToDocx from 'html-to-docx';

@Injectable()
export class WordExportService {
  async exportToWord(htmlContent: string): Promise<Buffer> {
    return await htmlToDocx(htmlContent);
  }
}
