import io
import re
from fastapi import FastAPI, UploadFile, File, HTTPException
import easyocr
from PIL import Image, ImageOps, ImageEnhance
import numpy as np

app = FastAPI()

# ประกาศตัวอ่านโดยระบุภาษาไทย ['th'] และอังกฤษ ['en']
reader = easyocr.Reader(['th', 'en']) 

def find_value_in_row(results, keywords):
    """ ค้นหาตัวเลขที่อยู่ทางขวาของ Keyword ในระดับบรรทัด (Y) เดียวกัน """
    for i, (bbox, text, prob) in enumerate(results):
        # ถ้าเจอ Keyword ที่เราสนใจ
        if any(kw.lower() in text.lower() for kw in keywords):
            # หาค่ากึ่งกลางแนวตั้ง (Center Y) ของ Keyword
            center_y = (bbox[0][1] + bbox[2][1]) / 2
            
            # ค้นหาไอเทมอื่นๆ ที่อยู่ในระดับ Y ใกล้เคียงกัน และอยู่ทางขวา (X มากกว่า)
            row_items = []
            for (b, t, p) in results:
                b_center_y = (b[0][1] + b[2][1]) / 2
                # ถ้าระดับ Y ห่างกันไม่เกิน 15 พิกเซล และอยู่ทางขวาของ Keyword
                if abs(center_y - b_center_y) < 15 and b[0][0] > bbox[0][0]:
                    # เน้นเอาเฉพาะก้อนที่มีตัวเลข (ยอดเงิน)
                    if any(char.isdigit() for char in t):
                        row_items.append((b[0][0], t))
            
            # ถ้าเจอหลายตัว เอาตัวที่อยู่ขวาสุด (X มากที่สุด) ซึ่งมักจะเป็นยอดรวมสุทธิ
            if row_items:
                row_items.sort(key=lambda x: x[0], reverse=True)
                return row_items[0][1]
    return "0.00"

ALLOWED_TYPES = {"image/jpeg", "image/png", "application/octet-stream"}

@app.post("/extract")
async def extract_text(file: UploadFile = File(...)):
    if file.content_type not in ALLOWED_TYPES:
        raise HTTPException(status_code=400, detail="Unsupported file type")

    contents = await file.read()
    image = Image.open(io.BytesIO(contents)).convert("RGB")
    
    # --- Image Pre-processing ---
    image = ImageOps.grayscale(image)
    enhancer = ImageEnhance.Contrast(image)
    image = enhancer.enhance(2.0) 
    img_array = np.array(image)

    # 1. อ่านข้อมูลพร้อมพิกัดแบบละเอียด (ไม่ใช้ paragraph เพื่อความแม่นยำของพิกัด)
    result = reader.readtext(img_array)

    # 2. สกัดข้อมูลสำคัญเข้า Summary ตามรูปแบใบแจ้งหนี้
    summary = {
        "rental": find_value_in_row(result, ["ค่าเช่า", "Rental"]),
        "electricity": find_value_in_row(result, ["ค่าไฟฟ้า", "Electricity"]),
        "water": find_value_in_row(result, ["ค่าน้ำ", "Water"]),
        "cable_tv": find_value_in_row(result, ["เคเบิล", "Cable"]),
        "total": find_value_in_row(result, ["รวมทั้งสิ้น", "Total"]),
    }

    # 3. เตรียม Raw Data สำหรับส่งกลับไปโชว์ทั้งหมด (เผื่อต้องการตรวจทาน)
    extracted_data = []
    for (bbox, text, prob) in result:
        extracted_data.append({
            "text": text,
            "confidence": float(prob),
            "box": [[float(p[0]), float(p[1])] for p in bbox]
        })

    return {
        "summary": summary,
        "data": extracted_data
    }

@app.get("/ping")
async def ping():
    return {"message": "OCR Service is ready with Enhanced EasyOCR!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)