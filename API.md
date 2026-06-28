# 社区诊所预约管理系统 API 文档

## 通用约定

### 基础地址

```
http://localhost:8080
```

### 认证方式

除 `GET /ping` 外，所有接口均需在请求头携带 Bearer Token：

```
Authorization: Bearer clinic-secret-token-2024
```

未携带或 Token 无效时返回：

```json
{
  "code": 10002,
  "msg": "未授权或token无效"
}
```

### 统一响应格式

成功：

```json
{
  "code": 0,
  "msg": "成功",
  "data": { ... }
}
```

失败：

```json
{
  "code": 10001,
  "msg": "参数错误"
}
```

`data` 字段仅在成功时返回，失败时省略。

### 时间格式约定

| 字段 | 格式 | 示例 |
|------|------|------|
| app_date | YYYY-MM-DD | 2024-07-15 |
| start_time / end_time | HH:mm | 08:00 |
| birth_date | RFC3339 | 1990-01-15T00:00:00Z |
| created_at / updated_at | RFC3339 | 2024-07-15T10:30:00+08:00 |

---

## 一、患者管理

### 1.1 创建患者

`POST /patients` | 需要Token

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 姓名，最长50字符 |
| phone | string | 是 | 手机号，唯一，最长20字符 |
| gender | string | 否 | 性别 |
| birth_date | string | 否 | 出生日期，RFC3339格式 |
| address | string | 否 | 地址，最长200字符 |

**请求示例**

```json
{
  "name": "张三",
  "phone": "13800138000",
  "gender": "男",
  "birth_date": "1990-01-15T00:00:00Z",
  "address": "幸福路1号"
}
```

**成功响应（201）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "张三",
    "phone": "13800138000",
    "gender": "男",
    "birth_date": "1990-01-15T00:00:00Z",
    "address": "幸福路1号",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 缺少必填字段、JSON格式错误、姓名或手机号为空 |
| 10004 | 手机号已被其他患者使用 |
| 10005 | 服务器内部错误 |

---

### 1.2 更新患者

`PUT /patients/:id` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 患者ID |

**请求参数（JSON Body）**

仅传入需要更新的字段，未传入的字段保持不变。

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 否 | 姓名 |
| phone | string | 否 | 手机号 |
| gender | string | 否 | 性别 |
| birth_date | string | 否 | 出生日期 |
| address | string | 否 | 地址 |

**请求示例**

```json
{
  "phone": "13900139000",
  "address": "新地址88号"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "张三",
    "phone": "13900139000",
    "gender": "男",
    "birth_date": "1990-01-15T00:00:00Z",
    "address": "新地址88号",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T11:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字、JSON格式错误 |
| 10003 | 患者ID不存在 |
| 10004 | 新手机号与其他患者重复 |
| 10005 | 服务器内部错误 |

---

### 1.3 删除患者

`DELETE /patients/:id` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 患者ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功"
}
```

> 删除为软删除，数据库中保留记录但不再出现在查询结果中。

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字 |
| 10003 | 患者ID不存在 |

---

### 1.4 查询患者详情

`GET /patients/:id` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 患者ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "张三",
    "phone": "13800138000",
    "gender": "男",
    "birth_date": "1990-01-15T00:00:00Z",
    "address": "幸福路1号",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字 |
| 10003 | 患者ID不存在 |

---

### 1.5 患者列表

`GET /patients` | 需要Token

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| phone | string | 否 | 按手机号模糊搜索 |
| name | string | 否 | 按姓名模糊搜索 |

两个参数均可选，同时传入时取交集。

**请求示例**

