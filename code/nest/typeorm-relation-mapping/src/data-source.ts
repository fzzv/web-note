import "reflect-metadata"
import { DataSource } from "typeorm"
import { User } from "./entity/User"
import { IdCard } from "./entity/IdCard"
import {Department} from "./entity/Department";
import {Employee} from "./entity/Employee";
import {Article} from "./entity/Article";
import {Tag} from "./entity/Tag";

export const AppDataSource = new DataSource({
    type: "mysql",
    host: "localhost",
    port: 3306,
    username: "root",
    password: "password",
    database: "typeorm_relation_mapping",
    synchronize: true, // 根据同步建表，也就是当 database 里没有和 Entity 对应的表的时候，会自动生成建表 sql 语句并执行
    logging: true, // 打印生成的 sql 语句
    entities: [User, IdCard, Department, Employee, Article, Tag], // 指定有哪些和数据库的表对应的 Entity
    migrations: [], // 修改表结构之类的 sql
    subscribers: [], // 一些 Entity 生命周期的订阅者，比如 insert、update、remove 前后，可以加入一些逻辑
    poolSize: 10, // 指定数据库连接池中连接的最大数量
    connectorPackage: 'mysql2', // 驱动包
    extra: { // 额外发送给驱动包的一些选项
        authPlugin: 'sha256_password',
    }
})
