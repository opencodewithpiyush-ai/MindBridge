import { useState, useEffect, useRef, type FormEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";
import { fetchModels, sendChatMessage } from "../../lib/chat";
import type { Message, Model, UploadedFile } from "../../lib/chat";
import { uploadFile } from "../../lib/api";
import styles from "./Chat.module.css";

function ImageWithShimmer({ url }: { url: string }) {
  const [loaded, setLoaded] = useState(false);

  return (
    <div className={styles["image-container"]}>
      {!loaded && <div className={styles["image-shimmer"]} />}
      <img
        src={url}
        alt="Generated image"
        className={styles["message-image"]}
        onLoad={() => setLoaded(true)}
      />
    </div>
  );
}

function Chat() {
  const { user, token, logout, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [models, setModels] = useState<Model[]>([]);
  const [selectedModel, setSelectedModel] = useState("");
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [error, setError] = useState("");
  const [showModelDropdown, setShowModelDropdown] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState<UploadedFile[]>([]);
  const [uploading, setUploading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLTextAreaElement>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    fetchModels().then((m) => {
      setModels(m);
      if (m.length > 0 && !selectedModel) {
        setSelectedModel(m.find((m) => m.name.includes("opus"))?.name || m[0].name);
      }
    });
  }, []);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  if (!isAuthenticated) {
    navigate("/login");
    return null;
  }

  const handleFileUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file || !token) return;

    setUploading(true);

    try {
      const result = await uploadFile(file, token);
      setUploadedFiles((prev) => [
        ...prev,
        { name: result.name, key: result.key, url: result.url, type: result.type, preview: result.preview },
      ]);
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Upload failed";
      setError(msg);
    } finally {
      setUploading(false);
      if (fileInputRef.current) {
        fileInputRef.current.value = "";
      }
    }
  };

  const removeUploadedFile = (index: number) => {
    setUploadedFiles((prev) => prev.filter((_, i) => i !== index));
  };

  const handleSend = async (e: FormEvent) => {
    e.preventDefault();
    if ((!input.trim() && uploadedFiles.length === 0) || loading) return;

    setError("");
    const userMessage: Message = {
      id: Date.now().toString(),
      role: "user",
      content: input.trim() || "See the attached file",
      files: uploadedFiles.length > 0 ? [...uploadedFiles] : undefined,
    };

    setMessages((prev) => [...prev, userMessage]);
    const query = input.trim();
    setInput("");
    setLoading(true);
    const currentFiles = [...uploadedFiles];
    setUploadedFiles([]);

      const assistantId = (Date.now() + 1).toString();
      setMessages((prev) => [
        ...prev,
        { id: assistantId, role: "assistant", content: "", _images: [], _expectingImage: false },
      ]);
      let accumulatedText = "";
      let images: string[] = [];
      let expectingImage = currentFiles.some((f) => f.type.startsWith("image/"));

      try {
        if (!token) return;
        await sendChatMessage(query, selectedModel, token, (chunkType, data) => {
          if (chunkType === "text-delta") {
            accumulatedText += data;
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantId ? { ...m, content: accumulatedText } : m
              )
            );
          } else if (chunkType === "image") {
            images = [...images, data];
            expectingImage = true;
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantId ? { ...m, _images: images, _expectingImage: expectingImage } : m
              )
            );
          }
        }, currentFiles.length > 0 ? currentFiles : undefined);

        setMessages((prev) =>
          prev.map((m) =>
            m.id === assistantId ? { ...m, _expectingImage: expectingImage } : m
          )
        );
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : "Failed to get response";
      setError(msg);
      setMessages((prev) =>
        prev.map((m) =>
          m.id === assistantId
            ? { ...m, content: `Error: ${msg}` }
            : m
        )
      );
    } finally {
      setLoading(false);
      inputRef.current?.focus();
    }
  };

  const handleNewChat = () => {
    setMessages([]);
    setError("");
    setUploadedFiles([]);
    inputRef.current?.focus();
  };

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  const selectedModelDisplay =
    models.find((m) => m.name === selectedModel)?.display || selectedModel;

  return (
    <div className={styles.chat}>
      {sidebarOpen && (
        <aside className={styles.sidebar}>
          <div className={styles["sidebar-header"]}>
            <div className={styles.logo}>
              <div className={styles["logo-icon"]}>M</div>
              <span className={styles["logo-text"]}>MindBridge</span>
            </div>
            <button
              onClick={handleNewChat}
              className={`${styles.btn} ${styles["btn-new"]}`}
            >
              + New Chat
            </button>
          </div>

          <div className={styles["sidebar-content"]}>
            {messages.length > 0 && (
              <div className={styles["current-chat"]}>
                <span className={styles["sidebar-label"]}>Current Chat</span>
                <button className={styles["chat-item"]}>
                  {messages[0]?.content.slice(0, 30)}
                  {messages[0].content.length > 30 ? "..." : ""}
                </button>
              </div>
            )}
          </div>

          <div className={styles["sidebar-footer"]}>
            <div className={styles["user-info"]}>
              <div className={styles["user-avatar"]}>
                {user?.name?.charAt(0).toUpperCase()}
              </div>
              <div className={styles["user-details"]}>
                <span className={styles["user-name"]}>{user?.name}</span>
                <span className={styles["user-email"]}>{user?.email}</span>
              </div>
            </div>
            <button onClick={handleLogout} className={`${styles.btn} ${styles["btn-logout"]}`}>
              Logout
            </button>
          </div>
        </aside>
      )}

      <main className={styles.main}>
        <header className={styles.header}>
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className={styles["sidebar-toggle"]}
          >
            {sidebarOpen ? "◀" : "▶"}
          </button>

          <div className={styles["model-selector"]}>
            <button
              className={styles["model-button"]}
              onClick={() => setShowModelDropdown(!showModelDropdown)}
            >
              {selectedModelDisplay}
              <span className={styles.chevron}>▼</span>
            </button>
            {showModelDropdown && (
              <div className={styles["model-dropdown"]}>
                {models.map((model) => (
                  <button
                    key={model.id}
                    className={`${styles["model-option"]} ${
                      model.name === selectedModel ? styles.active : ""
                    }`}
                    onClick={() => {
                      setSelectedModel(model.name);
                      setShowModelDropdown(false);
                    }}
                  >
                    {model.display}
                  </button>
                ))}
              </div>
            )}
          </div>

          <div className={styles["header-spacer"]} />
        </header>

        <div className={styles["messages-container"]}>
          {messages.length === 0 ? (
            <div className={styles.empty}>
              <h2 className={styles["empty-title"]}>
                What can I help with?
              </h2>
              <p className={styles["empty-description"]}>
                Start a conversation with any AI model.
              </p>
            </div>
          ) : (
            <div className={styles.messages}>
              {messages.map((msg) => {
                const content = msg.content || "";
                const imageUrls = msg._images || [];
                const expectingImage = msg._expectingImage || false;

                return (
                  <div
                    key={msg.id}
                    className={`${styles.message} ${styles[msg.role]}`}
                  >
                    <div className={styles["message-content"]}>
                      {msg.files && msg.files.length > 0 && (
                        <div className={styles["message-files"]}>
                           {msg.files.map((f, i) =>
                             f.type.startsWith("image/") ? (
                               <img
                                 key={i}
                                 src={f.preview || f.url}
                                 alt={f.name}
                                 className={styles["message-image"]}
                               />
                             ) : (
                              <div key={i} className={styles["message-file"]}>
                                📎 {f.name}
                              </div>
                            )
                          )}
                        </div>
                      )}
                      {imageUrls.length > 0 && (
                        <div className={styles["message-images"]}>
                          {imageUrls.map((url, i) => (
                            <ImageWithShimmer key={i} url={url} />
                          ))}
                        </div>
                      )}
                      {expectingImage && imageUrls.length === 0 && content === "" && (
                        <div className={styles.shimmer} />
                      )}
                      {content && <div className={styles["message-text"]}>{content}</div>}
                      {!content && !expectingImage && (
                        <span className={styles.typing}>Thinking...</span>
                      )}
                    </div>
                  </div>
                );
              })}
              <div ref={messagesEndRef} />
            </div>
          )}

          {error && <div className={styles["error-bar"]}>{error}</div>}
        </div>

        <form onSubmit={handleSend} className={styles["input-area"]}>
          {uploadedFiles.length > 0 && (
            <div className={styles["uploaded-files"]}>
              {uploadedFiles.map((file, index) => (
                <div key={index} className={styles["uploaded-file"]}>
                  {file.type.startsWith("image/") ? (
                    <img src={file.preview || file.url} alt={file.name} className={styles["file-preview"]} />
                  ) : (
                    <div className={styles["file-icon"]}>📎</div>
                  )}
                  <span className={styles["file-name"]}>{file.name}</span>
                  <button
                    type="button"
                    onClick={() => removeUploadedFile(index)}
                    className={styles["file-remove"]}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}

          <div className={styles["input-wrapper"]}>
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              className={styles["attach-button"]}
              disabled={uploading}
            >
              {uploading ? "..." : "+"}
            </button>
            <input
              ref={fileInputRef}
              type="file"
              onChange={handleFileUpload}
              className={styles["file-input"]}
              hidden
            />
            <textarea
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  handleSend(e);
                }
              }}
              className={styles.input}
              placeholder="Type your message..."
              rows={1}
              disabled={loading}
            />
            <button
              type="submit"
              className={styles["send-button"]}
              disabled={!input.trim() || loading}
            >
              {loading ? "..." : "↑"}
            </button>
          </div>
          <p className={styles["input-footer"]}>
            {selectedModelDisplay} · Press Enter to send, Shift+Enter for new line
          </p>
        </form>
      </main>
    </div>
  );
}

export default Chat;
