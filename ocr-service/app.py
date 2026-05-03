import io
from fastapi import FastAPI, UploadFile, File, HTTPException
import easyocr
from processor import InvoiceProcessor, OCRResponse

app = FastAPI(title="Invoice OCR Service")

# Initialize Reader & Processor
reader = easyocr.Reader(['th', 'en'])
processor = InvoiceProcessor(reader)

ALLOWED_TYPES = {"image/jpeg", "image/png", "application/octet-stream"}

@app.post("/extract", response_model=OCRResponse)
async def extract_text(file: UploadFile = File(...)):
    if file.content_type not in ALLOWED_TYPES:
        raise HTTPException(status_code=400, detail="Unsupported file type")

    contents = await file.read()
    try:
        result = processor.process(contents)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"OCR Processing failed: {str(e)}")

@app.get("/ping")
async def ping():
    return {"status": "ok", "message": "OCR Service is ready with Modular Architecture!"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)