import { apiClient } from './client';
import { UploadInvoiceResponse } from '../types/invoice';

export const invoiceApi = {
  /**
   * Upload an invoice file to be processed by the backend and OCR service
   */
  upload: async (file: File): Promise<UploadInvoiceResponse> => {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.post<UploadInvoiceResponse>('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });

    return response.data;
  },
};