```
GET /patients?phone=138&name=张
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "张三",
      "phone": "13800138000",
      "gender": "男",
      "birth_date": "1990-01-15T00:00:00Z",
      "address": "幸福路1号",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    }
  ]
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10005 | 服务器内部错误 |

---

### 1.6 查询患者就诊历史

`GET /patients/:id/history` | 需要Token

返回该患者所有已完成就诊的记录，按就诊日期倒序排列。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 患者ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "app_date": "2024-07-15",
      "start_time": "08:00",
      "end_time": "08:30",
      "doctor_name": "李医生",
      "dept": "全科",
      "diagnosis": "上呼吸道感染",
      "prescription": "阿莫西林 0.5g tid×5天"
    },
    {
      "app_date": "2024-07-08",
      "start_time": "09:00",
      "end_time": "09:30",
      "doctor_name": "王医生",
      "dept": "内科",
      "diagnosis": "高血压",
      "prescription": "硝苯地平缓释片 30mg qd"
    }
  ]
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字 |
| 10003 | 患者ID不存在 |
| 10005 | 服务器内部错误 |

**curl 示例**

```bash
curl -s http://localhost:8080/patients/1/history \
  -H "Authorization: Bearer clinic-secret-token-2024"
```

---

## 二、医生管理

### 2.1 创建医生

`POST /doctors` | 需要Token

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| name | string | 是 | 姓名，最长50字符 |
| phone | string | 否 | 手机号，最长20字符 |
| dept | string | 否 | 科室，最长50字符 |
| title | string | 否 | 职称，最长50字符 |

**请求示例**

```json
{
  "name": "李医生",
  "phone": "13900139000",
  "dept": "全科",
  "title": "主治医师"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "李医生",
    "phone": "13900139000",
    "dept": "全科",
    "title": "主治医师",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | JSON格式错误、姓名为空 |
| 10005 | 服务器内部错误 |

---

### 2.2 医生列表

`GET /doctors` | 需要Token

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "name": "李医生",
      "phone": "13900139000",
      "dept": "全科",
      "title": "主治医师",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    }
  ]
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10005 | 服务器内部错误 |

---

### 2.3 查询医生详情

`GET /doctors/:id` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 医生ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "name": "李医生",
    "phone": "13900139000",
    "dept": "全科",
    "title": "主治医师",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10003 | 医生ID不存在 |

---

## 三、医生排班

### 3.1 创建排班

`POST /doctors/:doctor_id/schedules` | 需要Token

为指定医生设置每周固定出诊时段。同一医生同一星期几的排班时段不可重叠。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| doctor_id | uint | 医生ID |

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| weekday | int | 是 | 星期几，1=周一 2=周二 ... 7=周日 |
| start_time | string | 是 | 开始时间，格式 HH:mm |
| end_time | string | 是 | 结束时间，格式 HH:mm |

**请求示例**

```json
{
  "weekday": 1,
  "start_time": "08:00",
  "end_time": "12:00"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "doctor_id": 1,
    "weekday": 1,
    "start_time": "08:00",
    "end_time": "12:00",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | doctor_id为0、weekday不在1-7范围、时间格式错误、开始时间≥结束时间 |
| 10003 | 医生不存在 |
| 10004 | 同一医生同一星期几同一时段已存在排班（唯一索引冲突） |
| 20001 | 新排班与该医生同一天已有排班的时段重叠 |
| 20003 | start_time 或 end_time 不符合 HH:mm 格式 |
| 10005 | 服务器内部错误 |

**curl 示例**

```bash
curl -s -X POST http://localhost:8080/doctors/1/schedules \
  -H "Authorization: Bearer clinic-secret-token-2024" \
  -H "Content-Type: application/json" \
  -d '{"weekday":1,"start_time":"08:00","end_time":"12:00"}'
