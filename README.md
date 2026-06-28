# 社区诊所预约管理系统 - 后端 API

Go + Gin + GORM + MySQL 实现的社区诊所预约管理后端服务。

## 快速启动

### 1. 启动 MySQL

```bash
docker-compose up -d
```

等待 MySQL 就绪（约 10-20 秒），可通过以下命令确认：

```bash
docker-compose logs -f mysql
# 看到 "ready for connections" 即可
```

### 2. 启动服务

```bash
go run main.go
```

默认监听 `http://localhost:8080`，启动时自动建表。

### 3. 验证服务

```bash
curl http://localhost:8080/ping
# 返回 {"message":"pong"}
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_HOST | 127.0.0.1 | MySQL 地址 |
| DB_PORT | 3306 | MySQL 端口 |
| DB_USER | clinic | MySQL 用户 |
| DB_PASSWORD | clinic123 | MySQL 密码 |
| DB_NAME | clinic | 数据库名 |
| SERVER_PORT | 8080 | 服务端口 |
| AUTH_TOKEN | clinic-secret-token-2024 | API Token |

## 认证方式

所有业务接口需在请求头携带 Token：

```
Authorization: Bearer clinic-secret-token-2024
```

## API 接口

### 统一响应格式

```json
{
  "code": 0,
  "msg": "成功",
  "data": {}
}
```

错误响应：

```json
{
  "code": 10001,
  "msg": "参数错误"
}
```

### 错误码

| 错误码 | 含义 |
|--------|------|
| 0 | 成功 |
| 10001 | 参数错误 |
| 10002 | 未授权或token无效 |
| 10003 | 资源不存在 |
| 10004 | 数据重复 |
| 10005 | 服务器内部错误 |
| 20001 | 排班时间冲突 |
| 20002 | 该时段已有预约 |
| 20003 | 无效的时间段 |
| 20004 | 该医生在此时段无排班 |

---

### 患者管理

**创建患者** `POST /patients`

```json
{
  "name": "张三",
  "phone": "13800138000",
  "gender": "男",
  "birth_date": "1990-01-15T00:00:00Z",
  "address": "幸福路1号"
}
```

**更新患者** `PUT /patients/:id`

**删除患者** `DELETE /patients/:id`

**查询患者** `GET /patients/:id`

**患者列表** `GET /patients?phone=138&name=张`

- `phone` 和 `name` 均为模糊搜索，可选

---

### 医生管理

**创建医生** `POST /doctors`

```json
{
  "name": "李医生",
  "phone": "13900139000",
  "dept": "全科",
  "title": "主治医师"
}
```

**医生列表** `GET /doctors`

**查询医生** `GET /doctors/:id`

---

### 医生排班

**创建排班** `POST /doctors/:doctor_id/schedules`

```json
{
  "weekday": 1,
  "start_time": "08:00",
  "end_time": "12:00"
}
```

- `weekday`: 1=周一 ... 7=周日
- 同一医生的同一星期几，排班时段不可重叠

**更新排班** `PUT /schedules/:id`

**删除排班** `DELETE /schedules/:id`

**查询排班** `GET /doctors/:doctor_id/schedules`

---

### 预约管理

**创建预约** `POST /appointments`

```json
{
  "patient_id": 1,
  "doctor_id": 1,
  "app_date": "2024-07-15",
  "start_time": "08:00",
  "end_time": "08:30",
  "remark": "头疼"
}
```

- 创建时自动校验：医生当天是否有排班覆盖该时段
- 同一医生同一天同一时间段只允许一个患者预约（状态为 booked）
- 创建后状态默认 `booked`

**取消预约** `PUT /appointments/:id/cancel`

- 仅 `booked` 状态可取消，取消后变为 `cancelled`，时段自动释放

**标记完成** `PUT /appointments/:id/complete`

- 仅 `booked` 状态可标记完成，变为 `completed`

**查询某天预约** `GET /appointments?date=2024-07-15`

- 返回当天所有预约，包含患者和医生信息

---

## 示例调用流程

```bash
TOKEN="Bearer clinic-secret-token-2024"
BASE="http://localhost:8080"

# 创建患者
curl -s -X POST "$BASE/patients" -H "Authorization: $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"张三","phone":"13800138000","gender":"男"}'

# 创建医生
curl -s -X POST "$BASE/doctors" -H "Authorization: $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"李医生","dept":"全科","title":"主治医师"}'

# 设置排班（周一 08:00-12:00）
curl -s -X POST "$BASE/doctors/1/schedules" -H "Authorization: $TOKEN" -H "Content-Type: application/json" \
  -d '{"weekday":1,"start_time":"08:00","end_time":"12:00"}'

# 创建预约（假设 2024-07-15 是周一）
curl -s -X POST "$BASE/appointments" -H "Authorization: $TOKEN" -H "Content-Type: application/json" \
  -d '{"patient_id":1,"doctor_id":1,"app_date":"2024-07-15","start_time":"08:00","end_time":"08:30","remark":"头疼"}'

# 查询当天预约
curl -s "$BASE/appointments?date=2024-07-15" -H "Authorization: $TOKEN"

# 取消预约
curl -s -X PUT "$BASE/appointments/1/cancel" -H "Authorization: $TOKEN"

# 标记完成
curl -s -X PUT "$BASE/appointments/1/complete" -H "Authorization: $TOKEN"
```

## 项目结构

```
├── main.go              # 入口，路由注册
├── config/config.go     # 配置加载
├── middleware/auth.go   # Token 认证中间件
├── model/models.go      # 数据模型 + 自动迁移
├── handler/
│   ├── patient.go       # 患者增删改查
│   ├── doctor.go        # 医生增删查
│   ├── schedule.go      # 排班管理
│   ├── appointment.go   # 预约管理
│   └── util.go          # 工具函数
├── response/response.go # 统一响应 + 错误码
├── docker-compose.yml   # MySQL 容器
└── go.mod
```
