# 🧾 Smart Invoice OCR System

ระบบจัดการและสกัดข้อมูลจากใบแจ้งหนี้อัตโนมัติ (Invoice OCR) ที่ทำงานแบบ Microservices รองรับภาษาไทยและอังกฤษอย่างแม่นยำด้วย **EasyOCR**

---

## 🏗️ System Architecture
- **Frontend:** Next.js 15 (App Router), TailwindCSS, Axios
- **Backend:** Go (Gin Framework), GORM, PostgreSQL
- **OCR Service:** Python (FastAPI), EasyOCR (Thai/English Support)
- **Storage:** MinIO (S3 Compatible Object Storage)
- **Infrastructure:** Docker Compose

---

## 🚀 Getting Started

### 1. Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- Python 3.10 (Recommended)

### 2. Infrastructure Setup (Docker)
สร้างไฟล์ `.env` จากข้อมูลใน `.env.dev` แล้วรัน Database และ Storage:
```bash
docker compose --env-file .env.dev up -d
```

### 3. OCR Service Setup (Python)
เข้าสู่โฟลเดอร์ `ocr-service` และสร้างสภาพแวดล้อมจำลอง:
```bash
cd ocr-service
py -3.10 -m venv venv310
.\venv310\Scripts\activate

# ติดตั้ง Library ที่จำเป็น
pip install numpy==1.26.4
pip install easyocr fastapi uvicorn python-multipart Pillow
```
รัน OCR Service:
```bash
python app.py
```

### 4. Backend Setup (Go)
เข้าสู่โฟลเดอร์ `backend` และติดตั้ง Dependencies:
```bash
cd backend
go mod tidy
```
รัน Backend Server:
```bash
go run cmd/api/main.go
```

### 5. Frontend Setup (Next.js)
เข้าสู่โฟลเดอร์ `frontend` และติดตั้ง Dependencies:
```bash
cd frontend
npm install
```
รัน Frontend Development Server:
```bash
npm run dev
```

---

## 🌟 Key Features
- **Thai Language Support:** รองรับการอ่านภาษาไทยได้อย่างแม่นยำ
- **Automatic Summary:** สกัดยอดรวม (Total Amount) และแยกประเภทค่าใช้จ่ายอัตโนมัติ
- **Secure File Storage:** เก็บไฟล์ใบแจ้งหนี้ต้นฉบับไว้บน MinIO
- **Clean Architecture:** ออกแบบโค้ดแยกส่วนกันชัดเจน ง่ายต่อการดูแลและสเกลต่อ

---

## 🔒 Security Note
- ไฟล์ `.env` และข้อมูลความลับถูกยกเว้นจากการเก็บเข้าระบบ Git (ผ่าน `.gitignore`)
- การเชื่อมต่อ Backend ใช้ Middleware สำหรับจัดการ CORS อย่างปลอดภัย