```

---

### 3.2 更新排班

`PUT /schedules/:id` | 需要Token

仅传入需要修改的字段，未传入的字段保持不变。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 排班记录ID |

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| weekday | int | 否 | 星期几，1-7 |
| start_time | string | 否 | 开始时间，HH:mm |
| end_time | string | 否 | 结束时间，HH:mm |

**请求示例**

```json
{
  "start_time": "08:30",
  "end_time": "12:30"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "doctor_id": 1,
    "weekday": 1,
    "start_time": "08:30",
    "end_time": "12:30",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T11:00:00+08:00"
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | JSON格式错误 |
| 10003 | 排班记录ID不存在 |
| 10004 | 修改后与已有排班唯一索引冲突 |
| 20001 | 修改后时段与该医生同一天其他排班重叠 |
| 20003 | 时间格式无效、开始时间≥结束时间 |
| 10005 | 服务器内部错误 |

---

### 3.3 删除排班

`DELETE /schedules/:id` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 排班记录ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功"
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字 |
| 10003 | 排班记录ID不存在 |

---

### 3.4 查询医生排班列表

`GET /doctors/:doctor_id/schedules` | 需要Token

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| doctor_id | uint | 医生ID |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "doctor_id": 1,
      "weekday": 1,
      "start_time": "08:00",
      "end_time": "12:00",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    },
    {
      "id": 2,
      "doctor_id": 1,
      "weekday": 1,
      "start_time": "14:00",
      "end_time": "17:00",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    }
  ]
}
```

结果按 weekday、start_time 升序排列。

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径doctor_id非数字 |
| 10005 | 服务器内部错误 |

---

## 四、预约管理

### 4.1 创建预约

`POST /appointments` | 需要Token

创建预约时系统自动执行以下校验：
1. 患者和医生是否存在
2. 该医生在预约日期的星期几是否有排班覆盖请求时段
3. 同一医生同一天该时段是否已有 `booked` 状态的预约（事务 + 行锁保证并发安全）
4. 数据库条件唯一索引兜底防重复

校验通过后创建预约，状态默认为 `booked`。

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| patient_id | uint | 是 | 患者ID |
| doctor_id | uint | 是 | 医生ID |
| app_date | string | 是 | 预约日期，格式 YYYY-MM-DD |
| start_time | string | 是 | 开始时间，格式 HH:mm |
| end_time | string | 是 | 结束时间，格式 HH:mm |
| remark | string | 否 | 备注，最长500字符 |

**请求示例**

```json
{
  "patient_id": 1,
  "doctor_id": 1,
  "app_date": "2024-07-15",
  "start_time": "08:00",
  "end_time": "08:30",
  "remark": "头疼、发热两天"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "patient_id": 1,
    "doctor_id": 1,
    "app_date": "2024-07-15",
    "start_time": "08:00",
    "end_time": "08:30",
    "status": "booked",
    "remark": "头疼、发热两天",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T10:00:00+08:00",
    "patient": {
      "id": 1,
      "name": "张三",
      "phone": "13800138000",
      "gender": "男",
      "birth_date": "1990-01-15T00:00:00Z",
      "address": "幸福路1号",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    },
    "doctor": {
      "id": 1,
      "name": "李医生",
      "phone": "13900139000",
      "dept": "全科",
      "title": "主治医师",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T10:00:00+08:00"
    }
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 必填字段缺失、app_date格式错误、JSON解析失败 |
| 10003 | 患者ID不存在、医生ID不存在 |
| 20002 | 该医生该时段已有 booked 状态的预约（含并发场景） |
| 20003 | start_time/end_time 格式无效、开始时间≥结束时间 |
| 20004 | 该医生在预约日期的星期几没有覆盖该时段的排班 |
| 10005 | 服务器内部错误 |

**curl 示例**

```bash
curl -s -X POST http://localhost:8080/appointments \
  -H "Authorization: Bearer clinic-secret-token-2024" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1,
    "doctor_id": 1,
    "app_date": "2024-07-15",
    "start_time": "08:00",
    "end_time": "08:30",
    "remark": "头疼、发热两天"
  }'
```

---

### 4.2 取消预约

`PUT /appointments/:id/cancel` | 需要Token

将预约状态从 `booked` 改为 `cancelled`，取消后该时段可被重新预约。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 预约ID |

**请求体**：无

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "id": 1,
    "patient_id": 1,
    "doctor_id": 1,
    "app_date": "2024-07-15",
    "start_time": "08:00",
    "end_time": "08:30",
    "status": "cancelled",
    "remark": "头疼、发热两天",
    "created_at": "2024-07-15T10:00:00+08:00",
    "updated_at": "2024-07-15T11:00:00+08:00",
    "patient": { "..." : "..." },
    "doctor": { "..." : "..." }
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字、预约状态不是 booked |
| 10003 | 预约ID不存在 |
| 10005 | 服务器内部错误 |

**curl 示例**

```bash
curl -s -X PUT http://localhost:8080/appointments/1/cancel \
  -H "Authorization: Bearer clinic-secret-token-2024"
```

---

### 4.3 完成就诊

`PUT /appointments/:id/complete` | 需要Token

将预约状态改为 `completed`，同时创建就诊记录（含诊断结论和处方）。操作在事务中执行，保证预约状态与就诊记录的一致性。

**路径参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| id | uint | 预约ID |

**请求参数（JSON Body）**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| diagnosis | string | 是 | 诊断结论，最长500字符 |
| prescription | string | 否 | 处方内容 |

**请求示例**

```json
{
  "diagnosis": "上呼吸道感染",
  "prescription": "阿莫西林 0.5g tid×5天"
}
```

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": {
    "appointment": {
      "id": 1,
      "patient_id": 1,
      "doctor_id": 1,
      "app_date": "2024-07-15",
      "start_time": "08:00",
      "end_time": "08:30",
      "status": "completed",
      "remark": "头疼、发热两天",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T12:00:00+08:00",
      "patient": { "..." : "..." },
      "doctor": { "..." : "..." }
    },
    "visit_record": {
      "id": 1,
      "appointment_id": 1,
      "diagnosis": "上呼吸道感染",
      "prescription": "阿莫西林 0.5g tid×5天",
      "created_at": "2024-07-15T12:00:00+08:00",
      "updated_at": "2024-07-15T12:00:00+08:00",
      "appointment": { "..." : "..." }
    }
  }
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | 路径id非数字、JSON格式错误、diagnosis为空、预约状态不是 booked |
| 10003 | 预约ID不存在 |
| 10005 | 服务器内部错误（含事务失败回滚） |

