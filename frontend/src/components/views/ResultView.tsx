import React from "react";
import { UploadInvoiceResponse } from "../../types/invoice";
import { CheckCircle2, FileText, RefreshCcw } from "lucide-react";

interface ResultViewProps {
  data: UploadInvoiceResponse;
  onReset: () => void;
}

export function ResultView({ data, onReset }: ResultViewProps) {
  return (
    <div className="w-full max-w-4xl mx-auto mt-8 bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="p-6 border-b border-gray-100 flex items-center justify-between bg-gray-50/50">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-green-100 text-green-600 rounded-full">
            <CheckCircle2 className="w-6 h-6" />
          </div>
          <div>
            <h2 className="text-xl font-semibold text-gray-800">Processing Complete</h2>
            <p className="text-sm text-gray-500">Document ID: {data.invoice.id}</p>
          </div>
        </div>
        <button
          onClick={onReset}
          className="flex items-center gap-2 px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 rounded-lg hover:bg-blue-100 transition-colors"
        >
          <RefreshCcw className="w-4 h-4" />
          Upload Another
        </button>
      </div>

      <div className="p-6 grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* ฝั่งซ้าย: ข้อมูลไฟล์ */}
        <div className="bg-gray-50 p-4 rounded-lg border border-gray-100">
          <div className="flex items-center gap-2 mb-4 text-gray-700 font-medium">
            <FileText className="w-5 h-5 text-gray-400" />
            File Metadata
          </div>
          <dl className="space-y-3 text-sm">
            <div>
              <dt className="text-gray-500">File Name (MinIO)</dt>
              <dd className="font-mono text-gray-800 break-all">{data.invoice.file_url}</dd>
            </div>
            <div>
              <dt className="text-gray-500">Status</dt>
              <dd className="inline-flex px-2 py-1 bg-green-100 text-green-700 rounded-md mt-1">
                {data.invoice.status}
              </dd>
            </div>
          </dl>
        </div>

        {/* ฝั่งขวา: ผลลัพธ์จาก AI */}
        <div className="bg-gray-50 p-4 rounded-lg border border-gray-100 h-[400px] flex flex-col">
          <div className="flex items-center gap-2 mb-4 text-gray-700 font-medium">
            <h3 className="flex items-center gap-2">Extracted Data</h3>
            <span className="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded-full">
              {data.ocr_result.data.length} items found
            </span>
          </div>

          <div className="flex-1 overflow-y-auto pr-2 space-y-2">
            {data.ocr_result.data.map((item, index) => (
              <div key={index} className="bg-white p-3 rounded border border-gray-100 shadow-sm flex justify-between items-start">
                <span className="text-gray-800 font-medium">{item.text}</span>
                <span className="text-xs text-gray-400 font-mono">
                  {(item.confidence * 100).toFixed(1)}%
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
