#!/bin/bash

# 日志颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_ROOT=$(cd "$(dirname "$0")" && pwd)
BACKEND_DIR="$PROJECT_ROOT/server"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查依赖
check_dependencies() {
    log_info "检查依赖..."

    # 检查Go
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.23或更高版本"
        exit 1
    fi

    # 检查Node.js
    if ! command -v node &> /dev/null; then
        log_error "Node.js未安装，请先安装Node.js 24或更高版本"
        exit 1
    fi

    # 检查npm
    if ! command -v npm &> /dev/null; then
        log_error "npm未安装，请先安装npm"
        exit 1
    fi

    # 检查nginx
    if ! command -v nginx &> /dev/null; then
        log_warn "nginx未安装，前端将使用Node.js serve运行"
    fi

    log_info "依赖检查完成"
}

# 构建后端
build_backend() {
    log_info "构建后端..."
    cd "$BACKEND_DIR" || exit 1

    # 设置Go代理
    export GOPROXY=https://mirrors.aliyun.com/goproxy/

    # 下载依赖
    go mod download

    # 构建
    go build -o main .

    if [ $? -eq 0 ]; then
        log_info "后端构建成功"
    else
        log_error "后端构建失败"
        exit 1
    fi

    cd "$PROJECT_ROOT"
}

# 构建前端
build_frontend() {
    log_info "构建前端..."
    cd "$FRONTEND_DIR" || exit 1

    # 安装依赖
    npm install

    # 构建
    npm run build

    if [ $? -eq 0 ]; then
        log_info "前端构建成功"
    else
        log_error "前端构建失败"
        exit 1
    fi

    cd "$PROJECT_ROOT"
}

# 创建必要目录
create_directories() {
    log_info "创建必要目录..."

    # 创建日志目录
    mkdir -p /var/log/app
    mkdir -p "$PROJECT_ROOT/logs"

    log_info "目录创建完成"
}

# 启动后端
start_backend() {
    log_info "启动后端服务..."
    cd "$BACKEND_DIR" || exit 1

    # 设置环境变量
    export GIN_MODE=release
    export PORT=10001

    # 后台启动
    nohup ./main --port=$PORT > ../logs/backend.log 2>&1 &
    BACKEND_PID=$!

    echo $BACKEND_PID > ../logs/backend.pid

    # 等待服务启动
    sleep 3

    if kill -0 $BACKEND_PID 2>/dev/null; then
        log_info "后端服务启动成功 (PID: $BACKEND_PID)"
    else
        log_error "后端服务启动失败"
        exit 1
    fi

    cd "$PROJECT_ROOT"
}

# 启动前端
start_frontend() {
    log_info "启动前端服务..."
    cd "$FRONTEND_DIR" || exit 1

    # 使用serve
    log_info "使用serve启动前端..."

    # 安装serve（如果没有）
    if ! command -v serve &> /dev/null; then
        npm install -g serve
    fi

    # 后台启动
    nohup serve -s dist -l 10002 > ../logs/frontend.log 2>&1 &
    FRONTEND_PID=$!

    echo $FRONTEND_PID > ../logs/frontend.pid

    # 等待服务启动
    sleep 2

    if kill -0 $FRONTEND_PID 2>/dev/null; then
        log_info "前端服务启动成功 (PID: $FRONTEND_PID)"
    else
         log_error "前端服务启动失败"
        exit 1
    fi

    cd "$PROJECT_ROOT"
}

# 停止服务
stop_services() {
    log_info "停止服务..."

    # 停止后端
    if [ -f logs/backend.pid ]; then
        BACKEND_PID=$(cat logs/backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            kill $BACKEND_PID
            log_info "后端服务已停止"
        fi
        rm -f logs/backend.pid
    fi

    # 停止前端
    if [ -f logs/frontend.pid ]; then
        FRONTEND_PID=$(cat logs/frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            kill $FRONTEND_PID
            log_info "前端服务已停止"
        fi
        rm -f logs/frontend.pid
    fi
}

# 显示服务状态
show_status() {
    log_info "服务状态:"

    # 后端状态
    if [ -f logs/backend.pid ]; then
        BACKEND_PID=$(cat logs/backend.pid)
        if kill -0 $BACKEND_PID 2>/dev/null; then
            echo -e "  后端: ${GREEN}运行中${NC} (PID: $BACKEND_PID, 端口: 10001)"
        else
            echo -e "  后端: ${RED}已停止${NC}"
        fi
    else
        echo -e "  后端: ${RED}已停止${NC}"
    fi

    # 前端状态
    if [ -f logs/frontend.pid ]; then
        FRONTEND_PID=$(cat logs/frontend.pid)
        if kill -0 $FRONTEND_PID 2>/dev/null; then
            echo -e "  前端: ${GREEN}运行中${NC} (PID: $FRONTEND_PID, 端口: 10002)"
        else
            echo -e "  前端: ${RED}已停止${NC}"
        fi
    else
        echo -e "  前端: ${RED}已停止${NC}"
    fi
}

# 主函数
main() {
    case "$1" in
        start)
            log_info "启动日志查看器项目..."
            check_dependencies
            create_directories
            build_backend
            build_frontend
            start_backend
            start_frontend
            log_info "项目启动完成!"
            echo "访问地址:"
            echo "  前端: http://your-server-ip:10002"
            echo "  后端API: http://your-server-ip:10001"
            ;;
        stop)
            stop_services
            log_info "项目已停止"
            ;;
        restart)
            stop_services
            sleep 2
            main start
            ;;
        status)
            show_status
            ;;
        build)
            log_info "重新构建项目..."
            check_dependencies
            build_backend
            build_frontend
            log_info "构建完成"
            ;;
        *)
            echo "用法: $0 {start|stop|restart|status|build}"
            echo "  start   - 启动服务"
            echo "  stop    - 停止服务"
            echo "  restart - 重启服务"
            echo "  status  - 查看状态"
            echo "  build   - 重新构建"
            exit 1
            ;;
    esac
}

main "$@"