**curl 示例**

```bash
curl -s -X PUT http://localhost:8080/appointments/1/complete \
  -H "Authorization: Bearer clinic-secret-token-2024" \
  -H "Content-Type: application/json" \
  -d '{"diagnosis":"上呼吸道感染","prescription":"阿莫西林 0.5g tid×5天"}'
```

---

### 4.4 查询某天的预约列表

`GET /appointments?date=YYYY-MM-DD` | 需要Token

返回指定日期所有状态的预约，包含关联的患者和医生信息，按开始时间升序排列。

**查询参数**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| date | string | 是 | 日期，格式 YYYY-MM-DD |

**成功响应（200）**

```json
{
  "code": 0,
  "msg": "成功",
  "data": [
    {
      "id": 1,
      "patient_id": 1,
      "doctor_id": 1,
      "app_date": "2024-07-15",
      "start_time": "08:00",
      "end_time": "08:30",
      "status": "completed",
      "remark": "头疼、发热两天",
      "created_at": "2024-07-15T10:00:00+08:00",
      "updated_at": "2024-07-15T12:00:00+08:00",
      "patient": {
        "id": 1,
        "name": "张三",
        "phone": "13800138000",
        "gender": "男",
        "birth_date": "1990-01-15T00:00:00Z",
        "address": "幸福路1号",
        "created_at": "2024-07-15T10:00:00+08:00",
        "updated_at": "2024-07-15T10:00:00+08:00"
      },
      "doctor": {
        "id": 1,
        "name": "李医生",
        "phone": "13900139000",
        "dept": "全科",
        "title": "主治医师",
        "created_at": "2024-07-15T10:00:00+08:00",
        "updated_at": "2024-07-15T10:00:00+08:00"
      }
    }
  ]
}
```

**可能返回的错误码**

| 错误码 | 触发条件 |
|--------|----------|
| 10001 | date参数缺失或格式错误 |
| 10005 | 服务器内部错误 |

**curl 示例**

```bash
curl -s "http://localhost:8080/appointments?date=2024-07-15" \
  -H "Authorization: Bearer clinic-secret-token-2024"
```

---

## 五、就诊记录

就诊记录不提供独立的外部接口，通过以下方式访问：

- **创建**：通过「完成就诊」接口（`PUT /appointments/:id/complete`）自动创建，一条预约对应一条就诊记录
- **查询**：通过「查询患者就诊历史」接口（`GET /patients/:id/history`）查看

