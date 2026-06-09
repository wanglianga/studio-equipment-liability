# 影棚器材损坏责任判定服务

## 原始需求

> 影棚需要器材损坏责任判定服务，Go 接口处理器材借出、归还检查、损坏登记、维修报价、押金扣除和客户申诉。业务内容包括器材编号、镜头型号、灯具功率、借用时间、使用棚位、押金、借前照片、还场照片、故障点位、维修费用和责任结论。客户租用相机、镜头、闪光灯或背景架后，管理员在归还时检查外观和功能；发现损坏后根据借前状态、使用记录和现场证据判断责任。服务要区分正常磨损、客户损坏、前序遗留、运输碰撞、配件缺失和无法判定。
> 借前状态要参与判定。器材出借前已有划痕或轻微松动时，接口保存照片和备注，归还时不能把旧问题全部算给本次客户。
> 配件缺失要单独扣费。镜头盖、电池、柔光罩或收纳包未归还时，服务按配件价目生成扣款，不必进入完整维修流程。

## 项目简介

基于 Go 实现的影棚器材损坏责任判定 HTTP API 服务，提供器材管理、借用登记、归还检查、损坏报告、维修报价、押金扣除和客户申诉的完整业务流程。

### 责任判定类型

| 类型 | 英文标识 | 说明 |
|------|---------|------|
| 正常磨损 | `normal_wear` | 轻微划痕或磨损，属于正常使用损耗 |
| 客户损坏 | `customer_damage` | 借前照片显示完好，归还时发现明显损坏 |
| 前序遗留 | `previous_remnant` | 器材此前已有未判定损坏记录或历史损坏 |
| 运输碰撞 | `transport_impact` | 灯具/背景架在棚位间移动中出现碰撞 |
| 配件缺失 | `accessory_missing` | 归还时发现配件缺失 |
| 无法判定 | `undetermined` | 缺乏借前照片对比或证据不足 |

## 技术栈

- Go 1.22+
- Go 标准库 `net/http`（路由使用 Go 1.22 方法路由）
- 内存存储（并发安全）

## 目录结构

```
├── main.go              # 入口，启动 HTTP 服务
├── model/model.go       # 数据模型定义
├── store/store.go       # 内存存储层
├── service/service.go   # 业务逻辑层（含责任判定引擎）
├── handler/handler.go   # HTTP 处理层
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── go.mod
└── README.md
```

## API 接口

### 器材管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/equipments` | 创建器材 |
| GET | `/api/equipments` | 器材列表 |
| GET | `/api/equipments/{id}` | 器材详情 |

### 借用与归还

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/borrow` | 器材借出 |
| POST | `/api/borrow/return` | 归还检查 |
| GET | `/api/borrow` | 借用记录列表 |
| GET | `/api/borrow/{id}` | 借用记录详情 |

### 损坏与维修

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/damage` | 损坏登记 |
| GET | `/api/damage` | 损坏报告列表 |
| GET | `/api/damage/{id}` | 损坏报告详情 |
| POST | `/api/repair-quote` | 维修报价 |

### 押金与申诉

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/deduction` | 押金扣除 |
| POST | `/api/appeal` | 客户申诉 |
| POST | `/api/appeal/review` | 申诉审核 |
| GET | `/api/appeal` | 申诉列表 |

### 健康检查

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 服务健康检查 |

## 启动方式

### 前置要求

- Docker 和 Docker Compose（推荐方式）
- 或 Go 1.22+（本地开发方式）

### Docker 一键启动（推荐）

#### 1. 构建并启动

```bash
docker compose up --build
```

后台运行：

```bash
docker compose up --build -d
```

#### 2. 停止服务

```bash
docker compose down
```

访问地址：http://localhost:8090

### 本地启动

#### 1. 安装依赖

```bash
go mod download
```

#### 2. 启动服务

```bash
go run .
```

访问地址：http://localhost:8090

## 使用示例

### 1. 创建器材

```bash
curl -X POST http://localhost:8090/api/equipments \
  -H "Content-Type: application/json" \
  -d '{
    "category": "lens",
    "brand": "Canon",
    "model": "EF 70-200mm",
    "lens_model": "EF 70-200mm f/2.8L IS III USM",
    "pre_borrow_photo": "photo_before_001.jpg"
  }'
```

### 2. 借出器材

```bash
curl -X POST http://localhost:8090/api/borrow \
  -H "Content-Type: application/json" \
  -d '{
    "equipment_id": "EQ-0001",
    "customer_name": "张三",
    "customer_phone": "13800138000",
    "studio_position": "A-3",
    "deposit": 2000,
    "pre_borrow_photos": ["borrow_before_1.jpg", "borrow_before_2.jpg"]
  }'
```

### 3. 归还检查

```bash
curl -X POST http://localhost:8090/api/borrow/return \
  -H "Content-Type: application/json" \
  -d '{
    "borrow_record_id": "BR-0001",
    "return_photos": ["return_photo_1.jpg", "return_photo_2.jpg"]
  }'
```

### 4. 登记损坏

```bash
curl -X POST http://localhost:8090/api/damage \
  -H "Content-Type: application/json" \
  -d '{
    "borrow_record_id": "BR-0001",
    "fault_points": [
      {"location": "镜头前端", "description": "镜片划痕", "severity": "severe"}
    ],
    "return_photos": ["damage_photo_1.jpg"]
  }'
```

### 5. 维修报价

```bash
curl -X POST http://localhost:8090/api/repair-quote \
  -H "Content-Type: application/json" \
  -d '{
    "damage_report_id": "DM-0001",
    "repair_cost": 800,
    "labor_cost": 200,
    "description": "更换前端镜片组"
  }'
```

### 6. 押金扣除

```bash
curl -X POST http://localhost:8090/api/deduction \
  -H "Content-Type: application/json" \
  -d '{
    "borrow_record_id": "BR-0001",
    "repair_quote_id": "RQ-0001"
  }'
```

### 7. 客户申诉

```bash
curl -X POST http://localhost:8090/api/appeal \
  -H "Content-Type: application/json" \
  -d '{
    "borrow_record_id": "BR-0001",
    "customer_name": "张三",
    "reason": "镜头划痕借出时已存在，借前照片可证明",
    "evidence": ["borrow_before_closeup.jpg"]
  }'
```

### 8. 申诉审核

```bash
curl -X POST http://localhost:8090/api/appeal/review \
  -H "Content-Type: application/json" \
  -d '{
    "appeal_id": "AP-0001",
    "accepted": true,
    "review_note": "借前照片确认划痕已存在，申诉成立"
  }'
```
