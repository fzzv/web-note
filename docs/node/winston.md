# Winston 日志库

Winston 是 Node.js 中最流行的日志库之一，提供了灵活的日志记录功能，支持多种传输方式、日志级别和格式化选项。

## 基本概念

Winston 的核心概念包括：
- **Logger（记录器）**：日志记录的主要接口
- **Transport（传输）**：日志输出的目标（控制台、文件等）
- **Format（格式）**：日志的格式化方式
- **Level（级别）**：日志的重要程度

## 安装依赖

```bash
npm install winston
```

## 基础使用

### 1. 简单日志记录器

```javascript
const winston = require('winston');
const path = require('path');

// 创建一个日志记录器的实例
const logger = winston.createLogger({
  // 设置日志的级别
  level: 'info',
  // 设置日志的格式
  format: winston.format.json(),
  // 配置日志的传输方式
  transports: [
    new winston.transports.Console(),    // 写入控制台
    new winston.transports.File({        // 写入文件
      filename: path.join(__dirname, './combined.log')
    })
  ]
});

logger.info('这是一条info日志');
```

**特点**：
- 同时输出到控制台和文件
- 使用 JSON 格式
- 设置最低日志级别为 info

### 2. 日志级别详解

Winston 支持以下日志级别（从高到低）：

```javascript
const winston = require('winston');
const path = require('path');

const logger = winston.createLogger({
  level: 'info',
  format: winston.format.json(),
  transports: [
    new winston.transports.Console(),
    new winston.transports.File({
      filename: path.join(__dirname, './combined.log')
    })
  ]
});

// 不同级别的日志
logger.error('这是一条error日志');     // 错误级别 - 优先级 0
logger.warn('这是一条warn日志');       // 警告级别 - 优先级 1
logger.info('这是一条info日志');       // 消息级别 - 优先级 2
logger.debug('这是一条debug日志');     // 详细级别 - 优先级 3
logger.verbose('这是一条verbose日志'); // 调试级别 - 优先级 4
logger.silly('这是一条silly日志');     // 无意义级别 - 优先级 5
```

**日志级别说明**：
- `error`：错误信息，需要立即关注
- `warn`：警告信息，可能出现问题
- `info`：一般信息，正常运行状态
- `debug`：调试信息，开发阶段使用
- `verbose`：详细信息，更多调试细节
- `silly`：最详细信息，通常不在生产环境使用

### 3. 自定义格式化

```javascript
const winston = require('winston');
const path = require('path');
const { combine, timestamp, printf } = winston.format;

// 自定义格式化函数
const myFormat = printf(({ level, message, timestamp }) => {
  return `${timestamp} ${level} ${message}`;
});

const logger = winston.createLogger({
  level: 'info',
  // 结合时间戳和自定义格式化函数
  format: combine(
    timestamp({ format: 'YYYY-MM-DD HH:mm:ss' }), // 添加时间戳
    myFormat // 应用自定义格式
  ),
  transports: [
    new winston.transports.Console(),
    new winston.transports.File({
      filename: path.join(__dirname, './combined.log')
    }),
    // 错误日志单独文件
    new winston.transports.File({
      filename: path.join(__dirname, './error.log'), 
      level: 'error'
    })
  ]
});

logger.info('这是一条info日志');
```

**输出格式**：
```
2024-01-20 10:30:45 info 这是一条info日志
```

### 4. 日志轮转（Daily Rotate）

```javascript
const winston = require('winston');
const path = require('path');
require('winston-daily-rotate-file');

const logger = winston.createLogger({
  transports: [
    new winston.transports.DailyRotateFile({
      // 指定日志文件的文件名模式
      filename: "app-%DATE%.log",
      // 指定文件的目录
      dirname: path.join(__dirname, './logs'),
      // 指定日期的格式
      datePattern: 'YYYY-MM-DD',
      // 指定日志的级别
      level: 'info',
      // 设置日志的最大的文件大小
      maxSize: '1k',
      // 设置日志文件的最大保留天数
      maxFiles: '14d', // 14天后会自动删除
      // 指定是否压缩旧的日志文件
      zippedArchive: true
    })
  ]
});

logger.info('这是一条测试日志');
```

**特点**：
- 按日期自动创建新文件（app-2024-01-20.log）
- 文件大小超过限制时自动轮转
- 自动删除过期日志文件
- 可选择压缩旧文件节省空间