### 就诊记录数据结构

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 记录ID |
| appointment_id | uint | 关联预约ID，唯一 |
| diagnosis | string | 诊断结论 |
| prescription | string | 处方内容 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |

---

## 六、健康检查

### 6.1 服务连通性检测

`GET /ping` | 不需要Token

**成功响应（200）**

```json
{
  "message": "pong"
}
```

---

## 七、错误码汇总

### 系统级错误（10001-10005）

| Code | Message | HTTP Status | 触发条件 |
|------|---------|-------------|----------|
| 0 | 成功 | 200 | 请求成功 |
| 10001 | 参数错误 | 400 | 必填字段缺失、JSON解析失败、格式不符合要求、路径参数非数字 |
| 10002 | 未授权或token无效 | 401 | 未携带 Authorization 头、Token 格式错误、Token 值不匹配 |
| 10003 | 资源不存在 | 404 | 查询的患者/医生/排班/预约ID在数据库中不存在 |
| 10004 | 数据重复 | 409 | 违反唯一约束（如手机号重复、排班唯一索引冲突） |
| 10005 | 服务器内部错误 | 500 | 数据库操作异常、事务失败等不可预期错误 |

### 业务级错误（20001-20004）

| Code | Message | HTTP Status | 触发条件 |
|------|---------|-------------|----------|
| 20001 | 排班时间冲突 | 409 | 新增或修改排班时，时段与该医生同一天已有排班重叠 |
| 20002 | 该时段已有预约 | 409 | 创建预约时，同一医生同一天该时段已有 booked 状态的预约（含并发竞争） |
| 20003 | 无效的时间段 | 400 | 时间格式不符合 HH:mm、开始时间≥结束时间 |
| 20004 | 该医生在此时段无排班 | 400 | 创建预约时，预约日期的星期几该医生没有覆盖请求时段的排班 |

### 并发安全说明

创建预约接口使用事务 + `SELECT ... FOR UPDATE` 行锁保护，配合数据库条件唯一索引（`booked_slot_key` 生成列），确保同一医生同一时段即使在高并发下也只会产生一条 `booked` 预约。并发请求中仅第一个成功，其余返回 `20002`。

---

## 八、完整调用流程示例

以下命令按顺序执行，演示从创建基础数据到完成就诊的完整闭环：

```bash
TOKEN="Bearer clinic-secret-token-2024"
BASE="http://localhost:8080"

# 1. 创建患者
curl -s -X POST "$BASE/patients" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"张三","phone":"13800138000","gender":"男"}'

# 2. 创建医生
curl -s -X POST "$BASE/doctors" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"李医生","dept":"全科","title":"主治医师"}'

# 3. 设置排班（周一 08:00-12:00）
curl -s -X POST "$BASE/doctors/1/schedules" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"weekday":1,"start_time":"08:00","end_time":"12:00"}'

# 4. 创建预约（2024-07-15 是周一）
curl -s -X POST "$BASE/appointments" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"patient_id":1,"doctor_id":1,"app_date":"2024-07-15","start_time":"08:00","end_time":"08:30","remark":"头疼"}'

# 5. 查询当天预约
curl -s "$BASE/appointments?date=2024-07-15" \
  -H "Authorization: $TOKEN"

# 6. 完成就诊（填写诊断和处方）
curl -s -X PUT "$BASE/appointments/1/complete" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"diagnosis":"上呼吸道感染","prescription":"阿莫西林 0.5g tid×5天"}'

# 7. 查询患者就诊历史
curl -s "$BASE/patients/1/history" \
  -H "Authorization: $TOKEN"

# === 或者取消预约 ===

# 6'. 取消预约（替代步骤6）
curl -s -X PUT "$BASE/appointments/1/cancel" \
  -H "Authorization: $TOKEN"

# 7'. 取消后可重新预约该时段
curl -s -X POST "$BASE/appointments" \
  -H "Authorization: $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"patient_id":2,"doctor_id":1,"app_date":"2024-07-15","start_time":"08:00","end_time":"08:30"}'
```
