import { API_BASE_URL, triggerSessionExpired } from "./api";

export interface UploadedFile {
  name: string;
  key: string;
  url: string;
  type: string;
  preview?: string;
}

export interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  files?: UploadedFile[];
  _images?: string[];
  _expectingImage?: boolean;
}

export interface Model {
  id: number;
  name: string;
  display: string;
}

export async function fetchModels(): Promise<Model[]> {
  const response = await fetch(`${API_BASE_URL}/models`);
  const data = await response.json();
  return data.models;
}

export async function sendChatMessage(
  query: string,
  model: string,
  token: string,
  onChunk: (type: "text-delta" | "image" | "tool-prompt", data: string) => void,
  files?: UploadedFile[]
): Promise<{ title: string; response: string }> {
  const body: { query: string; model: string; files?: Array<{ name: string; type: string; url: string }> } = {
    query,
    model,
  };

  if (files && files.length > 0) {
    body.files = files.map((f) => ({
      name: f.name,
      type: f.type,
      url: f.url,
    }));
  }

  const response = await fetch(`${API_BASE_URL}/chat/stream-raw`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      Accept: "text/event-stream",
    },
    body: JSON.stringify(body),
  });

  if (!response.ok) {
    if (response.status === 401) {
      triggerSessionExpired();
    }
    const data = await response.json();
    throw new Error(data.error || "Failed to send message");
  }

  const reader = response.body?.getReader();
  if (!reader) throw new Error("Streaming not supported");

  const decoder = new TextDecoder();
  let buffer = "";
  let title = "";
  let fullResponse = "";

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split("\n");
    buffer = lines.pop() || "";

    let i = 0;
    while (i < lines.length) {
      const line = lines[i];

      if (line.startsWith("event: chunk") && i + 1 < lines.length) {
        const dataLine = lines[i + 1];
        if (dataLine.startsWith("data: ")) {
          const parsed = JSON.parse(dataLine.slice(6));
          const chunkType = parsed.chunk?.type;

          if (chunkType === "text-delta") {
            const delta = parsed.chunk?.delta || "";
            onChunk("text-delta", delta);
          } else if (chunkType === "tool-input-available") {
            const input = parsed.chunk?.input;
            if (input?.prompt) {
              onChunk("tool-prompt", input.prompt);
            }
          } else if (chunkType === "tool-image-google") {
            const output = parsed.chunk?.output;
            if (output?.images?.[0]?.url) {
              onChunk("image", output.images[0].url);
            }
          } else if (chunkType === "output-available") {
            const output = parsed.chunk?.output;
            if (output?.images?.[0]?.url) {
              onChunk("image", output.images[0].url);
            }
          }
        }
        i += 2;
        continue;
      }

      if (line.startsWith("event: done") && i + 1 < lines.length) {
        const dataLine = lines[i + 1];
        if (dataLine.startsWith("data: ")) {
          const parsed = JSON.parse(dataLine.slice(6));
          title = parsed.title || "";
          fullResponse = parsed.response || "";
        }
        i += 2;
        continue;
      }

      if (line.startsWith("event: error") && i + 1 < lines.length) {
        const dataLine = lines[i + 1];
        if (dataLine.startsWith("data: ")) {
          const parsed = JSON.parse(dataLine.slice(6));
          throw new Error(parsed.error || "Stream error");
        }
      }

      i++;
    }
  }

  return { title, response: fullResponse };
}
