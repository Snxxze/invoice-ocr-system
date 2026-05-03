import io
import re
from PIL import Image, ImageOps, ImageEnhance
import numpy as np
from pydantic import BaseModel
from typing import List, Dict

class OCRBox(BaseModel):
    text: str
    confidence: float
    box: List[List[float]]

class OCRSummary(BaseModel):
    rental: str
    electricity: str
    water: str
    cable_tv: str
    total: str

class OCRResponse(BaseModel):
    summary: OCRSummary
    data: List[OCRBox]

class InvoiceProcessor:
    def __init__(self, reader):
        self.reader = reader
        # ในอนาคตสามารถดึง keywords เหล่านี้มาจาก Config หรือ DB ได้
        self.keywords_map = {
            "rental": ["ค่าเช่า", "Rental"],
            "electricity": ["ค่าไฟฟ้า", "Electricity"],
            "water": ["ค่าน้ำ", "Water"],
            "cable_tv": ["เคเบิล", "Cable"],
            "total": ["รวมทั้งสิ้น", "Total"]
        }

    def preprocess_image(self, contents: bytes) -> np.ndarray:
        """ ทำความสะอาดภาพเพื่อให้ OCR อ่านได้ง่ายขึ้น """
        image = Image.open(io.BytesIO(contents)).convert("RGB")
        image = ImageOps.grayscale(image)
        enhancer = ImageEnhance.Contrast(image)
        image = enhancer.enhance(2.0)
        return np.array(image)

    def find_value_in_row(self, results: list, keywords: List[str]) -> str:
        """ ค้นหาตัวเลขที่อยู่ทางขวาของ Keyword ในระดับบรรทัดเดียวกัน """
        for i, (bbox, text, prob) in enumerate(results):
            if any(kw.lower() in text.lower() for kw in keywords):
                center_y = (bbox[0][1] + bbox[2][1]) / 2
                
                row_items = []
                for (b, t, p) in results:
                    b_center_y = (b[0][1] + b[2][1]) / 2
                    if abs(center_y - b_center_y) < 15 and b[0][0] > bbox[0][0]:
                        if any(char.isdigit() for char in t):
                            row_items.append((b[0][0], t))
                
                if row_items:
                    row_items.sort(key=lambda x: x[0], reverse=True)
                    return row_items[0][1]
        return "0.00"

    def process(self, contents: bytes) -> Dict:
        """ รวมกระบวนการทั้งหมดตั้งแต่รับ bytes จนถึงสรุปข้อมูล """
        img_array = self.preprocess_image(contents)
        result = self.reader.readtext(img_array)

        summary = {key: self.find_value_in_row(result, kws) for key, kws in self.keywords_map.items()}
        
        extracted_data = [
            OCRBox(
                text=text,
                confidence=float(prob),
                box=[[float(p[0]), float(p[1])] for p in bbox]
            )
            for (bbox, text, prob) in result
        ]

        return {
            "summary": summary,
            "data": extracted_data
        }
