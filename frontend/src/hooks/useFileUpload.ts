import { useState } from 'react';
import { supabase } from '../supabaseClient';
import { API_BASE_URL } from '../config';

interface FileUploadOptions {
  maxSizeMB?: number;
  allowedTypes?: string[];
}

export const useFileUpload = (options: FileUploadOptions = {}) => {
  const { 
    maxSizeMB = 3, 
    allowedTypes = ['application/pdf', 'image/jpeg', 'image/jpg'] 
  } = options;

  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [error, setError] = useState<string | null>(null);

  const uploadFile = async (file: File) => {
    // 1. Validate file size
    const MAX_FILE_SIZE = maxSizeMB * 1024 * 1024;
    if (file.size > MAX_FILE_SIZE) {
      const err = `File is too large. Maximum size is ${maxSizeMB}MB.`;
      setError(err);
      throw new Error(err);
    }

    // 2. Validate file type
    if (!allowedTypes.includes(file.type)) {
      const err = `Only ${allowedTypes.join(', ')} files are allowed`;
      setError(err);
      throw new Error(err);
    }

    setUploading(true);
    setProgress(0);
    setError(null);

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      // 3. Get upload URL
      const urlResponse = await fetch(`${API_BASE_URL}/api/certificates/upload-url`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({
          filename: file.name,
          'content-type': file.type
        }),
      });

      if (!urlResponse.ok) {
        throw new Error('Failed to get upload URL');
      }

      const { 'upload-url': uploadUrl, 'file-key': fileKey } = await urlResponse.json();

      // 4. Upload to R2 with progress tracking using XMLHttpRequest
      const fileKeyResult = await new Promise<{ fileKey: string, fileName: string }>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('PUT', uploadUrl, true);
        xhr.setRequestHeader('Content-Type', file.type);

        xhr.upload.onprogress = (progressEvent) => {
          if (progressEvent.lengthComputable) {
            const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
            setProgress(percentCompleted);
          }
        };

        xhr.onload = () => {
          if (xhr.status >= 200 && xhr.status < 300) {
            resolve({ fileKey, fileName: file.name });
          } else {
            reject(new Error(`Upload failed with status ${xhr.status}`));
          }
        };

        xhr.onerror = () => {
          reject(new Error('Network error during upload'));
        };

        xhr.send(file);
      });

      return fileKeyResult;
    } catch (err: any) {
      const errorMessage = err.message || 'Failed to upload document';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setUploading(false);
    }
  };

  return { uploadFile, uploading, progress, error, setError };
};