## 高级配置

### 1. 完整的 Logger 配置

```javascript
const winston = require('winston');
const path = require('path');

const logger = winston.createLogger({
  // 日志级别
  level: process.env.LOG_LEVEL || 'info',
  
  // 默认元数据
  defaultMeta: {
    service: 'user-service',
    version: '1.0.0'
  },
  
  // 格式化
  format: winston.format.combine(
    winston.format.timestamp({
      format: 'YYYY-MM-DD HH:mm:ss'
    }),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  
  // 传输方式
  transports: [
    // 错误日志
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/error.log'),
      level: 'error',
      maxsize: 5242880, // 5MB
      maxFiles: 5
    }),
    
    // 组合日志
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/combined.log'),
      maxsize: 5242880,
      maxFiles: 5
    })
  ],
  
  // 异常处理
  exceptionHandlers: [
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/exceptions.log')
    })
  ],
  
  // Promise 拒绝处理
  rejectionHandlers: [
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/rejections.log')
    })
  ]
});

// 开发环境添加控制台输出
if (process.env.NODE_ENV !== 'production') {
  logger.add(new winston.transports.Console({
    format: winston.format.combine(
      winston.format.colorize(),
      winston.format.simple()
    )
  }));
}
```

### 2. 多种格式化选项

```javascript
const winston = require('winston');

// JSON 格式
const jsonLogger = winston.createLogger({
  format: winston.format.json(),
  transports: [new winston.transports.Console()]
});

// 简单格式
const simpleLogger = winston.createLogger({
  format: winston.format.simple(),
  transports: [new winston.transports.Console()]
});

// 带颜色的格式
const colorLogger = winston.createLogger({
  format: winston.format.combine(
    winston.format.colorize(),
    winston.format.simple()
  ),
  transports: [new winston.transports.Console()]
});

// 自定义格式
const customLogger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.printf(({ timestamp, level, message, ...meta }) => {
      return `[${timestamp}] ${level}: ${message} ${
        Object.keys(meta).length ? JSON.stringify(meta) : ''
      }`;
    })
  ),
  transports: [new winston.transports.Console()]
});

// 预定义格式
const prettyLogger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.align(),
    winston.format.printf(info => `${info.timestamp} ${info.level}: ${info.message}`)
  ),
  transports: [new winston.transports.Console()]
});
```

### 3. 条件日志记录

```javascript
const winston = require('winston');

const logger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    // 过滤敏感信息
    winston.format.printf(({ timestamp, level, message, ...meta }) => {
      // 移除密码字段
      if (meta.password) {
        meta.password = '***';
      }
      return `[${timestamp}] ${level}: ${message} ${JSON.stringify(meta)}`;
    })
  ),
  transports: [
    new winston.transports.Console(),
    // 只记录 error 级别到文件
    new winston.transports.File({
      filename: 'error.log',
      level: 'error'
    })
  ]
});

// 带元数据的日志
logger.info('用户登录', {
  userId: 123,
  username: 'john',
  password: 'secret123', // 会被过滤
  ip: '192.168.1.1'
});
```

## 实际应用场景

### 1. Express 应用日志

```javascript
const express = require('express');
const winston = require('winston');
const path = require('path');

// 创建日志记录器
const logger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.errors({ stack: true }),
    winston.format.json()
  ),
  defaultMeta: { service: 'web-app' },
  transports: [
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/error.log'),
      level: 'error'
    }),
    new winston.transports.File({
      filename: path.join(__dirname, 'logs/combined.log')
    })
  ]
});

if (process.env.NODE_ENV !== 'production') {
  logger.add(new winston.transports.Console({
    format: winston.format.simple()
  }));
}

const app = express();

// 请求日志中间件
app.use((req, res, next) => {
  logger.info('HTTP Request', {
    method: req.method,
    url: req.url,
    ip: req.ip,
    userAgent: req.get('User-Agent')
  });
  next();
});

// 错误处理中间件
app.use((err, req, res, next) => {
  logger.error('应用错误', {
    error: err.message,
    stack: err.stack,
    url: req.url,
    method: req.method
  });
  
  res.status(500).json({ error: '服务器内部错误' });
});

// 路由示例
app.get('/users/:id', async (req, res) => {
  try {
    const userId = req.params.id;
    logger.info('获取用户信息', { userId });
    
    // 模拟数据库查询
    const user = await getUserById(userId);
    
    if (!user) {
      logger.warn('用户未找到', { userId });
      return res.status(404).json({ error: '用户未找到' });
    }
    
    logger.info('用户信息获取成功', { userId, username: user.username });
    res.json(user);
  } catch (error) {
    logger.error('获取用户信息失败', {
      userId: req.params.id,
      error: error.message
    });
    res.status(500).json({ error: '服务器错误' });
  }
});

app.listen(3000, () => {
  logger.info('服务器启动', { port: 3000 });
});
```

