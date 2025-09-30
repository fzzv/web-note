// 导入 Injectable 装饰器，表示该类可以被依赖注入
import { Injectable } from '@nestjs/common';
// 导入 systeminformation 模块，用于获取系统信息
import si from 'systeminformation';

// 使用 Injectable 装饰器，将此类标记为可注入服务
@Injectable()
export class SystemService {
  // 私有方法，将字节值格式化为GB，保留两位小数
  private formatToGB(value: number) {
    return (value / (1024 ** 3)).toFixed(2);
  }

  // 异步方法，获取系统信息
  async getSystemInfo() {
    // 获取当前CPU负载信息
    const cpu = await si.currentLoad();
    // 获取内存使用情况
    const memory = await si.mem();
    // 获取磁盘使用情况
    const disk = await si.fsSize();
    // 获取操作系统信息
    const osInfo = await si.osInfo();
    // 获取网络接口信息
    const networkInterfaces = await si.networkInterfaces();

    // 返回格式化后的系统信息
    return {
      // CPU信息
      cpu: {
        // CPU核心数
        cores: cpu.cpus.length,
        // 用户进程占用的CPU负载百分比
        userLoad: cpu.currentLoadUser.toFixed(2),
        // 系统进程占用的CPU负载百分比
        systemLoad: cpu.currentLoadSystem.toFixed(2),
        // 空闲的CPU负载百分比
        idle: cpu.currentLoadIdle.toFixed(2),
      },
      // 内存信息
      memory: {
        // 总内存，单位GB
        total: this.formatToGB(memory.total),
        // 已使用内存，单位GB
        used: this.formatToGB(memory.used),
        // 空闲内存，单位GB
        free: this.formatToGB(memory.free),
        // 内存使用率，单位百分比
        usage: ((memory.used / memory.total) * 100).toFixed(2),
      },
      // 磁盘信息
      disks: disk.map(d => ({
        // 挂载点
        mount: d.mount,
        // 文件系统类型
        filesystem: d.fs,
        // 磁盘类型
        type: d.type,
        // 磁盘总大小，单位GB
        size: this.formatToGB(d.size),
        // 已使用空间，单位GB
        used: this.formatToGB(d.used),
        // 可用空间，单位GB
        available: this.formatToGB(d.available),
        // 磁盘使用率，单位百分比
        usage: d.use.toFixed(2),
      })),
      // 服务器信息
      server: {
        // 主机名
        hostname: osInfo.hostname,
        // IP地址，若无网络接口则返回 'N/A'
        ip: networkInterfaces[0]?.ip4 || 'N/A',
        // 操作系统发行版
        os: osInfo.distro,
        // 系统架构类型
        arch: osInfo.arch,
      }
    };
  }
}
