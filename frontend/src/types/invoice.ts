export interface Invoice {
  id: string;
  file_url: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface OCRWord {
  text: string;
  confidence: number;
  box: number[][];
}

export interface OCRResult {
  data: OCRWord[];
}

export interface UploadInvoiceResponse {
  invoice: Invoice;
  ocr_result: OCRResult;
}