### 2. 数据库操作日志

```javascript
const winston = require('winston');

const dbLogger = winston.createLogger({
  level: 'debug',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.printf(({ timestamp, level, message, ...meta }) => {
      return `[${timestamp}] ${level}: ${message} ${
        meta.sql ? `\nSQL: ${meta.sql}` : ''
      } ${
        meta.duration ? `\nDuration: ${meta.duration}ms` : ''
      }`;
    })
  ),
  transports: [
    new winston.transports.File({
      filename: 'logs/database.log'
    })
  ]
});

// 数据库查询包装函数
async function executeQuery(sql, params = []) {
  const startTime = Date.now();
  
  try {
    dbLogger.debug('执行数据库查询', { sql, params });
    
    // 模拟数据库查询
    const result = await database.query(sql, params);
    const duration = Date.now() - startTime;
    
    dbLogger.info('查询执行成功', {
      sql,
      rowCount: result.rowCount,
      duration
    });
    
    return result;
  } catch (error) {
    const duration = Date.now() - startTime;
    
    dbLogger.error('查询执行失败', {
      sql,
      params,
      error: error.message,
      duration
    });
    
    throw error;
  }
}
```

### 3. 业务逻辑日志

```javascript
const winston = require('winston');

const businessLogger = winston.createLogger({
  level: 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  defaultMeta: { module: 'business' },
  transports: [
    new winston.transports.File({
      filename: 'logs/business.log'
    })
  ]
});

class UserService {
  async createUser(userData) {
    const correlationId = generateCorrelationId();
    
    businessLogger.info('开始创建用户', {
      correlationId,
      email: userData.email,
      operation: 'createUser'
    });
    
    try {
      // 验证用户数据
      await this.validateUserData(userData);
      businessLogger.debug('用户数据验证通过', { correlationId });
      
      // 检查邮箱是否已存在
      const existingUser = await this.findUserByEmail(userData.email);
      if (existingUser) {
        businessLogger.warn('邮箱已存在', {
          correlationId,
          email: userData.email
        });
        throw new Error('邮箱已存在');
      }
      
      // 创建用户
      const user = await this.saveUser(userData);
      
      businessLogger.info('用户创建成功', {
        correlationId,
        userId: user.id,
        email: user.email,
        operation: 'createUser'
      });
      
      // 发送欢迎邮件
      await this.sendWelcomeEmail(user);
      businessLogger.info('欢迎邮件发送成功', { correlationId, userId: user.id });
      
      return user;
    } catch (error) {
      businessLogger.error('用户创建失败', {
        correlationId,
        email: userData.email,
        error: error.message,
        operation: 'createUser'
      });
      throw error;
    }
  }
}

function generateCorrelationId() {
  return Math.random().toString(36).substr(2, 9);
}
```

## 性能优化

### 1. 异步日志记录

```javascript
const winston = require('winston');

// 使用异步传输
const logger = winston.createLogger({
  format: winston.format.json(),
  transports: [
    new winston.transports.File({
      filename: 'logs/app.log',
      // 启用异步写入
      options: { flags: 'a' }
    })
  ]
});

// 批量日志记录
class BatchLogger {
  constructor() {
    this.logs = [];
    this.batchSize = 100;
    this.flushInterval = 5000; // 5秒
    
    setInterval(() => {
      this.flush();
    }, this.flushInterval);
  }
  
  log(level, message, meta = {}) {
    this.logs.push({
      timestamp: new Date().toISOString(),
      level,
      message,
      ...meta
    });
    
    if (this.logs.length >= this.batchSize) {
      this.flush();
    }
  }
  
  flush() {
    if (this.logs.length === 0) return;
    
    const logsToFlush = [...this.logs];
    this.logs = [];
    
    logsToFlush.forEach(log => {
      logger.log(log.level, log.message, log);
    });
  }
}

const batchLogger = new BatchLogger();
```

