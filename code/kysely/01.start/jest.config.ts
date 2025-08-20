import type { Config } from 'jest'

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  moduleFileExtensions: ['js', 'json', 'ts'],
  rootDir: '.',
  testRegex: '.*\\.test\\.ts$',
  coverageDirectory: './coverage',
  collectCoverageFrom: [
    'src/**/*.(t|j)s',
    '!src/main.ts', // 跳过入口文件
  ],
}

export default config
