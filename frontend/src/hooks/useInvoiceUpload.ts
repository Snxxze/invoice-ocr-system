import { useState } from "react";
import { invoiceApi } from "../api/invoice";
import { UploadInvoiceResponse } from "../types/invoice";

interface UseInvoiceUploadProps {
  onSuccess?: (data: UploadInvoiceResponse) => void;
  onError?: (error: string) => void;
}

export function useInvoiceUpload({ onSuccess, onError }: UseInvoiceUploadProps = {}) {
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const upload = async (file: File) => {
    if (isUploading) return;

    // Client-side validation
    if (!file.type.includes("image") && file.type !== "application/pdf") {
      const msg = "Please upload an image (JPG/PNG) or PDF file.";
      setError(msg);
      onError?.(msg);
      return;
    }

    setIsUploading(true);
    setError(null);

    try {
      const response = await invoiceApi.upload(file);
      onSuccess?.(response);
      return response;
    } catch (err: any) {
      const msg = err.response?.data?.error || "Failed to upload invoice.";
      setError(msg);
      onError?.(msg);
      throw err;
    } finally {
      setIsUploading(false);
    }
  };

  return {
    upload,
    isUploading,
    error,
    clearError: () => setError(null)
  };
}