### 2. 条件日志和采样

```javascript
const winston = require('winston');

// 采样日志记录器
class SamplingLogger {
  constructor(baseLogger, sampleRate = 0.1) {
    this.baseLogger = baseLogger;
    this.sampleRate = sampleRate;
  }
  
  log(level, message, meta = {}) {
    // 错误日志总是记录
    if (level === 'error' || Math.random() < this.sampleRate) {
      this.baseLogger.log(level, message, meta);
    }
  }
  
  info(message, meta) { this.log('info', message, meta); }
  warn(message, meta) { this.log('warn', message, meta); }
  error(message, meta) { this.log('error', message, meta); }
  debug(message, meta) { this.log('debug', message, meta); }
}

const baseLogger = winston.createLogger({
  transports: [new winston.transports.File({ filename: 'app.log' })]
});

// 10% 采样率
const samplingLogger = new SamplingLogger(baseLogger, 0.1);
```

### 3. 内存使用优化

```javascript
const winston = require('winston');

// 限制日志对象大小
const createSafeLogger = () => {
  const MAX_LOG_SIZE = 1000; // 最大字符数
  
  return winston.createLogger({
    format: winston.format.combine(
      winston.format.timestamp(),
      winston.format.printf(({ timestamp, level, message, ...meta }) => {
        let logString = `[${timestamp}] ${level}: ${message}`;
        
        if (Object.keys(meta).length > 0) {
          let metaString = JSON.stringify(meta);
          if (metaString.length > MAX_LOG_SIZE) {
            metaString = metaString.substring(0, MAX_LOG_SIZE) + '...[truncated]';
          }
          logString += ` ${metaString}`;
        }
        
        return logString;
      })
    ),
    transports: [
      new winston.transports.File({ filename: 'app.log' })
    ]
  });
};

const safeLogger = createSafeLogger();
```

## 与其他框架集成

### 1. Express + Morgan 集成

```javascript
const express = require('express');
const morgan = require('morgan');
const winston = require('winston');

const logger = winston.createLogger({
  format: winston.format.json(),
  transports: [
    new winston.transports.File({ filename: 'access.log' })
  ]
});

const app = express();

// Morgan 流式输出到 Winston
app.use(morgan('combined', {
  stream: {
    write: (message) => {
      logger.info(message.trim());
    }
  }
}));
```

### 2. NestJS 集成

```javascript
// winston.service.ts
import { Injectable, LoggerService } from '@nestjs/common';
import * as winston from 'winston';

@Injectable()
export class WinstonLoggerService implements LoggerService {
  private logger: winston.Logger;

  constructor() {
    this.logger = winston.createLogger({
      level: 'info',
      format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json()
      ),
      transports: [
        new winston.transports.File({ filename: 'logs/app.log' }),
        new winston.transports.Console()
      ]
    });
  }

  log(message: string, context?: string) {
    this.logger.info(message, { context });
  }

  error(message: string, trace?: string, context?: string) {
    this.logger.error(message, { trace, context });
  }

  warn(message: string, context?: string) {
    this.logger.warn(message, { context });
  }

  debug(message: string, context?: string) {
    this.logger.debug(message, { context });
  }

  verbose(message: string, context?: string) {
    this.logger.verbose(message, { context });
  }
}

// main.ts
import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { WinstonLoggerService } from './winston.service';

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    logger: new WinstonLoggerService()
  });
  await app.listen(3000);
}
bootstrap();
```

## 最佳实践

### 1. 日志结构化

```javascript
const winston = require('winston');

// 结构化日志格式
const structuredLogger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  defaultMeta: {
    service: 'user-service',
    version: '1.0.0'
  },
  transports: [
    new winston.transports.File({ filename: 'structured.log' })
  ]
});

// 标准化日志记录函数
function logUserAction(action, userId, details = {}) {
  structuredLogger.info('用户操作', {
    action,
    userId,
    timestamp: new Date().toISOString(),
    ...details
  });
}

function logError(error, context = {}) {
  structuredLogger.error('应用错误', {
    message: error.message,
    stack: error.stack,
    timestamp: new Date().toISOString(),
    ...context
  });
}

// 使用示例
logUserAction('login', 123, { ip: '192.168.1.1' });
logUserAction('purchase', 123, { productId: 456, amount: 99.99 });
```

