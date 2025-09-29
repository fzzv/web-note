// 导入 Injectable 装饰器，用于标记一个服务类
import { Injectable } from '@nestjs/common';
// 导入 PptxGenJS 库，用于生成 PPTX 文件
import PptxGenJS from 'pptxgenjs';
// 导入 html-pptxgenjs 库，用于将 HTML 转换为 PPTX 内容
import * as html2ppt from 'html-pptxgenjs';
// 使用 Injectable 装饰器将 PptExportService 标记为可注入的服务
@Injectable()
export class PptExportService {
  // 异步方法，用于将文章列表导出为 PPTX 文件
  async exportToPpt(articles: any[]) {
    // 创建一个新的 PPTX 对象
    const pptx = new (PptxGenJS as any)();
    // 遍历每篇文章，将其内容添加到 PPTX 幻灯片中
    for (const article of articles) {
      // 添加一个新的幻灯片到 PPTX
      const slide = pptx.addSlide();
      // 构建 HTML 内容，包含文章标题、状态、分类、标签和正文内容
      const htmlContent = `
                <h1>${article.title}</h1>
                <p><strong>状态:</strong> ${article.state}</p>
                <p><strong>分类:</strong> ${article.categories.map(c => c.name).join(', ')}</p>
                <p><strong>标签:</strong> ${article.tags.map(t => t.name).join(', ')}</p>
                <hr/>
                ${article.content}
            `;
      // 使用 html-pptxgenjs 将 HTML 内容转换为 PPTX 可用的文本项
      const items = html2ppt.htmlToPptxText(htmlContent);
      // 将生成的文本项添加到幻灯片中，设置其位置和大小
      slide.addText(items, { x: 0.5, y: 0.5, w: 9.5, h: 6, valign: 'top' });
    }
    // 将生成的 PPTX 文件以 nodebuffer 的形式输出
    return await pptx.write({ outputType: 'nodebuffer' });
  }
}
