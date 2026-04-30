"use client";

import { useState } from "react";
import { UploadZone } from "../components/forms/UploadZone";
import { ResultView } from "../components/views/ResultView";
import { UploadInvoiceResponse } from "../types/invoice";

export default function Home() {
  const [resultData, setResultData] = useState<UploadInvoiceResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleReset = () => {
    setResultData(null);
    setError(null);
  };

  return (
    <main className="min-h-screen bg-gray-50/50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto space-y-8">
        
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold tracking-tight text-gray-900">
            Smart Invoice OCR
          </h1>
          <p className="text-gray-500">
            Upload your invoices and let AI extract the data instantly.
          </p>
        </div>

        {/* แสดง Error Message ถ้ามีปัญหา */}
        {error && (
          <div className="p-4 bg-red-50 text-red-600 rounded-lg border border-red-100 flex justify-between items-center">
            <p className="text-sm font-medium">{error}</p>
            <button onClick={() => setError(null)} className="text-red-400 hover:text-red-600">
              ✕
            </button>
          </div>
        )}

        {/* สลับหน้าจอระหว่าง Upload กับ Result */}
        {!resultData ? (
          <div className="bg-white p-8 rounded-2xl shadow-sm border border-gray-100">
            <UploadZone 
              onSuccess={(data) => {
                setResultData(data);
                setError(null);
              }} 
              onError={(err) => setError(err)} 
            />
          </div>
        ) : (
          <ResultView 
            data={resultData} 
            onReset={handleReset} 
          />
        )}

      </div>
    </main>
  );
}
