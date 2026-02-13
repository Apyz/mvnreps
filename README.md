# mvnreps - Maven 仓库搜索工具

mvnreps 命令行工具，用于快速搜索 Maven 中央仓库中的依赖库。

## 功能特性

- 模糊搜索 Maven 仓库中的库
- 精确查询指定 groupId:artifactId 的历史版本
- 支持多种输出格式（default、pom、gradle、gradle.kts、format）
- 可自定义返回结果数量

## 安装

### 下载

下载 `mvnreps.exe` 二进制文件。

### 使用

将 `mvnreps.exe` 放到任意目录，或添加到系统 PATH 环境变量中以便全局使用。

## 使用方法

### 基本语法

```bash
mvnreps <query> [format]
```

### 搜索模式

#### 1. 模糊搜索

搜索包含关键词的 Maven 库：

```bash
mvnreps okhttp
```

输出示例：
```
🔍 库搜索结果 ("okhttp", 前 5 条):

com.squareup.okhttp3:okhttp:5.0.0-alpha.16
com.avito.android:okhttp:2024.32
io.github.qsy7.java.dependencies:okhttp:0.3.3
io.github.sunny-chung:okhttp:4.11.0-patch-1
com.lightningkite.rx:okhttp:1.0.7
```

#### 2. 精确查询

查看指定库的历史版本（按时间倒序）：

```bash
mvnreps com.squareup.okhttp3:okhttp
```

输出示例：
```
🔍 历史版本结果 ("com.squareup.okhttp3:okhttp", 前 5 条):

com.squareup.okhttp3:okhttp:5.0.0-alpha.16
com.squareup.okhttp3:okhttp:5.0.0-alpha.15
com.squareup.okhttp3:okhttp:5.0.0-alpha.14
com.squareup.okhttp3:okhttp:5.0.0-alpha.13
com.squareup.okhttp3:okhttp:5.0.0-alpha.12
```

### 指定返回数量

在查询参数后添加 `,N` 来指定返回结果数量：

```bash
mvnreps okhttp,10
mvnreps com.squareup.okhttp3:okhttp,3
```

### 输出格式

支持五种输出格式，在查询后添加格式参数：

#### default 格式（默认）

```bash
mvnreps okhttp
```

输出：
```
com.squareup.okhttp3:okhttp:5.0.0-alpha.16
```

#### pom 格式

Maven XML 依赖格式：

```bash
mvnreps com.squareup.okhttp3:okhttp,3 pom
```

输出：
```xml
<dependency>
  <groupId>com.squareup.okhttp3</groupId>
  <artifactId>okhttp</artifactId>
  <version>5.0.0-alpha.16</version>
</dependency>
```

#### gradle 格式

Gradle 依赖格式：

```bash
mvnreps com.squareup.okhttp3:okhttp,10 gradle
```

输出：
```gradle
implementation 'com.squareup.okhttp3:okhttp:5.0.0-alpha.16'
```

#### gradle.kts 格式

Kotlin DSL 依赖格式：

```bash
mvnreps com.squareup.okhttp3:okhttp,10 gradle.kts
```

输出：
```kotlin
implementation("com.squareup.okhttp3:okhttp:5.0.0-alpha.16")
```

#### format 格式

表格格式，便于阅读：

```bash
mvnreps okhttp,3 format
```

输出：
```
GroupId                                  | ArtifactId                     | Version
-------------------------------------------------------------------------------------
com.squareup.okhttp3                     | okhttp                         | 5.0.0-alpha.16
com.avito.android                        | okhttp                         | 2024.32
io.github.qsy7.java.dependencies         | okhttp                         | 0.3.3
```

## 使用示例

### 查找 Spring Boot 相关库

```bash
mvnreps spring-boot
```

### 查找特定版本的 Spring Boot

```bash
mvnreps org.springframework.boot:spring-boot,5 gradle
```

### 查找 Gson 库并输出 Kotlin DSL 格式

```bash
mvnreps com.google.code.gson:gson gradle.kts
```

### 查找常用工具库

```bash
mvnreps lombok
mvnreps commons-lang3
mvnreps guava
```
