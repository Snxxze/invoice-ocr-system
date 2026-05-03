"use client";

import React, { useState, useRef } from "react";
import { UploadCloud, Loader2 } from "lucide-react";
import { useInvoiceUpload } from "../../hooks/useInvoiceUpload";
import { UploadInvoiceResponse } from "../../types/invoice";

interface UploadZoneProps {
  onSuccess: (data: UploadInvoiceResponse) => void;
  onError: (error: string) => void;
}

export function UploadZone({ onSuccess, onError }: UploadZoneProps) {
  const [isDragging, setIsDragging] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const { upload, isUploading } = useInvoiceUpload({
    onSuccess,
    onError
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      upload(e.target.files[0]);
    }
  };

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files.length > 0) {
      upload(e.dataTransfer.files[0]);
    }
  };

  return (
    <div
      onClick={() => !isUploading && fileInputRef.current?.click()}
      onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
      onDragLeave={() => setIsDragging(false)}
      onDrop={onDrop}
      className={`
        relative flex flex-col items-center justify-center p-12 mt-8
        border-2 border-dashed rounded-xl transition-all duration-200 cursor-pointer
        ${isUploading ? "opacity-50 cursor-not-allowed bg-gray-50 border-gray-300" : ""}
        ${isDragging ? "border-blue-500 bg-blue-50" : "border-gray-300 hover:border-gray-400 hover:bg-gray-50"}
      `}
    >
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        accept="image/jpeg, image/png, application/pdf"
        className="hidden"
      />

      <div className="mb-4 text-gray-500">
        {isUploading ? (
          <Loader2 className="w-12 h-12 animate-spin text-blue-500" />
        ) : (
          <UploadCloud className={`w-12 h-12 ${isDragging ? "text-blue-500" : "text-gray-400"}`} />
        )}
      </div>

      <div className="text-center space-y-2">
        <h3 className="text-lg font-semibold text-gray-800">
          {isUploading ? "Processing Document..." : "Click or drag file to upload"}
        </h3>
        <p className="text-sm text-gray-500">
          {isUploading 
            ? "Extracting data with AI, please wait." 
            : "Supported formats: JPEG, PNG, PDF (Max size: 5MB)"}
        </p>
      </div>
    </div>
  );
}
