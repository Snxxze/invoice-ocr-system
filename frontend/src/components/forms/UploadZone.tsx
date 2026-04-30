"use client";

import React, { useState, useRef } from "react";
import { UploadCloud, FileType, Loader2 } from "lucide-react";
import { invoiceApi } from "../../api/invoice";
import { UploadInvoiceResponse } from "../../types/invoice";

// 1. กำหนด Props (ช่องทางคุยกับ Component แม่)
interface UploadZoneProps {
  onSuccess: (data: UploadInvoiceResponse) => void;
  onError: (error: string) => void;
}

export function UploadZone({ onSuccess, onError }: UploadZoneProps) {
  // 2. State ควบคุม UI
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  
  // ใช้ useRef เพื่ออ้างอิงถึง <input type="file"> ที่ซ่อนอยู่
  const fileInputRef = useRef<HTMLInputElement>(null);

  // 3. ฟังก์ชันจัดการการ Upload จริงๆ
  const handleUpload = async (file: File) => {
    // ป้องกันอัปโหลดไฟล์ซ้ำซ้อนตอนกำลังโหลด
    if (isUploading) return;
    
    // ตรวจสอบประเภทไฟล์คร่าวๆ ฝั่งหน้าบ้าน (ป้องกันการโหลดไฟล์ผิดประเภทให้เสียเวลา)
    if (!file.type.includes("image") && file.type !== "application/pdf") {
      onError("Please upload an image (JPG/PNG) or PDF file.");
      return;
    }

    setIsUploading(true);
    try {
      // เรียก API Service ที่เราสร้างไว้
      const response = await invoiceApi.upload(file);
      onSuccess(response); // แจ้งแม่ว่าทำงานเสร็จแล้ว พร้อมส่งผลลัพธ์
    } catch (err: any) {
      onError(err.response?.data?.error || "Failed to upload invoice.");
    } finally {
      setIsUploading(false); // ปิดสถานะโหลด
    }
  };

  // 4. ฟังก์ชันจัดการ Drag & Drop
  const onDragOver = (e: React.DragEvent) => {
    e.preventDefault(); // บังคับให้เบราว์เซอร์อนุญาตให้ drop ของลงมาได้
    setIsDragging(true);
  };

  const onDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    // ดึงไฟล์ที่ถูกลากมาวาง
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleUpload(files[0]);
    }
  };

  // 5. เรนเดอร์ UI
  return (
    <div
      onClick={() => !isUploading && fileInputRef.current?.click()}
      onDragOver={onDragOver}
      onDragLeave={onDragLeave}
      onDrop={onDrop}
      className={`
        relative flex flex-col items-center justify-center p-12 mt-8
        border-2 border-dashed rounded-xl transition-all duration-200 cursor-pointer
        ${isUploading ? "opacity-50 cursor-not-allowed bg-gray-50 border-gray-300" : ""}
        ${isDragging ? "border-blue-500 bg-blue-50" : "border-gray-300 hover:border-gray-400 hover:bg-gray-50"}
      `}
    >
      {/* Input ไฟล์ที่ถูกซ่อนไว้ แต่จะทำงานเมื่อกดกล่องนี้ */}
      <input
        type="file"
        ref={fileInputRef}
        onChange={(e) => {
          if (e.target.files && e.target.files[0]) {
            handleUpload(e.target.files[0]);
          }
        }}
        accept="image/jpeg, image/png, application/pdf"
        className="hidden"
      />

      {/* ควบคุมการแสดงผลไอคอน (หมุนติ้วๆ ตอนโหลด vs ก้อนเมฆตอนว่าง) */}
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
