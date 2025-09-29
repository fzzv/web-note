// 导入 Injectable 装饰器，用于将服务类标记为可注入的依赖
import { Injectable } from '@nestjs/common';
// 导入 ExcelJS 库，用于创建和操作 Excel 文件
import * as ExcelJS from 'exceljs';
@Injectable()
export class ExcelExportService {
  // 异步方法，用于将数据导出为 Excel 文件
  async exportAsExcel(data: any[], columns: { header: string, key: string, width: number }[]) {
    // 创建一个新的 Excel 工作簿
    const workbook = new ExcelJS.Workbook();
    // 添加一个新的工作表，并命名为 'Data'
    const worksheet = workbook.addWorksheet('Data');
    // 设置工作表的列，根据传入的列定义数组
    worksheet.columns = columns;
    // 遍历数据数组，将每一项数据作为一行添加到工作表中
    data.forEach(item => {
      worksheet.addRow(item);
    });
    // 将工作簿内容写入缓冲区，并返回该缓冲区（用于进一步处理或保存）
    return workbook.xlsx.writeBuffer();
  }
}