### 2. 环境配置

```javascript
const winston = require('winston');

// 环境相关配置
const getLogConfig = () => {
  const env = process.env.NODE_ENV || 'development';
  
  const baseConfig = {
    level: process.env.LOG_LEVEL || 'info',
    format: winston.format.combine(
      winston.format.timestamp(),
      winston.format.errors({ stack: true })
    ),
    defaultMeta: {
      service: process.env.SERVICE_NAME || 'app',
      environment: env
    }
  };

  if (env === 'production') {
    return {
      ...baseConfig,
      format: winston.format.combine(
        baseConfig.format,
        winston.format.json()
      ),
      transports: [
        new winston.transports.File({
          filename: 'logs/error.log',
          level: 'error'
        }),
        new winston.transports.File({
          filename: 'logs/combined.log'
        })
      ]
    };
  } else {
    return {
      ...baseConfig,
      transports: [
        new winston.transports.Console({
          format: winston.format.combine(
            winston.format.colorize(),
            winston.format.simple()
          )
        })
      ]
    };
  }
};

const logger = winston.createLogger(getLogConfig());
```

### 3. 日志监控和告警

```javascript
const winston = require('winston');

// 自定义传输用于监控
class AlertTransport extends winston.Transport {
  constructor(opts) {
    super(opts);
    this.name = 'alert';
    this.level = opts.level || 'error';
  }

  log(info, callback) {
    // 发送告警
    if (info.level === 'error') {
      this.sendAlert(info);
    }
    
    callback();
  }

  sendAlert(logInfo) {
    // 实现告警逻辑（邮件、短信、Slack等）
    console.log('🚨 告警:', logInfo.message);
  }
}

const logger = winston.createLogger({
  transports: [
    new winston.transports.File({ filename: 'app.log' }),
    new AlertTransport({ level: 'error' })
  ]
});

// 错误计数和频率限制
class ErrorCounter {
  constructor() {
    this.errors = new Map();
    this.resetInterval = 60000; // 1分钟
    
    setInterval(() => {
      this.errors.clear();
    }, this.resetInterval);
  }
  
  shouldAlert(errorKey) {
    const count = this.errors.get(errorKey) || 0;
    this.errors.set(errorKey, count + 1);
    
    // 相同错误1分钟内只告警一次
    return count === 0;
  }
}

const errorCounter = new ErrorCounter();

function logErrorWithAlert(error, context = {}) {
  const errorKey = `${error.message}-${context.module || 'unknown'}`;
  
  logger.error('应用错误', {
    message: error.message,
    stack: error.stack,
    ...context
  });
  
  if (errorCounter.shouldAlert(errorKey)) {
    // 发送告警
    logger.error('需要告警的错误', {
      message: error.message,
      context,
      alertRequired: true
    });
  }
}
```

### 4. 日志分析和搜索

```javascript
const winston = require('winston');

// 带有搜索标签的日志
const searchableLogger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.printf(({ timestamp, level, message, tags, ...meta }) => {
      const tagString = tags ? `[${tags.join(',')}]` : '';
      return `[${timestamp}] ${level}: ${tagString} ${message} ${JSON.stringify(meta)}`;
    })
  ),
  transports: [
    new winston.transports.File({ filename: 'searchable.log' })
  ]
});

// 使用标签进行分类
function logWithTags(level, message, tags = [], meta = {}) {
  searchableLogger.log(level, message, { tags, ...meta });
}

// 使用示例
logWithTags('info', '用户登录', ['auth', 'user'], { userId: 123 });
logWithTags('error', '支付失败', ['payment', 'error'], { orderId: 456 });
logWithTags('info', '数据同步完成', ['sync', 'database'], { recordCount: 1000 });
```

## 总结

Winston 是一个功能强大且灵活的 Node.js 日志库，具有以下特点：

1. **多传输支持**：控制台、文件、数据库等多种输出方式
2. **灵活格式化**：支持 JSON、自定义格式等多种格式
3. **日志级别管理**：完整的日志级别系统
4. **高性能**：支持异步写入和批量处理
5. **可扩展性**：丰富的插件生态系统
6. **生产就绪**：支持日志轮转、异常处理等企业级功能

在实际项目中，合理使用 Winston 的各种特性，可以构建出完善的日志系统，提高应用的可观测性和可维护性。
