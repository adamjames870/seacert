export interface ApiError {
  error: string;
  request_id?: string;
  status?: number;
}

export const getErrorMessage = async (response: Response): Promise<string> => {
  try {
    const data: ApiError = await response.json();
    let message = data.error || response.statusText;
    if (data.request_id) {
      message += ` (Request ID: ${data.request_id})`;
    }
    return message;
  } catch {
    return `An unexpected error occurred: ${response.statusText}`;
  }
};

export const handleApiError = async (response: Response): Promise<never> => {
  const message = await getErrorMessage(response);
  const error: any = new Error(message);
  error.status = response.status;
  try {
    const data = await response.clone().json();
    error.requestId = data.request_id;
  } catch {
    // ignore
  }
  throw error;
};
